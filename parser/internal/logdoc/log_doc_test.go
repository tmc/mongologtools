package logdoc_test

import (
	"encoding/json"
	"testing"

	"github.com/tmc/mongologtools/parser/internal/logdoc"
)

func TestLogDocParser(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`{ foo: [ 42 ] }`, `{"foo":[42]}`},
		{`{ foo: [ 42, 43, 44 ] }`, `{"foo":[42,43,44]}`},
		{`{ _updated_at: { $lte: new Date(1412941647719) } }`, `{"_updated_at":{"$lte":{"$date":"2014-10-10T11:47:27.719Z"}}}`},
		{`{ _id: ObjectId("54e792daf1845f045f4c000e"), data: BinData(0,"aGVsbG8K") }`, `{"_id":{"$oid":"54e792daf1845f045f4c000e"},"data":{"$binary":"aGVsbG8K","$type":"00"}}`},
		{`{ t: Timestamp(1420000000, 1) }`, `{"t":{"$timestamp":{"t":1420000000,"i":1}}}`},
		{`{ some_text: /ese/i }`, `{"some_text":{"$regex":"ese","$options":"i"}}`},
		{`{ n: NumberLong(-9223372036854775808) }`, `{"n":{"$numberLong":"-9223372036854775808"}}`},
		{`{ $query: { _p_user: "_User$XOWpxhM06t", _created_at: { $gte: new Date(1430758114453) }, _p_swiped_user: { $in: [ "_User$H5uW6XZTM4", "_User$MSjlL4vwIs", "_User$J04ObUTPzc", "_User$j2ToPWWr9b", "_User$AM6zi0ZSbk", "_User$B7NycSBHGN", "_User$qntkZ1RGMS", "_User$yZi1WUgh53", "_User$C45ZweDNei", "_User$sYVqPmtWtv", "_User$moJoWODQzq", "_User$YUybni9dpa", "_User$t01gQNtkf4", "_User$w3afYlWFSB", "_User$uQc0wUkiBl", "_User$3C0D6G1Qvh", "_User$YtythybLMl", "_User$nNojFdVEt9", "_User$1w7YF9Sou2", "_User$xmYaMkbeZK", "_User$WWXc9MZSSG", "_User$T5logFfrTT", "_User$VIYqMyy7UT", "_User$YeQZ6P1YeY", "_User$rh49GFwRfc", "_User$F0nN6AbaFw" ] } }, $maxScan: 500000 }`, `{"$maxScan":500000,"$query":{"_created_at":{"$gte":{"$date":"2015-05-04T16:48:34.453Z"}},"_p_swiped_user":{"$in":["_User$H5uW6XZTM4","_User$MSjlL4vwIs","_User$J04ObUTPzc","_User$j2ToPWWr9b","_User$AM6zi0ZSbk","_User$B7NycSBHGN","_User$qntkZ1RGMS","_User$yZi1WUgh53","_User$C45ZweDNei","_User$sYVqPmtWtv","_User$moJoWODQzq","_User$YUybni9dpa","_User$t01gQNtkf4","_User$w3afYlWFSB","_User$uQc0wUkiBl","_User$3C0D6G1Qvh","_User$YtythybLMl","_User$nNojFdVEt9","_User$1w7YF9Sou2","_User$xmYaMkbeZK","_User$WWXc9MZSSG","_User$T5logFfrTT","_User$VIYqMyy7UT","_User$YeQZ6P1YeY","_User$rh49GFwRfc","_User$F0nN6AbaFw"]},"_p_user":"_User$XOWpxhM06t"}}`},
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
	for i, testcase := range cases {
		doc, err := logdoc.ConvertLogToExtendedPigeon([]byte(testcase.input))
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
