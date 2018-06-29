package oscar

// Skip is an error, emitted, when test case was skipped
type Skip struct {
	Failed, Skipped string
}

func (s Skip) Error() string {
	if s.Failed == SuiteSetUp {
		return `Test case "` + s.Skipped + `" is skipped, because test suite initializer failed`
	}
	return `Test case "` + s.Skipped + `" is skipped, because "` + s.Failed + `" fails"`
}

// IsSkip returns true if provided error is skip
func IsSkip(e error) bool {
	if e != nil {
		_, ok := e.(Skip)
		return ok
	}

	return false
}
