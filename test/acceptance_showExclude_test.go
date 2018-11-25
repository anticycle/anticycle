package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnticycleShowExcluded(t *testing.T) {
	tests := []struct {
		name, expected string
		args           []string
	}{
		{
			name:     "Show default excluded",
			args:     []string{"-showExcluded"},
			expected: ".git\n.idea\n.vscode\nbin\ndist\ntestdata\nvendor\n",
		},
		{
			name:     "Set -exclude and -showExcluded",
			args:     []string{"-exclude='aaa bbb ccc'", "-showExcluded"},
			expected: ".git\n.idea\n.vscode\naaa\nbbb\nbin\nccc\ndist\ntestdata\nvendor\n",
		},
		{
			name:     "Set empty -excludeDefault and -showExcluded",
			args:     []string{"-excludeDefault=", "-showExcluded"},
			expected: "",
		},
		{
			name:     "Set -excludeDefault and -showExcluded",
			args:     []string{"-excludeDefault='foo bar'", "-showExcluded"},
			expected: "bar\nfoo\n",
		},
		{
			name:     "Set -exclude -excludeDefault and -showExcluded",
			args:     []string{"-exclude='aaa bbb ccc'", "-excludeDefault='foo bar'", "-showExcluded"},
			expected: "aaa\nbar\nbbb\nccc\nfoo\n",
		},
		{
			name:     "Set -exclude and -excludeDefault both empty",
			args:     []string{"-exclude=", "-excludeDefault=", "-showExcluded"},
			expected: "",
		},
		{
			name:     "Set -exclude and -excludeDefault both the same",
			args:     []string{"-exclude='foo bar'", "-excludeDefault='foo bar'", "-showExcluded"},
			expected: "bar\nbar\nfoo\nfoo\n", // TODO(pawelzny) remove duplicates
		},
		{
			name:     "Set -exclude and -showExcluded with -format=json",
			args:     []string{"-exclude='foo bar'", "-format=json", "-showExcluded"},
			expected: "{\"excluded\":[\".git\",\".idea\",\".vscode\",\"bar\",\"bin\",\"dist\",\"foo\",\"testdata\",\"vendor\"]}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, err := exec.Command("anticycle", test.args...).Output()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(stdOut))
		})
	}
}
