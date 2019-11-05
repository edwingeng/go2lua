package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
	"sync"

	"github.com/edwingeng/go2lua/walker"
	"golang.org/x/tools/go/packages"
)

var (
	TotalErrors int
)

type Parser struct {
	pkgPaths   []string
	fileFilter func(file string) bool

	pkgs []*packages.Package
}

func NewParser(pkgPaths []string, opts ...Option) *Parser {
	p := &Parser{
		pkgPaths: pkgPaths,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (this *Parser) commonPrefix() string {
	var allFiles []string
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			allFiles = append(allFiles, f1)
		}
	}

	var leN int
	for done := false; len(allFiles) > 1 && !done; {
		a := strings.Split(allFiles[0], string(filepath.Separator))
		for i := 1; i < len(a); i++ {
			prefix := filepath.Join(a[:i]...)
			for _, f := range allFiles {
				if !strings.HasPrefix(f, prefix) {
					done = true
					break
				}
			}
			leN = len(prefix) + 1
			if a[0] == "" {
				leN++
			}
		}
	}

	if leN > 0 {
		return allFiles[0][:leN]
	}
	return ""
}

func (this *Parser) Parse() error {
	var err error
	cfg := packages.Config{Mode: math.MaxInt64}
	this.pkgs, err = packages.Load(&cfg, this.pkgPaths...)
	return err
}

func (this *Parser) Output(dir string) {
	type item struct {
		fset *token.FileSet
		file *ast.File
	}

	ch := make(chan item)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			close(ch)
			wg.Done()
		}()
		for _, pkg := range this.pkgs {
			for _, syn := range pkg.Syntax {
				ch <- item{
					fset: pkg.Fset,
					file: syn,
				}
			}
		}
	}()
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				item, ok := <-ch
				if !ok {
					break
				}

				f1 := item.fset.Position(item.file.Package).Filename
				if this.fileFilter != nil && !this.fileFilter(f1) {
					continue
				}

				w := walker.NewWalker(item.fset, item.file)
				w.Walk()
				TotalErrors += w.NumErrors

				f2 := filepath.Base(f1)
				f3 := strings.TrimSuffix(f2, ".go") + ".lua"
				f4 := filepath.Join(dir, f3)
				if err := ioutil.WriteFile(f4, w.BufferBytes(), 0644); err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()
}

func (this *Parser) PrintDetails(astTree, luaCode bool) {
	commonPrefix := this.commonPrefix()
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if this.fileFilter != nil && !this.fileFilter(f1) {
				continue
			}

			if astTree {
				fmt.Println("=======###", f1[len(commonPrefix):])
				fmt.Println()
				var buf bytes.Buffer
				if err := ast.Fprint(&buf, pkg.Fset, syn, nil); err != nil {
					panic(err)
				}
				str := buf.String()
				if len(commonPrefix) > 0 {
					str = strings.ReplaceAll(str, commonPrefix, "")
				}
				fmt.Println(str)
			}

			if luaCode {
				w := walker.NewWalker(pkg.Fset, syn)
				w.Walk()

				fmt.Println("==========", f1[len(commonPrefix):])
				fmt.Println()
				fmt.Println(w.BufferString())
			}
		}
	}
}

type Option func(p *Parser)

func WithFileFilter(f func(file string) bool) Option {
	return func(p *Parser) {
		p.fileFilter = f
	}
}
