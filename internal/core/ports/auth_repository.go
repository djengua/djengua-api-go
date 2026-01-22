package ports

import (
	"context"

	"github.com/djengua/djengua-api-go/internal/core/domain"
)

type AuthUserRepository interface {
	CreateUser(ctx context.Context, u domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
