package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(dish *DishHandler, order *OrderHandler) *gin.Engine {
	r := gin.Default()

	dish.Register(r)
	order.Register(r)

	return r
}
