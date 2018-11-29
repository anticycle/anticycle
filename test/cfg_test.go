// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testCase struct {
		isJSON       bool
		name, golden string
		args         []string
	}
	testScenario struct {
		name, testdata string
	}
)

// Update flag for sanity and acceptance tests
var update = flag.Bool("update", false, "update .golden files")

func updateGolden(filename string, data []byte) {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func readGolden(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func runTestGolden(t *testing.T, scenario testScenario, test testCase) {
	testName := fmt.Sprintf(test.name, scenario.name)

	t.Run(testName, func(t *testing.T) {
		goldenFile := filepath.Join("testdata", scenario.testdata, test.golden)
		test.args = append(test.args, filepath.Join("testdata", scenario.testdata))

		stdOut, err := exec.Command("anticycle", test.args...).Output()
		assert.NoError(t, err)
		if *update {
			updateGolden(goldenFile, stdOut)
		}

		if test.isJSON {
			var expected, result map[string]interface{}
			golden := readGolden(goldenFile)
			if err := json.Unmarshal(stdOut, &result); err != nil {
				panic(err)
			}
			if err := json.Unmarshal(golden, &expected); err != nil {
				panic(err)
			}

			assert.Equal(t, expected, result)
		} else {
			expected := readGolden(goldenFile)
			assert.Equal(t, string(expected), string(stdOut))
		}
	})
}
