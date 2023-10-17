package router

import (
	customersvc "auth-se/internal/service/customer"
	"auth-se/internal/ucase/customer"
	"net/http"

	"auth-se/internal/handler"
)

func (rtr *router) mountCustomer(customerSvc customersvc.Customer) {
	rtr.router.HandleFunc("/internal/v1/customers", rtr.handle(
		handler.HttpRequest,
		customer.NewCustomerGetAll(customerSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/customers", rtr.handle(
		handler.HttpRequest,
		customer.NewCustomerCreate(customerSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/customers/{customer_id}", rtr.handle(
		handler.HttpRequest,
		customer.NewCustomerGetByID(customerSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/customers/{customer_id}", rtr.handle(
		handler.HttpRequest,
		customer.NewCustomerUpdate(customerSvc),
	)).Methods(http.MethodPut)

}
