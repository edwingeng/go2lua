package unsupported

import (
	"fmt"
	"go/types"
	"sort"

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

	ordered := orderedTypes(make([]*typeWithPos, 0, len(pass.TypesInfo.Types)))
	for e := range pass.TypesInfo.Types {
		ordered = append(ordered, newTypeWithPos(pass, e))
	}
	sort.Sort(ordered)

	var total int
	m := make(map[int]struct{})
	var a []utils.NodeError
	for _, v := range ordered {
		e, x := v.typ, pass.TypesInfo.Types[v.typ]
		err := checkUnsupportedTypeImpl(x.Type)
		if err != nil {
			if _, ok := m[v.pos.Line]; ok {
				continue
			}
			m[v.pos.Line] = struct{}{}
			a = append(a, utils.NewNodeError(err, e))
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
	newError := func() error {
		t := typ.Underlying()
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
			return newError()
		} else if t.Kind() == types.UnsafePointer {
			return newError()
		}
	case *types.Array:
		const maxLen = 64
		if t.Len() > maxLen {
			return fmt.Errorf("the length of an array should not exceed %d", maxLen)
		}
	case *types.Slice:
	case *types.Struct:
	case *types.Pointer:
		return checkUnsupportedPointers(t)
	case *types.Tuple:
	case *types.Signature:
	case *types.Interface:
	case *types.Map:
		switch t.Key().Underlying().(type) {
		case nil, *types.Array, *types.Struct, *types.Tuple, *types.Interface:
			return newError()
		}
	case *types.Chan:
		return newError()
	case *types.Named:
		panic("impossible")
	default:
		return newError()
	}

	return nil
}

func checkUnsupportedPointers(typ *types.Pointer) error {
	newError := func() error {
		t := typ.Underlying()
		if typ.String() != t.String() {
			return fmt.Errorf("unsupported pointer type: %s. underlying: %s", typ.String(), t.String())
		} else {
			return fmt.Errorf("unsupported pointer type: %s", t.String())
		}
	}

	switch typ.Elem().Underlying().(type) {
	case nil:
	case *types.Basic:
		return newError()
	case *types.Array:
		return newError()
	case *types.Slice:
		return newError()
	case *types.Struct:
	case *types.Pointer:
		return newError()
	case *types.Tuple:
	case *types.Signature:
		return newError()
	case *types.Interface:
		return newError()
	case *types.Map:
		return newError()
	case *types.Chan:
		return newError()
	case *types.Named:
		return newError()
	default:
		return newError()
	}

	return nil
}
