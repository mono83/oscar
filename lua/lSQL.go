package lua

import (
	"database/sql"
	"errors"

	lua "github.com/yuin/gopher-lua"
)

func lSQLWrapper(L *lua.LState) (*sql.Rows, error) {
	tc := lContext(L)
	sql := L.CheckString(2)

	// Obtaining database connection
	db := tc.GetDatabase()
	if db == nil {
		return nil, errors.New("no database connection configured")
	}

	// Reading arguments
	var args []interface{}
	if cnt := L.GetTop(); cnt > 2 {
		for i := 3; i <= cnt; i++ {
			args = append(args, L.ToString(i))
		}
	}

	// Performing query
	return db.Query(sql, args...)
}
