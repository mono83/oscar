package events

import "time"

// Emitted contains emitted, fired event with time and owner ID
type Emitted struct {
	OwnerID int
	Time    time.Time
	Data    interface{}
}
