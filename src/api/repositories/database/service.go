package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"iycds2025_api/src/api/core/entities"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(ctx context.Context, serviceReq *entities.ServiceCreate, userID int64) (*entities.Service, error) {
	// Convertir availability y zones a JSON
	availabilityJSON, err := json.Marshal(serviceReq.Availability)
	if err != nil {
		return nil, err
	}

	zonesJSON, err := json.Marshal(serviceReq.Zones)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO services (title, description, user_id, category, price, availability, zones, status, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		serviceReq.Title,
		serviceReq.Description,
		userID,
		serviceReq.Category,
		serviceReq.Price,
		string(availabilityJSON),
		string(zonesJSON),
		serviceReq.ImageURL,
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

func (r *ServiceRepository) GetByID(ctx context.Context, id int64) (*entities.Service, error) {
	query := `
		SELECT id, title, description, user_id, category, price, availability, zones, status, image_url, created_at, updated_at
		FROM services
		WHERE id = ?
	`

	var service entities.Service
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Title,
		&service.Description,
		&service.UserID,
		&service.Category,
		&service.Price,
		&service.Availability,
		&service.Zones,
		&service.Status,
		&service.ImageURL,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepository) GetByUserID(ctx context.Context, userID int64) ([]*entities.Service, error) {
	query := `
		SELECT id, title, description, user_id, category, price, availability, zones, status, image_url, created_at, updated_at
		FROM services
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*entities.Service
	for rows.Next() {
		var service entities.Service
		err := rows.Scan(
			&service.ID,
			&service.Title,
			&service.Description,
			&service.UserID,
			&service.Category,
			&service.Price,
			&service.Availability,
			&service.Zones,
			&service.Status,
			&service.ImageURL,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}

func (r *ServiceRepository) GetAllActive(ctx context.Context) ([]*entities.Service, error) {
	query := `
		SELECT id, title, description, user_id, category, price, availability, zones, status, image_url, created_at, updated_at
		FROM services
		WHERE status = 'active'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*entities.Service
	for rows.Next() {
		var service entities.Service
		err := rows.Scan(
			&service.ID,
			&service.Title,
			&service.Description,
			&service.UserID,
			&service.Category,
			&service.Price,
			&service.Availability,
			&service.Zones,
			&service.Status,
			&service.ImageURL,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}

func (r *ServiceRepository) Update(ctx context.Context, id int64, serviceReq *entities.ServiceUpdate, userID int64) (*entities.Service, error) {
	// Verificar que el servicio pertenece al usuario
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil || existing.UserID != userID {
		return nil, nil
	}

	// Construir query dinámico basado en campos presentes
	setParts := []string{}
	args := []interface{}{}

	if serviceReq.Title != "" {
		setParts = append(setParts, "title = ?")
		args = append(args, serviceReq.Title)
	}
	if serviceReq.Description != "" {
		setParts = append(setParts, "description = ?")
		args = append(args, serviceReq.Description)
	}
	if serviceReq.Category != "" {
		setParts = append(setParts, "category = ?")
		args = append(args, serviceReq.Category)
	}
	if serviceReq.Price != nil {
		setParts = append(setParts, "price = ?")
		args = append(args, *serviceReq.Price)
	}
	if serviceReq.Availability != nil {
		availabilityJSON, err := json.Marshal(serviceReq.Availability)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, "availability = ?")
		args = append(args, string(availabilityJSON))
	}
	if serviceReq.Zones != nil {
		zonesJSON, err := json.Marshal(serviceReq.Zones)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, "zones = ?")
		args = append(args, string(zonesJSON))
	}
	if serviceReq.ImageURL != "" {
		setParts = append(setParts, "image_url = ?")
		args = append(args, serviceReq.ImageURL)
	}

	if len(setParts) == 0 {
		return existing, nil // No hay cambios
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	query := "UPDATE services SET "
	for i, part := range setParts {
		if i > 0 {
			query += ", "
		}
		query += part
	}
	query += " WHERE id = ?"

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *ServiceRepository) UpdateStatus(ctx context.Context, id int64, status string, userID int64) error {
	query := `
		UPDATE services 
		SET status = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, id, userID)
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

func (r *ServiceRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM services 
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
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
