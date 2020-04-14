package rsa

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/mono83/oscar/util"
	"github.com/yuin/gopher-lua"
)

// MetaTableName contains name of entry in Lua meta space for RSA structure
const MetaTableName = "RSAWrapper"

// RegisterType registers RSA user data class in Lua space
func RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(MetaTableName)

	L.SetGlobal("RSA", mt)
	L.SetField(mt, "new", L.NewFunction(lCreate))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"exportCertificate":   lExportCertificate,
		"certificateExport":   lExportCertificate,
		"exportCertificate64": lExportCertificateB64,
		"certificateExport64": lExportCertificateB64,
		"signSHA256String64":  lSign256StringB64,
		"signSHA256Int64":     lSign256BLongSliceB64,
		"signSHA256Long64":    lSign256BLongSliceB64,
	}))
}

func luaToCertificate(L *lua.LState) *RSA {
	v := L.CheckUserData(1).Value
	if crt, ok := v.(*RSA); ok {
		return crt
	}

	return nil
}

func lCreate(L *lua.LState) int {
	keyLen := L.CheckInt(1)

	crt, err := Generate(keyLen, time.Hour)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	ud := L.NewUserData()
	ud.Value = crt
	L.SetMetatable(ud, L.GetTypeMetatable(MetaTableName))
	L.Push(ud)
	return 1
}

func lExportCertificate(L *lua.LState) int {
	crt := luaToCertificate(L)
	pem := crt.EncodedCertificate()
	L.Push(lua.LString(pem))
	return 1
}

func lExportCertificateB64(L *lua.LState) int {
	crt := luaToCertificate(L)
	pem := crt.EncodedCertificate()
	pem = base64.StdEncoding.EncodeToString([]byte(pem))
	L.Push(lua.LString(pem))
	return 1
}

func lSign256StringB64(L *lua.LState) int {
	crt := luaToCertificate(L)
	payload := L.CheckString(2)

	value, err := crt.SignSHA256B64([]byte(payload))
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(lua.LString(value))
	return 1
}

func lSign256BLongSliceB64(L *lua.LState) int {
	crt := luaToCertificate(L)
	hash := ""
	if cnt := L.GetTop(); cnt > 1 {
		// Packing long values to bytes
		buf := make([]byte, (cnt-1)*8)
		for i := 2; i <= cnt; i++ {
			chunk := make([]byte, 8)
			i64, err := util.LuaToInt64(L, i)
			if err != nil {
				L.Error(lua.LString(err.Error()), 1)
				return 0
			}
			binary.BigEndian.PutUint64(chunk, uint64(i64))
			for j := 0; j < 8; j++ {
				buf[(i-2)*8+j] = chunk[j]
			}
		}

		var err error
		hash, err = crt.SignSHA256B64(buf)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	}

	L.Push(lua.LString(hash))
	return 1
}
