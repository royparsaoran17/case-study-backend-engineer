package auth

import (
	"context"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/auth"
)

type Auth interface {
	VerifyToken(ctx context.Context, input presentations.Verify) (*auth.CustomerDetail, error)
	Login(ctx context.Context, input presentations.Login) (*auth.CustomerDetailToken, error)
	Register(ctx context.Context, input presentations.Register) (*auth.CustomerDetail, error)
}
