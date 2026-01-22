// internal/core/usecase/orders/service.go
package orders

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	orders ports.OrderRepository
	// opcional: products ports.ProductRepository si quieres validar precios contra catálogo
}

func NewService(orders ports.OrderRepository) *Service {
	return &Service{orders: orders}
}

func (s *Service) Create(ctx context.Context, in ports.OrderCreateInput) (domain.Order, error) {
	if strings.TrimSpace(in.UserID) == "" {
		return domain.Order{}, errors.New("user_id required")
	}
	if strings.TrimSpace(in.Currency) == "" {
		in.Currency = "MXN"
	}
	if len(in.Items) == 0 {
		return domain.Order{}, errors.New("items required")
	}

	// calc totals
	var subtotal float64
	for i := range in.Items {
		if in.Items[i].Qty <= 0 {
			return domain.Order{}, errors.New("qty must be > 0")
		}
		in.Items[i].LineTotal = round2(float64(in.Items[i].Qty) * in.Items[i].UnitPrice)
		subtotal += in.Items[i].LineTotal
	}
	subtotal = round2(subtotal)

	// tax: si quieres, configurable. Aquí 0 por defecto:
	tax := 0.0
	total := round2(subtotal + tax)

	// folio simple (puedes mejorarlo a secuencia)
	number := fmt.Sprintf("ORD-%d", time.Now().UTC().Unix())

	o := domain.Order{
		Number:   number,
		UserID:   in.UserID,
		Status:   domain.OrderStatusPending,
		Currency: in.Currency,
		Items:    in.Items,
		Subtotal: subtotal,
		Tax:      tax,
		Total:    total,
		Notes:    in.Notes,
	}
	return s.orders.CreateOrder(ctx, o)
}

func (s *Service) ListMine(ctx context.Context, userID string, page, pageSize int) ([]domain.Order, error) {
	return s.orders.ListOrdersByUser(ctx, userID, page, pageSize)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
