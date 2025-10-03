package appointment

import (
	"context"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type GetServiceAvailability interface {
	Execute(ctx context.Context, serviceID int64, date string) (*entities.AvailabilityResponse, error)
}

type GetServiceAvailabilityImpl struct {
	Service     interfaces.Service
	Appointment interfaces.Appointment
}

func (uc *GetServiceAvailabilityImpl) Execute(ctx context.Context, serviceID int64, date string) (*entities.AvailabilityResponse, error) {
	// Validar formato de fecha
	if !utils.ValidateDateFormat(date) {
		return nil, errors.NewBadRequest("Invalid date format. Use YYYY-MM-DD")
	}

	// Verificar que la fecha no sea pasada
	if !utils.IsDateInFuture(date) {
		return nil, errors.NewBadRequest("Cannot check availability for past dates")
	}

	// Obtener el servicio
	service, err := uc.Service.GetByID(ctx, serviceID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if service == nil {
		return nil, errors.NewNotFound("Service not found")
	}
	if service.Status != "active" {
		return nil, errors.NewBadRequest("Service is not active")
	}

	// Obtener el día de la semana
	dayOfWeek, err := utils.GetDayOfWeek(date)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid date: " + err.Error())
	}

	// Obtener la disponibilidad del día desde el JSON del servicio
	dayAvailability, err := utils.GetDayAvailabilityFromJSON(service.Availability, dayOfWeek)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to parse service availability: " + err.Error())
	}

	if !dayAvailability.Available {
		// El servicio no está disponible este día de la semana
		return &entities.AvailabilityResponse{
			Date:      date,
			DayOfWeek: dayOfWeek,
			TimeSlots: []entities.TimeSlot{},
		}, nil
	}

	// Generar todos los slots posibles para el día
	allSlots, err := utils.GenerateTimeSlots(dayAvailability.Start, dayAvailability.End)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid time range in service availability: " + err.Error())
	}

	// Obtener citas ocupadas para esta fecha y servicio
	occupiedAppointments, err := uc.Appointment.GetByServiceIDAndDate(ctx, serviceID, date)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get occupied slots: " + err.Error())
	}

	// Crear mapa de slots ocupados para búsqueda rápida
	occupiedSlots := make(map[string]bool)
	for _, appointment := range occupiedAppointments {
		occupiedSlots[appointment.TimeSlot] = true
	}

	// Construir respuesta con disponibilidad de cada slot
	timeSlots := make([]entities.TimeSlot, len(allSlots))
	for i, slot := range allSlots {
		timeSlots[i] = entities.TimeSlot{
			Time:      slot,
			Available: !occupiedSlots[slot],
		}
	}

	return &entities.AvailabilityResponse{
		Date:      date,
		DayOfWeek: dayOfWeek,
		TimeSlots: timeSlots,
	}, nil
}
