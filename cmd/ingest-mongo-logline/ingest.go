package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/tmc/mongologtools/parser"
)

func ingest(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	out := json.NewEncoder(w)
	for s.Scan() {
		r, err := parser.ParseLogLine(s.Text())
		if err != nil {
			log.Printf("line parsing err on `%s..`\n", string(s.Bytes()[:min(len(s.Text()), 30)]))
		}
		out.Encode(r)
	}
	return nil
}

func min(n, m int) int {
	if n < m {
		return n
	}
	return m
}
