package service

import (
	"context"

	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type UpdateServiceStatus interface {
	Execute(ctx context.Context, id int64, status string, userID int64) error
}

type UpdateServiceStatusImpl struct {
	Service interfaces.Service
}

func (uc *UpdateServiceStatusImpl) Execute(ctx context.Context, id int64, status string, userID int64) error {
	// Validar estados permitidos
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
	}
	if !validStatuses[status] {
		return errors.NewBadRequest("Invalid status. Allowed values: active, inactive")
	}

	// Verificar que el servicio existe y pertenece al usuario
	existing, err := uc.Service.GetByID(ctx, id)
	if err != nil {
		return errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if existing == nil {
		return errors.NewNotFound("Service not found")
	}
	if existing.UserID != userID {
		return errors.NewUnauthorized("You don't have permission to update this service")
	}

	// Verificar si el estado ya es el mismo
	if existing.Status == status {
		return errors.NewBadRequest("Service is already " + status)
	}

	// Actualizar el estado del servicio
	err = uc.Service.UpdateStatus(ctx, id, status, userID)
	if err != nil {
		return errors.NewInternalServerError("Failed to update service status: " + err.Error())
	}

	return nil
}
