package oscar

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/mono83/oscar/events"
)

// NewContext builds and returns context to be used in tests
func NewContext() *Context {
	c := &Context{
		values: make(map[string]string),
		events: make(chan *events.Emitted),
	}

	go c.listenEvents()

	return c
}

// Context is test invocation context
type Context struct {
	parent  *Context
	m       sync.Mutex
	values  map[string]string
	exports map[string]string

	ownerID int
	wg      sync.WaitGroup
	events  chan *events.Emitted
	onEvent func(*events.Emitted)
}

// Register registers new event listener
func (c *Context) Register(f func(*events.Emitted)) {
	if f == nil {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	if c.onEvent == nil {
		c.onEvent = f
	} else {
		prev := c.onEvent
		c.onEvent = func(event *events.Emitted) {
			f(event)
			prev(event)
		}
	}
}

// Fork builds and returns new child test context
func (c *Context) Fork(id int) *Context {
	c2 := &Context{
		parent:  c,
		ownerID: id,
		values:  make(map[string]string),
		events:  make(chan *events.Emitted),
	}

	go c2.listenEvents()

	return c2
}

func (c *Context) listenEvents() {
	for e := range c.events {
		c.wg.Done()
		if c.onEvent != nil {
			c.onEvent(e)
		}
	}
}

// Sleep freezes goroutine for required amount of time
func (c *Context) Sleep(duration time.Duration) {
	if duration.Nanoseconds() > 0 {
		c.Tracef("Entering sleep for %s", duration.String())
		c.Emit(events.Sleep(duration))
		time.Sleep(duration)
	}
}

func (c *Context) realEmit(e *events.Emitted) {
	c.wg.Add(1)
	c.events <- e

	if c.parent != nil {
		c.parent.realEmit(e)
	}
}

// Emit emits event to registered consumers
func (c *Context) Emit(t interface{}) {
	if t != nil {
		c.realEmit(&events.Emitted{OwnerID: c.ownerID, Time: time.Now(), Data: t})
	}
}

// Wait locks goroutine and waits for all events to be delivered
func (c *Context) Wait() {
	if c.parent != nil {
		c.parent.Wait()
	}
	c.wg.Wait()
}

// Tracef sends TRACE event without interpolation but with sprintf formatting
func (c *Context) Tracef(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: events.LogLevelTrace, Pattern: fmt.Sprintf(pattern, a...)})
}

// Debug sends DEBUG event with variables interpolation
func (c *Context) Debug(message string) {
	c.Emit(events.LogEvent{Level: events.LogLevelDebug, Pattern: c.Interpolate(message)})
}

// Debugf sends DEBUG event without interpolation but with sprintf formatting
func (c *Context) Debugf(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: events.LogLevelDebug, Pattern: fmt.Sprintf(pattern, a...)})
}

// Info sends INFO event with variables interpolation
func (c *Context) Info(message string) {
	c.Emit(events.LogEvent{Level: events.LogLevelInfo, Pattern: c.Interpolate(message)})
}

// Infof sends INFO event without interpolation but with sprintf formatting
func (c *Context) Infof(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: events.LogLevelInfo, Pattern: fmt.Sprintf(pattern, a...)})
}

// Errorf sends ERROR event without interpolation but with sprintf formatting
func (c *Context) Errorf(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: events.LogLevelError, Pattern: fmt.Sprintf(pattern, a...)})
}

// Fail emits Failure event
func (c *Context) Fail(message string) {
	c.Emit(events.Failure(message))
}

// Get returns variable value
func (c *Context) Get(key string) string {
	c.m.Lock()
	defer c.m.Unlock()

	// Reading own values
	if len(c.values) > 0 {
		if value, ok := c.values[key]; ok {
			return value
		}
	}

	// Reading values from parent
	if c.parent != nil {
		return c.parent.Get(key)
	}

	// Reading values from export variables
	if len(c.exports) > 0 {
		if value, ok := c.exports[key]; ok {
			return value
		}
	}

	return ""
}

// GetExport returns map of export variables
func (c *Context) GetExport() map[string]string {
	return c.exports
}

// Set places new variable value
func (c *Context) Set(key, value string) {
	c.m.Lock()
	defer c.m.Unlock()

	var prev *string
	if v, ok := c.values[key]; ok {
		prev = &v
	}

	c.values[key] = value
	c.Emit(events.SetVar{Key: key, Value: value, Previous: prev})
}

// SetExport adds variable on top scope level
func (c *Context) SetExport(key, value string) {
	if c.parent != nil {
		c.parent.SetExport(key, value)
	} else {
		c.m.Lock()
		defer c.m.Unlock()

		if len(c.exports) == 0 {
			c.exports = map[string]string{}
		}
		c.exports[key] = value
	}
}

// Import writes variable values and do not emit event about this
// This method may be used to fill initial data or to copy data from SetUp func
func (c *Context) Import(m map[string]string) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(m) > 0 {
		for k, v := range m {
			c.values[k] = v
		}
	}
}

// AssertFinished sends event about finished assertion
func (c *Context) AssertFinished(success bool, operation, qualifier, doc string, actual, expected interface{}) {
	c.Emit(events.AssertDone{
		Success:   success,
		Operation: operation,
		Actual:    actual,
		Expected:  expected,
		Qualifier: qualifier,
		Doc:       doc,
	})
}

var iregex = regexp.MustCompile(`\${([\w.-]+)}`)

// Interpolate replaces all placeholders in provided string using vars from test case or
// global runner
func (c *Context) Interpolate(value string) string {
	return iregex.ReplaceAllStringFunc(value, func(i string) string {
		m := iregex.FindStringSubmatch(i)
		return c.Get(m[1])
	})
}
