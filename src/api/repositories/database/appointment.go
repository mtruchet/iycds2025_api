package database

import (
	"context"
	"database/sql"

	"iycds2025_api/src/api/core/entities"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(ctx context.Context, appointmentReq *entities.AppointmentCreate, clientID int64) (*entities.Appointment, error) {
	// Primero obtener el provider_id del servicio
	var providerID int64
	serviceQuery := `SELECT user_id FROM services WHERE id = ? AND status = 'active'`
	err := r.db.QueryRowContext(ctx, serviceQuery, appointmentReq.ServiceID).Scan(&providerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows // Servicio no encontrado o inactivo
		}
		return nil, err
	}

	// Verificar que no existe una cita en el mismo horario y fecha
	conflictQuery := `
		SELECT COUNT(*) FROM appointments 
		WHERE service_id = ? AND date = ? AND time_slot = ? 
		AND status IN ('pending', 'accepted')
	`
	var count int
	err = r.db.QueryRowContext(ctx, conflictQuery, 
		appointmentReq.ServiceID, appointmentReq.Date, appointmentReq.TimeSlot).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, sql.ErrNoRows // Usar como indicador de conflicto
	}

	// Crear la cita
	query := `
		INSERT INTO appointments (service_id, client_id, provider_id, date, time_slot, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'pending', ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		appointmentReq.ServiceID, clientID, providerID,
		appointmentReq.Date, appointmentReq.TimeSlot, appointmentReq.Notes,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id int64) (*entities.Appointment, error) {
	var appointment entities.Appointment
	query := `
		SELECT id, service_id, client_id, provider_id, date, time_slot, status, notes, created_at, updated_at
		FROM appointments WHERE id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appointment.ID, &appointment.ServiceID, &appointment.ClientID, &appointment.ProviderID,
		&appointment.Date, &appointment.TimeSlot, &appointment.Status, &appointment.Notes,
		&appointment.CreatedAt, &appointment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &appointment, nil
}

func (r *AppointmentRepository) GetByClientID(ctx context.Context, clientID int64) ([]*entities.Appointment, error) {
	query := `
		SELECT id, service_id, client_id, provider_id, date, time_slot, status, notes, created_at, updated_at
		FROM appointments WHERE client_id = ?
		ORDER BY date DESC, time_slot DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		var appointment entities.Appointment
		err := rows.Scan(
			&appointment.ID, &appointment.ServiceID, &appointment.ClientID, &appointment.ProviderID,
			&appointment.Date, &appointment.TimeSlot, &appointment.Status, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByServiceID(ctx context.Context, serviceID int64) ([]*entities.Appointment, error) {
	query := `
		SELECT id, service_id, client_id, provider_id, date, time_slot, status, notes, created_at, updated_at
		FROM appointments WHERE service_id = ?
		ORDER BY date DESC, time_slot DESC
	`

	rows, err := r.db.QueryContext(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		var appointment entities.Appointment
		err := rows.Scan(
			&appointment.ID, &appointment.ServiceID, &appointment.ClientID, &appointment.ProviderID,
			&appointment.Date, &appointment.TimeSlot, &appointment.Status, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByServiceIDAndDate(ctx context.Context, serviceID int64, date string) ([]*entities.Appointment, error) {
	query := `
		SELECT id, service_id, client_id, provider_id, date, time_slot, status, notes, created_at, updated_at
		FROM appointments WHERE service_id = ? AND date = ? AND status IN ('pending', 'accepted')
		ORDER BY time_slot ASC
	`

	rows, err := r.db.QueryContext(ctx, query, serviceID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		var appointment entities.Appointment
		err := rows.Scan(
			&appointment.ID, &appointment.ServiceID, &appointment.ClientID, &appointment.ProviderID,
			&appointment.Date, &appointment.TimeSlot, &appointment.Status, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id int64, status string, userID int64) error {
	query := `
		UPDATE appointments 
		SET status = ?, updated_at = NOW()
		WHERE id = ? AND (client_id = ? OR provider_id = ?)
	`

	result, err := r.db.ExecContext(ctx, query, status, id, userID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, id int64, clientID int64) error {
	query := `
		DELETE FROM appointments 
		WHERE id = ? AND client_id = ? AND status = 'pending'
	`

	result, err := r.db.ExecContext(ctx, query, id, clientID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
