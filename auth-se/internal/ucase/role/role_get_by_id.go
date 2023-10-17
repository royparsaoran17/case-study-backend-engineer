package role

import (
	"net/http"

	"auth-se/pkg/tracer"

	"auth-se/internal/appctx"
	"auth-se/internal/consts"
	"auth-se/internal/service/role"
	"auth-se/internal/ucase/contract"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type roleGetByID struct {
	service role.Role
}

func (r roleGetByID) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.role_get_by_id")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("roleGetByID")

	roleID := mux.Vars(data.Request)["role_id"]
	if _, err := uuid.Parse(roleID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	result, err := r.service.GetRoleByID(ctx, roleID)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrRoleNotFound:
			return *responder.
				WithCode(http.StatusNotFound).
				WithMessage(causer.Error())

		default:
			tracer.SpanError(ctx, err)
			return *responder.
				WithCode(http.StatusInternalServerError).
				WithMessage(http.StatusText(http.StatusInternalServerError))
		}

	}

	return *responder.
		WithData(result).
		WithCode(http.StatusOK).
		WithMessage("role fetched")
}

func NewRoleGetByID(service role.Role) contract.UseCase {
	return &roleGetByID{service: service}
}
