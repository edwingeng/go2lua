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

local hashes = {
	init_20af73cf = true,
	init_dd73c24f = true,
	init_5fb63259 = true,
	init_fb0077f9 = true,
	init_ad493904 = true,
	init_2abb52c7 = true,
	init_06c2b30c = true,
	init_04c8fc80 = true,
	init_588a9963 = true,
	init_9146c7e3 = true,
	init_56f38f74 = true,
	init_0cb6e51f = true,
}
local hashCount = 0
for _ in pairs(hashes) do
	hashCount = hashCount + 1
end
if hashCount ~= 12 then
	error("hash collision detected")
end

local init_20af73cf = require("defer")
if type(init_20af73cf) == "function" then
    init_20af73cf()
end
local init_dd73c24f = require("for")
if type(init_dd73c24f) == "function" then
    init_dd73c24f()
end
local init_5fb63259 = require("func")
if type(init_5fb63259) == "function" then
    init_5fb63259()
end
local init_fb0077f9 = require("hello")
if type(init_fb0077f9) == "function" then
    init_fb0077f9()
end
local init_ad493904 = require("if")
if type(init_ad493904) == "function" then
    init_ad493904()
end
local init_2abb52c7 = require("operator")
if type(init_2abb52c7) == "function" then
    init_2abb52c7()
end
local init_06c2b30c = require("panic")
if type(init_06c2b30c) == "function" then
    init_06c2b30c()
end
local init_04c8fc80 = require("range")
if type(init_04c8fc80) == "function" then
    init_04c8fc80()
end
local init_588a9963 = require("slice")
if type(init_588a9963) == "function" then
    init_588a9963()
end
local init_9146c7e3 = require("string")
if type(init_9146c7e3) == "function" then
    init_9146c7e3()
end
local init_56f38f74 = require("switch")
if type(init_56f38f74) == "function" then
    init_56f38f74()
end
local init_0cb6e51f = require("var")
if type(init_0cb6e51f) == "function" then
    init_0cb6e51f()
end

return gopkg
