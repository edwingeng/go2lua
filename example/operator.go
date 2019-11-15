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
