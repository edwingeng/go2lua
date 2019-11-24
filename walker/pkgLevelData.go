package walker

import (
	"go/token"
	"sync"
)

type PkgLevelData struct {
	sync.Mutex
	Vars map[token.Pos]string
}

func NewPkgLevelData() *PkgLevelData {
	return &PkgLevelData{
		Vars: make(map[token.Pos]string),
	}
}
