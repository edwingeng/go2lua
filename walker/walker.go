package walker

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/edwingeng/go2lua/utils"
	"github.com/edwingeng/go2lua/walker/codeprinter"
	"golang.org/x/tools/go/analysis"
)

type gotoLabelInfo struct {
	funcNode ast.Node
	name     string
}

type Walker struct {
	codeprinter.Printer
	root    ast.Node
	shadows map[token.Pos]int

	nextIds        map[string]int
	funcScopeNames map[ast.Node]map[string]struct{}

	NumErrors        int
	FuncInit         bool
	ElseIfs          map[ast.Node]struct{}
	BreakLabels      map[ast.Node]string
	ContinueLabels   map[ast.Node]string
	GotoLabels       map[gotoLabelInfo]struct{}
	ForShadows       map[ast.Node]struct{}
	Fallthroughs     map[ast.Node]ast.Node
	FallthroughCases map[ast.Node]string
}

func NewWalker(pass *analysis.Pass, node ast.Node, opts ...Option) (w *Walker) {
	w = &Walker{
		root:             node,
		nextIds:          make(map[string]int),
		funcScopeNames:   make(map[ast.Node]map[string]struct{}),
		ElseIfs:          make(map[ast.Node]struct{}),
		BreakLabels:      make(map[ast.Node]string),
		ContinueLabels:   make(map[ast.Node]string),
		GotoLabels:       make(map[gotoLabelInfo]struct{}),
		ForShadows:       make(map[ast.Node]struct{}),
		Fallthroughs:     make(map[ast.Node]ast.Node),
		FallthroughCases: make(map[ast.Node]string),
	}
	for _, opt := range opts {
		opt(w)
	}

	file := pass.Fset.File(node.Pos()).Name()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	w.Printer = *codeprinter.NewPrinter(pass, data)
	return
}

func (this *Walker) printError(err error, node ast.Node) {
	utils.PrintErrors(this.Pass, utils.NewNodeError(err, node))
	this.NumErrors++
}

func (this *Walker) makeUniqueName(key string) string {
	this.nextIds[key]++
	return fmt.Sprintf("__unique_%s_%d", key, this.nextIds[key])
}

func (this *Walker) makeFuncScopeUniqueName(funcNode ast.Node, key string) string {
	nm, ok := this.funcScopeNames[funcNode]
	if !ok {
		nm = make(map[string]struct{})
		this.funcScopeNames[funcNode] = nm
	}

	newName := fmt.Sprintf("__%s", key)
	if _, ok := nm[newName]; ok {
		i := 2
		for ; i < 999; i++ {
			str := fmt.Sprintf("%s_x%d", newName, i)
			if _, ok := nm[str]; !ok {
				newName = str
				break
			}
		}
		if i > 999 {
			panic("IMPOSSIBLE")
		}
	}

	nm[newName] = struct{}{}
	return newName
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
	var switchStack []*ast.SwitchStmt
	var caseStack []*ast.CaseClause
	ast.Inspect(this.root, func(node ast.Node) bool {
		if node == nil {
			n := stack[len(stack)-1]
			switch n.(type) {
			case *ast.FuncLit, *ast.FuncDecl:
				funcStack = funcStack[:len(funcStack)-1]
			case *ast.SwitchStmt:
				switchStack = switchStack[:len(switchStack)-1]
			case *ast.CaseClause:
				caseStack = caseStack[:len(caseStack)-1]
			}
			stack = stack[:len(stack)-1]
			return true
		}

		if this.shadows != nil {
			if n1, ok := this.shadows[node.Pos()]; ok {
				var forNode ast.Node
				for i := len(stack) - 1; i >= 0; i-- {
					if _, ok := stack[i].(*ast.ForStmt); ok {
						forNode = stack[i]
						break
					}
				}
				if forNode != nil {
					n2 := this.Pass.Fset.Position(forNode.Pos()).Line
					if n2 == n1 {
						this.ForShadows[forNode] = struct{}{}
					}
				}
			}
		}

		stack = append(stack, node)
		switch n := node.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			funcStack = append(funcStack, node)
		case *ast.SwitchStmt:
			switchStack = append(switchStack, n)
		case *ast.CaseClause:
			caseStack = append(caseStack, n)
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
					return true
				}
				if _, ok := this.ContinueLabels[loopNode]; !ok {
					if loopNodeLabel == "" {
						funcNode := funcStack[len(funcStack)-1]
						this.ContinueLabels[loopNode] = this.makeFuncScopeUniqueName(funcNode, "continue")
					} else {
						this.ContinueLabels[loopNode] = loopNodeLabel + "_continue"
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
					return true
				}
				if loopNode == nil {
					this.printError(fmt.Errorf("unexpected token: %s", n.Tok), node)
					return true
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

			case n.Tok == token.FALLTHROUGH:
				var targetCase *ast.CaseClause
				curSwitch := switchStack[len(switchStack)-1]
				curCase := caseStack[len(caseStack)-1]
				if curSwitch.Body != nil {
					for i, caseClause := range curSwitch.Body.List {
						if curCase == caseClause {
							targetCase = curSwitch.Body.List[i+1].(*ast.CaseClause)
						}
					}
				}
				if targetCase == nil {
					panic("IMPOSSIBLE")
				}
				this.Fallthroughs[node] = targetCase
				this.FallthroughCases[targetCase] = ""
			}
		}

		return true
	})
}

