package domain

import "fmt"

type DishRepository interface {
	GetByID(id uint) (Dish, error)
}

type MenuService struct {
	Repo DishRepository
}

func (s *MenuService) MonitorNewDish(ch chan uint) {
	for id := range ch {
		d, err := s.Repo.GetByID(id)
		if err != nil {
			fmt.Println("Ошибка загрузки блюда:", err)
			continue
		}
		fmt.Printf("Новое блюдо: %s — $%.2f (%.2fг)\n", d.Name, d.Price, d.Weight)
	}
}
