package impact

// List of known impact levels
const (
	None    Level = 0 // No impact on testing infrastructure, no remote requests
	Read    Level = 1 // Test case performs shared data reading, generating load on infrastructure
	Default Level = 2 // Default impact level
	Create  Level = 3 // Test case creates new entries for own purposes and may (or may not) delete them
	Modify  Level = 4 // Test case performs modifications on shared area
)
