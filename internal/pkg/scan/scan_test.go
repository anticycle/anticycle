// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package scan

import (
	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchPackages(t *testing.T) {
	expected := []*model.Pkg{
		{
			Name: "main",
			Path: "testdata/nocycle",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar",
					NameShort: "bar",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz",
					NameShort: "baz",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/main.go",
					Imports: []string{
						"fmt",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
		{
			Name: "bar",
			Path: "testdata/nocycle/pkg/bar",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/pkg/bar/bar.go",
					Imports: []string{
						"fmt",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
		{
			Name: "baz",
			Path: "testdata/nocycle/pkg/bar/baz",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/pkg/bar/baz/baz.go",
					Imports: []string{
						"fmt",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
		{
			Name: "foo",
			Path: "testdata/nocycle/pkg/foo",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/pkg/foo/foo.go",
					Imports: []string{
						"fmt",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
	}
	packages, err := FetchPackages("testdata/nocycle", []string{})
	assert.NoError(t, err)
	assert.EqualValues(t, expected, packages)
}

func TestFetchPackages_WithExcludedDir(t *testing.T) {
	expected := []*model.Pkg{
		{
			Name: "main",
			Path: "testdata/nocycle",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar",
					NameShort: "bar",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz",
					NameShort: "baz",
					Alias:     nil,
				},
				"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo": {
					Name:      "github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					NameShort: "foo",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/main.go",
					Imports: []string{
						"fmt",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz",
						"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
		{
			Name: "foo",
			Path: "testdata/nocycle/pkg/foo",
			Imports: model.Imports{
				"fmt": {
					Name:      "fmt",
					NameShort: "fmt",
					Alias:     nil,
				},
			},
			Files: []*model.File{
				{
					Path: "testdata/nocycle/pkg/foo/foo.go",
					Imports: []string{
						"fmt",
					},
				},
			},
			Cycles: model.ImportsCycle{},
		},
	}

	excluded := []string{"bar"}
	packages, err := FetchPackages("testdata/nocycle", excluded)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, packages)
}
