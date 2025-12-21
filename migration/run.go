package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	rbac "github.com/skyapps-id/casdoor-test/migration/module"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env")
	}

	// Load cert (optional)
	var cert string
	if certBytes, err := os.ReadFile("./cert.pem"); err == nil {
		log.Println("📄 Certificate loaded from cert.pem")
		cert = string(certBytes)
	}

	// Config Casdoor
	config := rbac.CasdoorConfig{
		Endpoint:         getEnv("CASDOOR_ENDPOINT", ""),
		ClientID:         getEnv("CASDOOR_CLIENT_ID", ""),
		ClientSecret:     getEnv("CASDOOR_CLIENT_SECRET", ""),
		Certificate:      getEnv("CASDOOR_CERTIFICATE", cert),
		OrganizationName: getEnv("CASDOOR_ORGANIZATION", ""),
		ApplicationName:  getEnv("CASDOOR_APPLICATION", ""),
	}

	migration, err := rbac.NewCasdoorMigration(config)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run run.go [up|down]")
	}
	switch os.Args[1] {
	case "up":
		log.Println("🚀 Running migration UP...")
		if err := migration.Run(); err != nil {
			log.Fatalf("Migration UP failed: %v", err)
		}
		log.Println("✅ Migration UP completed")

	case "down":
		log.Println("↩️ Running migration DOWN (rollback)...")
		if err := migration.Rollback(); err != nil {
			log.Fatalf("Migration DOWN failed: %v", err)
		}
		log.Println("✅ Rollback completed")

	default:
		log.Fatalf("Unknown command: %s (use up or down)", os.Args[1])
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
