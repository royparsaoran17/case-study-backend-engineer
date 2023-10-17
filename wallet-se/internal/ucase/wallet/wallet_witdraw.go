package wallet

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"wallet-se/internal/consts"
	"wallet-se/internal/service/wallet"
	"wallet-se/pkg/logger"
	"wallet-se/pkg/tracer"

	"github.com/pkg/errors"
	"wallet-se/internal/appctx"
	"wallet-se/internal/presentations"
	"wallet-se/internal/ucase/contract"
)

type walletWithdraw struct {
	service wallet.Wallet
}

func (c walletWithdraw) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.wallet_withdraw")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("walletWithdraw")

	var input presentations.WalletWithdraw
	if err := data.Cast(&input); err != nil {
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	wallets, err := c.service.WalletWithdraw(ctx, input)
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
				logger.Info(err)
				tracer.SpanError(ctx, err)
				return *responder.
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError))
			}
		}
	}

	return *responder.
		WithCode(http.StatusOK).
		WithData(wallets).
		WithMessage("wallet withdrawd")
}

func NewWalletWithdraw(service wallet.Wallet) contract.UseCase {
	return &walletWithdraw{service: service}
}
