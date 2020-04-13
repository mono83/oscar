package lua

import "github.com/yuin/gopher-lua"

// lSetVar sets variable value
func lSetVar(L *lua.LState) int {
	tc := lContext(L)
	key := L.CheckString(2)
	value := L.ToString(3)

	tc.Set(key, value)

	return 0
}

// lExportVar sets export variable
func lExportVar(L *lua.LState) int {
	tc := lContext(L)
	key := L.CheckString(2)
	value := L.ToString(3)

	tc.SetExport(key, value)

	return 0
}
