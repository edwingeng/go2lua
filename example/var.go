package example

var myNumber1 int
var MyNumber2 int
var myNumber3, myNumber4 int
var MyNumber5, MyNumber6 int
var myNumber7, MyNumber8 int

var (
	myString1            string
	myString2, MyString3 string
)

func var1() {
	var myInt8 int8
	var myInt16 int16
	var myInt32 int32
	var myInt64 int64
	var myInt int
	var myUint8 uint8
	var myUint16 uint16
	var myUint32 uint32
	var myUint64 uint64
	var myUint uint
	var myRune rune
	var myByte byte
	var myUintptr uintptr
	var myFloat32 float32
	var myFloat64 float64
	var myBool bool
	var myString string

	_ = myInt8
	_ = myInt16
	_ = myInt32
	_ = myInt64
	_ = myInt
	_ = myUint8
	_ = myUint16
	_ = myUint32
	_ = myUint64
	_ = myUint
	_ = myRune
	_ = myByte
	_ = myUintptr
	_ = myFloat32
	_ = myFloat64
	_ = myBool
	_ = myString

	_ = myNumber1
	_ = MyNumber2
	_ = myNumber3
	_ = myNumber4
	_ = MyNumber5
	_ = MyNumber6
	_ = myNumber7
	_ = MyNumber8

	_, _, _ = myString1, myString2, MyString3
}
