package events

// Start is an event, emitted, when some test holder (case, suite, etc.) is started
type Start struct {
	ID   int
	Type string
	Name string
}

// Finish is an event, emitted, when some test holder (case, suite, etc.) is finished
type Finish struct {
	ID    int
	Type  string
	Name  string
	Error error
}
