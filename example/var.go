package example

var myNumber1 int
var MyNumber2 int = 200
var myNumber3, myNumber4 int
var MyNumber5, MyNumber6 = 500, 600
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
	var myUint8 uint8 = 0b1011
	var myUint16 uint16 = 0o660
	var myUint32 uint32 = 0660
	var myUint64 uint64 = 0x01_FF
	var myUint uint = 1_000_000
	var myRune rune
	var myByte byte
	var myUintptr uintptr
	var myFloat32 float32 = 3.14
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

type Foo1 struct {
	String1 string
	Num1    int
}

type Foo2 struct {
	Foo1
	String2 string
	Num2    int
}

type Foo3 struct {
	string
	int
}

func var2() {
	var myArray1 [3]int
	var myArray2 [3][3]int

	var mySlice1 []int
	var mySlice2 [][]int

	var foo1 Foo1
	var foo2 Foo2
	var foo3 Foo3
	var foo4 struct {
		Num int
		string
	}
	var foo5 = new(Foo1)

	var ptr4 *Foo1

	var fn1 func()
	var fn2 func(int)
	var fn3 func() int

	var obj1 interface{}
	var obj2 interface {
		Print()
	}

	var map1 map[string]int
	var map2 map[int]struct{}
	var map3 map[*Foo1]int
	var map4 map[*Foo1]Foo3
	var map5 map[int]*Foo1
	var map6 map[int]interface{}
	var map7 = map[int]int{
		1: 10,
		2: 20,
	}
	var map8 = map[string]int{
		"1": 10,
		"2": 20,
	}
	var map9 = map[*Foo1]struct{}{
		&foo1: {},
	}

	_, _ = myArray1, myArray2
	_, _ = mySlice1, mySlice2
	_, _, _, _, _ = foo1, foo2, foo3, foo4, foo5
	_ = ptr4
	_, _, _ = fn1, fn2, fn3
	_, _ = obj1, obj2
	_, _, _, _, _, _ = map1, map2, map3, map4, map5, map6
	_, _, _ = map7, map8, map9
}
