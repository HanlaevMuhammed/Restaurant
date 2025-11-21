package order

import (
	menu "Day8/menu"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID    int         `gorm:"primaryKey" json:"id"`
	Items []menu.Dish `gorm:"many2many:order_dishes" json:"items"`
	Total int         `json:"total"`
}
