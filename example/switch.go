package example

func switch1(n int) {
	switch n {
	}
}

func switch2(n int) {
	switch n {
	case 1:
		println("a", n)
		break
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

func switch4(n int) {
	switch n {
	case 1:
		println("a", n)
		fallthrough
	case 2:
		println("b", n)
		fallthrough
	case 3:
		println("c", n)
		break
	case 4:
		println("d", n)

	case 5:
		println("e", n)
		fallthrough
	default:
		println("f", n)
		fallthrough
	case 6, 7:
		println("g", n)
	}
}

func switch5(n int) {
	switch {
	case n == 1:
		println("a", n)
		fallthrough
	case n == 2:
		println("b", n)
		fallthrough
	case n == 3:
		println("c", n)
		break
	case n == 4:
		println("d", n)

	case n == 5:
		println("e", n)
		fallthrough
	default:
		println("f", n)
		fallthrough
	case n == 6 || n == 7:
		println("g", n)
	}
}
