package lua

import (
	"github.com/mono83/oscar"
	"github.com/yuin/gopher-lua"
)

// lContext is the most important function in Go-Lua intercommunication
// It reads one of stack arguments and converts it into *core.Context
// It will panic on failure
func lContext(L *lua.LState) *oscar.Context {
	v := L.CheckUserData(1).Value
	if tc, ok := v.(*oscar.Context); ok {
		return tc
	}

	panic("Unable to read testing context from Lua stack")
}

// lRaiseContextError emits error message to logs
// and raises standard lua error
func lRaiseContextError(L *lua.LState, ctx *oscar.Context, pattern string, args ...interface{}) {
	// Emitting error message
	ctx.Errorf(pattern, args...)

	// Raising Lua error
	L.RaiseError(pattern, args...)
}
