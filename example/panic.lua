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
        __defered.args = {}
        __defered.f = function ()
            print(b1)
        end
        if b1 then
            error("world")
        end
        return 1, 2, 3
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    __defered.f(table.unpack(__defered.args))
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end
