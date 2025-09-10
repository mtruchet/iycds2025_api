package dependencies

import (
	"os"

	"iycds2025_api/configs"
	"iycds2025_api/src/api/core/usecases/login"
	"iycds2025_api/src/api/core/usecases/password"
	"iycds2025_api/src/api/core/usecases/register"
	"iycds2025_api/src/api/core/usecases/service"
	"iycds2025_api/src/api/infrastructure/entrypoints/api"
	apiHandlers "iycds2025_api/src/api/infrastructure/entrypoints/api/handlers"
	"iycds2025_api/src/api/repositories/database"
)

type HandlerContainer struct {
	Ping            api.Handler
	UserLogin       api.Handler
	UserRegister    api.Handler
	PasswordForgot  api.Handler
	PasswordReset   api.Handler
	ServiceCreate   api.Handler
	ServiceUpdate   api.Handler
	ServiceDelete   api.Handler
	ServiceList     api.Handler
	ServiceListAll  api.Handler
	Categories      api.Handler
}

func Start() *HandlerContainer {
	// Database
	db := configs.ConnectDatabase()

	// Repositories
	userRepo := &database.UserRepository{
		DB: db,
	}

	serviceRepo := database.NewServiceRepository(db)

	// Services
	emailService := configs.NewEmailService()

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

	// Service use cases
	createServiceUseCase := &service.CreateServiceImpl{
		Service: serviceRepo,
	}

	updateServiceUseCase := &service.UpdateServiceImpl{
		Service: serviceRepo,
	}

	deleteServiceUseCase := &service.DeleteServiceImpl{
		Service: serviceRepo,
	}

	listMyServicesUseCase := &service.ListMyServicesImpl{
		Service: serviceRepo,
	}

	listAllServicesUseCase := &service.ListAllServicesImpl{
		Service: serviceRepo,
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
	handlers.ServiceCreate = &apiHandlers.ServiceCreateHandler{
		CreateService: createServiceUseCase,
	}
	handlers.ServiceUpdate = &apiHandlers.ServiceUpdateHandler{
		UpdateService: updateServiceUseCase,
	}
	handlers.ServiceDelete = &apiHandlers.ServiceDeleteHandler{
		DeleteService: deleteServiceUseCase,
	}
	handlers.ServiceList = &apiHandlers.ServiceListHandler{
		ListMyServices: listMyServicesUseCase,
	}
	handlers.ServiceListAll = &apiHandlers.ServiceListAllHandler{
		ListAllServices: listAllServicesUseCase,
	}
	handlers.Categories = &apiHandlers.CategoriesHandler{}

	return &handlers
}
