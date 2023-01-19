package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
)

/*
 	1 - Statement is a thing which controls
		the execution
*/

func NewFunc(fs *token.FileSet, functionName string) *Func {
	return &Func{Fset: fs, Name: functionName}
}

type Func struct {
	Fset    *token.FileSet
	Name    string
	ASTDecl *ast.FuncDecl
}

func (f *Func) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return f
	}
	if f.Name != funcDecl.Name.Name {
		return f
	}
	f.ASTDecl = funcDecl
	return nil
}

func extractFromFile(filename string, startScope, endScope int64) string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Seek(startScope-1, 0)
	buf := make([]byte, (endScope - startScope))
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func stringfyFunction(f func(x int, y int) int) string {
	// gets the program pointer of the given function
	p := reflect.ValueOf(f).Pointer()
	fc := runtime.FuncForPC(p)

	fileName, _ := fc.FileLine(p)
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	fDec := NewFunc(fset, strings.Split(fc.Name(), ".")[1])
	ast.Walk(fDec, node)

	if fDec.ASTDecl == nil {
		return "not found"
	}

	return extractFromFile(fileName, int64(fDec.ASTDecl.Pos()), int64(fDec.ASTDecl.End()))
}

func sum(x int, y int) int {
	return x + y
}

func main() {
	fmt.Println(stringfyFunction(sum))
}
