package codeprinter

import (
	"bytes"
	"fmt"
	"go/ast"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

type Printer struct {
	bytes.Buffer
	Pass        *analysis.Pass
	srcFile     []byte
	indentBytes []byte

	Indent      int
	CurrentNode ast.Node
}

func NewPrinter(pass *analysis.Pass, srcFile []byte, opts ...Option) (p *Printer) {
	p = &Printer{
		Pass:        pass,
		srcFile:     srcFile,
		indentBytes: []byte("    "),
	}
	for _, opt := range opts {
		opt(p)
	}
	return
}

func (this *Printer) printIndent(blank bool) {
	bts := this.Buffer.Bytes()
	n := len(bts)
	if n == 0 {
		for i := 0; i < this.Indent; i++ {
			this.Buffer.Write(this.indentBytes)
		}
		return
	}
	if bts[n-1] != '\n' {
		return
	}

	var c int
	pos := this.Pass.Fset.Position(this.CurrentNode.Pos())
	for i := pos.Offset; i >= 0; {
		r, size := utf8.DecodeLastRune(this.srcFile[:i])
		i -= size
		if r == '\n' {
			if c++; c >= 2 {
				break
			}
		} else if !unicode.IsSpace(r) {
			break
		}
	}
	if c >= 2 {
		for i := n - 1; i >= 0; {
			r, size := utf8.DecodeLastRune(bts[:i])
			i -= size
			if r == '\n' {
				break
			} else if !unicode.IsSpace(r) {
				if !blank {
					this.Buffer.WriteByte('\n')
				}
				break
			}
		}
	}

	for i := 0; i < this.Indent; i++ {
		this.Buffer.Write(this.indentBytes)
	}
}

func (this *Printer) Print(a ...interface{}) {
	this.printIndent(false)
	_, _ = fmt.Fprint(&this.Buffer, a...)
}

func (this *Printer) Println(a ...interface{}) {
	this.printIndent(len(a) == 0)
	_, _ = fmt.Fprintln(&this.Buffer, a...)
}

func (this *Printer) Printf(format string, a ...interface{}) {
	this.printIndent(false)
	_, _ = fmt.Fprintf(&this.Buffer, format, a...)
}

type Option func(p *Printer)

func WithIndentString(str string) Option {
	return func(p *Printer) {
		p.indentBytes = []byte(str)
	}
}
