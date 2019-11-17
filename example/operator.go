package example

func operator1(n int, ok bool) {
	a := +n
	b := -a
	_ = !ok
	d := ^b
	e := &struct{}{}
	f := *e
	_, _, _ = d, e, f
}

func operator2(n1, n2 int, b1, b2 bool) {
	_ = n1 * n2
	_ = n1 / n2
	_ = n1 % n2
	_ = n1 << uint(n2)
	_ = n1 >> uint(n2)
	_ = n1 & n2
	_ = n1 &^ n2
	_ = n1 + n2
	_ = n1 - n2
	_ = n1 | n2
	_ = n1 ^ n2
	_ = n1 == n2
	_ = n1 != n2
	_ = n1 < n2
	_ = n1 <= n2
	_ = n1 > n2
	_ = n1 >= n2
	_ = b1 && b2
	_ = b1 || b2
}

func operator3(n1, n2 int) {
	_ = n1*n2 + n1/n2
	_ = (n1 * n2) + (n1 / n2)
	_ = n1 * (n2 + n1) / n2
	_ = (n1 * n2) + (n1/n2)*(n1+n2)
	_ = ((n1 * n2) + (n1 / n2)) * (n1 + n2)
	_ = n1 + n2 + n1 + n2 + n1
	_ = n1*n2 + n1/n2 + n1%n2 + n1&n2
}

func operator4(str1, str2 string, b1, b2 byte, r1, r2 rune) {
	_ = "x" + "y"
	_ = str1 + str2
	_ = "x" + str1 + "y"
	_ = str1 + "x" + str2
	_ = b1 + b2
	_ = r1 + r2
}

func operator5(n1, n2 uint32) {
	_ = n1 | n2 ^ n2
	_ = n1 | n2 + n1 | n2
	_ = n1 + n2 | n1 + n2
	_ = n1 | n2 ^ n1
	_ = n1<<n2 + n1<<n2*n2>>n1
}