func (this *Walker) walkIdentList(list []*ast.Ident, funcNode ast.Node) {
	for i, x := range list {
		if i > 0 {
			this.Print(", ")
		}
		this.walkImpl(x, funcNode)
	}
}

func (this *Walker) walkExprList(list []ast.Expr, funcNode ast.Node) {
	for i, x := range list {
		if i > 0 {
			this.Print(", ")
		}
		if this.isCallExpr_MakeMap(x) {
			this.Print("{}")
		} else {
			this.walkImpl(x, funcNode)
		}
	}
}

func (this *Walker) walkStmtList(list []ast.Stmt, newline bool, funcNode ast.Node) {
	for _, x := range list {
		this.walkImpl(x, funcNode)
		if newline {
			this.Println()
		}
	}
}

func (this *Walker) walkDeclList(list []ast.Decl, funcNode ast.Node) {
	for _, x := range list {
		this.walkImpl(x, funcNode)
		this.Println()
	}
}

func (this *Walker) Walk() {
	if this.Buffer.Len() > 0 {
		return
	}

	this.initialize()
	this.walkImpl(this.root, nil)

	bts := this.Buffer.Bytes()
	bts = bytes.TrimRight(bts, "\n")
	if n := len(bts); n < this.Buffer.Len() {
		this.Buffer.Truncate(n + 1)
	}

	var newBuf bytes.Buffer
	newBuf.Grow(this.Buffer.Len())
	var blankLine bool
	scanner := bufio.NewScanner(&this.Buffer)
	for scanner.Scan() {
		line := scanner.Bytes()
		newLine := bytes.TrimRightFunc(line, unicode.IsSpace)
		n := len(newLine)
		if n == 0 {
			if blankLine {
				continue
			}
		}
		newBuf.Write(newLine)
		newBuf.WriteByte('\n')
		blankLine = n == 0
	}
	this.Buffer = newBuf
}

