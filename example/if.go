package example

func If(n int) {
	if n == 0 {
		println("a1:", 0)
	}

	if n == 1 {
		println("b1:", 1)
	} else {
		println("b2:", n)
	}

	if n == 1 {
		println("c1:", 1)
	} else if n == 2 {
		println("c2:", 2)
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
			println("e2:", 1)
		} else if n == 2 {
			println("e3:", 2)
		} else if n == 3 {
			println("e4:", 3)
		} else {
			println("e5:", n)
		}
	}
}
