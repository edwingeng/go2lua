require "undef"
expose_slice_metatable = true
local slice = require("goslice")
local inspect = require("inspect")

local always = function(s, full, suffix)
    suffix = suffix or ""
    if s == nil then
        error("nil slice" .. suffix)
    end
    if s == undef then
        error("undefined slice" .. suffix)
    end
    if getmetatable(s) ~= slice.mt then
        error("getmetatable(s) ~= slice.mt" .. suffix)
    end
    local fields = {"data", "len", "off"}
    for _, f in ipairs(fields) do
        if s[f] == nil then
            error(string.format("no .%s. s: %s" .. suffix, f, inspect(s)))
        end
    end
    if full and s.len ~= #s.data then
        error(string.format("invalid .len: %d. actual: %d" .. suffix, s.len, #s.data))
    end
    if s.off + s.len > #s.data then
        error("s.off + s.len > #s.data" .. suffix)
    end
end

local equals = function(s1, s2, suffix)
    suffix = suffix or ""
    if type(s1) ~= type(s2) then
        error("type(s1) ~= type(s2)" .. suffix)
    end

    always(s1, false, suffix)
    always(s2, false, suffix)

    if #s1 ~= #s2 then
        error("#s1 ~= #s2" .. suffix)
    end
end

local test_slice_make = function()
    local init = function(n)
        local a = {}
        for i = 1, n do
            table.insert(a, 0)
        end
        return a
    end

    for _, v in ipairs({0, 1, 2, 3, 10}) do
        local suffix = string.format(", v: %d", v)
        local s = slice.make(init, v)
        always(s, true, suffix)
        if #s ~= v then
            error("#s ~= v" .. suffix)
        end
    end
end

local test_slice_fromArray = function()
    for i = 0, 3 do
        local suffix = string.format(", i: %d", i)
        local a = {}
        for j = 1, i do
            table.insert(a, i)
        end
        local s = slice.fromArray(a)
        always(s, true, suffix)
        if #s ~= i then
            error("#s ~= i" .. suffix)
        end
    end

    local s1a = slice.fromArray({1, 2, 3, 4, 5})
    s1a[3] = slice.fromArray({7, 8, 9})
    local arr = {1, 2, {7, 8, 9}, 4, 5}
    local s1b = slice.fromArray(arr, true)
    if inspect(s1a) ~= inspect(s1b) then
        error("inspect(s1a) ~= inspect(s1b)")
    end
end

local test_slice_toArray = function()
    for i = 0, 3 do
        local suffix = string.format(", i: %d", i)
        local a = {}
        for j = 1, i do
            table.insert(a, i)
        end
        local s = slice.fromArray(a)
        always(s, true, suffix)
        local b = slice.toArray(s)
        if b == nil then
            error("b == nil")
        end
        if #b ~= i then
            error("#b ~= i" .. suffix)
        end
        for j, v in ipairs(b) do
            if v ~= a[j] then
                error("v ~= a[j]" .. string.format(". j: %d", j) .. suffix)
            end
        end
    end

    local s1a = slice.fromArray({1, 2, 3, 4, 5})
    s1a[3] = slice.fromArray({7, 8, 9})
    local s1b = slice.toArray(s1a, true)
    local arr = {1, 2, {7, 8, 9}, 4, 5}
    if inspect(s1b) ~= inspect(arr) then
        error("inspect(s1b) ~= inspect(arr)")
    end

    local s2a = slice.slice(s1a, 3)
    local s2b = slice.toArray(s2a, true)
    arr = {{7, 8, 9}, 4, 5}
    if inspect(s2b) ~= inspect(arr) then
        error("inspect(s2b) ~= inspect(arr)")
    end
end

local test_slice_append = function()
    local s1a = slice.append(undef, undef)
    always(s1a, true)
    local s1b = slice.fromArray({undef})
    always(s1b, true)
    equals(s1a, s1b)

    local s2a = slice.append(undef, 100)
    always(s2a, true)
    local s2b = slice.fromArray({100})
    always(s2b, true)
    equals(s2a, s2b)

    s2a = slice.append(s2a, 200)
    always(s2a, true)
    s2b = slice.fromArray({100, 200})
    always(s2b, true)
    equals(s2a, s2b)

    s2a = slice.append(s2a, undef)
    always(s2a, true)
    s2b = slice.fromArray({100, 200, undef})
    always(s2b, true)
    equals(s2a, s2b)

    local s3a = slice.append(s2a, 500)
    always(s3a, true)
    local s3b = slice.fromArray({100, 200, undef, 500})
    always(s3b, true)
    equals(s3a, s3b)

    if #s2a ~= 3 then
        error("#s2a ~= 3")
    end
    if #s3a ~= 4 then
        error("#s3a ~= 4")
    end
end

local test_slice_appendArray = function()
    local s2a = slice.fromArray({1, 2, 3, 4, 5})
    always(s2a, true)
    local s2b = slice.fromArray({6, 7})
    always(s2b, true)
    local s2c = slice.appendArray(s2a, slice.toArray(s2b))
    always(s2c, true)
    local s2d = slice.fromArray({1, 2, 3, 4, 5, 6, 7})
    always(s2d, true)
    equals(s2c, s2d)

    if #s2a ~= 5 then
        error("#s2a ~= 5")
    end
    if #s2c ~= 7 then
        error("#s2c ~= 7")
    end

    local s3a = slice.appendArray(undef, slice.toArray(s2b))
    always(s3a, true)
    equals(s3a, s2b)

    local s5a = slice.appendArray(s2a, slice.toArray(s2a))
    always(s5a, true)
    local s5b = slice.fromArray({1, 2, 3, 4, 5, 1, 2, 3, 4, 5})
    always(s5b, true)
    equals(s5a, s5b)

    if #s2a ~= 5 then
        error("#s2a ~= 5")
    end

    local s6a = slice.slice(s2a, 2, 5)
    always(s6a, false)
    local s6b = slice.appendArray(s2a, slice.toArray(s6a))
    always(s6b, false)
    local s6c = slice.fromArray({1, 2, 3, 4, 5, 2, 3, 4})
    always(s6c, true)
    equals(s6b, s6c)

    local s7a = slice.appendArray(s6a, slice.toArray(s2a))
    always(s7a, false)
    local s7b = slice.fromArray({2, 3, 4, 1, 2, 3, 4, 5})
    always(s7b, true)
    equals(s7a, s7b)
end

local test_slice_appendSlice = function()
    local s1a = slice.appendSlice(undef, undef)
    always(s1a, true)
    local s1b = slice.fromArray({})
    always(s1b, true)
    equals(s1a, s1b)

    local s2a = slice.fromArray({1, 2, 3, 4, 5})
    always(s2a, true)
    local s2b = slice.fromArray({6, 7})
    always(s2b, true)
    local s2c = slice.appendSlice(s2a, s2b)
    always(s2c, true)
    local s2d = slice.fromArray({1, 2, 3, 4, 5, 6, 7})
    always(s2d, true)
    equals(s2c, s2d)

    if #s2a ~= 5 then
        error("#s2a ~= 5")
    end
    if #s2c ~= 7 then
        error("#s2c ~= 7")
    end

    local s3a = slice.appendSlice(undef, s2b)
    always(s3a, true)
    equals(s3a, s2b)

    local s4a = slice.appendSlice(s2b, undef)
    always(s4a, true)
    equals(s4a, s2b)

    local s5a = slice.appendSlice(s2a, s2a)
    always(s5a, true)
    local s5b = slice.fromArray({1, 2, 3, 4, 5, 1, 2, 3, 4, 5})
    always(s5b, true)
    equals(s5a, s5b)

    if #s2a ~= 5 then
        error("#s2a ~= 5")
    end

    local s6a = slice.slice(s2a, 2, 5)
    always(s6a, false)
    local s6b = slice.appendSlice(s2a, s6a)
    always(s6b, false)
    local s6c = slice.fromArray({1, 2, 3, 4, 5, 2, 3, 4})
    always(s6c, true)
    equals(s6b, s6c)

    local s7a = slice.appendSlice(s6a, s2a)
    always(s7a, false)
    local s7b = slice.fromArray({2, 3, 4, 1, 2, 3, 4, 5})
    always(s7b, true)
    equals(s7a, s7b)
end

local test_slice_slice = function()
    local s1a = slice.slice(undef, 1, 1)
    always(s1a, true)
    local s1b = slice.fromArray({})
    always(s1b, true)
    equals(s1a, s1b)

    local s2a = slice.slice(s1b, 1, 1)
    always(s2a, true)
    equals(s2a, s1b)
    if s2a.off ~= 0 then
        error("s2b.off ~= 0")
    end

    local s3a = slice.fromArray({1, 2, 3, 4, 5})
    always(s3a, true)
    local s3b = slice.slice(s3a)
    always(s3b, true)
    equals(s3a, s3b)

    local s4a = slice.slice(s3a, 2)
    always(s4a, false)
    local s4b = slice.fromArray({2, 3, 4, 5})
    always(s4b, true)
    equals(s4a, s4b)
    if s4a.off ~= 1 then
        error("s4a.off ~= 1")
    end

    local s5a = slice.slice(s3a, nil, 4)
    always(s5a, false)
    local s5b = slice.fromArray({1, 2, 3})
    always(s5b, true)
    equals(s5a, s5b)

    local s6a = slice.slice(s3a, 2, 5)
    always(s6a, false)
    local s6b = slice.fromArray({2, 3, 4})
    always(s6b, true)
    equals(s6a, s6b)

    if #s3a ~= 5 then
        error("#s3a ~= 5")
    end
    if s6a.off ~= 1 then
        error("s6a.off ~= 1")
    end

    local s7a = slice.slice(s3a, 1, 6)
    always(s7a, false)
    equals(s7a, s3a)

    local s8a = slice.slice(s3a, 6)
    always(s8a, false)
    equals(s8a, slice.fromArray({}))

    local s9a = slice.slice(s6a, 2, 3)
    always(s9a, false)
    equals(s9a, slice.fromArray({3}))
    if s9a.off ~= 2 then
        error("s9a.off ~= 2")
    end

    local s10a = slice.appendSlice(slice.slice(s3a, nil, 3), slice.slice(s3a, 4))
    always(s10a, false)
    local s10b = slice.fromArray({1, 2, 4, 5})
    always(s10b, true)
    equals(s10a, s10b)
    local s10c = slice.fromArray({1, 2, 4, 5, 5})
    always(s10c, true)
    equals(s3a, s10c)
end

local test_slice_copy = function()
    local n1 = slice.copy(undef, undef)
    if n1 ~= 0 then
        error("n1 ~= 0")
    end

    local n2 = slice.copy(slice.fromArray{1, 2, 3}, undef)
    if n2 ~= 0 then
        error("n2 ~= 0")
    end

    local n3 = slice.copy(undef, slice.fromArray({1, 2, 3}))
    if n3 ~= 0 then
        error("n3 ~= 0")
    end

    local s4a = slice.fromArray({1, 2, 3})
    always(s4a, true)
    local s4b = slice.fromArray({4, 5, 6})
    always(s4b, true)
    local n4 = slice.copy(s4a, s4b)
    always(s4a, true)
    if n4 ~= 3 then
        error("n4 ~= 3")
    end
    local s4c = slice.fromArray({1, 2, 3})
    always(s4c, true)
    equals(s4a, s4c)

    local s5a = slice.fromArray({5, 4, 3, 2, 1})
    always(s5a, true)
    local n5 = slice.copy(s4a, s5a)
    always(s4a, true)
    if n5 ~= 3 then
        error("n5 ~= 3")
    end
    local s5b = slice.slice(s5a, 1, 4)
    always(s5b, false)
    equals(s4a, s5b)

    local s6a = slice.slice(s5a, 2)
    always(s6a, false)
    local n6 = slice.copy(s6a, s4a)
    always(s6a, false)
    local s6b = slice.fromArray({5, 4, 3, 1})
    always(s6b, true)
    equals(s6a, s6b)
    local s6c = slice.fromArray({5, 5, 4, 3, 1})
    always(s6c, false)
    equals(s5a, s6c)
end

local test_slice_clone = function()
    if slice.clone(undef) ~= undef then
        error("slice.clone(undef) ~= undef")
    end

    local s2a = slice.fromArray({1, 2, 3, 4, 5})
    always(s2a, true)
    local s2b = slice.clone(s2a)
    always(s2b, true)
    equals(s2a, s2b)
    if s2a.data == s2b.data then
        error("s2a.data == s2b.data")
    end

    local s3a = slice.clone(s2a, 2)
    always(s3a, true)
    local s3b = slice.fromArray({2, 3, 4, 5})
    always(s3b, true)
    equals(s3a, s3b)

    local s4a = slice.clone(s2a, nil, 5)
    always(s4a, true)
    local s4b = slice.fromArray({1, 2, 3, 4})
    always(s4b, true)
    equals(s4a, s4b)

    local s5a = slice.clone(s2a, 6)
    always(s5a, true)
    equals(s5a, slice.fromArray({}))

    local s6a = slice.clone(s2a, 2, 5)
    always(s6a, true)
    local s6b = slice.fromArray({2, 3, 4})
    always(s6b, true)
    equals(s6a, s6b)

    local s7a = slice.slice(s2a, 2, 5)
    always(s7a, false)
    local s7b = slice.clone(s7a, 2)
    always(s7b, true)
    equals(s7b, slice.fromArray({3, 4}))
end

test_slice_make()
test_slice_fromArray()
test_slice_toArray()
test_slice_append()
test_slice_appendArray()
test_slice_appendSlice()
test_slice_slice()
test_slice_copy()
test_slice_clone()
