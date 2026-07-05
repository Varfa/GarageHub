package models

import "time"

type RepairOrder struct {
	ID       int
	ClientID int
	CarID    int

	Number       int
	AcceptedAt   time.Time
	Status       string
	MechanicID   int
	WorkType     string
	Works        string
	Parts        string
	PhotoAfter   string
	PhotosBefore string
	TotalPrice   float64
}
