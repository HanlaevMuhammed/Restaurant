package domain

import "fmt"

type OrderRepository interface {
	GetAll() ([]Order, error)
	Create(o *Order) error
	AttachItems(o *Order, dishes []Dish) error
	GetByID(id uint) (Order, error)
	Save(o *Order) error
}

type OrderService struct {
	Repo OrderRepository
}

func (s *OrderService) MonitorNewOrders(ch chan Order) {
	for incoming := range ch {
		o, err := s.Repo.GetByID(incoming.ID)
		if err != nil {
			fmt.Println("Ошибка загрузки заказа:", err)
			continue
		}

		fmt.Printf("Получен заказ: ID=%d, Сумма=%.2f, блюд=%d\n",
			o.ID, o.Total, len(o.Items))
	}
}
