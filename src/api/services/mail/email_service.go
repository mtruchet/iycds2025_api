package mail

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	gomail "gopkg.in/mail.v2"
)

// EmailService define la interfaz para enviar correos electrónicos
type EmailService interface {
	SendPasswordResetEmail(to string, resetLink string) error
}

// SMTPEmailService implementa el servicio de correo usando SMTP
type SMTPEmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// NewSMTPEmailService crea una nueva instancia del servicio de correo
func NewSMTPEmailService() *SMTPEmailService {
	port := 587 // Puerto SMTP estándar para TLS

	// Validar que las variables SMTP necesarias estén configuradas
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpFromEmail := os.Getenv("SMTP_FROM_EMAIL")
	smtpFromName := os.Getenv("SMTP_FROM_NAME")

	// Verificar variables críticas
	if smtpHost == "" {
		fmt.Println("WARNING: SMTP_HOST no está configurado, usando valor por defecto")
		smtpHost = "smtp.gmail.com"
	}
	if smtpUsername == "" {
		fmt.Println("WARNING: SMTP_USERNAME no está configurado")
		smtpUsername = "test@gmail.com"
	}
	if smtpPassword == "" {
		fmt.Println("WARNING: SMTP_PASSWORD no está configurado")
		smtpPassword = "test-password"
	}
	if smtpFromEmail == "" {
		fmt.Println("WARNING: SMTP_FROM_EMAIL no está configurado, usando valor por defecto")
		smtpFromEmail = "noreply@iycds2025.com"
	}
	if smtpFromName == "" {
		fmt.Println("WARNING: SMTP_FROM_NAME no está configurado, usando valor por defecto")
		smtpFromName = "IYCDS 2025"
	}

	return &SMTPEmailService{
		Host:     smtpHost,
		Port:     port,
		Username: smtpUsername,
		Password: smtpPassword,
		From:     smtpFromEmail,
		FromName: smtpFromName,
	}
}

// SendPasswordResetEmail envía un correo de restablecimiento de contraseña
func (s *SMTPEmailService) SendPasswordResetEmail(to string, resetLink string) error {
	m := gomail.NewMessage()

	// Configurar remitente con formato "Nombre <email@dominio.com>"
	if s.FromName != "" {
		m.SetHeader("From", fmt.Sprintf("%s <%s>", s.FromName, s.From))
	} else {
		m.SetHeader("From", s.From)
	}

	m.SetHeader("To", to)
	m.SetHeader("Subject", "Restablecimiento de contraseña - IYCDS 2025")

	// Cuerpo HTML del correo
	body := fmt.Sprintf(`
		<html>
			<head>
				<style>
					body {
						font-family: Arial, sans-serif;
						line-height: 1.6;
						color: #333;
					}
					.container {
						max-width: 600px;
						margin: 0 auto;
						padding: 20px;
						border: 1px solid #ddd;
						border-radius: 5px;
					}
					h2 {
						color: #2c5282;
					}
					.button {
						display: inline-block;
						background-color: #2c5282;
						color: white;
						padding: 12px 24px;
						text-decoration: none;
						border-radius: 4px;
						margin: 15px 0;
					}
					.footer {
						margin-top: 30px;
						font-size: 12px;
						color: #666;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h2>Restablecimiento de contraseña</h2>
					<p>Has solicitado restablecer tu contraseña en IYCDS 2025.</p>
					<p>Haz clic en el siguiente enlace para crear una nueva contraseña:</p>
					<p>
						<a href="%s" class="button">Restablecer contraseña</a>
					</p>
					<p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
					<p>Este enlace expirará en 1 hora por motivos de seguridad.</p>
					<div class="footer">
						<p>Este correo fue enviado automáticamente. Por favor, no respondas a este mensaje.</p>
					</div>
				</div>
			</body>
		</html>
	`, resetLink)

	m.SetBody("text/html", body)

	// Configurar dialer SMTP
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// Enviar el correo
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error al enviar correo: %v", err)
	}

	fmt.Printf("Correo de restablecimiento enviado a: %s\n", to)
	return nil
}

