// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"github.com/anticycle/anticycle/pkg/model"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

func shouldSkip(pathName string, excluded []string) bool {
	i := sort.SearchStrings(excluded, pathName)
	if i < len(excluded) && excluded[i] == pathName {
		return true
	}
	return false
}

func newPackages(root map[string]*ast.Package, path string) []*model.Pkg {
	packages := make([]*model.Pkg, 0)

	for name, astPkg := range root {
		pkg := model.NewPkg()
		pkg.Name = name
		pkg.Path = path

		for path, astFile := range astPkg.Files {
			file := model.NewFile()
			file.Path = path

			for _, importSpec := range astFile.Imports {
				importInfo := model.NewImportInfo(importSpec)
				file.Imports = append(file.Imports, importInfo)
				pkg.Imports[importInfo.Name] = importInfo
			}
			pkg.Files = append(pkg.Files, file)
		}
		packages = append(packages, pkg)
	}
	return packages
}

func walkDir(dir string, excluded []string) ([]*model.Pkg, error) {
	packages := make([]*model.Pkg, 0, 16)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkip(info.Name(), excluded) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			fset := token.NewFileSet()
			parsedDir, err := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			packages = append(packages, newPackages(parsedDir, path)...)
		}

		return nil
	})

	return packages, err
}
