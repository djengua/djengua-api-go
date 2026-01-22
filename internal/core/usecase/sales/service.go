// internal/core/usecase/sales/service.go
package sales

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	sales  ports.SaleRepository
	orders ports.OrderRepository
}

func NewService(sales ports.SaleRepository, orders ports.OrderRepository) *Service {
	return &Service{sales: sales, orders: orders}
}

func (s *Service) RegisterSale(ctx context.Context, in ports.SaleCreateInput) (domain.Sale, domain.Order, error) {
	if strings.TrimSpace(in.UserID) == "" || strings.TrimSpace(in.OrderID) == "" {
		return domain.Sale{}, domain.Order{}, errors.New("user_id and order_id required")
	}
	if in.Amount <= 0 {
		return domain.Sale{}, domain.Order{}, errors.New("amount must be > 0")
	}
	if strings.TrimSpace(in.Currency) == "" {
		in.Currency = "MXN"
	}
	if in.Method == "" {
		in.Method = domain.PayOther
	}

	ord, err := s.orders.GetOrder(ctx, in.OrderID)
	if err != nil {
		return domain.Sale{}, domain.Order{}, err
	}
	if ord.UserID != in.UserID {
		return domain.Sale{}, domain.Order{}, errors.New("order does not belong to user")
	}
	if ord.Status == domain.OrderStatusPaid {
		return domain.Sale{}, domain.Order{}, errors.New("order already paid")
	}

	soldAt := time.Now().UTC()
	if in.SoldAt != nil && !in.SoldAt.IsZero() {
		soldAt = in.SoldAt.UTC()
	}

	sale, err := s.sales.CreateSale(ctx, domain.Sale{
		OrderID:   in.OrderID,
		UserID:    in.UserID,
		Method:    in.Method,
		Amount:    in.Amount,
		Currency:  in.Currency,
		Reference: in.Reference,
		SoldAt:    soldAt,
	})
	if err != nil {
		return domain.Sale{}, domain.Order{}, err
	}

	ord, err = s.orders.UpdateOrderStatus(ctx, ord.ID, domain.OrderStatusPaid)
	if err != nil {
		return domain.Sale{}, domain.Order{}, err
	}

	return sale, ord, nil
}
