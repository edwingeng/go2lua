-- package: example

local gopkg = {}
_G["github.com/edwingeng/go2lua/example"] = gopkg
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


local init_001 = require("defer")
if type(init_001) == "function" then
    init_001()
end
local init_002 = require("for")
if type(init_002) == "function" then
    init_002()
end
local init_003 = require("func")
if type(init_003) == "function" then
    init_003()
end
local init_004 = require("hello")
if type(init_004) == "function" then
    init_004()
end
local init_005 = require("if")
if type(init_005) == "function" then
    init_005()
end
local init_006 = require("operator")
if type(init_006) == "function" then
    init_006()
end
local init_007 = require("panic")
if type(init_007) == "function" then
    init_007()
end
local init_008 = require("range")
if type(init_008) == "function" then
    init_008()
end
local init_009 = require("slice")
if type(init_009) == "function" then
    init_009()
end
local init_010 = require("string")
if type(init_010) == "function" then
    init_010()
end
local init_011 = require("switch")
if type(init_011) == "function" then
    init_011()
end
local init_012 = require("var")
if type(init_012) == "function" then
    init_012()
end

return gopkg
