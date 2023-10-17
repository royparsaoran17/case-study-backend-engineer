package router

import (
	"net/http"

	"auth-se/internal/handler"
	rolesvc "auth-se/internal/service/role"
	"auth-se/internal/ucase/role"
)

func (rtr *router) mountRoles(roleSvc rolesvc.Role) {
	rtr.router.HandleFunc("/internal/v1/roles", rtr.handle(
		handler.HttpRequest,
		role.NewRoleGetAll(roleSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/roles", rtr.handle(
		handler.HttpRequest,
		role.NewRoleCreate(roleSvc),
	)).Methods(http.MethodPost)

	rtr.router.HandleFunc("/internal/v1/roles/{role_id}", rtr.handle(
		handler.HttpRequest,
		role.NewRoleGetByID(roleSvc),
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/internal/v1/roles/{role_id}", rtr.handle(
		handler.HttpRequest,
		role.NewRoleUpdate(roleSvc),
	)).Methods(http.MethodPut)

}
