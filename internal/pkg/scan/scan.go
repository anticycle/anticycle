// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"github.com/anticycle/anticycle/pkg/model"
)

func FetchPackages(dir string, excluded []string) ([]*model.Pkg, error) {
	packages, err := walkDir(dir, excluded)
	if err != nil {
		return nil, err
	}
	return packages, nil
}

func FindCycles(packages []*model.Pkg) ([]*model.Pkg, error) {
	// TODO implement find cycles
	return packages, nil
}
