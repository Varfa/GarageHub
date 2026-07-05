package models

import "time"

type Part struct {
	ID int

	Name         string
	Manufacturer string
	PartNumber   string

	Note string

	CreatedAt time.Time
	UpdatedAt time.Time
}
