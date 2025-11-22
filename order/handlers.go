package order

import (
	db "Day8/database"
	"Day8/menu"
	"fmt"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"

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
			Items: items,
			Total: total,
		}

		if err = db.DB.Omit("Items.*").Create(&newOrder).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка создания заказа"})
			return
		}

		db.DB.Model(&newOrder).Association("Items").Append(items)

		orderCh <- newOrder

		ctx.JSON(200, gin.H{"stauts": "Заказ создан", "order": newOrder, "Блюда отсутствующие в меню": notFound})
	})

	router.POST("/orders/pay", func(ctx *gin.Context) {
		var req struct {
			ID uint `json:"id"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"error": "Ошибка чтения JSON"})
			return
		}

		var o Order

		if err := db.DB.Preload("Items").First(&o, req.ID).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "Заказ не найден"})
			return
		}

		if o.Paid {
			ctx.JSON(400, gin.H{"error": "Заказ уже оплачен"})
			return
		}

		o.Paid = true

		if err := db.DB.Save(&o).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка сохранения заказа"})
			return
		}

		ctx.JSON(200, gin.H{"status": "Оплата прошла успешно", "order": o})

	})

	router.POST("/orders/tip", func(ctx *gin.Context) {
		var req struct {
			ID     uint    `json:"id"`
			Amount float64 `json:"amount"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"error": "Ошибка чтения JSON"})
			return
		}

		var o Order
		if err := db.DB.First(&o, req.ID).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "Заказ не найден"})
			return
		}

		o.Tip += req.Amount
		if err := db.DB.Save(&o).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка сохранения заказа"})
			return
		}

		ctx.JSON(200, gin.H{"status": "Чаевые добавлены", "order": o})

	})

	router.GET("/orders/qrcode/:id", func(ctx *gin.Context) {

		idstr := ctx.Param("id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Не удалось преобразовать в str"})
			return
		}

		var o Order
		if err := db.DB.First(&o, id).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "Заказ не найден"})
			return
		}

		url := fmt.Sprintf("http://localhost:8080/orders/qrcode/%d", id)

		png, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка генерации QRcode"})
			return
		}

		ctx.Data(200, "image/png", png)

	})

}

func MonitorNewOrders(orderCh chan Order) {
	for newOrder := range orderCh {
		var o Order
		if err := db.DB.Preload("Items").First(&o, newOrder.ID).Error; err != nil {
			fmt.Println("Ошибка загрузки заказа:", err)
			return
		}
		fmt.Printf("Получен новый заказ: ID=%d, Сумма=%.2f, Количество блюд=%d\n",
			o.ID, o.Total, len(o.Items))
	}
}
