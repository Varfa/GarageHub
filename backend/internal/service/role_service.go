package service

import (
	"context"
	"errors"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var ErrRoleInUse = errors.New(
	"role is assigned to users",
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(
	repo *repository.RoleRepository,
) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (s *RoleService) ListActive(
	ctx context.Context,
) ([]models.Role, error) {
	return s.repo.ListActive(ctx)
}

func (s *RoleService) List(
	ctx context.Context,
) ([]models.RoleListItem, error) {
	return s.repo.List(ctx)
}

func (s *RoleService) GetByID(
	ctx context.Context,
	id int64,
) (*models.RoleDetails, error) {
	if id <= 0 {
		return nil, repository.ErrRoleNotFound
	}

	return s.repo.GetByID(
		ctx,
		id,
	)
}

func (s *RoleService) ListPermissions(
	ctx context.Context,
	roleID int64,
) ([]models.RolePermissionItem, error) {
	if roleID <= 0 {
		return nil, repository.ErrRoleNotFound
	}

	return s.repo.ListPermissions(
		ctx,
		roleID,
	)
}

func (s *RoleService) UpdatePermissions(
	ctx context.Context,
	roleID int64,
	permissionIDs []int64,
) error {
	if roleID <= 0 {
		return repository.ErrRoleNotFound
	}

	_, err := s.repo.GetByID(
		ctx,
		roleID,
	)
	if err != nil {
		return err
	}

	return s.repo.UpdatePermissions(
		ctx,
		roleID,
		permissionIDs,
	)
}
func (s *RoleService) Create(
	ctx context.Context,
	code string,
	name string,
	description string,
) (int64, error) {
	if code == "" {
		return 0, errors.New("role code is required")
	}

	if name == "" {
		return 0, errors.New("role name is required")
	}

	return s.repo.Create(
		ctx,
		code,
		name,
		description,
	)
}
func (s *RoleService) Update(
	ctx context.Context,
	roleID int64,
	name string,
	description string,
	isActive bool,
) error {
	if roleID <= 0 {
		return repository.ErrRoleNotFound
	}

	role, err := s.repo.GetByID(
		ctx,
		roleID,
	)
	if err != nil {
		return err
	}

	if !role.IsSystem && name == "" {
		return errors.New("role name is required")
	}

	if role.IsSystem {
		name = role.Name
	}

	return s.repo.Update(
		ctx,
		roleID,
		name,
		description,
		isActive,
	)
}
func (s *RoleService) Delete(
	ctx context.Context,
	roleID int64,
) error {
	if roleID <= 0 {
		return repository.ErrRoleNotFound
	}

	role, err := s.repo.GetByID(
		ctx,
		roleID,
	)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("system role cannot be deleted")
	}

	assigned, err := s.repo.IsAssignedToUsers(
		ctx,
		roleID,
	)
	if err != nil {
		return err
	}

	if assigned {
		return ErrRoleInUse
	}

	return s.repo.Delete(
		ctx,
		roleID,
	)
}
