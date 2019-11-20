package example

func func1() {
	var f1 = func() {
		println("f1")
	}
	f1()

	f2 := func() {
		println("f2")
	}
	f2()

	func() {
		println("f3")
	}()

	func(n1, n2, n3 int) {
		println(n1, n2, n3)
	}(1, 2, 3)
}

func func2() {
	func(cb func(n int)) {
		if cb != nil {
			cb(100)
		}
	}(func(n int) {
		println(n)
	})
}
