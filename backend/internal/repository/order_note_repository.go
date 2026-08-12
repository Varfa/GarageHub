package repository

import (
	"context"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
)

type OrderNoteRepository struct {
	db DBTX
}

func NewOrderNoteRepository(
	db DBTX,
) *OrderNoteRepository {
	return &OrderNoteRepository{
		db: db,
	}
}
func (r *OrderNoteRepository) Create(
	ctx context.Context,
	note models.OrderNote,
) error {
	query := `
	INSERT INTO order_notes (
		order_id,
		employee_id,
		text
	)
	VALUES (
		$1, $2, $3
	)
`
	_, err := r.db.Exec(ctx, query, note.OrderID, note.EmployeeID, note.Text)
	if err != nil {
		return fmt.Errorf("создание заметки заказа: %w", err)
	}
	return nil
}
func (r *OrderNoteRepository) ListByOrderID(
	ctx context.Context,
	orderID int64,
) ([]models.OrderNote, error) {
	query := `
		SELECT
			id,
			order_id,
			employee_id,
			text,
			created_at
		FROM order_notes
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получение заметок заказа: %w",
			err,
		)
	}
	defer rows.Close()

	var notes []models.OrderNote

	for rows.Next() {
		var note models.OrderNote

		err := rows.Scan(
			&note.ID,
			&note.OrderID,
			&note.EmployeeID,
			&note.Text,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("чтение заметки заказа: %w", err)
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("перебор заметок заказа: %w", err)
	}

	return notes, nil
}
