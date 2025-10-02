package user

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type UpdateUser interface {
	Execute(ctx context.Context, userID int64, userUpdate *entities.UserUpdate) (*entities.UserResponse, error)
}

type UpdateUserImpl struct {
	User interfaces.User
}

func (uc *UpdateUserImpl) Execute(ctx context.Context, userID int64, userUpdate *entities.UserUpdate) (*entities.UserResponse, error) {
	// Verificar que el usuario existe
	existing, err := uc.User.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get user: " + err.Error())
	}
	if existing == nil {
		return nil, errors.NewNotFound("User not found")
	}

	// Si se está actualizando el email, verificar que no esté en uso
	if userUpdate.Email != "" && userUpdate.Email != existing.Email {
		existingUser, err := uc.User.GetByEmail(ctx, userUpdate.Email)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to check email availability: " + err.Error())
		}
		if existingUser != nil {
			return nil, errors.NewBadRequest("Email is already in use")
		}
	}

	// Actualizar el usuario
	updatedUser, err := uc.User.Update(ctx, userID, userUpdate)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to update user: " + err.Error())
	}

	// Convertir a response (sin password)
	return uc.toUserResponse(updatedUser), nil
}

func (uc *UpdateUserImpl) toUserResponse(user *entities.User) *entities.UserResponse {
	return &entities.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Locality:  user.Locality,
		Province:  user.Province,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
