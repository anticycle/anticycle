// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package model contains structures and constructors used by anticycle analyze functions.
package model

import (
	"go/ast"
	"path"
	"strings"
)

type (
	// ImportInfo holds information about import statements.
	ImportInfo struct {
		Name      string  `json:"name"`
		NameShort string  `json:"nameShort"`
		Alias     *string `json:"alias"`
	}

	// File is a representation of source file with its path and list of imports.
	File struct {
		Path    string   `json:"path"`
		Imports []string `json:"imports"`
	}

	// Pkg is a higher level structure which has all information about its files and imports.
	Pkg struct {
		Name    string       `json:"name"`
		Path    string       `json:"path"`
		Imports Imports      `json:"imports"`
		Files   []*File      `json:"files"`
		Cycles  ImportsCycle `json:"cycles,omitempty"`
	}

	// Imports type holds information about package level imports fetched from all files in a package.
	Imports      map[string]*ImportInfo
	// ImportCycle is a back reference to a file which takes a part in cycle.
	ImportsCycle map[string][]*File
)

// NewPkg creates new Pkg with empty imports, files and cycles arrays.
// Will return pointer to Pkg structure.
func NewPkg() *Pkg {
	return &Pkg{
		Imports: make(Imports),
		Files:   make([]*File, 0, 8),
		Cycles:  make(ImportsCycle),
	}
}

// NewFile creates new File with empty imports array.
// Will return pointer to File structure.
func NewFile() *File {
	return &File{
		Imports: make([]string, 0, 8),
	}
}

// NewImportInfo creates new ImportInfo from ast.ImportSpec.
// Will produce full name, short name and alias if exist.
func NewImportInfo(imp *ast.ImportSpec) *ImportInfo {
	if imp == nil {
		return nil
	}
	name := strings.Trim(imp.Path.Value, `"`)
	shortName := path.Base(name)
	info := &ImportInfo{
		Name:      name,
		NameShort: shortName,
	}
	if imp.Name != nil {
		alias := strings.Trim(imp.Name.Name, `"`)
		info.Alias = &alias
	}
	return info
}
