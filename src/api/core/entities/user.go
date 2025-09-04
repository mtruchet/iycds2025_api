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
