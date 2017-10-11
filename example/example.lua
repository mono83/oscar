-- Loading Oscar module
local o = require("oscar")

o.add("Base64 encoding transformations", function(tc)
    tc:set("value", "Hello, world")
    tc:assertEquals("SGVsbG8sIHdvcmxk", tc:stringToBase64("${value}"))

    tc:assertEquals("AAAAAAAAAAEAAAAAAAAAAv//////////", tc:packInt64ToBase64(1, 2, -1))
end)

-- Registering test cases
o.add("Simple POST with JSON", function(tc)
    tc:log("Using ${lua.engine}")

    tc:info("Sending request to https://httpbin.org")
    tc:httpPost("https://httpbin.org/post", 'Hello, world')

    tc:info("Checking received response")
    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.data", "Hello, world")
end)

