package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/mono83/oscar/util/jsonPath"
)

func lJSONPathExtract(L *lua.LState) int {
	tc := lContext(L)
	path := tc.Interpolate(L.CheckString(2))
	source := tc.Interpolate(L.CheckString(3))

	value, err := jsonPath.Extract([]byte(source), path)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(lua.LString(value))
	return 1
}
