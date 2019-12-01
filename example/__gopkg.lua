-- package: example

require("undef")
goslice = require("goslice")
godefer = require("godefer")
gostruct = require("gostruct")

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

local checkModName = function (name)
	if package.loaded[name] ~= nil then
		error("mod name collision detected. name: " .. name)
	end
end

checkModName("defer")
checkModName("for")
checkModName("func")
checkModName("hello")
checkModName("if")
checkModName("operator")
checkModName("order")
checkModName("panic")
checkModName("range")
checkModName("slice")
checkModName("struct")
checkModName("switch")
checkModName("var")
checkModName("xstring")

local hashes = {
	init_20af73cf = true,
	init_dd73c24f = true,
	init_5fb63259 = true,
	init_fb0077f9 = true,
	init_ad493904 = true,
	init_2abb52c7 = true,
	init_ed9fed76 = true,
	init_06c2b30c = true,
	init_04c8fc80 = true,
	init_588a9963 = true,
	init_6fe654c1 = true,
	init_56f38f74 = true,
	init_0cb6e51f = true,
	init_86e7cc0b = true,
}
local hashCount = 0
for _ in pairs(hashes) do
	hashCount = hashCount + 1
end
if hashCount ~= 14 then
	error("hash collision detected")
end

local init_20af73cf = require("defer")
local init_dd73c24f = require("for")
local init_5fb63259 = require("func")
local init_fb0077f9 = require("hello")
local init_ad493904 = require("if")
local init_2abb52c7 = require("operator")
local init_ed9fed76 = require("order")
local init_06c2b30c = require("panic")
local init_04c8fc80 = require("range")
local init_588a9963 = require("slice")
local init_6fe654c1 = require("struct")
local init_56f38f74 = require("switch")
local init_0cb6e51f = require("var")
local init_86e7cc0b = require("xstring")

-- Initializers
order3 = 100
order2 = order3 - 10
order1 = order2 - 10
order4 = order1 - 10
order5 = order6()
_ = order4
_ = order5
MyNumber2 = 200
MyNumber5 = 500
MyNumber6 = 600

init_20af73cf() -- defer
init_dd73c24f() -- for
init_5fb63259() -- func
init_fb0077f9() -- hello
init_ad493904() -- if
init_2abb52c7() -- operator
init_ed9fed76() -- order
init_06c2b30c() -- panic
init_04c8fc80() -- range
init_588a9963() -- slice
init_6fe654c1() -- struct
init_56f38f74() -- switch
init_0cb6e51f() -- var
init_86e7cc0b() -- xstring

return gopkg
