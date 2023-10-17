package router

import (
	"manage-se/internal/middleware"
	authsvc "manage-se/internal/service/auth"
	"manage-se/internal/ucase/customer"
	"net/http"

	"manage-se/internal/handler"
	customersvc "manage-se/internal/service/customer"
)

func (rtr *router) mountCustomer(customerSvc customersvc.Customer, authSvc authsvc.Auth) {
	rtr.router.HandleFunc("/external/v1/customers", rtr.handle(
		handler.HttpRequest,
		customer.NewCustomerGetAll(customerSvc),
		middleware.Authorize(authSvc, "Super Admin"),
	)).Methods(http.MethodGet)

}
