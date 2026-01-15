package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db          *mongo.Database
	products    *mongo.Collection
	categories  *mongo.Collection
	collections *mongo.Collection
}

func NewRepository(database *mongo.Database) *Repository {
	return &Repository{
		db:          database,
		products:    database.Collection("products"),
		categories:  database.Collection("categories"),
		collections: database.Collection("collections"),
	}
}

// --- Mongo documents ---

type productDoc struct {
	ID            string               `bson:"_id"`
	SKU           string               `bson:"sku"`
	Name          string               `bson:"name"`
	Description   *string              `bson:"description,omitempty"`
	Price         float64              `bson:"price"`
	Currency      string               `bson:"currency"`
	Status        domain.ProductStatus `bson:"status"`
	CategoryID    string               `bson:"category_id"`
	CollectionIDs []string             `bson:"collection_ids,omitempty"`
	Images        []string             `bson:"images,omitempty"`
	Stock         *int32               `bson:"stock,omitempty"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
	DeletedAt     *time.Time           `bson:"deleted_at,omitempty"`
}

type categoryDoc struct {
	ID          string                `bson:"_id"`
	Name        string                `bson:"name"`
	Slug        string                `bson:"slug"`
	Description *string               `bson:"description,omitempty"`
	Status      domain.CategoryStatus `bson:"status"`
	CreatedAt   time.Time             `bson:"created_at"`
	UpdatedAt   time.Time             `bson:"updated_at"`
	DeletedAt   *time.Time            `bson:"deleted_at,omitempty"`
}

type collectionDoc struct {
	ID          string                  `bson:"_id"`
	Name        string                  `bson:"name"`
	Slug        string                  `bson:"slug"`
	Description *string                 `bson:"description,omitempty"`
	Status      domain.CollectionStatus `bson:"status"`
	StartDate   *time.Time              `bson:"start_date,omitempty"`
	EndDate     *time.Time              `bson:"end_date,omitempty"`
	CreatedAt   time.Time               `bson:"created_at"`
	UpdatedAt   time.Time               `bson:"updated_at"`
	DeletedAt   *time.Time              `bson:"deleted_at,omitempty"`
}

func (d productDoc) toDomain() domain.Product {
	return domain.Product{
		ID:            d.ID,
		SKU:           d.SKU,
		Name:          d.Name,
		Description:   d.Description,
		Price:         d.Price,
		Currency:      d.Currency,
		Status:        d.Status,
		CategoryID:    d.CategoryID,
		CollectionIDs: d.CollectionIDs,
		Images:        d.Images,
		Stock:         d.Stock,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
		DeletedAt:     d.DeletedAt,
	}
}

func (d categoryDoc) toDomain() domain.Category {
	return domain.Category{
		ID:          d.ID,
		Name:        d.Name,
		Slug:        d.Slug,
		Description: d.Description,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
	}
}

func (d collectionDoc) toDomain() domain.Collection {
	return domain.Collection{
		ID:          d.ID,
		Name:        d.Name,
		Slug:        d.Slug,
		Description: d.Description,
		Status:      d.Status,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
	}
}

func notDeletedFilter() bson.M {
	// Because DeletedAt is omitempty, it won't exist unless soft-deleted.
	return bson.M{"deleted_at": bson.M{"$exists": false}}
}

// --- Products ---

type ProductFilters struct {
	Status       *domain.ProductStatus
	CategoryID   *string
	CollectionID *string
	MinPrice     *float64
	MaxPrice     *float64
	Q            *string
	Page         int
	PageSize     int
}

func (f *ProductFilters) normalize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
}

func (r *Repository) CreateProduct(ctx context.Context, in domain.Product) (domain.Product, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}

	doc := productDoc{
		ID:            id,
		SKU:           in.SKU,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		Currency:      in.Currency,
		Status:        in.Status,
		CategoryID:    in.CategoryID,
		CollectionIDs: in.CollectionIDs,
		Images:        in.Images,
		Stock:         in.Stock,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err := r.products.InsertOne(ctx, doc)
	if err != nil {
		return domain.Product{}, err
	}
	return r.GetProduct(ctx, id)
}

func (r *Repository) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	var d productDoc
	err := r.products.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Product{}, ErrNotFound
	}
	if err != nil {
		return domain.Product{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) ListProducts(ctx context.Context, f ProductFilters) ([]domain.Product, error) {
	f.normalize()

	filter := notDeletedFilter()

	if f.Status != nil {
		filter["status"] = *f.Status
	}
	if f.CategoryID != nil {
		filter["category_id"] = *f.CategoryID
	}
	if f.CollectionID != nil {
		// collection_ids is an array of strings; matching a scalar matches any array element.
		filter["collection_ids"] = *f.CollectionID
	}
	if f.MinPrice != nil || f.MaxPrice != nil {
		rangeQ := bson.M{}
		if f.MinPrice != nil {
			rangeQ["$gte"] = *f.MinPrice
		}
		if f.MaxPrice != nil {
			rangeQ["$lte"] = *f.MaxPrice
		}
		filter["price"] = rangeQ
	}
	if f.Q != nil && strings.TrimSpace(*f.Q) != "" {
		q := strings.TrimSpace(*f.Q)
		r := primitive.Regex{Pattern: q, Options: "i"}
		filter["$or"] = []bson.M{{"name": r}, {"sku": r}}
	}

	skip := int64((f.Page - 1) * f.PageSize)
	limit := int64(f.PageSize)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(limit)

	cur, err := r.products.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]domain.Product, 0)
	for cur.Next(ctx) {
		var d productDoc
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

func (r *Repository) UpdateProductPut(ctx context.Context, id string, in domain.Product) (domain.Product, error) {
	now := time.Now().UTC()

	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	update := bson.M{"$set": bson.M{
		"sku":            in.SKU,
		"name":           in.Name,
		"description":    in.Description,
		"price":          in.Price,
		"currency":       in.Currency,
		"status":         in.Status,
		"category_id":    in.CategoryID,
		"collection_ids": in.CollectionIDs,
		"images":         in.Images,
		"stock":          in.Stock,
		"updated_at":     now,
	}}

	res, err := r.products.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Product{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Product{}, ErrNotFound
	}
	return r.GetProduct(ctx, id)
}

func (r *Repository) DeleteProduct(ctx context.Context, id string) error {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	res, err := r.products.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Categories ---

func (r *Repository) CreateCategory(ctx context.Context, in domain.Category) (domain.Category, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}
	doc := categoryDoc{
		ID:          id,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		Status:      in.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.categories.InsertOne(ctx, doc)
	if err != nil {
		return domain.Category{}, err
	}
	return r.GetCategory(ctx, id)
}

func (r *Repository) GetCategory(ctx context.Context, id string) (domain.Category, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	var d categoryDoc
	err := r.categories.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Category{}, ErrNotFound
	}
	if err != nil {
		return domain.Category{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(limit)

	cur, err := r.categories.Find(ctx, notDeletedFilter(), opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]domain.Category, 0)
	for cur.Next(ctx) {
		var d categoryDoc
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

func (r *Repository) UpdateCategoryPut(ctx context.Context, id string, in domain.Category) (domain.Category, error) {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	update := bson.M{"$set": bson.M{
		"name":        in.Name,
		"slug":        in.Slug,
		"description": in.Description,
		"status":      in.Status,
		"updated_at":  now,
	}}
	res, err := r.categories.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Category{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Category{}, ErrNotFound
	}
	return r.GetCategory(ctx, id)
}

func (r *Repository) DeleteCategory(ctx context.Context, id string) error {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	res, err := r.categories.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Collections ---

func (r *Repository) CreateCollection(ctx context.Context, in domain.Collection) (domain.Collection, error) {
	now := time.Now().UTC()
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}
	doc := collectionDoc{
		ID:          id,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.collections.InsertOne(ctx, doc)
	if err != nil {
		return domain.Collection{}, err
	}
	return r.GetCollection(ctx, id)
}

func (r *Repository) GetCollection(ctx context.Context, id string) (domain.Collection, error) {
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	var d collectionDoc
	err := r.collections.FindOne(ctx, filter).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Collection{}, ErrNotFound
	}
	if err != nil {
		return domain.Collection{}, err
	}
	return d.toDomain(), nil
}

func (r *Repository) ListCollections(ctx context.Context, page, pageSize int) ([]domain.Collection, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(skip).SetLimit(limit)

	cur, err := r.collections.Find(ctx, notDeletedFilter(), opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]domain.Collection, 0)
	for cur.Next(ctx) {
		var d collectionDoc
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

func (r *Repository) UpdateCollectionPut(ctx context.Context, id string, in domain.Collection) (domain.Collection, error) {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	update := bson.M{"$set": bson.M{
		"name":        in.Name,
		"slug":        in.Slug,
		"description": in.Description,
		"status":      in.Status,
		"start_date":  in.StartDate,
		"end_date":    in.EndDate,
		"updated_at":  now,
	}}
	res, err := r.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Collection{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Collection{}, ErrNotFound
	}
	return r.GetCollection(ctx, id)
}

func (r *Repository) DeleteCollection(ctx context.Context, id string) error {
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}
	res, err := r.collections.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
