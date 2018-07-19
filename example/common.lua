-- Loading Oscar module
local o = require("oscar")
local w

o.add("Sleep", function(tc)
    tc:sleep(100)
end, {impact="none"})

o.add("Common logging", function(tc)
    tc:log("This is debug message")
    tc:info("This is info message")
    tc:set("foo", "bar")
    tc:log("Interpolating ${foo}")
end, {impact="none"})

o.add("Common environment variables", function(tc)
    tc:log("This is ${some}")
    tc:assertEquals("${some}", "xxx-foo-xxx")
end, {impact="none"})

o.add("Common library wrappers", function(tc)
    w:doAssertPositive(10)
end, {impact="none"})

o.add("Init variable read", function(tc)
    tc:assertEquals("${initvar}", "zzz")
end, {impact="none"})

o.setUp(function(tc)
    tc:info("SetUp function")
    w = Wrapper:create(tc)
    tc:set("initvar", "zzz")
end)

o.add("SHA256 test", function(tc)
    tc:assertEquals(tc:sha256Hex("test string"), "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b")
end, {impact="none"})