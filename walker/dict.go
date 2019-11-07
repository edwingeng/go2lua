package walker

var (
	go2LuaFuncMap = map[string]string{
		"print":   "io.write",
		"println": "print",
	}

	go2LuaOperMap = map[string]string{
		`&&`: "and",
		`||`: "or",
		`!`:  "not",
	}
)
