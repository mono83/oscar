-- Loading Oscar module
local o = require("oscar")

o.add("Failing assertions", function(tc)
    tc:assertEquals("1", "1")
    tc:assertEquals("1", "2")
end, {impact="none"})

o.add("Simple failing with no data", function(tc)
    tc:fail()
end, {impact="none"})

o.add("Simple failing with message", function(tc)
    tc:set("foo", "bar")
    tc:fail("Some example ${foo} message")
end, {impact="none"})

o.add("Simple failing with fmt", function(tc)
    tc:fail("Some example %s message", "baz")
end, {impact="none"})