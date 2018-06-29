package oscar

import "github.com/mono83/oscar/impact"

// Case represents common test case
type Case interface {
	// ID returns test case name and identifier
	ID() (int, string)

	// GetImpact returns impact level, induced by test case on remote infrastructure
	GetImpact() impact.Level

	// Performs test case assertions using provided test context
	Assert(*Context) error
}
