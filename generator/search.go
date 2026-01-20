package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// FindInterfaces searches for Go interfaces in the specified directory.
func FindInterfaces(root string) (map[string]string, error) {
	interfaces := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return nil
		}

		for _, decl := range node.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							interfaces[typeSpec.Name.Name] = path
						}
					}
				}
			}
		}
		return nil
	})

	return interfaces, err
}
