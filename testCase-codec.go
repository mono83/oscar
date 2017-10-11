package oscar

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/yuin/gopher-lua"
)

func lTestCaseStringBase64(L *lua.LState) int {
	tc := luaToTestCase(L)
	value := tc.Interpolate(L.ToString(2))

	b64 := base64.StdEncoding.EncodeToString([]byte(value))
	L.Push(lua.LString(b64))
	return 1
}

func lTestCasePackInt64Base64(L *lua.LState) int {
	b64 := ""
	if cnt := L.GetTop(); cnt > 1 {
		// Packing long values to bytes
		buf := make([]byte, (cnt-1)*8)
		for i := 2; i <= cnt; i++ {
			chunk := make([]byte, 8)
			binary.BigEndian.PutUint64(chunk, uint64(L.ToInt64(i)))
			for j := 0; j < 8; j++ {
				buf[(i-2)*8+j] = chunk[j]
			}
		}

		b64 = base64.StdEncoding.EncodeToString(buf)
	}

	L.Push(lua.LString(b64))
	return 1
}
