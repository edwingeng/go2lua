package example

func forLoop1() {
	for {
		break
	}
}

func forLoop2() {
	for i := 0; i < 3; i++ {
		// Empty
	}
}

func forLoop3() {
	i := 0
	for i < 1 {
		break
	}
}

func forLoop4() {
	for i := 0; ; {
		println(i)
		break
	}
}

func forLoop5() {
	i := 0
	for ; ; i++ {
		if i >= 1 {
			break
		}
	}
}

func forLoop6() {
	for i := 0; i < 3; i++ {
		if i >= 1 {
			continue
		}
		println(i)
	}

	for i := 0; i < 3; i++ {
		if i >= 1 {
			continue
		}
		println(i)
	}
}

func forLoop7() {
pos1:
	for i := 0; i < 3; i++ {
		for {
			break pos1
		}
	}
}

func forLoop8() {
pos1:
	for i := 0; i < 3; i++ {
		for {
			continue pos1
		}
	}
}

func forLoop9(n int) {
	if n > 0 {
		goto pos1
	}
	println(100)

pos1:
	for i := 0; i < 3; i++ {
		for {
			continue pos1
		}
	}
}

func forLoop10() {
	for i := 0; i < 3; i++ {
		i := i * 10
		println(i)
	}
}
