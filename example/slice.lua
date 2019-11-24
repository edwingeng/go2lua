-- package: example

local gopkg = _G["github.com/edwingeng/go2lua/example"]
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

slice1 = function()
    local myInt8s = goslice.make(nil, 0)
    local myInt16s = goslice.make(nil, 0)
    local myInt32s = goslice.make(nil, 0)
    local myInt64s = goslice.make(nil, 0)
    local myInts = goslice.make(nil, 0)
    local myUint8s = goslice.make(nil, 0)
    local myUint16s = goslice.make(nil, 0)
    local myUint32s = goslice.make(nil, 0)
    local myUint64s = goslice.make(nil, 0)
    local myUints = goslice.make(nil, 0)
    local myRunes = goslice.make(nil, 0)
    local myBytes = goslice.make(nil, 0)
    local myUintptrs = goslice.make(nil, 0)
    local myFloat32s = goslice.make(nil, 0)
    local myFloat64s = goslice.make(nil, 0)
    local myBools = goslice.make(nil, 0)
    local myStrings = goslice.make(nil, 0)

    local _ = myInt8s
    local _ = myInt16s
    local _ = myInt32s
    local _ = myInt64s
    local _ = myInts
    local _ = myUint8s
    local _ = myUint16s
    local _ = myUint32s
    local _ = myUint64s
    local _ = myUints
    local _ = myRunes
    local _ = myBytes
    local _ = myUintptrs
    local _ = myFloat32s
    local _ = myFloat64s
    local _ = myBools
    local _ = myStrings
end

slice2 = function()
    local myInt8s = goslice.make(goslice.newNumberArray, 10)
    local myInt16s = goslice.make(goslice.newNumberArray, 10)
    local myInt32s = goslice.make(goslice.newNumberArray, 10)
    local myInt64s = goslice.make(goslice.newNumberArray, 10)
    local myInts = goslice.make(goslice.newNumberArray, 10)
    local myUint8s = goslice.make(goslice.newNumberArray, 10)
    local myUint16s = goslice.make(goslice.newNumberArray, 10)
    local myUint32s = goslice.make(goslice.newNumberArray, 10)
    local myUint64s = goslice.make(goslice.newNumberArray, 10)
    local myUints = goslice.make(goslice.newNumberArray, 10)
    local myRunes = goslice.make(goslice.newNumberArray, 10)
    local myBytes = goslice.make(goslice.newNumberArray, 10)
    local myUintptrs = goslice.make(goslice.newNumberArray, 10)
    local myFloat32s = goslice.make(goslice.newNumberArray, 10)
    local myFloat64s = goslice.make(goslice.newNumberArray, 10)
    local myBools = goslice.make(goslice.newBoolArray, 10)
    local myStrings = goslice.make(goslice.newStringArray, 10)

    local _ = myInt8s
    local _ = myInt16s
    local _ = myInt32s
    local _ = myInt64s
    local _ = myInts
    local _ = myUint8s
    local _ = myUint16s
    local _ = myUint32s
    local _ = myUint64s
    local _ = myUints
    local _ = myRunes
    local _ = myBytes
    local _ = myUintptrs
    local _ = myFloat32s
    local _ = myFloat64s
    local _ = myBools
    local _ = myStrings
end

