package models

import "time"

type OrderEmployee struct {
	ID           int64
	OrderID      int64
	EmployeeID   int64
	AssignedAt   time.Time
	UnassignedAt *time.Time
}
