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

type walletGetByOwned struct {
	service wallet.Wallet
}

func (c walletGetByOwned) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.wallet_get_by_owned")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("walletGetByOwned")

	walletOwned := mux.Vars(data.Request)["owned_id"]
	if _, err := uuid.Parse(walletOwned); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	result, err := c.service.GetWalletByOwned(ctx, walletOwned)
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

func NewWalletGetByOwned(service wallet.Wallet) contract.UseCase {
	return &walletGetByOwned{service: service}
}
