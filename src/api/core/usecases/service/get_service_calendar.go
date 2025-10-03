package service

import (
	"context"
	"time"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
	"iycds2025_api/src/api/utils"
)

type GetServiceCalendarUseCase struct {
	serviceRepository    interfaces.Service
	appointmentRepository interfaces.Appointment
}

func NewGetServiceCalendarUseCase(serviceRepo interfaces.Service, appointmentRepo interfaces.Appointment) *GetServiceCalendarUseCase {
	return &GetServiceCalendarUseCase{
		serviceRepository:     serviceRepo,
		appointmentRepository: appointmentRepo,
	}
}

func (uc *GetServiceCalendarUseCase) Execute(serviceID int) (*entities.CalendarResponse, *errors.APIError) {
	ctx := context.Background()
	
	// Verificar que el servicio existe y está activo
	service, err := uc.serviceRepository.GetByID(ctx, int64(serviceID))
	if err != nil {
		return nil, errors.NewNotFound("Service not found")
	}

	if service.Status != "active" {
		return nil, errors.NewBadRequest("Service is not active")
	}

	// Generar los próximos 30 días desde hoy
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	var calendarDays []entities.CalendarDay
	
	for i := 0; i < 30; i++ {
		currentDate := today.AddDate(0, 0, i)
		dateStr := currentDate.Format("2006-01-02")
		dayOfWeek := utils.GetDayOfWeekInSpanish(currentDate.Weekday())

		// Obtener disponibilidad del día desde el JSON
		dayName, dayErr := utils.GetDayOfWeek(dateStr)
		if dayErr != nil {
			calendarDays = append(calendarDays, entities.CalendarDay{
				Date:            dateStr,
				DayOfWeek:       dayOfWeek,
				HasAvailability: false,
				AvailableSlots:  0,
				TotalSlots:      0,
			})
			continue
		}

		dayAvailability, availErr := utils.GetDayAvailabilityFromJSON(service.Availability, dayName)
		if availErr != nil || !dayAvailability.Available {
			calendarDays = append(calendarDays, entities.CalendarDay{
				Date:            dateStr,
				DayOfWeek:       dayOfWeek,
				HasAvailability: false,
				AvailableSlots:  0,
				TotalSlots:      0,
			})
			continue
		}

		// Generar time slots para este día
		timeSlots, slotsErr := utils.GenerateTimeSlots(dayAvailability.Start, dayAvailability.End)
		if slotsErr != nil {
			calendarDays = append(calendarDays, entities.CalendarDay{
				Date:            dateStr,
				DayOfWeek:       dayOfWeek,
				HasAvailability: false,
				AvailableSlots:  0,
				TotalSlots:      0,
			})
			continue
		}

		// Por ahora asumimos que todos los slots están disponibles
		// En el futuro podríamos verificar appointments existentes aquí
		totalSlots := len(timeSlots)
		availableSlots := totalSlots

		calendarDays = append(calendarDays, entities.CalendarDay{
			Date:            dateStr,
			DayOfWeek:       dayOfWeek,
			HasAvailability: totalSlots > 0,
			AvailableSlots:  availableSlots,
			TotalSlots:      totalSlots,
		})
	}

	response := &entities.CalendarResponse{
		ServiceID:    serviceID,
		ServiceTitle: service.Title,
		StartDate:    today.Format("2006-01-02"),
		EndDate:      today.AddDate(0, 0, 29).Format("2006-01-02"),
		Days:         calendarDays,
	}

	return response, nil
}
