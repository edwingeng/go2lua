package utils

import (
	"go/token"
)

func LuaOpPrecedenceFromGoOp(op token.Token) int {
	switch op {
	case token.LOR:
		return 1
	case token.LAND:
		return 2
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		return 3
	case token.OR:
		return 4
	case token.XOR:
		return 5
	case token.AND, token.AND_NOT:
		return 6
	case token.SHL, token.SHR:
		return 7
	case token.ADD, token.SUB:
		return 8
	case token.MUL, token.QUO, token.REM:
		return 9
	default:
		return 0
	}
}
