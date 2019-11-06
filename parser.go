package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/edwingeng/go2lua/walker"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var (
	SyntaxErrorDetected bool
)

var (
	rexShadow      = regexp.MustCompile(`shadows declaration at line (\d+)`)
	shadowFlagOnce sync.Once
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

func findShadows(pkg *packages.Package, syn *ast.File) map[token.Pos]int {
	m := make(map[token.Pos]int)
	pass := &analysis.Pass{
		Analyzer:   shadow.Analyzer,
		Fset:       pkg.Fset,
		Files:      []*ast.File{syn},
		OtherFiles: pkg.OtherFiles,
		Pkg:        pkg.Types,
		TypesInfo:  pkg.TypesInfo,
		TypesSizes: pkg.TypesSizes,
	}
	shadowFlagOnce.Do(func() {
		if err := shadow.Analyzer.Flags.Parse([]string{"-strict"}); err != nil {
			panic(err)
		}
	})
	pass.ResultOf = map[*analysis.Analyzer]interface{}{
		inspect.Analyzer: inspector.New(pass.Files),
	}
	pass.Report = func(d analysis.Diagnostic) {
		matches := rexShadow.FindStringSubmatch(d.Message)
		if len(matches) > 0 {
			n, _ := strconv.Atoi(matches[1])
			m[d.Pos] = n
		}
	}
	_, _ = pass.Analyzer.Run(pass)
	return m
}

func (this *Parser) Output(dir string) {
	type item struct {
		pkg *packages.Package
		syn *ast.File
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
					pkg: pkg,
					syn: syn,
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

				f1 := item.pkg.Fset.Position(item.syn.Package).Filename
				if this.fileFilter != nil && !this.fileFilter(f1) {
					continue
				}

				shadows := findShadows(item.pkg, item.syn)
				w := walker.NewWalker(item.pkg.Fset, item.syn, walker.WithShadows(shadows))
				w.Walk()

				runtime.Gosched()
				f2 := filepath.Base(f1)
				f3 := strings.TrimSuffix(f2, ".go") + ".lua"
				f4 := filepath.Join(dir, f3)
				if err := ioutil.WriteFile(f4, w.BufferBytes(), 0644); err != nil {
					panic(err)
				}

				if !SyntaxErrorDetected {
					SyntaxErrorDetected = w.NumErrors > 0
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

				if !SyntaxErrorDetected {
					SyntaxErrorDetected = w.NumErrors > 0
				}
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
