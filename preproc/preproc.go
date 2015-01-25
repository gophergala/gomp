package preproc

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/gophergala/gomp/gensym"
)

type Cond int

type Context struct {
	genSym        func() string
	runtimeCalled bool
	cmap          ast.CommentMap
}

const (
	COND_LT = iota
	COND_LE
	COND_GT
	COND_GE
)

func parseForInit(stmt *ast.Stmt) (variable *ast.Ident, initExpr *ast.Expr, ok bool) {
	if stmt == nil {
		return
	}
	var assignStmt *ast.AssignStmt
	if assignStmt, ok = (*stmt).(*ast.AssignStmt); !ok {
		return
	}
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return
	}
	if variable, ok = assignStmt.Lhs[0].(*ast.Ident); !ok {
		return
	}
	initExpr = &assignStmt.Rhs[0]
	return
}

func parseForCond(expr *ast.Expr) (variable *ast.Ident, op token.Token, bound *ast.Expr, ok bool) {
	if expr == nil {
		return
	}
	binaryExpr, ok := (*expr).(*ast.BinaryExpr)
	if !ok {
		return
	}
	switch binaryExpr.Op {
	case token.LEQ, token.LSS, token.GTR, token.GEQ:
		op = binaryExpr.Op
	default:
		return
	}
	if variable, ok = binaryExpr.X.(*ast.Ident); !ok {
		return
	}
	bound = &binaryExpr.Y
	return
}

func parseForPost(stmt *ast.Stmt) (variable *ast.Ident, op token.Token, postExpr ast.Expr, ok bool) {
	if stmt == nil {
		return
	}

	if incDecStmt, isIncDec := (*stmt).(*ast.IncDecStmt); isIncDec {
		variable, ok = incDecStmt.X.(*ast.Ident)
		postExpr = mkIntLit(1)
		switch incDecStmt.Tok {
		case token.INC:
			op = token.ADD_ASSIGN
		case token.DEC:
			op = token.SUB_ASSIGN
		default:
			panic("Unknown op in IncDecStmt")
		}

		return
	}
	if assignStmt, isAssignStmt := (*stmt).(*ast.AssignStmt); isAssignStmt {
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			return
		}
		if variable, ok = assignStmt.Lhs[0].(*ast.Ident); !ok {
			return
		}
		switch assignStmt.Tok {
		case token.ADD_ASSIGN, token.SUB_ASSIGN:
			op = assignStmt.Tok
		default:
			ok = false
			return
		}
		postExpr = assignStmt.Rhs[0]
	}
	return
}

func mkIntLit(n int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(n)}
}

func mkSym(context *Context) *ast.Ident {
	return &ast.Ident{Name: context.genSym()}
}

func mkIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func mkTypeConv(expr ast.Expr, t string) ast.Expr {
	return &ast.CallExpr{
		Fun:  mkIdent(t),
		Args: []ast.Expr{expr},
	}
}

func mkGoLambda(body *ast.BlockStmt, arg *ast.Ident) *ast.GoStmt {
	return &ast.GoStmt{
		Call: &ast.CallExpr{
			Fun: &ast.FuncLit{
				Type: &ast.FuncType{Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{arg},
							Type:  mkIdent("int")}}}},
				Body: body,
			},
			Args: []ast.Expr{mkTypeConv(arg, "int")},
		},
	}
}

