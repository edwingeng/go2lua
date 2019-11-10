slice_mt = function(data)
    return {
        __len = function(s)
            return s.len
        end,
        __index = function(s, i)
            if i >= s.len then
                error(string.format("index out of range [%d] with length %d", i, s.len))
            end
            return data[i]
        end,
        __newindex = function(s, i, v)
            if i < s.len then
                data[i] = v
            else
                data[s.len] = v
                s.len = s.len + 1
            end
        end,
        __metatable = false
    }
end

slice_make = function(len, fnInit)
    if len > 0 then
        local data = fnInit(len)
        local x = {data = data, len = len}
        return setmetatable(x, slice_mt(data))
    else
        local data = {}
        local x = {data = data, len = 0}
        return setmetatable(x, slice_mt(data))
    end
end
