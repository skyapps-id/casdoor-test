package middleware

import (
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
			// 1️⃣ Ambil user dari context
			user, ok := c.Get("casdoorUser").(*casdoorsdk.User)
			if !ok || user == nil {
				return echo.NewHTTPError(401, "Unauthorized")
			}

			// 2️⃣ Ambil role user
			if len(user.Roles) == 0 {
				return echo.NewHTTPError(403, "No role assigned")
			}

			// (asumsi single role, kalau multi role → loop)
			role := user.Roles[0].Name

			// 3️⃣ Ambil request info
			action := c.Request().Method
			resource := c.Path() // PENTING: pakai path echo, bukan raw URL

			// 4️⃣ Normalize path (id → *)
			pathParts := strings.Split(strings.Trim(resource, "/"), "/")
			if len(pathParts) > 2 {
				lastPart := pathParts[len(pathParts)-1]
				if _, err := strconv.Atoi(lastPart); err == nil {
					pathParts[len(pathParts)-1] = "*"
					resource = "/" + strings.Join(pathParts, "/")
				}
			}

			// 5️⃣ Build Casbin request
			req := casdoorsdk.CasbinRequest{
				user.Owner, // subOwner
				role,       // subName (ROLE)
				action,     // method
				resource,   // path
				user.Owner, // objOwner
				"*",        // objName
			}

			// 6️⃣ Enforce RBAC
			allowed, err := config.CasdoorClient.Enforce(
				"",
				"",
				"",
				user.Owner+"/rbac-enforcer",
				"",
				req,
			)
			if err != nil {
				return echo.NewHTTPError(500, "RBAC enforcement failed: "+err.Error())
			}

			// 7️⃣ Deny kalau tidak allowed
			if !allowed {
				return echo.NewHTTPError(403, "Forbidden")
			}

			return next(c)
		}
	}
}
