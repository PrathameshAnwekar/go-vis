package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a project path.")
		return
	}
	fmt.Println("Parsing project:", args)

	projectRoot := args[1]
	goFileList, err := getGoFiles(projectRoot)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range goFileList {
		parseFile(file)
	}

	goFileList, err = removeCommonPrefix(goFileList)
	fmt.Println(goFileList)
}

func getGoFiles(projectRoot string) ([]string, error) {
	fmt.Println("Getting go files from:", projectRoot)
	fileList := make([]string, 0)
	err := filepath.WalkDir(projectRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			fmt.Println("Adding: ", path)
			fileList = append(fileList, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func removeCommonPrefix(list []string) ([]string, error) {
	if len(list) <= 1 {
		return list, nil
	}
	prefix := list[0]
	for _, item := range list {
		for strings.Index(item, prefix) != 0 {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				break
			}
		}
	}
	for i, str := range list {
		list[i] = strings.TrimPrefix(str, prefix)
	}
	return list, nil
}

func parseFile(fileName string) error {
	fileSet := token.NewFileSet()

	node, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error while parsing file - %s: %e", fileName, err)
	}
	fmt.Println("Inspecting contents of file", fileName, "in package ", node.Name.Name)

	ast.Inspect(node, func(n ast.Node) bool {
		switch element := n.(type) {
		case *ast.FuncDecl:
			{
				fmt.Println("Recognised function", element.Name.Name)
			}
		case *ast.TypeSpec:
			{
				fmt.Println("Recognised typespec", element.Name.Name)

			}
		}
		return true
	})

	return nil
}
