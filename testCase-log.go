package oscar

import (
	"github.com/yuin/gopher-lua"
)

func lTestCaseDebug(L *lua.LState) int {
	if t := luaToTestCase(L); t != nil {
		t.Emit(TestLogEvent{Level: 0, Owner: t, Message: L.ToString(2)})
	}
	return 0
}

func lTestCaseInfo(L *lua.LState) int {
	if t := luaToTestCase(L); t != nil {
		t.Emit(TestLogEvent{Level: 1, Owner: t, Message: L.ToString(2)})
	}
	return 0
}
