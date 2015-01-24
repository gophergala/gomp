package preproc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func visitForInit(init ast.Stmt) (variable string, value int64, ok bool) {
	if init == nil {
		return
	}
	var assignStmt *ast.AssignStmt
	if assignStmt, ok = init.(*ast.AssignStmt); !ok {
		return
	}
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return
	}
	var lhs *ast.Ident
	if lhs, ok = assignStmt.Lhs[0].(*ast.Ident); !ok {
		return
	}
	variable = lhs.Name
	var rhs *ast.BasicLit
	if rhs, ok = assignStmt.Rhs[0].(*ast.BasicLit); !ok || rhs.Kind != token.INT {
		return
	}
	value, err := strconv.ParseInt(rhs.Value, 0, 64)
	if err != nil {
		return
	}
	return
}

func visitFor(stmt ast.ForStmt) {
	if stmt.Init == nil || stmt.Cond == nil || stmt.Post == nil {
		return
	}
	variable, value, ok := visitForInit(stmt.Init)
	if !ok {
		return
	}
	fmt.Println("Variable:", variable, ", value:", value)
}

func visitStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	if forStmt, ok := stmt.(*ast.ForStmt); ok {
		visitFor(*forStmt)
	}
}

func visitFunction(f *ast.FuncDecl) {
	fmt.Println("Visiting function: ", f.Name)
	if f.Body == nil {
		return
	}
	for _, stmt := range f.Body.List {
		visitStmt(stmt)
	}
}

// Run preprocessor on a source. filename is used for error reporting.
// This function is currently not implemented.
func PreprocFile(source, filename string) (result string, err error) {
	return
}

// This function should be used instead of PreprocFile.
func PreprocFileImpl(source, filename string) (result string, err error) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, source,
		parser.ParseComments|parser.AllErrors)
	if err != nil {
		return
	}
	for _, decl := range file.Decls {
		if fun, ok := decl.(*ast.FuncDecl); ok {
			visitFunction(fun)
		}
	}
	return
}
