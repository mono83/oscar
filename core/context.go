package core

import (
	"fmt"
	"github.com/mono83/oscar/core/events"
	"regexp"
	"sync"
	"time"
)

// NewContext builds and returns context to be used in tests
func NewContext() *Context {
	c := &Context{
		values: make(map[string]string),
		events: make(chan interface{}),
	}

	go c.listenEvents()

	return c
}

// Context is test invocation context
type Context struct {
	parent *Context
	m      sync.Mutex
	values map[string]string

	wg      sync.WaitGroup
	events  chan interface{}
	OnEvent func(interface{})
}

// Fork builds and returns new child test context
func (c *Context) Fork() *Context {
	c2 := &Context{
		parent: c,
		values: make(map[string]string),
		events: make(chan interface{}),
	}

	go c2.listenEvents()

	return c2
}

func (c *Context) listenEvents() {
	for e := range c.events {
		c.wg.Done()
		if c.OnEvent != nil {
			c.OnEvent(e)
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

// Emit emits event to registered consumers
func (c *Context) Emit(t interface{}) {
	if t != nil {
		c.wg.Add(1)
		c.events <- t

		if c.parent != nil {
			c.parent.Emit(t)
		}
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
	c.Emit(events.LogEvent{Level: 0, Pattern: fmt.Sprintf(pattern, a...)})
}

// Debug sends DEBUG event with variables interpolation
func (c *Context) Debug(message string) {
	c.Emit(events.LogEvent{Level: 1, Pattern: c.Interpolate(message)})
}

// Debugf sends DEBUG event without interpolation but with sprintf formatting
func (c *Context) Debugf(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: 1, Pattern: fmt.Sprintf(pattern, a...)})
}

// Info sends INFO event with variables interpolation
func (c *Context) Info(message string) {
	c.Emit(events.LogEvent{Level: 2, Pattern: c.Interpolate(message)})
}

// Infof sends INFO event without interpolation but with sprintf formatting
func (c *Context) Infof(pattern string, a ...interface{}) {
	c.Emit(events.LogEvent{Level: 2, Pattern: fmt.Sprintf(pattern, a...)})
}

// Get returns variable value
func (c *Context) Get(key string) string {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.values) > 0 {
		if value, ok := c.values[key]; ok {
			return value
		}
	}

	if c.parent != nil {
		return c.parent.Get(key)
	}

	return ""
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

// SetInitial writes variable values and do not emit event about this
func (c *Context) SetInitial(m map[string]string) {
	c.m.Lock()
	defer c.m.Unlock()

	for k, v := range m {
		c.values[k] = v
	}
}

// AssertFinished sends event about finished assertion
func (c *Context) AssertFinished(err error) {
	c.Emit(events.AssertDone{Error: err})
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
