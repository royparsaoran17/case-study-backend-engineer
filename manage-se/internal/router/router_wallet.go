package router

import (
	"manage-se/internal/consts"
	"manage-se/internal/middleware"
	authsvc "manage-se/internal/service/auth"
	"manage-se/internal/ucase/wallet"
	"net/http"

	"manage-se/internal/handler"
	walletsvc "manage-se/internal/service/wallet"
)

func (rtr *router) mountWallet(walletSvc walletsvc.Wallet, authSvc authsvc.Auth) {
	rtr.router.HandleFunc("/external/v1/wallet/init", rtr.handle(
		handler.HttpRequest,
		wallet.NewInitialize(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/external/v1/wallet/enabled", rtr.handle(
		handler.HttpRequest,
		wallet.NewEnabled(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/external/v1/wallet/disabled", rtr.handle(
		handler.HttpRequest,
		wallet.NewDisabled(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/external/v1/wallet/balance", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletGetBalance(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/external/v1/wallet/transaction", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletGetTransaction(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/external/v1/wallet/deposits", rtr.handle(
		handler.HttpRequest,
		wallet.NewDeposit(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/external/v1/wallet/withdraws", rtr.handle(
		handler.HttpRequest,
		wallet.NewWithdraw(walletSvc),
		middleware.Authorize(authSvc, consts.AllRoles),
	)).Methods(http.MethodPost)

}
