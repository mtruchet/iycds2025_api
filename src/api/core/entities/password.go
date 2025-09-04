package entities

import "time"

// ChangePassword representa la solicitud para cambiar la contraseña
type ChangePassword struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=16"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=16"`
}

// ForgotPassword representa la solicitud para olvidó su contraseña
type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPassword representa la solicitud para restablecer contraseña
type ResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=16"`
}

// PasswordResetToken representa un token para restablecer la contraseña
type PasswordResetToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}
