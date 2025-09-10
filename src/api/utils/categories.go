package utils

import (
	"strings"
)

// ValidCategories contiene las categorías permitidas para servicios
var ValidCategories = []string{
	"Limpieza",
	"Jardinería", 
	"Plomería",
	"Electricidad",
	"Carpintería",
	"Pintura",
	"Mecánica",
	"Tecnología",
	"Educación",
	"Salud",
	"Belleza",
	"Mascotas",
	"Transporte",
	"Eventos",
	"Fotografía",
	"Cocina",
	"Fitness",
	"Música",
	"Idiomas",
	"Otros",
}

// IsValidCategory verifica si una categoría es válida
// Ignora mayúsculas/minúsculas para ser más flexible
func IsValidCategory(category string) bool {
	for _, valid := range ValidCategories {
		if strings.EqualFold(valid, category) {
			return true
		}
	}
	return false
}

// NormalizeCategory normaliza una categoría a su formato correcto
// Devuelve la categoría con la capitalización correcta si es válida
func NormalizeCategory(category string) (string, bool) {
	for _, valid := range ValidCategories {
		if strings.EqualFold(valid, category) {
			return valid, true
		}
	}
	return "", false
}

// GetValidCategories devuelve la lista de categorías válidas
func GetValidCategories() []string {
	// Devolver una copia para evitar modificaciones accidentales
	categories := make([]string, len(ValidCategories))
	copy(categories, ValidCategories)
	return categories
}
