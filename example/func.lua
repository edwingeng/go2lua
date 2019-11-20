-- package: example

local func1 = function()
    local f1 = function ()
        print("f1")
    end
    f1()

    local f2 = function ()
        print("f2")
    end
    f2()

    __lambda = function ()
        print("f3")
    end
    __lambda()

    __lambda_x2 = function (n1, n2, n3)
        print(n1, n2, n3)
    end
    __lambda_x2(1, 2, 3)
end

local func2 = function()
    __lambda = function (cb)
        if cb ~= nil then
            cb(100)
        end
    end
    __lambda(function (n)
        print(n)
    end)
end
