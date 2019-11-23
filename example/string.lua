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

local string1 = function(b, r)
    local str = "hello"
    local _ = 97
    local _ = 20320
    local _ = string.byte(str, 1)
    local _ = utf8.char(97)
    local _ = utf8.char(b)
    local _ = utf8.char(r)
    local _ = string.len(str)
end

local string2 = function(str)
    do
        local i = 0
        while i < string.len(str) do
            print(i, string.byte(str, i + 1))
            i = i + 1
        end
    end

    for i, r in utf8.codes(str) do
        print(i, r)
    end
end
