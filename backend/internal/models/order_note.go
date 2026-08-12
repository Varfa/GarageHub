package models

import "time"

type OrderNote struct {
	ID         int64
	OrderID    int64
	EmployeeID *int64
	Text       string
	CreatedAt  time.Time
}
