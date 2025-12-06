package main

import (
	"log"
	"os"

	"github.com/skyapps-id/casdoor-test/migrations"
)

func main() {
	certBytes, err := os.ReadFile("./cert.pem")
	if err == nil {
		log.Println("📄 Certificate loaded from cert.pem")
	}
	// Load konfigurasi dari environment variables atau config file
	config := migrations.CasdoorConfig{
		Endpoint:         getEnv("CASDOOR_ENDPOINT", "http://localhost:8000"),
		ClientID:         getEnv("CASDOOR_CLIENT_ID", "26054a2dfb593fa0990c"),
		ClientSecret:     getEnv("CASDOOR_CLIENT_SECRET", "520986d9f34ddcee617d9dfe6c4e203c97c5bc04"),
		Certificate:      getEnv("CASDOOR_CERTIFICATE", string(certBytes)), // Path ke certificate file atau certificate string
		OrganizationName: getEnv("CASDOOR_ORGANIZATION", "skyapps"),
		ApplicationName:  getEnv("CASDOOR_APPLICATION", "application_i9irbv"),
	}

	// Buat instance migration
	migration, err := migrations.NewCasdoorMigration(config)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	// // Jalankan migration
	if err := migration.Run(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")

	// Untuk rollback (uncomment jika diperlukan)
	// if err := migration.Rollback(); err != nil {
	// 	log.Fatalf("Rollback failed: %v", err)
	// }
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
