package provider

import (
	"manage-se/internal/provider/dependencies"
)

type Option func(dependency *dependencies.Dependency)

func WithHttpClient(client dependencies.HttpClient) Option {
	return func(dependency *dependencies.Dependency) {
		dependency.HttpClient = client
	}
}
