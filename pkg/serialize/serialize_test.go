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
	analysis := &model.Analysis{
		Cycles:   []*model.Pkg{pkg},
		Metadata: &model.AnalysisMeta{Cycles: [][]string{}},
	}

	jsonStr, err := ToJSON(analysis)
	assert.NoError(t, err)

	expected := `{"cycles":[{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}],"metadata":{"cycles":[]}}`
	assert.Equal(t, expected, jsonStr)
}

func TestToJSON_WithEmptyInput(t *testing.T) {
	analysis := &model.Analysis{
		Cycles:   []*model.Pkg{},
		Metadata: &model.AnalysisMeta{Cycles: [][]string{}},
	}

	jsonStr, err := ToJSON(analysis)
	assert.NoError(t, err)
	assert.Equal(t, "{\"cycles\":[],\"metadata\":{\"cycles\":[]}}", jsonStr)
}

func ExampleToJSON() {
	pkg := model.NewPkg()
	pkg.Name = "test/pkg"
	pkg.Imports = map[string]*model.ImportInfo{"internal": model.NewImportInfo(nil)}
	analysis := &model.Analysis{
		Cycles:   []*model.Pkg{pkg},
		Metadata: &model.AnalysisMeta{Cycles: [][]string{}},
	}

	jsonStr, _ := ToJSON(analysis)
	fmt.Print(jsonStr)
	// Output: {"cycles":[{"name":"test/pkg","path":"","imports":{"internal":null},"files":[],"haveCycle":false}],"metadata":{"cycles":[]}}
}

func TestToTxt(t *testing.T) {
	analysis := &model.Analysis{
		Cycles: []*model.Pkg{
			{
				Name: "bar",
				Path: "/tmp/anticycle/skipNotAffected/bar",
				Imports: map[string]*model.ImportInfo{
					"/tmp/anticycle/skipNotAffected/baz": {
						Name:      "/tmp/anticycle/skipNotAffected/baz",
						NameShort: "baz",
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
		},
		Metadata: &model.AnalysisMeta{Cycles: [][]string{
			{"bar", "baz", "bar"},
			{"baz", "bar", "baz"},
		}},
	}
	expected := "Found 2 cycles\n\nbar -> baz -> bar\nbaz -> bar -> baz\n\nDetails\n\n[bar -> baz] \"/tmp/anticycle/skipNotAffected/baz\"\n   /tmp/anticycle/skipNotAffected/bar/bar.go\n\n[baz -> bar] \"/tmp/anticycle/skipNotAffected/bar\"\n   /tmp/anticycle/skipNotAffected/baz/baz.go"

	result, err := ToTxt(analysis)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
