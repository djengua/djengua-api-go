package application

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/domain"
	"github.com/djengua/djengua-api-go/internal/ports"
)

type ProductsService struct {
	repo ports.ProductRepository
}

func NewProductsService(repo ports.ProductRepository) *ProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) List(ctx context.Context, filters ports.ProductFilters) ([]domain.Product, error) {
	return s.repo.ListProducts(ctx, filters)
}

func (s *ProductsService) Get(ctx context.Context, id string) (domain.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *ProductsService) Create(ctx context.Context, in domain.Product) (domain.Product, error) {
	return s.repo.CreateProduct(ctx, in)
}

func (s *ProductsService) Update(ctx context.Context, id string, in domain.Product) (domain.Product, error) {
	return s.repo.UpdateProductPut(ctx, id, in)
}

func (s *ProductsService) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}
