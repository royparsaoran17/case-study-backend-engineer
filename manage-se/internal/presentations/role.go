package presentations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RoleUpdate struct {
	Name string `json:"name"`
}

func (r *RoleUpdate) Validate() error {
	return validation.Errors{
		"name": validation.Validate(&r.Name, validation.Required),
	}.Filter()
}

type RoleCreate struct {
	Name string `json:"name"`
}

func (r *RoleCreate) Validate() error {
	return validation.Errors{
		"name": validation.Validate(&r.Name, validation.Required),
	}.Filter()
}
