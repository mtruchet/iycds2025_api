package appointment

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type ListServiceAppointments interface {
	Execute(ctx context.Context, serviceID int64, providerID int64) ([]*entities.AppointmentResponse, error)
}

type ListServiceAppointmentsImpl struct {
	Service     interfaces.Service
	Appointment interfaces.Appointment
}

func (uc *ListServiceAppointmentsImpl) Execute(ctx context.Context, serviceID int64, providerID int64) ([]*entities.AppointmentResponse, error) {
	// Verificar que el servicio existe y pertenece al proveedor
	service, err := uc.Service.GetByID(ctx, serviceID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if service == nil {
		return nil, errors.NewNotFound("Service not found")
	}
	if service.UserID != providerID {
		return nil, errors.NewUnauthorized("You don't have permission to view appointments for this service")
	}

	// Obtener appointments del servicio
	appointments, err := uc.Appointment.GetByServiceID(ctx, serviceID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get appointments: " + err.Error())
	}

	// Convertir a response
	var responses []*entities.AppointmentResponse
	for _, appointment := range appointments {
		serviceResponse := &entities.ServiceResponse{
			ID:          service.ID,
			Title:       service.Title,
			Description: service.Description,
			Category:    service.Category,
			Price:       service.Price,
			Status:      service.Status,
			ImageURL:    service.ImageURL,
			CreatedAt:   service.CreatedAt,
			UpdatedAt:   service.UpdatedAt,
		}

		responses = append(responses, &entities.AppointmentResponse{
			ID:         appointment.ID,
			Service:    *serviceResponse,
			ClientID:   appointment.ClientID,
			ProviderID: appointment.ProviderID,
			Date:       appointment.Date,
			TimeSlot:   appointment.TimeSlot,
			Status:     appointment.Status,
			Notes:      appointment.Notes,
			CreatedAt:  appointment.CreatedAt,
			UpdatedAt:  appointment.UpdatedAt,
		})
	}

	return responses, nil
}
