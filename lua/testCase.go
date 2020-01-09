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
		events.IfFailure(emitted, func(f events.Failure) {
			t.err = f
		})
	})

	defer func() {
		if r := recover(); r != nil {
			// Unfolding panic
			var message string
			if e, ok := r.(error); ok {
				// This is an error
				message = e.Error()
			} else {
				// This is not an error
				message = fmt.Sprintf("%+v", r)
			}
			c.Fail(message)
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
