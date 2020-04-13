package lua

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mono83/oscar"
	"github.com/stretchr/testify/assert"
)

var sqlLongSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("MySQL single long", function(tc)
	tc:export("res", tc:sqlGetLong("SELECT count(1) FROM users WHERE type = ? AND createdAt < ?", "staff", "1000"))
end)
`

func TestLSQLLong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(
		"SELECT count\\(1\\) FROM users WHERE type = \\? AND createdAt < \\?",
	).WithArgs(
		"staff",
		"1000",
	).WillReturnRows(mock.NewRows([]string{"count(1)"}).AddRow(42))

	ctx := oscar.NewContext()
	ctx.SetDatabase(db)

	// Loading
	if suite, err := SuiteFromString(ctx, "sql", sqlLongSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			// we make sure that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())

			// Check response
			assert.Equal(t, "42", ctx.Get("res"))
		}
	}
}

var sqlStringSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("MySQL single long", function(tc)
	tc:export("res", tc:sqlGetString("SELECT name FROM users WHERE status = ? AND type = ? LIMIT 1", "enabled", "staff"))
end)
`

func TestLSQLString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(
		"SELECT name FROM users WHERE status = \\? AND type = \\? LIMIT 1",
	).WithArgs(
		"enabled",
		"staff",
	).WillReturnRows(mock.NewRows([]string{"name"}).AddRow("Hello, world"))

	ctx := oscar.NewContext()
	ctx.SetDatabase(db)

	// Loading
	if suite, err := SuiteFromString(ctx, "sql", sqlStringSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			// we make sure that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())

			// Check response
			assert.Equal(t, "Hello, world", ctx.Get("res"))
		}
	}
}

var sqlRowSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("MySQL single long", function(tc)
	local dat = tc:sqlGetRow("SELECT id, name FROM users WHERE status = ? AND type = ? LIMIT 1", "enabled", "staff")
	tc:export("recid", dat.id)
	tc:export("recname", dat.name)
end)
`

func TestLSQLRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(
		"SELECT id, name FROM users WHERE status = \\? AND type = \\? LIMIT 1",
	).WithArgs(
		"enabled",
		"staff",
	).WillReturnRows(mock.NewRows([]string{"id", "name"}).AddRow(31, "Admin"))

	ctx := oscar.NewContext()
	ctx.SetDatabase(db)

	// Loading
	if suite, err := SuiteFromString(ctx, "sql", sqlRowSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			// we make sure that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())

			// Check response
			assert.Equal(t, "31", ctx.Get("recid"))
			assert.Equal(t, "Admin", ctx.Get("recname"))
		}
	}
}
