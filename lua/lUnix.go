package lua

import (
	"time"

	"github.com/yuin/gopher-lua"
)

// lUnix returns Unix timestamp in seconds
func lUnix(L *lua.LState) int {
	L.Push(lua.LNumber(float64(time.Now().Unix())))
	return 1
}
