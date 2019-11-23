-- package: example

local gopkg = _G["github.com/edwingeng/go2lua/example"]
do
    local g = _G
    local newEnv = setmetatable({}, {
        __index = function (t, k)
            local v = gopkg[k]
            if v == nil then return g[k] end
            return v
        end,
        __newindex = gopkg,
    })
    _ENV = newEnv
end

local defer1 = function()
    local __defered = {}
    local __body = function ()
        __defered.args = {}
        __defered.f = function ()
            print(100)
        end
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    __defered.f(table.unpack(__defered.args))
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer2 = function()
    local __defered = {}
    local __body = function ()
        local __funcObj = {args = {}}
        table.insert(__defered, __funcObj)
        __funcObj.f = function ()
            print(100)
        end

        print(300)

        local __funcObj_x2 = {args = {}}
        table.insert(__defered, __funcObj_x2)
        __funcObj_x2.f = function ()
            print(200)
        end
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    for i = #__defered, 1, -1 do
        local x = __defered[i]
        x.f(table.unpack(x.args))
    end
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer3 = function()
    local __defered = {}
    local __body = function ()
        local f1 = function ()
            print(100)
        end
        local f2 = function ()
            print(200)
        end

        local __funcObj = {args = {}}
        table.insert(__defered, __funcObj)
        __funcObj.f = f1
        local __funcObj_x2 = {args = {}}
        table.insert(__defered, __funcObj_x2)
        __funcObj_x2.f = f2
        print(300)
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    for i = #__defered, 1, -1 do
        local x = __defered[i]
        x.f(table.unpack(x.args))
    end
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer4 = function()
    local __defered = {}
    local __body = function ()
        __defered.args = {}
        __defered.f = function ()
            local __defered = {}
            local __body = function ()
                __defered.args = {}
                __defered.f = function ()
                    local __defered = {}
                    local __body = function ()
                        __defered.args = {}
                        __defered.f = function ()
                            print(100)
                        end
                        print(200)
                    end

                    local r = table.pack(xpcall(__body, debug.traceback))
                    __defered.f(table.unpack(__defered.args))
                    if not r[1] then
                        print(r[2])
                        return
                    end
                    return table.unpack(r, 2)
                end
                print(300)
            end

            local r = table.pack(xpcall(__body, debug.traceback))
            __defered.f(table.unpack(__defered.args))
            if not r[1] then
                print(r[2])
                return
            end
            return table.unpack(r, 2)
        end
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    __defered.f(table.unpack(__defered.args))
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer5 = function()
    local __defered = {}
    local __body = function ()
        do
            local i = 0
            while i < 3 do
                local __funcObj = {args = {}}
                table.insert(__defered, __funcObj)
                __funcObj.f = function ()
                    print(i)
                end
                i = i + 1
            end
        end

        do
            local i = 0
            while i < 3 do
                local __funcObj_x2 = {args = {i}}
                table.insert(__defered, __funcObj_x2)
                __funcObj_x2.f = function (i)
                    print(i)
                end
                i = i + 1
            end
        end
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    for i = #__defered, 1, -1 do
        local x = __defered[i]
        x.f(table.unpack(x.args))
    end
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer6 = function(n1, n2)
    local __defered = {}
    local __body = function ()
        __defered.args = {}
        __defered.f = function ()
            print(n1, n2)
        end
    end

    local r = table.pack(xpcall(__body, debug.traceback))
    __defered.f(table.unpack(__defered.args))
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

return function() end
