package events

import (
	"fmt"
	"strings"
	"time"
)

// Emitted contains emitted, fired event with time and owner ID
type Emitted struct {
	OwnerID int
	Time    time.Time
	Data    interface{}
}

// TypeString return type string
func (e Emitted) TypeString() string {
	ts := fmt.Sprintf("%T", e.Data)
	if strings.HasPrefix(ts, "events.") {
		ts = ts[7:]
	}
	return ts
}

// TimeString returns string representation of log event time
func (e Emitted) TimeString() string {
	return e.Time.Format("15:04:05.000")
}
