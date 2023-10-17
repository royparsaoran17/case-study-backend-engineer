package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
	"wallet-se/internal/consts"
)

type TransactionCreate struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	TransactionAt time.Time `json:"transaction_at"`
	TransactionBy string    `json:"transaction_by"`
	WalletID      string    `json:"wallet_id"`
}

func (r *TransactionCreate) Validate() error {
	return validation.Errors{
		"type":           validation.Validate(&r.Type, validation.Required, validation.In(consts.TypeDebit.String(), consts.TypeCredit.String())),
		"amount":         validation.Validate(&r.Amount, validation.Required),
		"status":         validation.Validate(&r.Status, validation.Required, validation.In(consts.StatusTransactionSuccess.String(), consts.StatusTransactionFailed.String())),
		"transaction_by": validation.Validate(&r.TransactionBy, validation.Required, is.UUID),
		"transaction_at": validation.Validate(&r.TransactionAt, validation.Required),
		"wallet_id":      validation.Validate(&r.WalletID, validation.Required, is.UUID),
	}.Filter()
}
