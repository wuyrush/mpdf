package main

import (
	"fmt"
	"os"
	"path/filepath"

	ftype "github.com/h2non/filetype"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	// prepare args

}

type filetype int

const (
	TypeAbsent filetype = iota
	TypeFile
	TypeDir
	TypeUnknown

	MIMETypePDF       = "application/pdf"
	HeaderLengthBytes = 261
)

// os file system path
type path string

// absolute path
func (p path) abs() (string, error) {
	abs, err := filepath.Abs(string(p))
	if err != nil {
		return "", err
	}
	return abs, nil
}

// path of parent directory
func (p path) parent() (path, error) {
	abs, err := p.abs()
	if err != nil {
		return path(""), err
	}
	return path(filepath.Dir(abs)), nil
}

func (p path) fileType() (filetype, error) {
	abs, err := p.abs()
	if err != nil {
		return TypeUnknown, err
	}
	info, err := os.Stat(abs)
	if err != nil && os.IsNotExist(err) {
		return TypeAbsent, nil
	} else if err != nil {
		return TypeUnknown, err
	}
	ft := TypeFile
	if info.IsDir() {
		ft = TypeDir
	}
	return ft, nil
}

type merger struct {
	in        path
	out       path
	chrono    bool
	reversed  bool
	overwrite bool
}

// Ok to write to the destination?
func (m merger) canWriteOut() (bool, error) {
	ft, err := m.out.fileType()
	if err != nil {
		return false, err
	}
	if ft == TypeAbsent {
		parent, err := m.out.parent()
		if err != nil {
			return false, err
		}
		ft, err := parent.fileType()
		if err != nil {
			return false, err
		}
		return ft == TypeDir, nil
	}
	if ft == TypeFile && !m.overwrite {
		return false, nil
	}
	return true, nil
}

func (m merger) merge() error {
	// read paths of files from input directory
	// sort the files according to specified ordering
	// call pdfcpu's merge API to merge the file
	// output the merged file to specified destination
}

// Returns files which we infer as type PDF under the given path. Only files which are
// the direct children of dir are inspected.
func listPDFs(dir path) ([]path, error) {
	ft, err := dir.fileType()
	if err != nil {
		return nil, err
	}
	if ft != TypeDir {
		return nil, fmt.Errorf("listPDFs: path %s is not a directory", dir)
	}
	r, err := dir.abs()
	if err != nil {
		return nil, fmt.Errorf("listPDFs: error getting absolute path of %s: %w", dir, err)
	}
	var selected []path
	var errWalk error
	filepath.Walk(r, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			errWalk = fmt.Errorf("error when walking directory %s: %w", r, err)
			return errWalk
		}
		if p == r {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		// This seems to be a file, infer its type.
		// TODO: parallelize the inferring process
		f, err := os.Open(p)
		defer f.Close()
		if err != nil {
			errWalk = fmt.Errorf("error opening file at path %s: %w", p, err)
			return errWalk
		}
		head := make([]byte, HeaderLengthBytes)
		_, err = f.Read(head)
		if err != nil {
			errWalk = fmt.Errorf("error reading file at path %s: %w", p, err)
			return errWalk
		}
		if ftype.IsMIME(head, MIMETypePDF) {
			selected = append(selected, path(p))
		}
		return nil
	})
	return selected, errWalk
}

func isPDF(p string) (bool, error) {
	f, err := os.Open(p)
	defer f.Close()
	if err != nil {
		return false, fmt.Errorf("isPDF: error opening file at path %s: %w", p, err)
	}
	head := make([]byte, HeaderLengthBytes)
	_, err = f.Read(head)
	if err != nil {
		return false, fmt.Errorf("isPDF: error reading file at path %s: %w", p, err)
	}
	if ftype.IsMIME(head, MIMETypePDF) {
		return true, nil
	}
	return false, nil

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
