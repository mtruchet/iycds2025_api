package login

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type UserLogin interface {
	Execute(ctx context.Context, userRequest *entities.Login) (string, error)
}

type UserLoginImpl struct {
	User interfaces.User
}

func (uc *UserLoginImpl) Execute(ctx context.Context, userRequest *entities.Login) (string, error) {
	user, err := uc.User.GetByEmail(ctx, userRequest.Email)
	if err != nil || user == nil {
		return "", errors.NewUnauthorized("Invalid credentials")
	}

	checkPasswordHash := utils.CheckPasswordHash(userRequest.Password, user.Password)
	if !checkPasswordHash {
		return "", errors.NewUnauthorized("Invalid credentials")
	}

	permissions, err := uc.User.GetPermissions(ctx, user.ID)
	if err != nil {
		return "", errors.NewInternalServerError("Failed to fetch permissions")
	}

	role, err := uc.User.GetRole(ctx, user.ID)
	if err != nil {
		return "", errors.NewInternalServerError("Failed to fetch role")
	}

	token, err := utils.GenerateJWT(user.ID, role, permissions, user.FirstLogin)
	if err != nil {
		return "", errors.NewInternalServerError("Failed to generate JWT token")
	}
	return token, nil
}
