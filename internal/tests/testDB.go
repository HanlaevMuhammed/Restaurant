package tests

import (
	"restaurant_service/database"
	"restaurant_service/internal/domain"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	err = db.AutoMigrate(&domain.Dish{}, &domain.Order{})
	if err != nil {
		t.Fatalf("Миграция не прошла: %v", err)
	}

	database.DB = db
	return db
}
