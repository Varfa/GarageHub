package service

import (
	"context"
	"errors"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (s *ClientService) Create(ctx context.Context, client models.Client) error {
	if client.Name == "" {
		return errors.New("необходимо указать имя")
	}

	if client.Phone == "" {
		return errors.New("необходимо указать номер телефона")
	}

	exists, err := s.repo.ExistsByNameAndPhone(
		ctx,
		client.Name,
		client.Phone,
	)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("клиент с таким именем и телефоном уже существует")
	}

	return s.repo.Create(ctx, client)
}

func (s *ClientService) List(ctx context.Context, search string) ([]models.Client, error) {
	return s.repo.List(ctx, search)
}

func (s *ClientService) GetByID(ctx context.Context, id int) (*models.Client, error) {
	if id <= 0 {
		return nil, errors.New("некорректный id клиента")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ClientService) Update(ctx context.Context, client models.Client) error {
	if client.ID <= 0 {
		return errors.New("некорректный id клиента")
	}

	if client.Name == "" {
		return errors.New("необходимо указать имя")
	}

	if client.Phone == "" {
		return errors.New("необходимо указать номер телефона")
	}

	return s.repo.Update(ctx, client)
}

func (s *ClientService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный id клиента")
	}

	return s.repo.Delete(ctx, id)
}
