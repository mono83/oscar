package lua

import (
	"github.com/mono83/oscar"
	"github.com/stretchr/testify/assert"
	"testing"
)

var shaSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("SHA256 hash", function(tc)
	tc:info(tc:sha256Hex("test string"));
end)
`

func TestSHA256(t *testing.T) {
	// Initializing events collector
	col := &collector{}

	// Building context
	ctx := oscar.NewContext()
	ctx.Register(col.OnEvent)

	// Loading
	if suite, err := SuiteFromString(ctx, "http", shaSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			if assert.Equal(t, col.TestCaseCount, 1) {
				if assert.Len(t, col.Logs.Info, 1, "collector.Logs.Info") {
					assert.Equal(t, "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b", col.Logs.Info[0])
				}
			}
		}
	}
}
