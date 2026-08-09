package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrItemNotFound = errors.New("складская позиция не найдена")

type WarehouseRepository struct {
	db DBTX
}

func NewWarehouseRepository(db DBTX) *WarehouseRepository {
	return &WarehouseRepository{
		db: db,
	}
}
func (r *WarehouseRepository) List(
	ctx context.Context,
	search string,
) ([]models.WarehouseItem, error) {
	query := `
		SELECT
			id,
			name,
			sku,
			manufacturer,
			purchase_price_cents,
			sale_price_cents,
			quantity,
			min_quantity,
			location,
			note,
			is_active,
			created_at,
			updated_at
		FROM warehouse_items
		WHERE is_active = TRUE
		AND (
				$1 = ''
				OR name ILIKE '%' || $1 || '%'
				OR sku ILIKE '%' || $1 || '%'
				OR manufacturer ILIKE '%' || $1 || '%'
		)
		ORDER BY name, sku
	`
	rows, err := r.db.Query(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf(
			"получение списка: %w",
			err,
		)
	}
	defer rows.Close()

	var items []models.WarehouseItem

	for rows.Next() {

		var item models.WarehouseItem

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.SKU,
			&item.Manufacturer,
			&item.PurchasePriceCents,
			&item.SalePriceCents,
			&item.Quantity,
			&item.MinQuantity,
			&item.Location,
			&item.Note,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование складской позиции: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"итерация по складу: %w",
			err,
		)
	}
	return items, nil
}
func (r *WarehouseRepository) Create(ctx context.Context, item models.WarehouseItem) error {
	query := `
		INSERT INTO warehouse_items (
		name,
		sku,
		manufacturer,
		purchase_price_cents,
		sale_price_cents,
		quantity,
		min_quantity,
		location,
		note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		item.Name,
		item.SKU,
		item.Manufacturer,
		item.PurchasePriceCents,
		item.SalePriceCents,
		item.Quantity,
		item.MinQuantity,
		item.Location,
		item.Note,
	)
	if err != nil {
		return fmt.Errorf("создание складской позиции: %w", err)
	}
	return nil
}
func (r *WarehouseRepository) GetByID(ctx context.Context, id int64) (*models.WarehouseItem, error) {
	query := `
		SELECT
			id,
			name,
			sku,
			manufacturer,
			purchase_price_cents,
			sale_price_cents,
			quantity,
			min_quantity,
			location,
			note,
			is_active,
			created_at,
			updated_at
		FROM warehouse_items
		WHERE id = $1
	`
	var item models.WarehouseItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.SKU,
		&item.Manufacturer,
		&item.PurchasePriceCents,
		&item.SalePriceCents,
		&item.Quantity,
		&item.MinQuantity,
		&item.Location,
		&item.Note,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}

		return nil, fmt.Errorf("получение складской позиции по id: %w", err)
	}

	return &item, nil
}

func (r *WarehouseRepository) Update(
	ctx context.Context,
	item models.WarehouseItem,
) error {
	query := `
		UPDATE warehouse_items
		SET
			name = $1,
			sku = $2,
			manufacturer = $3,
			purchase_price_cents = $4,
			sale_price_cents = $5,
			quantity = $6,
			min_quantity = $7,
			location = $8,
			note = $9,
			updated_at = NOW()
		WHERE id = $10
	`
	commandTag, err := r.db.Exec(
		ctx,
		query,
		item.Name,
		item.SKU,
		item.Manufacturer,
		item.PurchasePriceCents,
		item.SalePriceCents,
		item.Quantity,
		item.MinQuantity,
		item.Location,
		item.Note,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление складской позиции: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *WarehouseRepository) SetActive(ctx context.Context, id int64, isActive bool) error {

	query := `
		UPDATE warehouse_items
		SET
			is_active = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	commandTag, err := r.db.Exec(ctx, query, isActive, id)
	if err != nil {
		return fmt.Errorf("изменение статуса складской позиции: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil

}
func (r *WarehouseRepository) ListArchived(
	ctx context.Context,
	search string,
) ([]models.WarehouseItem, error) {
	query := `
		SELECT
			id,
			name,
			sku,
			manufacturer,
			purchase_price_cents,
			sale_price_cents,
			quantity,
			min_quantity,
			location,
			note,
			is_active,
			created_at,
			updated_at
		FROM warehouse_items
		WHERE is_active = FALSE
		AND (
			$1 = ''
			OR name ILIKE '%' || $1 || '%'
			OR sku ILIKE '%' || $1 || '%'
			OR manufacturer ILIKE '%' || $1 || '%'
		)
		ORDER BY name, sku
	`

	rows, err := r.db.Query(
		ctx,
		query,
		search,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получение архива складских позиций: %w",
			err,
		)
	}
	defer rows.Close()

	var items []models.WarehouseItem

	for rows.Next() {
		var item models.WarehouseItem

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.SKU,
			&item.Manufacturer,
			&item.PurchasePriceCents,
			&item.SalePriceCents,
			&item.Quantity,
			&item.MinQuantity,
			&item.Location,
			&item.Note,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"сканирование архивной складской позиции: %w",
				err,
			)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"итерация по архиву складских позиций: %w",
			err,
		)
	}

	return items, nil
}
