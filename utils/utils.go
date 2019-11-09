package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type ErrItem struct {
	Err  error
	Node ast.Node
}

func NewErrItem(err error, node ast.Node) ErrItem {
	return ErrItem{
		Err:  err,
		Node: node,
	}
}

func PrintErrors(pass *analysis.Pass, a ...ErrItem) {
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
		_, _ = fmt.Fprintf(&buf, "%s [%d:%d]\n", pos.Filename, pos.Line, pos.Column)
		_, _ = fmt.Fprintf(&buf, "    %+v\n", strings.ReplaceAll(v.Err.Error(), "\n", "\n    "))
	}
	buf.WriteByte('\n')
	_, _ = os.Stderr.Write(buf.Bytes())
}

func DefaultValue(typ types.Type) string {
	switch t := typ.Underlying().(type) {
	case nil:
		return "undef"

	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return "false"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			return "0"
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return "0"
		case types.Uintptr, types.Float32, types.Float64, types.UnsafePointer:
			return "0"
		case types.String:
			return `""`
		default:
			panic("IMPOSSIBLE")
		}

	case *types.Array:
		var sb strings.Builder
		sb.WriteString("{")
		for i := int64(0); i < t.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(DefaultValue(t.Elem()))
		}
		sb.WriteString("}")
		return sb.String()

	case *types.Slice:
		return "{}"

	case *types.Struct:
		var sb strings.Builder
		sb.WriteString("{")
		for i := 0; i < t.NumFields(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			f := t.Field(i)
			_, _ = fmt.Fprintf(&sb, "%s = %s", f.Name(), DefaultValue(f.Type()))
		}
		sb.WriteString("}")
		return sb.String()

	case *types.Pointer:
		return "undef"

	case *types.Tuple:
		var sb strings.Builder
		sb.WriteString("{")
		for i := 0; i < t.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			f := t.At(i)
			_, _ = fmt.Fprintf(&sb, "%s = %s", f.Name(), DefaultValue(f.Type()))
		}
		sb.WriteString("}")
		return sb.String()

	case *types.Signature:
		return "undef"

	case *types.Interface:
		return "undef"

	case *types.Map:
		return "undef"

	case *types.Chan:
		panic("IMPOSSIBLE")

	case *types.Named:
		return DefaultValue(t.Underlying())

	default:
		panic("IMPOSSIBLE")
	}
}
