package main

import (
	"restaurant_service/database"
	"restaurant_service/internal/api"
	"restaurant_service/internal/domain"
	"restaurant_service/internal/storage"
)

func main() {
	database.Connect()

	orderCh := make(chan domain.Order, 10)
	menuCh := make(chan uint, 10)

	dishRepo := storage.NewDishRepository(database.DB)
	orderRepo := storage.NewOrderRepository(database.DB) 

	menuService := &domain.MenuService{Repo: dishRepo}
	orderService := &domain.OrderService{Repo: orderRepo}

	go menuService.MonitorNewDish(menuCh)
	go orderService.MonitorNewOrders(orderCh)

	dishHandler := &api.DishHandler{
		Repo:   dishRepo,
		Notify: menuCh,
	}

	orderHandler := &api.OrderHandler{
		OrderRepo: orderRepo,
		DishRepo:  dishRepo,
		Notify:    orderCh,
	}

	router := api.InitRouter(dishHandler, orderHandler)

	router.Run(":8080")
}
