package walker

var (
	go2LuaFuncMap = map[string]string{
		"print":   "io.write",
		"println": "print",
	}

	go2LuaBinaryOperMap = map[string]string{
		`&&`: " and ",
		`||`: " or ",
		`!`:  " not ",
		"^":  " ~ ",
		"!=": " ~= ",
		"&^": " & ~",
	}
)
