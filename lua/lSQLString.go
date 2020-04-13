package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// lSQLString reads single string and returns it's value as string
// If multiple lines were in resultset, only last value is returned
func lSQLString(L *lua.LState) int {
	rows, err := lSQLWrapper(L)
	if err != nil {
		L.RaiseError("SQL Error: %s", err.Error())
		return 0
	}
	defer rows.Close()

	// Reading data
	var s string
	success := false
	for rows.Next() {
		if err := rows.Scan(&s); err != nil {
			L.RaiseError("SQL RowScan Error: %s", err.Error())
			return 0
		}

		success = true
	}

	if !success {
		L.RaiseError("No data received from database")
		return 0
	}

	// Pushing as string, because Lua may fail with int64 values
	L.Push(lua.LString(s))
	return 1
}
