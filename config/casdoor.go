package config

import (
	"log"
	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var CasdoorClient *casdoorsdk.Client

func InitCasdoor() {
	endpoint := "http://localhost:8000"
	clientId := "26054a2dfb593fa0990c"
	clientSecret := "520986d9f34ddcee617d9dfe6c4e203c97c5bc04"
	appName := "application_i9irbv"
	orgName := "skyapps_6gurhn"
	redirectURL := "http://localhost:9000/callback"

	// Load certificate
	certificate := loadCertificate()

	if certificate != "" {
		// Initialize dengan certificate untuk JWT verification
		CasdoorClient = casdoorsdk.NewClientWithConf(&casdoorsdk.AuthConfig{
			Endpoint:         endpoint,
			ClientId:         clientId,
			ClientSecret:     clientSecret,
			Certificate:      certificate,
			OrganizationName: orgName,
			ApplicationName:  appName,
		})
		log.Println("✅ Casdoor client initialized with certificate")
	} else {
		// Fallback ke basic client (tanpa JWT verification)
		CasdoorClient = casdoorsdk.NewClient(
			endpoint,
			clientId,
			clientSecret,
			appName,
			orgName,
			redirectURL,
		)
		log.Println("⚠️  Casdoor client initialized without certificate")
	}
}

// loadCertificate mencoba load certificate dengan berbagai cara
func loadCertificate() string {
	// Priority 1: Load dari file cert.pem
	if certBytes, err := os.ReadFile("./cert.pem"); err == nil {
		log.Println("📄 Certificate loaded from cert.pem")
		return string(certBytes)
	}

	// Priority 2: Load dari environment variable
	if cert := os.Getenv("CASDOOR_CERTIFICATE"); cert != "" {
		log.Println("📄 Certificate loaded from environment")
		return cert
	}

	return ""
}
