package interfaces

import (
	"context"
	"iycds2025_api/src/api/core/entities"
)

type Appointment interface {
	Create(ctx context.Context, appointment *entities.AppointmentCreate, clientID int64) (*entities.Appointment, error)
	GetByID(ctx context.Context, id int64) (*entities.Appointment, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*entities.Appointment, error)
	GetByServiceID(ctx context.Context, serviceID int64) ([]*entities.Appointment, error)
	GetByServiceIDAndDate(ctx context.Context, serviceID int64, date string) ([]*entities.Appointment, error)
	UpdateStatus(ctx context.Context, id int64, status string, userID int64) error
	Delete(ctx context.Context, id int64, clientID int64) error
}
