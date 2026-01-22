package ports

import (
	"context"
	"errors"

	"github.com/djengua/djengua-api-go/internal/core/domain"
)

var ErrConflict = errors.New("conflict")

type UserRepository interface {
	AuthUserRepository
	ListUsers(ctx context.Context, page, pageSize int) ([]domain.User, error) // si la quieres
	// otros métodos futuros...
}
