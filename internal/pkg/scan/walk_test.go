// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"go/ast"
	"testing"

	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestShouldSkipDir(t *testing.T) {
	excluded := []string{".git", "vendor"}

	assert.True(t, shouldSkip("vendor", excluded))
	assert.True(t, shouldSkip(".git", excluded))
	assert.False(t, shouldSkip("pkg", excluded))
	assert.False(t, shouldSkip("foo", excluded))
}

func TestMakePackages(t *testing.T) {
	expected := []*model.Pkg{
		{
			Name: "foo",
			Path: "internal/pkg/foo",
			Imports: map[string]*model.ImportInfo{
				"pkg/foo": {Name: "pkg/foo", NameShort: "foo", Alias: nil},
			},
			Files: []*model.File{
				{
					Path: "internal/pkg/foo/foo.go",
					Imports: []*model.ImportInfo{
						{Name: "pkg/foo", NameShort: "foo", Alias: nil},
					},
				},
			},
			Cycles: []*model.Cycle{},
		},
	}

	root := map[string]*ast.Package{
		"foo": {
			Files: map[string]*ast.File{
				"internal/pkg/foo/foo.go": {
					Imports: []*ast.ImportSpec{{Path: &ast.BasicLit{Value: `"pkg/foo"`}}},
				},
			},
		},
	}
	packages := newPackages(root, "internal/pkg/foo")
	assert.EqualValues(t, expected, packages)
}

func TestMakePackages_WithEmptyRoot(t *testing.T) {
	root := make(map[string]*ast.Package)
	packages := newPackages(root, "testpath")
	assert.Len(t, packages, 0)
}

func TestWalkDir(t *testing.T) {
	dir, remove := makeProjectNoCycles("walkDir")
	defer remove()

	expected := []string{"bar", "baz", "foo"}
	packages, err := walkDir(dir, []string{})
	assert.NoError(t, err)

	result := make([]string, 0)
	for _, pkg := range packages {
		result = append(result, pkg.Name)
	}
	assert.EqualValues(t, expected, result)
}

func TestWalkDir_WithExcludedDirs(t *testing.T) {
	dir, remove := makeProjectNoCycles("walkDirExcluded")
	defer remove()
	expected := []string{"baz", "foo"}

	packages, err := walkDir(dir, []string{"bar"})
	assert.NoError(t, err)

	result := make([]string, 0)
	for _, pkg := range packages {
		result = append(result, pkg.Name)
	}
	assert.EqualValues(t, expected, result)
}
