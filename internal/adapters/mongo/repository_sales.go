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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type saleDoc struct {
	ID        string               `bson:"_id"`
	OrderID   string               `bson:"order_id"`
	UserID    string               `bson:"user_id"`
	Method    domain.PaymentMethod `bson:"method"`
	Amount    float64              `bson:"amount"`
	Currency  string               `bson:"currency"`
	Reference *string              `bson:"reference,omitempty"`
	SoldAt    time.Time            `bson:"sold_at"`
	CreatedAt time.Time            `bson:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at"`
	DeletedAt *time.Time           `bson:"deleted_at,omitempty"`
}

func (d saleDoc) toDomain() domain.Sale {
	return domain.Sale{
		ID: d.ID, OrderID: d.OrderID, UserID: d.UserID, Method: d.Method,
		Amount: d.Amount, Currency: d.Currency, Reference: d.Reference, SoldAt: d.SoldAt,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt, DeletedAt: d.DeletedAt,
	}
}

func (r *Repository) CreateSale(ctx context.Context, s domain.Sale) (domain.Sale, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(s.ID)
	if id == "" {
		id = uuid.NewString()
	}
	soldAt := s.SoldAt
	if soldAt.IsZero() {
		soldAt = now
	}

	doc := saleDoc{
		ID: id, OrderID: s.OrderID, UserID: s.UserID, Method: s.Method,
		Amount: s.Amount, Currency: s.Currency, Reference: s.Reference,
		SoldAt: soldAt, CreatedAt: now, UpdatedAt: now,
	}

	_, err := r.sales.InsertOne(ctx, doc)
	if err != nil {
		return domain.Sale{}, err
	}
	return r.GetSale(ctx, id)
}

func (r *Repository) GetSale(ctx context.Context, id string) (domain.Sale, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	var d saleDoc
	err := r.sales.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return domain.Sale{}, ports.ErrNotFound
	}
	if err != nil {
		return domain.Sale{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) ListSalesByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Sale, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	filter := notDeletedFilter()
	filter["user_id"] = userID

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSort(bson.D{{Key: "sold_at", Value: -1}}).SetSkip(skip).SetLimit(limit)

	cur, err := r.sales.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]domain.Sale, 0)
	for cur.Next(ctx) {
		var d saleDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		out = append(out, d.toDomain())
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
