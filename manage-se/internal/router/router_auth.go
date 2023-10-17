package router

import (
	"manage-se/internal/ucase/auth"
	"net/http"

	"manage-se/internal/handler"
	authsvc "manage-se/internal/service/auth"
)

func (rtr *router) mountAuth(authSvc authsvc.Auth) {
	rtr.router.HandleFunc("/external/v1/verify", rtr.handle(
		handler.HttpRequest,
		auth.NewVerify(authSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/external/v1/login", rtr.handle(
		handler.HttpRequest,
		auth.NewLogin(authSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/external/v1/register", rtr.handle(
		handler.HttpRequest,
		auth.NewRegister(authSvc),
	)).Methods(http.MethodPost)

}
