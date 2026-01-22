// internal/core/domain/orders.go
package domain

import "time"

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "draft"
	OrderStatusPending   OrderStatus = "pending" // creada, esperando pago
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Qty       int32   `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
	LineTotal float64 `json:"line_total"`
}

type Order struct {
	ID        string      `json:"id"`
	Number    string      `json:"number"` // folio
	UserID    string      `json:"user_id"`
	Status    OrderStatus `json:"status"`
	Currency  string      `json:"currency"`
	Items     []OrderItem `json:"items"`
	Subtotal  float64     `json:"subtotal"`
	Tax       float64     `json:"tax"`
	Total     float64     `json:"total"`
	Notes     *string     `json:"notes,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty"`
}

type PaymentMethod string

const (
	PayCash     PaymentMethod = "cash"
	PayCard     PaymentMethod = "card"
	PayTransfer PaymentMethod = "transfer"
	PayOther    PaymentMethod = "other"
)

type Sale struct {
	ID        string        `json:"id"`
	OrderID   string        `json:"order_id"`
	UserID    string        `json:"user_id"`
	Method    PaymentMethod `json:"method"`
	Amount    float64       `json:"amount"`
	Currency  string        `json:"currency"`
	Reference *string       `json:"reference,omitempty"` // folio banco/stripe/etc
	SoldAt    time.Time     `json:"sold_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at,omitempty"`
}
