package events

// AssertDone is an event, emitted when assertion is done. Error is optional
type AssertDone struct {
	Success          bool
	Operation        string
	Actual, Expected interface{}
	Qualifier        string
	Doc              string
}

// IfIsAssertDone checks, if provided interface is instance of AssertDone event
// and if so, passes in into callback func
func IfIsAssertDone(e *Emitted, f func(AssertDone)) {
	if e != nil && e.Data != nil && f != nil {
		if a, ok := e.Data.(AssertDone); ok {
			f(a)
		}
	}
}
