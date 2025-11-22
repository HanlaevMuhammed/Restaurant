package database

import (
	"fmt"
	"log"
	"restaurant_service/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	config.InitConfig()
	dbConf := config.Config.Database

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConf.Host, dbConf.User, dbConf.Password, dbConf.Name, dbConf.Port, dbConf.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	DB = db
	log.Println("База данных успешно подключена!")
}
