package example

func forLoop1() {
	for {
		break
	}
}

func forLoop2() {
	for i := 0; i < 1; i++ {
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
}

func forRangeMap() {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200

	for k, v := range m {
		println(k, v)
	}
}
