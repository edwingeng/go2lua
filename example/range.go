package example

func rangeMap1() {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200

	for k, v := range m {
		println(k, v)
	}
}

func rangeMap2() {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200
	m["c"] = 300

	for k, v := range m {
		if k == "b" {
			continue
		} else {
			println(k, v)
		}
	}
}

func rangeMap3() {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200

pos1:
	for k, v := range m {
		println(k, v)
		for {
			break pos1
		}
	}
}

func rangeMap4() {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200

pos1:
	for k, v := range m {
		println(k, v)
		for {
			continue pos1
		}
	}
}
