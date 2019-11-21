-- package: example

local panic1 = function(b1, b2)
    if b1 then
        error("hello")
    end
    if b2 then
        error(100)
    end
end

local panic2 = function(b1)
    local __defered = {}
    local __body = function ()
        local __funcObj = {args = {}}
        table.insert(__defered, __funcObj)
        __funcObj.f = function ()
            print(b1)
        end
        if b1 then
            error("world")
        end
        return 1, 2, 3
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    for i = #__defered, 1, -1 do
        local x = __defered[i]
        x.f(table.unpack(x.args))
    end
    if r[1] then
        table.remove(r, 1)
        return table.unpack(r)
    else
        print(r[2])
    end
end
