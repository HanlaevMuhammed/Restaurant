package storage

import (
	"errors"
	"restaurant_service/internal/domain"

	"gorm.io/gorm"
)

type DishRepository struct {
	db *gorm.DB
}

func NewDishRepository(db *gorm.DB) *DishRepository {
	return &DishRepository{db: db}
}

func (r *DishRepository) GetByID(id uint) (domain.Dish, error) {
	var dish domain.Dish
	err := r.db.First(&dish, id).Error
	return dish, err
}

func (r *DishRepository) GetAll() ([]domain.Dish, error) {
	var dishes []domain.Dish
	err := r.db.Find(&dishes).Error
	return dishes, err
}

func (r *DishRepository) GetByName(name string) (domain.Dish, error) {
	var dish domain.Dish
	err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&dish).Error
	if err != nil {
		return domain.Dish{}, err
	}
	return dish, nil
}

func (r *DishRepository) Create(d domain.Dish) error {
	return r.db.Create(&d).Error
}

func (r *DishRepository) DeleteByName(name string) error {
	result := r.db.Where("LOWER(name) = LOWER(?)", name).Delete(&domain.Dish{})
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return result.Error
}

func (r *DishRepository) UpdatePrice(name string, price float64) (domain.Dish, error) {
	var dish domain.Dish
	if err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&dish).Error; err != nil {
		return domain.Dish{}, err
	}
	dish.Price = price
	if err := r.db.Save(&dish).Error; err != nil {
		return domain.Dish{}, err
	}
	return dish, nil
}
