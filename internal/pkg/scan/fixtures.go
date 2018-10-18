package scan

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	tmp = "anticycle"
)

func tmpDir(rootDir string) (string, func()) {
	acTmp := fmt.Sprintf("%v/anticycle", os.TempDir())
	if _, err := os.Stat(acTmp); os.IsNotExist(err) {
		if err := os.Mkdir(acTmp, 0700); err != nil {
			os.RemoveAll(acTmp)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	dir := fmt.Sprintf("%v/%v", acTmp, rootDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0700); err != nil {
			os.RemoveAll(acTmp)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	remove := func() { os.RemoveAll(dir) }
	return dir, remove
}

func tmpFile(dir, filename, data string) (*os.File, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0700); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	if _, err := f.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

func makeProjectNoCycles(testName string) (string, func()) {
	dir, remove := tmpDir(testName)
	packages := []struct{Name, Data string}{
		{"foo", "package foo\nimport \"bar/bar\""},
		{"bar", "package bar\nimport \"fmt\""},
		{"baz", "package baz\nimport \"bar/bar\""},
	}

	if err:= _generateProject(dir, packages); err != nil {
		remove()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return dir, remove
}

func makeProjectOneToOne(testName string) (string, func()) {
	dir, remove := tmpDir(testName)
	packages := []struct{Name, Data string}{
		{"foo", "package foo"},
		{"bar", "package bar\nimport \"baz/baz\""},
		{"baz", "package baz\nimport \"bar/bar\""},
	}

	if err:= _generateProject(dir, packages); err != nil {
		remove()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return dir, remove
}

func makeProjectTriangle(testName string) (string, func()) {
	dir, remove := tmpDir(testName)
	packages := []struct{Name, Data string}{
		{"foo", "package foo\nimport \"baz/baz\""},
		{"bar", "package bar\nimport \"foo/foo\""},
		{"baz", "package baz\nimport \"bar/bar\""},
	}

	if err:= _generateProject(dir, packages); err != nil {
		remove()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return dir, remove
}

func makeProjectDiagonalSquare(testName string) (string, func()) {
	dir, remove := tmpDir(testName)
	packages := []struct{Name, Data string}{
		{"foo", "package foo\nimport \"pas/pas\""},
		{"bar", "package bar\nimport \"foo/foo\""},
		{"baz", "package baz\nimport \"bar/bar\""},
		{"pas", "package pas\nimport \"baz/baz\""},
	}

	if err:= _generateProject(dir, packages); err != nil {
		remove()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return dir, remove
}

func _generateProject(dir string, packages []struct{Name, Data string}) error {
	for _, pkg := range packages {
		_, err := tmpFile(filepath.Join(dir, pkg.Name), fmt.Sprintf("%v.go", pkg.Name), pkg.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
