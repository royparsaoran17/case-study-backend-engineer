package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"wallet-se/internal/consts"
)

type WalletUpdate struct {
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
}

func (r *WalletUpdate) Validate() error {
	return validation.Errors{
		"balance": validation.Validate(&r.Balance),
		"status":  validation.Validate(&r.Status, validation.In(consts.StatusWalletEnabled.String(), consts.StatusWalletDisabled.String())),
	}.Filter()
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
	CustomerID string  `json:"customer_id"`
	WalletID   string  `json:"wallet_id"`
	Balance    float64 `json:"balance"`
}

func (r *WalletDeposit) Validate() error {
	return validation.Errors{
		"customer_id": validation.Validate(&r.CustomerID, validation.Required, is.UUID),
		"wallet_id":   validation.Validate(&r.WalletID, validation.Required, is.UUID),
		"balance":     validation.Validate(&r.Balance, validation.Required),
	}.Filter()
}

type WalletWithdraw struct {
	CustomerID string  `json:"customer_id"`
	WalletID   string  `json:"wallet_id"`
	Balance    float64 `json:"balance"`
}

func (r *WalletWithdraw) Validate() error {
	return validation.Errors{
		"customer_id": validation.Validate(&r.CustomerID, validation.Required, is.UUID),
		"wallet_id":   validation.Validate(&r.WalletID, validation.Required, is.UUID),
		"balance":     validation.Validate(&r.Balance, validation.Required),
	}.Filter()
}
