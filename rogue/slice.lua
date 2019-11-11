local slice_mt = {
    __len = function(s)
        return s.len
    end,
    __index = function(s, i)
        if i <= s.len then
            return s.data[i + s.off]
        end
        error(string.format("index out of range [%d] with length %d", i, s.len))
    end,
    __newindex = function(s, i, v)
        local n = s.len
        if i <= n then
            s.data[i + s.off] = v
        elseif i == n + 1 then
            s.len = i
            s.data[i + s.off] = v
        else
            error(string.format("unexpected newindex [%d] with length %d", i, s.len))
        end
    end,
    __pairs = function(s)
        local function iter(d, k) return next(d, k) end
        return iter, s.data, nil
    end,
    __ipairs = function(s)
        local function iter(d, i)
            local j = i + 1
            local v = d[j]
            if v then
                return j, v
            end
        end
        return iter, s.data, 0
    end,
    __metatable = false
}

if show_slice_metatable then
    slice_mt.__metatable = nil
end

local slice_make = function(fnInit, len)
    len = len or 0
    if len > 0 then
        local a = fnInit(len)
        if a == nil then
            error("nil slice data")
        end
        local s = {data = a, len = len, off = 0}
        return setmetatable(s, slice_mt)
    else
        return setmetatable({data = {}, len = 0, off = 0}, slice_mt)
    end
end

local function slice_fromArray(a, deep, depth)
    if depth > 99 then
        error("depth > 99")
    end

    if deep then
        local n = #a
        for i = 1, n do
            if type(a[i]) == "table" and getmetatable(a[i]) ~= slice_mt then
                a[i] = slice_fromArray(a[i], true, depth + 1)
            end
        end
    end

    local s = {data = a, len = #a, off = 0}
    return setmetatable(s, slice_mt)
end

local function slice_toArray(s, deep, depth)
    if depth > 99 then
        error("depth > 99")
    end

    local a
    local n = s.len
    if s.off == 0 and n == #s.data then
        a = s.data
    else
        a = {}
        local j = s.off
        for i = 1, n do
            j = j + 1
            a[i] = s.data[j]
        end
    end

    if deep then
        for i = 1, n do
            if type(a[i]) == "table" and getmetatable(a[i]) == slice_mt then
                a[i] = slice_toArray(a[i], true, depth + 1)
            end
        end
    end
    return a
end

local slice_append = function(s, v)
    if s ~= undef then
        local x = {data = s.data, len = s.len + 1, off = s.off}
        x.data[x.len + x.off] = v
        return setmetatable(x, slice_mt)
    else
        local a = {v}
        local x = {data = a, len = 1, off = 0}
        return setmetatable(x, slice_mt)
    end
end

local slice_appendSlice = function(s, appendix)
    if appendix == undef then
        if s ~= undef then
            local x = {data = s.data, len = s.len, off = s.off}
            return setmetatable(x, slice_mt)
        else
            return setmetatable({data = {}, len = 0, off = 0}, slice_mt)
        end
    else
        if s ~= undef then
            local n = appendix.len
            local x = {data = s.data, len = s.len + n, off = s.off}
            local j = s.len + s.off
            local w = n + appendix.off
            for i = 1 + appendix.off, w do
                j = j + 1
                x.data[j] = appendix.data[i]
            end
            return setmetatable(x, slice_mt)
        else
            local a = {}
            local n = appendix.len
            local j = appendix.off
            for i = 1, n do
                j = j + 1
                a[i] = appendix.data[j]
            end
            local x = {data = a, len = n, off = 0}
            return setmetatable(x, slice_mt)
        end
    end
end

local slice_slice = function(s, start, eNd)
    if s == undef then
        return setmetatable({data = {}, len = 0, off = 0}, slice_mt)
    end

    local n = s.len
    local beyond = n + 1
    start = start or 1
    eNd = eNd or beyond
    if start > beyond then
        error(string.format("'start' out of range [%d] with length %d, beyond %d", start, n, beyond))
    end
    if eNd > beyond then
        error(string.format("'eNd' out of range [%d] with length %d, beyond %d", eNd, n, beyond))
    end
    if eNd < start then
        error(string.format("invalid 'eNd': %d < %d", eNd, start))
    end

    local x = {data = s.data, len = eNd - start, off = s.off + start - 1}
    return setmetatable(x, slice_mt)
end

local slice_copy = function(dst, src)
    if src == undef or src.len == 0 then
        return 0
    else
        if dst ~= undef then
            local dstLen, srcLen = dst.len, src.len
            local n = srcLen
            if dstLen < srcLen then
                n = dstLen
            end
            local off1, off2 = dst.off, src.off
            if dst.data ~= src.data or dst.off <= src.off then
                for i = 1, n do
                    dst.data[i + off1] = src.data[i + off2]
                end
            else
                for i = n, 1, -1 do
                    dst.data[i + off1] = src.data[i + off2]
                end
            end
            return n
        else
            return 0
        end
    end
end

local slice_clone = function(s, start, eNd)
    if s == undef then
        return undef
    end

    local n = s.len
    local beyond = n + 1
    start = start or 1
    eNd = eNd or beyond
    if start > beyond then
        error(string.format("'start' out of range [%d] with length %d, beyond %d", start, n, beyond))
    end
    if eNd > beyond then
        error(string.format("'eNd' out of range [%d] with length %d, beyond %d", eNd, n, beyond))
    end
    if eNd < start then
        error(string.format("invalid 'eNd': %d < %d", eNd, start))
    end

    local a = {}
    local j = eNd - start
    local off = s.off
    for i = eNd - 1, start, -1 do
        a[j] = s.data[i + off]
        j = j - 1
    end
    return setmetatable({data = a, len = eNd - start, off = 0}, slice_mt)
end

return {
    mt = slice_mt,
    make = slice_make,
    fromArray = function(s, deep) return slice_fromArray(s, deep, 1) end,
    toArray = function(s, deep) return slice_toArray(s, deep, 1) end,
    append = slice_append,
    appendSlice = slice_appendSlice,
    slice = slice_slice,
    copy = slice_copy,
    clone = slice_clone,
}
