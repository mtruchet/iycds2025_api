package entities

import "time"

// Appointment representa una cita/reserva
type Appointment struct {
	ID         int64     `json:"id"`
	ServiceID  int64     `json:"service_id"`
	ClientID   int64     `json:"client_id"`
	ProviderID int64     `json:"provider_id"`
	Date       string    `json:"date"`       // YYYY-MM-DD
	TimeSlot   string    `json:"time_slot"`  // "HH:MM-HH:MM" ej: "09:00-09:30"
	Status     string    `json:"status"`     // "pending", "accepted", "rejected", "cancelled", "completed"
	Notes      string    `json:"notes"`      // Notas del cliente
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AppointmentCreate representa la solicitud de creación de cita
type AppointmentCreate struct {
	ServiceID int64  `json:"service_id" validate:"required"`
	Date      string `json:"date" validate:"required"`        // YYYY-MM-DD
	TimeSlot  string `json:"time_slot" validate:"required"`   // "HH:MM-HH:MM"
	Notes     string `json:"notes" validate:"omitempty,max=500"`
}

// AppointmentUpdate representa la actualización de estado de cita
type AppointmentUpdate struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected cancelled completed"`
}

// AppointmentResponse representa la respuesta de una cita con información del servicio
type AppointmentResponse struct {
	ID         int64           `json:"id"`
	Service    ServiceResponse `json:"service"`
	ClientID   int64           `json:"client_id"`
	ProviderID int64           `json:"provider_id"`
	Date       string          `json:"date"`
	TimeSlot   string          `json:"time_slot"`
	Status     string          `json:"status"`
	Notes      string          `json:"notes"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// TimeSlot representa un slot de tiempo disponible
type TimeSlot struct {
	Time      string `json:"time"`       // "HH:MM-HH:MM"
	Available bool   `json:"available"`  // true si está libre
}

// AvailabilityResponse representa la disponibilidad de un servicio en una fecha específica
type AvailabilityResponse struct {
	Date      string     `json:"date"`
	DayOfWeek string     `json:"day_of_week"`
	TimeSlots []TimeSlot `json:"time_slots"`
}

// CalendarDay representa un día en el calendario con su disponibilidad
type CalendarDay struct {
	Date            string `json:"date"`              // "YYYY-MM-DD"
	DayOfWeek       string `json:"day_of_week"`       // "monday", "tuesday", etc.
	HasAvailability bool   `json:"has_availability"`  // true si tiene horarios configurados
	AvailableSlots  int    `json:"available_slots"`   // número de slots disponibles
	TotalSlots      int    `json:"total_slots"`       // número total de slots configurados
}

// CalendarResponse representa la respuesta del calendario de un servicio
type CalendarResponse struct {
	ServiceID    int           `json:"service_id"`
	ServiceTitle string        `json:"service_title"`
	StartDate    string        `json:"start_date"`    // "YYYY-MM-DD"
	EndDate      string        `json:"end_date"`      // "YYYY-MM-DD"
	Days         []CalendarDay `json:"days"`
}
