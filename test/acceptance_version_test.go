package test

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fullVersionRegex = `^Anticycle version v\d\.\d\.\d(-(alpha|beta)(-\d+)?)?, build .{7}\n$`
	versionRegex     = `^v\d\.\d\.\d(-(alpha|beta)(-\d+)?)?$`
	buildRegex       = `^.{7}$`
)

func TestAnticycleVersion_ShouldPrintVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Show only version",
			args: []string{"-version"},
		},
		{
			name: "Show version and exclude",
			args: []string{"-version", "-exclude=''"},
		},
		{
			name: "Show version and excludeOnly",
			args: []string{"-version", "-excludeDefault=''"},
		},
		{
			name: "Show version and exclude and excludeOnly",
			args: []string{"-version", "-exclude=''", "-excludeDefault=''"},
		},
		{
			name: "Show version and showExclude",
			args: []string{"-version", "-showExcluded"},
		},
		{
			name: "Show version and all",
			args: []string{"-version", "-all"},
		},
		{
			name: "Show version with directory",
			args: []string{"-version", "./testdata"},
		},
		{
			name: "Show version with everything",
			args: []string{"-version", "-exclude=''", "-excludeDefault=''", "-all", "./testdata"},
		},
	}
	validVersion := regexp.MustCompile(fullVersionRegex)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, err := exec.Command("anticycle", test.args...).Output()
			assert.NoError(t, err)
			assert.True(t, validVersion.MatchString(string(stdOut)))
		})
	}
}

func TestAnticycleVersionJSONFormat_ShouldPrintInJSON(t *testing.T) {
	stdOut, err := exec.Command("anticycle", "-version", "-format=json").Output()
	assert.NoError(t, err)

	var output map[string]string
	err = json.Unmarshal(stdOut, &output)
	assert.NoError(t, err)

	assert.Len(t, output, 2)
	version, ok := output["version"]
	assert.True(t, ok)
	validVersion := regexp.MustCompile(versionRegex)
	assert.True(t, validVersion.MatchString(version))

	build, ok := output["build"]
	assert.True(t, ok)
	validBuild := regexp.MustCompile(buildRegex)
	assert.True(t, validBuild.MatchString(build))
}

func TestAnticycleVersionWithHelpFlag_ShouldPrintHelp(t *testing.T) {
	stdOut, err := exec.Command("anticycle", "-version", "-help").Output()
	assert.NoError(t, err)
	validVersion := regexp.MustCompile(fullVersionRegex)
	assert.False(t, validVersion.MatchString(string(stdOut)))
}
