package role

import (
	"auth-se/internal/common"
	"auth-se/pkg/tracer"
	"net/http"

	"auth-se/internal/appctx"
	"auth-se/internal/consts"
	"auth-se/internal/service/role"
	"auth-se/internal/ucase/contract"
	"github.com/pkg/errors"
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

		switch causer := errors.Cause(err); causer {
		case common.ErrInvalidMetadata:
			return *responder.
				WithCode(http.StatusBadRequest).
				WithMessage(err.Error())

		default:
			switch causer.(type) {
			case consts.Error:
				return *responder.
					WithCode(http.StatusBadRequest).
					WithMessage(causer.Error())

			default:
				tracer.SpanError(ctx, err)
				return *responder.
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError))
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
