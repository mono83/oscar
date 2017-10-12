-- Loading Oscar module
local o = require("oscar")

o.add("Base64 encoding transformations", function(tc)
    tc:set("value", "Hello, world")
    tc:assertEquals("SGVsbG8sIHdvcmxk", tc:stringToBase64("${value}"))

    tc:assertEquals("AAAAAAAAAAEAAAAAAAAAAv//////////", tc:packInt64ToBase64(1, 2, -1))
end)
