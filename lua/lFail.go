package lua

import (
	"github.com/mono83/oscar"
	"github.com/yuin/gopher-lua"
)

// lFail emits test case assertion failure with custom text
func lFail(L *lua.LState) int {
	tc := lContext(L)

	var msg string
	var args []interface{}

	msg = "Test case failed"

	if L.GetTop() == 2 {
		msg = tc.Interpolate(L.ToString(2))
	} else if l := L.GetTop(); l > 2 {
		msg = L.ToString(2)
		for i := 3; i <= l; i++ {
			args = append(args, L.ToString(i))
		}
	}

	throwLua(L, tc, msg, args...)

	return 0
}

// throwLua emits error message to logs and raises standard lua error
func throwLua(L *lua.LState, ctx *oscar.Context, pattern string, args ...interface{}) {
	// Emitting error message
	ctx.Errorf(pattern, args...)

	// Raising Lua error
	L.RaiseError(pattern, args...)
}
