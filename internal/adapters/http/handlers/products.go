// internal/adapters/http/handlers/products.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
	"github.com/djengua/djengua-api-go/internal/core/usecase/products"

	"github.com/go-chi/chi/v5"
)

type ProductsHandler struct {
	service *products.Service
}

func NewProductsHandler(service *products.Service) *ProductsHandler {
	return &ProductsHandler{service: service}
}

type productPayload struct {
	ID            string             `json:"id,omitempty"`
	SKU           string             `json:"sku"`
	Name          string             `json:"name"`
	Description   *string            `json:"description,omitempty"`
	Price         float64            `json:"price"`
	Cost          float64            `json:"cost"`
	Currency      string             `json:"currency"`
	Type          domain.ProductType `json:"type,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
	Status        domain.Status      `json:"status"`
	Size          domain.Size        `json:"size"`
	CategoryID    string             `json:"category_id"`
	CollectionIDs []string           `json:"collection_ids,omitempty"`
	Images        []string           `json:"images,omitempty"`
	Stock         *int32             `json:"stock,omitempty"`
	Attributes    map[string]any     `json:"attributes,omitempty"`
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page := parseIntQuery(r, "page", 1)
	pageSize := parseIntQuery(r, "page_size", 20)

	var filters ports.ProductFilters
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

	items, err := h.service.List(r.Context(), filters)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.service.Get(r.Context(), id)
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

	p, err := h.service.Create(r.Context(), domain.Product{
		ID:            strings.TrimSpace(in.ID),
		SKU:           strings.TrimSpace(in.SKU),
		Name:          strings.TrimSpace(in.Name),
		Description:   in.Description,
		Price:         in.Price,
		Cost:          in.Cost,
		Currency:      strings.ToUpper(strings.TrimSpace(in.Currency)),
		Type:          in.Type,
		Tags:          in.Tags,
		Status:        in.Status,
		Size:          in.Size,
		CategoryID:    strings.TrimSpace(in.CategoryID),
		CollectionIDs: in.CollectionIDs,
		Images:        in.Images,
		Stock:         in.Stock,
		Attributes:    in.Attributes,
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

	p, err := h.service.Update(r.Context(), id, domain.Product{
		SKU:           strings.TrimSpace(in.SKU),
		Name:          strings.TrimSpace(in.Name),
		Description:   in.Description,
		Price:         in.Price,
		Cost:          in.Cost,
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
	current, err := h.service.Get(r.Context(), id)
	if mapDBErr(w, err) {
		return
	}
	var patch map[string]any
	if err := decodeJSON(r, &patch); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

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
	applyFloat("cost", &current.Cost)
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
		Cost:          current.Cost,
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

	p, err := h.service.Update(r.Context(), id, current)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); mapDBErr(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
