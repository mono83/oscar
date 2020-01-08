package lua

import (
	"github.com/mono83/oscar/util/jsonpath"
	"github.com/yuin/gopher-lua"
)

func lJSONPathExtract(L *lua.LState) int {
	tc := lContext(L)
	xpath := tc.Interpolate(L.CheckString(2))
	source := tc.Interpolate(L.CheckString(3))

	tc.Tracef(`Reading JSON XPath "%s"`, xpath)

	value, err := jsonpath.Extract([]byte(source), xpath)
	if err != nil {
		lRaiseContextError(L, tc, "JSON XPATH error: %s", err.Error())
		return 0
	}

	L.Push(lua.LString(value))
	return 1
}
