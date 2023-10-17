package router

import (
	"net/http"
	"wallet-se/internal/handler"
	walletsvc "wallet-se/internal/service/wallet"
	"wallet-se/internal/ucase/wallet"
)

func (rtr *router) mountWallet(walletSvc walletsvc.Wallet) {
	rtr.router.HandleFunc("/internal/v1/wallet/withdraw", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletWithdraw(walletSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/wallet/deposit", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletDeposit(walletSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/wallets", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletGetAll(walletSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/wallets", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletCreate(walletSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/wallets/{wallet_id}", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletGetByID(walletSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/wallets/owned/{owned_id}", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletGetByOwned(walletSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/wallets/{wallet_id}", rtr.handle(
		handler.HttpRequest,
		wallet.NewWalletUpdate(walletSvc),
	)).Methods(http.MethodPut)

}
