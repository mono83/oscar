package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// lSQLRow reads single row and returns it's value as map
// If multiple lines were in resultset, only last value is returned
func lSQLRow(L *lua.LState) int {
	rows, err := lSQLWrapper(L)
	if err != nil {
		L.RaiseError("SQL Error: %s", err.Error())
		return 0
	}
	defer rows.Close()

	// Reading data
	cols, err := rows.Columns()
	if err != nil {
		L.RaiseError("SQL Columns Error: %s", err.Error())
		return 0
	}
	pointers := make([]interface{}, len(cols))
	container := make([]string, len(cols))
	for i := range pointers {
		pointers[i] = &container[i]
	}
	success := false
	for rows.Next() {
		if err := rows.Scan(pointers...); err != nil {
			L.RaiseError("SQL RowScan Error: %s", err.Error())
			return 0
		}

		success = true
	}

	if !success {
		L.RaiseError("No data received from database")
		return 0
	}

	// Converting to table
	tab := &lua.LTable{}
	for i := 0; i < len(cols); i++ {
		tab.RawSetString(cols[i], lua.LString(container[i]))
	}
	L.Push(tab)
	return 1
}
