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

local operator4 = function(str1, str2, b1, b2, r1, r2)
    local _ = "x" .. "y"
    local _ = str1 .. str2
    local _ = "x" .. str1 .. "y"
    local _ = str1 .. "x" .. str2
end
