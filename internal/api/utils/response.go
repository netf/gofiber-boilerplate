package utils

import (
	"math"
)

// Response represents a generic response structure
type Response[T any] struct {
	Data interface{} `json:"data"`
	Meta *MetaData   `json:"meta,omitempty"`
}

// MetaData contains pagination metadata
type MetaData struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// CreateResponse generates a Response for either a single item or a paginated list
func CreateResponse[T any](data interface{}, opts ...int) Response[T] {
	response := Response[T]{
		Data: data,
	}

	if _, ok := data.([]T); ok {
		// It's a slice, so create a paged response
		if len(opts) != 3 {
			return response
		}
		page, pageSize, total := opts[0], opts[1], int64(opts[2])
		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		response.Meta = &MetaData{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		}
	}

	return response
}

// ErrorResponse represents a standardized error response structure
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// CreateErrorResponse generates a standardized error response
func CreateErrorResponse(message string, statusCode int) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  statusCode,
	}
}
