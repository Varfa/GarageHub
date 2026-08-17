package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderEmployeeRepository struct {
	db DBTX
}

func NewOrderEmployeeRepository(
	db DBTX,
) *OrderEmployeeRepository {
	return &OrderEmployeeRepository{
		db: db,
	}
}

var ErrOrderEmployeeAlreadyAssigned = errors.New(
	"order employee already assigned",
)

func (r *OrderEmployeeRepository) Assign(
	ctx context.Context,
	orderID int64,
	employeeID int64,
) error {
	query := `
		INSERT INTO order_employees (
			order_id,
			employee_id
		)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		orderID,
		employeeID,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" {
			return ErrOrderEmployeeAlreadyAssigned
		}

		return fmt.Errorf(
			"назначение сотрудника на заказ: %w",
			err,
		)
	}

	return nil
}
func (r *OrderEmployeeRepository) ListActiveByOrderID(
	ctx context.Context,
	orderID int64,
) ([]models.OrderEmployeeListItem, error) {
	query := `
	SELECT
		oe.id,
		oe.order_id,
		oe.employee_id,
		e.first_name,
		e.last_name,
		p.name,
		p.code,
		oe.assigned_at,
		oe.unassigned_at
	FROM order_employees oe
	JOIN employees e
		ON e.id = oe.employee_id
	JOIN employee_positions p
		ON p.id = e.position_id
	WHERE oe.order_id = $1
		AND oe.unassigned_at IS NULL
	ORDER BY oe.assigned_at ASC
`

	rows, err := r.db.Query(
		ctx,
		query,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получение назначенных сотрудников заказа: %w",
			err,
		)
	}
	defer rows.Close()

	var employees []models.OrderEmployeeListItem

	for rows.Next() {
		var employee models.OrderEmployeeListItem
		err := rows.Scan(
			&employee.ID,
			&employee.OrderID,
			&employee.EmployeeID,
			&employee.FirstName,
			&employee.LastName,
			&employee.PositionName,
			&employee.PositionCode,
			&employee.AssignedAt,
			&employee.UnassignedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"чтение назначенного сотрудника заказа: %w",
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
			"перебор назначенных сотрудников заказа: %w",
			err,
		)
	}

	return employees, nil
}
func (r *OrderEmployeeRepository) Unassign(
	ctx context.Context,
	orderID int64,
	employeeID int64,
) error {
	query := `
		UPDATE order_employees
		SET unassigned_at = NOW()
		WHERE order_id = $1
			AND employee_id = $2
			AND unassigned_at IS NULL
	`

	result, err := r.db.Exec(
		ctx,
		query,
		orderID,
		employeeID,
	)
	if err != nil {
		return fmt.Errorf(
			"снятие сотрудника с заказа: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"активное назначение сотрудника не найдено",
		)
	}

	return nil
}
