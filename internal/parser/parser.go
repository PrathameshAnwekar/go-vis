package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/PrathameshAnwekar/go-vis/internal/log"
)

func ParseGoProject(fileList []string) {
	for _, file := range fileList {
		parseFile(file)
	}
}

func parseFile(fileName string) error {
	fileSet := token.NewFileSet()

	node, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error while parsing file - \"%s\".\n Error: %e", fileName, err)
	}
	log.D("Inspecting contents of file", fileName, "in package ", node.Name.Name)

	ast.Inspect(node, func(n ast.Node) bool {
		switch element := n.(type) {
		case *ast.FuncDecl:
			{
				log.D("Recognised function", element.Name.Name)
			}
		case *ast.TypeSpec:
			{
				log.D("Recognised typespec", element.Name.Name)

			}
		}
		return true
	})

	return nil
}
