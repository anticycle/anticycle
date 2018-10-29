package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anticycle/anticycle/pkg/anticycle"
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

func TestShowAllPackages(t *testing.T) {
	cases := []struct {
		Name, Golden string
		Args         []string
	}{
		// no go file in directory scenario
		{
			Name:   "No GO files in directory, output all as text",
			Args:   []string{"-all", "-format=text", "./testdata/empty"},
			Golden: filepath.Join("testdata", "empty", "sanity.txt.golden"),
		},
		{
			Name:   "No GO files in directory, output all as JSON",
			Args:   []string{"-all", "-format=json", "./testdata/empty"},
			Golden: filepath.Join("testdata", "empty", "sanity.json.golden"),
		},

		// no cycles scenario
		{
			Name:   "No cycles, output all as text",
			Args:   []string{"-all", "-format=text", "./testdata/nocycle"},
			Golden: filepath.Join("testdata", "nocycle", "sanity.txt.golden"),
		},
		{
			Name:   "No cycles, output all as JSON",
			Args:   []string{"-all", "-format=json", "./testdata/nocycle"},
			Golden: filepath.Join("testdata", "nocycle", "sanity.json.golden"),
		},

		// one-to-one cycle scenario
		{
			Name:   "One-to-one cycle, output all as text",
			Args:   []string{"-all", "-format=text", "./testdata/onetoone"},
			Golden: filepath.Join("testdata", "onetoone", "sanity.txt.golden"),
		},
		{
			Name:   "One-to-one cycle, output all as JSON",
			Args:   []string{"-all", "-format=json", "./testdata/onetoone"},
			Golden: filepath.Join("testdata", "onetoone", "sanity.json.golden"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cmd := exec.Command("anticycle", c.Args...)
			stdoutStderr, err := cmd.CombinedOutput()
			assert.NoError(t, err)
			if *update {
				ioutil.WriteFile(c.Golden, stdoutStderr, 0644)
			}
			expected, _ := ioutil.ReadFile(c.Golden)
			assert.Equal(t, expected, stdoutStderr)
		})
	}
}

func TestFalsePositiveExternalPackage(t *testing.T) {
	t.Skip("Skip until issue #10 is not resolved")
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
	cmd := exec.Command("anticycle", "-showExclude")
	stdoutStderr, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	expected := fmt.Sprintf("%v\n", strings.Join(anticycle.DefaultExcluded, "\n"))
	assert.Equal(t, expected, string(stdoutStderr))
}
