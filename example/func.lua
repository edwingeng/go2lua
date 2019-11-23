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

local func1 = function()
    local f1 = function ()
        print("f1")
    end
    f1()

    local f2 = function ()
        print("f2")
    end
    f2()

    local __lambda = function ()
        print("f3")
    end
    __lambda()

    local __lambda_x2 = function (n1, n2, n3)
        print(n1, n2, n3)
    end
    __lambda_x2(1, 2, 3)
end

local func2 = function()
    local __lambda = function (cb)
        if cb ~= nil then
            cb(100)
        end
    end
    __lambda(function (n)
        print(n)
    end)
end

return function() end
