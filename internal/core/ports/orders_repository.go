package ports

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/core/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o domain.Order) (domain.Order, error)
	GetOrder(ctx context.Context, id string) (domain.Order, error)
	ListOrdersByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, st domain.OrderStatus) (domain.Order, error)
}

type SaleRepository interface {
	CreateSale(ctx context.Context, s domain.Sale) (domain.Sale, error)
	GetSale(ctx context.Context, id string) (domain.Sale, error)
	ListSalesByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Sale, error)
}
