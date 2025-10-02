package entities

import "time"

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Locality   string    `json:"locality"`
	Province   string    `json:"province"`
	Phone      string    `json:"phone"`
	FirstLogin bool      `json:"first_login"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserRegister representa la solicitud de registro de usuario
type UserRegister struct {
	Name            string `json:"name" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=16"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=16"`
	Locality        string `json:"locality" validate:"required,min=2,max=100"`
	Province        string `json:"province" validate:"required,min=2,max=100"`
	Phone           string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// UserUpdate representa la solicitud de actualización de usuario
type UserUpdate struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Locality string `json:"locality" validate:"omitempty,min=2,max=100"`
	Province string `json:"province" validate:"omitempty,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// UserResponse representa la respuesta de información de usuario (sin password)
type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Locality  string `json:"locality"`
	Province  string `json:"province"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
