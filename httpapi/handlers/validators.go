package handlers

import (
	"regexp"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/domain"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateProduct(in productPayload) string {
	if strings.TrimSpace(in.SKU) == "" {
		return "sku is required"
	}
	if strings.TrimSpace(in.Name) == "" {
		return "name is required"
	}
	if in.Price <= 0 {
		return "price must be > 0"
	}
	cur := strings.TrimSpace(in.Currency)
	if len(cur) != 3 {
		return "currency must be a 3-letter code (e.g. USD)"
	}
	switch in.Status {
	case domain.ProductActive, domain.ProductInactive, domain.ProductDraft:
		// ok
	default:
		return "status must be active|inactive|draft"
	}
	if strings.TrimSpace(in.CategoryID) == "" {
		return "category_id is required"
	}
	return ""
}

type categoryPayload struct {
	ID          string                `json:"id,omitempty"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Description *string               `json:"description,omitempty"`
	Status      domain.CategoryStatus `json:"status"`
}

func validateCategory(in categoryPayload) string {
	if strings.TrimSpace(in.Name) == "" {
		return "name is required"
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		return "slug is required"
	}
	if !slugRe.MatchString(slug) {
		return "slug must be URL-friendly (lowercase letters, numbers, hyphens)"
	}
	switch in.Status {
	case domain.CategoryActive, domain.CategoryInactive:
		// ok
	default:
		return "status must be active|inactive"
	}
	return ""
}

type collectionPayload struct {
	ID          string                  `json:"id,omitempty"`
	Name        string                  `json:"name"`
	Slug        string                  `json:"slug"`
	Description *string                 `json:"description,omitempty"`
	Status      domain.CollectionStatus `json:"status"`
	StartDate   *time.Time              `json:"start_date,omitempty"`
	EndDate     *time.Time              `json:"end_date,omitempty"`
}

func validateCollection(in collectionPayload) string {
	if strings.TrimSpace(in.Name) == "" {
		return "name is required"
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		return "slug is required"
	}
	if !slugRe.MatchString(slug) {
		return "slug must be URL-friendly (lowercase letters, numbers, hyphens)"
	}
	switch in.Status {
	case domain.CollectionActive, domain.CollectionInactive:
		// ok
	default:
		return "status must be active|inactive"
	}
	if in.StartDate != nil && in.EndDate != nil && in.StartDate.After(*in.EndDate) {
		return "start_date must be <= end_date"
	}
	return ""
}
