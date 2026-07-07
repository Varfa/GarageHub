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
	return s.repo.Create(ctx, client)
}
func (s *ClientService) List(ctx context.Context) ([]models.Client, error) {
	return s.repo.List(ctx)
}

func (s *ClientService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
