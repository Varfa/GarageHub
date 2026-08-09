package models

import "time"

type WarehouseItem struct {
	ID                 int64
	Name               string
	SKU                string
	Manufacturer       string
	PurchasePriceCents int64
	SalePriceCents     int64
	Quantity           int
	MinQuantity        int
	Location           string
	Note               string
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
