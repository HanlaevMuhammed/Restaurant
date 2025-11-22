package menu

import "gorm.io/gorm"

type Dish struct {
	gorm.Model
	Name   string  `gorm:"uniqueIndex:idx_name_deleted_at,where:deleted_at IS NULL" json:"name"`
	Price  float64 `gorm:"not null" json:"price"`
	Weight float64 `json:"weight"`
}

type UpdateDish struct {
	Price float64 `json:"price"`
}
