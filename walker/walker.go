package walker

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"os"
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

type Walker struct {
	Fset *token.FileSet
	root ast.Node

	nextNums map[string]int
	indent   int

	buffer         bytes.Buffer
	FuncInit       bool
	ElseIfs        map[ast.Node]struct{}
	BreakLabels    map[ast.Node]string
	ContinueLabels map[ast.Node]string
}

func NewWalker(fset *token.FileSet, node ast.Node) *Walker {
	w := &Walker{
		Fset:           fset,
		root:           node,
		nextNums:       make(map[string]int),
		ElseIfs:        make(map[ast.Node]struct{}),
		BreakLabels:    make(map[ast.Node]string),
		ContinueLabels: make(map[ast.Node]string),
	}
	return w
}

func (this *Walker) printIndent() {
	bts := this.buffer.Bytes()
	if n := len(bts); n > 0 && bts[n-1] == '\n' {
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
	_, _ = fmt.Fprintln(&buf, err.Error()+".")
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
	ast.Inspect(this.root, func(node ast.Node) bool {
		if node == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, node)
		switch n := node.(type) {
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

			case n.Label != nil && (n.Tok == token.BREAK || n.Tok == token.CONTINUE):
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
			}
		}

		return true
	})
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

func (this *Walker) walkIdentList(list []*ast.Ident) {
	for i, x := range list {
		if i > 0 {
			this.print(", ")
		}
		this.walkImpl(x)
	}
}

func (this *Walker) walkExprList(list []ast.Expr) {
	for i, x := range list {
		if i > 0 {
			this.print(", ")
		}
		if this.isCallExpr_MakeMap(x) {
			this.print("{}")
		} else {
			this.walkImpl(x)
		}
	}
}

func (this *Walker) walkStmtList(list []ast.Stmt, sep string) {
	for _, x := range list {
		this.walkImpl(x)
		this.print(sep)
	}
}

func (this *Walker) walkDeclList(list []ast.Decl) {
	for _, x := range list {
		this.walkImpl(x)
	}
}

func (this *Walker) Walk() {
	if this.buffer.Len() > 0 {
		return
	}

	this.initialize()
	this.walkImpl(this.root)
	this.trim()
}

