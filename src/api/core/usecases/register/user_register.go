package register

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type UserRegister interface {
	Execute(ctx context.Context, userRequest *entities.UserRegister) (*entities.User, error)
}

type UserRegisterImpl struct {
	User interfaces.User
}

func (uc *UserRegisterImpl) Execute(ctx context.Context, userRequest *entities.UserRegister) (*entities.User, error) {
	// Validar que las contraseñas coincidan
	if userRequest.Password != userRequest.ConfirmPassword {
		return nil, errors.NewBadRequest("Passwords do not match")
	}

	// Verificar si el email ya existe
	existingUser, err := uc.User.GetByEmail(ctx, userRequest.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check existing user")
	}
	if existingUser != nil {
		return nil, errors.NewBadRequest("Email already exists")
	}

	// Encriptar la contraseña
	hashedPassword, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid password: " + err.Error())
	}

	// Actualizar la contraseña en el request con la versión encriptada
	userRequest.Password = hashedPassword

	// Crear el usuario
	user, err := uc.User.Create(ctx, userRequest)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to create user")
	}

	return user, nil
}
