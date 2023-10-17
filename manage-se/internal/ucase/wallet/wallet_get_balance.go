package wallet

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/tracer"
	"net/http"

	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/consts"
	"manage-se/internal/service/wallet"
	"manage-se/internal/ucase/contract"
)

type walletGetBalance struct {
	service wallet.Wallet
}

func (r walletGetBalance) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.wallet_get_balance")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("walletGetBalance")

	wallets, err := r.service.ViewWalletBalance(ctx)
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
		WithMessage("get wallets balance success").
		WithData(wallets)
}

func NewWalletGetBalance(service wallet.Wallet) contract.UseCase {
	return &walletGetBalance{service: service}
}
