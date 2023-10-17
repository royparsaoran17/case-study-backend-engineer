// Package middleware
package middleware

import (
	"auth-se/internal/appctx"
	"net/http"
)

// MiddlewareFunc is contract for middleware and must implement this type for http if need middleware http request
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, conf *appctx.Config) error

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(w http.ResponseWriter, r *http.Request, conf *appctx.Config, mfs []MiddlewareFunc) error {
	for _, mf := range mfs {
		if err := mf(w, r, conf); err != nil {
			return err
		}
	}

	return nil
}
