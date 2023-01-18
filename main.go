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
)

/*
 	1 - Statement is a thing which controls
		the execution
*/

func New(fs *token.FileSet, firstStmtLine int) *BlockStatement {
	return &BlockStatement{Fset: fs, Line: firstStmtLine}
}

type BlockStatement struct {
	Fset  *token.FileSet
	Line  int
	Block *ast.BlockStmt
}

func (f *BlockStatement) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if bs, ok := n.(*ast.BlockStmt); ok {
		stmtstart := bs.Pos()
		stmtline := f.Fset.Position(stmtstart).Line
		if stmtline == f.Line {
			f.Block = bs
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
	buf := make([]byte, (endScope-startScope)+1)
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
	fileName, line := fc.FileLine(p)

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	functionBlockStatment := New(fset, line)
	ast.Walk(functionBlockStatment, node)

	if functionBlockStatment.Block == nil {
		return "not found"
	}

	return extractFromFile(fileName, int64(functionBlockStatment.Block.Lbrace), int64(functionBlockStatment.Block.Rbrace))
}

func sum(x int, y int) int {
	for {
		x += y
	}

	return x + y
}

func main() {
	fmt.Println(coisa(sum))
}
