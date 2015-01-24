package preproc

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/gophergala/gomp/gensym"
)

type Cond int

type Context struct {
	symGen func() string
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

func parseForPost(stmt *ast.Stmt) (variable *ast.Ident, op token.Token, ok bool) {
	if stmt == nil {
		return
	}

	if incDecStmt, isIncDec := (*stmt).(*ast.IncDecStmt); isIncDec {
		variable, ok = incDecStmt.X.(*ast.Ident)
		op = incDecStmt.Tok
		return
	}
	return
}

func visitFor(stmt *ast.ForStmt, context *Context) *ast.BlockStmt {
	initVar, _, initOk := parseForInit(&stmt.Init)
	condVar, _, _, condOk := parseForCond(&stmt.Cond)
	postVar, _, postOk := parseForPost(&stmt.Post)

	if !initOk || !condOk || !postOk {
		return nil
	}
	if initVar.Name != condVar.Name || initVar.Name != postVar.Name {
		return nil
	}

	block := new(ast.BlockStmt)
	block.List = []ast.Stmt{ast.Stmt(stmt)}
	return block
}

func visitStmt(stmt *ast.Stmt, context *Context) {
	if stmt == nil {
		return
	}
	if forStmt, ok := (*stmt).(*ast.ForStmt); ok {
		if block := visitFor(forStmt, context); block != nil {
			*stmt = block
		}
	}
}

func visitFunction(f *ast.FuncDecl, context *Context) {
	if f.Body == nil {
		return
	}
	for i, _ := range f.Body.List {
		visitStmt(&f.Body.List[i], context)
	}
}

// Run preprocessor on a source. filename is used for error reporting.
// This function is currently not implemented.
func PreprocFile(source, filename string) (result string, err error) {
	return
}

// This function should be used instead of PreprocFile.
func PreprocFileImpl(source, filename string) (result string, err error) {
	context := Context{gensym.MkGen(source)}

	file, err := parser.ParseFile(token.NewFileSet(), filename, source,
		parser.ParseComments|parser.AllErrors)
	if err != nil {
		return
	}
	for _, decl := range file.Decls {
		if fun, ok := decl.(*ast.FuncDecl); ok {
			visitFunction(fun, &context)
		}
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), file)
	result = buf.String()
	return
}
