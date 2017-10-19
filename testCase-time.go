package oscar

import (
	"github.com/yuin/gopher-lua"
	"time"
)

// lUnix returns Unix timestamp in seconds
func lUnix(L *lua.LState) int {
	L.Push(lua.LNumber(float64(time.Now().Unix())))
	return 1
}

// lSleep pauses execution for requested amount of milliseconds
func lSleep(L *lua.LState) int {
	tc := luaToTestCase(L)
	milliseconds := L.ToInt(2)
	duration := time.Millisecond * time.Duration(milliseconds)
	tc.Trace("Sleeping for %d milliseconds", milliseconds)
	time.Sleep(duration)
	tc.Emit(SleepEvent(duration))

	return 0
}
