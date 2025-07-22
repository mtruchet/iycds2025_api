package app

import (
	"fmt"
	"log"
	"os"
	"iycds2025_api/src/api/infrastructure/dependencies"
	"iycds2025_api/src/api/middleware"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Start() {
	fmt.Println("Starting IYCDS2025 API")

	router = createRouter()

	// Obtener el puerto desde la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	fmt.Printf("Running API in %s environment on port %s\n", appEnv, port)

	err := router.Run(":" + port)
	if err != nil {
		fmt.Printf("Error when starting the API: %s \n", err)
	}
}

func createRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSConfig())
	configureEnv()
	handlers := dependencies.Start()
	configureURLMappings(router, handlers)
	return router
}

func configureEnv() {
	// Configuración del entorno
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
		os.Setenv("APP_ENV", appEnv)
		log.Println("APP_ENV no configurado, usando valor por defecto:", appEnv)
	}
	fmt.Printf("Running in %s environment\n", appEnv)

	// Configurar modo Gin según el entorno
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
}
