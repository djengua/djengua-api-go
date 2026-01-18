package application

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/domain"
	"github.com/djengua/djengua-api-go/internal/ports"
)

type CategoriesService struct {
	repo ports.CategoryRepository
}

func NewCategoriesService(repo ports.CategoryRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) List(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx, page, pageSize)
}

func (s *CategoriesService) Get(ctx context.Context, id string) (domain.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *CategoriesService) Create(ctx context.Context, in domain.Category) (domain.Category, error) {
	return s.repo.CreateCategory(ctx, in)
}

func (s *CategoriesService) Update(ctx context.Context, id string, in domain.Category) (domain.Category, error) {
	return s.repo.UpdateCategoryPut(ctx, id, in)
}

func (s *CategoriesService) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteCategory(ctx, id)
}
