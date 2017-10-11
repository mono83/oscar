package oscar

import (
	"fmt"
	"regexp"
	"time"
)

// TestContext is nested structure, that holds test invocation context
type TestContext struct {
	Parent *TestContext

	Vars  map[string]string
	Error error

	CountAssertSuccess  int
	CountRemoteRequests int

	OnFinish func(*TestContext) error
	OnEvent  func(interface{})

	startedAt, finishedAt time.Time
}

// Get returns variable value from vars map
func (t *TestContext) Get(key string) string {
	if len(t.Vars) > 0 {
		if v, ok := t.Vars[key]; ok {
			return v
		}
	}
	if t.Parent != nil {
		return t.Parent.Get(key)
	}

	return ""
}

// Set assigns new variable value
func (t *TestContext) Set(key, value string) {
	t.Trace(`Setting "%s" := "%s"`, key, value)
	if len(t.Vars) == 0 {
		t.Vars = map[string]string{}
	}

	t.Vars[key] = value
}

var iregex = regexp.MustCompile(`\${([\w.-]+)}`)

// Interpolate replaces all placeholders in provided string using vars from test case or
// global runner
func (t *TestContext) Interpolate(value string) string {
	return iregex.ReplaceAllStringFunc(value, func(i string) string {
		m := iregex.FindStringSubmatch(i)
		return t.Get(m[1])
	})
}

// Elapsed returns elapsed time
func (t *TestContext) Elapsed() time.Duration {
	return t.finishedAt.Sub(t.startedAt)
}

// Emit publishes new event into nested test context
func (t *TestContext) Emit(event interface{}) {
	if s, ok := event.(StartEvent); ok {
		if t.startedAt.IsZero() {
			t.startedAt = s.Time
		}
	} else if s, ok := event.(FinishEvent); ok {
		t.finishedAt = s.Time
	} else if _, ok := event.(AssertionSuccess); ok {
		t.CountAssertSuccess++
	} else if a, ok := event.(AssertionFailure); ok {
		if t.Error == nil {
			t.Error = a
		}
	} else if _, ok := event.(RemoteRequestEvent); ok {
		t.CountRemoteRequests++
	}

	if t.OnEvent != nil {
		t.OnEvent(event)
	}
	if s, ok := event.(FinishEvent); ok && t.OnFinish != nil {
		if t, ok := s.Owner.(*TestCase); ok {
			t.OnFinish(t.TestContext)
		}
	}

	if t.Parent != nil {
		t.Parent.Emit(event)
	}
}

// Trace emits tracing event
func (t *TestContext) Trace(pattern string, args ...interface{}) {
	t.Emit(TraceEvent(fmt.Sprintf(pattern, args...)))
}

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
