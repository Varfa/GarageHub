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
) error {
	order.Complaint = strings.TrimSpace(
		order.Complaint)

	if order.ClientID <= 0 {
		return ErrOrderClientRequired
	}

	if order.CarID <= 0 {
		return ErrOrderCarRequired
	}

	if order.Complaint == "" {
		return ErrOrderComplaintRequired
	}

	order.Status = "new"

	return s.repo.Create(
		ctx,
		order,
	)
}
func (s *OrderService) List(
	ctx context.Context,
) ([]models.OrderListItem, error) {

	return s.repo.List(
		ctx,
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
		"completed":
	default:
		return ErrOrderInvalidStatus
	}

	return s.repo.UpdateStatus(
		ctx,
		id,
		status,
	)
}
