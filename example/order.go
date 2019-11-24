package example

var order1 = order2 - 10
var (
	order2 = order3 - 10
	order3 = 100
	order4 = order1 - 10
	order5 = order6()
)

func order6() int {
	return 500
}

var _, _ = order4, order5
