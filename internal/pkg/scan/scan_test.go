// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"testing"

	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestFetchPackages(t *testing.T) {
	dir, remove := makeProjectNoCycles("fetchNoCycle")
	defer remove()

	expected := []*model.Pkg{
		{
			Name: "bar",
			Path: "/tmp/anticycle/fetchNoCycle/bar",
			Imports: map[string]*model.ImportInfo{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchNoCycle/bar/bar.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "fmt",
							NameShort: "fmt",
							Alias:     nil,
						},
					},
				},
			},
			Cycles:    []*model.Cycle{},
			HaveCycle: false,
		},
		{
			Name: "baz",
			Path: "/tmp/anticycle/fetchNoCycle/baz",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/fetchNoCycle/bar": {
					Name:      "/tmp/anticycle/fetchNoCycle/bar",
					NameShort: "bar",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchNoCycle/baz/baz.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/fetchNoCycle/bar",
							NameShort: "bar",
							Alias:     nil,
						},
					},
				},
			},
			Cycles:    []*model.Cycle{},
			HaveCycle: false,
		},
		{
			Name: "foo",
			Path: "/tmp/anticycle/fetchNoCycle/foo",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/fetchNoCycle/bar": {
					Name:      "/tmp/anticycle/fetchNoCycle/bar",
					NameShort: "bar",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchNoCycle/foo/foo.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/fetchNoCycle/bar",
							NameShort: "bar",
							Alias:     nil,
						},
					},
				},
			},
			Cycles:    []*model.Cycle{},
			HaveCycle: false,
		},
	}
	packages, err := FetchPackages(dir, []string{})
	assert.NoError(t, err)
	assert.EqualValues(t, expected, packages)
}

func TestFetchPackages_WithExcludedDir(t *testing.T) {
	dir, remove := makeProjectNoCycles("fetchExcludedNoCycle")
	defer remove()

	expected := []*model.Pkg{
		{
			Name: "bar",
			Path: "/tmp/anticycle/fetchExcludedNoCycle/bar",
			Imports: map[string]*model.ImportInfo{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchExcludedNoCycle/bar/bar.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "fmt",
							NameShort: "fmt",
							Alias:     nil,
						},
					},
				},
			},
			Cycles:    []*model.Cycle{},
			HaveCycle: false,
		},
	}
	packages, err := FetchPackages(dir, []string{"baz", "foo"})
	assert.NoError(t, err)
	assert.EqualValues(t, expected, packages)
}

func TestFindCycles_NoCycles(t *testing.T) {
	packages := []*model.Pkg{
		{
			Name: "bar",
			Path: "/tmp/anticycle/fetchNoCycle/bar",
			Imports: map[string]*model.ImportInfo{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchNoCycle/bar/bar.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "fmt",
							NameShort: "fmt",
							Alias:     nil,
						},
					},
				},
			},
		},
		{
			Name: "foo",
			Path: "/tmp/anticycle/fetchNoCycle/foo",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/fetchNoCycle/bar": {
					Name:      "/tmp/anticycle/fetchNoCycle/bar",
					NameShort: "bar",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/fetchNoCycle/foo/foo.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/fetchNoCycle/bar",
							NameShort: "bar",
							Alias:     nil,
						},
					},
				},
			},
		},
	}
	expected := []struct{ files, imports, cycles int }{
		{files: 1, imports: 1, cycles: 0},
		{files: 1, imports: 1, cycles: 0},
	}

	cycles, err := FindCycles(packages)
	assert.NoError(t, err)

	for i, cycle := range cycles {
		t.Run(cycle.Name, func(t *testing.T) {
			assert.Len(t, cycle.Files, expected[i].files)
			assert.Len(t, cycle.Imports, expected[i].imports)
			assert.Len(t, cycle.Cycles, expected[i].cycles)
		})
	}
}

func TestFindCycles(t *testing.T) {
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
	expected := []struct{ files, imports, cycles int }{
		{files: 2, imports: 2, cycles: 1},
		{files: 1, imports: 2, cycles: 1},
		{files: 1, imports: 0, cycles: 0},
	}

	cycles, err := FindCycles(packages)
	assert.NoError(t, err)

	for i, cycle := range cycles {
		t.Run(cycle.Name, func(t *testing.T) {
			assert.Len(t, cycle.Files, expected[i].files)
			assert.Len(t, cycle.Imports, expected[i].imports)
			assert.Len(t, cycle.Cycles, expected[i].cycles)
		})
	}
}

func TestFindCycles_FalsePositiveExternalPkgNameOverlap(t *testing.T) {
	packages := []*model.Pkg{
		{
			Name: "foo",
			Path: "/tmp/anticycle/falsePositive/foo",
			Imports: map[string]*model.ImportInfo{
				"/tmp/anticycle/falsePositive/bar": {
					Name:      "/tmp/anticycle/falsePositive/bar",
					NameShort: "bar",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/falsePositive/foo/foo.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "/tmp/anticycle/falsePositive/bar",
							NameShort: "bar",
							Alias:     nil,
						},
					},
				},
			},
		},
		{
			Name: "bar",
			Path: "/tmp/anticycle/falsePositive/bar",
			Imports: map[string]*model.ImportInfo{
				"github.com/external/fake/foo": {
					Name:      "github.com/external/fake/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "/tmp/anticycle/falsePositive/bar/bar.go",
					Imports: []*model.ImportInfo{
						{
							Name:      "github.com/external/fake/foo",
							NameShort: "foo",
							Alias:     nil,
						},
					},
				},
			},
		},
	}
	expected := []struct{ files, imports, cycles int }{
		{files: 1, imports: 1, cycles: 0},
		{files: 1, imports: 1, cycles: 0},
	}
	cycles, err := FindCycles(packages)
	assert.NoError(t, err)

	for i, cycle := range cycles {
		t.Run(cycle.Name, func(t *testing.T) {
			assert.Len(t, cycle.Files, expected[i].files)
			assert.Len(t, cycle.Imports, expected[i].imports)
			assert.Len(t, cycle.Cycles, expected[i].cycles)
		})
	}
}

func BenchmarkFindCycles_NoCycles(b *testing.B) {
	dir, remove := makeProjectNoCycles("benchFindNoCycle")
	defer remove()
	packages, err := FetchPackages(dir, []string{})
	assert.NoError(b, err)
	assert.Len(b, packages, 3)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindCycles(packages)
	}
}

func BenchmarkFindCycles_OneToOne(b *testing.B) {
	dir, remove := makeProjectOneToOne("benchFindOneToOne")
	defer remove()
	packages, err := FetchPackages(dir, []string{})
	assert.NoError(b, err)
	assert.Len(b, packages, 3)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindCycles(packages)
	}
}

func BenchmarkFindCycles_Triangle(b *testing.B) {
	dir, remove := makeProjectTriangle("benchFindTriangle")
	defer remove()
	packages, err := FetchPackages(dir, []string{})
	assert.NoError(b, err)
	assert.Len(b, packages, 3)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindCycles(packages)
	}
}

func BenchmarkFindCycles_DiagonalSquare(b *testing.B) {
	dir, remove := makeProjectDiagonalSquare("benchFindDiagonalSquare")
	defer remove()
	packages, err := FetchPackages(dir, []string{})
	assert.NoError(b, err)
	assert.Len(b, packages, 4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindCycles(packages)
	}
}
