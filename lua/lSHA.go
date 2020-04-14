package lua

import (
	"crypto/sha256"
	"fmt"

	"github.com/yuin/gopher-lua"
)

// lSHA256Hex returns SHA256 hash of provided string in hexadecimal format
func lSHA256Hex(L *lua.LState) int {
	h := sha256.New()
	h.Write([]byte(L.CheckString(2)))
	L.Push(lua.LString(fmt.Sprintf("%x", h.Sum(nil))))
	return 1
}
