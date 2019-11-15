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