// MockEmailService implementa el servicio de correo para testing
type MockEmailService struct{}

func (m *MockEmailService) SendPasswordResetEmail(email, resetLink string) error {
	// Por ahora simulamos el envío de email
	fmt.Printf("Mock: Sending password reset email to %s with link %s\n", email, resetLink)
	return nil
}

func NewMockEmailService() EmailService {
	return &MockEmailService{}
}

// SendGridEmailService implementa el servicio de correo usando SendGrid
type SendGridEmailService struct {
	APIKey    string
	FromEmail string
	FromName  string
}

// NewSendGridEmailService crea una nueva instancia del servicio de SendGrid
func NewSendGridEmailService() *SendGridEmailService {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	// Validar configuración
	if apiKey == "" {
		fmt.Println("WARNING: SENDGRID_API_KEY no está configurado")
		apiKey = "test-api-key"
	}
	if fromEmail == "" {
		fmt.Println("WARNING: SENDGRID_FROM_EMAIL no está configurado, usando valor por defecto")
		fromEmail = "noreply@iycds2025.com"
	}
	if fromName == "" {
		fmt.Println("WARNING: SENDGRID_FROM_NAME no está configurado, usando valor por defecto")
		fromName = "IYCDS 2025"
	}

	return &SendGridEmailService{
		APIKey:    apiKey,
		FromEmail: fromEmail,
		FromName:  fromName,
	}
}

// SendPasswordResetEmail envía un correo de restablecimiento de contraseña usando SendGrid
func (s *SendGridEmailService) SendPasswordResetEmail(to string, resetLink string) error {
	from := mail.NewEmail(s.FromName, s.FromEmail)
	toEmail := mail.NewEmail("", to)
	subject := "Restablecimiento de contraseña - IYCDS 2025"

	// Cuerpo HTML del correo
	htmlContent := fmt.Sprintf(`
		<html>
			<head>
				<style>
					body {
						font-family: Arial, sans-serif;
						line-height: 1.6;
						color: #333;
					}
					.container {
						max-width: 600px;
						margin: 0 auto;
						padding: 20px;
						border: 1px solid #ddd;
						border-radius: 5px;
					}
					h2 {
						color: #2c5282;
					}
					.button {
						display: inline-block;
						background-color: #2c5282;
						color: white;
						padding: 12px 24px;
						text-decoration: none;
						border-radius: 4px;
						margin: 15px 0;
					}
					.footer {
						margin-top: 30px;
						font-size: 12px;
						color: #666;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h2>Restablecimiento de contraseña</h2>
					<p>Has solicitado restablecer tu contraseña en IYCDS 2025.</p>
					<p>Haz clic en el siguiente enlace para crear una nueva contraseña:</p>
					<p>
						<a href="%s" class="button">Restablecer contraseña</a>
					</p>
					<p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
					<p>Este enlace expirará en 1 hora por motivos de seguridad.</p>
					<div class="footer">
						<p>Este correo fue enviado automáticamente. Por favor, no respondas a este mensaje.</p>
					</div>
				</div>
			</body>
		</html>
	`, resetLink)

	// Texto plano como alternativa
	plainTextContent := fmt.Sprintf(`
Restablecimiento de contraseña - IYCDS 2025

Has solicitado restablecer tu contraseña en IYCDS 2025.

Visita el siguiente enlace para crear una nueva contraseña:
%s

Si no solicitaste este cambio, puedes ignorar este correo.
Este enlace expirará en 1 hora por motivos de seguridad.

Este correo fue enviado automáticamente. Por favor, no respondas a este mensaje.
	`, resetLink)

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(s.APIKey)
	response, err := client.Send(message)

	if err != nil {
		return fmt.Errorf("error al enviar correo con SendGrid: %v", err)
	}

	// SendGrid devuelve 202 para éxito
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		fmt.Printf("Correo de restablecimiento enviado exitosamente a: %s (Status: %d)\n", to, response.StatusCode)
		return nil
	}

	return fmt.Errorf("SendGrid devolvió error: Status %d, Body: %s", response.StatusCode, response.Body)
}
