package domain

import "time"

type Status string
type Size string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Draft    Status = "draft"
)

const (
	XL Size = "XL"
	L  Size = "L"
	M  Size = "M"
	S  Size = "S"
	XS Size = "XS"
)

type ProductType string

const (
	ProductPhysical ProductType = "physical"
	ProductDigital  ProductType = "digital"
	ProductService  ProductType = "service"
)

type ProductOption struct {
	Name   string   `json:"name"`             // e.g. "color", "size"
	Values []string `json:"values,omitempty"` // e.g. ["S","M","L"] or ["red","blue"]
}

type Product struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	// Money
	Price    float64 `json:"price"` // venta al público
	Cost     float64 `json:"cost"`  // costo del producto
	Currency string  `json:"currency"`

	// Generic classification
	Type   ProductType `json:"type,omitempty"`
	Tags   []string    `json:"tags,omitempty"`
	Status Status      `json:"status"`

	// Variants/basic
	Size Size `json:"size"`

	// Relations
	CategoryID    string   `json:"category_id"`
	CollectionIDs []string `json:"collection_ids,omitempty"`

	// Media & stock
	Images []string `json:"images,omitempty"`
	Stock  *int32   `json:"stock,omitempty"`

	// Generic attributes for arbitrary products (ropa, juguetes, regalos, personalizados, etc.)
	Attributes map[string]any `json:"attributes,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Category struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Collection struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	Status      Status     `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
