package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrSessionNotFound = errors.New(
	"session not found",
)

type SessionRepository struct {
	db DBTX
}

func NewSessionRepository(
	db DBTX,
) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {
	query := `
	INSERT INTO user_sessions (
		user_id,
		token_hash,
		expires_at
	)
	VALUES ($1, $2, $3)
`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf(
			"create user session: %w",
			err)
	}

	return nil
}
func (r *SessionRepository) Delete(
	ctx context.Context,
	tokenHash string,
) error {
	query := `
		DELETE FROM user_sessions
		WHERE token_hash = $1

		`
	_, err := r.db.Exec(ctx, query, tokenHash)

	if err != nil {
		return fmt.Errorf("delete user session: %w", err)
	}
	return nil
}
func (r *SessionRepository) GetUserByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*models.User, error) {

	var user models.User

	query := `
	SELECT
		u.id,
		u.employee_id,
		u.role_id,
		u.email,
		u.password_hash,
		u.is_owner,
		u.is_active,
		u.last_login_at,
		u.created_at,
		u.updated_at

	FROM users u

	JOIN user_sessions s
		ON s.user_id = u.id

	WHERE s.token_hash = $1
		AND s.expires_at > NOW()
		AND u.is_active = TRUE
`
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&user.ID,
		&user.EmployeeID,
		&user.RoleID,
		&user.Email,
		&user.PasswordHash,
		&user.IsOwner,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("get user by session: %w", err)

	}
	return &user, nil
}
