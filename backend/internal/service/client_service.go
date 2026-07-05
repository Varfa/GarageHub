package service

import (
	"context"

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
	return s.repo.Create(ctx, client)
}
