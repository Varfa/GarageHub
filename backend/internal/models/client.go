package models

import "time"

type Client struct {
	ID     int
	Number int

	Name    string
	Phone   string
	Email   string
	Address string
	Note    string

	LastVisitAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
