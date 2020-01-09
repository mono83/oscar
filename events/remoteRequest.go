package events

import (
	"fmt"
	"time"
)

// RemoteRequest is TestContext event used to track remote requests (i.e. HTTP requests)
type RemoteRequest struct {
	Type    string
	URI     string
	Ray     string
	Elapsed time.Duration
	Success bool
}

// ElapsedString returns string representation of elapsed time in seconds
func (r RemoteRequest) ElapsedString() string {
	return fmt.Sprintf("%.3fs", r.Elapsed.Seconds())
}
