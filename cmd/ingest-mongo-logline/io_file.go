package main

import (
	"io"
	"os"
)

type fileio struct {
	path string
}

func (f *fileio) Reader() (io.Reader, error) {
	if f.path == "-" {
		return os.Stdin, nil
	}
	return os.Open(f.path)
}

func (f *fileio) Writer() (io.Writer, error) {
	if f.path == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY, 0660)
}

func init() {
	RegisterIO("file", func(path string) IO {
		return &fileio{path: path}
	})
}
