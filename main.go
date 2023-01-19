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

func New(fs *token.FileSet, functionName string) *FunctionDefinition {
	return &FunctionDefinition{Fset: fs, FunctionName: functionName}
}

type FunctionDefinition struct {
	Fset         *token.FileSet
	FunctionName string
	FuncDec      *ast.FuncDecl
}

func (f *FunctionDefinition) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if bs, ok := n.(*ast.FuncDecl); ok {
		if f.FunctionName == bs.Name.Name {
			f.FuncDec = bs
			return nil
		}
	}
	return f
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

func coisa(f func(x int, y int) int) string {
	// gets the program pointer of the given function
	p := reflect.ValueOf(f).Pointer()
	fc := runtime.FuncForPC(p)
	fileName, _ := fc.FileLine(p)

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	functionBlockStatment := New(fset, strings.Split(fc.Name(), ".")[1])
	ast.Walk(functionBlockStatment, node)

	if functionBlockStatment.FuncDec == nil {
		return "not found"
	}

	return extractFromFile(fileName, int64(functionBlockStatment.FuncDec.Pos()), int64(functionBlockStatment.FuncDec.End()))
}

func sum(x int, y int) int {
	return x + y
}

func main() {
	fmt.Println(coisa(sum))
}
