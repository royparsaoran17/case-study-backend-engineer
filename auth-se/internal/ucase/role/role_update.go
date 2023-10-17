package role

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"auth-se/internal/appctx"
	"auth-se/internal/consts"
	"auth-se/internal/presentations"
	"auth-se/internal/service/role"
	"auth-se/internal/ucase/contract"
	"auth-se/pkg/logger"
	"auth-se/pkg/tracer"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type roleUpdate struct {
	service role.Role
}

func (r roleUpdate) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.role_update")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("roleUpdate")
	var input presentations.RoleUpdate

	if err := data.Cast(&input); err != nil {
		logger.Warn(fmt.Sprintf("error cast to roleUpdate presentation %+v", err))
		tracer.SpanError(ctx, err)
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	roleID := mux.Vars(data.Request)["role_id"]
	if _, err := uuid.Parse(roleID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	err := r.service.UpdateRoleByID(ctx, roleID, input)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrRoleNotFound:
			return *responder.
				WithCode(http.StatusNotFound).
				WithMessage(causer.Error())

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
		WithCode(http.StatusOK).
		WithMessage("role updated")
}

func NewRoleUpdate(service role.Role) contract.UseCase {
	return &roleUpdate{service: service}
}
