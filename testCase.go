package oscar

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"time"
)

// TestCaseMeta contains metatable name for userdata structure TestCase in lua
const TestCaseMeta = "TestCaseType"

// TestCase is main holder of test steps and assertions
type TestCase struct {
	*TestContext
	Name     string
	Function *lua.LFunction
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

// Run starts all assertions and operations within test case
func (t *TestCase) Run(L *lua.LState) (err error) {
	t.Emit(StartEvent{Time: time.Now(), Owner: t})
	t.Trace("Invoking %s", t.Name)
	defer func() {
		if r := recover(); r != nil {
			t.Trace("Recovered panic %+v", r)
			err = fmt.Errorf("%+v", r)
			if t.Error == nil {
				t.assertDone(err)
			} else {
				t.Error = fmt.Errorf("%s\n%s", t.Error.Error(), err)
			}
		}
	}()
	L.Push(t.Function)
	L.Push(t.lSelf(L))
	L.Call(1, 0)

	t.Emit(FinishEvent{Time: time.Now(), Owner: t})
	return nil
}

// assertDone registers assert attempt
func (t *TestCase) assertDone(err error) {
	if err == nil {
		t.Emit(AssertionSuccess(""))
	} else {
		t.Emit(AssertionFailure(err))
		t.Emit(FinishEvent{Time: time.Now(), Owner: t})
	}
}
