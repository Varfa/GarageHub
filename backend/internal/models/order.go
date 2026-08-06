package models

import "time"

type Order struct {
	ID       int
	ClientID int
	CarID    int

	Complaint string
	Diagnosis string
	Note      string

	Status string

	EstimatedCostCents int64
	FinalCostCents     int64

	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	UpdatedAt   time.Time
}

type OrderWork struct {
	ID      int
	OrderID int

	Name        string
	Description string
	Quantity    float64
	PriceCents  int64

	IsAdditional bool
	IsApproved   bool
	ApprovedAt   *time.Time
}

type OrderPart struct {
	ID      int
	OrderID int

	Name       string
	PartNumber string
	Quantity   float64
	PriceCents int64

	IsAdditional bool
	IsApproved   bool
	ApprovedAt   *time.Time
}

type Notification struct {
	ID       int
	OrderID  int
	ClientID int

	Channel string
	Message string
	Status  string

	SentAt    *time.Time
	CreatedAt time.Time
}
