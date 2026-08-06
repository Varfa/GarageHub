package service

import (
	"context"
	"errors"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
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

	if order.ClientID <= 0 {
		return errors.New("необходимо выбрать клиента")
	}

	if order.CarID <= 0 {
		return errors.New("необходимо выбрать автомобиль")
	}

	if order.Complaint == "" {
		return errors.New("необходимо описать проблему клиента")
	}

	order.Status = "new"

	return s.repo.Create(ctx, order)
}
