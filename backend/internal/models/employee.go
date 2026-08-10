package models

import "time"

type EmployeePosition struct {
	ID          int64
	Code        string
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EmployeeListItem struct {
	ID           int64
	FirstName    string
	LastName     string
	Phone        string
	Email        *string
	PositionName string
	PositionCode string
	IsActive     bool
}

type Employee struct {
	ID         int64
	FirstName  string
	LastName   string
	Phone      string
	Email      *string
	PositionID int64
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type EmployeePhone struct {
	ID         int64
	EmployeeID int64
	Phone      string
	Label      string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
