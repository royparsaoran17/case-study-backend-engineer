package router

import (
	"net/http"
	"wallet-se/internal/ucase/transaction"

	"wallet-se/internal/handler"
	transactionsvc "wallet-se/internal/service/transaction"
)

func (rtr *router) mountTransactions(transactionSvc transactionsvc.Transaction) {
	rtr.router.HandleFunc("/internal/v1/transactions", rtr.handle(
		handler.HttpRequest,
		transaction.NewTransactionGetAll(transactionSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/transactions", rtr.handle(
		handler.HttpRequest,
		transaction.NewTransactionCreate(transactionSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/transactions/{transaction_id}", rtr.handle(
		handler.HttpRequest,
		transaction.NewTransactionGetByID(transactionSvc),
	)).Methods(http.MethodGet)

}
