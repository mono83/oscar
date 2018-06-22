package lua

import (
	"errors"
	"github.com/mono83/oscar/core"
	"github.com/mono83/oscar/util/jsonencoder"
	"github.com/mono83/oscar/util/rsa"
	"github.com/yuin/gopher-lua"
)

// TestCaseMeta contains metatable name for userdata structure TestCase in lua
const TestCaseMeta = "TestCaseType"

// SuiteFromFiles builds suite using Lua sources file
func SuiteFromFiles(files ...string) (core.Suite, error) {
	if len(files) == 0 {
		return nil, errors.New("empty files list to load")
	}

	// Building Lua state
	L := lua.NewState()

	// Building test suite
	s := &fileTestSuite{
		id:    id(),
		name:  files[len(files)-1],
		state: L,
	}

	s.InjectModule(L)

	// Reading files sequentially
	for _, file := range files {
		if err := L.DoFile(file); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type fileTestSuite struct {
	id    int
	name  string
	state *lua.LState
	cases []*testcase
}

// ID returns test suite name and identifier
func (f *fileTestSuite) ID() (int, string) {
	return f.id, f.name
}

// GetCases returns test cases
func (f *fileTestSuite) GetCases() []core.Case {
	cs := make([]core.Case, len(f.cases))
	for i := range f.cases {
		cs[i] = f.cases[i]
	}

	return cs
}

// InjectModule injects TestSuite module (named "oscar") into lua engine
func (f *fileTestSuite) InjectModule(L *lua.LState) {
	L.PreloadModule("oscar", func(L *lua.LState) int {
		// Registering test case type
		mt := L.NewTypeMetatable(TestCaseMeta)
		// methods
		L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"assert":          lAssertEquals,
			"assertEq":        lAssertEquals,
			"assertEquals":    lAssertEquals,
			"assertJSONPath":  lAssertJSONPath,
			"assertJSONXPath": lAssertJSONPath,
			"get":             lGetVar,
			"httpGet":         lHTTPGet,
			"httpPost":        lHTTPPost,
			"jsonPath":        lJSONPathExtract,
			"jsonXPath":       lJSONPathExtract,
			"set":             lSetVar,

			"unix":  lUnix,
			"sleep": lSleep,

			"stringToBase64":    lStringToBase64,
			"packInt64ToBase64": lPackSliceInt64ToBase64,

			"log":   lDebug,
			"debug": lDebug,
			"info":  lInfo,
		}))

		// Making lambdas
		reg := func(L *lua.LState) int {
			name := L.CheckString(1)
			clb := L.CheckFunction(2)

			f.cases = append(
				f.cases,
				&testcase{
					id:       id(),
					name:     name,
					function: clb,
					state:    f.state,
				},
			)

			return 0
		}

		// register functions to the table
		mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"add":   reg,
			"init":  reg,
			"setUp": reg,
			"unix":  lUnix,
		})

		// Adding RSA module
		rsa.RegisterType(L)

		// Adding JSON module
		jsonencoder.RegisterType(L)

		// returns the module
		L.Push(mod)
		return 1
	})
}
