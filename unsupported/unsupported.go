package unsupported

import (
	"fmt"
	"go/types"

	"github.com/edwingeng/go2lua/utils"
	"golang.org/x/tools/go/analysis"
)

func CheckUnsupported(pass *analysis.Pass) error {
	if pass == nil {
		return nil
	}

	return checkUnsupportedTypes(pass)
}

func checkUnsupportedTypes(pass *analysis.Pass) error {
	if pass.TypesInfo == nil {
		return nil
	}

	var total int
	m := make(map[int]struct{})
	var a []utils.ErrItem
	for e, tnv := range pass.TypesInfo.Types {
		err := checkUnsupportedTypeImpl(tnv.Type)
		if err != nil {
			pos := pass.Fset.Position(e.Pos())
			if _, ok := m[pos.Line]; ok {
				continue
			}
			m[pos.Line] = struct{}{}
			a = append(a, utils.NewErrItem(err, e))
			if total++; total > 10 {
				break
			}
		}
	}

	if len(a) > 0 {
		utils.PrintErrors(pass, a...)
		return a[0].Err
	}
	return nil
}

func checkUnsupportedTypeImpl(typ types.Type) error {
	newError := func(t types.Type) error {
		if typ.String() != t.String() {
			return fmt.Errorf("unsupported data type: %s. underlying: %s", typ.String(), t.String())
		} else {
			return fmt.Errorf("unsupported data type: %s", t.String())
		}
	}

	switch t := typ.Underlying().(type) {
	case nil:
	case *types.Basic:
		if t.Info()&types.IsComplex != 0 {
			return newError(t)
		} else if t.Kind() == types.UnsafePointer {
			return newError(t)
		}
	case *types.Array:
		const maxLen = 64
		if t.Len() > maxLen {
			return fmt.Errorf("the length of an array should not exceed %d", maxLen)
		}
	case *types.Slice:
	case *types.Struct:
	case *types.Pointer:
	case *types.Tuple:
	case *types.Signature:
	case *types.Interface:
	case *types.Map:
	case *types.Chan:
		return newError(t)
	case *types.Named:
	default:
		return newError(t)
	}

	return nil
}
