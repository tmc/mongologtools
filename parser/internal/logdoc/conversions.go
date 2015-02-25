package logdoc

import (
	"encoding/json"
	"strconv"
	"strings"

	mongo_json "github.com/mongodb/mongo-tools/common/json"
)

func (p *LogDoc) Numeric(value string) interface{} {
	n := json.Number(value)
	if i64, err := n.Int64(); err == nil {
		return i64
	}
	f64, _ := n.Float64()
	return f64
}

func (d *LogDoc) Date(value string) mongo_json.Date {
	n, _ := strconv.Atoi(value)
	return mongo_json.Date(n)
}

func (d *LogDoc) ObjectId(value string) mongo_json.ObjectId {
	return mongo_json.ObjectId(value)
}

func (d *LogDoc) Bindata(value string) mongo_json.BinData {
	// example: BinData(0,"aGVsbG8K")
	parts := strings.Split(value, ",")
	binType, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	data := strings.Trim(parts[1], `"`)
	return mongo_json.BinData{
		Type:   byte(binType),
		Base64: data,
	}

}

func (d *LogDoc) Timestamp(value string) mongo_json.Timestamp {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		parts = strings.Split(value, "|")
	}
	if len(parts) != 2 {
		return mongo_json.Timestamp{}
	}
	p1, p2 := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	seconds, _ := strconv.ParseUint(p1, 10, 0)
	increment, _ := strconv.ParseUint(p2, 10, 0)
	return mongo_json.Timestamp{
		Seconds:   uint32(seconds),
		Increment: uint32(increment),
	}
}

func (d *LogDoc) Numberlong(value string) mongo_json.NumberLong {
	n, _ := strconv.ParseInt(value, 10, 0)
	return mongo_json.NumberLong(n)
}
func (d *LogDoc) Regex(value string) mongo_json.RegExp {
	slashIdx := strings.LastIndex(value, "/")
	pattern, options := value[:slashIdx], value[slashIdx+1:]
	return mongo_json.RegExp{
		Pattern: pattern,
		Options: options,
	}
}

func (d *LogDoc) Minkey() mongo_json.MinKey {
	return mongo_json.MinKey{}
}

func (d *LogDoc) Maxkey() mongo_json.MaxKey {
	return mongo_json.MaxKey{}
}

func (d *LogDoc) Undefined() mongo_json.Undefined {
	return mongo_json.Undefined{}
}
