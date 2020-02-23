package main

import (
	"go/ast"
	"log"
	"sort"
)

type inspector struct {
	typeMap map[string]bool
}

func newInspector() *inspector {
	return &inspector{
		typeMap: make(map[string]bool, 0),
	}
}

func (i *inspector) processPackage(pkg *ast.Package) {
	for fileName, file := range pkg.Files {
		log.Printf("Collecting types from %q", fileName)
		ast.FileExports(file)
		ast.Inspect(file, i.collectTypes)
	}

	log.Printf("Collected %d types in total", len(i.typeMap))

	for _, file := range pkg.Files {
		ast.Inspect(file, i.findExistingUnmarshalers)
	}
}

func (i *inspector) collectTypes(node ast.Node) bool {
	if node == nil {
		return false
	}

	if ts, ok := node.(*ast.TypeSpec); ok {
		i.typeMap[ts.Name.String()] = false
		return false
	}

	return true
}

func (i *inspector) findExistingUnmarshalers(node ast.Node) bool {
	if node == nil {
		return false
	}

	if fd, ok := node.(*ast.FuncDecl); ok {
		if fd.Recv == nil {
			// regular function
			return false
		}

		if fd.Name.Name != "UnmarshalJSON" {
			return false
		}

		recvType := fd.Recv.List[0].Type

		if ident, ok := recvType.(*ast.Ident); ok {
			if _, ok := i.typeMap[ident.Name]; ok {
				i.typeMap[ident.Name] = true
			}
		}

		if se, ok := recvType.(*ast.StarExpr); ok {
			if ident, ok := se.X.(*ast.Ident); ok {
				if _, ok := i.typeMap[ident.Name]; ok {
					i.typeMap[ident.Name] = true
				}
			}
		}
	}
	return true
}

func (i *inspector) typesWithoutUnmarshaler() []string {
	types := make([]string, 0)

	for name, hasUnmarshaler := range i.typeMap {
		if !hasUnmarshaler {
			types = append(types, name)
		}
	}

	sort.Strings(types)

	return types
}
