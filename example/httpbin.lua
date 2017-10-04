-- Loading Oscar module
local o = require("oscar")

-- Registering test cases
o.add("Simple POST with JSON", function (tc)
    tc:log("Using ${lua.engine}")

    tc:log("Sending request to https://httpbin.org")
    tc:httpPost("https://httpbin.org/post", 'Hello, world')

    tc:log("Checking received response")
    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.data", "Hello, world")
end)