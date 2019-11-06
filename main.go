package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	astTree := flag.Bool("astTree", false, "print ast trees for debug purpose")
	outputDir := flag.String("outputDir", "", "the output directory")
	filter := flag.String("filter", "", "file filter, for debug purpose")
	flag.Parse()

	if *outputDir != "" {
		if stat, err := os.Stat(*outputDir); err != nil {
			panic(err)
		} else if !stat.IsDir() {
			panic(fmt.Errorf("%s is not a directory", *outputDir))
		}
	}

	fileFilter := func(file string) bool {
		if *filter != "" {
			matched, err := filepath.Match(*filter, filepath.Base(file))
			if err != nil || !matched {
				return false
			}
		}
		return true
	}

	var pkgPaths []string
	for i := 0; i < flag.NArg(); i++ {
		pkgPaths = append(pkgPaths, flag.Arg(i))
	}
	p := NewParser(pkgPaths, WithFileFilter(fileFilter))
	if err := p.Parse(); err != nil {
		panic(err)
	}

	if *outputDir != "" {
		p.Output(*outputDir)
		if *astTree {
			p.PrintDetails(true, false)
		}
	} else {
		p.PrintDetails(*astTree, true)
	}

	if SyntaxErrorDetected {
		os.Exit(1)
	}
}
