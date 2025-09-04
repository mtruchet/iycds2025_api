package dependencies

import (
	"os"

	"iycds2025_api/configs"
	"iycds2025_api/src/api/core/usecases/login"
	"iycds2025_api/src/api/core/usecases/password"
	"iycds2025_api/src/api/core/usecases/register"
	"iycds2025_api/src/api/infrastructure/entrypoints/api"
	apiHandlers "iycds2025_api/src/api/infrastructure/entrypoints/api/handlers"
	"iycds2025_api/src/api/repositories/database"
	"iycds2025_api/src/api/services/mail"
)

type HandlerContainer struct {
	Ping           api.Handler
	UserLogin      api.Handler
	UserRegister   api.Handler
	PasswordForgot api.Handler
	PasswordReset  api.Handler
}

func Start() *HandlerContainer {
	// Database
	db := configs.ConnectDatabase()

	// Repositories
	userRepo := &database.UserRepository{
		DB: db,
	}

	// Services
	var emailService mail.EmailService

	// Determinar qué servicio de correo usar según el entorno
	useRealEmail := os.Getenv("USE_REAL_EMAIL")

	// Por defecto en desarrollo usar mock, en producción usar real
	appEnv := os.Getenv("APP_ENV")
	if useRealEmail == "" {
		if appEnv == "production" {
			useRealEmail = "true"
		} else {
			useRealEmail = "false"
		}
	}

	if useRealEmail == "true" {
		emailService = mail.NewSMTPEmailService()
	} else {
		emailService = mail.NewMockEmailService()
	}

	// URL del frontend para los enlaces
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // URL por defecto
	}

	// Use cases
	userLoginUseCase := &login.UserLoginImpl{
		User: userRepo,
	}

	forgotPasswordUseCase := &password.ForgotPasswordImpl{
		User:         userRepo,
		EmailService: emailService,
		FrontendURL:  frontendURL,
	}

	resetPasswordUseCase := &password.ResetPasswordImpl{
		User: userRepo,
	}

	userRegisterUseCase := &register.UserRegisterImpl{
		User: userRepo,
	}

	// Handlers
	handlers := HandlerContainer{}

	handlers.Ping = &apiHandlers.Ping{}
	handlers.UserLogin = &apiHandlers.UserLogin{
		UseCase: userLoginUseCase,
	}
	handlers.UserRegister = &apiHandlers.UserRegister{
		UseCase: userRegisterUseCase,
	}
	handlers.PasswordForgot = &apiHandlers.PasswordForgot{
		UseCase: forgotPasswordUseCase,
	}
	handlers.PasswordReset = &apiHandlers.PasswordReset{
		UseCase: resetPasswordUseCase,
	}

	return &handlers
}
