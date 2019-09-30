package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPathString(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Log("Error when getting current working dir")
		t.Fail()
	}
	absWd, err := filepath.Abs(wd)
	if err != nil {
		t.Log("Error when getting absolute path of current working dir")
		t.Fail()
	}

	cases := []struct {
		p        path
		expected string
	}{
		{path("./foo/bar.pdf"), filepath.Join(absWd, "foo", "bar.pdf")},
		{path("./../foo/bar.pdf"), filepath.Join(filepath.Dir(absWd), "foo", "bar.pdf")},
		{path("../foo/bar.pdf"), filepath.Join(filepath.Dir(absWd), "foo", "bar.pdf")},
	}

	for _, c := range cases {
		cs := c
		t.Run(string(cs.p), func(t *testing.T) {
			actual, err := cs.p.abs()
			if err != nil {
				t.Errorf("error getting absolute path of %s: %v", cs.p, err)
			}
			if actual != cs.expected {
				t.Errorf("got: %s want: %s", actual, cs.expected)
			}
		})
	}
}

func TestNewMerger(t *testing.T) {

	inPath, outPath := path("."), path(".")
	_ = merger{
		in:        inPath,
		out:       outPath,
		overwrite: true,
		chrono:    true,
		reversed:  true,
	}
}

/*
 prior condition:
 a. out is 0) non-existent 1) a file 2) a directory
 b. overwrite is 0) false 2) true
 if out is non-existent:
	return |out's parent directory exist?|
 if out is a file:
	return |can overwrite?|
 if out is a directory:
	return true
*/
func TestCanWriteOut(t *testing.T) {
	// get project root, or we can pass this as env var
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	// generate a temporary fs hierarchy
	outDir, err := setupOutDir(cwd)
	defer os.RemoveAll(outDir)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		name      string
		p         path
		overwrite bool
		expected  bool
	}{
		{
			name:      "[out=NonExistentFileInExistingDir,overwrite=false]",
			p:         path(filepath.Join(outDir, "bar")),
			overwrite: false,
			expected:  true,
		},
		{
			name:      "[out=NonExistentFileInExistingDir,overwrite=true]",
			p:         path(filepath.Join(outDir, "bar")),
			overwrite: true,
			expected:  true,
		},
		{
			name:      "[out=NonExistentFileInNonExistentDir,overwrite=false]",
			p:         path(filepath.Join(outDir, "bar", "qux")),
			overwrite: false,
			expected:  false,
		},
		{
			name:      "[out=NonExistentFileInNonExistentDir,overwrite=true]",
			p:         path(filepath.Join(outDir, "bar", "qux")),
			overwrite: true,
			expected:  false,
		},
		{
			name:      "[out=ExistingDir,overwrite=false]",
			p:         path(filepath.Join(outDir, "foo")),
			overwrite: false,
			expected:  true,
		},
		{
			name:      "[out=ExistingDir,overwrite=true]",
			p:         path(filepath.Join(outDir, "foo")),
			overwrite: true,
			expected:  true},
		{
			name:      "[out=ExistingFile,overwrite=false]",
			p:         path(filepath.Join(outDir, "foo", "bar")),
			overwrite: false,
			expected:  false,
		},
		{
			name:      "[out=ExistingFile,overwrite=true]",
			p:         path(filepath.Join(outDir, "foo", "bar")),
			overwrite: true,
			expected:  true,
		},
	}

	for _, c := range cases {
		cs := c
		t.Run(cs.name, func(t *testing.T) {
			m := merger{
				out:       cs.p,
				overwrite: cs.overwrite,
			}
			actual, err := m.canWriteOut()
			if err != nil {
				t.Errorf("error when calling canWriteOut() on path %s: %v", cs.p, err)
			}
			if actual != cs.expected {
				t.Errorf("got: %t want: %t", actual, cs.expected)
			}
		})
	}
}

func TestMerge() {

}

/*
	Create a directory for testing, with following hierarchy:
	project-root
		|-<random-name>
			|-foo
				|-bar (file)
*/
func setupOutDir(root string) (dir string, err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()
	dir, err = ioutil.TempDir(root, "test")
	if err != nil {
		log.Println("error creating output dir for testing")
		return
	}
	nestedDir := filepath.Join(dir, "foo")
	err = os.Mkdir(nestedDir, 0700)
	if err != nil {
		log.Printf("error creating nested dir %s\n", nestedDir)
		return
	}
	filename := filepath.Join(nestedDir, "bar")
	err = ioutil.WriteFile(filename, []byte(""), 0640)
	if err != nil {
		log.Printf("error creating temporary file %s: %v\n", filename, err)
		return
	}
	return dir, err
}