slice3 = function()
    local myInt8s = goslice.append(goslice.make(nil, 0), 1)
    local myInt16s = goslice.append(goslice.make(nil, 0), 1)
    local myInt32s = goslice.append(goslice.make(nil, 0), 1)
    local myInt64s = goslice.append(goslice.make(nil, 0), 1)
    local myInts = goslice.append(goslice.make(nil, 0), 1)
    local myUint8s = goslice.append(goslice.make(nil, 0), 1)
    local myUint16s = goslice.append(goslice.make(nil, 0), 1)
    local myUint32s = goslice.append(goslice.make(nil, 0), 1)
    local myUint64s = goslice.append(goslice.make(nil, 0), 1)
    local myUints = goslice.append(goslice.make(nil, 0), 1)
    local myRunes = goslice.append(goslice.make(nil, 0), 1)
    local myBytes = goslice.append(goslice.make(nil, 0), 1)
    local myUintptrs = goslice.append(goslice.make(nil, 0), 1)
    local myFloat32s = goslice.append(goslice.make(nil, 0), 1)
    local myFloat64s = goslice.append(goslice.make(nil, 0), 1)
    local myBools = goslice.append(goslice.make(nil, 0), false)
    local myStrings = goslice.append(goslice.make(nil, 0), "hello")

    local _ = myInt8s
    local _ = myInt16s
    local _ = myInt32s
    local _ = myInt64s
    local _ = myInts
    local _ = myUint8s
    local _ = myUint16s
    local _ = myUint32s
    local _ = myUint64s
    local _ = myUints
    local _ = myRunes
    local _ = myBytes
    local _ = myUintptrs
    local _ = myFloat32s
    local _ = myFloat64s
    local _ = myBools
    local _ = myStrings
end

slice4 = function()
    local myInt8s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myInt16s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myInt32s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myInt64s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myInts = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUint8s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUint16s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUint32s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUint64s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUints = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myRunes = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myBytes = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myUintptrs = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myFloat32s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myFloat64s = goslice.appendArray(goslice.make(nil, 0), {1, 2, 3})
    local myBools = goslice.appendArray(goslice.make(nil, 0), {false, true, false})
    local myStrings = goslice.appendArray(goslice.make(nil, 0), {"hello", "world", "!"})

    local _ = myInt8s
    local _ = myInt16s
    local _ = myInt32s
    local _ = myInt64s
    local _ = myInts
    local _ = myUint8s
    local _ = myUint16s
    local _ = myUint32s
    local _ = myUint64s
    local _ = myUints
    local _ = myRunes
    local _ = myBytes
    local _ = myUintptrs
    local _ = myFloat32s
    local _ = myFloat64s
    local _ = myBools
    local _ = myStrings
end

slice5 = function()
    local myInts = goslice.fromArray({1, 2, 3})
    local myBools = goslice.fromArray({false, true, false})
    local myStrings = goslice.fromArray({"hello", "world", "!"})

    myInts = goslice.appendSlice(myInts, myInts)
    myBools = goslice.appendSlice(myBools, myBools)
    myStrings = goslice.appendSlice(myStrings, myStrings)
end

slice6 = function()
    local myInts1 = goslice.fromArray({1, 2, 3})
    local myInts2 = goslice.make(goslice.newNumberArray, 10)
    local myInts3 = goslice.make(goslice.newNumberArray, 1)
    local myInts4 = goslice.make(nil, 0)
    local myInts5 = myInts1

    print(goslice.copy(myInts2, myInts1))
    print(goslice.copy(myInts3, myInts1))
    print(goslice.copy(myInts4, myInts1))
    print(goslice.copy(myInts5, myInts4))
end

slice7 = function()
    local myInts1 = goslice.fromArray({1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
    local myInts2 = goslice.slice(myInts1, 1)
    local myInts3 = goslice.slice(myInts1, nil, myInts1.len + 1)
    local myInts4 = goslice.slice(myInts1, 2)
    local myInts5 = goslice.slice(myInts1, nil, myInts1.len - 1 + 1)
    local myInts6 = goslice.slice(myInts1, 3, 6)
    local myInts7 = goslice.slice(myInts1)

    myInts6[1] = 0
    myInts6[3] = 0

    local myInts8 = goslice.make(nil, 0)
    local myInts9 = goslice.slice(myInts8, 1)

    local _, _, _, _, _, _ = myInts2, myInts3, myInts4, myInts5, myInts7, myInts9
end

slice8 = function()
    local myStrings1 = goslice.make(goslice.newStringArray, 10)
    print(goslice.cap(myStrings1))
    local myStrings2 = goslice.slice(myStrings1, 3, 8)
    print(goslice.cap(myStrings2))
    local myStrings3 = goslice.slice(myStrings2, 2, 4)
    print(goslice.cap(myStrings3))
end

return function() end
