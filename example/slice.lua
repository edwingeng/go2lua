-- package: example

local slice1 = function()
    local myInt8s = slice.make(nil, 0)
    local myInt16s = slice.make(nil, 0)
    local myInt32s = slice.make(nil, 0)
    local myInt64s = slice.make(nil, 0)
    local myInts = slice.make(nil, 0)
    local myUint8s = slice.make(nil, 0)
    local myUint16s = slice.make(nil, 0)
    local myUint32s = slice.make(nil, 0)
    local myUint64s = slice.make(nil, 0)
    local myUints = slice.make(nil, 0)
    local myRunes = slice.make(nil, 0)
    local myBytes = slice.make(nil, 0)
    local myUintptrs = slice.make(nil, 0)
    local myFloat32s = slice.make(nil, 0)
    local myFloat64s = slice.make(nil, 0)
    local myBools = slice.make(nil, 0)
    local myStrings = slice.make(nil, 0)

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

local slice2 = function()
    local myInt8s = slice.make(slice.newNumberArray, 10)
    local myInt16s = slice.make(slice.newNumberArray, 10)
    local myInt32s = slice.make(slice.newNumberArray, 10)
    local myInt64s = slice.make(slice.newNumberArray, 10)
    local myInts = slice.make(slice.newNumberArray, 10)
    local myUint8s = slice.make(slice.newNumberArray, 10)
    local myUint16s = slice.make(slice.newNumberArray, 10)
    local myUint32s = slice.make(slice.newNumberArray, 10)
    local myUint64s = slice.make(slice.newNumberArray, 10)
    local myUints = slice.make(slice.newNumberArray, 10)
    local myRunes = slice.make(slice.newNumberArray, 10)
    local myBytes = slice.make(slice.newNumberArray, 10)
    local myUintptrs = slice.make(slice.newNumberArray, 10)
    local myFloat32s = slice.make(slice.newNumberArray, 10)
    local myFloat64s = slice.make(slice.newNumberArray, 10)
    local myBools = slice.make(slice.newBoolArray, 10)
    local myStrings = slice.make(slice.newStringArray, 10)

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

local slice3 = function()
    local myInt8s = slice.append(slice.make(nil, 0), 1)
    local myInt16s = slice.append(slice.make(nil, 0), 1)
    local myInt32s = slice.append(slice.make(nil, 0), 1)
    local myInt64s = slice.append(slice.make(nil, 0), 1)
    local myInts = slice.append(slice.make(nil, 0), 1)
    local myUint8s = slice.append(slice.make(nil, 0), 1)
    local myUint16s = slice.append(slice.make(nil, 0), 1)
    local myUint32s = slice.append(slice.make(nil, 0), 1)
    local myUint64s = slice.append(slice.make(nil, 0), 1)
    local myUints = slice.append(slice.make(nil, 0), 1)
    local myRunes = slice.append(slice.make(nil, 0), 1)
    local myBytes = slice.append(slice.make(nil, 0), 1)
    local myUintptrs = slice.append(slice.make(nil, 0), 1)
    local myFloat32s = slice.append(slice.make(nil, 0), 1)
    local myFloat64s = slice.append(slice.make(nil, 0), 1)
    local myBools = slice.append(slice.make(nil, 0), false)
    local myStrings = slice.append(slice.make(nil, 0), "hello")

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

local slice4 = function()
    local myInt8s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myInt16s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myInt32s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myInt64s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myInts = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUint8s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUint16s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUint32s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUint64s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUints = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myRunes = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myBytes = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myUintptrs = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myFloat32s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myFloat64s = slice.appendArray(slice.make(nil, 0), {1, 2, 3})
    local myBools = slice.appendArray(slice.make(nil, 0), {false, true, false})
    local myStrings = slice.appendArray(slice.make(nil, 0), {"hello", "world", "!"})

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

local slice5 = function()
    local myInts = slice.fromArray({1, 2, 3})
    local myBools = slice.fromArray({false, true, false})
    local myStrings = slice.fromArray({"hello", "world", "!"})

    myInts = slice.appendSlice(myInts, myInts)
    myBools = slice.appendSlice(myBools, myBools)
    myStrings = slice.appendSlice(myStrings, myStrings)
end

local slice6 = function()
    local myInts1 = slice.fromArray({1, 2, 3})
    local myInts2 = slice.make(slice.newNumberArray, 10)
    local myInts3 = slice.make(slice.newNumberArray, 1)
    local myInts4 = slice.make(nil, 0)
    local myInts5 = myInts1

    print(slice.copy(myInts2, myInts1))
    print(slice.copy(myInts3, myInts1))
    print(slice.copy(myInts4, myInts1))
    print(slice.copy(myInts5, myInts4))
end
