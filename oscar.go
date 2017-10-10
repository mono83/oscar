package oscar

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/yuin/gopher-lua"
	"io"
	"time"
)

// Oscar is main test runner
type Oscar struct {
	Include []string
	Debug   bool
	Vars    map[string]string

	Output io.Writer
	Cases  []*TestCase

	CaseSelector func(*TestCase) bool
}

// StartFile starts test cases from lua file
func (o *Oscar) StartFile(file string) error {
	L := lua.NewState()
	defer L.Close()

	if len(o.Vars) == 0 {
		o.Vars = map[string]string{}
	}

	o.Vars["lua.engine"] = "Oscar"

	// Loading module
	o.InjectModule(L)

	if len(o.Include) > 0 {
		for _, h := range o.Include {
			o.tracef("Reading header %s", h)
			if err := L.DoFile(h); err != nil {
				return err
			}
		}
	}

	// Loading file
	o.tracef("Reading file %s", file)
	before := time.Now()
	if err := L.DoFile(file); err != nil {
		return err
	}
	o.tracef("File parsed in %.1fms", time.Now().Sub(before).Seconds()*1000)

	// Running tests
	return o.Start(L)
}

// InjectModule injects Oscar module (named "oscar") into lua engine
func (o *Oscar) InjectModule(L *lua.LState) {
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
			"log":             lTestCaseDebug,
			"info":            lTestCaseInfo,
		}))

		// register functions to the table
		mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"add": o.lAdd,
		})

		// returns the module
		L.Push(mod)
		return 1
	})
}

// Start begins all tests
func (o *Oscar) Start(L *lua.LState) (err error) {
	o.tracef("Starting tests")
	for i, s := range o.Cases {
		if o.CaseSelector == nil || o.CaseSelector(s) {
			o.tracef("Starting test case #%d - \"%s\"", i+1, s.Name)
			s.Run(L)
		} else {
			o.tracef("Test case %s skipped by case selector predicate", s.Name)
		}
	}

	// Choosing error
	for _, s := range o.Cases {
		if s.CntAssertFail > 0 {
			err = fmt.Errorf("at least one test case failure in %s", s.Name)
			break
		}
	}

	if err != nil {
		// Printing error details
		fmt.Fprintln(o.output())
		fmt.Fprintln(o.output(), " Errors:")
		i := 1
		for _, s := range o.Cases {
			if s.Error != nil {
				fmt.Fprintf(o.output(), "  %d. %s\n", i, s.Name)
				fmt.Fprintln(o.output(), "     ", s.Error)
				for k, v := range s.Vars {
					fmt.Fprintln(o.output(), "      ", k, ":=", v)
				}
				fmt.Fprintln(o.output())
				i++
			}
		}
		fmt.Fprintln(o.output())
	}

	// Building global aftermath
	longest := len("Test suite")
	for _, s := range o.Cases {
		if s.CntAssertFail > 0 || s.CntAssertSuccess > 0 {
			if l := len(s.Name); l > longest {
				longest = l
			}
		}
	}

	namePattern := fmt.Sprintf(" %%-%ds", longest)
	fullPattern := "%s" + namePattern + "  %5d   %5d     %5d   %7.1fms\n"

	fmt.Fprintln(o.output())
	fmt.Fprintln(o.output())
	fmt.Fprintf(
		o.output(),
		"      "+namePattern+" Success  Failed  Requests  Time spent\n",
		"Test suite",
	)
	fmt.Fprintln(o.output())

	for _, s := range o.Cases {
		if s.CntAssertFail > 0 || s.CntAssertSuccess > 0 {
			status := colorOscarSummarySuccess.Sprint("  OK  ")
			if s.CntAssertFail > 0 {
				status = colorOscarSummaryFailed.Sprint(" FAIL ")
			}

			fmt.Fprintf(
				o.output(),
				fullPattern,
				status,
				s.Name,
				s.CntAssertSuccess,
				s.CntAssertFail,
				s.CntRemote,
				s.Elapsed().Seconds()*1000,
			)
		}
	}

	fmt.Fprintln(o.output())
	fmt.Fprintln(o.output())

	return
}

// Get returns variable value from vars map
func (o *Oscar) Get(key string) string {
	if len(o.Vars) > 0 {
		if v, ok := o.Vars[key]; ok {
			return v
		}
	}

	return ""
}

// tracef output debug information (if enabled) using printf syntax
func (o *Oscar) tracef(message string, args ...interface{}) {
	if o.Debug {
		fmt.Fprintf(o.output(), message+"\n", args...)
	}
}

// lAdd registers test case from lua callback function and name
func (o *Oscar) lAdd(L *lua.LState) int {
	name := L.CheckString(1)
	clb := L.CheckFunction(2)

	o.tracef("Registering test case %s", name)
	o.Cases = append(o.Cases, &TestCase{Name: name, Function: clb, oscar: o})
	return 0
}

func (o *Oscar) output() io.Writer {
	if o.Output == nil {
		return nop{}
	}

	return o.Output
}

var colorOscarSummarySuccess = color.New(color.FgHiGreen)
var colorOscarSummaryFailed = color.New(color.FgHiRed)

type nop struct {
}

func (nop) Write(p []byte) (n int, err error) { return }
