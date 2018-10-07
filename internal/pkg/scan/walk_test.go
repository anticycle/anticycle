// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestShouldSkipDir(t *testing.T) {
	excluded := []string{".git", "vendor"}

	assert.True(t, shouldSkip("vendor", excluded))
	assert.True(t, shouldSkip(".git", excluded))
	assert.False(t, shouldSkip("pkg", excluded))
	assert.False(t, shouldSkip("foo", excluded))
}

func TestMakePackages(t *testing.T) {
	root := map[string]*ast.Package{
		"foo": {
			Files: map[string]*ast.File{
				"internal/pkg/foo/foo.go": {
					Imports: []*ast.ImportSpec{
						{
							Path: &ast.BasicLit{Value: `"pkg/foo"`},
						},
					},
				},
			},
		},
	}

	packages := newPackages(root, "internal/pkg/foo")
	expected := []*model.Pkg{
		{
			Name: "foo",
			Path: "internal/pkg/foo",
			Imports: model.Imports{
				"pkg/foo": {
					Name:      "pkg/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path:    "internal/pkg/foo/foo.go",
					Imports: []string{"pkg/foo"},
				},
			},
			Cycles: make(model.ImportsCycle),
		},
	}

	assert.EqualValues(t, expected, packages)
}

func TestMakePackages_WithEmptyRoot(t *testing.T) {
	root := make(map[string]*ast.Package)
	packages := newPackages(root, "testpath")
	assert.Len(t, packages, 0)
}

func TestWalkDir(t *testing.T) {
	packages, err := walkDir("testdata/nocycle", []string{})
	assert.NoError(t, err)

	result := make([]string, 0)
	for _, pkg := range packages {
		result = append(result, pkg.Name)
	}
	expected := []string{"main", "bar", "baz", "foo"}
	assert.EqualValues(t, expected, result)
}

func TestWalkDir_WithExcludedDirs(t *testing.T) {
	excluded := []string{"bar"}
	packages, err := walkDir("testdata/nocycle", excluded)
	assert.NoError(t, err)

	result := make([]string, 0)
	for _, pkg := range packages {
		result = append(result, pkg.Name)
	}
	expected := []string{"main", "foo"}
	assert.EqualValues(t, expected, result)
}
