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

func TestToTxt(t *testing.T) {
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
							NameShort: "foo",
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
	expected := "[bar -> baz] \"/tmp/anticycle/skipNotAffected/baz\"\n   /tmp/anticycle/skipNotAffected/bar/bar.go\n[bar -> foo] \"/tmp/anticycle/skipNotAffected/foo\"\n   /tmp/anticycle/skipNotAffected/bar/clean.go\n\n[baz -> bar] \"/tmp/anticycle/skipNotAffected/bar\"\n   /tmp/anticycle/skipNotAffected/baz/baz.go\n[baz -> foo] \"/tmp/anticycle/skipNotAffected/foo\"\n   /tmp/anticycle/skipNotAffected/baz/baz.go\n\n"

	result, err := ToTxt(packages)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
