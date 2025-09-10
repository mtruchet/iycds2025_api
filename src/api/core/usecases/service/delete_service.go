package service

import (
	"context"

	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type DeleteService interface {
	Execute(ctx context.Context, id int64, userID int64) error
}

type DeleteServiceImpl struct {
	Service interfaces.Service
}

func (uc *DeleteServiceImpl) Execute(ctx context.Context, id int64, userID int64) error {
	// Verificar que el servicio existe y pertenece al usuario
	existing, err := uc.Service.GetByID(ctx, id)
	if err != nil {
		return errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if existing == nil {
		return errors.NewNotFound("Service not found")
	}
	if existing.UserID != userID {
		return errors.NewUnauthorized("You don't have permission to delete this service")
	}

	// Realizar baja lógica (cambiar estado a inactivo)
	err = uc.Service.UpdateStatus(ctx, id, "inactive", userID)
	if err != nil {
		return errors.NewInternalServerError("Failed to deactivate service: " + err.Error())
	}

	return nil
}
