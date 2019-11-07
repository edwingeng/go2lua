package example

func If1(n int) {
	if n == 0 {
		println("a1:", n)
	}

	if n == 1 {
		println("b1:", n)
	} else {
		println("b2:", n)
	}

	if n == 1 {
		println("c1:", n)
	} else if n == 2 {
		println("c2:", n)
	} else {
		println("c3:", n)
	}

	if n > 10 {
		if n > 100 {
			println("d1:", n)
		} else {
			println("d2:", n)
		}
	} else {
		if n < 1 {
			println("d3:", n)
		} else {
			println("d4:", n)
		}
	}

	if n > 10 {
		println("e1:", n)
	} else {
		if n == 1 {
			println("e2:", n)
		} else if n == 2 {
			println("e3:", n)
		} else if n == 3 {
			println("e4:", n)
		} else {
			println("e5:", n)
		}
	}
}

func If2(n int) {
	if x1 := n * 10; x1 > 10 {
		println(x1)
	}

	x2 := 0
	if x2 = n * 10; x2 > 10 {
		println(x2)
	}
}
