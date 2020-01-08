// +build http

package lua

import (
	"github.com/mono83/oscar"
	"github.com/stretchr/testify/assert"
	"testing"
)

var httpSuite = `
-- Loading Oscar module
local o = require("oscar")

o.add("Simple GET with JSON", function(tc)
    tc:httpGet("https://httpbin.org/get")

    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.url", "https://httpbin.org/get")
end)


o.add("Simple POST with JSON", function(tc)
    tc:httpPost("https://httpbin.org/post", 'Hello, world')

    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.data", "Hello, world")
end, {depends="Simple GET with JSON"})
`

func TestHTTPMethods(t *testing.T) {
	// Initializing events collector
	col := &collector{Debug: true}

	// Building context
	ctx := oscar.NewContext()
	ctx.Register(col.OnEvent)
	ctx.Import(map[string]string{
		"who": "Test runner",
	})

	// Loading
	if suite, err := SuiteFromString(ctx, "http", httpSuite); assert.NoError(t, err) {
		// Running suite
		if assert.NoError(t, oscar.RunSequential(ctx, []oscar.Suite{suite})) {
			// Checking test cases count
			if assert.Equal(t, col.TestCaseCount, 2) {
				if assert.Len(t, col.Logs.Trace, 14, "collector.Logs.Trace") {
					assert.Equal(t, `Preparing HTTP GET request to https://httpbin.org/get`, col.Logs.Trace[0])
					assert.Equal(t, `HTTP request done in `, col.Logs.Trace[1][0:21])
					assert.Equal(t, `Assert "200" (actual, left) equals "200"`, col.Logs.Trace[2])
					assert.Equal(t, `Assert "true" (actual, left) equals "true"`, col.Logs.Trace[3])
					assert.Equal(t, `Reading JSON XPath "$.url"`, col.Logs.Trace[4])
					assert.Equal(t, `Assert "https://httpbin.org/get" (actual, left) equals "https://httpbin.org/get"`, col.Logs.Trace[5])
					assert.Equal(t, `Assertion OK. "$.url" == "https://httpbin.org/get"`, col.Logs.Trace[6])
					assert.Equal(t, `Preparing HTTP POST request to https://httpbin.org/post`, col.Logs.Trace[7])
					assert.Equal(t, `HTTP request done in `, col.Logs.Trace[8][0:21])
					assert.Equal(t, `Assert "200" (actual, left) equals "200"`, col.Logs.Trace[9])
					assert.Equal(t, `Assert "true" (actual, left) equals "true"`, col.Logs.Trace[10])
					assert.Equal(t, `Reading JSON XPath "$.data"`, col.Logs.Trace[11])
					assert.Equal(t, `Assert "Hello, world" (actual, left) equals "Hello, world"`, col.Logs.Trace[12])
					assert.Equal(t, `Assertion OK. "$.data" == "Hello, world"`, col.Logs.Trace[13])
				}
			}
		}
	}
}
