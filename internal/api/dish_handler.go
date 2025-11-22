package api

import (
	"sync"

	"restaurant_service/internal/domain"
	"restaurant_service/internal/storage"

	"github.com/gin-gonic/gin"
)

var menuMutex sync.RWMutex

type DishHandler struct {
	Repo   *storage.DishRepository
	Notify chan uint
}

func (h *DishHandler) Register(router *gin.Engine) {

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Главная страница")
	})

	router.GET("/menu", func(ctx *gin.Context) {
		menuMutex.RLock()
		defer menuMutex.RUnlock()

		dishes, _ := h.Repo.GetAll()
		if len(dishes) == 0 {
			ctx.String(200, "Меню пустое")
			return
		}
		ctx.JSON(200, gin.H{"menu": dishes})
	})

	router.GET("/menu/:name", func(ctx *gin.Context) {
		menuMutex.RLock()
		defer menuMutex.RUnlock()

		name := ctx.Param("name")
		dish, err := h.Repo.GetByName(name)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
			return
		}
		ctx.JSON(200, gin.H{"dish": dish})
	})

	router.POST("/menu", func(ctx *gin.Context) {
		var newDish domain.Dish

		if err := ctx.BindJSON(&newDish); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.Repo.Create(newDish); err != nil {
			ctx.JSON(500, gin.H{"error": "Ошибка сохранения"})
			return
		}

		h.Notify <- newDish.ID

		ctx.JSON(200, gin.H{"status": "Добавлено", "dish": newDish})
	})

	router.DELETE("/menu/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		if err := h.Repo.DeleteByName(name); err != nil {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
			return
		}
		ctx.JSON(200, gin.H{"status": "Удалено"})
	})

	router.PUT("/menu/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		var update domain.UpdateDish
		if err := ctx.BindJSON(&update); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		dish, err := h.Repo.UpdatePrice(name, update.Price)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
			return
		}

		ctx.JSON(200, gin.H{"status": "Цена обновлена", "dish": dish})
	})
}
