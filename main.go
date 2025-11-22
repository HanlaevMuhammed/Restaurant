package main

import (
	"Day8/database"
	menu "Day8/menu"
	order "Day8/order"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	database.DB.AutoMigrate(&menu.Dish{}, &order.Order{})

	orderCh := make(chan order.Order, 10)
	menuCh := make(chan uint, 10)
	go order.MonitorNewOrders(orderCh)

	go menu.MonitorNewDish(menuCh)

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
