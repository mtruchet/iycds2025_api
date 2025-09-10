package service

import (
	"context"
	"encoding/json"
	"strings"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type UpdateService interface {
	Execute(ctx context.Context, id int64, serviceReq *entities.ServiceUpdate, userID int64) (*entities.ServiceResponse, error)
}

type UpdateServiceImpl struct {
	Service interfaces.Service
}

func (uc *UpdateServiceImpl) Execute(ctx context.Context, id int64, serviceReq *entities.ServiceUpdate, userID int64) (*entities.ServiceResponse, error) {
	// Verificar que el servicio existe y pertenece al usuario
	existing, err := uc.Service.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if existing == nil {
		return nil, errors.NewNotFound("Service not found")
	}
	if existing.UserID != userID {
		return nil, errors.NewUnauthorized("You don't have permission to update this service")
	}

	// Validar categoría si se está actualizando
	if serviceReq.Category != "" {
		normalizedCategory, isValid := utils.NormalizeCategory(serviceReq.Category)
		if !isValid {
			return nil, errors.NewBadRequest("Invalid category. Valid categories are: " + strings.Join(utils.GetValidCategories(), ", "))
		}
		serviceReq.Category = normalizedCategory
	}

	// Actualizar el servicio
	service, err := uc.Service.Update(ctx, id, serviceReq, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to update service: " + err.Error())
	}

	// Convertir a response
	return uc.toServiceResponse(service)
}

func (uc *UpdateServiceImpl) toServiceResponse(service *entities.Service) (*entities.ServiceResponse, error) {
	// Parsear availability JSON
	var availability map[string]interface{}
	if err := json.Unmarshal([]byte(service.Availability), &availability); err != nil {
		return nil, errors.NewInternalServerError("Failed to parse availability")
	}

	// Parsear zones JSON
	var zones []entities.Zone
	if err := json.Unmarshal([]byte(service.Zones), &zones); err != nil {
		return nil, errors.NewInternalServerError("Failed to parse zones")
	}

	return &entities.ServiceResponse{
		ID:           service.ID,
		Title:        service.Title,
		Description:  service.Description,
		Category:     service.Category,
		Price:        service.Price,
		Availability: availability,
		Zones:        zones,
		Status:       service.Status,
		ImageURL:     service.ImageURL,
		CreatedAt:    service.CreatedAt,
		UpdatedAt:    service.UpdatedAt,
	}, nil
}
