package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
)

type CarRepository struct {
	db *pgxpool.Pool
}

func NewCarRepository(db *pgxpool.Pool) *CarRepository {
	return &CarRepository{
		db: db,
	}
}
func (r *CarRepository) Create(ctx context.Context, car models.Car) error {
	query := `
		INSERT INTO cars (
			client_id,
			brand,
			model,
			year,
			vin,
			plate_number,
			engine,
			power_kw,
			color,
			mileage,
			note,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		car.ClientID,
		car.Brand,
		car.Model,
		car.Year,
		car.VIN,
		car.PlateNumber,
		car.Engine,
		car.PowerKW,
		car.Color,
		car.Mileage,
		car.Note,
		car.CreatedAt,
		car.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("создание автомобиля: %w", err)
	}

	return nil
}

func (r *CarRepository) ListByClientID(
	ctx context.Context,
	clientID int,
) ([]models.Car, error) {
	query := `
		SELECT
			id,
			client_id,
			brand,
			model,
			year,
			vin,
			plate_number,
			engine,
			power_kw,
			color,
			mileage,
			note,
			created_at,
			updated_at
		FROM cars
		WHERE client_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("получение списка автомобилей: %w", err)
	}
	defer rows.Close()

	var cars []models.Car

	for rows.Next() {
		var car models.Car

		err := rows.Scan(
			&car.ID,
			&car.ClientID,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.VIN,
			&car.PlateNumber,
			&car.Engine,
			&car.PowerKW,
			&car.Color,
			&car.Mileage,
			&car.Note,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("сканирование автомобиля: %w", err)
		}

		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("итерация по автомобилям: %w", err)
	}

	return cars, nil

}
func (r *CarRepository) List(
	ctx context.Context,
	search string,
) ([]models.Car, error) {
	query := `
		SELECT
			id,
			client_id,
			brand,
			model,
			year,
			vin,
			plate_number,
			engine,
			power_kw,
			color,
			mileage,
			note,
			created_at,
			updated_at
		FROM cars
		WHERE $1 = ''
			OR brand ILIKE '%' || $1 || '%'
			OR model ILIKE '%' || $1 || '%'
			OR vin ILIKE '%' || $1 || '%'
			OR plate_number ILIKE '%' || $1 || '%'
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("получение списка автомобилей: %w", err)
	}
	defer rows.Close()

	var cars []models.Car

	for rows.Next() {
		var car models.Car

		err := rows.Scan(
			&car.ID,
			&car.ClientID,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.VIN,
			&car.PlateNumber,
			&car.Engine,
			&car.PowerKW,
			&car.Color,
			&car.Mileage,
			&car.Note,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("сканирование автомобиля: %w", err)
		}

		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("итерация по автомобилям: %w", err)
	}

	return cars, nil
}

var ErrCarAlreadyExists = errors.New("автомобиль уже существует")

func (r *CarRepository) ExistsByPlateOrVIN(
	ctx context.Context,
	plateNumber string,
	vin string,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM cars
			WHERE plate_number = $1
			   OR ($2 <> '' AND vin = $2)
		)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		plateNumber,
		vin,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("проверка автомобиля на дубликат: %w", err)
	}

	return exists, nil
}

var ErrCarNotFound = errors.New("автомобиль не найден")

func (r *CarRepository) GetByID(
	ctx context.Context,
	id int,
) (*models.Car, error) {
	query := `
		SELECT
			id,
			client_id,
			brand,
			model,
			year,
			vin,
			plate_number,
			engine,
			power_kw,
			color,
			mileage,
			note,
			created_at,
			updated_at
		FROM cars
		WHERE id = $1
	`

	var car models.Car

	err := r.db.QueryRow(ctx, query, id).Scan(
		&car.ID,
		&car.ClientID,
		&car.Brand,
		&car.Model,
		&car.Year,
		&car.VIN,
		&car.PlateNumber,
		&car.Engine,
		&car.PowerKW,
		&car.Color,
		&car.Mileage,
		&car.Note,
		&car.CreatedAt,
		&car.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarNotFound
		}

		return nil, fmt.Errorf("получение автомобиля по id: %w", err)
	}

	return &car, nil
}
func (r *CarRepository) Update(ctx context.Context, car models.Car) error {
	query := `
		UPDATE cars
		SET
			brand = $1,
			model = $2,
			year = $3,
			vin = $4,
			plate_number = $5,
			engine = $6,
			power_kw = $7,
			color = $8,
			mileage = $9,
			note = $10,
			updated_at = NOW()
		WHERE id = $11
	`

	_, err := r.db.Exec(
		ctx,
		query,
		car.Brand,
		car.Model,
		car.Year,
		car.VIN,
		car.PlateNumber,
		car.Engine,
		car.PowerKW,
		car.Color,
		car.Mileage,
		car.Note,
		car.ID,
	)

	if err != nil {
		return fmt.Errorf("обновление автомобиля: %w", err)
	}

	return nil
}
func (r *CarRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM cars WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление автомобиля: %w", err)
	}

	return nil
}
