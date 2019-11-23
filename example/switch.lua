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

local switch1 = function(n)
    repeat
        local __switch = n
    until true
end

local switch2 = function(n)
    repeat
        local __switch = n
        if __switch == 1 then
            print("a", n)
            break
        elseif __switch == 2 then
            print("b", n)
        else
            -- default
            print("c", n)
        end
    until true
end

local switch3 = function(n)
    local a = 3
    local b = 2
    repeat
        local __switch = n
        if __switch == 1 or __switch == 3 then
            print("a", n)

        elseif __switch == a + b then
            print("b", n)
        else
            -- default
            print("c", n)
        end
    until true
end

local switch4 = function(n)
    repeat
        local __switch = n
        local __fall = false
        if __switch == 1 then
            print("a", n)
            __fall = true
            goto __case_2
        end

    ::__case_2::
        if  __fall or __switch == 2 then
            __fall = false
            print("b", n)
            __fall = true
            goto __case_3
        end

    ::__case_3::
        if  __fall or __switch == 3 then
            __fall = false
            print("c", n)
            break
        end

        if __switch == 4 then
            print("d", n)
            goto __switch_break
        end

        if __switch == 5 then
            print("e", n)
            __fall = true
            goto __case_6
        end

    ::__case_7::
        if  __fall or __switch == 6 or __switch == 7 then
            __fall = false
            print("g", n)
            goto __switch_break
        end

    ::__case_6::
        do
            -- default
            __fall = false
            print("f", n)
            __fall = true
            goto __case_7
        end
    until true
::__switch_break::
end

local switch5 = function(n)
    repeat
        local __fall = false
        if n == 1 then
            print("a", n)
            __fall = true
            goto __case_2
        end

    ::__case_2::
        if  __fall or n == 2 then
            __fall = false
            print("b", n)
            __fall = true
            goto __case_3
        end

    ::__case_3::
        if  __fall or n == 3 then
            __fall = false
            print("c", n)
            break
        end

        if n == 4 then
            print("d", n)
            goto __switch_break
        end

        if n == 5 then
            print("e", n)
            __fall = true
            goto __case_6
        end

    ::__case_7::
        if  __fall or n == 6 or n == 7 then
            __fall = false
            print("g", n)
            goto __switch_break
        end

    ::__case_6::
        do
            -- default
            __fall = false
            print("f", n)
            __fall = true
            goto __case_7
        end
    until true
::__switch_break::
end

return function() end
