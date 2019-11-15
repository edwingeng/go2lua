-- package: example

local operator1 = function(n, ok)
    local a = n
    local b = -a
    local _ = not ok
    local d = ~b
    local e = {}
    local f = e
    local _, _, _ = d, e, f
end