func emitSchedulerLoop(originVar, begin, end, step *ast.Ident,
	context *Context, originBody *ast.BlockStmt) (code []ast.Stmt) {
	// taskSize := (end - begin + 1) / (numCPU * step)
	taskSize := mkSym(context)
	nom := ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  end,
			Op: token.SUB,
			Y:  begin,
		},
		Op: token.ADD,
		Y:  mkIntLit(1),
	}
	denom := ast.BinaryExpr{
		X:  step,
		Op: token.MUL,
		Y: &ast.CallExpr{Fun: &ast.SelectorExpr{
			X:   mkIdent("runtime"),
			Sel: mkIdent("NumCPU")}},
	}
	taskSizeStmt := ast.AssignStmt{
		Lhs: []ast.Expr{taskSize},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.BinaryExpr{X: &nom, Op: token.QUO, Y: &denom}},
	}
	context.runtimeCalled = true
	code = append(code, &taskSizeStmt)

	routineId, perRoutineCounter := mkSym(context), mkSym(context)

	routineIdBoundExpr := ast.BinaryExpr{
		X:  begin,
		Op: token.ADD,
		Y: &ast.BinaryExpr{
			X:  routineId,
			Op: token.MUL,
			Y: &ast.BinaryExpr{
				X:  taskSize,
				Op: token.MUL,
				Y:  step,
			},
		},
	}

	// begin + routineId * taskSize * step <= end
	routineIdExpr := ast.BinaryExpr{
		X:  &routineIdBoundExpr,
		Op: token.LEQ,
		Y:  end,
	}

	channel, channelSize := mkSym(context), mkSym(context)
	channelSizeDecl := ast.AssignStmt{
		Lhs: []ast.Expr{channelSize},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X: &ast.BinaryExpr{
						X:  end,
						Op: token.SUB,
						Y:  begin,
					},
					Op: token.QUO,
					Y: &ast.BinaryExpr{
						X:  taskSize,
						Op: token.MUL,
						Y:  step,
					},
				},
				Op: token.ADD,
				Y:  mkIntLit(1)}}}
	emptyStruct := ast.StructType{
		Fields: &ast.FieldList{},
	}
	channelDecl := ast.AssignStmt{
		Lhs: []ast.Expr{channel},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: mkIdent("make"),
				Args: []ast.Expr{
					&ast.ChanType{
						Value: &emptyStruct,
						Dir:   ast.SEND | ast.RECV,
					},
					channelSize,
				},
			},
		},
	}
	code = append(code, &channelSizeDecl)
	code = append(code, &channelDecl)

	{
		nestedLoop := ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{originVar, perRoutineCounter},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{&routineIdBoundExpr, mkIntLit(0)},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  originVar,
					Op: token.LEQ,
					Y:  end,
				},
				Op: token.LAND,
				Y: &ast.BinaryExpr{
					X:  perRoutineCounter,
					Op: token.LSS,
					Y:  taskSize,
				},
			},
			Post: &ast.AssignStmt{
				Lhs: []ast.Expr{originVar, perRoutineCounter},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{&ast.BinaryExpr{
					X:  originVar,
					Op: token.ADD,
					Y:  step},
					&ast.BinaryExpr{
						X:  perRoutineCounter,
						Op: token.ADD,
						Y:  mkIntLit(1)}},
			},
			Body: originBody,
		}
		loop := ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{routineId},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{mkIntLit(0)},
			},
			Cond: &routineIdExpr,
			Post: &ast.IncDecStmt{
				X:   routineId,
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{
				mkGoLambda(
					&ast.BlockStmt{
						List: []ast.Stmt{
							&nestedLoop,
							&ast.SendStmt{
								Chan: channel,
								Value: &ast.CompositeLit{
									Type: &emptyStruct,
								},
							}}},
					routineId),
			}},
		}
		code = append(code, &loop)
	}
	{
		loopVar := mkSym(context)
		code = append(code, &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{loopVar},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{mkIntLit(0)},
			},
			Cond: &ast.BinaryExpr{
				X:  loopVar,
				Op: token.LSS,
				Y:  channelSize,
			},
			Post: &ast.IncDecStmt{
				X:   loopVar,
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.UnaryExpr{
							X:  channel,
							Op: token.ARROW,
						},
					},
				},
			},
		})
	}
	return code
}

