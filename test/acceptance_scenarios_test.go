// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnticycleScenarioCorruptedFiles(t *testing.T) {
	tests := []struct {
		name, golden string
		args         []string
	}{
		{
			name:   "Broken import statement in txt format",
			args:   []string{"-format=text", "-exclude=brokenPackage", "./testdata/corruptedFiles"},
			golden: filepath.Join("testdata", "corruptedFiles", "broken-import.golden"),
		},
		{
			name:   "Broken import statement in json format",
			args:   []string{"-format=json", "-exclude=brokenPackage", "./testdata/corruptedFiles"},
			golden: filepath.Join("testdata", "corruptedFiles", "broken-import.golden"),
		},
		{
			name:   "Broken package statement in txt format",
			args:   []string{"-format=text", "-exclude=brokenImport", "./testdata/corruptedFiles"},
			golden: filepath.Join("testdata", "corruptedFiles", "broken-package.golden"),
		},
		{
			name:   "Broken package statement in json format",
			args:   []string{"-format=json", "-exclude=brokenImport", "./testdata/corruptedFiles"},
			golden: filepath.Join("testdata", "corruptedFiles", "broken-package.golden"),
		},
		{
			name:   "Fail on first error should fail on broken import (alphabetical order)",
			args:   []string{"./testdata/corruptedFiles"},
			golden: filepath.Join("testdata", "corruptedFiles", "fail-on-first.golden"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdErr, err := exec.Command("anticycle", test.args...).CombinedOutput()
			assert.Error(t, err)
			if *update {
				updateGolden(test.golden, stdErr)
			}
			expected := readGolden(test.golden)
			assert.Equal(t, string(expected), string(stdErr))
		})
	}
}

func TestAnticycleScenarioMatrix(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "Diagonal scenario",
			testdata: "diagonal",
		},
		{
			name:     "No Cycle scenario",
			testdata: "nocycle",
		},
		{
			name:     "One to One scenario",
			testdata: "onetoone",
		},
		{
			name:     "Triangle scenario",
			testdata: "triangle",
		},
	}

	baseTests := []testCase{
		{
			name:   "get all packages in %s format",
			args:   []string{"-all"},
			golden: "all",
		},
		{
			name:   "without foo in %s format",
			args:   []string{"-exclude=foo"},
			golden: "without-foo",
		},
		{
			name:   "without default foo in %s format",
			args:   []string{"-excludeDefault=foo"},
			golden: "without-default-foo",
		},
		{
			name:   "all but foo in %s format",
			args:   []string{"-all", "-exclude=foo"},
			golden: "all-but-foo",
		},
		{
			name:   "without all in %s format",
			args:   []string{"-exclude='foo bar baz pas'"},
			golden: "without-all",
		},
	}

	tests := make([]testCase, 0, len(baseTests)*2)
	for _, tt := range baseTests {
		tests = append(tests,
			testCase{
				name:   "%s " + fmt.Sprintf(tt.name, "text"),
				args:   append(tt.args, "-format=text"),
				golden: tt.golden + ".txt.golden",
				isJSON: false,
			},
			testCase{
				name:   "%s " + fmt.Sprintf(tt.name, "JSON"),
				args:   append(tt.args, "-format=json"),
				golden: tt.golden + ".json.golden",
				isJSON: true,
			},
		)

	}

	for _, scenario := range scenarios {
		for _, test := range tests {
			runTestGolden(t, scenario, test)
		}
	}
}
