-- package: example

Fibs = function(n)
    if n == 1 or n == 2 then
        return 1
    end
    return Fibs(n - 1) + Fibs(n - 2)
end

