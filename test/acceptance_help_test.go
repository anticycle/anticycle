package test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	helpPrefix = `Usage: anticycle [options] [directory]`
)

func TestAnticycleHelp_ShouldPrintHelpText(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Show only help",
			args: []string{"-help"},
		},
		{ // TODO(pawelzny) Add support for help text in json format
			name: "Show help with -format=json",
			args: []string{"-help", "-format=json"},
		},
		{
			name: "Show help and version",
			args: []string{"-help", "-version"},
		},
		{
			name: "Show help and excludeOnly",
			args: []string{"-help", "-excludeDefault=''"},
		},
		{
			name: "Show help and exclude and excludeOnly",
			args: []string{"-help", "-exclude=''", "-excludeDefault=''"},
		},
		{
			name: "Show help and showExclude",
			args: []string{"-help", "-showExcluded"},
		},
		{
			name: "Show help and all",
			args: []string{"-help", "-all"},
		},
		{
			name: "Show help with directory",
			args: []string{"-help", "./testdata"},
		},
		{
			name: "Show help with everything",
			args: []string{"-help", "-exclude=''", "-excludeDefault=''", "-all", "./testdata"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, err := exec.Command("anticycle", test.args...).Output()
			assert.NoError(t, err)
			assert.True(t, bytes.HasPrefix(stdOut, []byte(helpPrefix)))
		})
	}
}
