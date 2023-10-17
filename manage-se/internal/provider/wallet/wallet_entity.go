package wallet

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

type Wallet struct {
	ID        string     `json:"id" db:"id"`
	Balance   float64    `json:"balance" db:"balance"`
	Status    string     `json:"status" db:"status"`
	OwnedBy   string     `json:"owned_by" db:"owned_by"`
	EnabledAt time.Time  `json:"enabled_at" db:"enabled_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type WalletDetail struct {
	Wallet
	Transaction []Transaction `json:"transactions"`
}

type WalletDeposit struct {
	CustomerID string  `json:"customer_id"`
	WalletID   string  `json:"wallet_id"`
	Balance    float64 `json:"balance"`
}

type WalletWithdraw struct {
	CustomerID string  `json:"customer_id"`
	WalletID   string  `json:"wallet_id"`
	Balance    float64 `json:"balance"`
}
