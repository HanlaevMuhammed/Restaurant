package storage

import (
	"restaurant_service/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Create(o *domain.Order) error {
	return r.db.Omit("Items.*").Create(o).Error
}

func (r *OrderRepository) AttachItems(o *domain.Order, dishes []domain.Dish) error {
	return r.db.Model(o).Association("Items").Append(dishes)
}

func (r *OrderRepository) GetByID(id uint) (domain.Order, error) {
	var o domain.Order
	err := r.db.Preload("Items").First(&o, id).Error
	return o, err
}

func (r *OrderRepository) Save(o *domain.Order) error {
	return r.db.Save(o).Error
}
