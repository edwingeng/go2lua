package example

func defer1() {
	defer func() {
		println(100)
	}()
}

func defer2() {
	defer func() {
		println(100)
	}()

	println(300)

	defer func() {
		println(200)
	}()
}

func defer3() {
	f1 := func() {
		println(100)
	}
	f2 := func() {
		println(200)
	}

	defer f1()
	defer f2()
	println(300)
}

func defer4() {
	defer func() {
		defer func() {
			defer func() {
				println(100)
			}()
			println(200)
		}()
		println(300)
	}()
}

func defer5() {
	for i := 0; i < 3; i++ {
		defer func() {
			println(i)
		}()
	}

	for i := 0; i < 3; i++ {
		defer func(i int) {
			println(i)
		}(i)
	}
}

func defer6(n1, n2 int) {
	defer func() {
		println(n1, n2)
	}()
}
