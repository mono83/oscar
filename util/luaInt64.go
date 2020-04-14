package util

import (
	"strconv"

	"github.com/yuin/gopher-lua"
)

// LuaToInt64 reads Lua function argument as string, and then parses it as int64
func LuaToInt64(L *lua.LState, arg int) (int64, error) {
	str := L.ToString(arg)
	return strconv.ParseInt(str, 10, 64)
}
