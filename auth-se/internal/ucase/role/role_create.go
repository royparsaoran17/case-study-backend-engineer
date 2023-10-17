package role

import (
	"auth-se/internal/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"auth-se/pkg/tracer"

	"auth-se/internal/appctx"
	"auth-se/internal/presentations"
	"auth-se/internal/service/role"
	"auth-se/internal/ucase/contract"
	"github.com/pkg/errors"
)

type roleCreate struct {
	service role.Role
}

func (r roleCreate) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.role_create")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("roleCreate")

	var input presentations.RoleCreate
	if err := data.Cast(&input); err != nil {
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	err := r.service.CreateRole(ctx, input)
	if err != nil {
		causer := errors.Cause(err)
		switch causer {
		default:
			switch cause := causer.(type) {
			case consts.Error:
				return *responder.
					WithCode(http.StatusBadRequest).
					WithMessage(cause.Error())

			case validation.Errors:
				return *responder.
					WithCode(http.StatusUnprocessableEntity).
					WithError(cause).
					WithMessage("Validation error(s)")

			default:
				tracer.SpanError(ctx, err)
				return *responder.
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError))
			}
		}
	}

	return *responder.
		WithCode(http.StatusCreated).
		WithMessage("role created")
}

func NewRoleCreate(service role.Role) contract.UseCase {
	return &roleCreate{service: service}
}
