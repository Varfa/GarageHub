package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserRoleInvalid = errors.New(
		"user.role_invalid",
	)
	ErrUserInvalidID = errors.New(
		"user.invalid_id",
	)
	ErrUserEmployeeRequired = errors.New(
		"user.employee_required",
	)

	ErrUserRoleRequired = errors.New(
		"user.role_required",
	)
	ErrInvalidCredentials = errors.New(
		"user.invalid_credentials",
	)
	ErrUserEmailRequired = errors.New(
		"user.email_required",
	)

	ErrUserPasswordRequired = errors.New(
		"user.password_required",
	)

	ErrOwnerAlreadyExists = errors.New(
		"user.owner_already_exists",
	)
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(
	repo *repository.UserRepository,
) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) HasOwner(
	ctx context.Context,
) (bool, error) {
	return s.repo.HasOwner(ctx)
}

func (s *UserService) CreateOwner(
	ctx context.Context,
	email string,
	password string,
) error {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" {
		return ErrUserEmailRequired
	}

	if password == "" {
		return ErrUserPasswordRequired
	}

	hasOwner, err := s.repo.HasOwner(ctx)
	if err != nil {
		return err
	}

	if hasOwner {
		return ErrOwnerAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		IsOwner:      true,
		IsActive:     true,
	}

	return s.repo.CreateOwner(
		ctx,
		user,
	)
}

func (s *UserService) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (*models.User, error) {
	email = strings.TrimSpace(email)

	user, err := s.repo.GetByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(
			err,
			repository.ErrUserNotFound,
		) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
func (s *UserService) CreateUser(
	ctx context.Context,
	employeeID int64,
	roleID int64,
	email string,
	password string,
) error {
	email = strings.TrimSpace(email)

	if employeeID <= 0 {
		return ErrUserEmployeeRequired
	}
	if roleID <= 0 {
		return ErrUserRoleRequired
	}
	if email == "" {
		return ErrUserEmailRequired
	}
	if password == "" {
		return ErrUserPasswordRequired
	}
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	user := models.User{
		EmployeeID:   &employeeID,
		RoleID:       &roleID,
		Email:        email,
		PasswordHash: string(passwordHash),
		IsOwner:      false,
		IsActive:     true,
	}
	return s.repo.CreateUser(ctx, user)
}
func (s *UserService) HasPermission(
	ctx context.Context,
	userID int64,
	permissionCode string,
) (bool, error) {

	return s.repo.HasPermission(
		ctx,
		userID,
		permissionCode,
	)
}
func (s *UserService) List(
	ctx context.Context,
) ([]models.UserListItem, error) {
	return s.repo.List(ctx)

}
func (s *UserService) GetByID(
	ctx context.Context,
	id int64,
) (*models.UserListItem, error) {
	if id <= 0 {
		return nil, ErrUserInvalidID
	}

	return s.repo.GetByID(
		ctx,
		id,
	)
}
func (s *UserService) UpdateRole(
	ctx context.Context,
	userID int64,
	roleID int64,
) error {
	if userID <= 0 {
		return ErrUserInvalidID

	}
	if roleID <= 0 {
		return ErrUserRoleInvalid
	}
	user, err := s.repo.GetByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	if user.IsOwner {
		return ErrUserRoleInvalid
	}
	return s.repo.UpdateRole(
		ctx,
		userID,
		roleID,
	)

}
func (s *UserService) SetActive(
	ctx context.Context,
	userID int64,
	isActive bool,
) error {
	if userID <= 0 {
		return ErrUserInvalidID
	}

	user, err := s.repo.GetByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	if user.IsOwner {
		return ErrUserRoleInvalid
	}

	return s.repo.SetActive(
		ctx,
		userID,
		isActive,
	)
}
