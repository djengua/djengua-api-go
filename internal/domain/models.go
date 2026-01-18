// internal/domain/models.go
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

type Product struct {
	ID            string     `json:"id"`
	SKU           string     `json:"sku"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Price         float64    `json:"price"`
	Cost          float64    `json:"cost"`
	Currency      string     `json:"currency"`
	Status        Status     `json:"status"`
	Size          Size       `json:"size"`
	CategoryID    string     `json:"category_id"`
	CollectionIDs []string   `json:"collection_ids,omitempty"`
	Images        []string   `json:"images,omitempty"`
	Stock         *int32     `json:"stock,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
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
