// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package serialize provides API to serialize analyzed output.
package serialize

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/anticycle/anticycle/pkg/model"
)

// ToJSON takes list of packages and produces JSON string.
func ToJSON(packages []*model.Pkg) (string, error) {
	jsonBytes, err := json.Marshal(packages)
	if err != nil {
		return "", err
	}
	output := string(jsonBytes)
	if output == "null" || output == "[]" {
		output = ""
	}
	return output, nil
}

// ToTxt takes list of packages and produces human friendly text output.
func ToTxt(packages []*model.Pkg) (string, error) {
	// Packages and Imports order are required for deterministic output.
	// Without it we can't predict result due to dynamic nature of hash maps
	// used to store input data.
	pkgOrder := make([]string, 0, len(packages))
	impsOrder := make(map[string][]string)

	// To easy create list of files grouped by imports per package
	// we need to create input data based on slice of packages.
	input := make(map[string]map[string][]string)
	for _, pkg := range packages {
		pkgOrder = append(pkgOrder, pkg.Name)
		impsOrder[pkg.Name] = make([]string, 0)

		if _, ok := input[pkg.Name]; !ok {
			input[pkg.Name] = make(map[string][]string)
			// Prepare upfront all imports with empty list of files
			for imp := range pkg.Imports {
				input[pkg.Name][imp] = make([]string, 0)
			}
		}
		// append all files to corresponding imports without duplicates
		for _, file := range pkg.Files {
			for _, imp := range file.Imports {
				if !sliceContains(impsOrder[pkg.Name], imp.Name) {
					impsOrder[pkg.Name] = append(impsOrder[pkg.Name], imp.Name)
				}
				if sliceContains(input[pkg.Name][imp.Name], file.Path) {
					continue
				}
				input[pkg.Name][imp.Name] = append(input[pkg.Name][imp.Name], file.Path)
			}
		}
	}

	output := make([]string, 0)
	for _, pkg := range pkgOrder {
		var out []string
		for _, imp := range impsOrder[pkg] {
			files := input[pkg][imp]
			impSegments := strings.Split(imp, "/")
			out = append(out, fmt.Sprintf("[%s -> %s] \"%s\"\n", pkg, impSegments[len(impSegments)-1], imp))

			var filesList []string
			for _, file := range files {
				filesList = append(filesList, fmt.Sprintf("   %s", file))
			}
			sort.Strings(filesList)
			out = append(out, strings.Join(filesList, "\n"), "\n")
		}
		output = append(output, strings.Join(out, ""))
	}

	return strings.Join(output, "\n"), nil
}

func sliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
