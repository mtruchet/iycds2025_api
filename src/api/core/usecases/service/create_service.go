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

type CreateService interface {
	Execute(ctx context.Context, serviceReq *entities.ServiceCreate, userID int64) (*entities.ServiceResponse, error)
}

type CreateServiceImpl struct {
	Service interfaces.Service
}

func (uc *CreateServiceImpl) Execute(ctx context.Context, serviceReq *entities.ServiceCreate, userID int64) (*entities.ServiceResponse, error) {
	// Validar categoría
	normalizedCategory, isValid := utils.NormalizeCategory(serviceReq.Category)
	if !isValid {
		return nil, errors.NewBadRequest("Invalid category. Valid categories are: " + strings.Join(utils.GetValidCategories(), ", "))
	}
	serviceReq.Category = normalizedCategory

	// Crear el servicio
	service, err := uc.Service.Create(ctx, serviceReq, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to create service: " + err.Error())
	}

	// Convertir a response
	return uc.toServiceResponse(service)
}

func (uc *CreateServiceImpl) toServiceResponse(service *entities.Service) (*entities.ServiceResponse, error) {
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
