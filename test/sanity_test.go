package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anticycle/anticycle/pkg/anticycle"
	"github.com/anticycle/anticycle/pkg/model"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

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
			cmd := exec.Command("anticycle", test.args...)
			stdoutStderr, err := cmd.CombinedOutput()
			assert.NoError(t, err)
			if *update {
				ioutil.WriteFile(test.golden, stdoutStderr, 0644)
			}
			if test.isJSON {
				var expected, result []model.Pkg
				golden, _ := ioutil.ReadFile(test.golden)
				json.Unmarshal(stdoutStderr, &result)
				json.Unmarshal(golden, &expected)

				assert.Equal(t, expected, result)
			} else {
				expected, _ := ioutil.ReadFile(test.golden)
				assert.Equal(t, string(expected), string(stdoutStderr))
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
	cases := []struct {
		Name, Golden string
		Args         []string
	}{
		{
			Name:   "Broken import clause",
			Args:   []string{"./testdata/corruptedFiles/brokenImport"},
			Golden: filepath.Join("testdata", "corruptedFiles", "sanity-brokenImport.golden"),
		},
		{
			Name:   "Broken package clause",
			Args:   []string{"./testdata/corruptedFiles/brokenPackage"},
			Golden: filepath.Join("testdata", "corruptedFiles", "sanity-brokenPackage.golden"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cmd := exec.Command("anticycle", c.Args...)
			stdoutStderr, err := cmd.CombinedOutput()
			assert.NotNil(t, err)
			if *update {
				ioutil.WriteFile(c.Golden, stdoutStderr, 0644)
			}
			expected, _ := ioutil.ReadFile(c.Golden)
			assert.Equal(t, expected, stdoutStderr)
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
