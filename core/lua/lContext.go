package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/mono83/oscar/core"
)

// lContext is the most important function in Go-Lua intercommunication
// It reads one of stack arguments and converts it into *core.Context
// It will panic on failure
func lContext(L *lua.LState) *core.Context {
	v := L.CheckUserData(1).Value
	if tc, ok := v.(*core.Context); ok {
		return tc
	}

	panic("Unable to read testing context from Lua stack")

	return nil
}
