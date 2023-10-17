package wallet

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"wallet-se/internal/appctx"
	"wallet-se/internal/consts"
	"wallet-se/internal/presentations"
	"wallet-se/internal/service/wallet"
	"wallet-se/internal/ucase/contract"
	"wallet-se/pkg/logger"
	"wallet-se/pkg/tracer"
)

type walletUpdate struct {
	service wallet.Wallet
}

func (c walletUpdate) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.wallet_update")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("walletUpdate")
	var input presentations.WalletUpdate

	if err := data.Cast(&input); err != nil {
		logger.Warn(fmt.Sprintf("error cast to walletUpdate presentation %+v", err))
		tracer.SpanError(ctx, err)
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	walletID := mux.Vars(data.Request)["wallet_id"]
	if _, err := uuid.Parse(walletID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	err := c.service.UpdateWalletByID(ctx, walletID, input)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrWalletNotFound:
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
		WithMessage("wallet updated")
}

func NewWalletUpdate(service wallet.Wallet) contract.UseCase {
	return &walletUpdate{service: service}
}
