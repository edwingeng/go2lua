package example

func panic1(b1, b2 bool) {
	if b1 {
		panic("hello")
	}
	if b2 {
		panic(100)
	}
}

func panic2(b1 bool) (int, int, int) {
	defer func() {
		println(b1)
	}()
	if b1 {
		panic("world")
	}
	return 1, 2, 3
}
