package oscar

// Case represents common test case
type Case interface {
	// ID returns test case name and identifier
	ID() (int, string)

	// GetDependsOn returns slice of identifiers, that must succeed before case will run
	GetDependsOn() []int

	// Performs test case assertions using provided test context
	Assert(*Context) error
}
