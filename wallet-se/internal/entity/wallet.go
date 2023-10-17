package entity

import (
	"time"
)

type Wallet struct {
	ID        string     `json:"id" db:"id"`
	Balance   float64    `json:"balance" db:"balance"`
	Status    string     `json:"status" db:"status"`
	OwnedBy   string     `json:"owned_by" db:"owned_by"`
	EnabledAt string     `json:"enabled_at" db:"enabled_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type WalletDetail struct {
	Wallet
	Transaction []Transaction `json:"transactions"`
}
