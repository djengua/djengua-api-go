// internal/adapters/mongo/repository_users.go
package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type userDoc struct {
	ID           string          `bson:"_id"`
	Email        string          `bson:"email"`
	Name         string          `bson:"name"`
	Role         domain.UserRole `bson:"role"`
	PasswordHash string          `bson:"password_hash"`
	CreatedAt    time.Time       `bson:"created_at"`
	UpdatedAt    time.Time       `bson:"updated_at"`
	DeletedAt    *time.Time      `bson:"deleted_at,omitempty"`
}

func (d userDoc) toDomain() domain.User {
	return domain.User{
		ID:           d.ID,
		Email:        d.Email,
		Name:         d.Name,
		Role:         d.Role,
		PasswordHash: d.PasswordHash,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		DeletedAt:    d.DeletedAt,
	}
}

func (r *Repository) CreateUser(ctx context.Context, u domain.User) (domain.User, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(u.ID)
	if id == "" {
		id = uuid.NewString()
	}

	doc := userDoc{
		ID:           id,
		Email:        strings.ToLower(strings.TrimSpace(u.Email)),
		Name:         strings.TrimSpace(u.Name),
		Role:         u.Role,
		PasswordHash: u.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err := r.users.InsertOne(ctx, doc)
	if err != nil {
		// Duplicate key error (email) requiere índice unique en "email"
		var we mongodriver.WriteException
		if errors.As(err, &we) {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					return domain.User{}, ports.ErrConflict
				}
			}
		}
		return domain.User{}, err
	}

	return r.GetUserByID(ctx, id)
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	var d userDoc
	err := r.users.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return domain.User{}, ports.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	filter := bson.M{"email": strings.ToLower(strings.TrimSpace(email))}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	var d userDoc
	err := r.users.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return domain.User{}, ports.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return d.toDomain(), nil
}
