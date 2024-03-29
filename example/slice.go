package example

func slice1() {
	var myInt8s []int8
	var myInt16s []int16
	var myInt32s []int32
	var myInt64s []int64
	var myInts []int
	var myUint8s []uint8
	var myUint16s []uint16
	var myUint32s []uint32
	var myUint64s []uint64
	var myUints []uint
	var myRunes []rune
	var myBytes []byte
	var myUintptrs []uintptr
	var myFloat32s []float32
	var myFloat64s []float64
	var myBools []bool
	var myStrings []string

	_ = myInt8s
	_ = myInt16s
	_ = myInt32s
	_ = myInt64s
	_ = myInts
	_ = myUint8s
	_ = myUint16s
	_ = myUint32s
	_ = myUint64s
	_ = myUints
	_ = myRunes
	_ = myBytes
	_ = myUintptrs
	_ = myFloat32s
	_ = myFloat64s
	_ = myBools
	_ = myStrings
}

func slice2() {
	var myInt8s = make([]int8, 10)
	var myInt16s = make([]int16, 10)
	var myInt32s = make([]int32, 10)
	var myInt64s = make([]int64, 10)
	var myInts = make([]int, 10)
	var myUint8s = make([]uint8, 10)
	var myUint16s = make([]uint16, 10)
	var myUint32s = make([]uint32, 10)
	var myUint64s = make([]uint64, 10)
	var myUints = make([]uint, 10)
	var myRunes = make([]rune, 10)
	var myBytes = make([]byte, 10)
	var myUintptrs = make([]uintptr, 10)
	var myFloat32s = make([]float32, 10)
	var myFloat64s = make([]float64, 10)
	var myBools = make([]bool, 10)
	var myStrings = make([]string, 10)

	_ = myInt8s
	_ = myInt16s
	_ = myInt32s
	_ = myInt64s
	_ = myInts
	_ = myUint8s
	_ = myUint16s
	_ = myUint32s
	_ = myUint64s
	_ = myUints
	_ = myRunes
	_ = myBytes
	_ = myUintptrs
	_ = myFloat32s
	_ = myFloat64s
	_ = myBools
	_ = myStrings
}

func slice3() {
	myInt8s := append([]int8(nil), 1)
	myInt16s := append([]int16(nil), 1)
	myInt32s := append([]int32(nil), 1)
	myInt64s := append([]int64(nil), 1)
	myInts := append([]int(nil), 1)
	myUint8s := append([]uint8(nil), 1)
	myUint16s := append([]uint16(nil), 1)
	myUint32s := append([]uint32(nil), 1)
	myUint64s := append([]uint64(nil), 1)
	myUints := append([]uint(nil), 1)
	myRunes := append([]rune(nil), 1)
	myBytes := append([]byte(nil), 1)
	myUintptrs := append([]uintptr(nil), 1)
	myFloat32s := append([]float32(nil), 1)
	myFloat64s := append([]float64(nil), 1)
	myBools := append([]bool(nil), false)
	myStrings := append([]string(nil), "hello")

	_ = myInt8s
	_ = myInt16s
	_ = myInt32s
	_ = myInt64s
	_ = myInts
	_ = myUint8s
	_ = myUint16s
	_ = myUint32s
	_ = myUint64s
	_ = myUints
	_ = myRunes
	_ = myBytes
	_ = myUintptrs
	_ = myFloat32s
	_ = myFloat64s
	_ = myBools
	_ = myStrings
}

func slice4() {
	myInt8s := append([]int8(nil), 1, 2, 3)
	myInt16s := append([]int16(nil), 1, 2, 3)
	myInt32s := append([]int32(nil), 1, 2, 3)
	myInt64s := append([]int64(nil), 1, 2, 3)
	myInts := append([]int(nil), 1, 2, 3)
	myUint8s := append([]uint8(nil), 1, 2, 3)
	myUint16s := append([]uint16(nil), 1, 2, 3)
	myUint32s := append([]uint32(nil), 1, 2, 3)
	myUint64s := append([]uint64(nil), 1, 2, 3)
	myUints := append([]uint(nil), 1, 2, 3)
	myRunes := append([]rune(nil), 1, 2, 3)
	myBytes := append([]byte(nil), 1, 2, 3)
	myUintptrs := append([]uintptr(nil), 1, 2, 3)
	myFloat32s := append([]float32(nil), 1, 2, 3)
	myFloat64s := append([]float64(nil), 1, 2, 3)
	myBools := append([]bool(nil), false, true, false)
	myStrings := append([]string(nil), "hello", "world", "!")

	_ = myInt8s
	_ = myInt16s
	_ = myInt32s
	_ = myInt64s
	_ = myInts
	_ = myUint8s
	_ = myUint16s
	_ = myUint32s
	_ = myUint64s
	_ = myUints
	_ = myRunes
	_ = myBytes
	_ = myUintptrs
	_ = myFloat32s
	_ = myFloat64s
	_ = myBools
	_ = myStrings
}

func slice5() {
	myInts := []int{1, 2, 3}
	myBools := []bool{false, true, false}
	myStrings := []string{"hello", "world", "!"}

	myInts = append(myInts, myInts...)
	myBools = append(myBools, myBools...)
	myStrings = append(myStrings, myStrings...)
}

func slice6() {
	myInts1 := []int{1, 2, 3}
	myInts2 := make([]int, 10)
	myInts3 := make([]int, 1)
	var myInts4 []int
	myInts5 := myInts1

	println(copy(myInts2, myInts1))
	println(copy(myInts3, myInts1))
	println(copy(myInts4, myInts1))
	println(copy(myInts5, myInts4))
}

func slice7() {
	myInts1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	myInts2 := myInts1[0:]
	myInts3 := myInts1[:len(myInts1)]
	myInts4 := myInts1[1:]
	myInts5 := myInts1[:len(myInts1)-1]
	myInts6 := myInts1[2:5]
	myInts7 := myInts1[:]

	myInts6[0] = 0
	myInts6[2] = 0

	var myInts8 []int
	myInts9 := myInts8[0:]

	_, _, _, _, _, _ = myInts2, myInts3, myInts4, myInts5, myInts7, myInts9
}

func slice8() {
	myStrings1 := make([]string, 10)
	println(cap(myStrings1))
	myStrings2 := myStrings1[2:7]
	println(cap(myStrings2))
	myStrings3 := myStrings2[1:3]
	println(cap(myStrings3))
}
