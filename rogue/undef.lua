local undefined = function()
    error("undefined")
end

undef = setmetatable({}, {
    __len = undefined,
    __index = undefined,
    __newindex = undefined,
    __metatable = false
})
