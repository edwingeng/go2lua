package parser

import (
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

func (this *Parser) commonPrefixLen() int {
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

	return leN
}

func (this *Parser) Parse() error {
	var err error
	cfg := packages.Config{Mode: math.MaxInt64}
	this.pkgs, err = packages.Load(&cfg, this.pkgPaths...)
	return err
}

func (this *Parser) Output(dir string) error {
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if this.fileFilter != nil && !this.fileFilter(f1) {
				continue
			}

			w := walker.NewWalker(pkg.Fset)
			w.Walk(syn)

			f2 := filepath.Base(f1)
			f3 := strings.TrimSuffix(f2, ".go") + ".lua"
			f4 := filepath.Join(dir, f3)
			if err := ioutil.WriteFile(f4, w.Bytes(), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Parser) PrintDetails(debugMode bool) {
	commonPrefixLen := this.commonPrefixLen()
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if this.fileFilter != nil && !this.fileFilter(f1) {
				continue
			}

			if debugMode {
				fmt.Println("=======###", f1[commonPrefixLen:])
				fmt.Println()
				if err := ast.Print(pkg.Fset, syn); err != nil {
					panic(err)
				}
				fmt.Println()
			}

			w := walker.NewWalker(pkg.Fset)
			w.Walk(syn)

			fmt.Println("==========", f1[commonPrefixLen:])
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
