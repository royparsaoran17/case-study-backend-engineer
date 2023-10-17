package role

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/tracer"
	"net/http"

	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/consts"
	"manage-se/internal/service/role"
	"manage-se/internal/ucase/contract"
)

type roleGetAll struct {
	service role.Role
}

func (r roleGetAll) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.role_get_all")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("roleGetAll")

	roles, err := r.service.GetAllRole(ctx)
	if err != nil {
		errCause := errors.Cause(err)
		switch errCause {
		default:
			switch causer := errCause.(type) {
			case consts.Error:
				return *responder.WithContext(ctx).WithCode(http.StatusBadRequest).WithMessage(errCause.Error())

			case providererrors.Error:
				return *responder.WithContext(ctx).WithCode(causer.Code).WithError(causer.Errors).WithMessage(causer.Message)

			case validation.Errors:
				return *responder.
					WithContext(ctx).
					WithCode(http.StatusUnprocessableEntity).
					WithMessage("Validation Error(s)").
					WithError(errCause)

			default:
				return *responder.WithContext(ctx).WithCode(http.StatusInternalServerError).WithMessage(http.StatusText(http.StatusInternalServerError))
			}
		}
	}

	return *responder.
		WithCode(http.StatusOK).
		WithMessage("get all roles success").
		WithData(roles)
}

func NewRoleGetAll(service role.Role) contract.UseCase {
	return &roleGetAll{service: service}
}
