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
            print(1)
        elseif switchTag == 2 then 
            print(2)
        else
            print(n)
        end
    end
end

local switch3 = function(n)
    local a = 3
    local b = 2
    do
        local switchTag = n
        if switchTag == 1 or switchTag == 3 then 
            print(1)
        elseif switchTag == (a + b) then 
            print(2)
        else
            print(n)
        end
    end
end
