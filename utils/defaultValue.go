package utils

import (
	"fmt"
	"go/types"
	"strings"
)

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
			panic("impossible")
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
		return "goslice.make(nil, 0)"

	case *types.Struct:
		if t.NumFields() == 0 {
			return "gostruct.empty"
		}

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
		return "{}"

	case *types.Chan:
		panic("impossible")

	case *types.Named:
		return DefaultValue(t.Underlying())

	default:
		panic("impossible")
	}
}
