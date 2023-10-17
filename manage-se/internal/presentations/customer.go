package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CustomerUpdate struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	RoleID string `json:"role_id"`
}

func (r *CustomerUpdate) Validate() error {
	return validation.Errors{
		"name":    validation.Validate(&r.Name, validation.Required),
		"phone":   validation.Validate(&r.Phone, validation.Required, is.E164),
		"role_id": validation.Validate(&r.RoleID, validation.Required, is.UUID),
	}.Filter()
}

type CustomerCreate struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

func (r *CustomerCreate) Validate() error {
	return validation.Errors{
		"name":     validation.Validate(&r.Name, validation.Required),
		"phone":    validation.Validate(&r.Phone, validation.Required, is.E164),
		"password": validation.Validate(&r.Password, validation.Required),
		"role_id":  validation.Validate(&r.RoleID, validation.Required, is.UUID),
	}.Filter()
}
