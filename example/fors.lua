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
                goto xxx_continue_1
            end
            print(i)
        ::xxx_continue_1::
            i = i + 1
        end
    end
end

local forLoop7 = function()
::pos1::
    do
        local i = 0
        while i < 3 do
            while true do
                goto pos1_break
            end
            i = i + 1
        end
    end
::pos1_break::
end

local forLoop8 = function()
::pos1::
    do
        local i = 0
        while i < 3 do
            while true do
                goto pos1_continue
            end
        ::pos1_continue::
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
