package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/skyapps-id/casdoor-test/config"
	"github.com/skyapps-id/casdoor-test/handlers"
	"github.com/skyapps-id/casdoor-test/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using OS env")
	}

	// Initialize Casdoor
	config.InitCasdoor()

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
	api := e.Group("/api",
		middleware.CasdoorAuthRequired(),
		middleware.CasdoorRBAC(),
	)
	{
		// User info
		api.GET("/me", handlers.GetCurrentUser)

		// User management (requires permission)
		api.GET("/users", handlers.ListUsers)
		api.POST("/users", handlers.AddUser)
		api.PUT("/users/:username", handlers.UpdateUser)
		api.DELETE("/users/:username", handlers.DeleteUser)

		// Role management (admin only)
		api.GET("/roles", handlers.ListRoles)
		api.POST("/roles", handlers.AddRole)
		api.PUT("/roles/:role", handlers.UpdateRole)
		api.DELETE("/roles/:role", handlers.DeleteRole)

		// Assign role to user
		api.POST("/users/:username/roles", handlers.AssignRole)
		api.DELETE("/users/:username/roles/:role", handlers.RemoveRole)

		// RBAC sync
		api.POST("/rbac/sync", handlers.SyncRBAC)
	}

	e.Logger.Fatal(e.Start(":9000"))
}
