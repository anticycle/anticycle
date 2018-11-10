// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package serialize provides API to serialize analyzed output.
package serialize

import (
	"encoding/json"
	"fmt"
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
	output := make([]string, 0, len(packages))
	for _, pkg := range packages {
		fOutput := make([]string, 0)
		cycles := make(map[string]bool, len(pkg.Cycles))
		if pkg.HaveCycle == true {
			for _, cycle := range pkg.Cycles {
				cycles[cycle.AffectedImport.Name] = true
			}
		}

		for _, file := range pkg.Files {
			fOutput = append(fOutput, fmt.Sprintf("  %v", file.Path))
			for _, imp := range file.Imports {
				var format string
				if _, ok := cycles[imp.Name]; ok {
					format = "    C \"%v\""
				} else {
					format = "      \"%v\""
				}
				fOutput = append(fOutput, fmt.Sprintf(format, imp.Name))
			}
		}
		output = append(output, fmt.Sprintf("%v\n%v\n", pkg.Name, strings.Join(fOutput, "\n")))
	}

	return strings.Join(output, "\n"), nil
}
