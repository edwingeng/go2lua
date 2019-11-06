-- package: example

local switch1 = function(n)
    repeat
        local switchTag = n
    until true
end

local switch2 = function(n)
    repeat
        local switchTag = n
        if switchTag == 1 then
            print("a", n)
            break
        elseif switchTag == 2 then
            print("b", n)
        else
            print("c", n)
        end
    until true
end

local switch3 = function(n)
    local a = 3
    local b = 2
    repeat
        local switchTag = n
        if switchTag == 1 or switchTag == 3 then
            print("a", n)

        elseif switchTag == (a + b) then
            print("b", n)
        else
            print("c", n)
        end
    until true
end

local switch4 = function(n)
    repeat
        local switchTag = n
        local __fall = false
        if switchTag == 1 then
            print("a", n)
            __fall = true
            goto __switch_2
        end
        
    ::__switch_2::
        if  __fall or switchTag == 2 then
            __fall = false
            print("b", n)
            __fall = true
            goto __switch_3
        end
        
    ::__switch_3::
        if  __fall or switchTag == 3 then
            __fall = false
            print("c", n)
            break
        end
        
        if switchTag == 4 then
            print("d", n)
            goto __switch_break
        end
        
        if switchTag == 5 then
            print("e", n)
            __fall = true
            goto __switch_6
        end
        
    ::__switch_7::
        if  __fall or switchTag == 6 or switchTag == 7 then
            __fall = false
            print("g", n)
            goto __switch_break
        end
        
    ::__switch_6::
        do
            __fall = false
            print("f", n)
            __fall = true
            goto __switch_7
        end
    until true
::__switch_break::
end
