// internal/core/usecase/users/service.go
package users

import (
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type Service struct {
	repo ports.UserRepository
}

func NewService(repo ports.UserRepository) *Service {
	return &Service{repo: repo}
}

// func (s *Service) List(ctx context.Context, filters ports.UsersFilter) ([]domain.User, error) {
// 	return s.repo.ListUsers(ctx, filters)
// }

// func (s *Service) Get(ctx context.Context, id string) (domain.Product, error) {
// 	return s.repo.GetProduct(ctx, id)
// }

// func (s *Service) Create(ctx context.Context, in domain.Product) (domain.Product, error) {
// 	return s.repo.CreateProduct(ctx, in)
// }

// func (s *Service) Update(ctx context.Context, id string, in domain.Product) (domain.Product, error) {
// 	return s.repo.UpdateProductPut(ctx, id, in)
// }

// func (s *Service) Delete(ctx context.Context, id string) error {
// 	return s.repo.DeleteProduct(ctx, id)
// }
