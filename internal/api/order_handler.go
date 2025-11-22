package api

import (
	"fmt"
	"strconv"
	"strings"

	"restaurant_service/internal/domain"
	"restaurant_service/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type OrderHandler struct {
	OrderRepo *storage.OrderRepository
	DishRepo  *storage.DishRepository
	Notify    chan domain.Order
}

func (h *OrderHandler) Register(router *gin.Engine) {

	router.GET("/orders", func(ctx *gin.Context) {
		orders, _ := h.OrderRepo.GetAll()
		if len(orders) == 0 {
			ctx.String(200, "Заказов нет")
			return
		}
		ctx.JSON(200, gin.H{"orders": orders})
	})

	router.POST("/orders", func(ctx *gin.Context) {
		body, _ := ctx.GetRawData()
		names := strings.Split(strings.TrimSpace(string(body)), ",")

		var items []domain.Dish
		var notFound []string
		total := 0.0

		for _, name := range names {
			name = strings.TrimSpace(name)
			d, err := h.DishRepo.GetByName(name)
			if err != nil {
				notFound = append(notFound, name)
				continue
			}

			items = append(items, d)
			total += d.Price
		}

		newOrder := domain.Order{Items: items, Total: total}
		if err := h.OrderRepo.Create(&newOrder); err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка создания заказа"})
			return
		}

		if err := h.OrderRepo.AttachItems(&newOrder, items); err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка привязки блюд к заказу"})
			return
		}

		h.Notify <- newOrder

		ctx.JSON(200, gin.H{
			"status":    "Заказ создан",
			"order":     newOrder,
			"not_found": notFound,
		})
	})

	router.POST("/orders/pay", func(ctx *gin.Context) {
		var req struct{ ID uint }

		if ctx.BindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Ошибка JSON"})
			return
		}

		o, err := h.OrderRepo.GetByID(req.ID)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Заказ не найден"})
			return
		}

		if o.Paid {
			ctx.JSON(400, gin.H{"error": "Уже оплачен"})
			return
		}

		o.Paid = true
		h.OrderRepo.Save(&o)

		ctx.JSON(200, gin.H{"status": "Оплачено", "order": o})
	})

	router.POST("/orders/tip", func(ctx *gin.Context) {
		var req struct {
			ID     uint
			Amount float64
		}

		ctx.BindJSON(&req)

		o, err := h.OrderRepo.GetByID(req.ID)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Заказ не найден"})
			return
		}

		o.Tip += req.Amount
		h.OrderRepo.Save(&o)
		ctx.JSON(200, gin.H{"status": "Чаевые добавлены", "order": o})
	})

	router.GET("/orders/qrcode/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))

		o, err := h.OrderRepo.GetByID(uint(id))
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Не найден"})
			return
		}

		url := fmt.Sprintf("http://localhost:8080/orders/%d", o.ID)

		png, _ := qrcode.Encode(url, qrcode.Medium, 256)
		ctx.Data(200, "image/png", png)
	})
}
