package unsupported

import (
	"go/ast"
	"go/token"

	"github.com/edwingeng/go2lua/utils"
	"golang.org/x/tools/go/analysis"
)

type typeWithPos struct {
	typ ast.Expr
	pos token.Position
}

func newTypeWithPos(pass *analysis.Pass, typ ast.Expr) *typeWithPos {
	pos := pass.Fset.Position(typ.Pos())
	return &typeWithPos{typ, pos}
}

type orderedTypes []*typeWithPos

func (p orderedTypes) Len() int      { return len(p) }
func (p orderedTypes) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p orderedTypes) Less(i, j int) bool {
	return utils.PositionLess(p[i].pos, p[j].pos)
}
