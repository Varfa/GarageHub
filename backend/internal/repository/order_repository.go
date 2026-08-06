package repository

import (
	"context"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
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
