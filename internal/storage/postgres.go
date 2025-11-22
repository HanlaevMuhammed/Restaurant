package storage

import (
	"restaurant_service/database"

	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return database.DB
}
