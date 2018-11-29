// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnticycleExcludeDirs(t *testing.T) {
	tests := []struct {
		name, golden string
		args         []string
	}{
		{
			name:   "No exclusion arguments",
			args:   []string{"-format=text", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "no-flags.golden"),
		},
		{
			name:   "Exclude left package",
			args:   []string{"-format=text", "-exclude='left'", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "ex-left.golden"),
		},
		{
			name:   "Exclude is empty",
			args:   []string{"-format=text", "-exclude=''", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "ex-empty.golden"),
		},

		// next tests require -all flag to actually see something
		{
			name:   "Exclude top and bottom package",
			args:   []string{"-format=text", "-all", "-exclude='top bottom'", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "ex-top_bottom.golden"),
		},
		{
			name:   "Disable default excluded",
			args:   []string{"-format=text", "-all", "-excludeDefault=''", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "disable-default-ex.golden"),
		},
		{
			name:   "Both exclude and excludeDefault are empty",
			args:   []string{"-format=text", "-all", "-exclude=''", "-excludeDefault=''", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "both-ex-empty.golden"),
		},
		{
			name:   "Exclude all packages",
			args:   []string{"-format=text", "-all", "-exclude='top bottom left right cmd'", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "ex-all.golden"),
		},
		{
			name:   "Both exclude and excludeDefault are the same",
			args:   []string{"-format=text", "-all", "-exclude='top bottom cmd'", "-excludeDefault='top bottom cmd'", "./testdata/excludeDirs"},
			golden: filepath.Join("testdata", "excludeDirs", "both-ex-the-same.golden"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, err := exec.Command("anticycle", test.args...).Output()
			assert.NoError(t, err)
			if *update {
				updateGolden(test.golden, stdOut)
			}
			expected := readGolden(test.golden)
			assert.Equal(t, string(expected), string(stdOut))
		})
	}
}
