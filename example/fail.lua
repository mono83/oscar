-- Loading Oscar module
local o = require("oscar")

o.add("Failing case", function(tc)
    tc:assertEquals("1", "1")
    tc:assertEquals("1", "2")
end)