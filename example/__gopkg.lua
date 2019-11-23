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

local init_20af73cf = require("/defer")
init_20af73cf()
local init_dd73c24f = require("/for")
init_dd73c24f()
local init_5fb63259 = require("/func")
init_5fb63259()
local init_fb0077f9 = require("/hello")
init_fb0077f9()
local init_ad493904 = require("/if")
init_ad493904()
local init_2abb52c7 = require("/operator")
init_2abb52c7()
local init_06c2b30c = require("/panic")
init_06c2b30c()
local init_04c8fc80 = require("/range")
init_04c8fc80()
local init_588a9963 = require("/slice")
init_588a9963()
local init_9146c7e3 = require("/string")
init_9146c7e3()
local init_56f38f74 = require("/switch")
init_56f38f74()
local init_0cb6e51f = require("/var")
init_0cb6e51f()

return gopkg