func visitFor(stmt *ast.ForStmt, context *Context) *ast.BlockStmt {
	initVar, initExpr, initOk := parseForInit(&stmt.Init)
	condVar, condOp, condExpr, condOk := parseForCond(&stmt.Cond)
	postVar, postOp, postExpr, postOk := parseForPost(&stmt.Post)

	if !initOk || !condOk || !postOk {
		return nil
	}
	if initVar.Name != condVar.Name || initVar.Name != postVar.Name {
		return nil
	}

	block := new(ast.BlockStmt)
	block.List = []ast.Stmt{}
	initVarSym, condVarSym, incVarSym := mkSym(context), mkSym(context), mkSym(context)
	{
		boundsDecl := ast.AssignStmt{
			Lhs: []ast.Expr{initVarSym, condVarSym, incVarSym},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{*initExpr, *condExpr, postExpr},
		}

		*initExpr, *condExpr = ast.Expr(initVarSym), ast.Expr(condVarSym)
		stmt.Post = &ast.AssignStmt{
			Lhs: []ast.Expr{initVar},
			Tok: postOp,
			Rhs: []ast.Expr{incVarSym},
		}

		block.List = append(block.List, &boundsDecl)
	}

	if condOp == token.LEQ {
		block.List = append(
			block.List,
			emitSchedulerLoop(initVar, initVarSym, condVarSym, incVarSym, context, stmt.Body)...)
	} else {
		block.List = append(block.List, ast.Stmt(stmt))
	}
	return block
}

func visitExpr(e *ast.Expr, context *Context) {
	if e == nil {
		return
	}
	switch t := (*e).(type) {
	case *ast.FuncLit:
		if t.Body == nil {
			return
		}
		for _, s := range t.Body.List {
			visitStmt(&s, context)
		}
	}
}

func shouldParalellize(stmt *ast.Stmt, context *Context) bool {
	commentGroups := ((*context).cmap)[(*stmt).(ast.Node)]
	length := len(commentGroups)
	if length == 0 {
		return false
	}
	commentGroup := *commentGroups[length-1]
	length = len(commentGroup.List)
	if length == 0 {
		return false
	}
	if !strings.HasPrefix(commentGroup.List[length-1].Text, "//gomp") {
		return false
	}

	return true
}

func visitStmt(stmt *ast.Stmt, context *Context) {
	if stmt == nil {
		return
	}
	switch t := (*stmt).(type) {
	case *ast.AssignStmt:
		for _, e := range t.Rhs {
			visitExpr(&e, context)
		}
	case *ast.ForStmt:
		if shouldParalellize(stmt, context) {
			if block := visitFor(t, context); block != nil {
				*stmt = block
				//TODO: save old comments here
			}
		}
	case *ast.BlockStmt:
		visitBlock(t, context)
	case *ast.IfStmt:
		visitBlock(t.Body, context)
	case *ast.SwitchStmt:
		visitBlock(t.Body, context)
	case *ast.TypeSwitchStmt:
		visitBlock(t.Body, context)
	case *ast.CaseClause:
		for i, _ := range t.Body {
			visitStmt(&t.Body[i], context)
		}
	}
}

func visitBlock(stmt *ast.BlockStmt, context *Context) {
	if stmt != nil {
		for i, _ := range stmt.List {
			visitStmt(&stmt.List[i], context)
		}
	}
}

func visitFunction(f *ast.FuncDecl, context *Context) {
	if f.Body != nil {
		visitBlock(f.Body, context)
	}
}

// Run preprocessor on a source. filename is used for error reporting.
// This function is currently not implemented.
func PreprocFile(source, filename string) (result string, err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, source,
		parser.ParseComments|parser.AllErrors)
	if err != nil {
		return
	}
	context := Context{gensym.MkGen(source), false, ast.NewCommentMap(fset, file, file.Comments)}

	for _, decl := range file.Decls {
		switch t := decl.(type) {
		case *ast.FuncDecl:
			visitFunction(t, &context)
		}
	}

	if context.runtimeCalled {
		const runtimePath = `"runtime"`
		runtimeImported := false
		for _, spec := range file.Imports {
			if spec.Path != nil && spec.Path.Value == runtimePath {
				runtimeImported = true
				break
			}
		}
		if !runtimeImported {
			runtimeImport := ast.ImportSpec{
				Path: &ast.BasicLit{Value: runtimePath, Kind: token.STRING}}
			runtimeDecl := ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{&runtimeImport}}
			file.Decls = append([]ast.Decl{&runtimeDecl}, file.Decls...)
			file.Imports = append(file.Imports, &runtimeImport)
		}
	}
	file.Imports = []*ast.ImportSpec{}

	//Delete all comments from file
	file.Comments = nil
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), file)
	result = buf.String()
	return
}
