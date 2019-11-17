package example

func string1(b byte, r rune) {
	str := "hello"
	_ = 'a'
	_ = '你'
	_ = str[0]
	_ = string(97)
	_ = string(b)
	_ = string(r)
	_ = len(str)
}

func string2(str string) {
	for i := 0; i < len(str); i++ {
		println(i, str[i])
	}

	for i, r := range str {
		println(i, r)
	}
}
