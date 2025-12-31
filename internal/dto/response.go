package dto

// APIResponse is the standard response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ListResponse is a paginated list response
type ListResponse[T any] struct {
	Items       []T   `json:"items"`
	TotalItems  int64 `json:"totalItems"`
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int64 `json:"totalPages"`
}
