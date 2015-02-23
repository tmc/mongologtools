package parser

import "github.com/tmc/mongologtools/parser/internal/logdoc"

// ConvertLogToExtended converts MongoDB log line formatted documents to an extended JSON representation
func ConvertLogToExtended(input []byte) (map[string]interface{}, error) {
	return logdoc.ConvertLogToExtended(input)
}
