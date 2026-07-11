package service

import (
	"context"
	"errors"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

type CarService struct {
	repo *repository.CarRepository
}

func NewCarService(repo *repository.CarRepository) *CarService {
	return &CarService{
		repo: repo,
	}
}

func (s *CarService) Create(ctx context.Context, car models.Car) error {
	if car.Brand == "" {
		return errors.New("необходимо указать бренд автомобиля")
	}

	if car.Model == "" {
		return errors.New("необходимо указать модель автомобиля")
	}
	if car.PlateNumber == "" {
		return errors.New("необходимо указать номер автомобиля")
	}
	if car.ClientID <= 0 {
		return errors.New("необходимо выбрать клиента")
	}

	return s.repo.Create(ctx, car)
}

func (s *CarService) ListByClientID(
	ctx context.Context,
	clientID int,
) ([]models.Car, error) {

	if clientID <= 0 {
		return nil, errors.New("некорректный id клиента")
	}

	return s.repo.ListByClientID(ctx, clientID)
}
func (s *CarService) List(
	ctx context.Context,
	search string,
) ([]models.Car, error) {
	return s.repo.List(ctx, search)
}

func (s *CarService) GetByID(ctx context.Context, id int) (*models.Car, error) {
	if id <= 0 {
		return nil, errors.New("некорректный id автомобиля")
	}

	return s.repo.GetByID(ctx, id)
}
func (s *CarService) Update(ctx context.Context, car models.Car) error {
	if car.ID <= 0 {
		return errors.New("некорректный id автомобиля")
	}

	if car.Brand == "" {
		return errors.New("необходимо указать бренд автомобиля")
	}

	if car.Model == "" {
		return errors.New("необходимо указать модель автомобиля")
	}
	if car.PlateNumber == "" {
		return errors.New("необходимо указать номер автомобиля")
	}
	return s.repo.Update(ctx, car)
}

func (s *CarService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный id автомобиля")
	}

	return s.repo.Delete(ctx, id)
}
