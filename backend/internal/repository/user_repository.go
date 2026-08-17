package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New(
	"user not found",
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(
	db DBTX,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) HasOwner(
	ctx context.Context,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE is_owner = TRUE
				AND is_active = TRUE
		)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
	).Scan(
		&exists,
	)
	if err != nil {
		return false, fmt.Errorf(
			"проверка наличия owner: %w",
			err,
		)
	}

	return exists, nil
}

func (r *UserRepository) CreateOwner(
	ctx context.Context,
	user models.User,
) error {
	query := `
		INSERT INTO users (
			employee_id,
			role_id,
			email,
			password_hash,
			is_owner,
			is_active
		)
		VALUES (
			NULL,
			NULL,
			$1,
			$2,
			TRUE,
			TRUE
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf(
			"создание owner: %w",
			err,
		)
	}

	return nil
}
func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	var user models.User

	query :=
		`SELECT
		id,
		employee_id,
		role_id,
		email,
		password_hash,
		is_owner,
		is_active,
		last_login_at,
		created_at,
		updated_at

		FROM users
			WHERE LOWER(email) = LOWER($1)


	`

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
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
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}
func (r *UserRepository) CreateUser(
	ctx context.Context,
	user models.User,
) error {
	query := `
	INSERT INTO users (
		employee_id,
		role_id,
		email,
		password_hash,
		is_owner,
		is_active
	)
	VALUES ($1, $2, $3, $4, FALSE, TRUE)
`
	_, err := r.db.Exec(
		ctx,
		query,
		user.EmployeeID,
		user.RoleID,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}
	return nil
}

func (r *UserRepository) HasPermission(
	ctx context.Context,
	userID int64,
	permissionCode string,
) (bool, error) {

	query := `
	SELECT EXISTS (
		SELECT 1
		FROM users u
		JOIN role_permissions rp
			ON rp.role_id = u.role_id
		JOIN permissions p
			ON p.id = rp.permission_id
		WHERE u.id = $1
		  AND p.code = $2
	)
`
	var hasPermission bool

	err := r.db.QueryRow(ctx, query, userID, permissionCode).Scan(&hasPermission)
	if err != nil {
		return false, err

	}

	return hasPermission, nil
}
func (r *UserRepository) List(
	ctx context.Context,
) ([]models.UserListItem, error) {
	query := `
		SELECT
			u.id,
			COALESCE(u.employee_id, 0),
			COALESCE(e.first_name, ''),
			COALESCE(e.last_name, ''),
			COALESCE(ep.name, ''),
			u.email,
			u.role_id,
			COALESCE(r.code, ''),
			COALESCE(r.name, ''),
			u.is_active,
			u.is_owner
		FROM users u

		LEFT JOIN employees e
			ON e.id = u.employee_id

		LEFT JOIN employee_positions ep
			ON ep.id = e.position_id

		LEFT JOIN roles r
			ON r.id = u.role_id

		ORDER BY
			u.is_owner DESC,
			e.last_name ASC,
			e.first_name ASC,
			u.id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list users: %w",
			err,
		)
	}
	defer rows.Close()

	var users []models.UserListItem

	for rows.Next() {
		var user models.UserListItem

		err := rows.Scan(
			&user.ID,
			&user.EmployeeID,
			&user.FirstName,
			&user.LastName,
			&user.PositionName,
			&user.Email,
			&user.RoleID,
			&user.RoleCode,
			&user.RoleName,
			&user.IsActive,
			&user.IsOwner,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan user: %w",
				err,
			)
		}

		users = append(
			users,
			user,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate users: %w",
			err,
		)
	}

	return users, nil
}
func (r *UserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.UserListItem, error) {
	query := `
	SELECT
		u.id,
		COALESCE(u.employee_id, 0),
		COALESCE(e.first_name, ''),
		COALESCE(e.last_name, ''),
		COALESCE(ep.name, ''),
		u.email,
		u.role_id,
		COALESCE(r.code, ''),
		COALESCE(r.name, ''),
		u.is_active,
		u.is_owner
	FROM users u

	LEFT JOIN employees e
		ON e.id = u.employee_id

	LEFT JOIN employee_positions ep
		ON ep.id = e.position_id

	LEFT JOIN roles r
		ON r.id = u.role_id

	WHERE u.id = $1
`
	var user models.UserListItem
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.EmployeeID,
		&user.FirstName,
		&user.LastName,
		&user.PositionName,
		&user.Email,
		&user.RoleID,
		&user.RoleCode,
		&user.RoleName,
		&user.IsActive,
		&user.IsOwner,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf(
			"get user by id: %w",
			err,
		)
	}

	return &user, nil
}
func (r *UserRepository) UpdateRole(
	ctx context.Context,
	userID int64,
	roleID int64,
) error {
	query := `
		UPDATE users
		SET
			role_id = $1,
			updated_at = NOW()
		WHERE id = $2
		  AND is_owner = FALSE
	`

	result, err := r.db.Exec(
		ctx,
		query,
		roleID,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"update user role: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
func (r *UserRepository) SetActive(
	ctx context.Context,
	userID int64,
	isActive bool,
) error {
	query := `
		UPDATE users
		SET
			is_active = $1,
			updated_at = NOW()
		WHERE id = $2
		  AND is_owner = FALSE
	`

	result, err := r.db.Exec(
		ctx,
		query,
		isActive,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"update user active status: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
