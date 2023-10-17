// Package contract
package contract

import (
	"auth-se/internal/appctx"
	"context"
)

// UseCase is a use case contract
type UseCase interface {
	Serve(data *appctx.Data) appctx.Response
}

// MessageProcessor is use case queue message processor contract
type MessageProcessor interface {
	Serve(ctx context.Context, data *appctx.ConsumerData) error
}
