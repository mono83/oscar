package events

// Start is an event, emitted, when some test holder (case, suite, etc.) is started
type Start struct {
	Type string
	Name string
}

// Finish is an event, emitted, when some test holder (case, suite, etc.) is finished
type Finish struct {
	Type  string
	Name  string
	Error error
}
