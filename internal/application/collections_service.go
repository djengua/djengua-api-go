package application

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/domain"
	"github.com/djengua/djengua-api-go/internal/ports"
)

type CollectionsService struct {
	repo ports.CollectionRepository
}

func NewCollectionsService(repo ports.CollectionRepository) *CollectionsService {
	return &CollectionsService{repo: repo}
}

func (s *CollectionsService) List(ctx context.Context, page, pageSize int) ([]domain.Collection, error) {
	return s.repo.ListCollections(ctx, page, pageSize)
}

func (s *CollectionsService) Get(ctx context.Context, id string) (domain.Collection, error) {
	return s.repo.GetCollection(ctx, id)
}

func (s *CollectionsService) Create(ctx context.Context, in domain.Collection) (domain.Collection, error) {
	return s.repo.CreateCollection(ctx, in)
}

func (s *CollectionsService) Update(ctx context.Context, id string, in domain.Collection) (domain.Collection, error) {
	return s.repo.UpdateCollectionPut(ctx, id, in)
}

func (s *CollectionsService) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteCollection(ctx, id)
}
