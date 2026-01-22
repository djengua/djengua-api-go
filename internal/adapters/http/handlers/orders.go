// internal/adapters/http/handlers/orders.go
package handlers

import (
	"net/http"

	"github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type OrdersHandler struct {
	uc ports.OrderService
}

func NewOrdersHandler(uc ports.OrderService) *OrdersHandler {
	return &OrdersHandler{uc: uc}
}

type createOrderReq struct {
	Currency string             `json:"currency"`
	Items    []domain.OrderItem `json:"items"`
	Notes    *string            `json:"notes"`
}

func (h *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing auth")
		return
	}

	var in createOrderReq
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.uc.Create(r.Context(), ports.OrderCreateInput{
		UserID: uid, Currency: in.Currency, Items: in.Items, Notes: in.Notes,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, o)
}

func (h *OrdersHandler) MyOrders(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing auth")
		return
	}
	// usa defaults simples; si ya tienes helpers de paginación, aplícalos
	out, err := h.uc.ListMine(r.Context(), uid, 1, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, out)
}
