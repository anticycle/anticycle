// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package serialize provides API to serialize analyzed output.
package serialize

import (
	"encoding/json"
	"github.com/anticycle/anticycle/pkg/model"
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
