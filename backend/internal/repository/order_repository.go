package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db DBTX
}

func NewOrderRepository(db DBTX) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	order models.Order,
) (int64, error) {
	query := `
		INSERT INTO orders (
			client_id,
			car_id,
			complaint,
			diagnosis,
			note,
			status,
			estimated_cost_cents,
			final_cost_cents,
			started_at,
			completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var orderID int64

	err := r.db.QueryRow(
		ctx,
		query,
		order.ClientID,
		order.CarID,
		order.Complaint,
		order.Diagnosis,
		order.Note,
		order.Status,
		order.EstimatedCostCents,
		order.FinalCostCents,
		order.StartedAt,
		order.CompletedAt,
	).Scan(
		&orderID,
	)

	if err != nil {
		return 0, fmt.Errorf(
			"создание заказа: %w",
			err,
		)
	}

	return orderID, nil
}

func (r *OrderRepository) List(
	ctx context.Context,
	search string,
) ([]models.OrderListItem, error) {
	query := `
		SELECT
			o.id,
			o.client_id,
			c.name,
			o.car_id,
			car.brand,
			car.model,
			car.plate_number,

			COALESCE(
				string_agg(
					DISTINCT e.first_name || ' ' || e.last_name,
					', '
				),
				''
			) AS assigned_employees,

			o.complaint,
			o.status,
			o.estimated_cost_cents,
			o.final_cost_cents,
			o.created_at

		FROM orders o

		JOIN clients c
			ON c.id = o.client_id

		JOIN cars car
			ON car.id = o.car_id

		LEFT JOIN order_employees oe
			ON oe.order_id = o.id
			AND oe.unassigned_at IS NULL

		LEFT JOIN employees e
			ON e.id = oe.employee_id

		WHERE o.status NOT IN ('completed', 'cancelled')
		AND (
	$1 = ''
	OR CAST(o.id AS TEXT) ILIKE '%' || $1 || '%'
	OR c.name ILIKE '%' || $1 || '%'
	OR car.plate_number ILIKE '%' || $1 || '%'
	OR car.vin ILIKE '%' || $1 || '%'
	OR car.brand ILIKE '%' || $1 || '%'
	OR car.model ILIKE '%' || $1 || '%'
)

		GROUP BY
			o.id,
			o.client_id,
			c.name,
			o.car_id,
			car.brand,
			car.model,
			car.plate_number,
			o.complaint,
			o.status,
			o.estimated_cost_cents,
			o.final_cost_cents,
			o.created_at

		ORDER BY o.id DESC

	`

	rows, err := r.db.Query(
		ctx,
		query,
		search,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list orders: %w",
			err,
		)
	}
	defer rows.Close()

	var orders []models.OrderListItem

	for rows.Next() {
		var order models.OrderListItem

		err := rows.Scan(
			&order.ID,
			&order.ClientID,
			&order.ClientName,
			&order.CarID,
			&order.CarBrand,
			&order.CarModel,
			&order.CarPlateNumber,
			&order.AssignedEmployees,
			&order.Complaint,
			&order.Status,
			&order.EstimatedCostCents,
			&order.FinalCostCents,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan order: %w",
				err,
			)
		}

		orders = append(
			orders,
			order,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate orders: %w",
			err,
		)
	}

	return orders, nil
}

var ErrOrderNotFound = errors.New(
	"order not found",
)

func (r *OrderRepository) ListClosed(
	ctx context.Context,
	search string,
) ([]models.OrderListItem, error) {
	query := `
		SELECT
			o.id,
			o.client_id,
			c.name,
			o.car_id,
			car.brand,
			car.model,
			car.plate_number,

			COALESCE(
				string_agg(
					DISTINCT e.first_name || ' ' || e.last_name,
					', '
				),
				''
			) AS assigned_employees,

			o.complaint,
			o.status,
			o.estimated_cost_cents,
			o.final_cost_cents,
			o.created_at

		FROM orders o

		JOIN clients c
			ON c.id = o.client_id

		JOIN cars car
			ON car.id = o.car_id

		LEFT JOIN order_employees oe
			ON oe.order_id = o.id
			AND oe.unassigned_at IS NULL

		LEFT JOIN employees e
			ON e.id = oe.employee_id

		WHERE o.status IN ('completed', 'cancelled')
	AND (
		$1 = ''
		OR CAST(o.id AS TEXT) ILIKE '%' || $1 || '%'
		OR c.name ILIKE '%' || $1 || '%'
		OR car.plate_number ILIKE '%' || $1 || '%'
		OR car.vin ILIKE '%' || $1 || '%'
		OR car.brand ILIKE '%' || $1 || '%'
		OR car.model ILIKE '%' || $1 || '%'
	)

		GROUP BY
			o.id,
			o.client_id,
			c.name,
			o.car_id,
			car.brand,
			car.model,
			car.plate_number,
			o.complaint,
			o.status,
			o.estimated_cost_cents,
			o.final_cost_cents,
			o.created_at

		ORDER BY o.id DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		search,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list closed orders: %w",
			err,
		)
	}
	defer rows.Close()

	var orders []models.OrderListItem

	for rows.Next() {
		var order models.OrderListItem

		err := rows.Scan(
			&order.ID,
			&order.ClientID,
			&order.ClientName,
			&order.CarID,
			&order.CarBrand,
			&order.CarModel,
			&order.CarPlateNumber,
			&order.AssignedEmployees,
			&order.Complaint,
			&order.Status,
			&order.EstimatedCostCents,
			&order.FinalCostCents,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan closed order: %w",
				err,
			)
		}

		orders = append(
			orders,
			order,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate closed orders: %w",
			err,
		)
	}

	return orders, nil
}

func (r *OrderRepository) GetByID(
	ctx context.Context,
	id int,
) (*models.Order, error) {
	query := `
		SELECT
			id,
			client_id,
			car_id,
			complaint,
			diagnosis,
			note,
			status,
			estimated_cost_cents,
			final_cost_cents,
			created_at,
			started_at,
			completed_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&order.ID,
		&order.ClientID,
		&order.CarID,
		&order.Complaint,
		&order.Diagnosis,
		&order.Note,
		&order.Status,
		&order.EstimatedCostCents,
		&order.FinalCostCents,
		&order.CreatedAt,
		&order.StartedAt,
		&order.CompletedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf(
			"get order by id: %w",
			err,
		)
	}

	return &order, nil
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	id int,
	status string,
) error {
	query := `
		UPDATE orders
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"обновление статуса заказа: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}
func (r *OrderRepository) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE FROM orders
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}
