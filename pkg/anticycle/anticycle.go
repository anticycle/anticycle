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
func Analyze(dir string, excludedDir []string, all bool, related bool) (result []*model.Pkg, err error) {
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
		if !related {
			result = withoutRelated(result)
		}
	}

	return result, err
}

func onlyAffected(packages []*model.Pkg) (result []*model.Pkg) {
	result = make([]*model.Pkg, 0, len(packages))
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
						if len(fileImports) > 0 {
							file.Imports = fileImports
							files = append(files, file)
						}
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

func withoutRelated(packages []*model.Pkg) (result []*model.Pkg) {
	pkgToInt := make(map[string]int, len(packages))
	cycleRefs := make(map[string][]string, len(packages))
	for i, pkg := range packages {
		if !pkg.HaveCycle {
			continue
		}

		pkgToInt[pkg.Name] = i

		if _, ok := cycleRefs[pkg.Name]; !ok {
			cycleRefs[pkg.Name] = make([]string, 0, len(pkg.Cycles))
		}
		for _, cycle := range pkg.Cycles {
			cycleRefs[pkg.Name] = append(cycleRefs[pkg.Name], cycle.AffectedImport.NameShort)
		}
	}

	isInResult := make(map[int]interface{}, len(packages))
	output := make([]int, 0, len(packages))
	if len(cycleRefs) > 1 {
		// Check if it is possible to return to the same package.
		for from, to := range cycleRefs {
			visited := map[string]interface{}{from: struct{}{}}
			for _, imp := range to {
				for _, next := range cycleRefs[imp] {
					if _, ok := visited[next]; ok {
						if _, ok := isInResult[pkgToInt[from]]; !ok {
							output = append(output, pkgToInt[from])
							isInResult[pkgToInt[from]] = struct{}{}
						}
						break
					}
				}
			}
		}
	}

	sort.Ints(output)
	result = make([]*model.Pkg, 0, len(output))
	for _, idx := range output {
		result = append(result, packages[idx])
	}
	return result
}
