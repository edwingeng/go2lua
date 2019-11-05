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
