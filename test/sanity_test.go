// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package test

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anticycle/anticycle/pkg/anticycle"
	"github.com/stretchr/testify/assert"
)

func TestNoGoFiles(t *testing.T) {
	cmd := exec.Command("anticycle", "./testdata/empty")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "", string(stdoutStderr))
}

func TestNoCycle(t *testing.T) {
	cmd := exec.Command("anticycle", "./testdata/nocycle")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "", string(stdoutStderr))
}

func TestAnticycle(t *testing.T) {
	tests := []struct {
		isJSON       bool
		name, golden string
		args         []string
	}{
		// no go file in directory scenario
		{
			name:   "No GO files in directory, output all as text",
			args:   []string{"-all", "-format=text", "./testdata/empty"},
			golden: filepath.Join("testdata", "empty", "sanity-all.txt.golden"),
		},
		{
			isJSON: true,
			name:   "No GO files in directory, output all as JSON",
			args:   []string{"-all", "-format=json", "./testdata/empty"},
			golden: filepath.Join("testdata", "empty", "sanity-all.json.golden"),
		},

		// no cycles scenario
		{
			name:   "No cycles, output all as text",
			args:   []string{"-all", "-format=text", "./testdata/nocycle"},
			golden: filepath.Join("testdata", "nocycle", "sanity-all.txt.golden"),
		},
		{
			isJSON: true,
			name:   "No cycles, output all as JSON",
			args:   []string{"-all", "-format=json", "./testdata/nocycle"},
			golden: filepath.Join("testdata", "nocycle", "sanity-all.json.golden"),
		},

		// one-to-one cycle scenario
		{
			name:   "One-to-one cycle, output all as text",
			args:   []string{"-all", "-format=text", "./testdata/onetoone"},
			golden: filepath.Join("testdata", "onetoone", "sanity-all.txt.golden"),
		},
		{
			name:   "One-to-one cycle, output as text",
			args:   []string{"-format=text", "./testdata/onetoone"},
			golden: filepath.Join("testdata", "onetoone", "sanity.txt.golden"),
		},
		{
			isJSON: true,
			name:   "One-to-one cycle, output all as JSON",
			args:   []string{"-all", "-format=json", "./testdata/onetoone"},
			golden: filepath.Join("testdata", "onetoone", "sanity-all.json.golden"),
		},
		{
			isJSON: true,
			name:   "One-to-one cycle, output as JSON",
			args:   []string{"-format=json", "./testdata/onetoone"},
			golden: filepath.Join("testdata", "onetoone", "sanity.json.golden"),
		},

		// not-affected-files cycle scenario
		{
			name:   "Not-affected-files cycle, output as text",
			args:   []string{"-format=text", "./testdata/notAffectedFiles"},
			golden: filepath.Join("testdata", "notAffectedFiles", "sanity.txt.golden"),
		},
		{
			isJSON: true,
			name:   "Not-affected-files cycle, output as JSON",
			args:   []string{"-format=json", "./testdata/notAffectedFiles"},
			golden: filepath.Join("testdata", "notAffectedFiles", "sanity.json.golden"),
		},

		// external false positive cycle scenario
		{
			name:   "External false positive, output as text",
			args:   []string{"-format=text", "./testdata/externalFalsePositive"},
			golden: filepath.Join("testdata", "externalFalsePositive", "sanity.txt.golden"),
		},
		{
			isJSON: true,
			name:   "External false positive, output as JSON",
			args:   []string{"-format=json", "./testdata/externalFalsePositive"},
			golden: filepath.Join("testdata", "externalFalsePositive", "sanity.json.golden"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, err := exec.Command("anticycle", test.args...).Output()
			assert.NoError(t, err)
			if *update {
				updateGolden(test.golden, stdOut)
			}
			if test.isJSON {
				var expected, result map[string]interface{}
				golden := readGolden(test.golden)
				if err := json.Unmarshal(stdOut, &result); err != nil {
					panic(err)
				}
				if err := json.Unmarshal(golden, &expected); err != nil {
					panic(err)
				}

				assert.Equal(t, expected, result)
			} else {
				expected := readGolden(test.golden)
				assert.Equal(t, string(expected), string(stdOut))
			}
		})
	}
}

func TestFalsePositiveExternalPackage(t *testing.T) {
	cmd := exec.Command("anticycle", "./testdata/externalFalsePositive")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "", string(stdoutStderr)) // should not find any cycles
}

func TestDefaultExclude(t *testing.T) {
	cmd := exec.Command("anticycle", "-all", "./testdata/defaultExclude")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "", string(stdoutStderr))
}

func TestCorruptedFiles(t *testing.T) {
	tests := []struct {
		name, golden string
		args         []string
	}{
		{
			name:   "Broken import clause",
			args:   []string{"./testdata/corruptedFiles/brokenImport"},
			golden: filepath.Join("testdata", "corruptedFiles", "sanity-brokenImport.golden"),
		},
		{
			name:   "Broken package clause",
			args:   []string{"./testdata/corruptedFiles/brokenPackage"},
			golden: filepath.Join("testdata", "corruptedFiles", "sanity-brokenPackage.golden"),
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
			assert.Equal(t, expected, stdErr)
		})
	}
}

func TestShowExclude(t *testing.T) {
	cmd := exec.Command("anticycle", "-showExcluded")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	expected := fmt.Sprintf("%v\n", strings.Join(anticycle.DefaultExcluded, "\n"))
	assert.Equal(t, expected, string(stdoutStderr))
}
