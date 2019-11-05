-- package: example

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
            goto xxx_continue_1
        else
            print(k, v)
        end
    ::xxx_continue_1::
    end
end

local rangeMap3 = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200
::pos1::
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
::pos1::
    for k, v in pairs(m) do
        print(k, v)
        while true do
            goto pos1_continue
        end
    ::pos1_continue::
    end
end
