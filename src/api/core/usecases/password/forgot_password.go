package password

import (
	"context"
	"fmt"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/services/mail"
)

type ForgotPassword interface {
	Execute(ctx context.Context, request *entities.ForgotPassword) (bool, error)
}

type ForgotPasswordImpl struct {
	User         interfaces.User
	EmailService mail.EmailService
	FrontendURL  string
}

func (uc *ForgotPasswordImpl) Execute(ctx context.Context, request *entities.ForgotPassword) (bool, error) {
	// Verificar si el usuario existe
	user, err := uc.User.GetByEmail(ctx, request.Email)
	if err != nil {
		return false, err
	}
	if user == nil {
		// Por seguridad, no revelar si un email existe o no
		return true, nil
	}

	// Crear un token de restablecimiento
	token, expiresAt, err := uc.User.CreatePasswordResetToken(ctx, user.ID)
	if err != nil {
		return false, errors.NewInternalServerError("Failed to create reset token")
	}

	// Construir el enlace de restablecimiento
	resetLink := fmt.Sprintf("%s/reset-password/%s", uc.FrontendURL, token)

	// Enviar correo electrónico
	err = uc.EmailService.SendPasswordResetEmail(request.Email, resetLink)
	if err != nil {
		// Registrar el error pero no devolverlo al cliente por seguridad
		fmt.Printf("Error sending password reset email to %s: %v\n", request.Email, err)
	} else {
		fmt.Printf("Password reset email sent successfully to %s. Expires at: %s\n",
			request.Email, expiresAt.Format("2006-01-02 15:04:05"))
	}

	// Por seguridad, siempre devolvemos éxito aunque el email no exista
	return true, nil
}
