package walker

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"
)

var (
	indentBytes = []byte("    ")
)

var (
	go2LuaFuncMap = map[string]string{
		"print":   "io.write",
		"println": "print",
	}

	go2LuaOperMap = map[string]string{
		`&&`: "and",
		`||`: "or",
		`!`:  "not",
	}
)

type gotoLabelInfo struct {
	funcNode ast.Node
	name     string
}

type Walker struct {
	Fset *token.FileSet
	root ast.Node

	fileData []byte
	current  ast.Node
	nextNums map[string]int
	indent   int

	buffer         bytes.Buffer
	FuncInit       bool
	ElseIfs        map[ast.Node]struct{}
	BreakLabels    map[ast.Node]string
	ContinueLabels map[ast.Node]string
	GotoLabels     map[gotoLabelInfo]struct{}
}

func NewWalker(fset *token.FileSet, node ast.Node) *Walker {
	w := &Walker{
		Fset:           fset,
		root:           node,
		nextNums:       make(map[string]int),
		ElseIfs:        make(map[ast.Node]struct{}),
		BreakLabels:    make(map[ast.Node]string),
		ContinueLabels: make(map[ast.Node]string),
		GotoLabels:     make(map[gotoLabelInfo]struct{}),
	}
	return w
}

func (this *Walker) printIndent() {
	bts := this.buffer.Bytes()
	if n := len(bts); n > 0 && bts[n-1] == '\n' {
		var c int
		pos := this.Fset.Position(this.current.Pos())
		for i := pos.Offset; i >= 0; {
			r, size := utf8.DecodeLastRune(this.fileData[:i])
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
					this.buffer.WriteByte('\n')
					break
				}
			}
		}
		for i := 0; i < this.indent; i++ {
			this.buffer.Write(indentBytes)
		}
	}
}

func (this *Walker) print(a ...interface{}) {
	this.printIndent()
	_, _ = fmt.Fprint(&this.buffer, a...)
}

func (this *Walker) println(a ...interface{}) {
	this.printIndent()
	_, _ = fmt.Fprintln(&this.buffer, a...)
}

func (this *Walker) printf(format string, a ...interface{}) {
	this.printIndent()
	_, _ = fmt.Fprintf(&this.buffer, format, a...)
}

func (this *Walker) printError(err error, node ast.Node) {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "%+v\n", err.Error())
	if err := ast.Fprint(&buf, this.Fset, node, nil); err != nil {
		panic(err)
	}
	_, _ = os.Stderr.Write(buf.Bytes())
}

func (this *Walker) makeUniqueName(key string) string {
	this.nextNums[key]++
	return fmt.Sprintf("xxx_%s_%d", key, this.nextNums[key])
}

