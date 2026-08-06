package models

import (
	"fmt"
	"time"
)

type Car struct {
	ID       int
	ClientID int

	Brand string
	Model string
	Year  int

	VIN         string
	PlateNumber string

	Engine  string
	PowerKW int

	Color   string
	Mileage int
	Note    string

	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerName string
}

func (c Car) Code() string {
	return fmt.Sprintf("CAR-%06d", c.ID)
}
