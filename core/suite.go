package core

// Suite is collection of test cases
type Suite interface {
	// ID returns test suite name and identifier
	ID() (int, string)

	// GetCases returns slice of cases, registered within suite
	GetCases() []Case
}
