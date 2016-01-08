package parser_test

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/mongologtools/parser"
)

func ExampleParseLogLine() {
	line := "Mon Feb 23 03:20:19.670 [TTLMonitor] query local.system.indexes query: { expireAfterSeconds: { $exists: true } } ntoreturn:0 ntoskip:0 nscanned:0 keyUpdates:0 locks(micros) r:86 nreturned:0 reslen:20 0ms"
	doc, _ := parser.ParseLogLine(line)
	buf, _ := json.Marshal(doc)
	fmt.Print(string(buf))
	// output:
	// {"context":"TTLMonitor","duration_ms":"0","keyUpdates":0,"nreturned":0,"ns":"local.system.indexes","nscanned":0,"ntoreturn":0,"ntoskip":0,"op":"query","query":{"expireAfterSeconds":{"$exists":true}},"r":86,"reslen":20,"timestamp":"Mon Feb 23 03:20:19.670"}
}
