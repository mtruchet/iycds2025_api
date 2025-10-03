package appointment

import (
	"context"
	"database/sql"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type CreateAppointment interface {
	Execute(ctx context.Context, appointmentReq *entities.AppointmentCreate, clientID int64) (*entities.AppointmentResponse, error)
}

type CreateAppointmentImpl struct {
	Service     interfaces.Service
	Appointment interfaces.Appointment
}

func (uc *CreateAppointmentImpl) Execute(ctx context.Context, appointmentReq *entities.AppointmentCreate, clientID int64) (*entities.AppointmentResponse, error) {
	// Validar formato de fecha y hora
	if !utils.ValidateDateFormat(appointmentReq.Date) {
		return nil, errors.NewBadRequest("Invalid date format. Use YYYY-MM-DD")
	}

	if !utils.IsDateInFuture(appointmentReq.Date) {
		return nil, errors.NewBadRequest("Cannot create appointments for past dates")
	}

	// Validar formato de time slot (debe ser HH:MM-HH:MM)
	if !uc.validateTimeSlotFormat(appointmentReq.TimeSlot) {
		return nil, errors.NewBadRequest("Invalid time slot format. Use HH:MM-HH:MM")
	}

	// Verificar que el servicio existe y está activo
	service, err := uc.Service.GetByID(ctx, appointmentReq.ServiceID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get service: " + err.Error())
	}
	if service == nil {
		return nil, errors.NewNotFound("Service not found")
	}
	if service.Status != "active" {
		return nil, errors.NewBadRequest("Service is not active")
	}

	// Verificar que el cliente no sea el mismo que el proveedor
	if service.UserID == clientID {
		return nil, errors.NewBadRequest("Cannot create appointment for your own service")
	}

	// Verificar que el horario esté dentro de la disponibilidad del servicio
	dayOfWeek, err := utils.GetDayOfWeek(appointmentReq.Date)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid date: " + err.Error())
	}

	dayAvailability, err := utils.GetDayAvailabilityFromJSON(service.Availability, dayOfWeek)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to parse service availability: " + err.Error())
	}

	if !dayAvailability.Available {
		return nil, errors.NewBadRequest("Service is not available on " + dayOfWeek + "s")
	}

	// Verificar que el time slot está dentro del rango de disponibilidad
	if !uc.isTimeSlotInRange(appointmentReq.TimeSlot, dayAvailability.Start, dayAvailability.End) {
		return nil, errors.NewBadRequest("Time slot is outside service availability hours")
	}

	// Crear la cita
	appointment, err := uc.Appointment.Create(ctx, appointmentReq, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewBadRequest("Time slot is already occupied or service not found")
		}
		return nil, errors.NewInternalServerError("Failed to create appointment: " + err.Error())
	}

	// Convertir a response con información del servicio
	return uc.toAppointmentResponse(appointment, service)
}

func (uc *CreateAppointmentImpl) validateTimeSlotFormat(timeSlot string) bool {
	// Formato debe ser HH:MM-HH:MM
	if len(timeSlot) != 11 || timeSlot[5] != '-' {
		return false
	}

	startTime := timeSlot[:5]
	endTime := timeSlot[6:]

	return utils.ValidateTimeFormat(startTime) && utils.ValidateTimeFormat(endTime)
}

func (uc *CreateAppointmentImpl) isTimeSlotInRange(timeSlot, start, end string) bool {
	// Generar todos los slots válidos para el rango
	validSlots, err := utils.GenerateTimeSlots(start, end)
	if err != nil {
		return false
	}

	// Verificar si el slot solicitado está en la lista de válidos
	for _, validSlot := range validSlots {
		if validSlot == timeSlot {
			return true
		}
	}

	return false
}

func (uc *CreateAppointmentImpl) toAppointmentResponse(appointment *entities.Appointment, service *entities.Service) (*entities.AppointmentResponse, error) {
	// Convertir servicio a response (reutilizar lógica existente)
	serviceResponse, err := uc.convertServiceToResponse(service)
	if err != nil {
		return nil, err
	}

	return &entities.AppointmentResponse{
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
	}, nil
}

func (uc *CreateAppointmentImpl) convertServiceToResponse(service *entities.Service) (*entities.ServiceResponse, error) {
	// Esta función duplica la lógica de los otros use cases de servicio
	// En un refactor futuro, se podría extraer a un utility común
	
	// Por simplicidad, returnamos una versión básica sin parsear availability y zones
	return &entities.ServiceResponse{
		ID:          service.ID,
		Title:       service.Title,
		Description: service.Description,
		Category:    service.Category,
		Price:       service.Price,
		Status:      service.Status,
		ImageURL:    service.ImageURL,
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
		// availability y zones se pueden parsear aquí si es necesario
	}, nil
}
