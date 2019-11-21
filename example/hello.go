package example

func init() {
	print("Hello ")
	println("World!")
}

func Add(n1, n2 int) int {
	return n1 + n2
}

func Sub(n1 int, n2 int) int {
	n3 := n1 - n2
	return n3
}

func Fibs(n int) int {
	if n == 1 || n == 2 {
		return 1
	}

	return Fibs(n-1) + Fibs(n-2)
}
