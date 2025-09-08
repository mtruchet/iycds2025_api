package configs

import (
	"fmt"
	"os"

	"iycds2025_api/src/api/services/mail"
)

// EmailConfig contiene la configuración para los servicios de email
type EmailConfig struct {
	ServiceType string
	// SendGrid
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string
	// SMTP
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string
}

// GetEmailConfig retorna la configuración de email según el entorno
func GetEmailConfig() EmailConfig {
	appEnv := os.Getenv("APP_ENV")
	emailServiceType := os.Getenv("EMAIL_SERVICE_TYPE")

	// Configuración por defecto según el entorno
	if emailServiceType == "" {
		if appEnv == "production" {
			emailServiceType = "sendgrid" // Por defecto SendGrid en producción
		} else {
			emailServiceType = "mock" // Por defecto Mock en desarrollo
		}
	}

	config := EmailConfig{
		ServiceType: emailServiceType,
		// SendGrid
		SendGridAPIKey:    os.Getenv("SENDGRID_API_KEY"),
		SendGridFromEmail: os.Getenv("SENDGRID_FROM_EMAIL"),
		SendGridFromName:  os.Getenv("SENDGRID_FROM_NAME"),
		// SMTP
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      587, // Puerto por defecto
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		SMTPFromName:  os.Getenv("SMTP_FROM_NAME"),
	}

	return config
}

// NewEmailService crea el servicio de email apropiado según la configuración
func NewEmailService() mail.EmailService {
	config := GetEmailConfig()

	fmt.Printf("Configurando servicio de email: %s\n", config.ServiceType)

	switch config.ServiceType {
	case "sendgrid":
		fmt.Println("Inicializando SendGrid Email Service")
		return mail.NewSendGridEmailService()
	case "smtp":
		fmt.Printf("Inicializando SMTP Email Service (Host: %s)\n", config.SMTPHost)
		return mail.NewSMTPEmailService()
	case "mock":
		fmt.Println("Inicializando Mock Email Service (solo logs)")
		return mail.NewMockEmailService()
	default:
		fmt.Printf("Tipo de servicio de email desconocido '%s', usando mock por defecto\n", config.ServiceType)
		return mail.NewMockEmailService()
	}
}
