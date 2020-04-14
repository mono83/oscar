package jsonencoder

import (
	"encoding/json"
	"errors"

	"github.com/yuin/gopher-lua"
)

type jsonValue struct {
	lua.LValue
}

// MarshalJSON marshal themselves into JSON
func (j jsonValue) MarshalJSON() ([]byte, error) {
	return toJSON(j.LValue)
}

// toJSON marshal into JSON
func toJSON(value lua.LValue) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool, lua.LString, lua.LNumber, *lua.LNilType:
		data, err = json.Marshal(converted)

	case *lua.LTable:
		var array []jsonValue
		var hashmap map[string]jsonValue
		var isHashmap bool

		// Detecting hashmap type
		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			if _, numberKey := k.(lua.LNumber); !numberKey {
				isHashmap = true
				return
			}
		})

		if isHashmap {
			hashmap = make(map[string]jsonValue)
		}

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			if isHashmap {
				hashmap[k.String()] = jsonValue{v}
			} else {
				array = append(array, jsonValue{v})
			}
		})

		if isHashmap {
			data, err = json.Marshal(hashmap)
		} else {
			data, err = json.Marshal(array)
		}

	default:
		err = errors.New("cannot encode data")
	}
	return
}

// fromJSON unmarshal data
func fromJSON(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromJSON(L, item))
		}
		return arr
	case map[string]interface{}:
		table := L.CreateTable(0, len(converted))
		for key, item := range converted {
			table.RawSetH(lua.LString(key), fromJSON(L, item))
		}
		return table
	}

	return lua.LNil
}
