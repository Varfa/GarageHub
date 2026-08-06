package service

import "github.com/Varfa/GarageHub/internal/models"

type CreateOrderInput struct {
	NewClient bool
	NewCar    bool

	ClientID int
	CarID    int

	Client models.Client
	Car    models.Car
	Order  models.Order
}
