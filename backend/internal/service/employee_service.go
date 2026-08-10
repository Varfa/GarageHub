package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var (
	ErrEmployeeInvalidID = errors.New(
		"employee.invalid_id",
	)

	ErrEmployeePhoneInvalidID = errors.New(
		"employee.phone.invalid_id",
	)

	ErrEmployeeFirstNameRequired = errors.New(
		"employee.validation.first_name_required",
	)

	ErrEmployeeLastNameRequired = errors.New(
		"employee.validation.last_name_required",
	)

	ErrEmployeePhoneRequired = errors.New(
		"employee.validation.phone_required",
	)

	ErrEmployeePositionRequired = errors.New(
		"employee.validation.position_required",
	)

	ErrEmployeeAlreadyArchived = errors.New(
		"employee.already_archived",
	)

	ErrEmployeeAlreadyActive = errors.New(
		"employee.already_active",
	)

	ErrEmployeeArchived = errors.New(
		"employee.archived",
	)
	ErrEmployeePhoneAlreadyExists = repository.ErrEmployeePhoneAlreadyExists

	ErrEmployeePhoneNotFound = repository.ErrEmployeePhoneNotFound

	ErrEmployeeLastPhone = repository.ErrEmployeeLastPhone

	ErrEmployeePrimaryPhoneDelete = repository.ErrEmployeePrimaryPhoneDelete
)

type EmployeeService struct {
	employeeRepository         *repository.EmployeeRepository
	employeePositionRepository *repository.EmployeePositionRepository
	employeePhoneRepository    *repository.EmployeePhoneRepository
}

func NewEmployeeService(
	employeeRepository *repository.EmployeeRepository,
	employeePositionRepository *repository.EmployeePositionRepository,
	employeePhoneRepository *repository.EmployeePhoneRepository,
) *EmployeeService {
	return &EmployeeService{
		employeeRepository:         employeeRepository,
		employeePositionRepository: employeePositionRepository,
		employeePhoneRepository:    employeePhoneRepository,
	}
}

func (s *EmployeeService) ListPositions(
	ctx context.Context,
) ([]models.EmployeePosition, error) {
	return s.employeePositionRepository.ListActive(ctx)
}

func (s *EmployeeService) Create(
	ctx context.Context,
	employee models.Employee,
) error {
	employee = normalizeEmployee(employee)

	if err := validateEmployee(employee); err != nil {
		return err
	}

	return s.employeeRepository.Create(
		ctx,
		employee,
	)
}

func (s *EmployeeService) ListActive(
	ctx context.Context,
	search string,
) ([]models.EmployeeListItem, error) {
	search = strings.TrimSpace(search)

	return s.employeeRepository.ListActive(
		ctx,
		search,
	)
}

func (s *EmployeeService) ListArchived(
	ctx context.Context,
	search string,
) ([]models.EmployeeListItem, error) {
	search = strings.TrimSpace(search)

	return s.employeeRepository.ListArchived(
		ctx,
		search,
	)
}

func (s *EmployeeService) GetByID(
	ctx context.Context,
	id int64,
) (*models.Employee, error) {
	if id <= 0 {
		return nil, ErrEmployeeInvalidID
	}

	return s.employeeRepository.GetByID(
		ctx,
		id,
	)
}

func (s *EmployeeService) ListPhones(
	ctx context.Context,
	employeeID int64,
) ([]models.EmployeePhone, error) {
	if employeeID <= 0 {
		return nil, ErrEmployeeInvalidID
	}

	return s.employeePhoneRepository.ListByEmployeeID(
		ctx,
		employeeID,
	)
}

