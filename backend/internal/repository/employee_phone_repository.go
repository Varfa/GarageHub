package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmployeePhoneAlreadyExists = errors.New(
		"employee_phone.already_exists",
	)

	ErrEmployeePhoneNotFound = errors.New(
		"employee_phone.not_found",
	)

	ErrEmployeeLastPhone = errors.New(
		"employee_phone.last_phone",
	)

	ErrEmployeePrimaryPhoneDelete = errors.New(
		"employee_phone.primary_delete",
	)
)

type EmployeePhoneRepository struct {
	db DBTX
}

func NewEmployeePhoneRepository(
	db DBTX,
) *EmployeePhoneRepository {
	return &EmployeePhoneRepository{
		db: db,
	}
}

func (r *EmployeePhoneRepository) Create(
	ctx context.Context,
	phone models.EmployeePhone,
) (*models.EmployeePhone, error) {
	query := `
		INSERT INTO employee_phones (
			employee_id,
			phone,
			label,
			is_primary
		)
		VALUES ($1, $2, $3, FALSE)
		RETURNING
			id,
			employee_id,
			phone,
			label,
			is_primary,
			created_at,
			updated_at
	`

	var createdPhone models.EmployeePhone

	err := r.db.QueryRow(
		ctx,
		query,
		phone.EmployeeID,
		phone.Phone,
		phone.Label,
	).Scan(
		&createdPhone.ID,
		&createdPhone.EmployeeID,
		&createdPhone.Phone,
		&createdPhone.Label,
		&createdPhone.IsPrimary,
		&createdPhone.CreatedAt,
		&createdPhone.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" {
			return nil, ErrEmployeePhoneAlreadyExists
		}

		return nil, fmt.Errorf(
			"создание телефона сотрудника: %w",
			err,
		)
	}

	if phone.IsPrimary {
		if err := r.SetPrimary(
			ctx,
			createdPhone.EmployeeID,
			createdPhone.ID,
		); err != nil {
			return nil, err
		}

		createdPhone.IsPrimary = true
	}

	return &createdPhone, nil
}

func (r *EmployeePhoneRepository) ListByEmployeeID(
	ctx context.Context,
	employeeID int64,
) ([]models.EmployeePhone, error) {
	query := `
		SELECT
			id,
			employee_id,
			phone,
			label,
			is_primary,
			created_at,
			updated_at
		FROM employee_phones
		WHERE employee_id = $1
		ORDER BY
			is_primary DESC,
			id
	`

	rows, err := r.db.Query(
		ctx,
		query,
		employeeID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получение телефонов сотрудника: %w",
			err,
		)
	}

	defer rows.Close()

	var phones []models.EmployeePhone

	for rows.Next() {
		var phone models.EmployeePhone

		err := rows.Scan(
			&phone.ID,
			&phone.EmployeeID,
			&phone.Phone,
			&phone.Label,
			&phone.IsPrimary,
			&phone.CreatedAt,
			&phone.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"сканирование телефона сотрудника: %w",
				err,
			)
		}

		phones = append(
			phones,
			phone,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"итерация телефонов сотрудника: %w",
			err,
		)
	}

	return phones, nil
}

func (r *EmployeePhoneRepository) SetPrimary(
	ctx context.Context,
	employeeID int64,
	phoneID int64,
) error {
	checkQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM employee_phones
			WHERE id = $1
			  AND employee_id = $2
		)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		checkQuery,
		phoneID,
		employeeID,
	).Scan(
		&exists,
	)
	if err != nil {
		return fmt.Errorf(
			"проверка телефона сотрудника: %w",
			err,
		)
	}

	if !exists {
		return ErrEmployeePhoneNotFound
	}

	resetQuery := `
		UPDATE employee_phones
		SET
			is_primary = FALSE,
			updated_at = NOW()
		WHERE employee_id = $1
		  AND is_primary = TRUE
	`

	_, err = r.db.Exec(
		ctx,
		resetQuery,
		employeeID,
	)
	if err != nil {
		return fmt.Errorf(
			"сброс основного телефона: %w",
			err,
		)
	}

	setPrimaryQuery := `
		UPDATE employee_phones
		SET
			is_primary = TRUE,
			updated_at = NOW()
		WHERE id = $1
		  AND employee_id = $2
	`

	commandTag, err := r.db.Exec(
		ctx,
		setPrimaryQuery,
		phoneID,
		employeeID,
	)
	if err != nil {
		return fmt.Errorf(
			"назначение основного телефона: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrEmployeePhoneNotFound
	}

	syncEmployeeQuery := `
		UPDATE employees e
		SET
			phone = ep.phone,
			updated_at = NOW()
		FROM employee_phones ep
		WHERE e.id = $1
		  AND ep.id = $2
		  AND ep.employee_id = e.id
	`

	_, err = r.db.Exec(
		ctx,
		syncEmployeeQuery,
		employeeID,
		phoneID,
	)
	if err != nil {
		return fmt.Errorf(
			"синхронизация основного телефона сотрудника: %w",
			err,
		)
	}

	return nil
}

func (r *EmployeePhoneRepository) Delete(
	ctx context.Context,
	employeeID int64,
	phoneID int64,
) error {
	infoQuery := `
		SELECT
			ep.is_primary,
			(
				SELECT COUNT(*)
				FROM employee_phones
				WHERE employee_id = $1
			)
		FROM employee_phones ep
		WHERE ep.id = $2
		  AND ep.employee_id = $1
	`

	var isPrimary bool
	var phonesCount int

	err := r.db.QueryRow(
		ctx,
		infoQuery,
		employeeID,
		phoneID,
	).Scan(
		&isPrimary,
		&phonesCount,
	)
	if err != nil {
		return ErrEmployeePhoneNotFound
	}

	if phonesCount <= 1 {
		return ErrEmployeeLastPhone
	}

	if isPrimary {
		return ErrEmployeePrimaryPhoneDelete
	}

	deleteQuery := `
		DELETE FROM employee_phones
		WHERE id = $1
		  AND employee_id = $2
		  AND is_primary = FALSE
	`

	commandTag, err := r.db.Exec(
		ctx,
		deleteQuery,
		phoneID,
		employeeID,
	)
	if err != nil {
		return fmt.Errorf(
			"удаление телефона сотрудника: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrEmployeePhoneNotFound
	}

	return nil
}
