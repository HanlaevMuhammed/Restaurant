package order

import (
	db "Day8/database"
	"Day8/menu"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, orderCh chan Order) {

	router.GET("/orders", func(ctx *gin.Context) {
		var orders []Order
		db.DB.Preload("Items").Find(&orders)
		if len(orders) == 0 {
			ctx.String(200, "Заказов пока нет")
			return
		}
		ctx.JSON(200, gin.H{"Orders": orders})
	})

	router.POST("/orders", func(ctx *gin.Context) {
		body, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Не удалось прочиать файл"})
			return
		}

		input := strings.TrimSpace(string(body))
		if input == "" {
			ctx.JSON(400, gin.H{"error": "Пустой заказ"})
			return
		}

		names := strings.Split(input, ",")
		var items []menu.Dish
		total := 0.0
		var notFound []string

		for _, name := range names {
			name = strings.TrimSpace(name)
			var dish menu.Dish

			if err := db.DB.Where("LOWER(name) = LOWER(?)", name).First(&dish).Error; err != nil {
				notFound = append(notFound, name)
				continue
			}

			items = append(items, dish)
			total += dish.Price
		}

		newOrder := Order{
			ID:    len(items) + 1,
			Items: items,
			Total: int(total),
		}

		if err = db.DB.Omit("Items.*").Create(&newOrder).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка создания заказа"})
			return
		}

		db.DB.Model(&newOrder).Association("Items").Append(items)

		orderCh <- newOrder

		ctx.JSON(200, gin.H{"stauts": "Заказ создан", "order": newOrder, "Блюда отсутствующие в меню": notFound})
	})

}

func MonitorNewOrders(orderCh chan Order) {
	for newOrder := range orderCh {
		fmt.Printf("Получен новый заказ: ID=%d, Сумма=%d, Количество блюд=%d\n",
			newOrder.ID, newOrder.Total, len(newOrder.Items))
	}
}
