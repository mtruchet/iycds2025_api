package service

import (
	"context"
	"database/sql"

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

	// Eliminar el servicio de manera definitiva
	err = uc.Service.Delete(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFound("Service not found or already deleted")
		}
		return errors.NewInternalServerError("Failed to delete service: " + err.Error())
	}

	return nil
}
