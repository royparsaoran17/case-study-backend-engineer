package transaction

import (
	"net/http"
	"wallet-se/internal/service/transaction"
	"wallet-se/pkg/logger"

	"wallet-se/pkg/tracer"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"wallet-se/internal/appctx"
	"wallet-se/internal/consts"
	"wallet-se/internal/ucase/contract"
)

type transactionGetByID struct {
	service transaction.Transaction
}

func (r transactionGetByID) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.transaction_get_by_id")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("transactionGetByID")

	transactionID := mux.Vars(data.Request)["transaction_id"]
	if _, err := uuid.Parse(transactionID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	result, err := r.service.GetTransactionByID(ctx, transactionID)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrTransactionNotFound:
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
		WithMessage("transaction fetched")
}

func NewTransactionGetByID(service transaction.Transaction) contract.UseCase {
	return &transactionGetByID{service: service}
}
