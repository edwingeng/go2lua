undef = setmetatable({}, {
    __index = function()
        error("undefined")
    end,
    __newindex = function()
        error("undefined")
    end,
    __metatable = false
})
