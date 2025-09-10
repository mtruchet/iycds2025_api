package service

import (
	"context"
	"encoding/json"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type ListMyServices interface {
	Execute(ctx context.Context, userID int64) (*entities.ServiceListResponse, error)
}

type ListMyServicesImpl struct {
	Service interfaces.Service
}

func (uc *ListMyServicesImpl) Execute(ctx context.Context, userID int64) (*entities.ServiceListResponse, error) {
	// Obtener servicios del usuario
	services, err := uc.Service.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get services: " + err.Error())
	}

	// Convertir a response
	serviceResponses := make([]entities.ServiceResponse, len(services))
	for i, service := range services {
		response, err := uc.toServiceResponse(service)
		if err != nil {
			return nil, err
		}
		serviceResponses[i] = *response
	}

	return &entities.ServiceListResponse{
		Services: serviceResponses,
		Total:    len(serviceResponses),
	}, nil
}

func (uc *ListMyServicesImpl) toServiceResponse(service *entities.Service) (*entities.ServiceResponse, error) {
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
