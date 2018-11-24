// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// anticycle: Static code analysis for dependency cycles.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/anticycle/anticycle/pkg/anticycle"
	"github.com/anticycle/anticycle/pkg/serialize"
)

var version = "undefined"
var build = "undefined"

const helpText = `Usage: anticycle [options] [directory]

  Anticycle is a tool for static code analysis which search for 
  dependency cycles. It scans recursively all source files and 
  parses theirs dependencies. Anticycle does not compile the code, 
  so it is ideal for searching for complex, difficult to debug cycles.

Options:
  -all               Output all packages, with and without cycles.

  -format="text"     Output format. Available: text, json.

  -exclude=""        A space-separated list of directories that should 
                     not be scanned. The list will be added to the 
                     default list of directories.
  -excludeOnly=""    A space-separated list of directories that should 
                     not be scanned. The list will override the default.
  -showExclude       Shows default list of excluded directories.

  -help              Shows this help text.
  -version           Shows version and build hash.

Directory:
  An optional path to the analyzed project. If the directory is not 
  defined, the current working directory will be used.

Output:
  The output of the Anticycle is a text by default with human friendly
  format of cycle affected package, import and filename.
  
  You can specify json format with -format=json flag.
  The JSON contains package names, 
  all dependencies in the package, a list of files belonging to 
  the package and their individual dependencies.

  By default output will contain only cycles to reduce a clutter.
  If you want to print all packages you can use -all flag.

  If Anticycle does not find any source files, it will exit 
  with code 0 without output in text format,
  or with empty structure in JSON format.
  In case of error, the program will exit with code 1, and the error 
  message will be sent to stderr.`

func trap(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	showHelp := flag.Bool("help", false, "Show help text")
	showVersion := flag.Bool("version", false, "Show version tag")
	showEx := flag.Bool("showExclude", false, "Show default list of excluded directories")
	ex := flag.String("exclude", "", "A space-separated list of directories")
	exO := flag.String("excludeOnly", "", "A space-separated list of directories")
	format := flag.String("format", "text", "Output format. Available: text,json")
	all := flag.Bool("all", false, "Output all packages, with and without cycles")
	flag.Parse()

	// SHOW HELP TEXT
	if *showHelp == true {
		fmt.Fprintln(os.Stdout, helpText)
		os.Exit(0)
	}

	// SHOW VERSION
	if *showVersion == true {
		var ver string
		if strings.ToLower(*format) == "json" {
			verJSON, _ := json.Marshal(map[string]string{"version": version, "build": build})
			ver = string(verJSON)
		} else {
			ver = fmt.Sprintf("Anticycle version %s, build %s", version, build)
		}

		fmt.Fprintln(os.Stdout, ver)
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
		paths := strings.Trim(*exO, "\"'")
		anticycle.DefaultExcluded = strings.Split(paths, " ")
	}

	// EXCLUDE
	if *ex != "" {
		paths := strings.Trim(*ex, "\"'")
		exclude = append(exclude, strings.Split(paths, " ")...)
	}

	// DIRECTORY
	var dir string
	if len(flag.Args()) > 0 {
		dir = path.Clean(flag.Args()[0])
	} else {
		dir = "."
	}

	var err error
	var output string

	excluded := anticycle.ExcludeDirs(exclude)
	cycles, err := anticycle.Collect(dir, excluded, *all)
	trap(err)

	analysis := anticycle.Analyze(cycles)

	switch strings.ToLower(*format) {
	case "json":
		output, err = serialize.ToJSON(analysis)
	case "text":
		output, err = serialize.ToTxt(analysis)
	default:
		err = fmt.Errorf("-format='%v' is not available, try 'text' or 'json'", *format)
	}
	trap(err)
	if output != "" {
		fmt.Fprintln(os.Stdout, output)
	}
	os.Exit(0)
}
