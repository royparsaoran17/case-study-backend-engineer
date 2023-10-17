package entity

import (
	"time"
)

type Transaction struct {
	ID            string     `json:"id" db:"id"`
	Type          string     `json:"type" db:"type"`
	ReferenceID   string     `json:"reference_id" db:"reference_id"`
	Amount        float64    `json:"amount" db:"amount"`
	Status        string     `json:"status" db:"status"`
	TransactionAt time.Time  `json:"transaction_at" db:"transaction_at"`
	TransactionBy string     `json:"transaction_by" db:"transaction_by"`
	WalletID      string     `json:"wallet_id" db:"wallet_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}
