package lua

import (
	"github.com/mono83/oscar"
	lua "github.com/yuin/gopher-lua"
)

type testSuite struct {
	id    int
	name  string
	state *lua.LState
	setup *testcase
	cases []*testcase
}

// ID returns test suite name and identifier
func (t *testSuite) ID() (int, string) {
	return t.id, t.name
}

// GetSetUp returns optional setup function, that will be invoked before any other test cases
func (t *testSuite) GetSetUp() oscar.Case {
	if t.setup == nil {
		return nil
	}
	return t.setup
}

// GetCases returns test cases
func (t *testSuite) GetCases() []oscar.Case {
	cs := make([]oscar.Case, len(t.cases))
	for i := range t.cases {
		cs[i] = t.cases[i]
	}

	return cs
}
