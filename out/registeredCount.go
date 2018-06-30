package out

import "github.com/mono83/oscar/events"

// RegisteredCount is special event listener container, used to provide amount of registered test entities
type RegisteredCount struct {
	Value int
}

// BuildListener builds and returns event listener
func (r *RegisteredCount) BuildListener() func(e *events.Emitted) {
	er := events.EventRouter{
		RegistrationIn: func(events.RegistrationBegin, *events.Emitted) {
			r.Value = r.Value + 1
		},
	}

	return er.OnEvent
}
