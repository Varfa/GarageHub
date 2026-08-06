package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
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
		return nil, errors.New(
			"некорректный id сотрудника",
		)
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
		return nil, errors.New(
			"некорректный id сотрудника",
		)
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
		return errors.New(
			"некорректный id сотрудника",
		)
	}

	if err := s.ensureEmployeeActive(
		ctx,
		phone.EmployeeID,
	); err != nil {
		return err
	}

	if phone.Phone == "" {
		return errors.New(
			"необходимо указать номер телефона",
		)
	}

	if phone.Label == "" {
		phone.Label = "Дополнительный"
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
		return errors.New(
			"некорректный id сотрудника",
		)
	}

	if phoneID <= 0 {
		return errors.New(
			"некорректный id телефона",
		)
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
		return errors.New(
			"некорректный id сотрудника",
		)
	}

	if phoneID <= 0 {
		return errors.New(
			"некорректный id телефона",
		)
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
		return errors.New(
			"некорректный id сотрудника",
		)
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
		return errors.New(
			"некорректный id сотрудника",
		)
	}

	employee, err := s.employeeRepository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if !employee.IsActive {
		return errors.New(
			"сотрудник уже находится в архиве",
		)
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
		return errors.New(
			"некорректный id сотрудника",
		)
	}

	employee, err := s.employeeRepository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if employee.IsActive {
		return errors.New(
			"сотрудник уже активен",
		)
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
		return errors.New(
			"сотрудник находится в архиве: сначала восстановите его",
		)
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
		return errors.New(
			"необходимо указать имя",
		)
	}

	if employee.LastName == "" {
		return errors.New(
			"необходимо указать фамилию",
		)
	}

	if employee.Phone == "" {
		return errors.New(
			"необходимо указать основной номер телефона",
		)
	}

	if employee.PositionID <= 0 {
		return errors.New(
			"необходимо выбрать должность",
		)
	}

	return nil
}
