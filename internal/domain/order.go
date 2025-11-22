package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Items []Dish  `gorm:"many2many:order_dishes" json:"items"`
	Total float64 `json:"total"`
	Tip   float64 `json:"tip"`
	Paid  bool    `json:"paid"`
}
