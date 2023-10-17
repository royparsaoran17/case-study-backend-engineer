package transaction

import (
	"net/http"
	"wallet-se/internal/common"
	"wallet-se/internal/service/transaction"
	"wallet-se/pkg/logger"
	"wallet-se/pkg/tracer"

	"github.com/pkg/errors"
	"wallet-se/internal/appctx"
	"wallet-se/internal/consts"
	"wallet-se/internal/ucase/contract"
)

type transactionGetAll struct {
	service transaction.Transaction
}

func (r transactionGetAll) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.transaction_get_all")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("transactionGetAll")

	transactions, err := r.service.GetAllTransaction(ctx)
	if err != nil {

		switch causer := errors.Cause(err); causer {
		case common.ErrInvalidMetadata:
			return *responder.
				WithCode(http.StatusBadRequest).
				WithMessage(err.Error())

		default:
			switch causer.(type) {
			case consts.Error:
				return *responder.
					WithCode(http.StatusBadRequest).
					WithMessage(causer.Error())

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
		WithMessage("get all transactions success").
		WithData(transactions)
}

func NewTransactionGetAll(service transaction.Transaction) contract.UseCase {
	return &transactionGetAll{service: service}
}
