package router

import (
	"manage-se/internal/ucase/role"
	"net/http"

	"manage-se/internal/handler"
	rolesvc "manage-se/internal/service/role"
)

func (rtr *router) mountRole(roleSvc rolesvc.Role) {
	rtr.router.HandleFunc("/external/v1/roles", rtr.handle(
		handler.HttpRequest,
		role.NewRoleGetAll(roleSvc),
	)).Methods(http.MethodGet)

}
