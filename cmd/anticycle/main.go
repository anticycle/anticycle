// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// anticycle: Static code analysis for dependency cycles.

package main

import (
	"flag"
	"fmt"
	"github.com/anticycle/anticycle/pkg/anticycle"
	"github.com/anticycle/anticycle/pkg/serialize"
	"os"
	"path"
	"strings"
)

const helpText = `Usage: anticycle [options] [directory]

  Anticycle is a tool for static code analysis which search for 
  dependency cycles. It scans recursively all source files and 
  parses theirs dependencies. Anticycle does not compile the code, 
  so it is ideal for searching for complex, difficult to debug cycles.

Options:
  -exclude=""        A comma separated list of directories that should 
                     not be scanned. The list will be added to the 
                     default list of directories.
  -excludeOnly=""    A comma separated list of directories that should 
                     not be scanned. The list will override the default.
  -showExclude       Shows default list of excluded directories.
  -help              Shows this help text.

Directory:
  An optional path to the analyzed project. If the directory is not 
  defined, the current working directory will be used.

Output:
  The output of the Anticycle is JSON string containing package names, 
  all dependencies in the package, a list of files belonging to 
  the package and their individual dependencies. If Anticycle does not 
  find any source files, it will exit with code 0 without output.
  In case of error, the program will exit with code 1, and the error 
  message will be sent to stderr.`

func trap(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	help := flag.Bool("help", false, "Show help text")
	showEx := flag.Bool("showExclude", false, "Show default list of excluded directories")
	ex := flag.String("exclude", "", "A comma separated list of directories")
	exO := flag.String("excludeOnly", "", "A comma separated list of directories")
	flag.Parse()

	// SHOW HELP TEXT
	if *help == true {
		fmt.Fprintln(os.Stdout, helpText)
		os.Exit(0)
	}

	// SHOW DEFAULT EXCLUDED DIRECTORIES
	if *showEx == true {
		fmt.Fprintln(os.Stdout, strings.Join(anticycle.DefaultExcluded, "\n"))
		os.Exit(0)
	}

	// EXCLUDE ONLY
	exclude := make([]string, 0)
	if *exO != "" {
		anticycle.DefaultExcluded = strings.Split(*exO, ",")
	}
	// EXCLUDE
	if *ex != "" {
		exclude = append(exclude, strings.Split(*ex, ",")...)
	}

	// DIRECTORY
	var dir string
	if len(flag.Args()) > 0 {
		dir = path.Clean(flag.Args()[0])
	} else {
		dir = "."
	}

	excluded := anticycle.ExcludeDirs(exclude)
	cycles, err := anticycle.Analyze(dir, excluded)
	trap(err)
	output, err := serialize.ToJson(cycles)
	trap(err)
	if output != "" {
		fmt.Fprintln(os.Stdout, output)
	}
	os.Exit(0)
}