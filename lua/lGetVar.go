package lua

import "github.com/yuin/gopher-lua"

// lGetVar returns variable value from context
func lGetVar(L *lua.LState) int {
	tc := lContext(L)
	key := L.CheckString(2)

	v := tc.Get(key)

	L.Push(lua.LString(v))
	return 1
}
