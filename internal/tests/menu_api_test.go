package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restaurant_service/database"
	"restaurant_service/internal/api"
	domain "restaurant_service/internal/domain"
	"restaurant_service/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMenuRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	setupTestDB(t)

	dishRepo := storage.NewDishRepository(database.DB)
	orderRepo := storage.NewOrderRepository(database.DB)

	dishHandler := &api.DishHandler{
		Repo:   dishRepo,
		Notify: make(chan uint, 10),
	}

	orderHandler := &api.OrderHandler{
		OrderRepo: orderRepo,
	}

	return api.InitRouter(dishHandler, orderHandler)
}

func TestGetEmptyMenu(t *testing.T) {
	r := setupMenuRouter(t)
	req, _ := http.NewRequest("GET", "/menu", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Меню пустое")
}

func TestCreateDishAndGetMenu(t *testing.T) {
	r := setupMenuRouter(t)

	newDish := domain.Dish{
		Name:   "Borsh",
		Price:  150.00,
		Weight: 125.00,
	}

	jsonBody, _ := json.Marshal(newDish)
	req, _ := http.NewRequest("POST", "/menu", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Добавлено")

	req2, _ := http.NewRequest("GET", "/menu", nil)
	rr2 := httptest.NewRecorder()

	r.ServeHTTP(rr2, req2)
	assert.Equal(t, 200, rr2.Code)
	assert.Contains(t, rr2.Body.String(), "Borsh")
}
