package oscar

import "github.com/yuin/gopher-lua"

func lTestCaseGet(L *lua.LState) int {
	tc := luaToTestCase(L)
	key := L.CheckString(2)

	v := tc.Get(key)

	L.Push(lua.LString(v))
	return 1
}
