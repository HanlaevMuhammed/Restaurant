package main

import (
	"Day8/database"
	menu "Day8/menu"
	order "Day8/order"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func monitorOrders(orders []order.Order) {
	for {
		fmt.Println("Количество заказов:", len(orders))
		time.Sleep(10 * time.Second)
	}
}

func monitorMenu(menu []menu.Dish) {
	for {
		fmt.Println("Список меню:", menu)
		time.Sleep(5 * time.Second)
	}
}

func main() {

	database.Connect()
	database.DB.AutoMigrate(&menu.Dish{}, &order.Order{})

	orderCh := make(chan order.Order, 10)
	menuCh := make(chan menu.Dish, 3)
	go order.MonitorNewOrders(orderCh)

	// go monitorOrders(orderData)
	// go monitorMenu(menuData)
	// go menu.MonitoringNewDish(menuCh)

	// go menu.AddRandomDishes(menuCh)
	// go menu.PrintCurrentMenu()

	router := gin.Default()

	menu.RegisterRoutes(router, menuCh)
	order.RegisterRoutes(router, orderCh)

	router.GET("/debug/routes", func(ctx *gin.Context) {
		routes := router.Routes()
		list := []gin.H{}
		for _, r := range routes {
			list = append(list, gin.H{
				"method": r.Method,
				"path":   r.Path,
			})
		}
		ctx.JSON(200, list)
	})
	router.Run()

}
