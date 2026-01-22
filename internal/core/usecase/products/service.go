// internal/core/usecase/products/service.go
package products

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	repo ports.ProductRepository
}

func NewService(repo ports.ProductRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filters ports.ProductFilters) ([]domain.Product, error) {
	return s.repo.ListProducts(ctx, filters)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *Service) Create(ctx context.Context, in domain.Product) (domain.Product, error) {
	return s.repo.CreateProduct(ctx, in)
}

func (s *Service) Update(ctx context.Context, id string, in domain.Product) (domain.Product, error) {
	return s.repo.UpdateProductPut(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}
