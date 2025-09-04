package password

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/interfaces"
)

type ResetPassword interface {
	Execute(ctx context.Context, request *entities.ResetPassword) error
}

type ResetPasswordImpl struct {
	User interfaces.User
}

func (uc *ResetPasswordImpl) Execute(ctx context.Context, request *entities.ResetPassword) error {
	// Restablecer la contraseña usando el token
	err := uc.User.ResetPassword(ctx, request.Token, request.NewPassword)
	if err != nil {
		return err
	}

	return nil
}
