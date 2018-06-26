package events

import "time"

// RemoteRequest is TestContext event used to track remote requests (i.e. HTTP requests)
type RemoteRequest struct {
	Type    string
	URI     string
	Elapsed time.Duration
	Success bool
}
