package core

// SuiteSetUp contains internal name for suite initializer func
const SuiteSetUp = "__INIT__"

// Suite is collection of test cases
type Suite interface {
	// ID returns test suite name and identifier
	ID() (int, string)

	// GetSetUp returns optional setup function, that will be invoked before any other test cases
	GetSetUp() Case

	// GetCases returns slice of cases, registered within suite
	GetCases() []Case
}
