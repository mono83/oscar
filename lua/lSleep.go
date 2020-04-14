package lua

import (
	"time"

	"github.com/yuin/gopher-lua"
)

// lSleep pauses execution for requested amount of milliseconds
func lSleep(L *lua.LState) int {
	tc := lContext(L)
	milliseconds := L.ToInt(2)
	duration := time.Millisecond * time.Duration(milliseconds)
	tc.Sleep(duration)

	return 0
}
