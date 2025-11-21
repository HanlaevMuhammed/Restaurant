package menu

import (
	db "Day8/database"
	"sync"

	"github.com/gin-gonic/gin"
)

var menuMutex sync.RWMutex

func RegisterRoutes(router *gin.Engine, menuCh chan Dish) {
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

		menuCh <- newDish

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

// func AddRandomDishes(menuCh chan Dish) {
// 	randomDishes := []Dish{
// 		{Name: "Случайный суп", Price: 150.0},
// 		{Name: "Случайный салат", Price: 120.0},
// 		{Name: "Случайное основное блюдо", Price: 250.0},
// 		{Name: "Случайный десерт", Price: 100.0},
// 		{Name: "Случайный напиток", Price: 80.0},
// 	}

// 	ticker := time.NewTicker(3 * time.Second)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		randomDish := randomDishes[rand.IntN(len(randomDishes))]
// 		randomDish.Price = 50 + float64(rand.IntN(300))

// 		menuMutex.Lock()
// 		*menu = append(*menu, randomDish)
// 		currentMenu := *menu
// 		menuMutex.Unlock()

// 		if err := storage.Save(filename, currentMenu); err != nil {
// 			fmt.Printf("Ошибка сохранения меню: %v\n", err)
// 		}

// 		menuCh <- randomDish
// 		fmt.Printf("Добавлено случайное блюдо: %s (%.2f руб)\n", randomDish.Name, randomDish.Price)
// 	}
// }

// func PrintCurrentMenu(menu *[]Dish) {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		menuMutex.RLock()
// 		currentMenu := *menu
// 		menuMutex.RUnlock()

// 		fmt.Println("\n====== ТЕКУЩЕЕ МЕНЮ ======")
// 		if len(currentMenu) == 0 {
// 			fmt.Println("Меню пустое")
// 		} else {
// 			for i, dish := range currentMenu {
// 				fmt.Printf("%d. %s - %.2f руб\n", i+1, dish.Name, dish.Price)
// 			}
// 		}
// 		fmt.Printf("==========================\nВсего блюд: %d\n", len(currentMenu))
// 	}
// }

// func MonitoringNewDish(menuCh chan Dish) {
// 	for newDish := range menuCh {
// 		fmt.Printf("Новое блюдо в меню: %s - %.2f руб\n", newDish.Name, newDish.Price)
// 	}
// }
