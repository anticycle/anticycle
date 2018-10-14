// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestNewFile(t *testing.T) {
	file := NewFile()
	assert.Empty(t, file.Imports)
	assert.IsType(t, []string{}, file.Imports)
}

func TestNewPkg(t *testing.T) {
	pkg := NewPkg()
	assert.Empty(t, pkg.Imports)
	assert.IsType(t, Imports{}, pkg.Imports)

	assert.Empty(t, pkg.Files)
	assert.IsType(t, []*File{}, pkg.Files)

	assert.Empty(t, pkg.Cycles)
	assert.IsType(t, ImportsCycle{}, pkg.Cycles)
}

func TestNewImportInfo(t *testing.T) {
	astImport := &ast.ImportSpec{
		Path: &ast.BasicLit{Value: `"test/pkg"`},
	}
	importInfo := NewImportInfo(astImport)
	assert.Equal(t, "test/pkg", importInfo.Name)
	assert.Equal(t, "pkg", importInfo.NameShort)
	assert.Nil(t, importInfo.Alias)
}

func TestNewImportInfo_WithAlias(t *testing.T) {
	astImport := &ast.ImportSpec{
		Name: &ast.Ident{Name: `"packaging"`},
		Path: &ast.BasicLit{Value: `"test/pkg"`},
	}
	importInfo := NewImportInfo(astImport)
	assert.Equal(t, "packaging", *importInfo.Alias)
}

func TestNewImportInfo_IfAstImportSpecIsNil(t *testing.T) {
	importInfo := NewImportInfo(nil)
	assert.Nil(t, importInfo)
}
