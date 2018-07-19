package lua

import (
	"errors"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/mono83/oscar/impact"
	"github.com/mono83/oscar/util/jsonencoder"
	"github.com/mono83/oscar/util/rsa"
	"github.com/yuin/gopher-lua"
	"strings"
)

// TestCaseMeta contains metatable name for userdata structure TestCase in lua
const TestCaseMeta = "TestCaseType"

// SuiteFromFiles builds suite using Lua sources file
func SuiteFromFiles(c *oscar.Context, files ...string) (oscar.Suite, error) {
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

	registered := map[string]int{}

	ctx := c.Fork(s.id)
	ctx.Register(func(emitted *events.Emitted) {
		if b, ok := emitted.Data.(events.RegistrationBegin); ok {
			registered[b.Name] = b.ID
		}
	})
	s.InjectModule(ctx, L)

	// Emitting registration start event
	ctx.Emit(events.RegistrationBegin{Type: "TestSuite", ID: s.id, Name: s.name})

	// Reading files sequentially
	for _, file := range files {
		if err := L.DoFile(file); err != nil {
			return nil, err
		}
	}

	// Resolving dependencies
	for _, tc := range s.cases {
		if len(tc.deps) > 0 {
			for _, dep := range tc.deps {
				f, ok := registered[dep]
				if !ok {
					return nil, errors.New("unable to find dependency " + dep + " for " + tc.name)
				}
				tc.dep = append(tc.dep, f)
			}
		}
	}

	// Emitting registration done event
	ctx.Emit(events.RegistrationEnd{Type: "TestSuite", Name: s.name})

	return s, nil
}

type fileTestSuite struct {
	id    int
	name  string
	state *lua.LState
	setup *testcase
	cases []*testcase
}

// ID returns test suite name and identifier
func (f *fileTestSuite) ID() (int, string) {
	return f.id, f.name
}

// GetSetUp returns optional setup function, that will be invoked before any other test cases
func (f *fileTestSuite) GetSetUp() oscar.Case {
	if f.setup == nil {
		return nil
	}
	return f.setup
}

// GetCases returns test cases
func (f *fileTestSuite) GetCases() []oscar.Case {
	cs := make([]oscar.Case, len(f.cases))
	for i := range f.cases {
		cs[i] = f.cases[i]
	}

	return cs
}

// InjectModule injects TestSuite module (named "oscar") into lua engine
func (f *fileTestSuite) InjectModule(ctx *oscar.Context, L *lua.LState) {
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
				imp:      impact.Default,
				name:     name,
				function: clb,
				state:    f.state,
			}

			f.cases = append(f.cases, c)

			if L.GetTop() > 2 {
				L.CheckTable(3).ForEach(func(key lua.LValue, value lua.LValue) {
					if key.Type() == lua.LTString {
						keyStr := strings.ToLower(strings.TrimSpace(key.String()))

						switch keyStr {
						case "impact":
							c.imp = impact.ParseOrDefault(value.String())
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

			f.setup = &testcase{
				id:       id,
				imp:      impact.Default,
				name:     oscar.SuiteSetUp,
				function: clb,
				state:    f.state,
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
