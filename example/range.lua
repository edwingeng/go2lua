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

local rangeMap1 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200

    for k, v in pairs(m) do
        print(k, v)
    end
end

local rangeMap2 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200
    m["c"] = 300

    for k, v in pairs(m) do
        if k == "b" then
            goto __continue
        else
            print(k, v)
        end
    ::__continue::
    end
end

local rangeMap3 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200

    for k, v in pairs(m) do
        print(k, v)
        while true do
            goto pos1_break
        end
    end
::pos1_break::
end

local rangeMap4 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200

    for k, v in pairs(m) do
        print(k, v)
        while true do
            goto pos1_continue
        end
    ::pos1_continue::
    end
end

local rangeMap5 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200

    for _ in pairs(m) do
    end
end

return function() end
