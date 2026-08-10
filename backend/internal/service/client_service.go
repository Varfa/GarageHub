package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var (
	ErrClientInvalidID = errors.New(
		"client.invalid_id",
	)

	ErrClientNameRequired = errors.New(
		"client.validation.name_required",
	)

	ErrClientPhoneRequired = errors.New(
		"client.validation.phone_required",
	)

	ErrClientAlreadyExists = errors.New(
		"client.already_exists",
	)

	ErrClientNotFound = repository.ErrClientNotFound

	ErrClientHasCars = repository.ErrClientHasCars
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(
	repo *repository.ClientRepository,
) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (s *ClientService) Create(
	ctx context.Context,
	client models.Client,
) error {
	client = normalizeClient(client)

	if err := validateClient(client); err != nil {
		return err
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
		return ErrClientAlreadyExists
	}

	return s.repo.Create(
		ctx,
		client,
	)
}

func (s *ClientService) List(
	ctx context.Context,
	search string,
) ([]models.Client, error) {
	search = strings.TrimSpace(search)

	return s.repo.List(
		ctx,
		search,
	)
}

func (s *ClientService) ListWithCarsCount(
	ctx context.Context,
	search string,
) ([]models.ClientListItem, error) {
	search = strings.TrimSpace(search)

	return s.repo.ListWithCarsCount(
		ctx,
		search,
	)
}

func (s *ClientService) GetByID(
	ctx context.Context,
	id int,
) (*models.Client, error) {
	if id <= 0 {
		return nil, ErrClientInvalidID
	}

	return s.repo.GetByID(
		ctx,
		id,
	)
}

func (s *ClientService) Update(
	ctx context.Context,
	client models.Client,
) error {
	if client.ID <= 0 {
		return ErrClientInvalidID
	}

	client = normalizeClient(client)

	if err := validateClient(client); err != nil {
		return err
	}

	return s.repo.Update(
		ctx,
		client,
	)
}

func (s *ClientService) Delete(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return ErrClientInvalidID
	}

	return s.repo.Delete(
		ctx,
		id,
	)
}

func normalizeClient(
	client models.Client,
) models.Client {
	client.Name = strings.TrimSpace(
		client.Name,
	)

	client.Phone = strings.TrimSpace(
		client.Phone,
	)

	client.Email = strings.TrimSpace(
		client.Email,
	)

	client.Address = strings.TrimSpace(
		client.Address,
	)

	client.Note = strings.TrimSpace(
		client.Note,
	)

	return client
}

func validateClient(
	client models.Client,
) error {
	if client.Name == "" {
		return ErrClientNameRequired
	}

	if client.Phone == "" {
		return ErrClientPhoneRequired
	}

	return nil
}
