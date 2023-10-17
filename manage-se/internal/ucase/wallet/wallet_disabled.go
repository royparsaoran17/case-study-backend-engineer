package wallet

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"manage-se/internal/consts"
	"manage-se/internal/provider/providererrors"
	"manage-se/internal/service/wallet"
	"manage-se/pkg/logger"
	"net/http"

	"manage-se/pkg/tracer"

	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/ucase/contract"
)

type disabled struct {
	service wallet.Wallet
}

func (c disabled) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.disabled")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("disabled")

	wallet, err := c.service.DisableWallet(ctx)
	if err != nil {
		errCause := errors.Cause(err)
		switch errCause {
		case consts.ErrUserUnauthorized:
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
				logger.Info(err)
				return *responder.WithContext(ctx).WithCode(http.StatusInternalServerError).WithMessage(http.StatusText(http.StatusInternalServerError))
			}
		}
	}

	return *responder.
		WithCode(http.StatusOK).
		WithData(wallet).
		WithMessage("disabled")
}

func NewDisabled(service wallet.Wallet) contract.UseCase {
	return &disabled{service: service}
}
