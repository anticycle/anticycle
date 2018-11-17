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

// ToJSON takes cycle analysis and produces JSON string.
func ToJSON(analysis *model.Analysis) (string, error) {
	jsonBytes, err := json.Marshal(analysis)
	if err != nil {
		return "", err
	}
	output := string(jsonBytes)
	if output == "null" || output == "[]" {
		output = ""
	}
	return output, nil
}

// ToTxt takes cycle analysis and produces human friendly text output.
func ToTxt(analysis *model.Analysis) (string, error) {
	// Packages and Imports order are required for deterministic output.
	// Without it we can't predict result due to dynamic nature of hash maps
	// used to store input data.
	packages := analysis.Cycles
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

	var output strings.Builder
	if len(analysis.Metadata.Cycles) > 0 {
		output.WriteString(fmt.Sprintf("Found %d cycles\n\n", len(analysis.Metadata.Cycles)))
		for _, c := range analysis.Metadata.Cycles {
			output.WriteString(fmt.Sprintf("%s\n", strings.Join(c, " -> ")))
		}
		output.WriteString("\nDetails\n\n")
	}

	// output := make([]string, 0)
	for _, pkg := range pkgOrder {
		var out strings.Builder
		for _, imp := range impsOrder[pkg] {
			files := input[pkg][imp]
			impSegments := strings.Split(imp, "/")
			out.WriteString(fmt.Sprintf("[%s -> %s] \"%s\"\n", pkg, impSegments[len(impSegments)-1], imp))

			var filesList []string
			for _, file := range files {
				filesList = append(filesList, fmt.Sprintf("   %s", file))
			}
			sort.Strings(filesList)
			out.WriteString(fmt.Sprintf("%s\n", strings.Join(filesList, "\n")))
		}
		output.WriteString(fmt.Sprintf("%s\n", out.String()))
	}

	return strings.TrimRight(output.String(), "\r\n"), nil
}

func sliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
