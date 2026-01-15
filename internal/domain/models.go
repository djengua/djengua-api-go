package domain

import "time"

type ProductStatus string

type CategoryStatus string

type CollectionStatus string

const (
	ProductActive   ProductStatus = "active"
	ProductInactive ProductStatus = "inactive"
	ProductDraft    ProductStatus = "draft"
)

const (
	CategoryActive   CategoryStatus = "active"
	CategoryInactive CategoryStatus = "inactive"
)

const (
	CollectionActive   CollectionStatus = "active"
	CollectionInactive CollectionStatus = "inactive"
)

type Product struct {
	ID            string        `json:"id"`
	SKU           string        `json:"sku"`
	Name          string        `json:"name"`
	Description   *string       `json:"description,omitempty"`
	Price         float64       `json:"price"`
	Currency      string        `json:"currency"`
	Status        ProductStatus `json:"status"`
	CategoryID    string        `json:"category_id"`
	CollectionIDs []string      `json:"collection_ids,omitempty"`
	Images        []string      `json:"images,omitempty"`
	Stock         *int32        `json:"stock,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at,omitempty"`
}

type Category struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description *string        `json:"description,omitempty"`
	Status      CategoryStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

type Collection struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description *string          `json:"description,omitempty"`
	Status      CollectionStatus `json:"status"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}
