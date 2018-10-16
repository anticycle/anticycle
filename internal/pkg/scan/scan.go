// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"github.com/anticycle/anticycle/pkg/model"
	"strings"
)

func FetchPackages(dir string, excluded []string) ([]*model.Pkg, error) {
	packages, err := walkDir(dir, excluded)
	if err != nil {
		return nil, err
	}
	return packages, nil
}

func FindCycles(packages []*model.Pkg) ([]*model.Pkg, error) {
	// 2D array where rows are start nodes and columns are end nodes
	pkgToInt := make(map[string]int)
	for idx, pkg := range packages {
		pkgToInt[pkg.Name] = idx
	}

	// allocate composed 2d slice
	nodes := len(packages)
	graph := make([][]uint16, nodes, nodes)
	for i := range graph {
		graph[i] = make([]uint16, nodes, nodes)
	}

	// fill graph with nodes
	for idx, pkg := range packages {
		for _, imp := range pkg.Imports {
			if impIdx, ok := pkgToInt[imp.NameShort]; ok {
				graph[idx][impIdx] = 1
			}
		}
	}

	// Apply Roy-Warshall algorithm to find cycles
	for x := 0; x < nodes; x++ {
		for y := 0; y < nodes; y++ {
			for z := 0; z < nodes; z++ {
				if graph[y][z] == 0 {
					graph[y][z] = graph[y][x] * graph[x][z]
				}
			}
		}
	}

	// mark cycles
	for i := 0; i < nodes; i++ {
		// cycle when graph[x][x] is 1
		if graph[i][i] == 1 {
			packages[i].HaveCycle = true
			// look at links to other packages
			for j := 0; j < nodes; j++ {
				if j != i && graph[i][j] == 1 {
					// check which file is affected
					for _, file := range packages[i].Files {
						for _, imp := range file.Imports {
							if strings.Contains(imp.NameShort, packages[j].Name) {
								cycle := &model.Cycle{
									AffectedFile:   file.Path,
									AffectedImport: imp,
								}
								packages[i].Cycles = append(packages[i].Cycles, cycle)
							}
						}
					}
				}
			}
		}
	}

	return packages, nil
}
