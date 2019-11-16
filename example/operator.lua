-- package: example

local operator1 = function(n, ok)
    local a = n
    local b = -a
    local _ = not ok
    local d = ~b
    local e = {}
    local f = e
    local _, _, _ = d, e, f
end

local operator2 = function(n1, n2, b1, b2)
    local _ = n1 * n2
    local _ = n1 / n2
    local _ = n1 % n2
    local _ = n1 << n2
    local _ = n1 >> n2
    local _ = n1 & n2
    local _ = n1 & ~n2
    local _ = n1 + n2
    local _ = n1 - n2
    local _ = n1 | n2
    local _ = n1 ~ n2
    local _ = n1 == n2
    local _ = n1 ~= n2
    local _ = n1 < n2
    local _ = n1 <= n2
    local _ = n1 > n2
    local _ = n1 >= n2
    local _ = b1 and b2
    local _ = b1 or b2
end

local operator3 = function(n1, n2)
    local _ = n1 * n2 + n1 / n2
    local _ = (n1 * n2) + (n1 / n2)
    local _ = n1 * (n2 + n1) / n2
    local _ = (n1 * n2) + (n1 / n2) * (n1 + n2)
    local _ = ((n1 * n2) + (n1 / n2)) * (n1 + n2)
    local _ = n1 + n2 + n1 + n2 + n1
    local _ = n1 * n2 + n1 / n2 + n1 % n2 + n1 & n2
end
