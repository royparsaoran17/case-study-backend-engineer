package repooption

import (
	"time"
)

type QueryOptions struct {
	IncludeInactive  bool
	CreatedAt        *time.Time
	NotIncludeSample []string
	IncludeLimit     *int
}

type QueryOption func(*QueryOptions)

func WithInactive(option bool) QueryOption {
	return func(options *QueryOptions) {
		options.IncludeInactive = option
	}
}

func WithLimit(option *int) QueryOption {
	return func(options *QueryOptions) {
		options.IncludeLimit = option
	}
}

func WithSampleData(option []string) QueryOption {
	return func(options *QueryOptions) {
		options.NotIncludeSample = option
	}
}

func WithCreatedAt(t *time.Time) QueryOption {
	return func(options *QueryOptions) {
		options.CreatedAt = t
	}
}
