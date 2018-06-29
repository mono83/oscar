package out

import "github.com/mono83/oscar/events"

// Dispatcher is event receiver, that dispatches it to all listeners in List
type Dispatcher struct {
	List []func(*events.Emitted)
}

// OnEvent is method to be attached to TestContext
func (d *Dispatcher) OnEvent(e *events.Emitted) {
	if len(d.List) > 0 {
		for _, r := range d.List {
			r(e)
		}
	}
}
