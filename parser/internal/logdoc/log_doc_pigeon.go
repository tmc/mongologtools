//go:generate pigeon -o log_doc_pigeon_generated.go log_doc_pigeon.peg
//go:generate goimports -w log_doc_pigeon_generated.go

package logdoc

import "bytes"

// ConvertLogToExtendedPigeon converts MongoDB log line formatted documents to an extended JSON representation
func ConvertLogToExtendedPigeon(input []byte) (map[string]interface{}, error) {
	debug := false
	result, err := ParseReader("", bytes.NewReader(input), Debug(debug))
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

func ConvertLogToExtendedPigeonMemoized(input []byte) (map[string]interface{}, error) {
	result, err := ParseReader("", bytes.NewReader(input), Memoize(true))
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}
