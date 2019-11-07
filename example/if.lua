-- package: example

If1 = function(n)
    if n == 0 then
        print("a1:", n)
    end

    if n == 1 then
        print("b1:", n)
    else
        print("b2:", n)
    end

    if n == 1 then
        print("c1:", n)
    elseif n == 2 then
        print("c2:", n)
    else
        print("c3:", n)
    end

    if n > 10 then
        if n > 100 then
            print("d1:", n)
        else
            print("d2:", n)
        end
    else
        if n < 1 then
            print("d3:", n)
        else
            print("d4:", n)
        end
    end

    if n > 10 then
        print("e1:", n)
    else
        if n == 1 then
            print("e2:", n)
        elseif n == 2 then
            print("e3:", n)
        elseif n == 3 then
            print("e4:", n)
        else
            print("e5:", n)
        end
    end
end

If2 = function(n)
    do
        local x1 = n * 10
        if x1 > 10 then
            print(x1)
        end
    end

    local x2 = 0
    do
        x2 = n * 10
        if x2 > 10 then
            print(x2)
        end
    end
end
