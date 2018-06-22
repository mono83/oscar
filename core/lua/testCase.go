package lua

import (
	"fmt"
	"github.com/mono83/oscar/core"
	"github.com/mono83/oscar/core/events"
	"github.com/yuin/gopher-lua"
)

type testcase struct {
	id       int
	name     string
	state    *lua.LState
	function *lua.LFunction
	err      error
}

// ID returns test case name and identifier
func (t *testcase) ID() (int, string) {
	return t.id, t.name
}

func (t *testcase) Assert(c *core.Context) (err error) {
	proxy := c.Fork()
	proxy.OnEvent = t.assertDoneInterceptor

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			proxy.AssertFinished(err)
		}
	}()

	// Building wrapper over context
	udc := t.state.NewUserData()
	udc.Value = proxy
	t.state.SetMetatable(udc, t.state.GetTypeMetatable(TestCaseMeta))

	// Injecting and invoking
	t.state.Push(t.function)
	t.state.Push(udc)
	t.state.Call(1, 0)

	proxy.Wait()

	return t.err
}

func (t *testcase) assertDoneInterceptor(i interface{}) {
	if s, ok := i.(events.AssertDone); ok {
		if s.Error != nil && t.err == nil {
			t.err = s.Error
		}
	}
}
