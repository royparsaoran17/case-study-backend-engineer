package entity

import "time"

type WalletBalance struct {
	ID        string          `json:"id"`
	OwnedBy   string          `json:"owned_by"`
	Customer  CustomerContext `json:"customer"`
	Status    string          `json:"status"`
	EnabledAt time.Time       `json:"enabled_at"`
	Balance   float64         `json:"balance"`
}
