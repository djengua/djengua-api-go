// internal/adapters/http/handlers/collections.go
package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/usecase/collections"

	"github.com/go-chi/chi/v5"
)

type CollectionsHandler struct {
	service *collections.Service
}

func NewCollectionsHandler(service *collections.Service) *CollectionsHandler {
	return &CollectionsHandler{service: service}
}

func (h *CollectionsHandler) List(w http.ResponseWriter, r *http.Request) {
	page := parseIntQuery(r, "page", 1)
	pageSize := parseIntQuery(r, "page_size", 50)
	items, err := h.service.List(r.Context(), page, pageSize)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CollectionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.service.Get(r.Context(), id)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CollectionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in collectionPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := validateCollection(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Create(r.Context(), domain.Collection{
		ID:          strings.TrimSpace(in.ID),
		Name:        strings.TrimSpace(in.Name),
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *CollectionsHandler) Put(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in collectionPayload
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := validateCollection(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Update(r.Context(), id, domain.Collection{
		Name:        strings.TrimSpace(in.Name),
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CollectionsHandler) Patch(w http.ResponseWriter, r *http.Request) {
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
			current.Status = domain.Status(s)
		}
	}
	if v, ok := patch["start_date"]; ok {
		if v == nil {
			current.StartDate = nil
		} else if s, ok := v.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				u := t.UTC()
				current.StartDate = &u
			}
		}
	}
	if v, ok := patch["end_date"]; ok {
		if v == nil {
			current.EndDate = nil
		} else if s, ok := v.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				u := t.UTC()
				current.EndDate = &u
			}
		}
	}

	payload := collectionPayload{Name: current.Name, Slug: current.Slug, Description: current.Description, Status: current.Status, StartDate: current.StartDate, EndDate: current.EndDate}
	if msg := validateCollection(payload); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	c, err := h.service.Update(r.Context(), id, current)
	if mapDBErr(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CollectionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); mapDBErr(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
