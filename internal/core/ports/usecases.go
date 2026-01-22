package ports

import (
	"context"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
)

type AuthClaims struct {
	UserID string
	Role   domain.UserRole
}

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (domain.User, string, error)
	Login(ctx context.Context, email, password string) (domain.User, string, error)
	ParseToken(tokenString string) (AuthClaims, error)
	GetUser(ctx context.Context, id string) (domain.User, error)
}

type ProductService interface {
	List(ctx context.Context, filters ProductFilters) ([]domain.Product, error)
	Get(ctx context.Context, id string) (domain.Product, error)
	Create(ctx context.Context, in domain.Product) (domain.Product, error)
	Update(ctx context.Context, id string, in domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id string) error
}

type CategoryService interface {
	List(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	Get(ctx context.Context, id string) (domain.Category, error)
	Create(ctx context.Context, in domain.Category) (domain.Category, error)
	Update(ctx context.Context, id string, in domain.Category) (domain.Category, error)
	Delete(ctx context.Context, id string) error
}

type CollectionService interface {
	List(ctx context.Context, page, pageSize int) ([]domain.Collection, error)
	Get(ctx context.Context, id string) (domain.Collection, error)
	Create(ctx context.Context, in domain.Collection) (domain.Collection, error)
	Update(ctx context.Context, id string, in domain.Collection) (domain.Collection, error)
	Delete(ctx context.Context, id string) error
}

type OrderService interface {
	Create(ctx context.Context, in OrderCreateInput) (domain.Order, error)
	ListMine(ctx context.Context, userID string, page, pageSize int) ([]domain.Order, error)
}

type SaleService interface {
	RegisterSale(ctx context.Context, in SaleCreateInput) (domain.Sale, domain.Order, error)
}

type OrderCreateInput struct {
	UserID   string
	Currency string
	Items    []domain.OrderItem
	Notes    *string
}

type SaleCreateInput struct {
	UserID    string
	OrderID   string
	Method    domain.PaymentMethod
	Amount    float64
	Currency  string
	Reference *string
	SoldAt    *time.Time
}
