package oscar

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"regexp"
	"time"
)

// TestCaseMeta contains metatable name for userdata structure TestCase in lua
const TestCaseMeta = "TestCaseType"

// TestCase is main holder of test steps and assertions
type TestCase struct {
	Name     string
	Function *lua.LFunction
	Vars     map[string]string

	CntAssertSuccess, CntAssertFail, CntRemote int
	startedAt, finishedAt                      time.Time

	oscar *Oscar
}

func luaToTestCase(L *lua.LState) *TestCase {
	v := L.CheckUserData(1).Value
	if tc, ok := v.(*TestCase); ok {
		return tc
	}

	return nil
}

func (t *TestCase) lSelf(L *lua.LState) lua.LValue {
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable(TestCaseMeta))
	return ud
}

// Elapsed returns elapsed time
func (t *TestCase) Elapsed() time.Duration {
	return t.finishedAt.Sub(t.startedAt)
}

func (t *TestCase) printAftermath() {
	delta := t.finishedAt.Sub(t.startedAt)
	if t.CntAssertFail > 0 {
		t.logError(fmt.Sprintf("Test case failed. Success: %d, failed %d", t.CntAssertSuccess, t.CntAssertFail))
	} else {
		t.logInfo(fmt.Sprintf("Test case done in %.2f sec. Assertions: %d", delta.Seconds(), t.CntAssertSuccess))
	}
}

// Run starts all assertions and operations within test case
func (t *TestCase) Run(L *lua.LState) (err error) {
	t.logInfo("Running test case " + t.Name)
	t.oscar.tracef("Invoking %s", t.Name)
	defer func() {
		if r := recover(); r != nil {
			t.oscar.tracef("Recovered panic %+v", r)
			t.CntAssertFail++
			err = fmt.Errorf("%+v", r)
		}

		t.finishedAt = time.Now()
		t.printAftermath()
	}()
	t.startedAt = time.Now()
	L.Push(t.Function)
	L.Push(t.lSelf(L))
	L.Call(1, 0)

	return nil
}

// Get returns variable value from vars map
func (t *TestCase) Get(key string) string {
	if len(t.Vars) > 0 {
		if v, ok := t.Vars[key]; ok {
			return v
		}
	}

	return t.oscar.Get(key)
}

// Set assigns new variable value
func (t *TestCase) Set(key, value string) {
	t.oscar.tracef(`Setting "%s" := "%s"`, key, value)
	if len(t.Vars) == 0 {
		t.Vars = map[string]string{}
	}

	t.Vars[key] = value
}

var iregex = regexp.MustCompile(`\${([\w.-]+)}`)

// Interpolate replaces all placeholders in provided string using vars from test case or
// global runner
func (t *TestCase) Interpolate(value string) string {
	return iregex.ReplaceAllStringFunc(value, func(i string) string {
		m := iregex.FindStringSubmatch(i)
		return t.Get(m[1])
	})
}

// assertDone registers assert attempt
func (t *TestCase) assertDone(success bool) {
	if success {
		t.CntAssertSuccess++
	} else {
		t.CntAssertFail++
	}
}
