// internal/adapters/http/handlers/sales.go
package handlers

import (
	"net/http"
	"time"

	"github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/usecase/sales"
)

type SalesHandler struct {
	uc *sales.Service
}

func NewSalesHandler(uc *sales.Service) *SalesHandler {
	return &SalesHandler{uc: uc}
}

type createSaleReq struct {
	OrderID   string               `json:"order_id"`
	Method    domain.PaymentMethod `json:"method"`
	Amount    float64              `json:"amount"`
	Currency  string               `json:"currency"`
	Reference *string              `json:"reference"`
	SoldAt    *time.Time           `json:"sold_at"`
}

func (h *SalesHandler) Register(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing auth")
		return
	}

	var in createSaleReq
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sale, ord, err := h.uc.RegisterSale(r.Context(), sales.CreateInput{
		UserID: uid, OrderID: in.OrderID, Method: in.Method,
		Amount: in.Amount, Currency: in.Currency, Reference: in.Reference, SoldAt: in.SoldAt,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"sale":  sale,
		"order": ord,
	})
}
