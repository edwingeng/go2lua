-- package: example

local switch1 = function(n)
    do
        local switchTag = n
    end
end

local switch2 = function(n)
    do
        local switchTag = n
        if switchTag == 1 then 
            print("a", n)
        elseif switchTag == 2 then 
            print("b", n)
        else
            print("c", n)
        end
    end
end

local switch3 = function(n)
    local a = 3
    local b = 2
    do
        local switchTag = n
        if switchTag == 1 or switchTag == 3 then 
            print("a", n)
        elseif switchTag == (a + b) then 
            print("b", n)
        else
            print("c", n)
        end
    end
end
