package example

func switch1(n int) {
	switch n {
	}
}

func switch2(n int) {
	switch n {
	case 1:
		println("a", n)
	default:
		println("c", n)
	case 2:
		println("b", n)
	}
}

func switch3(n int) {
	a := 3
	b := 2
	switch n {
	case 1, 3:
		println("a", n)
	case a + b:
		println("b", n)
	default:
		println("c", n)
	}
}
