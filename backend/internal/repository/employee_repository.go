package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrEmployeeNotFound = errors.New("сотрудник не найден")

type EmployeeRepository struct {
	db DBTX
}

func NewEmployeeRepository(
	db DBTX,
) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Create(
	ctx context.Context,
	employee models.Employee,
) error {
	query := `
		INSERT INTO employees (
			first_name,
			last_name,
			phone,
			email,
			position_id
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		employee.FirstName,
		employee.LastName,
		employee.Phone,
		employee.Email,
		employee.PositionID,
	)
	if err != nil {
		return fmt.Errorf(
			"создание сотрудника: %w",
			err,
		)
	}

	return nil
}

func (r *EmployeeRepository) ListActive(
	ctx context.Context,
	search string,
) ([]models.EmployeeListItem, error) {
	return r.listByStatus(
		ctx,
		search,
		true,
	)
}

func (r *EmployeeRepository) ListArchived(
	ctx context.Context,
	search string,
) ([]models.EmployeeListItem, error) {
	return r.listByStatus(
		ctx,
		search,
		false,
	)
}

func (r *EmployeeRepository) listByStatus(
	ctx context.Context,
	search string,
	isActive bool,
) ([]models.EmployeeListItem, error) {
	query := `
		SELECT
			e.id,
			e.first_name,
			e.last_name,
			e.phone,
			e.email,
			p.name,
			p.code,
			e.is_active
		FROM employees e
		JOIN employee_positions p
			ON p.id = e.position_id
		WHERE e.is_active = $1
		  AND (
				$2 = ''
				OR e.first_name ILIKE '%' || $2 || '%'
				OR e.last_name ILIKE '%' || $2 || '%'
				OR e.phone ILIKE '%' || $2 || '%'
				OR COALESCE(e.email, '') ILIKE '%' || $2 || '%'
				OR p.name ILIKE '%' || $2 || '%'
				OR p.code ILIKE '%' || $2 || '%'
		  )
		ORDER BY
			e.last_name,
			e.first_name
	`

	rows, err := r.db.Query(
		ctx,
		query,
		isActive,
		search,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получение списка сотрудников: %w",
			err,
		)
	}

	defer rows.Close()

	var employees []models.EmployeeListItem

	for rows.Next() {
		var employee models.EmployeeListItem

		err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.Phone,
			&employee.Email,
			&employee.PositionName,
			&employee.PositionCode,
			&employee.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"сканирование сотрудника: %w",
				err,
			)
		}

		employees = append(
			employees,
			employee,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"итерация по сотрудникам: %w",
			err,
		)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Employee, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone,
			email,
			position_id,
			is_active,
			created_at,
			updated_at
		FROM employees
		WHERE id = $1
	`

	var employee models.Employee

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&employee.Phone,
		&employee.Email,
		&employee.PositionID,
		&employee.IsActive,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}

		return nil, fmt.Errorf(
			"получение сотрудника по id: %w",
			err,
		)
	}

	return &employee, nil
}

func (r *EmployeeRepository) Update(
	ctx context.Context,
	employee models.Employee,
) error {
	query := `
		UPDATE employees
		SET
			first_name = $1,
			last_name = $2,
			phone = $3,
			email = $4,
			position_id = $5,
			updated_at = NOW()
		WHERE id = $6
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		employee.FirstName,
		employee.LastName,
		employee.Phone,
		employee.Email,
		employee.PositionID,
		employee.ID,
	)
	if err != nil {
		return fmt.Errorf(
			"обновление сотрудника: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}

func (r *EmployeeRepository) SetActive(
	ctx context.Context,
	id int64,
	isActive bool,
) error {
	query := `
		UPDATE employees
		SET
			is_active = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		isActive,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"изменение статуса сотрудника: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}
