package customer

import (
	"auth-se/internal/consts"
	"auth-se/internal/service/customer"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"auth-se/pkg/tracer"

	"auth-se/internal/appctx"
	"auth-se/internal/presentations"
	"auth-se/internal/ucase/contract"
	"github.com/pkg/errors"
)

type customerCreate struct {
	service customer.Customer
}

func (c customerCreate) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.customer_create")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("customerCreate")

	var input presentations.CustomerCreate
	if err := data.Cast(&input); err != nil {
		return *responder.WithCode(http.StatusBadRequest).
			WithError(err.Error()).
			WithMessage(http.StatusText(http.StatusBadRequest))
	}

	customers, err := c.service.CreateCustomer(ctx, input)
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
				tracer.SpanError(ctx, err)
				return *responder.
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError))
			}
		}
	}

	return *responder.
		WithCode(http.StatusCreated).
		WithData(customers).
		WithMessage("customer created")
}

func NewCustomerCreate(service customer.Customer) contract.UseCase {
	return &customerCreate{service: service}
}
