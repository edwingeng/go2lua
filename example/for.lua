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

forLoop1 = function()
    while true do
        break
    end
end

forLoop2 = function()
    do
        local i = 0
        while i < 3 do
            i = i + 1
        end
    end
end

forLoop3 = function()
    local i = 0
    while i < 1 do
        break
    end
end

forLoop4 = function()
    do
        local i = 0
        while true do
            print(i)
            break
        end
    end
end

forLoop5 = function()
    local i = 0
    while true do
        if i >= 1 then
            break
        end
        i = i + 1
    end
end

forLoop6 = function()
    do
        local i = 0
        while i < 3 do
            if i >= 1 then
                goto __continue
            end
            print(i)
        ::__continue::
            i = i + 1
        end
    end

    do
        local i = 0
        while i < 3 do
            if i >= 1 then
                goto __continue_x2
            end
            print(i)
        ::__continue_x2::
            i = i + 1
        end
    end
end

forLoop7 = function()
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

forLoop8 = function()
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

forLoop9 = function(n)
    if n > 0 then
        goto pos1
    end
    print(100)

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

forLoop10 = function()
    do
        local i = 0
        while i < 3 do
            do
                local i = i * 10
                print(i)
            end
            i = i + 1
        end
    end
end

forLoop11 = function()
    do
        local i = 0
        while i < 4 do
            repeat
                local __switch = i
                if __switch == 0 then
                    goto outer_continue
                elseif __switch == 1 then
                elseif __switch == 2 then
                    print("a", i)
                else
                    -- default
                    goto outer_break
                end
            until true
            print("b", i)
        ::outer_continue::
            i = i + 1
        end
    end
::outer_break::
end

return function() end
