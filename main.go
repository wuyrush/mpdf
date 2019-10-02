package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/h2non/filetype"
	"github.com/jung-kurt/gofpdf"
)

const (
	// https://github.com/h2non/filetype
	MIMETypePDF           = "application/pdf"
	FileHeaderLengthBytes = 261
)

func main() {
	// setup cli
	in := flag.String("in", ".", "Path to the directory containing PDF files to merge")
	out := flag.String("out", ".", "Path to the merged PDF file")
	chrono := flag.Bool("c", false, "Merge PDF files in the order of last modification time")
	reversed := flag.Bool("r", false, "Reverse merge order")
	overwrite := flag.Bool("f", false, "Overwrite even if output file already existed")
	flag.Parse()

	outAbs, err := filepath.Abs(*out)
	if err != nil {
		fmt.Printf("Error when checking output path %s: %s\n", *out, err)
		os.Exit(1)
	}
	ok, err := canWriteOut(outAbs, *overwrite)
	if err != nil {
		fmt.Printf("Error when checking output path %s: %s\n", outAbs, err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("Cannot write to output path %s", outAbs)
		os.Exit(1)
	}
	inAbs, err := filepath.Abs(*in)
	if err != nil {
		fmt.Printf("Error when checking input path %s: %s\n", *in, err)
		os.Exit(1)
	}
	m := merger{
		in:        inAbs,
		out:       outAbs,
		chrono:    *chrono,
		reversed:  *reversed,
		overwrite: *overwrite,
	}
	err = m.merge()
	if err != nil {
		fmt.Printf("Error when merging PDF files: %s", err)
		os.Exit(1)
	}
}

/*
	Check if we can write to out;
		1. If out exists and is a file, canWrite = overwrite;
		2. If out exists and is a directory, canWrite = true. Use a random filename for out in this case
		3. If out doesn't exist:
			3.1 If out doesn't exist and its parent exist and the parent is a directory,
				then canWrite = true. Use a random filename for out in this case;
			3.2 If out doesn't exist and its parent doesn't exist / its parent exists but a is a non-dir,
				canWrite = false
*/
func canWriteOut(p string, overwrite bool) (bool, error) {
	info, err := os.Stat(p)
	isNotExist := os.IsNotExist(err)
	if err != nil && !isNotExist {
		return false, err
	}
	if isNotExist {
		// case 3
		parent := filepath.Dir(p)
		info, err := os.Stat(parent)
		// either the parent dir doesn't exist or any other errors occur, out is not writable anyways
		if err != nil {
			return false, err
		}
		return info.IsDir(), nil
	}
	if info.IsDir() {
		return true, nil
	}
	// assume out points to an existing file
	return overwrite, nil
}

type merger struct {
	in        string
	out       string
	reversed  bool
	chrono    bool
	overwrite bool
}

/*
	Then collect files which we believe are of type PDF;
	Then sort the files according to specified ordering;
	Call pdfcpu to merge files and output;
	Clean up if necessary, and exit.
*/
func (m merger) merge() error {
	pdfs, err := collectPDF(m.in)
	if err != nil {
		return err
	} else if len(pdfs) < 2 {
		fmt.Printf("Found %d PDF files under path %s. Skip merging", len(pdfs), m.in)
		return nil
	}
	// in-place sort
	m.sort(pdfs)
	// merge
	return m.doMerge(pdfs)
}

func isPDF(fp string) (bool, error) {
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		return false, err
	}
	head := make([]byte, FileHeaderLengthBytes)
	_, err = f.Read(head)
	if err != nil && err != io.EOF {
		return false, err
	}
	if filetype.IsMIME(head, MIMETypePDF) {
		return true, nil
	}
	return false, nil
}

/*
	Processes files which are the direct children of path p only. Assumes p is an absolute path.
*/
func collectPDF(p string) ([]pdfFile, error) {
	var wg sync.WaitGroup
	checked := make(chan struct {
		f   pdfFile
		err error
	})
	exit := make(chan bool)
	filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if path == p {
			return nil
		} else if info.IsDir() {
			return filepath.SkipDir
		}
		// path points to a file
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := isPDF(path)
			if err != nil || ok {
				select {
				case <-exit:
				case checked <- struct {
					f   pdfFile
					err error
				}{f: pdfFile{path, info.ModTime()}, err: err}:
				}
			}
		}()
		return nil
	})
	go func() {
		wg.Wait()
		// all goroutine finished
		close(checked)
	}()
	defer close(exit)
	var pdfFiles []pdfFile
	for i := range checked {
		if i.err != nil {
			return nil, fmt.Errorf("Error checking type of file in %s: %w", i.f.path, i.err)
		}
		pdfFiles = append(pdfFiles, i.f)
	}
	return pdfFiles, nil
}

type pdfFile struct {
	path    string
	modTime time.Time
}

func (m merger) sort(pdfs []pdfFile) {

}

func (m merger) doMerge(pdfs []pdfFile) error {
	return nil
}

func genPDF(n int) {
	for i := 1; i <= n; i++ {
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, fmt.Sprint(i))
		err := pdf.OutputFileAndClose(fmt.Sprintf("./data/file-%d.pdf", i))
		if err != nil {
			panic(err)
		}
	}
}
