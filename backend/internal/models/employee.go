package models

import "time"

type Employee struct {
	ID     int
	Number int

	Name    string
	Phone   string
	Email   string
	Address string
	Note    string

	CreatedAt time.Time
	UpdatedAt time.Time
	Percent   int
	Position  string
}
