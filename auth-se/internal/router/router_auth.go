package router

import (
	"auth-se/internal/ucase/auth"
	"net/http"

	"auth-se/internal/handler"
	authsvc "auth-se/internal/service/auth"
)

func (rtr *router) mountAuth(authSvc authsvc.Auth) {
	rtr.router.HandleFunc("/internal/v1/verify", rtr.handle(
		handler.HttpRequest,
		auth.NewVerify(authSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/login", rtr.handle(
		handler.HttpRequest,
		auth.NewLogin(authSvc),
	)).Methods(http.MethodPost)

}
