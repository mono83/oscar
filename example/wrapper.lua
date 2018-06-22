Wrapper = {}
Wrapper.__index = Wrapper

function Wrapper:create(context)
    local foo = {}
    setmetatable(foo, Wrapper)
    foo.tc = context
    return foo
end

function Wrapper:doAssertPositive(num)
    self.tc:debug("Making assert for positive value")
    if num > 0 then
        self.tc:assertEquals("true", "true")
    else
        self.tc:assertEquals("false", "true")
    end
end