// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package anticycle provides high level API for dependencies and cycles analyzing.
package anticycle

import (
	"github.com/anticycle/anticycle/internal/pkg/scan"
	"github.com/anticycle/anticycle/pkg/model"
	"sort"
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
// Directories should have only basic names, not a whole path, see example.
func ExcludeDirs(custom []string) []string {
	excluded := DefaultExcluded
	excluded = append(excluded, custom...)
	sort.Strings(excluded)
	return excluded
}

// Analyze is a high level function which will parse recursively all .go files.
// Will produce slice with information about packages, files and import statements.
// Analyze will not compile the code. It looks only at imports in AST data.
func Analyze(dir string, excludedDir []string, allPkgs bool) ([]*model.Pkg, error) {
	// TODO: make integration tests
	packages, err := scan.FetchPackages(dir, excludedDir)
	if err != nil {
		return nil, err
	}

	cycles, err := scan.FindCycles(packages)
	if err != nil {
		return nil, err
	}

	var result []*model.Pkg
	if allPkgs == false {
		for _, pkg := range cycles {
			if pkg.HaveCycle == true {
				result = append(result, pkg)
			}
		}
	} else {
		result = cycles
	}

	return result, nil
}
