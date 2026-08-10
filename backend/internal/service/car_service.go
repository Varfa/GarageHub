package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var (
	ErrCarInvalidID = errors.New(
		"car.invalid_id",
	)

	ErrCarClientInvalidID = errors.New(
		"car.client_invalid_id",
	)

	ErrCarOwnerInvalidID = errors.New(
		"car.owner_invalid_id",
	)

	ErrCarBrandRequired = errors.New(
		"car.validation.brand_required",
	)

	ErrCarModelRequired = errors.New(
		"car.validation.model_required",
	)

	ErrCarPlateRequired = errors.New(
		"car.validation.plate_required",
	)

	ErrCarClientRequired = errors.New(
		"car.validation.client_required",
	)

	ErrCarAlreadyExists = errors.New(
		"car.already_exists",
	)
)

type CarService struct {
	repo *repository.CarRepository
}

func NewCarService(
	repo *repository.CarRepository,
) *CarService {
	return &CarService{
		repo: repo,
	}
}

func (s *CarService) Create(
	ctx context.Context,
	car models.Car,
) error {
	car = normalizeCar(car)

	if err := validateCar(car); err != nil {
		return err
	}

	if car.ClientID <= 0 {
		return ErrCarClientRequired
	}

	exists, err := s.repo.ExistsByPlateOrVIN(
		ctx,
		car.PlateNumber,
		car.VIN,
	)
	if err != nil {
		return err
	}

	if exists {
		return ErrCarAlreadyExists
	}

	return s.repo.Create(
		ctx,
		car,
	)
}

func (s *CarService) ListByClientID(
	ctx context.Context,
	clientID int,
) ([]models.Car, error) {
	if clientID <= 0 {
		return nil, ErrCarClientInvalidID
	}

	return s.repo.ListByClientID(
		ctx,
		clientID,
	)
}

func (s *CarService) List(
	ctx context.Context,
	search string,
) ([]models.Car, error) {
	search = strings.TrimSpace(search)

	return s.repo.List(
		ctx,
		search,
	)
}

func (s *CarService) GetByID(
	ctx context.Context,
	id int,
) (*models.Car, error) {
	if id <= 0 {
		return nil, ErrCarInvalidID
	}

	return s.repo.GetByID(
		ctx,
		id,
	)
}

func (s *CarService) Update(
	ctx context.Context,
	car models.Car,
) error {
	if car.ID <= 0 {
		return ErrCarInvalidID
	}

	car = normalizeCar(car)

	if err := validateCar(car); err != nil {
		return err
	}

	return s.repo.Update(
		ctx,
		car,
	)
}

func (s *CarService) Delete(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return ErrCarInvalidID
	}

	return s.repo.Delete(
		ctx,
		id,
	)
}

func (s *CarService) ChangeOwner(
	ctx context.Context,
	carID int,
	clientID int,
) error {
	if carID <= 0 {
		return ErrCarInvalidID
	}

	if clientID <= 0 {
		return ErrCarOwnerInvalidID
	}

	return s.repo.ChangeOwner(
		ctx,
		carID,
		clientID,
	)
}

func normalizeCar(
	car models.Car,
) models.Car {
	car.Brand = strings.TrimSpace(
		car.Brand,
	)

	car.Model = strings.TrimSpace(
		car.Model,
	)

	car.PlateNumber = strings.TrimSpace(
		car.PlateNumber,
	)

	car.VIN = strings.TrimSpace(
		car.VIN,
	)

	car.Engine = strings.TrimSpace(
		car.Engine,
	)

	car.Color = strings.TrimSpace(
		car.Color,
	)

	car.Note = strings.TrimSpace(
		car.Note,
	)

	return car
}

func validateCar(
	car models.Car,
) error {
	if car.Brand == "" {
		return ErrCarBrandRequired
	}

	if car.Model == "" {
		return ErrCarModelRequired
	}

	if car.PlateNumber == "" {
		return ErrCarPlateRequired
	}

	return nil
}
