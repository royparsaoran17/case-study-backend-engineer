package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Login struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r *Login) Validate() error {
	return validation.Errors{
		"phone":    validation.Validate(&r.Phone, validation.Required, is.E164),
		"password": validation.Validate(&r.Password, validation.Required),
	}.Filter()
}

type Verify struct {
	Token string `json:"token"`
}

func (r *Verify) Validate() error {
	return validation.Errors{
		"token": validation.Validate(&r.Token, validation.Required),
	}.Filter()
}
