package customer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"manage-se/internal/common"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/tracer"
	"net/http"

	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/consts"
	"manage-se/internal/service/customer"
	"manage-se/internal/ucase/contract"
)

type customerGetAll struct {
	service customer.Customer
}

func (r customerGetAll) Serve(data *appctx.Data) appctx.Response {
	ctx := tracer.SpanStart(data.Request.Context(), "usecase.customer_get_all")
	defer tracer.SpanFinish(ctx)

	responder := appctx.NewResponse().WithState("customerGetAll")

	metaData := common.MetadataFromURL(data.Request.URL.Query())

	metaDateRange, err := common.DateRangeFromURL(data.Request.URL.Query(), "created_at", "created_from", "created_until")
	if err != nil {
		return *responder.WithCode(http.StatusBadRequest).WithMessage(err.Error())
	}

	metaData.DateRange = metaDateRange

	customers, err := r.service.GetAllCustomer(ctx, &metaData)
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
		WithMessage("get all customers success").
		WithMeta(metaData).
		WithData(customers)
}

func NewCustomerGetAll(service customer.Customer) contract.UseCase {
	return &customerGetAll{service: service}
}
