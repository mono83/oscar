package lua

import (
	"fmt"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/yuin/gopher-lua"
)

type testcase struct {
	id       int
	name     string
	dep      []int
	deps     []string
	state    *lua.LState
	function *lua.LFunction
	err      error
}

// ID returns test case name and identifier
func (t *testcase) ID() (int, string) {
	return t.id, t.name
}

// GetDependsOn returns slice of identifiers, that must succeed before case will run
func (t *testcase) GetDependsOn() []int {
	return t.dep
}

func (t *testcase) Assert(c *oscar.Context) (err error) {
	c.Register(func(emitted *events.Emitted) {
		events.IfIsAssertDone(emitted, func(s events.AssertDone) {
			if s.Error != nil && t.err == nil {
				t.err = s.Error
			}
		})
	})

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			c.AssertFinished(err)
		}
	}()

	// Building wrapper over context
	udc := t.state.NewUserData()
	udc.Value = c
	t.state.SetMetatable(udc, t.state.GetTypeMetatable(TestCaseMeta))

	// Injecting and invoking
	t.state.Push(t.function)
	t.state.Push(udc)
	t.state.Call(1, 0)

	c.Wait()

	return t.err
}
