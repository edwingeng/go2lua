-- package: example

local init = function()
    io.write("Hello ")
    print("World!")
end

Add = function(n1, n2)
    return n1 + n2
end

Sub = function(n1, n2)
    local n3 = n1 - n2
    return n3
end

Fibs = function(n)
    if n == 1 or n == 2 then
        return 1
    end

    return Fibs(n - 1) + Fibs(n - 2)
end

return init
