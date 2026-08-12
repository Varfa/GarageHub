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

func (r *OrderRepository) Create(ctx context.Context, order models.Order) error {
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
	`

	_, err := r.db.Exec(
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
	)
	if err != nil {
		return fmt.Errorf("создание заказа: %w", err)
	}

	return nil
}

// List возвращает список заказов вместе
// с данными клиента и автомобиля.
func (r *OrderRepository) List(
	ctx context.Context,
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
		ORDER BY o.id DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
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
		return fmt.Errorf("обновление статуса заказа: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}
