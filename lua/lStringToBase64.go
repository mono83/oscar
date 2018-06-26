package lua

import (
	"encoding/base64"
	"github.com/yuin/gopher-lua"
)

func lStringToBase64(L *lua.LState) int {
	tc := lContext(L)
	value := tc.Interpolate(L.ToString(2))

	b64 := base64.StdEncoding.EncodeToString([]byte(value))
	L.Push(lua.LString(b64))
	return 1
}
