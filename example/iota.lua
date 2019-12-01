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

iotaNumA1 = 0
iotaNumA2 = iotaNumA1 + 1
iotaNumA3 = iotaNumA2 + 1
iotaNumB1 = 10
iotaNumB2 = iotaNumB1
iotaNumC1 = iotaNumB1 + 10 + 5
iotaNumC2 = iotaNumC1 + 1

iotaNumC3 = iotaNumC2 + 1

iotaNumX1 = 10
iotaNumD1 = 1
iotaNumD2 = iotaNumD1 + 1
iotaNumD3 = iotaNumD2 + 1
iotaNumD4, iotaNumD5 = 10, 20

iota1 = function()
    local iotaNumA1 = 0
    local iotaNumA2 = iotaNumA1 + 1
    local iotaNumA3 = iotaNumA2 + 1
    local iotaNumB1 = 10
    local iotaNumB2 = iotaNumB1
    local iotaNumC1 = iotaNumB1 + 10 + 5
    local iotaNumC2 = iotaNumC1 + 1
end

local init = function()
    if iotaNumA1 ~= 0 then
        error("iotaNumA1 != 0")
    end
    if iotaNumA2 ~= 1 then
        error("iotaNumA2 != 1")
    end
    if iotaNumA3 ~= 2 then
        error("iotaNumA3 != 2")
    end
    if iotaNumB1 ~= 10 then
        error("iotaNumB1 != 10")
    end
    if iotaNumB2 ~= 10 then
        error("iotaNumB2 != 10")
    end
    if iotaNumC1 ~= 25 then
        error("iotaNumC1 != 25")
    end
    if iotaNumC2 ~= 26 then
        error("iotaNumC2 != 26")
    end
    if iotaNumC3 ~= 27 then
        error("iotaNumC3 != 27")
    end
    if iotaNumD1 ~= 1 then
        error("iotaNumD1 != 1")
    end
    if iotaNumD2 ~= 2 then
        error("iotaNumD2 != 2")
    end
    if iotaNumD3 ~= 3 then
        error("iotaNumD3 != 3")
    end
    if iotaNumD4 ~= 10 then
        error("iotaNumD4 != 10")
    end
    if iotaNumD5 ~= 20 then
        error("iotaNumD5 != 20")
    end
end

return init
