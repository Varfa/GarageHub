package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
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

	return s.repo.Create(ctx, item)
}

func (s *WarehouseService) List(
	ctx context.Context,
	search string,
) ([]models.WarehouseItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.List(ctx, search)
}

func (s *WarehouseService) GetByID(
	ctx context.Context,
	id int64,
) (*models.WarehouseItem, error) {
	if id <= 0 {
		return nil, errors.New(
			"некорректный id складской позиции",
		)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *WarehouseService) Update(
	ctx context.Context,
	item models.WarehouseItem,
) error {
	if item.ID <= 0 {
		return errors.New(
			"некорректный id складской позиции",
		)
	}

	item = normalizeWarehouseItem(item)

	if err := validateWarehouseItem(item); err != nil {
		return err
	}

	return s.repo.Update(ctx, item)
}

func normalizeWarehouseItem(
	item models.WarehouseItem,
) models.WarehouseItem {
	item.Name = strings.TrimSpace(item.Name)
	item.SKU = strings.TrimSpace(item.SKU)
	item.Manufacturer = strings.TrimSpace(item.Manufacturer)
	item.Location = strings.TrimSpace(item.Location)
	item.Note = strings.TrimSpace(item.Note)

	return item
}

func validateWarehouseItem(
	item models.WarehouseItem,
) error {
	if item.Name == "" {
		return errors.New(
			"необходимо указать название складской позиции",
		)
	}

	if item.SKU == "" {
		return errors.New(
			"необходимо указать артикул",
		)
	}

	if item.PurchasePriceCents < 0 {
		return errors.New(
			"цена закупки не может быть отрицательной",
		)
	}

	if item.SalePriceCents < 0 {
		return errors.New(
			"цена продажи не может быть отрицательной",
		)
	}

	if item.Quantity < 0 {
		return errors.New(
			"количество не может быть отрицательным",
		)
	}

	if item.MinQuantity < 0 {
		return errors.New(
			"минимальное количество не может быть отрицательным",
		)
	}

	return nil
}
func (s *WarehouseService) Archive(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New(
			"некорректный id складской позиции",
		)
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !item.IsActive {
		return errors.New(
			"позиция уже находится в архиве",
		)
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
		return errors.New(
			"некорректный id складской позиции",
		)
	}

	item, err := s.repo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if item.IsActive {
		return errors.New(
			"позиция уже активна",
		)
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
