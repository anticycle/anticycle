// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package serialize

import (
	"fmt"
	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToJson(t *testing.T) {
	pkg := model.NewPkg()
	pkg.Name = "test/pkg"
	pkg.Imports = map[string]*model.ImportInfo{"internal": model.NewImportInfo(nil)}

	packages := []*model.Pkg{pkg}
	jsonStr, err := ToJson(packages)
	assert.NoError(t, err)

	expected := `[{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}]`
	assert.Equal(t, expected, jsonStr)
}

func TestToJson_WithNilInput(t *testing.T) {
	jsonStr, err := ToJson(nil)
	assert.NoError(t, err)
	assert.Equal(t, "", jsonStr)
}

func TestToJson_WithEmptyInput(t *testing.T) {
	jsonStr, err := ToJson([]*model.Pkg{})
	assert.NoError(t, err)
	assert.Equal(t, "", jsonStr)
}

func ExampleToJson() {
	pkg := model.NewPkg()
	pkg.Name = "test/pkg"
	pkg.Imports = map[string]*model.ImportInfo{"internal": model.NewImportInfo(nil)}

	jsonStr, _ := ToJson([]*model.Pkg{pkg})
	fmt.Print(jsonStr)
	// Output: [{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}]
}
