package parser_test

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/mongologtools/parser"
)

func ExampleParseLogLine() {
	line := "Mon Feb 23 03:20:19.670 [TTLMonitor] query local.system.indexes query: { expireAfterSeconds: { $exists: true } } ntoreturn:0 ntoskip:0 nscanned:0 keyUpdates:0 locks(micros) r:86 nreturned:0 reslen:20 0ms"
	doc, err := parser.ParseLogLine(line)
	if err != nil {
		panic(err)
	}
	buf, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))
	// output:
	// {"duration_ms":"0","keyUpdates":0,"nreturned":0,"ns":"local.system.indexes","nscanned":0,"ntoreturn":0,"ntoskip":0,"op":"query","query":{"expireAfterSeconds":{"$exists":true}},"reslen":20,"thread":"TTLMonitor","timestamp":"Mon Feb 23 03:20:19.670"}
}
