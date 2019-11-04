package example

func Fibs(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	return Fibs(n-1) + Fibs(n-2)
}
