package models

import (
	"fmt"
	"time"
)

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

type ClientListItem struct {
	ID          int
	Name        string
	Phone       string
	CarsCount   int
	LastVisitAt time.Time
}

func (c Client) FormattedDate() string {
	if c.LastVisitAt.IsZero() {
		return "—"
	}

	return c.LastVisitAt.Format("02.01.2006")
}

func (c ClientListItem) FormattedDate() string {
	if c.LastVisitAt.IsZero() {
		return "—"
	}

	return c.LastVisitAt.Format("02.01.2006")
}

func (c Client) Code() string {
	return fmt.Sprintf("CL-%06d", c.ID)
}
