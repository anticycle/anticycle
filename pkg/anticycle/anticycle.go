// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package anticycle provides high level API for dependencies and cycles analyzing.
package anticycle

import (
	"sort"

	"github.com/anticycle/anticycle/internal/pkg/scan"
	"github.com/anticycle/anticycle/pkg/model"
)

// DefaultExcluded is a default list of directories which should be skipped while analyzing.
var DefaultExcluded = []string{
	".git",
	".idea",
	".vscode",
	"bin",
	"dist",
	"testdata",
	"vendor",
}

// ExcludeDirs takes list of directories excluded by default and appends directories defined by user.
// Directories should have only basic, catalogue names, see example.
func ExcludeDirs(custom []string) []string {
	excluded := DefaultExcluded
	excluded = append(excluded, custom...)
	sort.Strings(excluded)
	return excluded
}

// Analyze is a high level function which will parse recursively all .go files.
// Will produce slice with information about packages, files and import statements.
// Analyze will not compile the code. It looks only at imports in AST data.
func Analyze(dir string, excludedDir []string, all bool) (result []*model.Pkg, err error) {
	// TODO: make integration tests
	packages, err := scan.FetchPackages(dir, excludedDir)
	if err != nil {
		return result, err
	}

	cycles, err := scan.FindCycles(packages)
	if err != nil {
		return result, err
	}

	if all {
		result = cycles
	} else {
		result = onlyAffected(cycles)
	}
	return result, err
}

func onlyAffected(packages []*model.Pkg) []*model.Pkg {
	var result []*model.Pkg
	for _, pkg := range packages {
		if pkg.HaveCycle {

			imports := make(map[string]*model.ImportInfo, len(pkg.Cycles))
			for _, cycle := range pkg.Cycles {
				if imp, ok := pkg.Imports[cycle.AffectedImport.Name]; ok {
					imports[cycle.AffectedImport.Name] = imp
					break
				}
			}
			pkg.Imports = imports

			files := make([]*model.File, 0, len(pkg.Cycles))
			for _, cycle := range pkg.Cycles {
				for _, file := range pkg.Files {
					if file.Path == cycle.AffectedFile {
						fileImports := make([]*model.ImportInfo, 0, len(pkg.Cycles))
						for _, imp := range file.Imports {
							if _, ok := imports[imp.Name]; ok {
								fileImports = append(fileImports, imp)
							}
						}
						file.Imports = fileImports
						files = append(files, file)
						break
					}
				}
			}
			pkg.Files = files
			result = append(result, pkg)
		}
	}
	return result
}
