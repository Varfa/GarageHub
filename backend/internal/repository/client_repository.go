package repository

import (
	"context"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}
func (r *ClientRepository) Create(ctx context.Context, client models.Client) error {
	query := `
		INSERT INTO clients (
			number,
			name,
			phone,
			email,
			address,
			note,
			last_visit_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		client.Number,
		client.Name,
		client.Phone,
		client.Email,
		client.Address,
		client.Note,
		client.LastVisitAt,
		client.CreatedAt,
		client.UpdatedAt,
	)

	return err
}
