package lua

import (
	"testing"

	"github.com/mono83/oscar"
	"github.com/stretchr/testify/assert"
)

var assertEqualsFailSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("Simple fail", function(tc)
	tc:assertEquals("1", "2");
end)
`

func TestAssertEqualsFail(t *testing.T) {
	// Initializing events collector
	col := &collector{}

	// Building context
	ctx := oscar.NewContext()
	ctx.Register(col.OnEvent)
	ctx.Import(map[string]string{
		"who": "Test runner",
	})

	// Loading
	if suite, err := SuiteFromString(ctx, "fooo", assertEqualsFailSuite); assert.NoError(t, err) {
		// Running suite
		err := oscar.RunSequential(ctx, []oscar.Suite{suite})
		if assert.Error(t, err) {
			if assert.Equal(t, 1, col.AssertionsCount) {
				if assert.Len(t, col.Failures, 1) {
					assert.Equal(
						t,
						"<string>:6: Assertion failed. \"1\" (actual, left) != \"2\".\nstack traceback:\n\t[G]: in function 'assertEquals'\n\t<string>:6: in main chunk\n\t[G]: ?",
						col.Failures[0].Error(),
					)
				}
			}
		}
	}
}
