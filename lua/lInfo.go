package lua

import "github.com/yuin/gopher-lua"

// lInfo sends logs with INFO level and uses variables interpolation
func lInfo(L *lua.LState) int {
	if t := lContext(L); t != nil {
		t.Info(L.ToString(2))
	}
	return 0
}
