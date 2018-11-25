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
  -all                 Output all packages, with and without cycles.

  -format="text"       Output format. Available: text, json.

  -exclude=""          A space-separated list of directories that should 
                       not be scanned. The list will be added to the 
                       default list of directories.
  -excludeDefault=""   A space-separated list of directories that should 
                       not be scanned. The list will override the default.
  -showExcluded        Shows list of excluded directories.

  -help                Shows this help text.
  -version             Shows version and build hash.

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
		_, fatal := fmt.Fprintln(os.Stderr, err)
		if fatal != nil {
			panic(fatal)
		}
		os.Exit(1)
	}
}

func main() {
	showHelp := flag.Bool("help", false, "Show help text.")
	showVersion := flag.Bool("version", false, "Show version and build hash.")
	showExcluded := flag.Bool("showExcluded", false, "Show list of excluded directories.")

	setExclude := flag.String("exclude", "/", "A space-separated list of directories.")
	setExcludeDefault := flag.String("excludeDefault", "/", "A space-separated list of directories.")

	outputFormat := flag.String("format", "text", "Output format. Available: text,json.")
	outputAll := flag.Bool("all", false, "Output all packages, with and without cycles.")
	flag.Parse()

	var err error

	err = validateFormat(*outputFormat)
	trap(err)

	if *showHelp == true {
		err = printOutput(renderHelp())
		trap(err)
		os.Exit(0)
	}

	if *showVersion == true {
		err = printOutput(renderVersion(*outputFormat))
		trap(err)
		os.Exit(0)
	}

	excluded := excludedDirs(*setExclude, *setExcludeDefault)
	if *showExcluded == true {
		err = printOutput(renderExcluded(*outputFormat, excluded))
		trap(err)
		os.Exit(0)
	}

	dir := rootDir(flag.Args())
	output, err := findCycles(dir, excluded, *outputFormat, *outputAll)
	trap(err)

	err = printOutput(output)
	trap(err)

	os.Exit(0)
}

func validateFormat(format string) (err error) {
	if !strings.EqualFold(format, "json") && !strings.EqualFold(format, "text") {
		err = fmt.Errorf("-format='%v' is not available, try 'text' or 'json'", format)
	}
	return err
}

func renderHelp() string {
	// TODO(pawelzny) support for JSON output format.
	return helpText
}

func renderVersion(format string) string {
	if strings.ToLower(format) == "json" {
		verJSON, _ := json.Marshal(map[string]string{"version": version, "build": build})
		return string(verJSON)
	}
	return fmt.Sprintf("Anticycle version %s, build %s", version, build)
}

func renderExcluded(format string, excluded []string) string {
	if strings.ToLower(format) == "json" {
		exJSON, _ := json.Marshal(map[string][]string{"excluded": excluded})
		return string(exJSON)
	}
	if len(excluded) == 0 {
		return ""
	}
	return strings.Join(excluded, "\n")
}

func excludedDirs(setExclude, setExcludeDefault string) (exclude []string) {
	if setExcludeDefault == "" {
		anticycle.DefaultExcluded = []string{}
	} else if setExcludeDefault != "/" {
		paths := strings.Trim(setExcludeDefault, "\"'")
		anticycle.DefaultExcluded = strings.Split(paths, " ")
	}

	if setExclude != "/" {
		paths := strings.Trim(setExclude, "\"'")
		exclude = strings.Split(paths, " ")
	}

	return anticycle.ExcludeDirs(exclude)
}

func rootDir(args []string) string {
	if len(args) > 0 {
		return path.Clean(args[0])
	}
	return "."
}

func findCycles(dir string, excluded []string, format string, all bool) (output string, err error) {
	cycles, err := anticycle.Collect(dir, excluded, all)
	if err != nil {
		return output, err
	}
	analysis := anticycle.Analyze(cycles)

	switch strings.ToLower(format) {
	case "json":
		output, err = serialize.ToJSON(analysis)
	case "text":
		output, err = serialize.ToTxt(analysis)
	}

	return output, err
}

func printOutput(output string) (err error) {
	if output != "" {
		_, err = fmt.Fprintln(os.Stdout, output)
	}
	return err
}
