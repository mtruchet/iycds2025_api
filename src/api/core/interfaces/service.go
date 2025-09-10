package interfaces

import (
	"context"
	"iycds2025_api/src/api/core/entities"
)

type Service interface {
	Create(ctx context.Context, service *entities.ServiceCreate, userID int64) (*entities.Service, error)
	GetByID(ctx context.Context, id int64) (*entities.Service, error)
	GetByUserID(ctx context.Context, userID int64) ([]*entities.Service, error)
	GetAllActive(ctx context.Context) ([]*entities.Service, error)
	Update(ctx context.Context, id int64, service *entities.ServiceUpdate, userID int64) (*entities.Service, error)
	UpdateStatus(ctx context.Context, id int64, status string, userID int64) error
	Delete(ctx context.Context, id int64, userID int64) error
}
