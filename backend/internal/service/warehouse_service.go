package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var (
	ErrWarehouseInvalidID = errors.New(
		"warehouse.invalid_id",
	)

	ErrWarehouseNameRequired = errors.New(
		"warehouse.validation.name_required",
	)

	ErrWarehouseSKURequired = errors.New(
		"warehouse.validation.sku_required",
	)

	ErrWarehousePurchasePriceNegative = errors.New(
		"warehouse.validation.purchase_price_negative",
	)

	ErrWarehouseSalePriceNegative = errors.New(
		"warehouse.validation.sale_price_negative",
	)

	ErrWarehouseQuantityNegative = errors.New(
		"warehouse.validation.quantity_negative",
	)

	ErrWarehouseMinQuantityNegative = errors.New(
		"warehouse.validation.min_quantity_negative",
	)

	ErrWarehouseAlreadyArchived = errors.New(
		"warehouse.already_archived",
	)

	ErrWarehouseAlreadyActive = errors.New(
		"warehouse.already_active",
	)
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(
	repo *repository.WarehouseRepository,
) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

func (s *WarehouseService) Create(
	ctx context.Context,
	item models.WarehouseItem,
) error {
	item = normalizeWarehouseItem(item)

	if err := validateWarehouseItem(item); err != nil {
		return err
	}

	return s.repo.Create(
		ctx,
		item,
	)
}

func (s *WarehouseService) List(
	ctx context.Context,
	search string,
) ([]models.WarehouseItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.List(
		ctx,
		search,
	)
}

func (s *WarehouseService) GetByID(
	ctx context.Context,
	id int64,
) (*models.WarehouseItem, error) {
	if id <= 0 {
		return nil, ErrWarehouseInvalidID
	}

	return s.repo.GetByID(
		ctx,
		id,
	)
}

func (s *WarehouseService) Update(
	ctx context.Context,
	item models.WarehouseItem,
) error {
	if item.ID <= 0 {
		return ErrWarehouseInvalidID
	}

	item = normalizeWarehouseItem(item)

	if err := validateWarehouseItem(item); err != nil {
		return err
	}

	return s.repo.Update(
		ctx,
		item,
	)
}

func (s *WarehouseService) Archive(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrWarehouseInvalidID
	}

	item, err := s.repo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if !item.IsActive {
		return ErrWarehouseAlreadyArchived
	}

	return s.repo.SetActive(
		ctx,
		id,
		false,
	)
}

func (s *WarehouseService) Restore(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrWarehouseInvalidID
	}

	item, err := s.repo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if item.IsActive {
		return ErrWarehouseAlreadyActive
	}

	return s.repo.SetActive(
		ctx,
		id,
		true,
	)
}

func (s *WarehouseService) ListArchived(
	ctx context.Context,
	search string,
) ([]models.WarehouseItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.ListArchived(
		ctx,
		search,
	)
}

func normalizeWarehouseItem(
	item models.WarehouseItem,
) models.WarehouseItem {
	item.Name = strings.TrimSpace(
		item.Name,
	)

	item.SKU = strings.TrimSpace(
		item.SKU,
	)

	item.Manufacturer = strings.TrimSpace(
		item.Manufacturer,
	)

	item.Location = strings.TrimSpace(
		item.Location,
	)

	item.Note = strings.TrimSpace(
		item.Note,
	)

	return item
}

func validateWarehouseItem(
	item models.WarehouseItem,
) error {
	if item.Name == "" {
		return ErrWarehouseNameRequired
	}

	if item.SKU == "" {
		return ErrWarehouseSKURequired
	}

	if item.PurchasePriceCents < 0 {
		return ErrWarehousePurchasePriceNegative
	}

	if item.SalePriceCents < 0 {
		return ErrWarehouseSalePriceNegative
	}

	if item.Quantity < 0 {
		return ErrWarehouseQuantityNegative
	}

	if item.MinQuantity < 0 {
		return ErrWarehouseMinQuantityNegative
	}

	return nil
}
