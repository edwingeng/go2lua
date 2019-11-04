package walker

import (
	"bytes"
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
	bytes.Buffer
	fset *token.FileSet

	nextNumber int
	indent     int

	FuncInit bool
	Elseifs  map[ast.Node]struct{}
}

func NewWalker(fset *token.FileSet) *Walker {
	w := &Walker{
		fset:    fset,
		Elseifs: make(map[ast.Node]struct{}),
	}
	return w
}

func (this *Walker) PrintIndent() {
	bts := this.Buffer.Bytes()
	if n := len(bts); n > 0 && bts[n-1] == '\n' {
		for i := 0; i < this.indent; i++ {
			this.Buffer.Write(indentBytes)
		}
	}
}

func (this *Walker) Print(a ...interface{}) {
	this.PrintIndent()
	_, _ = fmt.Fprint(&this.Buffer, a...)
}

func (this *Walker) Println(a ...interface{}) {
	this.PrintIndent()
	_, _ = fmt.Fprintln(&this.Buffer, a...)
}

func (this *Walker) Printf(format string, a ...interface{}) {
	this.PrintIndent()
	_, _ = fmt.Fprintf(&this.Buffer, format, a...)
}

func (this *Walker) printError(x interface{}) {
	err := ast.Fprint(os.Stderr, this.fset, x, nil)
	if err != nil {
		panic(err)
	}
}

func (this *Walker) NextName_NonameFunc() string {
	this.nextNumber++
	return fmt.Sprintf("noname_func_%d", this.nextNumber)
}

func (this *Walker) walkIdentList(list []*ast.Ident) {
	for i, x := range list {
		if i > 0 {
			this.Print(", ")
		}
		this.Walk(x)
	}
}

func (this *Walker) walkExprList(list []ast.Expr) {
	for i, x := range list {
		if i > 0 {
			this.Print(", ")
		}
		this.Walk(x)
	}
}

func (this *Walker) walkStmtList(list []ast.Stmt, sep string) {
	for _, x := range list {
		this.Walk(x)
		this.Print(sep)
	}
}

func (this *Walker) walkDeclList(list []ast.Decl) {
	for _, x := range list {
		this.Walk(x)
	}
}

