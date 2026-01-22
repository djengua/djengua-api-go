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

type orderDoc struct {
	ID        string             `bson:"_id"`
	Number    string             `bson:"number"`
	UserID    string             `bson:"user_id"`
	Status    domain.OrderStatus `bson:"status"`
	Currency  string             `bson:"currency"`
	Items     []domain.OrderItem `bson:"items"`
	Subtotal  float64            `bson:"subtotal"`
	Tax       float64            `bson:"tax"`
	Total     float64            `bson:"total"`
	Notes     *string            `bson:"notes,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (d orderDoc) toDomain() domain.Order {
	return domain.Order{
		ID: d.ID, Number: d.Number, UserID: d.UserID, Status: d.Status,
		Currency: d.Currency, Items: d.Items, Subtotal: d.Subtotal, Tax: d.Tax, Total: d.Total,
		Notes: d.Notes, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt, DeletedAt: d.DeletedAt,
	}
}

func (r *Repository) CreateOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(o.ID)
	if id == "" {
		id = uuid.NewString()
	}

	doc := orderDoc{
		ID:        id,
		Number:    o.Number,
		UserID:    o.UserID,
		Status:    o.Status,
		Currency:  o.Currency,
		Items:     o.Items,
		Subtotal:  o.Subtotal,
		Tax:       o.Tax,
		Total:     o.Total,
		Notes:     o.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.orders.InsertOne(ctx, doc)
	if err != nil {
		return domain.Order{}, err
	}
	return r.GetOrder(ctx, id)
}

func (r *Repository) GetOrder(ctx context.Context, id string) (domain.Order, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	var d orderDoc
	err := r.orders.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return domain.Order{}, ports.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) ListOrdersByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.Order, error) {
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
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(limit)

	cur, err := r.orders.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]domain.Order, 0)
	for cur.Next(ctx) {
		var d orderDoc
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

func (r *Repository) UpdateOrderStatus(ctx context.Context, id string, st domain.OrderStatus) (domain.Order, error) {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	res, err := r.orders.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"status":     st,
		"updated_at": now,
	}})
	if err != nil {
		return domain.Order{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Order{}, ports.ErrNotFound
	}
	return r.GetOrder(ctx, id)
}
