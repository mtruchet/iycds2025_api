package service

import (
	"context"
	"encoding/json"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type GetServiceByID interface {
	Execute(ctx context.Context, id int64) (*entities.ServiceResponse, error)
}

type GetServiceByIDImpl struct {
	Service interfaces.Service
}

func (uc *GetServiceByIDImpl) Execute(ctx context.Context, id int64) (*entities.ServiceResponse, error) {
	// Obtener el servicio por ID
	service, err := uc.Service.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if service == nil {
		return nil, errors.NewNotFound("Service not found")
	}

	// Solo retornar servicios activos para consultas públicas
	if service.Status != "active" {
		return nil, errors.NewNotFound("Service not found")
	}

	// Convertir a response
	return uc.toServiceResponse(service)
}

func (uc *GetServiceByIDImpl) toServiceResponse(service *entities.Service) (*entities.ServiceResponse, error) {
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
