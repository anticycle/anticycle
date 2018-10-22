// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package serialize provides API to serialize analyzed output.
package serialize

import (
	"encoding/json"
	"fmt"
	"github.com/anticycle/anticycle/pkg/model"
	"sort"
	"strings"
)

// ToJson takes list of packages and produces JSON string.
func ToJson(packages []*model.Pkg) (string, error) {
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

func isCycle(filePath string, cycles []string) bool {
	i := sort.SearchStrings(cycles, filePath)
	if i < len(cycles) && cycles[i] == filePath {
		return true
	}
	return false
}

func ToTxt(packages []*model.Pkg) (string, error) {
	output := make([]string, 0, len(packages))
	for _, pkg := range packages {
		fOutput := make([]string, 0)
		cycles := make([]string, 0, len(pkg.Cycles))
		if pkg.HaveCycle == true {
			for _, cycle := range pkg.Cycles {
				cycles = append(cycles, cycle.AffectedFile)
			}
			sort.Strings(cycles)
		}

		for _, file := range pkg.Files {
			fOutput = append(fOutput, fmt.Sprintf("  %v", file.Path))
			for _, imp := range file.Imports {
				format := "    \"%v\""
				if isCycle(file.Path, cycles) == true {
					format = "    C \"%v\""
				}
				fOutput = append(fOutput, fmt.Sprintf(format, imp.Name))
			}
		}
		output = append(output, fmt.Sprintf("%v\n%v\n", pkg.Name, strings.Join(fOutput, "\n")))
	}

	return strings.Join(output, "\n"), nil
}
