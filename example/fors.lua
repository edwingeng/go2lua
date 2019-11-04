-- package: example

local forLoop1 = function()
    while true do
        break
    end
end

local forLoop2 = function()
    do
        local i = 0
        while i < 1 do
            i = i + 1
        end
    end
end

local forLoop3 = function()
    local i = 0
    while i < 1 do
        break
    end
end

local forLoop4 = function()
    do
        local i = 0
        while true do
            print(i)
            break
        end
    end
end

local forLoop5 = function()
    local i = 0
    while true do
        if i >= 1 then
            break
        end
        i = i + 1
    end
end

local forLoop6 = function()
    do
        local i = 0
        while i < 3 do
            if i >= 1 then
                goto xxx_loop_1
            end
            print(i)
        ::xxx_loop_1::
            i = i + 1
        end
    end
end

local forRangeMap = function()
    local m = {}
    m["a"] = 100
    m["b"] = 200
    for k, v in pairs(m) do
        print(k, v)
    end
end
