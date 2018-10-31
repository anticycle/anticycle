// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package serialize

import (
	"fmt"
	"testing"

	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	pkg := model.NewPkg()
	pkg.Name = "test/pkg"
	pkg.Imports = map[string]*model.ImportInfo{"internal": model.NewImportInfo(nil)}

	packages := []*model.Pkg{pkg}
	jsonStr, err := ToJSON(packages)
	assert.NoError(t, err)

	expected := `[{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}]`
	assert.Equal(t, expected, jsonStr)
}

func TestToJSON_WithNilInput(t *testing.T) {
	jsonStr, err := ToJSON(nil)
	assert.NoError(t, err)
	assert.Equal(t, "", jsonStr)
}

func TestToJSON_WithEmptyInput(t *testing.T) {
	jsonStr, err := ToJSON([]*model.Pkg{})
	assert.NoError(t, err)
	assert.Equal(t, "", jsonStr)
}

func ExampleToJSON() {
	pkg := model.NewPkg()
	pkg.Name = "test/pkg"
	pkg.Imports = map[string]*model.ImportInfo{"internal": model.NewImportInfo(nil)}

	jsonStr, _ := ToJSON([]*model.Pkg{pkg})
	fmt.Print(jsonStr)
	// Output: [{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}]
}
