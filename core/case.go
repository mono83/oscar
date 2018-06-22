package core

// Case represents common test case
type Case interface {
	// ID returns test case name and identifier
	ID() (int, string)

	// Performs test case assertions using provided test context
	Assert(*Context) error
}
