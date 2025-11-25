package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restaurant_service/database"
	"restaurant_service/internal/api"
	"restaurant_service/internal/domain"
	"restaurant_service/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupOrderRouter(t *testing.T) *gin.Engine {
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
		DishRepo:  dishRepo,
		Notify:    make(chan domain.Order, 10),
	}

	return api.InitRouter(dishHandler, orderHandler)
}

func TestCreateOrder(t *testing.T) {
	r := setupOrderRouter(t)

	newDish := domain.Dish{
		Name:   "Borsh",
		Price:  150.00,
		Weight: 125.00,
	}

	jsonDish, _ := json.Marshal(newDish)
	reqDish, _ := http.NewRequest("POST", "/menu", bytes.NewBuffer(jsonDish))
	reqDish.Header.Set("Content-Type", "application/json")
	rrDish := httptest.NewRecorder()
	r.ServeHTTP(rrDish, reqDish)

	orderBody := "Borsh"

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer([]byte(orderBody)))
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	t.Logf("Response: %d - %s", rr.Code, rr.Body.String())
	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Заказ создан")
}

func TestCreateOrderWithNonExistentDish(t *testing.T) {
	r := setupOrderRouter(t)

	orderBody := "NonExistentDish"

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer([]byte(orderBody)))
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	t.Logf("Response: %d - %s", rr.Code, rr.Body.String())
	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "not_found")
}

func TestGetAllOrders(t *testing.T) {
	r := setupOrderRouter(t)

	newDish := domain.Dish{
		Name:   "Salad",
		Price:  100.00,
		Weight: 200.00,
	}

	jsonDish, _ := json.Marshal(newDish)
	reqDish, _ := http.NewRequest("POST", "/menu", bytes.NewBuffer(jsonDish))
	reqDish.Header.Set("Content-Type", "application/json")
	rrDish := httptest.NewRecorder()
	r.ServeHTTP(rrDish, reqDish)

	orderBody := "Salad"

	reqCreate, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer([]byte(orderBody)))
	reqCreate.Header.Set("Content-Type", "text/plain")
	rrCreate := httptest.NewRecorder()
	r.ServeHTTP(rrCreate, reqCreate)

	t.Logf("Create order response: %d - %s", rrCreate.Code, rrCreate.Body.String())

	req, _ := http.NewRequest("GET", "/orders", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	t.Logf("Get orders response: %d - %s", rr.Code, rr.Body.String())
	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "orders")
}

func TestCreateOrderWithMultipleDishes(t *testing.T) {
	r := setupOrderRouter(t)
	dishes := []domain.Dish{
		{Name: "Borsh", Price: 150.00, Weight: 125.00},
		{Name: "Salad", Price: 100.00, Weight: 200.00},
	}

	for _, dish := range dishes {
		jsonDish, _ := json.Marshal(dish)
		reqDish, _ := http.NewRequest("POST", "/menu", bytes.NewBuffer(jsonDish))
		reqDish.Header.Set("Content-Type", "application/json")
		rrDish := httptest.NewRecorder()
		r.ServeHTTP(rrDish, reqDish)
	}

	orderBody := "Borsh, Salad"

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer([]byte(orderBody)))
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	t.Logf("Response: %d - %s", rr.Code, rr.Body.String())
	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Заказ создан")
}