func (this *Walker) walkImpl(node ast.Node, funcNode ast.Node) {
	this.CurrentNode = node
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
				this.Print(", ")
			}
			this.walkImpl(f, funcNode)
		}

	// Expressions
	case *ast.BadExpr:
		this.printError(errors.New("bad expression detected"), n)

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
			this.walkImpl(n.Elt, funcNode)
		}

	case *ast.FuncLit:
		this.walkImpl(n.Type, n)
		this.walkImpl(n.Body, n)

	case *ast.CompositeLit:
		if n.Type != nil {
			this.walkImpl(n.Type, funcNode)
		}
		this.Print("{")
		this.walkExprList(n.Elts, funcNode)
		this.Print("}")

	case *ast.ParenExpr:
		this.Print("(")
		this.walkImpl(n.X, funcNode)
		this.Print(")")

	case *ast.SelectorExpr:
		this.walkImpl(n.X, funcNode)
		this.walkImpl(n.Sel, funcNode)

	case *ast.IndexExpr:
		this.walkImpl(n.X, funcNode)
		this.Print("[")
		this.walkImpl(n.Index, funcNode)
		this.Print("]")

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
		bConvert := false
		if len(n.Args) == 1 {
			if f, ok := n.Fun.(*ast.Ident); ok {
				fType := this.Pass.TypesInfo.Types[f]
				if fType.IsType() {
					bConvert = true
				}
			}
		}
		if !bConvert {
			this.walkImpl(n.Fun, funcNode)
			this.Print("(")
		}

		this.walkExprList(n.Args, funcNode)
		if !bConvert {
			this.Print(")")
		}

	case *ast.StarExpr:
		this.walkImpl(n.X, funcNode)

	case *ast.UnaryExpr:
		switch n.Op {
		case token.ADD:
			// Ignore
		case token.SUB:
			this.Print("-")
		case token.NOT:
			this.Print("not ")
		case token.XOR:
			this.Print("~")
		case token.MUL:
			this.printError(fmt.Errorf("unexpected unary op: %v", n.Op), n)
			return
		case token.AND:
			// Ignore
		case token.ARROW:
			this.printError(fmt.Errorf("unexpected unary op: %v", n.Op), n)
			return
		default:
			this.printError(fmt.Errorf("unexpected unary op: %v", n.Op), n)
			return
		}
		this.walkImpl(n.X, funcNode)

	case *ast.BinaryExpr:
		this.printBinarySubexpr(n.X, n, funcNode)
		if str, ok := go2LuaBinaryOperMap[n.Op.String()]; ok {
			this.Print(str)
		} else {
			this.Printf(" %s ", n.Op)
		}
		this.printBinarySubexpr(n.Y, n, funcNode)

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
			this.Indent--
			this.Printf("::%s::\n", n.Label.Name)
			this.Indent++
		} else {
			this.Print()
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
			this.Print(" = ")
			this.walkImpl(n.X, funcNode)
			this.Print(" + 1")
		case token.DEC:
			this.Print(" = ")
			this.walkImpl(n.X, funcNode)
			this.Print(" - 1")
		default:
			panic(fmt.Errorf("unexpected token: %s", n.Tok))
		}

	case *ast.AssignStmt:
		if n.Tok == token.DEFINE {
			this.Print("local ")
		} else {
			local := true
			for _, v := range n.Lhs {
				if x, ok := v.(*ast.Ident); !ok || x.Name != "_" {
					local = false
					break
				}
			}
			if local {
				this.Print("local ")
			}
		}
		this.walkExprList(n.Lhs, funcNode)
		this.Print(" = ")
		this.walkExprList(n.Rhs, funcNode)

	case *ast.GoStmt:
		this.walkImpl(n.Call, funcNode)

	case *ast.DeferStmt:
		this.walkImpl(n.Call, funcNode)

	case *ast.ReturnStmt:
		this.Print("return ")
		this.walkExprList(n.Results, funcNode)

	case *ast.BranchStmt:
		switch n.Tok {
		case token.BREAK:
			if n.Label == nil {
				this.Print("break")
			} else {
				this.Printf("goto %s_break", n.Label)
			}
		case token.CONTINUE:
			if n.Label == nil {
				this.Printf("goto %s", this.ContinueLabels[n])
			} else {
				this.Printf("goto %s_continue", n.Label)
			}
		case token.GOTO:
			if n.Label == nil {
				this.printError(errors.New("missing label"), node)
				return
			} else {
				this.Printf("goto %s", n.Label)
			}
		case token.FALLTHROUGH:
			if n.Label == nil {
				caseNode := this.Fallthroughs[n]
				label := this.FallthroughCases[caseNode]
				this.Println("__fall = true")
				this.Printf("goto %s", label)
			} else {
				this.printError(errors.New("unexpected label"), node)
				return
			}
		}

	case *ast.BlockStmt:
		this.Indent++
		this.walkStmtList(n.List, true, funcNode)
		this.Indent--

	case *ast.IfStmt:
		if n.Init != nil {
			this.Println("do")
			this.Indent++
			this.walkImpl(n.Init, funcNode)
			this.Println()
		}

		this.Print("if ")
		this.walkImpl(n.Cond, funcNode)

		this.Println(" then")
		this.walkImpl(n.Body, funcNode)
		var elif ast.Node
		if n.Else != nil {
			if nn, ok := n.Else.(*ast.IfStmt); ok {
				this.Print("else")
				this.ElseIfs[nn] = struct{}{}
				elif = nn
			} else {
				this.Println("else")
			}
			this.walkImpl(n.Else, funcNode)
		}
		if _, ok := this.ElseIfs[n]; !ok {
			this.Print("end")
		}
		if elif != nil {
			delete(this.ElseIfs, elif)
		}

		if n.Init != nil {
			this.Indent--
			this.Println()
			this.Print("end")
		}

	case *ast.CaseClause:
		this.walkExprList(n.List, funcNode)
		this.walkStmtList(n.Body, false, funcNode)

	case *ast.SwitchStmt:
		if n.Init != nil {
			this.Println("do")
			this.Indent++
			this.walkImpl(n.Init, funcNode)
			this.Println()
		}

		var includeFallthrough bool
		var switchLabel string
		var caseLabel string

		this.Println("repeat")
		this.Indent++
		if n.Tag != nil {
			this.Printf("local __switch = ")
			this.walkImpl(n.Tag, funcNode)
			this.Println()
		}

		if n.Body != nil {
			for _, stmt := range n.Body.List {
				if _, ok := this.FallthroughCases[stmt]; ok {
					includeFallthrough = true
					break
				}
			}
			if includeFallthrough {
				switchLabel = this.makeFuncScopeUniqueName(funcNode, "switch")
				caseLabel = this.makeFuncScopeUniqueName(funcNode, "case")
				this.Println("local __fall = false")
				for i, stmt := range n.Body.List {
					if _, ok := this.FallthroughCases[stmt]; ok {
						this.FallthroughCases[stmt] = fmt.Sprintf("%s_%d", caseLabel, i+1)
					}
				}
			}

			var def *ast.CaseClause
			var c int
			for _, stmt := range n.Body.List {
				caseClause, ok := stmt.(*ast.CaseClause)
				if !ok {
					panic("IMPOSSIBLE")
				}
				if caseClause.List == nil {
					def = caseClause
					continue
				}

				c++
				this.CurrentNode = caseClause
				this.printCaseClauseLabel(includeFallthrough && c > 1, stmt)
				if c == 1 || includeFallthrough {
					this.Printf("if ")
				} else {
					this.Printf("elseif ")
				}
				this.walkCaseClause(caseClause, n.Tag != nil, switchLabel, caseLabel, funcNode)
				if includeFallthrough {
					this.Println("end")
				}
			}

			if def != nil {
				c++
				this.CurrentNode = def
				this.printCaseClauseLabel(includeFallthrough && c > 1, def)
				if c > 0 && !includeFallthrough {
					this.Println("else")
				} else {
					this.Println("do")
				}
				this.walkCaseClause(def, n.Tag != nil, switchLabel, caseLabel, funcNode)
				if includeFallthrough {
					this.Println("end")
				}
			}

			if len(n.Body.List) > 0 && !includeFallthrough {
				this.Println("end")
			}
		}

		this.Indent--
		this.Print("until true")

		if n.Init != nil {
			this.Indent--
			this.Println("end")
		}
		if includeFallthrough {
			this.Println()
			this.Indent--
			this.Printf("::%s_break::", switchLabel)
			this.Indent++
		}

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
			this.Println("do")
			this.Indent++
			this.walkImpl(n.Init, funcNode)
			this.Println()
		}
		if n.Cond != nil {
			this.Print("while ")
			this.walkImpl(n.Cond, funcNode)
			this.Println(" do")
		} else {
			this.Println("while true do")
		}

		if n.Post != nil && n.Body != nil && len(n.Body.List) > 0 {
			if _, ok := this.ForShadows[n]; ok {
				this.Indent++
				this.Println("do")
			}
		}
		this.walkImpl(n.Body, funcNode)
		if n.Post != nil && n.Body != nil && len(n.Body.List) > 0 {
			if _, ok := this.ForShadows[n]; ok {
				this.Println("end")
				this.Indent--
			}
		}

		if label, ok := this.ContinueLabels[n]; ok {
			this.Printf("::%s::\n", label)
		}
		if n.Post != nil {
			this.Indent++
			this.walkImpl(n.Post, funcNode)
			this.Indent--
			this.Println()
		}

		if n.Init != nil {
			this.Println("end")
			this.Indent--
		}

		if label, ok := this.BreakLabels[n]; ok {
			this.Println("end")
			this.Indent--
			this.Printf("::%s::", label)
			this.Indent++
		} else {
			this.Print("end")
		}

	case *ast.RangeStmt:
		this.Print("for ")
		if n.Key != nil {
			this.walkImpl(n.Key, funcNode)
			this.Print(", ")
		} else {
			this.Print("_, ")
		}
		if n.Value != nil {
			this.walkImpl(n.Value, funcNode)
		}

		this.Print(" in pairs(")
		this.walkImpl(n.X, funcNode)
		this.Println(") do")

		this.walkImpl(n.Body, funcNode)

		if label, ok := this.ContinueLabels[n]; ok {
			this.Printf("::%s::\n", label)
		}
		if label, ok := this.BreakLabels[n]; ok {
			this.Println("end")
			this.Indent--
			this.Printf("::%s::", label)
			this.Indent++
		} else {
			this.Print("end")
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
		switch n.Tok {
		case token.VAR, token.CONST:
			this.CurrentNode = n
			this.Print()
			for i, s := range n.Specs {
				if i > 0 {
					this.CurrentNode = s
					this.Println()
				}
				var mixedNames, lowerNames, upperNames []string
				spec, ok := s.(*ast.ValueSpec)
				if !ok {
					panic("IMPOSSIBLE")
				}
				if funcNode != nil {
					for _, name := range spec.Names {
						mixedNames = append(mixedNames, name.Name)
					}
				} else {
					for _, name := range spec.Names {
						if ast.IsExported(name.Name) {
							upperNames = append(upperNames, name.Name)
						} else {
							lowerNames = append(lowerNames, name.Name)
						}
					}
				}

				if len(mixedNames) > 0 {
					this.printVarDefinition(true, mixedNames, spec, funcNode)
				}
				if len(lowerNames) > 0 {
					this.printVarDefinition(true, lowerNames, spec, funcNode)
				}
				if len(lowerNames) > 0 && len(upperNames) > 0 {
					this.Println()
				}
				if len(upperNames) > 0 {
					this.printVarDefinition(false, upperNames, spec, funcNode)
				}
			}
		}

	case *ast.FuncDecl:
		if n.Name.Name == "init" && n.Recv == nil {
			this.FuncInit = true
		}

		this.Print()
		if n.Doc != nil {
			this.walkImpl(n.Doc, n)
		}
		if n.Recv != nil {
			this.walkImpl(n.Recv, n)
		}

		if !ast.IsExported(n.Name.Name) {
			this.Print("local ")
		}

		this.walkImpl(n.Name, n)
		this.Print(" = function(")

		this.walkImpl(n.Type, n)
		this.Println(")")

		if n.Body != nil {
			this.walkImpl(n.Body, n)
		}
		this.Println("end")

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			this.walkImpl(n.Doc, funcNode)
		}

		this.Print("-- package: ")
		this.walkImpl(n.Name, funcNode)
		this.Println()

		this.walkDeclList(n.Decls, funcNode)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

		if this.FuncInit {
			this.Println()
			this.Println("return init")
		}

	case *ast.Package:
		for _, f := range n.Files {
			this.walkImpl(f, funcNode)
		}

	default:
		panic(fmt.Sprintf("ast.walkImpl: unexpected node type %T", n))
	}
}

