// Package middleware
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"wallet-se/internal/appctx"
	"wallet-se/internal/consts"
	"wallet-se/pkg/logger"
)

// ValidateContentType header
func ValidateContentType(r *http.Request, conf *appctx.Config) int {

	if ct := strings.ToLower(r.Header.Get(`Content-Type`)); ct != `application/json` {
		logger.Warn(fmt.Sprintf("[middleware] invalid content-type %s", ct))

		return consts.CodeBadRequest
	}

	return consts.CodeSuccess
}
