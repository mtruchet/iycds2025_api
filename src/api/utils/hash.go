package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword genera un hash seguro para una contraseña.
func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password too long")
	}

	// Agregar validación de contraseña segura
	if len(password) < 8 || !containsSpecialOrNumber(password) {
		return "", errors.New("password must be at least 8 characters long and contain a number or special character")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Función auxiliar para verificar caracteres especiales o números
func containsSpecialOrNumber(password string) bool {
	for _, char := range password {
		if unicode.IsDigit(char) || unicode.IsPunct(char) {
			return true
		}
	}
	return false
}

// CheckPasswordHash compara una contraseña con su hash.
func CheckPasswordHash(password, hash string) bool {
	if len(password) > 72 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
