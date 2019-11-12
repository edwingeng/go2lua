package main

import (
	"bytes"
	"errors"
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

	"github.com/edwingeng/go2lua/unsupported"
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
	rexSyntaxError = regexp.MustCompile(`(^[^:]+:\d+:\d+:)\s*(.*)`)
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
			allFiles = append(allFiles, filepath.Clean(f1))
		}
	}
	if len(allFiles) <= 1 {
		return ""
	}

	leN := 0
	file0 := allFiles[0]
outer:
	for leN < len(file0) {
		pos := strings.IndexRune(file0[leN:], filepath.Separator)
		if pos < 0 {
			break
		}

		str := file0[:leN+pos+1]
		for _, f := range allFiles {
			if !strings.HasPrefix(f, str) {
				break outer
			}
		}
		leN += pos + 1
	}

	return file0[:leN]
}

func (this *Parser) Parse() error {
	var err error
	cfg := packages.Config{Mode: math.MaxInt64}
	this.pkgs, err = packages.Load(&cfg, this.pkgPaths...)
	if err != nil {
		return err
	}

	for _, pkg := range this.pkgs {
		if len(pkg.Errors) > 0 {
			SyntaxErrorDetected = true
			for _, err := range pkg.Errors {
				a := rexSyntaxError.FindStringSubmatch(err.Error())
				if len(a) > 0 {
					fmt.Printf("%s\n    %s\n", a[1], a[2])
				} else {
					fmt.Println(err.Error())
				}
			}
		}
	}
	if SyntaxErrorDetected {
		fmt.Println()
		return errors.New("syntax error detected")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	for _, pkg := range this.pkgs {
		for _, syn := range pkg.Syntax {
			wg.Add(1)
			pkg, syn := pkg, syn
			go func() {
				defer wg.Done()
				pass := newPass(pkg, syn)
				runtime.Gosched()
				err := unsupported.CheckUnsupported(pass)
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
				}
			}()
			break
		}
	}
	wg.Wait()

	if firstErr != nil {
		return errors.New("unsupported feature(s) detected")
	}
	return nil
}

func newPass(pkg *packages.Package, syn *ast.File) *analysis.Pass {
	return &analysis.Pass{
		Analyzer:   shadow.Analyzer,
		Fset:       pkg.Fset,
		Files:      []*ast.File{syn},
		OtherFiles: pkg.OtherFiles,
		Pkg:        pkg.Types,
		TypesInfo:  pkg.TypesInfo,
		TypesSizes: pkg.TypesSizes,
	}
}

func findShadows(pkg *packages.Package, syn *ast.File) map[token.Pos]int {
	m := make(map[token.Pos]int)
	pass := newPass(pkg, syn)
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

	replaceSuffix := func(str, suffix, replacement string) string {
		if strings.HasSuffix(str, suffix) {
			str = str[:len(str)-len(suffix)] + replacement
		}
		return str
	}

	commonPrefix := this.commonPrefix()
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

				pass := newPass(item.pkg, item.syn)
				shadows := findShadows(item.pkg, item.syn)
				w := walker.NewWalker(pass, item.syn, walker.WithShadows(shadows))
				w.Walk()

				runtime.Gosched()
				f2 := filepath.Base(f1)
				f3 := replaceSuffix(f2, ".go", ".lua")
				f4 := filepath.Join(dir, f3)
				if err := ioutil.WriteFile(f4, w.Buffer.Bytes(), 0644); err != nil {
					panic(err)
				}
				fmt.Printf("- %s\n", replaceSuffix(f1[len(commonPrefix):], ".go", ".lua"))

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
				w := walker.NewWalker(newPass(pkg, syn), syn)
				w.Walk()

				fmt.Println("==========", f1[len(commonPrefix):])
				fmt.Println()
				fmt.Println(w.Buffer.String())

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
