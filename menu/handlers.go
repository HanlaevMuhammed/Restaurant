package menu

import (
	db "Day8/database"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var menuMutex sync.RWMutex

func RegisterRoutes(router *gin.Engine, menuCh chan uint) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Главная страница!")
	})

	router.GET("/menu", func(ctx *gin.Context) {
		menuMutex.RLock()
		defer menuMutex.RUnlock()
		var dishes []Dish
		db.DB.Find(&dishes)
		if len(dishes) == 0 {
			ctx.String(200, "Меню пустое")
			return
		}
		ctx.JSON(200, gin.H{"Menu": dishes})
	})

	router.GET("/menu/:name", func(ctx *gin.Context) {
		menuMutex.RLock()
		defer menuMutex.RUnlock()
		var dish Dish
		name := ctx.Param("name")

		if err := db.DB.Where("LOWER(name) = LOWER(?)", name).First(&dish).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
		}
		ctx.JSON(200, gin.H{"dish": dish})
	})

	router.POST("/menu", func(ctx *gin.Context) {
		var newDish Dish
		if err := ctx.BindJSON(&newDish); err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if err := db.DB.Create(&newDish).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Не удалось сохранить меню"})
			return
		}

		menuCh <- newDish.ID

		ctx.JSON(200, gin.H{"Добавлено новое блюдо:": newDish})
	})

	router.DELETE("/menu/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		result := db.DB.Where("LOWER(name) = LOWER(?)", name).Delete(&Dish{})

		if result.RowsAffected == 0 {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
			return
		}

		ctx.JSON(200, gin.H{"status": "удалено"})
	})

	router.PUT("/menu/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		var update UpdateDish
		if err := ctx.BindJSON(&update); err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		var dish Dish

		if err := db.DB.Where("LOWER(name) = LOWER(?)", name).First(&dish).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "Блюдо не найдено"})
			return
		}

		dish.Price = update.Price
		db.DB.Save(&dish)

		ctx.JSON(200, gin.H{"status": "Цена обновлена", "dish": dish})
	})
}

func MonitorNewDish(menuCh chan uint) {
	for id := range menuCh {
		var d Dish
		if err := db.DB.First(&d, id).Error; err != nil {
			fmt.Println("Ошибка загрузки новой позиции:", err)
			continue
		}

		fmt.Printf("Новое блюдо: %s — $%.2f (%.2fг)\n", d.Name, d.Price, d.Weight)
	}
}
