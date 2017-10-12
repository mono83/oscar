package oscar

import (
	"fmt"
	"github.com/mono83/oscar/util/rsa"
	"github.com/yuin/gopher-lua"
	"time"
)

// TestSuite is main test runner
type TestSuite struct {
	*TestContext
	Include []string

	Cases []*TestCase

	CaseSelector func(*TestCase) bool
}

// StartFile starts test cases from lua file
func (o *TestSuite) StartFile(file string) error {
	L := lua.NewState()
	defer L.Close()

	if len(o.Vars) == 0 {
		o.Vars = map[string]string{}
	}

	o.Vars["lua.engine"] = "TestSuite"

	// Loading module
	o.InjectModule(L)

	if len(o.Include) > 0 {
		for _, h := range o.Include {
			o.Trace("Reading header %s", h)
			if err := L.DoFile(h); err != nil {
				return err
			}
		}
	}

	// Loading file
	o.Trace("Reading file %s", file)
	before := time.Now()
	if err := L.DoFile(file); err != nil {
		return err
	}
	o.Trace("File parsed in %.1fms", time.Now().Sub(before).Seconds()*1000)

	// Running tests
	return o.Start(L)
}

// InjectModule injects TestSuite module (named "oscar") into lua engine
func (o *TestSuite) InjectModule(L *lua.LState) {
	L.PreloadModule("oscar", func(L *lua.LState) int {
		// Registering test case type
		mt := L.NewTypeMetatable(TestCaseMeta)
		// methods
		L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"assert":          lTestCaseAssert,
			"assertEq":        lTestCaseAssert,
			"assertEquals":    lTestCaseAssert,
			"assertJSONPath":  lTestCaseAssertJSONXPath,
			"assertJSONXPath": lTestCaseAssertJSONXPath,
			"get":             lTestCaseGet,
			"httpPost":        lTestCaseHTTPPost,
			"jsonPath":        lJSONPathExtract,
			"jsonXPath":       lJSONPathExtract,
			"set":             lTestCaseSet,

			"stringToBase64":    lTestCaseStringBase64,
			"packInt64ToBase64": lTestCasePackInt64Base64,

			"log":  lTestCaseDebug,
			"info": lTestCaseInfo,
		}))

		// register functions to the table
		mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"add": o.lAdd,
		})

		// Adding RSA module
		rsa.RegisterType(L)

		// returns the module
		L.Push(mod)
		return 1
	})
}

// Start begins all tests
func (o *TestSuite) Start(L *lua.LState) error {
	o.Emit(StartEvent{Time: time.Now(), Owner: o})
	o.Trace("Starting tests")
	for i, s := range o.Cases {
		if o.CaseSelector == nil || o.CaseSelector(s) {
			o.Trace("Starting test case #%d - \"%s\"", i+1, s.Name)
			s.Run(L)
		} else {
			o.Trace("Test case %s skipped by case selector predicate", s.Name)
		}
	}

	o.Emit(FinishEvent{Time: time.Now(), Owner: o})

	return o.GetError()
}

// GetError returns overall total error for test runner
func (o *TestSuite) GetError() (err error) {
	// Choosing error
	for _, s := range o.Cases {
		if s.Error != nil {
			err = fmt.Errorf("at least one test case failure in %s", s.Name)
			break
		}
	}

	return
}

// lAdd registers test case from lua callback function and name
func (o *TestSuite) lAdd(L *lua.LState) int {
	name := L.CheckString(1)
	clb := L.CheckFunction(2)

	o.Trace("Registering test case %s", name)
	o.Cases = append(
		o.Cases,
		&TestCase{
			Name:     name,
			Function: clb,
			TestContext: &TestContext{
				Parent: o.TestContext,
			},
		},
	)
	return 0
}

type nop struct {
}

func (nop) Write(p []byte) (n int, err error) { return }
