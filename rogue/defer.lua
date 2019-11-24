local defer_run1 = function (__body)
    local __defered = {}
    local r = table.pack(xpcall(__body, debug.traceback, __defered))
    __defered.f(table.unpack(__defered.args))
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

local defer_runN = function (__body)
    local __defered = {}
    local r = table.pack(xpcall(__body, debug.traceback, __defered))
    for i = #__defered, 1, -1 do
        local x = __defered[i]
        x.f(table.unpack(x.args))
    end
    if not r[1] then
        print(r[2])
        return
    end
    return table.unpack(r, 2)
end

return {
    run1 = defer_run1,
    runN = defer_runN,
}
