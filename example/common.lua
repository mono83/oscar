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