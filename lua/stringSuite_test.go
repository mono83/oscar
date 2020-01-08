package lua

import (
	"fmt"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

var simpleSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("Simple log", function(tc)
	tc:info("Hello, world");
end)

o.add("Simple log with environment variable", function(tc)
	tc:info("Hello, ${who}");
end)

o.add("Variable setter", function(tc)
	tc:set("custom", "value-foo");
	tc:info("This is ${custom}");
end)

o.add("Variable getter and setter", function(tc)
	tc:info("Variable $custom contains ${custom}");
	tc:set("custom", tc:get("who") .. " again");
	tc:info("This is ${custom}");
end)

o.add("Assertions success", function(tc)
	tc:assert("foo", "foo");
	tc:assert(1, 1);
end)

o.add("Assertion JSON XPATH BODY", function(tc)
	tc:set("http.response.body", '{"user": [{"id": 42}]}');
	tc:assertJSONXPath("$.user[0].id", "42");
end)
`

func TestStringSuite(t *testing.T) {
	// Initializing events collector
	col := &collector{}

	// Building context
	ctx := oscar.NewContext()
	ctx.Register(col.OnEvent)
	ctx.Import(map[string]string{
		"who": "Test runner",
	})

	// Loading
	if suite, err := SuiteFromString(ctx, "fooo", simpleSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			// Checking test cases count
			if assert.Equal(t, col.TestCaseCount, 6) {
				// Checking data from logs
				if assert.Len(t, col.Logs.Info, 5, "collector.Logs.Info") {
					assert.Equal(t, "Hello, world", col.Logs.Info[0])
					assert.Equal(t, "Hello, Test runner", col.Logs.Info[1])
					assert.Equal(t, "This is value-foo", col.Logs.Info[2])
					assert.Equal(t, "Variable $custom contains ", col.Logs.Info[3])
					assert.Equal(t, "This is Test runner again", col.Logs.Info[4])
				}
				if assert.Len(t, col.Logs.Trace, 5, "collector.Logs.Trace") {
					assert.Equal(t, `Assert "foo" (actual, left) equals "foo"`, col.Logs.Trace[0])
					assert.Equal(t, `Assert "1" (actual, left) equals "1"`, col.Logs.Trace[1])
					assert.Equal(t, `Reading JSON XPath "$.user[0].id"`, col.Logs.Trace[2])
					assert.Equal(t, `Assert "42" (actual, left) equals "42"`, col.Logs.Trace[3])
					assert.Equal(t, `Assertion OK. "$.user[0].id" == "42"`, col.Logs.Trace[4])
				}
			}
		}
	}
}

type collector struct {
	Debug bool // If true outputs all logs to stdout

	TestCaseCount int // Contains amount of test cases been run

	Logs struct {
		Trace []string // Contains all messages with level 0 (trace)
		Info  []string // Contains all messages with level 2 (info)
	}
}

func (c *collector) OnEvent(e *events.Emitted) {
	if e != nil && e.Data != nil {
		switch e.Data.(type) {
		case events.Start:
			start := e.Data.(events.Start)
			if start.Type == "TestCase" {
				c.TestCaseCount++
			}
		case events.LogEvent:
			log := e.Data.(events.LogEvent)
			switch log.Level {
			case 0:
				c.Logs.Trace = append(c.Logs.Trace, log.Pattern)
			case 2:
				c.Logs.Info = append(c.Logs.Info, log.Pattern)
			default:
				// nothing, just ignore
			}

			if c.Debug {
				fmt.Println(log.Level, log.Pattern)
			}
		}
	}
}
