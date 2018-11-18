// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package anticycle

import (
	"fmt"
	"testing"

	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestExcludeDirs(t *testing.T) {
	custom := []string{"foo", "bar", "baz"}
	excluded := ExcludeDirs(custom)

	expected := []string{".git", ".idea", ".vscode", "bar", "baz", "bin", "dist", "foo", "testdata", "vendor"}
	assert.Equal(t, expected, excluded)
}

func TestExcludeDirs_WithoutCustom(t *testing.T) {
	excluded := ExcludeDirs([]string{})
	assert.Equal(t, DefaultExcluded, excluded)
}

func ExampleExcludeDirs() {
	custom := []string{"foo", "bar", "baz"}
	excluded := ExcludeDirs(custom)
	fmt.Print(excluded)
	// Output: [.git .idea .vscode bar baz bin dist foo testdata vendor]
}

func TestAnticycle_filterAffected(t *testing.T) {
	packages := []*model.Pkg{
		{
			Name: "bar",
			Path: "/tmp/anticycle/skipNotAffected/bar",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/skipNotAffected/baz": {
					Name:      "/tmp/anticycle/skipNotAffected/baz",
					NameShort: "baz",
					Alias:     nil,
				},
				"/tmp/anticycle/skipNotAffected/foo": {
					Name:      "/tmp/anticycle/skipNotAffected/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/skipNotAffected/bar/bar.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/skipNotAffected/baz",
							NameShort: "baz",
							Alias:     nil,
						},
					},
				},
				{
					Path: "/tmp/anticycle/skipNotAffected/bar/clean.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/skipNotAffected/foo",
							NameShort: "baz",
							Alias:     nil,
						},
					},
				},
			},
			HaveCycle: true,
			Cycles: []*model.Cycle{
				{
					AffectedFile: "/tmp/anticycle/skipNotAffected/bar/bar.go",
					AffectedImport: &model.ImportInfo{
						Name:      "/tmp/anticycle/skipNotAffected/baz",
						NameShort: "baz",
						Alias:     nil,
					},
				},
			},
		},
		{
			Name: "baz",
			Path: "/tmp/anticycle/skipNotAffected/baz",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/skipNotAffected/bar": {
					Name:      "/tmp/anticycle/skipNotAffected/bar",
					NameShort: "bar",
					Alias:     nil,
				},
				"/tmp/anticycle/skipNotAffected/foo": {
					Name:      "/tmp/anticycle/skipNotAffected/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/skipNotAffected/baz/baz.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/skipNotAffected/bar",
							NameShort: "bar",
							Alias:     nil,
						},
						{
							Name:      "/tmp/anticycle/skipNotAffected/foo",
							NameShort: "foo",
							Alias:     nil,
						},
					},
				},
			},
			HaveCycle: true,
			Cycles: []*model.Cycle{
				{
					AffectedFile: "/tmp/anticycle/skipNotAffected/baz/baz.go",
					AffectedImport: &model.ImportInfo{
						Name:      "/tmp/anticycle/skipNotAffected/bar",
						NameShort: "bar",
						Alias:     nil,
					},
				},
			},
		},
		{
			Name: "foo",
			Path: "/tmp/anticycle/skipNotAffected/foo",
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/skipNotAffected/foo/foo.go",
				},
			},
		},
	}
	result := onlyAffected(packages)
	assert.Len(t, result, 2)

	for _, pkg := range result {
		t.Run(pkg.Name, func(t *testing.T) {
			assert.Len(t, pkg.Files, 1)
			assert.Len(t, pkg.Files[0].Imports, 1)
			assert.Len(t, pkg.Imports, 1)
			assert.Len(t, pkg.Cycles, 1)
		})
	}
}

func TestSortMetaCycles(t *testing.T) {
	tests := []struct {
		name       string
		metaCycles [][]string
		expected   [][]string
	}{
		{
			name: "bar baz bar",
			metaCycles: [][]string{
				{"bar", "baz", "bar"},
				{"baz", "bar", "baz"},
			},
			expected: [][]string{
				{"bar", "baz", "bar"},
				{"baz", "bar", "baz"},
			},
		},
		{
			name: "baz bar baz",
			metaCycles: [][]string{
				{"baz", "bar", "baz"},
				{"bar", "baz", "bar"},
			},
			expected: [][]string{
				{"bar", "baz", "bar"},
				{"baz", "bar", "baz"},
			},
		},
		{
			name: "a b c ",
			metaCycles: [][]string{
				{"b", "c", "a", "c"},
				{"c", "a", "c"},
				{"a", "c", "b", "c"},
				{"c", "b", "c"},
			},
			expected: [][]string{
				{"a", "c", "b", "c"},
				{"b", "c", "a", "c"},
				{"c", "a", "c"},
				{"c", "b", "c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortMetaCycles(tt.metaCycles)
			assert.Equal(t, tt.expected, result)
		})
	}
}