func (this *Walker) isCallExpr_MakeMap(node ast.Node) bool {
	if n, ok := node.(*ast.CallExpr); ok {
		if funcExpr, ok := n.Fun.(*ast.Ident); ok {
			if funcExpr.Name == "make" {
				if len(n.Args) > 0 {
					if _, ok := n.Args[0].(*ast.MapType); ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func (this *Walker) initialize() {
	if this.root == nil {
		return
	}

	var stack []ast.Node
	var funcStack []ast.Node
	ast.Inspect(this.root, func(node ast.Node) bool {
		if node == nil {
			n := stack[len(stack)-1]
			switch n.(type) {
			case *ast.FuncLit, *ast.FuncDecl:
				funcStack = funcStack[:len(funcStack)-1]
			}
			stack = stack[:len(stack)-1]
			return true
		}

		stack = append(stack, node)
		switch n := node.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			funcStack = append(funcStack, node)
		case *ast.BranchStmt:
			switch {
			case n.Label == nil && n.Tok == token.CONTINUE:
				var loopNode ast.Node
				var loopNodeLabel string
				for i := len(stack) - 1; i >= 0 && loopNode == nil; i-- {
					switch sn := stack[i].(type) {
					case *ast.ForStmt, *ast.RangeStmt:
						loopNode = sn
						if stmt, ok := stack[i-1].(*ast.LabeledStmt); ok {
							loopNodeLabel = stmt.Label.Name
						}
					}
				}
				if loopNode == nil {
					this.printError(fmt.Errorf("unexpected token: %s", n.Tok), node)
					break
				}
				if _, ok := this.ContinueLabels[loopNode]; !ok {
					if loopNodeLabel == "" {
						this.ContinueLabels[loopNode] = this.makeUniqueName("continue")
					} else {
						this.ContinueLabels[loopNode] = loopNodeLabel
					}
				}
				loopNodeLabel = this.ContinueLabels[loopNode]
				this.ContinueLabels[node] = loopNodeLabel

			case n.Label != nil && n.Label.Name != "" && (n.Tok == token.BREAK || n.Tok == token.CONTINUE):
				var loopNode ast.Node
				var err error
				for i := len(stack) - 1; i >= 0 && loopNode == nil; i-- {
					switch sn := stack[i].(type) {
					case *ast.LabeledStmt:
						if sn.Label.Name == n.Label.Name {
							switch sn.Stmt.(type) {
							case *ast.ForStmt, *ast.RangeStmt:
								loopNode = sn.Stmt
							default:
								err = fmt.Errorf("%q is NOT a 'for' label", sn.Label.Name)
							}
						}
					}
				}
				if err != nil {
					this.printError(err, n)
					break
				}
				if loopNode == nil {
					this.printError(fmt.Errorf("unexpected token: %s", n.Tok), node)
					break
				}
				switch n.Tok {
				case token.BREAK:
					if _, ok := this.BreakLabels[loopNode]; !ok {
						this.BreakLabels[loopNode] = n.Label.Name + "_break"
					}
				case token.CONTINUE:
					if _, ok := this.ContinueLabels[loopNode]; !ok {
						this.ContinueLabels[loopNode] = n.Label.Name + "_continue"
					}
				default:
					panic("IMPOSSIBLE")
				}

			case n.Label != nil && n.Label.Name != "" && n.Tok == token.GOTO:
				key := gotoLabelInfo{
					funcNode: funcStack[len(funcStack)-1],
					name:     n.Label.Name,
				}
				this.GotoLabels[key] = struct{}{}
			}
		}

		return true
	})

	var err error
	file := this.Fset.File(this.root.Pos()).Name()
	this.fileData, err = ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
}

func (this *Walker) forStmt_Continues(node *ast.ForStmt) (found, immediate bool, labels []string) {
	m := make(map[string]struct{})
	ast.Inspect(node, func(node ast.Node) bool {
		if n, ok := node.(*ast.BranchStmt); ok {
			if n.Tok == token.CONTINUE {
				found = true

				m[n.Label.Name] = struct{}{}
			}
		}
		return true
	})
	for k := range m {
		labels = append(labels, k)
	}
	return
}

func (this *Walker) walkIdentList(list []*ast.Ident, funcNode ast.Node) {
	for i, x := range list {
		if i > 0 {
			this.print(", ")
		}
		this.walkImpl(x, funcNode)
	}
}

func (this *Walker) walkExprList(list []ast.Expr, funcNode ast.Node) {
	for i, x := range list {
		if i > 0 {
			this.print(", ")
		}
		if this.isCallExpr_MakeMap(x) {
			this.print("{}")
		} else {
			this.walkImpl(x, funcNode)
		}
	}
}

func (this *Walker) walkStmtList(list []ast.Stmt, newline bool, funcNode ast.Node) {
	for _, x := range list {
		this.walkImpl(x, funcNode)
		if newline {
			this.println()
		}
	}
}

func (this *Walker) walkDeclList(list []ast.Decl, funcNode ast.Node) {
	for _, x := range list {
		this.walkImpl(x, funcNode)
	}
}

func (this *Walker) Walk() {
	if this.buffer.Len() > 0 {
		return
	}

	this.initialize()
	this.walkImpl(this.root, nil)
	this.trim()
}

func (this *Walker) walkImpl(node ast.Node, funcNode ast.Node) {
	this.current = node
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		for _, c := range n.List {
			this.walkImpl(c, funcNode)
		}

	case *ast.Field:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}
		this.walkIdentList(n.Names, funcNode)
		if n.Tag != nil {
			this.walkImpl(n.Tag, funcNode)
		}
		if n.Comment != nil {
			this.walkImpl(n.Comment, funcNode)
		}

	case *ast.FieldList:
		for i, f := range n.List {
			if i > 0 {
				this.print(", ")
			}
			this.walkImpl(f, funcNode)
		}

	// Expressions
	case *ast.BadExpr:
		this.printError(errors.New("bad expression detected"), n)

	case *ast.Ident:
		if str, ok := go2LuaFuncMap[n.Name]; ok {
			this.print(str)
		} else {
			this.print(n.Name)
		}

	case *ast.BasicLit:
		this.print(n.Value)

	case *ast.Ellipsis:
		if n.Elt != nil {
			this.walkImpl(n.Elt, funcNode)
		}

	case *ast.FuncLit:
		this.walkImpl(n.Type, n)
		this.walkImpl(n.Body, n)

	case *ast.CompositeLit:
		if n.Type != nil {
			this.walkImpl(n.Type, funcNode)
		}
		this.walkExprList(n.Elts, funcNode)

	case *ast.ParenExpr:
		this.walkImpl(n.X, funcNode)

	case *ast.SelectorExpr:
		this.walkImpl(n.X, funcNode)
		this.walkImpl(n.Sel, funcNode)

	case *ast.IndexExpr:
		this.walkImpl(n.X, funcNode)
		this.print("[")
		this.walkImpl(n.Index, funcNode)
		this.print("]")

	case *ast.SliceExpr:
		this.walkImpl(n.X, funcNode)
		if n.Low != nil {
			this.walkImpl(n.Low, funcNode)
		}
		if n.High != nil {
			this.walkImpl(n.High, funcNode)
		}
		if n.Max != nil {
			this.walkImpl(n.Max, funcNode)
		}

	case *ast.TypeAssertExpr:
		this.walkImpl(n.X, funcNode)
		if n.Type != nil {
			this.walkImpl(n.Type, funcNode)
		}

	case *ast.CallExpr:
		this.walkImpl(n.Fun, funcNode)

		this.print("(")
		this.walkExprList(n.Args, funcNode)
		this.print(")")

	case *ast.StarExpr:
		this.walkImpl(n.X, funcNode)

	case *ast.UnaryExpr:
		this.walkImpl(n.X, funcNode)

	case *ast.BinaryExpr:
		this.walkImpl(n.X, funcNode)
		if str, ok := go2LuaOperMap[n.Op.String()]; ok {
			this.printf(" %s ", str)
		} else {
			this.printf(" %s ", n.Op)
		}
		this.walkImpl(n.Y, funcNode)

	case *ast.KeyValueExpr:
		this.walkImpl(n.Key, funcNode)
		this.walkImpl(n.Value, funcNode)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			this.walkImpl(n.Len, funcNode)
		}
		this.walkImpl(n.Elt, funcNode)

	case *ast.StructType:
		this.walkImpl(n.Fields, funcNode)

	case *ast.FuncType:
		if n.Params != nil {
			this.walkImpl(n.Params, funcNode)
		}

	case *ast.InterfaceType:
		this.walkImpl(n.Methods, funcNode)

	case *ast.MapType:
		this.walkImpl(n.Key, funcNode)
		this.walkImpl(n.Value, funcNode)

	case *ast.ChanType:
		this.walkImpl(n.Value, funcNode)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		this.walkImpl(n.Decl, funcNode)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		key := gotoLabelInfo{
			funcNode: funcNode,
			name:     n.Label.Name,
		}
		if _, ok := this.GotoLabels[key]; ok {
			this.indent--
			this.printf("::%s::\n", n.Label.Name)
			this.indent++
		} else {
			this.print()
		}
		this.walkImpl(n.Stmt, funcNode)

	case *ast.ExprStmt:
		this.walkImpl(n.X, funcNode)

	case *ast.SendStmt:
		this.walkImpl(n.Chan, funcNode)
		this.walkImpl(n.Value, funcNode)

	case *ast.IncDecStmt:
		this.walkImpl(n.X, funcNode)
		switch n.Tok {
		case token.INC:
			this.print(" = ")
			this.walkImpl(n.X, funcNode)
			this.print(" + 1")
		case token.DEC:
			this.print(" = ")
			this.walkImpl(n.X, funcNode)
			this.print(" - 1")
		default:
			panic(fmt.Errorf("unexpected token: %s", n.Tok))
		}

	case *ast.AssignStmt:
		if n.Tok == token.DEFINE {
			this.print("local ")
		}
		this.walkExprList(n.Lhs, funcNode)
		this.print(" = ")
		this.walkExprList(n.Rhs, funcNode)

	case *ast.GoStmt:
		this.walkImpl(n.Call, funcNode)

	case *ast.DeferStmt:
		this.walkImpl(n.Call, funcNode)

	case *ast.ReturnStmt:
		this.print("return ")
		this.walkExprList(n.Results, funcNode)

	case *ast.BranchStmt:
		if n.Label == nil {
			switch n.Tok {
			case token.BREAK:
				this.print("break")
			case token.CONTINUE:
				this.printf("goto %s", this.ContinueLabels[n])
			case token.GOTO:
				this.printError(errors.New("missing label"), node)
			}
		} else {
			switch n.Tok {
			case token.BREAK:
				this.printf("goto %s_break", n.Label)
			case token.CONTINUE:
				this.printf("goto %s_continue", n.Label)
			case token.GOTO:
				this.printf("goto %s", n.Label)
			}
		}

	case *ast.BlockStmt:
		this.indent++
		this.walkStmtList(n.List, true, funcNode)
		this.indent--

	case *ast.IfStmt:
		if n.Init != nil {
			this.walkImpl(n.Init, funcNode)
		}
		this.print("if ")
		this.walkImpl(n.Cond, funcNode)

		this.println(" then")
		this.walkImpl(n.Body, funcNode)
		var elif ast.Node
		if n.Else != nil {
			if nn, ok := n.Else.(*ast.IfStmt); ok {
				this.print("else")
				this.ElseIfs[nn] = struct{}{}
				elif = nn
			} else {
				this.println("else")
			}
			this.walkImpl(n.Else, funcNode)
		}
		if _, ok := this.ElseIfs[n]; !ok {
			this.print("end")
		}
		if elif != nil {
			delete(this.ElseIfs, elif)
		}

	case *ast.CaseClause:
		this.walkExprList(n.List, funcNode)
		this.walkStmtList(n.Body, false, funcNode)

	case *ast.SwitchStmt:
		if n.Init != nil {
			this.walkImpl(n.Init, funcNode)
		}
		if n.Tag != nil {
			this.walkImpl(n.Tag, funcNode)
		}
		this.walkImpl(n.Body, funcNode)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			this.walkImpl(n.Init, funcNode)
		}
		this.walkImpl(n.Assign, funcNode)
		this.walkImpl(n.Body, funcNode)

	case *ast.CommClause:
		if n.Comm != nil {
			this.walkImpl(n.Comm, funcNode)
		}
		this.walkStmtList(n.Body, false, funcNode)

	case *ast.SelectStmt:
		this.walkImpl(n.Body, funcNode)

	case *ast.ForStmt:
		if n.Init != nil {
			this.println("do")
			this.indent++
			this.walkImpl(n.Init, funcNode)
			this.println()
		}
		if n.Cond != nil {
			this.print("while ")
			this.walkImpl(n.Cond, funcNode)
			this.println(" do")
		} else {
			this.println("while true do")
		}

		if n.Post != nil && n.Body != nil && len(n.Body.List) > 0 {
			this.indent++
			this.println("do")
		}
		this.walkImpl(n.Body, funcNode)
		if n.Post != nil && n.Body != nil && len(n.Body.List) > 0 {
			this.println("end")
			this.indent--
		}

		if label, ok := this.ContinueLabels[n]; ok {
			this.printf("::%s::\n", label)
		}
		if n.Post != nil {
			this.indent++
			this.walkImpl(n.Post, funcNode)
			this.indent--
			this.println()
		}

		if n.Init != nil {
			this.println("end")
			this.indent--
		}

		if label, ok := this.BreakLabels[n]; ok {
			this.println("end")
			this.indent--
			this.printf("::%s::", label)
			this.indent++
		} else {
			this.print("end")
		}

	case *ast.RangeStmt:
		this.print("for ")
		if n.Key != nil {
			this.walkImpl(n.Key, funcNode)
			this.print(", ")
		} else {
			this.print("_, ")
		}
		if n.Value != nil {
			this.walkImpl(n.Value, funcNode)
		}

		this.print(" in pairs(")
		this.walkImpl(n.X, funcNode)
		this.println(") do")

		this.walkImpl(n.Body, funcNode)

		if label, ok := this.ContinueLabels[n]; ok {
			this.printf("::%s::\n", label)
		}
		if label, ok := this.BreakLabels[n]; ok {
			this.println("end")
			this.indent--
			this.printf("::%s::", label)
			this.indent++
		} else {
			this.print("end")
		}

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}
		if n.Name != nil {
			this.walkImpl(n.Name, funcNode)
		}
		this.walkImpl(n.Path, funcNode)
		if n.Comment != nil {
			this.walkImpl(n.Comment, funcNode)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}
		this.walkIdentList(n.Names, funcNode)
		if n.Type != nil {
			this.walkImpl(n.Type, funcNode)
		}
		this.walkExprList(n.Values, funcNode)
		if n.Comment != nil {
			this.walkImpl(n.Comment, funcNode)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}
		this.walkImpl(n.Name, funcNode)
		this.walkImpl(n.Type, funcNode)
		if n.Comment != nil {
			this.walkImpl(n.Comment, funcNode)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}
		for _, s := range n.Specs {
			this.walkImpl(s, funcNode)
		}

	case *ast.FuncDecl:
		if n.Name.Name == "init" && n.Recv == nil {
			this.FuncInit = true
		}

		this.print()
		if n.Doc != nil {
			this.walkImpl(n.Doc, n)
		}
		if n.Recv != nil {
			this.walkImpl(n.Recv, n)
		}

		first := n.Name.Name[0]
		if first >= 'a' && first <= 'z' {
			this.print("local ")
		}

		this.walkImpl(n.Name, n)
		this.print(" = function(")

		this.walkImpl(n.Type, n)
		this.println(")")

		if n.Body != nil {
			this.walkImpl(n.Body, n)
		}
		this.println("end")

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}

		this.print("-- package: ")
		this.walkImpl(n.Name, funcNode)
		this.println()

		this.walkDeclList(n.Decls, funcNode)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

		if this.FuncInit {
			this.println()
			this.println("return init")
		}

	case *ast.Package:
		for _, f := range n.Files {
			this.walkImpl(f, funcNode)
		}

	default:
		panic(fmt.Sprintf("ast.walkImpl: unexpected node type %T", n))
	}
}

func (this *Walker) trim() {
	bts := this.buffer.Bytes()
	bts = bytes.TrimRight(bts, "\n")
	if n := len(bts); n < this.buffer.Len() {
		this.buffer.Truncate(n + 1)
	}
}

func (this *Walker) BufferString() string {
	return this.buffer.String()
}

func (this *Walker) BufferBytes() []byte {
	return this.buffer.Bytes()
}
