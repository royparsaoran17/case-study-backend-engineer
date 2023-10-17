package auth

import "time"

type Metadata struct {
	Pagination
	Filtering
	*DateRange `json:"date_range,omitempty"`
}

type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Total   int `json:"total"`
}

// Filtering represents a filterable fields.
type Filtering struct {
	OrderBy   string `json:"order_by"`
	OrderType string `json:"order_type"`
	Search    string `json:"search,omitempty"`
	SearchBy  string `json:"search_by,omitempty"`
}

type DateRange struct {
	Field string    `json:"field"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
