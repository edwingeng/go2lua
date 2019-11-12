package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type NodeError struct {
	Err  error
	Node ast.Node
}

func NewNodeError(err error, node ast.Node) NodeError {
	return NodeError{
		Err:  err,
		Node: node,
	}
}

func PrintErrors(pass *analysis.Pass, a ...NodeError) {
	var buf bytes.Buffer
	for _, v := range a {
		pos := pass.Fset.Position(v.Node.Pos())
		_, _ = fmt.Fprintf(&buf, "%s:%d:%d\n", pos.Filename, pos.Line, pos.Column)
		_, _ = fmt.Fprintf(&buf, "    %+v\n", strings.ReplaceAll(v.Err.Error(), "\n", "\n    "))
	}
	buf.WriteByte('\n')
	_, _ = os.Stderr.Write(buf.Bytes())
}

func PositionLess(p1, p2 token.Position) bool {
	if p1.Filename < p2.Filename {
		return true
	} else if p1.Filename == p2.Filename {
		if p1.Offset < p2.Offset {
			return true
		}
	}
	return false
}
