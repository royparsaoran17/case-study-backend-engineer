// Package middleware
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"manage-se/internal/appctx"
	"manage-se/internal/consts"
	"manage-se/pkg/logger"
)

// ValidateContentType header
func ValidateContentType(r *http.Request, conf *appctx.Config) int {

	if ct := strings.ToLower(r.Header.Get(`Content-Type`)); ct != `application/json` {
		logger.Warn(fmt.Sprintf("[middleware] invalid content-type %s", ct))

		return consts.CodeBadRequest
	}

	return consts.CodeSuccess
}
