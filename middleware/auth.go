package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/labstack/echo/v4"
	"github.com/skyapps-id/casdoor-test/config"
)

func CasdoorAuthRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return echo.NewHTTPError(401, "Missing Bearer token")
			}
			token := auth[7:]

			// Parse token via initialized client
			claims, err := config.CasdoorClient.ParseJwtToken(token)
			if err != nil {
				return echo.NewHTTPError(401, "Invalid token")
			}

			user := claims.User
			if user.Name == "" {
				return echo.NewHTTPError(401, "User is nil in token")
			}

			// Optionally refresh full user info from Casdoor
			fullUser, err := config.CasdoorClient.GetUser(user.Name)
			if err != nil {
				return echo.NewHTTPError(401, "User not found at Casdoor")
			}

			c.Set("casdoorUser", fullUser)
			return next(c)
		}
	}
}

// Middleware untuk enforce permission menggunakan Casbin
func CasdoorRBAC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("casdoorUser").(*casdoorsdk.User)
			if !ok {
				return echo.NewHTTPError(401, "Unauthorized")
			}

			// Ambil dari context
			resource := c.Request().URL.Path
			// action := c.Request().Method

			// Normalize path dengan wildcard
			pathParts := strings.Split(strings.Trim(resource, "/"), "/")
			if len(pathParts) > 2 {
				lastPart := pathParts[len(pathParts)-1]
				if _, err := strconv.Atoi(lastPart); err == nil {
					pathParts[len(pathParts)-1] = "*"
					resource = "/" + strings.Join(pathParts, "/")
				}
			}

			req := casdoorsdk.CasbinRequest{
				"skyapps", "admin", "GET", "/api/users", "skyapps", "*",
			}

			fmt.Println(user.Owner)

			allowed, err := config.CasdoorClient.Enforce(
				"",                          // permissionId
				"",                          // modelId
				"",                          // resourceId
				user.Owner+"/rbac-enforcer", // enforcerId
				"",                          // owner - kosongkan karena sudah ada di permissionId
				req,
			)
			if err != nil {
				return echo.NewHTTPError(500, "RBAC enforcement failed: "+err.Error())
			}
			if allowed {
				return echo.NewHTTPError(403, "Forbidden")
			}

			return next(c)
		}
	}
}
