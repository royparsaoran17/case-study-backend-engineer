package wallet

import (
	"net/http"
	"wallet-se/internal/service/wallet"
	"wallet-se/pkg/logger"
	"wallet-se/pkg/tracer"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"wallet-se/internal/appctx"
	"wallet-se/internal/consts"
	"wallet-se/internal/ucase/contract"
)

type walletGetByID struct {
	service wallet.Wallet
}

func (c walletGetByID) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.wallet_get_by_id")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("walletGetByID")

	walletID := mux.Vars(data.Request)["wallet_id"]
	if _, err := uuid.Parse(walletID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	result, err := c.service.GetWalletByID(ctx, walletID)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrWalletNotFound:
			return *responder.
				WithCode(http.StatusNotFound).
				WithMessage(causer.Error())

		default:
			logger.Info(err)
			tracer.SpanError(ctx, err)
			return *responder.
				WithCode(http.StatusInternalServerError).
				WithMessage(http.StatusText(http.StatusInternalServerError))
		}

	}

	return *responder.
		WithData(result).
		WithCode(http.StatusOK).
		WithMessage("wallet fetched")
}

func NewWalletGetByID(service wallet.Wallet) contract.UseCase {
	return &walletGetByID{service: service}
}
