package example

func switch1(n int) {
	switch n {
	}
}

func switch2(n int) {
	switch n {
	case 1:
		println(1)
	default:
		println(n)
	case 2:
		println(2)
	}
}

func switch3(n int) {
	a := 3
	b := 2
	switch n {
	case 1, 3:
		println(1)
	case a + b:
		println(2)
	default:
		println(n)
	}
}
