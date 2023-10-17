package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"manage-se/internal/consts"
	"manage-se/internal/provider/providererrors"
	"manage-se/internal/service/auth"
	"manage-se/pkg/jwt"
	"net/http"

	"manage-se/pkg/tracer"

	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/presentations"
	"manage-se/internal/ucase/contract"
)

type register struct {
	service auth.Auth
}

func (c register) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.register")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("register")

	var input presentations.Register
	if err := data.Cast(&input); err != nil {
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	customers, err := c.service.Register(ctx, input)
	if err != nil {
		errCause := errors.Cause(err)
		switch errCause {
		case jwt.ErrTokenExpired:
			return *responder.
				WithCode(http.StatusUnauthorized).
				WithMessage(errCause.Error())

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
		WithData(customers).
		WithMessage("customer registered")
}

func NewRegister(service auth.Auth) contract.UseCase {
	return &register{service: service}
}