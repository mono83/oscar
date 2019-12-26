package lua

import (
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/mono83/oscar/util/jsonencoder"
	"github.com/mono83/oscar/util/rsa"
	lua "github.com/yuin/gopher-lua"
	"strings"
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

// InjectModule injects TestSuite module (named "oscar") into lua engine
func (t *testSuite) InjectModule(ctx *oscar.Context, L *lua.LState) {
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
			"fail":            lFail,
			"get":             lGetVar,
			"httpGet":         lHTTPGet,
			"httpPost":        lHTTPPost,
			"jsonPath":        lJSONPathExtract,
			"jsonXPath":       lJSONPathExtract,
			"set":             lSetVar,

			"unix":      lUnix,
			"sleep":     lSleep,
			"sha256Hex": lSHA256Hex,

			"stringToBase64":    lStringToBase64,
			"packInt64ToBase64": lPackSliceInt64ToBase64,

			"log":   lDebug,
			"debug": lDebug,
			"info":  lInfo,
		}))

		// Making lambdas
		clbRegCase := func(L *lua.LState) int {
			name := L.CheckString(1)
			clb := L.CheckFunction(2)

			id := id()

			ctx.Emit(events.RegistrationBegin{Type: "TestCase", ID: id, Name: name})

			c := &testcase{
				id:       id,
				name:     name,
				function: clb,
				state:    t.state,
			}

			t.cases = append(t.cases, c)

			if L.GetTop() > 2 {
				L.CheckTable(3).ForEach(func(key lua.LValue, value lua.LValue) {
					if key.Type() == lua.LTString {
						keyStr := strings.ToLower(strings.TrimSpace(key.String()))

						switch keyStr {
						case "depends_on", "dependson", "depends", "deps", "dep":
							if value.Type() == lua.LTString {
								c.deps = []string{value.String()}
							} else if value.Type() == lua.LTTable {
								value.(*lua.LTable).ForEach(func(key lua.LValue, value lua.LValue) {
									c.deps = append(c.deps, value.String())
								})
							} else {
								L.RaiseError("unsupported value for test dependency")
							}
						}
					}
				})
			}

			ctx.Emit(events.RegistrationEnd{Type: "TestCase", Name: name})

			return 0
		}

		clbRegSetup := func(L *lua.LState) int {
			clb := L.CheckFunction(1)

			id := id()

			ctx.Emit(events.RegistrationBegin{Type: "TestSuiteInit", ID: id, Name: oscar.SuiteSetUp})

			t.setup = &testcase{
				id:       id,
				name:     oscar.SuiteSetUp,
				function: clb,
				state:    t.state,
			}

			ctx.Emit(events.RegistrationEnd{Type: "TestSuiteInit", Name: oscar.SuiteSetUp})

			return 0
		}

		// register functions to the table
		mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"add":   clbRegCase,
			"init":  clbRegSetup,
			"setUp": clbRegSetup,
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
