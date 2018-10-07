// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package anticycle

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
