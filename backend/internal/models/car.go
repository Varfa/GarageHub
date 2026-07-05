package models

import "time"

type Car struct {
	ID       int
	ClientID int

	VIN         string
	PlateNumber string

	Brand string
	Model string

	PowerKW int
	Engine  string

	Mileage   int
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
