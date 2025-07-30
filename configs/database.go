package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConfig contiene la configuración para la conexión a la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDatabaseConfig retorna la configuración de la base de datos según el entorno
func GetDatabaseConfig() DatabaseConfig {
	appEnv := os.Getenv("APP_ENV")

	var config DatabaseConfig

	if appEnv == "production" {
		// Configuración para entorno de producción (si fuera necesario)
		config = DatabaseConfig{
			Host:     getEnvOrDefault("PROD_DB_HOST", "localhost"),
			Port:     getEnvOrDefault("PROD_DB_PORT", "3306"),
			User:     getEnvOrDefault("PROD_DB_USER", "root"),
			Password: getEnvOrDefault("PROD_DB_PASSWORD", ""),
			DBName:   getEnvOrDefault("PROD_DB_NAME", "iycds2025"),
		}
		log.Println("Using production database configuration")
	} else {
		// Configuración para entorno de desarrollo (local)
		config = DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "3306"),
			User:     getEnvOrDefault("DB_USER", "root"),
			Password: getEnvOrDefault("DB_PASSWORD", ""),
			DBName:   getEnvOrDefault("DB_NAME", "iycds2025"),
		}
		log.Println("Using development database configuration")
	}
	return config
}

// Helper para obtener variable de entorno con valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ConnectDatabase establece la conexión con la base de datos MySQL
func ConnectDatabase() *sql.DB {
	config := GetDatabaseConfig()

	// Construir el DSN (Data Source Name) para la conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	log.Printf("Connecting to database at %s:%s...\n", config.Host, config.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection successful!")
	return db
}

// GetDatabaseConnection retorna una instancia de conexión a la base de datos
// Esta función puede ser usada para health checks
func GetDatabaseConnection() *sql.DB {
	db := ConnectDatabase()
	if db == nil {
		log.Printf("Error connecting to database")
		return nil
	}
	return db
}
