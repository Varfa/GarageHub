package models

type StockItem struct {
	ID int

	PartName string
	Quantity int
	Unit     string

	Location string

	MinQuantity int

	Note string
}
