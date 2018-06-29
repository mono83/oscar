package lua

import (
	"fmt"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/mono83/oscar/impact"
	"github.com/yuin/gopher-lua"
)

type testcase struct {
	id       int
	name     string
	imp      impact.Level
	state    *lua.LState
	function *lua.LFunction
	err      error
}

// ID returns test case name and identifier
func (t *testcase) ID() (int, string) {
	return t.id, t.name
}

// GetImpact returns impact level, induced by test case on remote infrastructure
func (t *testcase) GetImpact() impact.Level {
	return t.imp
}

// GetDependsOn returns slice of identifiers, that must succeed before case will run
func (t *testcase) GetDependsOn() []int {
	return nil
}

func (t *testcase) Assert(c *oscar.Context) (err error) {
	c.OnEvent = buildAssertDoneInterceptor(t, c.OnEvent)

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

func buildAssertDoneInterceptor(t *testcase, o func(*events.Emitted)) func(*events.Emitted) {
	return func(i *events.Emitted) {
		events.IfIsAssertDone(i, func(s events.AssertDone) {
			if s.Error != nil && t.err == nil {
				t.err = s.Error
			}
		})

		if o != nil {
			o(i)
		}
	}
}
