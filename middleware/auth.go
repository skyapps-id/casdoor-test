package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/skyapps-id/casdoor-test/config"
)

func AuthRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Ambil token dari header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			// Extract token (format: Bearer <token>)
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format",
				})
			}

			token := parts[1]

			// Verifikasi token dengan Casdoor
			claims, err := config.CasdoorClient.ParseJwtToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token " + err.Error(),
				})
			}

			// Validasi organization
			if claims.User.Owner != "skyapps" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Invalid organization",
				})
			}

			// Simpan user info ke context
			c.Set("user_id", claims.User.Name)
			c.Set("user_email", claims.User.Email)
			c.Set("user_org", claims.User.Owner)
			c.Set("user_roles", claims.User.Roles)

			return next(c)
		}
	}
}

// Middleware untuk enforce permission menggunakan Casbin
func EnforcePermission(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Get("user_id").(string)

			// Check permission dengan Casbin
			allowed, err := config.CheckPermission(username, resource, action)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to check permission",
				})
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error":    "Insufficient permissions",
					"required": resource + ":" + action,
				})
			}

			return next(c)
		}
	}
}
