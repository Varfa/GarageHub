package service

import (
	"context"
	"errors"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var ErrOrderEmployeeAlreadyAssigned = repository.ErrOrderEmployeeAlreadyAssigned

var ErrOrderEmployeeInvalidOrderID = errors.New(
	"order_employee.invalid_order_id",
)

var ErrOrderEmployeeInvalidEmployeeID = errors.New(
	"order_employee.invalid_employee_id",
)

type OrderEmployeeService struct {
	repo *repository.OrderEmployeeRepository
}

func NewOrderEmployeeService(
	repo *repository.OrderEmployeeRepository,
) *OrderEmployeeService {
	return &OrderEmployeeService{
		repo: repo,
	}
}

func (s *OrderEmployeeService) Assign(
	ctx context.Context,
	orderID int64,
	employeeID int64,
) error {

	if orderID <= 0 {
		return ErrOrderEmployeeInvalidOrderID
	}
	if employeeID <= 0 {
		return ErrOrderEmployeeInvalidEmployeeID
	}
	return s.repo.Assign(
		ctx,
		orderID,
		employeeID,
	)
}
func (s *OrderEmployeeService) ListActiveByOrderID(
	ctx context.Context,
	orderID int64,
) ([]models.OrderEmployeeListItem, error) {
	if orderID <= 0 {
		return nil, ErrOrderEmployeeInvalidOrderID
	}

	return s.repo.ListActiveByOrderID(
		ctx,
		orderID,
	)
}
func (s *OrderEmployeeService) Unassign(
	ctx context.Context,
	orderID int64,
	employeeID int64,
) error {
	if orderID <= 0 {
		return ErrOrderEmployeeInvalidOrderID
	}
	if employeeID <= 0 {
		return ErrOrderEmployeeInvalidEmployeeID
	}
	return s.repo.Unassign(ctx, orderID, employeeID)
}
