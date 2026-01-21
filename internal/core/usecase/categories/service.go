package categories

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	repo ports.CategoryRepository
}

func NewService(repo ports.CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx, page, pageSize)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *Service) Create(ctx context.Context, in domain.Category) (domain.Category, error) {
	return s.repo.CreateCategory(ctx, in)
}

func (s *Service) Update(ctx context.Context, id string, in domain.Category) (domain.Category, error) {
	return s.repo.UpdateCategoryPut(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteCategory(ctx, id)
}