func (this *Walker) Walk(node ast.Node) {
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		for _, c := range n.List {
			this.Walk(c)
		}

	case *ast.Field:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		this.walkIdentList(n.Names)
		if n.Tag != nil {
			this.Walk(n.Tag)
		}
		if n.Comment != nil {
			this.Walk(n.Comment)
		}

	case *ast.FieldList:
		for i, f := range n.List {
			if i > 0 {
				this.Print(", ")
			}
			this.Walk(f)
		}

	// Expressions
	case *ast.BadExpr:
		this.printError(n)

	case *ast.Ident:
		if str, ok := go2LuaFuncMap[n.Name]; ok {
			this.Print(str)
		} else {
			this.Print(n.Name)
		}

	case *ast.BasicLit:
		this.Print(n.Value)

	case *ast.Ellipsis:
		if n.Elt != nil {
			this.Walk(n.Elt)
		}

	case *ast.FuncLit:
		this.Walk(n.Type)
		this.Walk(n.Body)

	case *ast.CompositeLit:
		if n.Type != nil {
			this.Walk(n.Type)
		}
		this.walkExprList(n.Elts)

	case *ast.ParenExpr:
		this.Walk(n.X)

	case *ast.SelectorExpr:
		this.Walk(n.X)
		this.Walk(n.Sel)

	case *ast.IndexExpr:
		this.Walk(n.X)
		this.Walk(n.Index)

	case *ast.SliceExpr:
		this.Walk(n.X)
		if n.Low != nil {
			this.Walk(n.Low)
		}
		if n.High != nil {
			this.Walk(n.High)
		}
		if n.Max != nil {
			this.Walk(n.Max)
		}

	case *ast.TypeAssertExpr:
		this.Walk(n.X)
		if n.Type != nil {
			this.Walk(n.Type)
		}

	case *ast.CallExpr:
		this.Walk(n.Fun)

		this.Print("(")
		this.walkExprList(n.Args)
		this.Print(")")

	case *ast.StarExpr:
		this.Walk(n.X)

	case *ast.UnaryExpr:
		this.Walk(n.X)

	case *ast.BinaryExpr:
		this.Walk(n.X)
		if str, ok := go2LuaOperMap[n.Op.String()]; ok {
			this.Printf(" %s ", str)
		} else {
			this.Printf(" %s ", n.Op)
		}
		this.Walk(n.Y)

	case *ast.KeyValueExpr:
		this.Walk(n.Key)
		this.Walk(n.Value)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			this.Walk(n.Len)
		}
		this.Walk(n.Elt)

	case *ast.StructType:
		this.Walk(n.Fields)

	case *ast.FuncType:
		if n.Params != nil {
			this.Walk(n.Params)
		}

	case *ast.InterfaceType:
		this.Walk(n.Methods)

	case *ast.MapType:
		this.Walk(n.Key)
		this.Walk(n.Value)

	case *ast.ChanType:
		this.Walk(n.Value)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		this.Walk(n.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		this.Walk(n.Label)
		this.Walk(n.Stmt)

	case *ast.ExprStmt:
		this.Walk(n.X)

	case *ast.SendStmt:
		this.Walk(n.Chan)
		this.Walk(n.Value)

	case *ast.IncDecStmt:
		this.Walk(n.X)
		switch n.Tok {
		case token.INC:
			this.Print(" = ")
			this.Walk(n.X)
			this.Print(" + 1")
		case token.DEC:
			this.Print(" = ")
			this.Walk(n.X)
			this.Print(" - 1")
		default:
			panic(fmt.Errorf("unexpected token: %v", n.Tok))
		}

	case *ast.AssignStmt:
		if n.Tok == token.DEFINE {
			this.Print("local ")
		}
		this.walkExprList(n.Lhs)
		this.Print(" = ")
		this.walkExprList(n.Rhs)

	case *ast.GoStmt:
		this.Walk(n.Call)

	case *ast.DeferStmt:
		this.Walk(n.Call)

	case *ast.ReturnStmt:
		this.Print("return ")
		this.walkExprList(n.Results)

	case *ast.BranchStmt:
		if n.Label != nil {
			this.Walk(n.Label)
		} else {
			this.Print("break")
		}

	case *ast.BlockStmt:
		this.indent++
		this.walkStmtList(n.List, "\n")
		this.indent--

	case *ast.IfStmt:
		if n.Init != nil {
			this.Walk(n.Init)
		}
		this.Print("if ")
		this.Walk(n.Cond)

		this.Println(" then")
		this.Walk(n.Body)
		var elif ast.Node
		if n.Else != nil {
			if nn, ok := n.Else.(*ast.IfStmt); ok {
				this.Print("else")
				this.Elseifs[nn] = struct{}{}
				elif = nn
			} else {
				this.Println("else")
			}
			this.Walk(n.Else)
		}
		if _, ok := this.Elseifs[n]; !ok {
			this.Print("end")
		}
		if elif != nil {
			delete(this.Elseifs, elif)
		}

	case *ast.CaseClause:
		this.walkExprList(n.List)
		this.walkStmtList(n.Body, "")

	case *ast.SwitchStmt:
		if n.Init != nil {
			this.Walk(n.Init)
		}
		if n.Tag != nil {
			this.Walk(n.Tag)
		}
		this.Walk(n.Body)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			this.Walk(n.Init)
		}
		this.Walk(n.Assign)
		this.Walk(n.Body)

	case *ast.CommClause:
		if n.Comm != nil {
			this.Walk(n.Comm)
		}
		this.walkStmtList(n.Body, "")

	case *ast.SelectStmt:
		this.Walk(n.Body)

	case *ast.ForStmt:
		if n.Init != nil {
			this.Println("do")
			this.indent++
			this.Walk(n.Init)
			this.Println()
		}
		if n.Cond != nil {
			this.Print("while ")
			this.Walk(n.Cond)
			this.Println(" do")
		} else {
			this.Println("while true do")
		}
		this.Walk(n.Body)
		if n.Post != nil {
			this.indent++
			this.Walk(n.Post)
			this.indent--
			this.Println()
		}

		if n.Init != nil {
			this.Println("end")
			this.indent--
		}
		this.Print("end")

	case *ast.RangeStmt:
		if n.Key != nil {
			this.Walk(n.Key)
		}
		if n.Value != nil {
			this.Walk(n.Value)
		}
		this.Walk(n.X)
		this.Walk(n.Body)

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		if n.Name != nil {
			this.Walk(n.Name)
		}
		this.Walk(n.Path)
		if n.Comment != nil {
			this.Walk(n.Comment)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		this.walkIdentList(n.Names)
		if n.Type != nil {
			this.Walk(n.Type)
		}
		this.walkExprList(n.Values)
		if n.Comment != nil {
			this.Walk(n.Comment)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		this.Walk(n.Name)
		this.Walk(n.Type)
		if n.Comment != nil {
			this.Walk(n.Comment)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		for _, s := range n.Specs {
			this.Walk(s)
		}

	case *ast.FuncDecl:
		if n.Name.Name == "init" && n.Recv == nil {
			this.FuncInit = true
		}

		if n.Doc != nil {
			this.Walk(n.Doc)
		}
		if n.Recv != nil {
			this.Walk(n.Recv)
		}

		first := n.Name.Name[0]
		if first >= 'a' && first <= 'z' {
			this.Print("local ")
		}

		this.Walk(n.Name)
		this.Print(" = function(")

		this.Walk(n.Type)
		this.Println(")")

		if n.Body != nil {
			this.Walk(n.Body)
		}
		this.Println("end")
		this.Println()

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			this.Walk(n.Doc)
		}

		this.Print("-- package: ")
		this.Walk(n.Name)
		this.Println()
		this.Println()

		this.walkDeclList(n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

		if this.FuncInit {
			this.Println("return init")
		}

	case *ast.Package:
		for _, f := range n.Files {
			this.Walk(f)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}
}
