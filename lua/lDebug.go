package lua

import "github.com/yuin/gopher-lua"

// lDebug sends logs with DEBUG level and uses variables interpolation
func lDebug(L *lua.LState) int {
	if t := lContext(L); t != nil {
		t.Debug(L.ToString(2))
	}
	return 0
}
