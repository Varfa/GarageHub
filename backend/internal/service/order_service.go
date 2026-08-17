package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var ErrOrderInvalidStatus = errors.New(
	"order.invalid_status",
)
var ErrOrderNotFound = repository.ErrOrderNotFound
var ErrOrderInvalidID = errors.New(
	"order.invalid_id",
)

var ErrOrderClientRequired = errors.New(
	"order.client_required",
)

var ErrOrderCarRequired = errors.New(
	"order.car_required",
)

var ErrOrderComplaintRequired = errors.New(
	"order.complaint_required",
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(
	ctx context.Context,
	order models.Order,
) (int64, error) {
	order.Complaint = strings.TrimSpace(
		order.Complaint)

	if order.ClientID <= 0 {
		return 0, ErrOrderClientRequired
	}

	if order.CarID <= 0 {
		return 0, ErrOrderCarRequired
	}

	if order.Complaint == "" {
		return 0, ErrOrderComplaintRequired
	}

	order.Status = "new"

	return s.repo.Create(
		ctx,
		order,
	)
}
func (s *OrderService) List(
	ctx context.Context,
	search string,
) ([]models.OrderListItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.List(
		ctx,
		search,
	)
}
func (s *OrderService) ListClosed(
	ctx context.Context,
	search string,
) ([]models.OrderListItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.ListClosed(
		ctx,
		search,
	)
}
func (s *OrderService) GetByID(
	ctx context.Context,
	id int,
) (*models.Order, error) {

	if id <= 0 {
		return nil, ErrOrderInvalidID
	}
	return s.repo.GetByID(
		ctx,
		id,
	)
}
func (s *OrderService) UpdateStatus(
	ctx context.Context,
	id int,
	status string,
) error {
	if id <= 0 {
		return ErrOrderInvalidID
	}

	switch status {
	case "new",
		"diagnostics",
		"waiting_approval",
		"in_progress",
		"completed",
		"cancelled":
	default:
		return ErrOrderInvalidStatus
	}

	return s.repo.UpdateStatus(
		ctx,
		id,
		status,
	)
}
func (s *OrderService) Restore(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return ErrOrderInvalidID
	}
	return s.repo.UpdateStatus(
		ctx, id, "new",
	)
}
func (s *OrderService) Delete(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return ErrOrderInvalidID
	}
	return s.repo.Delete(ctx, id)
}
