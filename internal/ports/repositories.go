package ports

import (
	"context"
	"errors"

	"github.com/djengua/djengua-api-go/internal/domain"
)

var ErrNotFound = errors.New("not found")

type ProductFilters struct {
	Status       *domain.Status
	CategoryID   *string
	CollectionID *string
	MinPrice     *float64
	MaxPrice     *float64
	Q            *string
	Page         int
	PageSize     int
}

func (f *ProductFilters) Normalize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, in domain.Product) (domain.Product, error)
	GetProduct(ctx context.Context, id string) (domain.Product, error)
	ListProducts(ctx context.Context, f ProductFilters) ([]domain.Product, error)
	UpdateProductPut(ctx context.Context, id string, in domain.Product) (domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, in domain.Category) (domain.Category, error)
	GetCategory(ctx context.Context, id string) (domain.Category, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	UpdateCategoryPut(ctx context.Context, id string, in domain.Category) (domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CollectionRepository interface {
	CreateCollection(ctx context.Context, in domain.Collection) (domain.Collection, error)
	GetCollection(ctx context.Context, id string) (domain.Collection, error)
	ListCollections(ctx context.Context, page, pageSize int) ([]domain.Collection, error)
	UpdateCollectionPut(ctx context.Context, id string, in domain.Collection) (domain.Collection, error)
	DeleteCollection(ctx context.Context, id string) error
}
