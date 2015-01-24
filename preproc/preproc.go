package preproc

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func visitStmt(stmt ast.Stmt) {
}

func visitFunction(f *ast.FuncDecl) {
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
