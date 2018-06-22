package lua

import (
	"github.com/yuin/gopher-lua"
	"time"
)

// lUnix returns Unix timestamp in seconds
func lUnix(L *lua.LState) int {
	L.Push(lua.LNumber(float64(time.Now().Unix())))
	return 1
}
