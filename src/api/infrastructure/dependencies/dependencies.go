package dependencies

import (
	"iycds2025_api/configs"
	"iycds2025_api/src/api/infrastructure/entrypoints/api"
	apiHandlers "iycds2025_api/src/api/infrastructure/entrypoints/api/handlers"
	"iycds2025_api/src/api/repositories/database"
)

type HandlerContainer struct {
	Ping api.Handler
}

func Start() *HandlerContainer {
	// Database
	db := configs.ConnectDatabase()

	// Repositories (por si los necesitas en el futuro)
	_ = &database.Repository{
		DB: db,
	}

	// Handlers
	handlers := HandlerContainer{}

	handlers.Ping = &apiHandlers.Ping{}

	return &handlers
}
