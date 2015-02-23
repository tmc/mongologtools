package parser

import "github.com/tmc/mongologtools/parser/internal/logline"

// ParseLogLine attempts to parse a mongodb log line into a structured representation
func ParseLogLine(input string) (map[string]interface{}, error) {
	return logline.ParseLogLine(input)
}