func (s *EmployeeService) AddPhone(
	ctx context.Context,
	phone models.EmployeePhone,
) error {
	phone.Phone = strings.TrimSpace(phone.Phone)
	phone.Label = strings.TrimSpace(phone.Label)

	if phone.EmployeeID <= 0 {
		return ErrEmployeeInvalidID
	}

	if err := s.ensureEmployeeActive(
		ctx,
		phone.EmployeeID,
	); err != nil {
		return err
	}

	if phone.Phone == "" {
		return ErrEmployeePhoneRequired
	}

	if phone.Label == "" {
		phone.Label = "additional"
	}

	phones, err := s.employeePhoneRepository.ListByEmployeeID(
		ctx,
		phone.EmployeeID,
	)
	if err != nil {
		return err
	}

	if len(phones) == 0 {
		phone.IsPrimary = true
	}

	_, err = s.employeePhoneRepository.Create(
		ctx,
		phone,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeService) SetPrimaryPhone(
	ctx context.Context,
	employeeID int64,
	phoneID int64,
) error {
	if employeeID <= 0 {
		return ErrEmployeeInvalidID
	}

	if phoneID <= 0 {
		return ErrEmployeePhoneInvalidID
	}

	if err := s.ensureEmployeeActive(
		ctx,
		employeeID,
	); err != nil {
		return err
	}

	return s.employeePhoneRepository.SetPrimary(
		ctx,
		employeeID,
		phoneID,
	)
}

func (s *EmployeeService) DeletePhone(
	ctx context.Context,
	employeeID int64,
	phoneID int64,
) error {
	if employeeID <= 0 {
		return ErrEmployeeInvalidID
	}

	if phoneID <= 0 {
		return ErrEmployeePhoneInvalidID
	}

	if err := s.ensureEmployeeActive(
		ctx,
		employeeID,
	); err != nil {
		return err
	}

	return s.employeePhoneRepository.Delete(
		ctx,
		employeeID,
		phoneID,
	)
}

func (s *EmployeeService) Update(
	ctx context.Context,
	employee models.Employee,
) error {
	employee = normalizeEmployee(employee)

	if employee.ID <= 0 {
		return ErrEmployeeInvalidID
	}

	if err := s.ensureEmployeeActive(
		ctx,
		employee.ID,
	); err != nil {
		return err
	}

	if err := validateEmployee(employee); err != nil {
		return err
	}

	return s.employeeRepository.Update(
		ctx,
		employee,
	)
}

func (s *EmployeeService) Archive(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrEmployeeInvalidID
	}

	employee, err := s.employeeRepository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if !employee.IsActive {
		return ErrEmployeeAlreadyArchived
	}

	return s.employeeRepository.SetActive(
		ctx,
		id,
		false,
	)
}

func (s *EmployeeService) Restore(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrEmployeeInvalidID
	}

	employee, err := s.employeeRepository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if employee.IsActive {
		return ErrEmployeeAlreadyActive
	}

	return s.employeeRepository.SetActive(
		ctx,
		id,
		true,
	)
}

func (s *EmployeeService) ensureEmployeeActive(
	ctx context.Context,
	id int64,
) error {
	employee, err := s.employeeRepository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if !employee.IsActive {
		return ErrEmployeeArchived
	}

	return nil
}

func normalizeEmployee(
	employee models.Employee,
) models.Employee {
	employee.FirstName = strings.TrimSpace(
		employee.FirstName,
	)

	employee.LastName = strings.TrimSpace(
		employee.LastName,
	)

	employee.Phone = strings.TrimSpace(
		employee.Phone,
	)

	if employee.Email != nil {
		email := strings.TrimSpace(
			*employee.Email,
		)

		if email == "" {
			employee.Email = nil
		} else {
			employee.Email = &email
		}
	}

	return employee
}

func validateEmployee(
	employee models.Employee,
) error {
	if employee.FirstName == "" {
		return ErrEmployeeFirstNameRequired
	}

	if employee.LastName == "" {
		return ErrEmployeeLastNameRequired
	}

	if employee.Phone == "" {
		return ErrEmployeePhoneRequired
	}

	if employee.PositionID <= 0 {
		return ErrEmployeePositionRequired
	}

	return nil
}
