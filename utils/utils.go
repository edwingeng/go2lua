package utils

import (
	"bytes"
	"fmt"
	"go/ast"
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
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[j].Node.Pos() > a[i].Node.Pos() {
				a[i], a[j] = a[j], a[i]
			}
		}
	}

	var buf bytes.Buffer
	for _, v := range a {
		pos := pass.Fset.Position(v.Node.Pos())
		_, _ = fmt.Fprintf(&buf, "%s:%d:%d\n", pos.Filename, pos.Line, pos.Column)
		_, _ = fmt.Fprintf(&buf, "    %+v\n", strings.ReplaceAll(v.Err.Error(), "\n", "\n    "))
	}
	buf.WriteByte('\n')
	_, _ = os.Stderr.Write(buf.Bytes())
}
