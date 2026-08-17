package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrRoleNotFound = errors.New("role not found")

type RoleRepository struct {
	db DBTX
}

func NewRoleRepository(
	db DBTX,
) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) ListActive(
	ctx context.Context,
) ([]models.Role, error) {
	query :=
		`
	SELECT
	id,
	code,
	name,
	description,
	is_system,
	is_active,
	created_at,
	updated_at

	FROM roles
	WHERE is_active = TRUE
	ORDER BY name ASC
	`
	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list active roles: %w",
			err,
		)
	}

	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
			&role.IsActive,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"role: %w",
				err,
			)
		}
		roles = append(
			roles,
			role,
		)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate roles: %w",
			err,
		)
	}
	return roles, nil
}
func (r *RoleRepository) List(
	ctx context.Context,
) ([]models.RoleListItem, error) {
	query := `
		SELECT
			r.id,
			r.code,
			r.name,
			COALESCE(r.description, ''),
			r.is_system,
			r.is_active,
			COUNT(rp.permission_id)
		FROM roles r

		LEFT JOIN role_permissions rp
			ON rp.role_id = r.id

		GROUP BY
			r.id,
			r.code,
			r.name,
			r.description,
			r.is_system,
			r.is_active

		ORDER BY
			r.is_system DESC,
			r.name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list roles: %w",
			err,
		)
	}
	defer rows.Close()

	var roles []models.RoleListItem

	for rows.Next() {
		var role models.RoleListItem

		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
			&role.IsActive,
			&role.PermissionsCount,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan role: %w",
				err,
			)
		}

		roles = append(
			roles,
			role,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate roles: %w",
			err,
		)
	}

	return roles, nil
}
func (r *RoleRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.RoleDetails, error) {
	query := `
		SELECT
			id,
			code,
			name,
			COALESCE(description, ''),
			is_system,
			is_active
		FROM roles
		WHERE id = $1
	`

	var role models.RoleDetails

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&role.ID,
		&role.Code,
		&role.Name,
		&role.Description,
		&role.IsSystem,
		&role.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}

		return nil, fmt.Errorf(
			"get role by id: %w",
			err,
		)
	}

	return &role, nil
}
func (r *RoleRepository) ListPermissions(
	ctx context.Context,
	roleID int64,
) ([]models.RolePermissionItem, error) {
	query := `
		SELECT
			p.id,
			p.code,
			p.module,
			p.name,
			COALESCE(p.description, ''),
			CASE
				WHEN rp.permission_id IS NOT NULL THEN TRUE
				ELSE FALSE
			END AS assigned
		FROM permissions p

		LEFT JOIN role_permissions rp
			ON rp.permission_id = p.id
			AND rp.role_id = $1

		ORDER BY
			p.module ASC,
			p.code ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		roleID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list role permissions: %w",
			err,
		)
	}
	defer rows.Close()

	var permissions []models.RolePermissionItem

	for rows.Next() {
		var permission models.RolePermissionItem

		err := rows.Scan(
			&permission.ID,
			&permission.Code,
			&permission.Module,
			&permission.Name,
			&permission.Description,
			&permission.Assigned,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan role permission: %w",
				err,
			)
		}

		permissions = append(
			permissions,
			permission,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate role permissions: %w",
			err,
		)
	}

	return permissions, nil
}

func (r *RoleRepository) UpdatePermissions(
	ctx context.Context,
	roleID int64,
	permissionIDs []int64,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"begin update role permissions transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		`
			DELETE FROM role_permissions
			WHERE role_id = $1
		`,
		roleID,
	)
	if err != nil {
		return fmt.Errorf(
			"delete role permissions: %w",
			err,
		)
	}

	for _, permissionID := range permissionIDs {
		_, err = tx.Exec(
			ctx,
			`
				INSERT INTO role_permissions (
					role_id,
					permission_id
				)
				VALUES ($1, $2)
			`,
			roleID,
			permissionID,
		)
		if err != nil {
			return fmt.Errorf(
				"insert role permission: %w",
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"commit update role permissions: %w",
			err,
		)
	}

	return nil
}
func (r *RoleRepository) Create(
	ctx context.Context,
	code string,
	name string,
	description string,
) (int64, error) {
	query := `
		INSERT INTO roles (
			code,
			name,
			description,
			is_system,
			is_active
		)
		VALUES (
			$1,
			$2,
			$3,
			FALSE,
			TRUE
		)
		RETURNING id
	`

	var roleID int64

	err := r.db.QueryRow(
		ctx,
		query,
		code,
		name,
		description,
	).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf(
			"create role: %w",
			err,
		)
	}

	return roleID, nil
}
func (r *RoleRepository) Update(
	ctx context.Context,
	roleID int64,
	name string,
	description string,
	isActive bool,
) error {
	query := `
		UPDATE roles
		SET
			name = $2,
			description = $3,
			is_active = $4
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		roleID,
		name,
		description,
		isActive,
	)
	if err != nil {
		return fmt.Errorf(
			"update role: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrRoleNotFound
	}

	return nil
}
func (r *RoleRepository) IsAssignedToUsers(
	ctx context.Context,
	roleID int64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE role_id = $1
		)
	`

	var assigned bool

	err := r.db.QueryRow(
		ctx,
		query,
		roleID,
	).Scan(&assigned)
	if err != nil {
		return false, fmt.Errorf(
			"check role users: %w",
			err,
		)
	}

	return assigned, nil
}
func (r *RoleRepository) Delete(
	ctx context.Context,
	roleID int64,
) error {
	query := `
		DELETE FROM roles
		WHERE id = $1
		  AND is_system = FALSE
	`

	result, err := r.db.Exec(
		ctx,
		query,
		roleID,
	)
	if err != nil {
		return fmt.Errorf(
			"delete role: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrRoleNotFound
	}

	return nil
}
