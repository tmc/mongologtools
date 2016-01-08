package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Be verbose"`
	Input   string `short:"i" long:"input" description:"Input" default:"file://-"`
	Output  string `short:"o" long:"output" description:"Output" default:"file://-"`
}

func main() {
	var opts Options
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument(s):", args)
		os.Exit(1)
	}
	input, err := Get(opts.Input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error configurting input:", err)
		os.Exit(1)
	}
	r, err := input.Reader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}

	output, err := Get(opts.Output)
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
