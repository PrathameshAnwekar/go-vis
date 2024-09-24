package model

import "go/ast"

type GoFile struct {
	Package *ast.Package
	Functions []*ast.FuncDecl
	Imports []*ast.ImportSpec
	Structs []*ast.StructType
}