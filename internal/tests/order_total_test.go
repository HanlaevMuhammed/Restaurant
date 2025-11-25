package tests

import (
	domain "restaurant_service/internal/domain"
	"testing"
)

func TestOrderTotal(t *testing.T) {
	tests := []struct {
		name  string
		items []domain.Dish
		want  float64
	}{
		{
			name:  "EmptyOrder",
			items: []domain.Dish{},
			want:  0,
		},
		{
			name:  "OneItem",
			items: []domain.Dish{{Price: 150}},
			want:  150,
		},
		{
			name:  "MultipleItems",
			items: []domain.Dish{{Price: 200}, {Price: 300}, {Price: 50.5}},
			want:  550.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := domain.Order{
				Items: tt.items,
			}
			got := order.TotalSum()
			if got != tt.want {
				t.Errorf("Total() = %.2f, want %.2f", got, tt.want)
			}
		})
	}

}
