-- Loading Oscar module
local o = require("oscar")

o.add("Skipped due init failure", function(tc)
    tc:fail("Expected fail")
end)

o.setUp(function(tc)
    tc:fail("Expected fail")
end)