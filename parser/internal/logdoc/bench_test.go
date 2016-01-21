package logdoc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	skipReason  string
	sampleLines [][]byte
)

func init() {
	sample, err := ioutil.ReadFile("benchmark_sample.log")
	if err != nil {
		skipReason = fmt.Sprintf("skipping because 'benchmark_sample.log' is missing: %v", err)
	}
	sampleLines = bytes.Split(sample, []byte("\n"))
}

func BenchmarkParsingPigeon(b *testing.B) {
	bench(b, ConvertLogToExtendedPigeon)
}

func BenchmarkParsingPigeonMemoized(b *testing.B) {
	bench(b, ConvertLogToExtendedPigeonMemoized)
}

func BenchmarkParsingPointlander(b *testing.B) {
	bench(b, ConvertLogToExtended)
}

type benchFunc func([]byte) (map[string]interface{}, error)

func bench(b *testing.B, benchFn benchFunc) {
	if skipReason != "" {
		b.Skip(skipReason)
	}
	for i := 0; i < b.N; i++ {
		line := sampleLines[i%len(sampleLines)]
		if len(line) == 0 {
			continue
		}
		_, err := benchFn(line)
		if err != nil {
			b.Error(err)
		}
	}
}
