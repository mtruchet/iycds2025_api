package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/utils"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, name, email, password, locality, province, phone, first_login, created_at, updated_at FROM users WHERE email = ?`

	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Locality,
		&user.Province, &user.Phone, &user.FirstLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuario no encontrado
		}
		return nil, err
	}

	return &user, nil
}

// Create crea un nuevo usuario en la base de datos
func (r *UserRepository) Create(ctx context.Context, userRegister *entities.UserRegister) (*entities.User, error) {
	user := &entities.User{
		Name:       userRegister.Name,
		Email:      userRegister.Email,
		Password:   userRegister.Password,
		Locality:   userRegister.Locality,
		Province:   userRegister.Province,
		Phone:      userRegister.Phone,
		FirstLogin: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	query := `INSERT INTO users (name, email, password, locality, province, phone, first_login, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.ExecContext(ctx, query,
		user.Name, user.Email, user.Password, user.Locality,
		user.Province, user.Phone, user.FirstLogin,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (r *UserRepository) GetPermissions(ctx context.Context, userID int64) ([]string, error) {
	query := `
		SELECT p.permission_name
		FROM user_permissions up
		JOIN permissions p ON up.permission_id = p.id
		WHERE up.user_id = ?
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	// Si no hay permisos específicos, devolver permisos básicos
	if len(permissions) == 0 {
		permissions = []string{"read"}
	}

	return permissions, nil
}

func (r *UserRepository) GetRole(ctx context.Context, userID int64) (string, error) {
	var role string
	query := `
		SELECT r.role_name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "user", nil // Rol por defecto
		}
		return "", err
	}

	return role, nil
}

func (r *UserRepository) CreatePasswordResetToken(ctx context.Context, userID int64) (string, time.Time, error) {
	// Generar token aleatorio
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(tokenBytes)

	// El token expira en 1 hora
	expiresAt := time.Now().Add(1 * time.Hour)

	// Insertar token en la base de datos
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, used)
		VALUES (?, ?, ?, false)
	`

	_, err = r.DB.ExecContext(ctx, query, userID, token, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (r *UserRepository) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// Verificar que el token existe y no ha expirado
	var userID int64
	var expiresAt time.Time
	var used bool

	query := `SELECT user_id, expires_at, used FROM password_reset_tokens WHERE token = ?`
	err := r.DB.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewBadRequest("Invalid or expired token")
		}
		return err
	}

	if used {
		return errors.NewBadRequest("Token has already been used")
	}

	if time.Now().After(expiresAt) {
		return errors.NewBadRequest("Token has expired")
	}

	// Hash de la nueva contraseña
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.NewBadRequest("Invalid password format")
	}

	// Iniciar transacción
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Actualizar la contraseña del usuario
	updateUserQuery := `UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?`
	_, err = tx.ExecContext(ctx, updateUserQuery, hashedPassword, userID)
	if err != nil {
		return err
	}

	// Marcar el token como usado
	updateTokenQuery := `UPDATE password_reset_tokens SET used = true WHERE token = ?`
	_, err = tx.ExecContext(ctx, updateTokenQuery, token)
	if err != nil {
		return err
	}

	// Commit de la transacción
	return tx.Commit()
}
