-- Loading Oscar module
local o = require("oscar")

o.add("Simple GET with JSON", function(tc)
    tc:info("Sending request to https://httpbin.org")
    tc:httpGet("https://httpbin.org/get")

    tc:info("Checking received response")
    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.url", "https://httpbin.org/get")
end, {impact="read"})


o.add("Simple POST with JSON", function(tc)
    tc:info("Sending request to https://httpbin.org")
    tc:httpPost("https://httpbin.org/post", 'Hello, world')

    tc:info("Checking received response")
    tc:assertEquals("${http.response.code}", "200")
    tc:assertEquals("${http.response.header.Access-Control-Allow-Credentials}", "true")
    tc:assertJSONXPath("$.data", "Hello, world")
end, {impact="read", depends="Simple GET with JSON"})
