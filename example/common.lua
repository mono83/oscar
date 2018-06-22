-- Loading Oscar module
local o = require("oscar")

o.add("Sleep", function(tc)
    tc:sleep(100)
end)

o.add("Common logging", function(tc)
    tc:log("This is debug message")
    tc:info("This is info message")
    tc:set("foo", "bar")
    tc:log("Interpolating ${foo}")
end)

o.add("Common environment variables", function(tc)
    tc:log("This is ${some}")
    tc:assertEquals("${some}", "xxx-foo-xxx")
end)

o.add("Common library wrappers", function(tc)
    local w = Wrapper:create(tc)
    w:doAssertPositive(10)
end)