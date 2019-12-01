package example

const (
	iotaNumA1 = iota
	iotaNumA2
	iotaNumA3
	iotaNumB1 = 10
	iotaNumB2
	iotaNumC1 = iotaNumB1 + 10 + iota
	iotaNumC2

	iotaNumC3
)

const (
	iotaNumX1 = 10
	iotaNumD1 = iota
	iotaNumD2
	iotaNumD3
	iotaNumD4, iotaNumD5 = 10, 20
)

func iota1() {
	const (
		iotaNumA1 = iota
		iotaNumA2
		iotaNumA3
		iotaNumB1 = 10
		iotaNumB2
		iotaNumC1 = iotaNumB1 + 10 + iota
		iotaNumC2
	)
}

func init() {
	if iotaNumA1 != 0 {
		panic("iotaNumA1 != 0")
	}
	if iotaNumA2 != 1 {
		panic("iotaNumA2 != 1")
	}
	if iotaNumA3 != 2 {
		panic("iotaNumA3 != 2")
	}
	if iotaNumB1 != 10 {
		panic("iotaNumB1 != 10")
	}
	if iotaNumB2 != 10 {
		panic("iotaNumB2 != 10")
	}
	if iotaNumC1 != 25 {
		panic("iotaNumC1 != 25")
	}
	if iotaNumC2 != 26 {
		panic("iotaNumC2 != 26")
	}
	if iotaNumC3 != 27 {
		panic("iotaNumC3 != 27")
	}
	if iotaNumD1 != 1 {
		panic("iotaNumD1 != 1")
	}
	if iotaNumD2 != 2 {
		panic("iotaNumD2 != 2")
	}
	if iotaNumD3 != 3 {
		panic("iotaNumD3 != 3")
	}
	if iotaNumD4 != 10 {
		panic("iotaNumD4 != 10")
	}
	if iotaNumD5 != 20 {
		panic("iotaNumD5 != 20")
	}
}
