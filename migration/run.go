package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/skyapps-id/casdoor-test/migrations"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using OS env")
	}

	certBytes, err := os.ReadFile("./cert.pem")
	if err == nil {
		log.Println("📄 Certificate loaded from cert.pem")
	}
	// Load konfigurasi dari environment variables atau config file
	config := migrations.CasdoorConfig{
		Endpoint:         getEnv("CASDOOR_ENDPOINT", ""),
		ClientID:         getEnv("CASDOOR_CLIENT_ID", ""),
		ClientSecret:     getEnv("CASDOOR_CLIENT_SECRET", ""),
		Certificate:      getEnv("CASDOOR_CERTIFICATE", string(certBytes)),
		OrganizationName: getEnv("CASDOOR_ORGANIZATION", ""),
		ApplicationName:  getEnv("CASDOOR_APPLICATION", ""),
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
