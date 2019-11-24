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

order1 = undef
order2 = undef
order3 = undef
order4 = undef
order5 = undef

order6 = function()
    return 500
end

_, _ = undef, undef

return function() end
