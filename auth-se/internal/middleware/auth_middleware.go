package middleware

import (
	"auth-se/internal/appctx"
	"net/http"
)

func Authorize() MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, conf *appctx.Config) error {
		return nil
	}
}
