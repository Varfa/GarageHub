package models

import "time"

type OrderStatusHistory struct {
	ID      int
	OrderID int

	OldStatus string
	NewStatus string
	Comment   string

	CreatedAt time.Time
}
