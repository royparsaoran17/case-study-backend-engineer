package auth

import (
	"auth-se/internal/consts"
	"auth-se/internal/service/auth"
	"auth-se/pkg/jwt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"auth-se/pkg/tracer"

	"auth-se/internal/appctx"
	"auth-se/internal/presentations"
	"auth-se/internal/ucase/contract"
	"github.com/pkg/errors"
)

type verify struct {
	service auth.Auth
}

func (c verify) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.verify")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("verify")

	var input presentations.Verify
	if err := data.Cast(&input); err != nil {
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	customers, err := c.service.Verify(input)
	if err != nil {
		causer := errors.Cause(err)
		switch causer {
		case jwt.ErrTokenExpired:
			return *responder.
				WithCode(http.StatusUnauthorized).
				WithMessage(causer.Error())

		default:
			switch cause := causer.(type) {
			case consts.Error, jwt.Error:
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
		WithData(customers).
		WithMessage("customer created")
}

func NewVerify(service auth.Auth) contract.UseCase {
	return &verify{service: service}
}
