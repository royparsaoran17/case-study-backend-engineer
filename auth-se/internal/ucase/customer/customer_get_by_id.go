package customer

import (
	"auth-se/internal/service/customer"
	"net/http"

	"auth-se/pkg/tracer"

	"auth-se/internal/appctx"
	"auth-se/internal/consts"
	"auth-se/internal/ucase/contract"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type customerGetByID struct {
	service customer.Customer
}

func (c customerGetByID) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.customer_get_by_id")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("customerGetByID")

	customerID := mux.Vars(data.Request)["customer_id"]
	if _, err := uuid.Parse(customerID); err != nil {
		return *responder.
			WithCode(http.StatusBadRequest).
			WithMessage(consts.ErrInvalidUUID.Error())
	}

	result, err := c.service.GetCustomerByID(ctx, customerID)
	if err != nil {
		switch causer := errors.Cause(err); causer {
		case consts.ErrCustomerNotFound:
			return *responder.
				WithCode(http.StatusNotFound).
				WithMessage(causer.Error())

		default:
			tracer.SpanError(ctx, err)
			return *responder.
				WithCode(http.StatusInternalServerError).
				WithMessage(http.StatusText(http.StatusInternalServerError))
		}

	}

	return *responder.
		WithData(result).
		WithCode(http.StatusOK).
		WithMessage("customer fetched")
}

func NewCustomerGetByID(service customer.Customer) contract.UseCase {
	return &customerGetByID{service: service}
}
