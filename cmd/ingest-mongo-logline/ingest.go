package main

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/tmc/mongologtools/parser"
)

func ingest(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	out := json.NewEncoder(w)
	for s.Scan() {
		r, err := parser.ParseLogLine(s.Text())
		if err != nil {
			return err
		}
		out.Encode(r)
	}
	return nil
}
