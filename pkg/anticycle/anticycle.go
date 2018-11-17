// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package anticycle provides high level API for dependencies and cycles analyzing.
package anticycle

import (
	"sort"
	"strings"

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

// Collect is a high level function which will parse recursively all .go files
// and capture information about packages, files and import statements.
// Collect will not compile the code. It looks only at imports in AST data.
func Collect(dir string, excludedDir []string, all bool) (result []*model.Pkg, err error) {
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

// Analyze goes through collected packages and computes metadata.
func Analyze(packages []*model.Pkg) *model.Analysis {
	if packages == nil {
		packages = make([]*model.Pkg, 0, 0)
	}
	analysis := &model.Analysis{
		Metadata: &model.AnalysisMeta{
			Cycles: make([][]string, 0, len(packages)*2),
		},
		Cycles: packages,
	}
	if len(packages) == 0 {
		return analysis
	}

	cycleRefs := make(map[string][]string, len(packages))
	for _, pkg := range packages {
		if !pkg.HaveCycle {
			continue
		}

		if _, ok := cycleRefs[pkg.Name]; !ok {
			cycleRefs[pkg.Name] = make([]string, 0, len(pkg.Cycles))
		}
		for _, cycle := range pkg.Cycles {
			cycleRefs[pkg.Name] = append(cycleRefs[pkg.Name], cycle.AffectedImport.NameShort)
		}
	}
	if len(cycleRefs) > 0 {
		for from := range cycleRefs {
			visited := make([]string, 0)
			visited = append(visited, walk(from, cycleRefs, visited)...)
			analysis.Metadata.Cycles = append(analysis.Metadata.Cycles, visited)
		}
		analysis.Metadata.Cycles = sortMetaCycles(analysis.Metadata.Cycles)
	}

	return analysis
}

// To provide deterministic output sort metadata cycles by first element of each cycle.
func sortMetaCycles(metaCycles [][]string) [][]string {
	sort.SliceStable(metaCycles, func(i, j int) bool {
		left := strings.Join(metaCycles[i], "")
		right := strings.Join(metaCycles[j], "")

		pair := []string{left, right}
		sort.Strings(pair)
		return pair[0] == left
	})
	return metaCycles
}

func walk(next string, refs map[string][]string, visited []string) []string {
	for _, vi := range visited {
		if vi == next {
			visited = append(visited, next)
			return visited
		}
	}
	visited = append(visited, next)
	return walk(refs[next][0], refs, visited)
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
