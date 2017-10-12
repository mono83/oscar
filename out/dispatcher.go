package out

// Dispatcher is event receiver, that dispatches it to all listeners in List
type Dispatcher struct {
	List []func(interface{})
}

// OnEvent is method to be attached to TestContext
func (d *Dispatcher) OnEvent(e interface{}) {
	if len(d.List) > 0 {
		for _, r := range d.List {
			r(e)
		}
	}
}
