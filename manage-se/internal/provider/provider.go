package provider

import (
	"manage-se/internal/appctx"
	"manage-se/internal/provider/auth"
	"manage-se/internal/provider/dependencies"
	"manage-se/internal/provider/wallet"
	"net/http"
	"time"
)

type Provider struct {
	Auth   Auth
	Wallet Wallet
}

func NewProviders(cfg *appctx.Provider, options ...Option) *Provider {
	dep := defaultDependency()

	for _, opt := range options {
		opt(dep)
	}

	return &Provider{
		Auth:   auth.NewClient(&cfg.Auth, dep),
		Wallet: wallet.NewClient(&cfg.Wallet, dep),
	}
}

func defaultDependency() *dependencies.Dependency {
	client := http.DefaultClient
	client.Timeout = time.Duration(10) * time.Second

	return &dependencies.Dependency{
		HttpClient: http.DefaultClient,
	}
}
