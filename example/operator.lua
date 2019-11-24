-- package: example

local gopkg = _G["github.com/edwingeng/go2lua/example"]
do
    local g = _G
    local newEnv = setmetatable({}, {
        __index = function (t, k)
            local v = gopkg[k]
            if v == nil then return g[k] end
            return v
        end,
        __newindex = gopkg,
    })
    _ENV = newEnv
end

operator1 = function(n, ok)
    local a = n
    local b = -a
    local _ = not ok
    local d = ~b
    local e = {}
    local f = e
    local _, _, _ = d, e, f
end

operator2 = function(n1, n2, b1, b2)
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

operator3 = function(n1, n2)
    local _ = n1 * n2 + n1 / n2
    local _ = (n1 * n2) + (n1 / n2)
    local _ = n1 * (n2 + n1) / n2
    local _ = (n1 * n2) + (n1 / n2) * (n1 + n2)
    local _ = ((n1 * n2) + (n1 / n2)) * (n1 + n2)
    local _ = n1 + n2 + n1 + n2 + n1
    local _ = n1 * n2 + n1 / n2 + n1 % n2 + (n1 & n2)
end

operator4 = function(str1, str2, b1, b2, r1, r2)
    local _ = "x" .. "y"
    local _ = str1 .. str2
    local _ = "x" .. str1 .. "y"
    local _ = str1 .. "x" .. str2
    local _ = b1 + b2
    local _ = r1 + r2
end

operator5 = function(n1, n2)
    local _ = (n1 | n2) ~ n2
    local _ = (n1 | n2) + n1 | n2
    local _ = (n1 + n2 | n1) + n2
    local _ = (n1 | n2) ~ n1
    local _ = (n1 << n2) + ((n1 << n2) * n2 >> n1)
end

return function() end
