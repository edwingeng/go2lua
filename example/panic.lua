-- package: example

local panic1 = function(b1, b2)
    if b1 then
        error("hello")
    end
    if b2 then
        error(100)
    end
end
