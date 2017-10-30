package jsonEncoder

import (
	"github.com/yuin/gopher-lua"
	"encoding/json"
)

func RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable("JSONWrapper")

	L.SetGlobal("JSON", mt)
	L.SetField(mt, "decode", L.NewFunction(lDecode))
	L.SetField(mt, "encode", L.NewFunction(lEncode))
}

func lDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}

	if err := json.Unmarshal([]byte(str), &value); err != nil {
		L.RaiseError("Unable to decode json - %s", err.Error())
	}

	L.Push(fromJSON(L, value))

	return 1
}

func lEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	data, err := toJSON(value)

	if err != nil {
		L.RaiseError("Unable to encode to json - %s", err.Error())
	}

	L.Push(lua.LString(string(data)))

	return 1
}
