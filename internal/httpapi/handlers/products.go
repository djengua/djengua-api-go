package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/djengua/djengua-api-go/internal/db"
	"github.com/djengua/djengua-api-go/internal/domain"

	"github.com/go-chi/chi/v5"
)

type ProductsHandler struct {
	repo *db.Repository
}

func NewProductsHandler(repo *db.Repository) *ProductsHandler {
	return &ProductsHandler{repo: repo}
}

type productPayload struct {
	ID            string        `json:"id,omitempty"`
	SKU           string        `json:"sku"`
	Name          string        `json:"name"`
	Description   *string       `json:"description,omitempty"`
	Price         float64       `json:"price"`
	Currency      string        `json:"currency"`
	Status        domain.Status `json:"status"`
	CategoryID    string        `json:"category_id"`
	CollectionIDs []string      `json:"collection_ids,omitempty"`
	Images        []string      `json:"images,omitempty"`
	Stock         *int32        `json:"stock,omitempty"`
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page := parseIntQuery(r, "page", 1)
	pageSize := parseIntQuery(r, "page_size", 20)

	var filters db.ProductFilters
	filters.Page = page
	filters.PageSize = pageSize

	if v := q.Get("status"); v != "" {
		s := domain.Status(v)
		filters.Status = &s
	}
	if v := q.Get("category_id"); v != "" {
		filters.CategoryID = &v
	}
	if v := q.Get("collection_id"); v != "" {
		filters.CollectionID = &v
	}
	if v := q.Get("q"); v != "" {
		filters.Q = &v
	}
	// min_price/max_price optional
	// Keep parse tolerant
	if v := q.Get("min_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MinPrice = &f
		}
	}
	if v := q.Get("max_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MaxPrice = &f
		}
	}

	items, err := h.repo.ListProducts(r.Context(), filters)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.repo.GetProduct(r.Context(), id)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in productPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateProduct(in); err != "" {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.repo.CreateProduct(r.Context(), domain.Product{
		ID:            strings.TrimSpace(in.ID),
		SKU:           strings.TrimSpace(in.SKU),
		Name:          strings.TrimSpace(in.Name),
		Description:   in.Description,
		Price:         in.Price,
		Currency:      strings.ToUpper(strings.TrimSpace(in.Currency)),
		Status:        in.Status,
		CategoryID:    strings.TrimSpace(in.CategoryID),
		CollectionIDs: in.CollectionIDs,
		Images:        in.Images,
		Stock:         in.Stock,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *ProductsHandler) Put(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in productPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateProduct(in); err != "" {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.repo.UpdateProductPut(r.Context(), id, domain.Product{
		SKU:           strings.TrimSpace(in.SKU),
		Name:          strings.TrimSpace(in.Name),
		Description:   in.Description,
		Price:         in.Price,
		Currency:      strings.ToUpper(strings.TrimSpace(in.Currency)),
		Status:        in.Status,
		CategoryID:    strings.TrimSpace(in.CategoryID),
		CollectionIDs: in.CollectionIDs,
		Images:        in.Images,
		Stock:         in.Stock,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	current, err := h.repo.GetProduct(r.Context(), id)
	if mapDBErr(w, err) {
		return
	}
	var patch map[string]any
	if err := decodeJSON(r, &patch); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Apply patch fields (very small patch engine)
	applyString := func(key string, dst *string) {
		if v, ok := patch[key]; ok {
			if s, ok := v.(string); ok {
				*dst = s
			}
		}
	}
	applyFloat := func(key string, dst *float64) {
		if v, ok := patch[key]; ok {
			switch t := v.(type) {
			case float64:
				*dst = t
			case int:
				*dst = float64(t)
			}
		}
	}
	applyInt32Ptr := func(key string, dst **int32) {
		if v, ok := patch[key]; ok {
			if v == nil {
				*dst = nil
				return
			}
			if f, ok := v.(float64); ok {
				i := int32(f)
				*dst = &i
			}
		}
	}

	applyString("sku", &current.SKU)
	applyString("name", &current.Name)
	if v, ok := patch["description"]; ok {
		if v == nil {
			current.Description = nil
		} else if s, ok := v.(string); ok {
			current.Description = &s
		}
	}
	applyFloat("price", &current.Price)
	applyString("currency", &current.Currency)
	if v, ok := patch["status"]; ok {
		if s, ok := v.(string); ok {
			current.Status = domain.Status(s)
		}
	}
	applyString("category_id", &current.CategoryID)
	if v, ok := patch["images"]; ok {
		if arr, ok := v.([]any); ok {
			imgs := make([]string, 0, len(arr))
			for _, x := range arr {
				if s, ok := x.(string); ok {
					imgs = append(imgs, s)
				}
			}
			current.Images = imgs
		}
	}
	applyInt32Ptr("stock", &current.Stock)
	if v, ok := patch["collection_ids"]; ok {
		if arr, ok := v.([]any); ok {
			ids := make([]string, 0, len(arr))
			for _, x := range arr {
				if s, ok := x.(string); ok {
					ids = append(ids, s)
				}
			}
			current.CollectionIDs = ids
		}
	}

	payload := productPayload{
		SKU:           current.SKU,
		Name:          current.Name,
		Description:   current.Description,
		Price:         current.Price,
		Currency:      current.Currency,
		Status:        current.Status,
		CategoryID:    current.CategoryID,
		CollectionIDs: current.CollectionIDs,
		Images:        current.Images,
		Stock:         current.Stock,
	}
	if err := validateProduct(payload); err != "" {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.repo.UpdateProductPut(r.Context(), id, current)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteProduct(r.Context(), id); mapDBErr(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// func validateProduct(in productPayload) string {
// 	if strings.TrimSpace(in.SKU) == "" {
// 		return "sku is required"
// 	}
// 	if strings.TrimSpace(in.Name) == "" {
// 		return "name is required"
// 	}
// 	if in.Price <= 0 {
// 		return "price must be > 0"
// 	}
// 	if strings.TrimSpace(in.Currency) == "" {
// 		return "currency is required"
// 	}
// 	if strings.TrimSpace(string(in.Status)) == "" {
// 		return "status is required"
// 	}
// 	switch in.Status {
// 	case domain.ProductActive, domain.ProductInactive, domain.ProductDraft:
// 	default:
// 		return "invalid status"
// 	}
// 	if strings.TrimSpace(in.CategoryID) == "" {
// 		return "category_id is required"
// 	}
// 	return ""
// }
