package repository

import (
	"context"
	"fmt"

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

func (r *ClientRepository) List(ctx context.Context) ([]models.Client, error) {
	rows, err := r.db.Query(ctx, "SELECT id, number, name, phone, email, address, note, last_visit_at, created_at, updated_at FROM clients")
	if err != nil {
		return nil, fmt.Errorf("получение списка клиентов %w", err)
	}

	defer rows.Close()
	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(
			&c.ID,
			&c.Number,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Address,
			&c.Note,
			&c.LastVisitAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("сканирование клиента: %w", err)
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("итерация по клиентам: %w", err)
	}

	return clients, nil

}
