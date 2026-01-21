package handlers

import (
	"regexp"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateProduct(in productPayload) string {
	if strings.TrimSpace(in.SKU) == "" {
		return "sku is required"
	}
	if strings.TrimSpace(in.Name) == "" {
		return "name is required"
	}

	// price = venta
	if in.Price <= 0 {
		return "price must be > 0"
	}
	// cost = costo (puede ser 0 si regalas/muestras, pero normalmente >0)
	if in.Cost < 0 {
		return "cost must be >= 0"
	}
	// opcional: asegurar margen (si quieres)
	// if in.Cost > 0 && in.Price < in.Cost { return "price should be >= cost" }

	cur := strings.TrimSpace(in.Currency)
	if len(cur) != 3 {
		return "currency must be a 3-letter code (e.g. USD)"
	}

	switch in.Status {
	case domain.Active, domain.Inactive, domain.Draft:
	default:
		return "status must be active|inactive|draft"
	}

	// type opcional, si viene debe ser válido
	if in.Type != "" {
		switch in.Type {
		case domain.ProductPhysical, domain.ProductDigital, domain.ProductService:
		default:
			return "type must be physical|digital|service"
		}
	}

	// size opcional (si quieres permitir vacío), si quieres forzar: quita el if
	if in.Size != "" {
		switch in.Size {
		case domain.XS, domain.S, domain.M, domain.L, domain.XL:
		default:
			return "size must be XS|S|M|L|XL"
		}
	}

	if strings.TrimSpace(in.CategoryID) == "" {
		return "category_id is required"
	}

	return ""
}

type categoryPayload struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description *string       `json:"description,omitempty"`
	Status      domain.Status `json:"status"`
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
	case domain.Active, domain.Inactive:
	default:
		return "status must be active|inactive"
	}
	return ""
}

type collectionPayload struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description *string       `json:"description,omitempty"`
	Status      domain.Status `json:"status"`
	StartDate   *time.Time    `json:"start_date,omitempty"`
	EndDate     *time.Time    `json:"end_date,omitempty"`
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
	case domain.Active, domain.Inactive:
	default:
		return "status must be active|inactive"
	}
	if in.StartDate != nil && in.EndDate != nil && in.StartDate.After(*in.EndDate) {
		return "start_date must be <= end_date"
	}
	return ""
}
