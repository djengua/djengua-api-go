// internal/core/usecase/collections/service.go
package collections

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	repo ports.CollectionRepository
}

func NewService(repo ports.CollectionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]domain.Collection, error) {
	return s.repo.ListCollections(ctx, page, pageSize)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Collection, error) {
	return s.repo.GetCollection(ctx, id)
}

func (s *Service) Create(ctx context.Context, in domain.Collection) (domain.Collection, error) {
	return s.repo.CreateCollection(ctx, in)
}

func (s *Service) Update(ctx context.Context, id string, in domain.Collection) (domain.Collection, error) {
	return s.repo.UpdateCollectionPut(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteCollection(ctx, id)
}
