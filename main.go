package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	outputDir := flag.String("outputDir", "", "the output directory")
	astTree := flag.Bool("astTree", false, "print ast tree(s) for debug purpose")
	filter := flag.String("filter", "", "file filter, for debug purpose")
	pkgRoot := flag.String("pkgRoot", "", "the root of all your package paths")
	flag.Parse()

	if *outputDir != "" {
		if stat, err := os.Stat(*outputDir); os.IsNotExist(err) {
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
	p := NewParser(pkgPaths, WithFileFilter(fileFilter), WithPkgRoot(*pkgRoot))
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

	if p.ErrorOccurred {
		os.Exit(1)
	}
}
