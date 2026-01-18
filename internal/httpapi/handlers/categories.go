package handlers

import (
	"net/http"
	"strings"

	"github.com/djengua/djengua-api-go/internal/application"
	"github.com/djengua/djengua-api-go/internal/domain"

	"github.com/go-chi/chi/v5"
)

type CategoriesHandler struct {
	service *application.CategoriesService
}

func NewCategoriesHandler(service *application.CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{service: service}
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	page := parseIntQuery(r, "page", 1)
	pageSize := parseIntQuery(r, "page_size", 50)
	items, err := h.service.List(r.Context(), page, pageSize)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.service.Get(r.Context(), id)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in categoryPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := validateCategory(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Create(r.Context(), domain.Category{
		ID:          strings.TrimSpace(in.ID),
		Name:        strings.TrimSpace(in.Name),
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Status:      in.Status,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *CategoriesHandler) Put(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in categoryPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := validateCategory(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Update(r.Context(), id, domain.Category{
		Name:        strings.TrimSpace(in.Name),
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Status:      in.Status,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CategoriesHandler) Patch(w http.ResponseWriter, r *http.Request) {
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

	if v, ok := patch["name"]; ok {
		if s, ok := v.(string); ok {
			current.Name = s
		}
	}
	if v, ok := patch["slug"]; ok {
		if s, ok := v.(string); ok {
			current.Slug = s
		}
	}
	if v, ok := patch["description"]; ok {
		if v == nil {
			current.Description = nil
		} else if s, ok := v.(string); ok {
			current.Description = &s
		}
	}
	if v, ok := patch["status"]; ok {
		if s, ok := v.(string); ok {
			current.Status = domain.CategoryStatus(s)
		}
	}

	payload := categoryPayload{Name: current.Name, Slug: current.Slug, Description: current.Description, Status: current.Status}
	if msg := validateCategory(payload); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Update(r.Context(), id, current)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); mapDBErr(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
