-- package: example

local myNumber1 = 0
MyNumber2 = 0
local myNumber3, myNumber4 = 0, 0
MyNumber5, MyNumber6 = 0, 0
local myNumber7 = 0
MyNumber8 = 0

local myString1 = ""
local myString2 = ""
MyString3 = ""

local var1 = function()
    local myInt8 = 0
    local myInt16 = 0
    local myInt32 = 0
    local myInt64 = 0
    local myInt = 0
    local myUint8 = 0
    local myUint16 = 0
    local myUint32 = 0
    local myUint64 = 0
    local myUint = 0
    local myRune = 0
    local myByte = 0
    local myUintptr = 0
    local myFloat32 = 0
    local myFloat64 = 0
    local myBool = false
    local myString = ""

    local _ = myInt8
    local _ = myInt16
    local _ = myInt32
    local _ = myInt64
    local _ = myInt
    local _ = myUint8
    local _ = myUint16
    local _ = myUint32
    local _ = myUint64
    local _ = myUint
    local _ = myRune
    local _ = myByte
    local _ = myUintptr
    local _ = myFloat32
    local _ = myFloat64
    local _ = myBool
    local _ = myString

    local _ = myNumber1
    local _ = MyNumber2
    local _ = myNumber3
    local _ = myNumber4
    local _ = MyNumber5
    local _ = MyNumber6
    local _ = myNumber7
    local _ = MyNumber8

    local _, _, _ = myString1, myString2, MyString3
end

local var2 = function()
    local myArray1 = {0, 0, 0}
    local myArray2 = {{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

    local mySlice1 = slice.make(nil, 0)
    local mySlice2 = slice.make(nil, 0)

    local foo1 = {String1 = "", Num1 = 0}
    local foo2 = {Foo1 = {String1 = "", Num1 = 0}, String2 = "", Num2 = 0}
    local foo3 = {string = "", int = 0}
    local foo4 = {Num = 0, string = ""}

    local ptr4 = undef

    local fn1 = undef
    local fn2 = undef
    local fn3 = undef

    local obj1 = undef
    local obj2 = undef

    local map1 = {}
    local map2 = {}
    local map3 = {}
    local map4 = {}
    local map5 = {}
    local map6 = {}

    local _, _ = myArray1, myArray2
    local _, _ = mySlice1, mySlice2
    local _, _, _, _ = foo1, foo2, foo3, foo4
    local _ = ptr4
    local _, _, _ = fn1, fn2, fn3
    local _, _ = obj1, obj2
    local _, _, _, _, _, _ = map1, map2, map3, map4, map5, map6
end
