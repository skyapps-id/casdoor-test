package main

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/skyapps-id/casdoor-test/config"
	"github.com/skyapps-id/casdoor-test/handlers"
	"github.com/skyapps-id/casdoor-test/middleware"
)

func main() {
	// Initialize Casdoor
	config.InitCasdoor()

	// Initialize Casbin
	if err := config.InitCasbin(); err != nil {
		panic("Failed to initialize Casbin: " + err.Error())
	}

	// Sync RBAC from Casdoor
	if err := config.SyncRBACFromCasdoor(); err != nil {
		panic("Failed to sync RBAC: " + err.Error())
	}

	// Setup Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Public routes
	e.GET("/login", handlers.GetLoginURL)
	e.GET("/callback", handlers.HandleCallback)
	e.GET("/health", handlers.HealthCheck)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// User info
		api.GET("/me", handlers.GetCurrentUser)

		// User management (requires permission)
		api.GET("/users", handlers.ListUsers, middleware.EnforcePermission("users", "read"))
		api.POST("/users", handlers.AddUser, middleware.EnforcePermission("users", "write"))
		api.PUT("/users/:username", handlers.UpdateUser, middleware.EnforcePermission("users", "write"))
		api.DELETE("/users/:username", handlers.DeleteUser, middleware.EnforcePermission("users", "delete"))

		// Role management (admin only)
		api.GET("/roles", handlers.ListRoles, middleware.EnforcePermission("roles", "read"))
		api.POST("/roles", handlers.AddRole, middleware.EnforcePermission("roles", "write"))
		api.PUT("/roles/:role", handlers.UpdateRole, middleware.EnforcePermission("roles", "write"))
		api.DELETE("/roles/:role", handlers.DeleteRole, middleware.EnforcePermission("roles", "delete"))

		// Assign role to user
		api.POST("/users/:username/roles", handlers.AssignRole, middleware.EnforcePermission("users", "write"))
		api.DELETE("/users/:username/roles/:role", handlers.RemoveRole, middleware.EnforcePermission("users", "write"))

		// RBAC sync
		api.POST("/rbac/sync", handlers.SyncRBAC, middleware.EnforcePermission("rbac", "write"))
	}

	e.Logger.Fatal(e.Start(":9000"))
}
