package oscar

import "time"

// TraceEvent is TestContext event for tracing purposes
type TraceEvent string

// StartEvent is TestContext event, emitted when something is started
type StartEvent struct {
	Owner interface{}
	Time  time.Time
}

// FinishEvent is TestContext event, emitted when something is done
type FinishEvent struct {
	Owner interface{}
	Time  time.Time
}

// RemoteRequestEvent is TestContext event used to track remote requests (i.e. HTTP requests)
type RemoteRequestEvent struct {
	Type    string
	Elapsed time.Duration
	Success bool
}

// SleepEvent is emitted on every sleep
type SleepEvent time.Duration

// AssertionSuccess is TestContext event, emitted on every success assertion
type AssertionSuccess string

// AssertionFailure is TestContext event, emitted on every failed assertion
type AssertionFailure error

// TestLogEvent is TestContext event for logging purposes
type TestLogEvent struct {
	Owner   *TestCase
	Message string
	Level   byte
}
