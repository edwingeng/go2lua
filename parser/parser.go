package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/edwingeng/go2lua/walker"
	"golang.org/x/tools/go/packages"
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
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if this.fileFilter != nil && !this.fileFilter(f1) {
				continue
			}

			w := walker.NewWalker(pkg.Fset)
			w.Initialize(syn)
			w.Walk(syn)
			w.Trim()

			f2 := filepath.Base(f1)
			f3 := strings.TrimSuffix(f2, ".go") + ".lua"
			f4 := filepath.Join(dir, f3)
			if err := ioutil.WriteFile(f4, w.Bytes(), 0644); err != nil {
				panic(err)
			}
		}
	}
}

func (this *Parser) PrintDetails(debugMode bool) {
	commonPrefix := this.commonPrefix()
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if this.fileFilter != nil && !this.fileFilter(f1) {
				continue
			}

			if debugMode {
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

			w := walker.NewWalker(pkg.Fset)
			w.Initialize(syn)
			w.Walk(syn)
			w.Trim()

			fmt.Println("==========", f1[len(commonPrefix):])
			fmt.Println()
			fmt.Println(w.String())
		}
	}
}

type Option func(p *Parser)

func WithFileFilter(f func(file string) bool) Option {
	return func(p *Parser) {
		p.fileFilter = f
	}
}
