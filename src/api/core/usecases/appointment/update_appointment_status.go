package appointment

import (
	"context"
	"database/sql"

	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/interfaces"
)

type UpdateAppointmentStatus interface {
	Execute(ctx context.Context, appointmentID int64, status string, userID int64) error
}

type UpdateAppointmentStatusImpl struct {
	Appointment interfaces.Appointment
}

func (uc *UpdateAppointmentStatusImpl) Execute(ctx context.Context, appointmentID int64, status string, userID int64) error {
	// Obtener la cita para verificar permisos
	appointment, err := uc.Appointment.GetByID(ctx, appointmentID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get appointment: " + err.Error())
	}
	if appointment == nil {
		return errors.NewNotFound("Appointment not found")
	}

	// Verificar permisos según el estado a actualizar
	switch status {
	case "accepted", "rejected":
		// Solo el proveedor puede aceptar/rechazar
		if appointment.ProviderID != userID {
			return errors.NewUnauthorized("Only the service provider can accept or reject appointments")
		}
		// Solo se pueden aceptar/rechazar citas pendientes
		if appointment.Status != "pending" {
			return errors.NewBadRequest("Can only accept or reject pending appointments")
		}

	case "cancelled":
		// Tanto cliente como proveedor pueden cancelar
		if appointment.ClientID != userID && appointment.ProviderID != userID {
			return errors.NewUnauthorized("You don't have permission to cancel this appointment")
		}
		// No se pueden cancelar citas ya completadas o rechazadas
		if appointment.Status == "completed" || appointment.Status == "rejected" {
			return errors.NewBadRequest("Cannot cancel completed or rejected appointments")
		}

	case "completed":
		// Solo el proveedor puede marcar como completada
		if appointment.ProviderID != userID {
			return errors.NewUnauthorized("Only the service provider can mark appointments as completed")
		}
		// Solo se pueden completar citas aceptadas
		if appointment.Status != "accepted" {
			return errors.NewBadRequest("Can only complete accepted appointments")
		}

	default:
		return errors.NewBadRequest("Invalid status. Valid values: accepted, rejected, cancelled, completed")
	}

	// Actualizar el estado
	err = uc.Appointment.UpdateStatus(ctx, appointmentID, status, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFound("Appointment not found or you don't have permission")
		}
		return errors.NewInternalServerError("Failed to update appointment status: " + err.Error())
	}

	return nil
}
