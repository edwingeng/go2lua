local empty = function()
    error("go empty struct")
end

local emptyStruct = setmetatable({}, {
    __len = empty,
    __index = empty,
    __newindex = empty,
    __metatable = false
})

return {
    empty = emptyStruct,
}
