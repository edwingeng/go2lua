-- package: example

Ifs = function(n)
    if n == 0 then
        print("a1:", 0)
    end
    if n == 1 then
        print("b1:", 1)
    else
        print("b2:", n)
    end
    if n == 1 then
        print("c1:", 1)
    elseif n == 2 then
        print("c2:", 2)
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
            print("e2:", 1)
        elseif n == 2 then
            print("e3:", 2)
        elseif n == 3 then
            print("e4:", 3)
        else
            print("e5:", n)
        end
    end
end
