package utils

import (
	"fmt"
	"math"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// PagedResponse represents a paginated response structure
type PagedResponse[T any] struct {
	Data  []T      `json:"data"`
	Meta  MetaData `json:"meta"`
	Links *Links   `json:"links,omitempty"`
}

// MetaData contains pagination metadata
type MetaData struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Links contains pagination links
type Links struct {
	Self  string `json:"self"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	First string `json:"first"`
	Last  string `json:"last"`
}

// CreatePagedResponse generates a PagedResponse for any slice of data
func CreatePagedResponse[T any](c *fiber.Ctx, data []T, page, pageSize int, total int64) PagedResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	meta := MetaData{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	links, err := generatePaginationLinks(c, page, pageSize, totalPages)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate pagination links")
	}

	return PagedResponse[T]{
		Data:  data,
		Meta:  meta,
		Links: links,
	}
}

// generatePaginationLinks creates and populates the Links structure for a PagedResponse
func generatePaginationLinks(c *fiber.Ctx, page, pageSize, totalPages int) (*Links, error) {
	baseURL := c.BaseURL()
	fullPath := c.Path()

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = fullPath

	q := u.Query()
	q.Set("page_size", fmt.Sprintf("%d", pageSize))

	links := &Links{}

	// Self link
	q.Set("page", fmt.Sprintf("%d", page))
	u.RawQuery = q.Encode()
	links.Self = u.String()

	// First link
	q.Set("page", "1")
	u.RawQuery = q.Encode()
	links.First = u.String()

	// Last link
	q.Set("page", fmt.Sprintf("%d", totalPages))
	u.RawQuery = q.Encode()
	links.Last = u.String()

	// Next link
	if page < totalPages {
		q.Set("page", fmt.Sprintf("%d", page+1))
		u.RawQuery = q.Encode()
		links.Next = u.String()
	}

	// Prev link
	if page > 1 {
		q.Set("page", fmt.Sprintf("%d", page-1))
		u.RawQuery = q.Encode()
		links.Prev = u.String()
	}

	return links, nil
}
