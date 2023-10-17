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

type Register struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	RoleID string `json:"role_id"`
}

func (r *Register) Validate() error {
	return validation.Errors{
		"name":    validation.Validate(&r.Name, validation.Required),
		"phone":   validation.Validate(&r.Phone, validation.Required, is.E164),
		"role_id": validation.Validate(&r.RoleID, validation.Required, is.UUID),
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

type ValidateImage struct {
	Url string `json:"url"`
}

func (r *ValidateImage) Validate() error {
	return validation.Errors{
		"url": validation.Validate(&r.Url, is.URL, validation.Required),
	}.Filter()
}
