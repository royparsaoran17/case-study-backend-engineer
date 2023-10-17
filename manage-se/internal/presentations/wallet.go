package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type WalletUpdate struct {
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
}

type WalletCreate struct {
	ID      string `json:"id"`
	OwnedBy string `json:"owned_by"`
}

func (r *WalletCreate) Validate() error {
	return validation.Errors{
		"owned_by": validation.Validate(&r.OwnedBy, validation.Required, is.UUID),
	}.Filter()
}

type WalletDeposit struct {
	Balance float64 `json:"balance"`
}

func (r *WalletDeposit) Validate() error {
	return validation.Errors{
		"balance": validation.Validate(&r.Balance, validation.Required),
	}.Filter()
}

type WalletWithdraw struct {
	Balance float64 `json:"balance"`
}

func (r *WalletWithdraw) Validate() error {
	return validation.Errors{
		"balance": validation.Validate(&r.Balance, validation.Required),
	}.Filter()
}
