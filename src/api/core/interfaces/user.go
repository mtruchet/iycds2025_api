package interfaces

import (
	"context"
	"time"

	"iycds2025_api/src/api/core/entities"
)

type User interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.UserRegister) (*entities.User, error)
	GetPermissions(ctx context.Context, userID int64) ([]string, error)
	GetRole(ctx context.Context, userID int64) (string, error)
	CreatePasswordResetToken(ctx context.Context, userID int64) (string, time.Time, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
}