func (this *Walker) printVarDefinition(local bool, names []string, spec *ast.ValueSpec, funcNode ast.Node) {
	if local {
		this.Print("local ")
	}
	leN := len(names)
	switch leN {
	case 0:
	case 1:
		this.Printf("%s ", names[0])
	default:
		this.Printf("%s ", strings.Join(names, ", "))
	}

	if len(spec.Values) > 0 {
		this.Print("= ")
		for i, v := range spec.Values {
			if i > 0 {
				this.Print(", ")
			}
			this.walkImpl(v, funcNode)
		}
	} else {
		this.Print("= ")
		typ := this.Pass.TypesInfo.Types[spec.Type].Type
		defVal := utils.DefaultValue(typ)
		for i := 0; i < leN; i++ {
			if i > 0 {
				this.Print(", ")
			}
			this.Print(defVal)
		}
	}
}

func (this *Walker) printCaseClauseLabel(newline bool, node ast.Node) {
	if newline {
		this.Println()
	}
	if str, ok := this.FallthroughCases[node]; ok {
		this.Indent--
		this.Printf("::%s::\n", str)
		this.Indent++
	}
}

func (this *Walker) walkCaseClause(node *ast.CaseClause, hasTag bool, switchLabel, caseLabel string, funcNode ast.Node) {
	this.Indent++
	if node.List == nil {
		this.Println("-- default")
	}
	_, fallthroughCase := this.FallthroughCases[node]
	for i, expr := range node.List {
		if i > 0 {
			this.Printf("or ")
		} else if fallthroughCase {
			this.Print(" __fall or ")
		}
		switch e := expr.(type) {
		case *ast.BasicLit:
			if hasTag {
				this.Printf("__switch == %s ", e.Value)
			} else {
				this.Printf("%s ", e.Value)
			}
		case *ast.Ident:
			if hasTag {
				this.Printf("__switch == %s ", e.Name)
			} else {
				this.Printf("%s ", e.Name)
			}
		default:
			if hasTag {
				this.Print("__switch == ")
			}
			this.walkImpl(e, funcNode)
			this.Print(" ")
		}
	}

	if node.List != nil {
		this.Println("then")
	}
	if fallthroughCase {
		this.Println("__fall = false")
	}
	this.walkStmtList(node.Body, true, funcNode)
	if caseLabel != "" {
		if n := len(node.Body); n > 0 {
			if _, ok := node.Body[n-1].(*ast.BranchStmt); !ok {
				this.Printf("goto %s_break\n", switchLabel)
			}
		}
	}
	this.Indent--
}

func (this *Walker) printBinarySubexpr(e ast.Expr, n *ast.BinaryExpr, funcNode ast.Node) {
	this.walkImpl(e, funcNode)
}

type Option func(w *Walker)

func WithShadows(shadows map[token.Pos]int) Option {
	return func(w *Walker) {
		w.shadows = shadows
	}
}
