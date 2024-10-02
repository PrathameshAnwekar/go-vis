package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/PrathameshAnwekar/go-vis/internal/log"
	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/tools/go/packages"
)

var FuncSet = mapset.NewSet[*ast.FuncDecl]()

func ParseGoProject(fileList []string) error {

	for _, file := range fileList {
		err := parseFile(file)
		if err != nil {
			return err
		}
	}

	err := getProjectPackageList()
	if err != nil {
		return err
	}

	return nil
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
				if FuncSet.Contains(element) {
					break
				}
				FuncSet.Add(element)
				calls := parseFunctionBody(element.Body)
				if len(calls) == 0 {
					log.D("No function calls detected.")
				} else {
					log.I("Function calls:")
					for _, call := range calls {
						fmt.Printf("    - %s\n", call)
					}
				}
			}
		}
		return true
	})
	log.I("total functions: ", FuncSet.Cardinality())
	return nil
}

func parseFunctionBody(body *ast.BlockStmt) []string {
	var calls []string

	ast.Inspect(body, func(n ast.Node) bool {
		switch e := n.(type) {
		// Check if the node is a function call expression.
		case *ast.CallExpr:
			{
				if fun, ok := e.Fun.(*ast.Ident); ok {
					calls = append(calls, fun.Name)
				} else if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
					// Handle method calls (e.g., fmt.Println)
					if pkg, ok := sel.X.(*ast.Ident); ok {
						calls = append(calls, fmt.Sprintf("%s.%s", pkg.Name, sel.Sel.Name))
					}
				}
			}
		default:
			log.D("unhandled case: ", e)
		}
		return true
	})

	return calls
}

func getProjectPackageList() error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return fmt.Errorf("error loading packages: %v", err)
	}

	// Identify the root module path (useful to filter out external dependencies)
	rootModule := getRootModule(pkgs)

	for _, pkg := range pkgs {
		packagePath := pkg.PkgPath
		log.W(packagePath)
		if strings.HasPrefix(packagePath, rootModule) {
			log.D("[PROJECT PACKAGE] ", packagePath)
		} else {
			log.D("[EXTERNAL PACKAGE] ", packagePath)
		}
	}
	return nil
}

// Helper function to get the root module from the loaded packages
func getRootModule(pkgs []*packages.Package) string {
	for _, pkg := range pkgs {
		if pkg.Module != nil && len(pkg.Module.Path) > 0 {
			log.D("Got project root module", pkg.Module.Path)
			return pkg.Module.Path
		}
	}
	return ""
}
