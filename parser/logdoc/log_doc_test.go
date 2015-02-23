//go:generate peg -inline -switch log_doc.peg

package logdoc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tmc/mongologtools/parser/logdoc"
)

func TestLogDocParser(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`{ foo: 42 }`, `{"foo":42}`},
		{`{ _updated_at: { $lte: new Date(1412941647719) } }`, `{"_updated_at":{"$lte":{"$date":"2014-10-10T11:47:27.719Z"}}}`},
		{`{ _id: ObjectId("54e792daf1845f045f4c000e"), data: BinData(0,"aGVsbG8K") }`, `{"_id":{"$oid":"54e792daf1845f045f4c000e"},"data":{"$binary":"aGVsbG8K","$type":"00"}}`},
		{`{ t: Timestamp(1420000000, 1) }`, `{"t":{"$timestamp":{"t":1420000000,"i":1}}}`},
		{`{ some_text: /ese/i }`, `{"some_text":{"$regex":"ese","$options":"i"}}`},
		{`{ n: NumberLong(-9223372036854775808) }`, `{"n":{"$numberLong":"-9223372036854775808"}}`},
	}
	for i, testcase := range cases {
		doc, err := logdoc.ConvertLogToExtended([]byte(testcase.input))
		if err != nil {
			t.Fatalf("case %d: error parsing: %v", i, err)
		}
		buf, err := json.Marshal(doc)
		if err != nil {
			t.Fatalf("case %d: error marshaling: %v", i, err)
		}
		result := string(buf)
		if result != testcase.expected {
			t.Errorf("case %d: expected '%s'\nbut got '%s'", i, testcase.expected, result)
		}
	}
}

func ExampleConvertLogToExtended() {
	doc, _ := logdoc.ConvertLogToExtended([]byte("{ x: Timestamp(13000000, 0)}"))
	buf, _ := json.Marshal(doc)
	fmt.Print(string(buf))
	// output:
	// {"x":{"$timestamp":{"t":13000000,"i":0}}}
}
