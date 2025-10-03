package appointment

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type ListMyAppointments interface {
	Execute(ctx context.Context, clientID int64) ([]*entities.AppointmentResponse, error)
}

type ListMyAppointmentsImpl struct {
	Service     interfaces.Service
	Appointment interfaces.Appointment
}

func (uc *ListMyAppointmentsImpl) Execute(ctx context.Context, clientID int64) ([]*entities.AppointmentResponse, error) {
	// Obtener appointments del cliente
	appointments, err := uc.Appointment.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get appointments: " + err.Error())
	}

	// Convertir a response con información de servicios
	var responses []*entities.AppointmentResponse
	for _, appointment := range appointments {
		// Obtener información del servicio
		service, err := uc.Service.GetByID(ctx, appointment.ServiceID)
		if err != nil {
			// Si no se puede obtener el servicio, continuar con información básica
			serviceResponse := &entities.ServiceResponse{
				ID:    appointment.ServiceID,
				Title: "Service not available",
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
			continue
		}

		// Convertir servicio a response
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
