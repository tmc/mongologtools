package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagInput  = flag.String("i", "file://-", "input io path")
	flagOutput = flag.String("o", "file://-", "output io path")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument(s):", flag.Args())
		os.Exit(1)
	}
	input, err := GetIO(*flagInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error configurting input:", err)
		os.Exit(1)
	}
	r, err := input.Reader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}

	output, err := GetIO(*flagOutput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error configurting output:", err)
		os.Exit(1)
	}
	w, err := output.Writer()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening output:", err)
		os.Exit(1)
	}

	if err := ingest(r, w); err != nil {
		fmt.Fprintln(os.Stderr, "error ingesting:", err)
		os.Exit(1)
	}
}
