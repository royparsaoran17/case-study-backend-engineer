package presentations

import "time"

// RequestState is the initial state for each request.
type RequestState struct {
	ID        string
	CreatedAt time.Time
}
