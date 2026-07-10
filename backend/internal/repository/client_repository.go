package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrClientNotFound = errors.New("клиент не найден")

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
	if err != nil {
		return fmt.Errorf("создание клиента: %w", err)
	}

	return nil
}

func (r *ClientRepository) List(ctx context.Context, search string) ([]models.Client, error) {
	query := `
		SELECT
			id,
			number,
			name,
			phone,
			email,
			address,
			note,
			last_visit_at,
			created_at,
			updated_at
		FROM clients
		WHERE $1 = ''
			OR name ILIKE '%' || $1 || '%'
			OR phone ILIKE '%' || $1 || '%'
			OR email ILIKE '%' || $1 || '%'
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("получение списка клиентов: %w", err)
	}
	defer rows.Close()

	var clients []models.Client

	for rows.Next() {
		var client models.Client

		err := rows.Scan(
			&client.ID,
			&client.Number,
			&client.Name,
			&client.Phone,
			&client.Email,
			&client.Address,
			&client.Note,
			&client.LastVisitAt,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("сканирование клиента: %w", err)
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("итерация по клиентам: %w", err)
	}

	return clients, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id int) (*models.Client, error) {
	query := `
		SELECT
			id,
			number,
			name,
			phone,
			email,
			address,
			note,
			last_visit_at,
			created_at,
			updated_at
		FROM clients
		WHERE id = $1
	`

	var client models.Client

	err := r.db.QueryRow(ctx, query, id).Scan(
		&client.ID,
		&client.Number,
		&client.Name,
		&client.Phone,
		&client.Email,
		&client.Address,
		&client.Note,
		&client.LastVisitAt,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClientNotFound
		}

		return nil, fmt.Errorf("получение клиента по id: %w", err)
	}

	return &client, nil
}

func (r *ClientRepository) Update(ctx context.Context, client models.Client) error {
	query := `
		UPDATE clients
		SET
			name = $1,
			phone = $2,
			email = $3,
			address = $4,
			note = $5,
			updated_at = NOW()
		WHERE id = $6
	`

	_, err := r.db.Exec(
		ctx,
		query,
		client.Name,
		client.Phone,
		client.Email,
		client.Address,
		client.Note,
		client.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление клиента: %w", err)
	}

	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM clients WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление клиента: %w", err)
	}

	return nil
}
