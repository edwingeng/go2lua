package example

func panic1(b1, b2 bool) {
	if b1 {
		panic("hello")
	}
	if b2 {
		panic(100)
	}
}
