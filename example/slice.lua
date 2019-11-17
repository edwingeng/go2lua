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
