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

defer1 = function()
    local __body = function (__defered)
        __defered.args = {}
        __defered.f = function ()
            print(100)
        end
    end

    return godefer.run1(__body)
end

defer2 = function()
    local __body = function (__defered)
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

    return godefer.runN(__body)
end

defer3 = function()
    local __body = function (__defered)
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

    return godefer.runN(__body)
end

defer4 = function()
    local __body = function (__defered)
        __defered.args = {}
        __defered.f = function ()
            local __body = function (__defered)
                __defered.args = {}
                __defered.f = function ()
                    local __body = function (__defered)
                        __defered.args = {}
                        __defered.f = function ()
                            print(100)
                        end
                        print(200)
                    end

                    return godefer.run1(__body)
                end
                print(300)
            end

            return godefer.run1(__body)
        end
    end

    return godefer.run1(__body)
end

defer5 = function()
    local __body = function (__defered)
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

    return godefer.runN(__body)
end

defer6 = function(n1, n2)
    local __body = function (__defered)
        __defered.args = {}
        __defered.f = function ()
            print(n1, n2)
        end
    end

    return godefer.run1(__body)
end

return function() end
