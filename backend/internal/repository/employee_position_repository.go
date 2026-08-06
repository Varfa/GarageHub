package repository

import (
	"context"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
)

type EmployeePositionRepository struct {
	db DBTX
}

func NewEmployeePositionRepository(db DBTX) *EmployeePositionRepository {
	return &EmployeePositionRepository{
		db: db,
	}
}
func (r *EmployeePositionRepository) ListActive(ctx context.Context) ([]models.EmployeePosition, error) {

	query :=
		`SELECT
		id,
		name,
		description,
		is_active,
		created_at,
		updated_at
		FROM employee_positions
		WHERE is_active = true
		ORDER BY name
`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("получение employee positions: %w", err)
	}

	defer rows.Close()
	var positions []models.EmployeePosition

	for rows.Next() {
		var position models.EmployeePosition
		err := rows.Scan(
			&position.ID,
			&position.Name,
			&position.Description,
			&position.IsActive,
			&position.CreatedAt,
			&position.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("получение employee positions: %w", err)
		}
		positions = append(positions, position)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}
	return positions, nil
}
