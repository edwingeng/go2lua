package main

import (
	"flag"
	"fmt"
	"go/ast"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/edwingeng/go2lua/walker"
	"golang.org/x/tools/go/packages"
)

func main() {
	astTree := flag.Bool("astTree", false, "print ast trees for debug purpose")
	outputDir := flag.String("outputDir", "", "the output directory")
	flag.Parse()

	if *outputDir != "" {
		if stat, err := os.Stat(*outputDir); err != nil {
			panic(err)
		} else if !stat.IsDir() {
			panic(fmt.Errorf("%s is not a directory", *outputDir))
		}
	}

	var pkgPaths []string
	for i := 0; i < flag.NArg(); i++ {
		pkgPaths = append(pkgPaths, flag.Arg(i))
	}
	cfg := packages.Config{Mode: math.MaxInt64}
	pkgs, err := packages.Load(&cfg, pkgPaths...)
	if err != nil {
		panic(err)
	}

	var allFiles []string
	for _, pkg := range pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			allFiles = append(allFiles, f1)
		}
	}
	var commonPrefixLen int
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
			commonPrefixLen = len(prefix) + 1
			if a[0] == "" {
				commonPrefixLen++
			}
		}
	}

	for _, pkg := range pkgs {
		for _, syn := range pkg.Syntax {
			f1 := pkg.Fset.Position(syn.Package).Filename
			if *astTree {
				fmt.Println("=======###", f1[commonPrefixLen:])
				fmt.Println()
				if err := ast.Print(pkg.Fset, syn); err != nil {
					panic(err)
				}
				fmt.Println()
			}

			w := walker.NewWalker(pkg.Fset)
			w.Walk(syn)

			if *outputDir != "" {
				f2 := filepath.Base(f1)
				f3 := strings.TrimSuffix(f2, ".go") + ".lua"
				f4 := filepath.Join(*outputDir, f3)
				if err := ioutil.WriteFile(f4, w.Bytes(), 0644); err != nil {
					panic(err)
				}
			} else {
				fmt.Println("==========", f1[commonPrefixLen:])
				fmt.Println()
				fmt.Println(w.String())
			}
		}
	}
}