func (this *Walker) walkImpl(node ast.Node) {
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		for _, c := range n.List {
			this.walkImpl(c)
		}

	case *ast.Field:
		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}
		this.walkIdentList(n.Names)
		if n.Tag != nil {
			this.walkImpl(n.Tag)
		}
		if n.Comment != nil {
			this.walkImpl(n.Comment)
		}

	case *ast.FieldList:
		for i, f := range n.List {
			if i > 0 {
				this.print(", ")
			}
			this.walkImpl(f)
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
			this.walkImpl(n.Elt)
		}

	case *ast.FuncLit:
		this.walkImpl(n.Type)
		this.walkImpl(n.Body)

	case *ast.CompositeLit:
		if n.Type != nil {
			this.walkImpl(n.Type)
		}
		this.walkExprList(n.Elts)

	case *ast.ParenExpr:
		this.walkImpl(n.X)

	case *ast.SelectorExpr:
		this.walkImpl(n.X)
		this.walkImpl(n.Sel)

	case *ast.IndexExpr:
		this.walkImpl(n.X)
		this.print("[")
		this.walkImpl(n.Index)
		this.print("]")

	case *ast.SliceExpr:
		this.walkImpl(n.X)
		if n.Low != nil {
			this.walkImpl(n.Low)
		}
		if n.High != nil {
			this.walkImpl(n.High)
		}
		if n.Max != nil {
			this.walkImpl(n.Max)
		}

	case *ast.TypeAssertExpr:
		this.walkImpl(n.X)
		if n.Type != nil {
			this.walkImpl(n.Type)
		}

	case *ast.CallExpr:
		this.walkImpl(n.Fun)

		this.print("(")
		this.walkExprList(n.Args)
		this.print(")")

	case *ast.StarExpr:
		this.walkImpl(n.X)

	case *ast.UnaryExpr:
		this.walkImpl(n.X)

	case *ast.BinaryExpr:
		this.walkImpl(n.X)
		if str, ok := go2LuaOperMap[n.Op.String()]; ok {
			this.printf(" %s ", str)
		} else {
			this.printf(" %s ", n.Op)
		}
		this.walkImpl(n.Y)

	case *ast.KeyValueExpr:
		this.walkImpl(n.Key)
		this.walkImpl(n.Value)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			this.walkImpl(n.Len)
		}
		this.walkImpl(n.Elt)

	case *ast.StructType:
		this.walkImpl(n.Fields)

	case *ast.FuncType:
		if n.Params != nil {
			this.walkImpl(n.Params)
		}

	case *ast.InterfaceType:
		this.walkImpl(n.Methods)

	case *ast.MapType:
		this.walkImpl(n.Key)
		this.walkImpl(n.Value)

	case *ast.ChanType:
		this.walkImpl(n.Value)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		this.walkImpl(n.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		this.indent--
		this.printf("::%s::\n", n.Label.Name)
		this.indent++
		this.walkImpl(n.Stmt)

	case *ast.ExprStmt:
		this.walkImpl(n.X)

	case *ast.SendStmt:
		this.walkImpl(n.Chan)
		this.walkImpl(n.Value)

	case *ast.IncDecStmt:
		this.walkImpl(n.X)
		switch n.Tok {
		case token.INC:
			this.print(" = ")
			this.walkImpl(n.X)
			this.print(" + 1")
		case token.DEC:
			this.print(" = ")
			this.walkImpl(n.X)
			this.print(" - 1")
		default:
			panic(fmt.Errorf("unexpected token: %s", n.Tok))
		}

	case *ast.AssignStmt:
		if n.Tok == token.DEFINE {
			this.print("local ")
		}
		this.walkExprList(n.Lhs)
		this.print(" = ")
		this.walkExprList(n.Rhs)

	case *ast.GoStmt:
		this.walkImpl(n.Call)

	case *ast.DeferStmt:
		this.walkImpl(n.Call)

	case *ast.ReturnStmt:
		this.print("return ")
		this.walkExprList(n.Results)

	case *ast.BranchStmt:
		if n.Label == nil {
			switch n.Tok {
			case token.BREAK:
				this.print("break")
			case token.CONTINUE:
				this.printf("goto %s", this.ContinueLabels[n])
			}
		} else {
			switch n.Tok {
			case token.BREAK:
				this.printf("goto %s_break", n.Label)
			case token.CONTINUE:
				this.printf("goto %s_continue", n.Label)
			}
		}

	case *ast.BlockStmt:
		this.indent++
		this.walkStmtList(n.List, "\n")
		this.indent--

	case *ast.IfStmt:
		if n.Init != nil {
			this.walkImpl(n.Init)
		}
		this.print("if ")
		this.walkImpl(n.Cond)

		this.println(" then")
		this.walkImpl(n.Body)
		var elif ast.Node
		if n.Else != nil {
			if nn, ok := n.Else.(*ast.IfStmt); ok {
				this.print("else")
				this.ElseIfs[nn] = struct{}{}
				elif = nn
			} else {
				this.println("else")
			}
			this.walkImpl(n.Else)
		}
		if _, ok := this.ElseIfs[n]; !ok {
			this.print("end")
		}
		if elif != nil {
			delete(this.ElseIfs, elif)
		}

	case *ast.CaseClause:
		this.walkExprList(n.List)
		this.walkStmtList(n.Body, "")

	case *ast.SwitchStmt:
		if n.Init != nil {
			this.walkImpl(n.Init)
		}
		if n.Tag != nil {
			this.walkImpl(n.Tag)
		}
		this.walkImpl(n.Body)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			this.walkImpl(n.Init)
		}
		this.walkImpl(n.Assign)
		this.walkImpl(n.Body)

	case *ast.CommClause:
		if n.Comm != nil {
			this.walkImpl(n.Comm)
		}
		this.walkStmtList(n.Body, "")

	case *ast.SelectStmt:
		this.walkImpl(n.Body)

	case *ast.ForStmt:
		if n.Init != nil {
			this.println("do")
			this.indent++
			this.walkImpl(n.Init)
			this.println()
		}
		if n.Cond != nil {
			this.print("while ")
			this.walkImpl(n.Cond)
			this.println(" do")
		} else {
			this.println("while true do")
		}
		this.walkImpl(n.Body)

		if label, ok := this.ContinueLabels[n]; ok {
			this.printf("::%s::\n", label)
		}
		if n.Post != nil {
			this.indent++
			this.walkImpl(n.Post)
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
			this.walkImpl(n.Key)
			this.print(", ")
		} else {
			this.print("_, ")
		}
		if n.Value != nil {
			this.walkImpl(n.Value)
		}

		this.print(" in pairs(")
		this.walkImpl(n.X)
		this.println(") do")

		this.walkImpl(n.Body)

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
			this.walkImpl(n.Doc)
		}
		if n.Name != nil {
			this.walkImpl(n.Name)
		}
		this.walkImpl(n.Path)
		if n.Comment != nil {
			this.walkImpl(n.Comment)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}
		this.walkIdentList(n.Names)
		if n.Type != nil {
			this.walkImpl(n.Type)
		}
		this.walkExprList(n.Values)
		if n.Comment != nil {
			this.walkImpl(n.Comment)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}
		this.walkImpl(n.Name)
		this.walkImpl(n.Type)
		if n.Comment != nil {
			this.walkImpl(n.Comment)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}
		for _, s := range n.Specs {
			this.walkImpl(s)
		}

	case *ast.FuncDecl:
		if n.Name.Name == "init" && n.Recv == nil {
			this.FuncInit = true
		}

		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}
		if n.Recv != nil {
			this.walkImpl(n.Recv)
		}

		first := n.Name.Name[0]
		if first >= 'a' && first <= 'z' {
			this.print("local ")
		}

		this.walkImpl(n.Name)
		this.print(" = function(")

		this.walkImpl(n.Type)
		this.println(")")

		if n.Body != nil {
			this.walkImpl(n.Body)
		}
		this.println("end")
		this.println()

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			this.walkImpl(n.Doc)
		}

		this.print("-- package: ")
		this.walkImpl(n.Name)
		this.println()
		this.println()

		this.walkDeclList(n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

		if this.FuncInit {
			this.println("return init")
		}

	case *ast.Package:
		for _, f := range n.Files {
			this.walkImpl(f)
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
