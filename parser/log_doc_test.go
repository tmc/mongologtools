package parser

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/mongologtools/parser/internal/logdoc"
)

func ExampleConvertLogToExtended() {
	doc, _ := logdoc.ConvertLogToExtended([]byte("{ x: Timestamp(13000000, 0)}"))
	buf, _ := json.Marshal(doc)
	fmt.Print(string(buf))
	// output:
	// {"x":{"$timestamp":{"t":13000000,"i":0}}}
}
