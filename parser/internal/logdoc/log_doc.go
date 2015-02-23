package logdoc

import "fmt"

// ConvertLogToExtended converts MongoDB log line formatted documents to an extended JSON representation
func ConvertLogToExtended(input []byte) (map[string]interface{}, error) {
	p := &LogDocParser{Buffer: string(input)}
	p.Init()
	p.LogDoc.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	if len(p.Values) == 0 {
		return nil, fmt.Errorf("log_doc: no values present after parsing")
	}
	if doc, ok := p.Values[0].(map[string]interface{}); ok {
		return doc, nil
	}
	return nil, fmt.Errorf("log_doc: got unexpected type %T", p.Values[0])
}

type LogDoc struct {
	Maps   []map[string]interface{}
	Lists  [][]interface{}
	Fields []string
	Values []interface{}
}

func (d *LogDoc) Init() {
	d.Maps = make([]map[string]interface{}, 0)
	d.Lists = make([][]interface{}, 0)
}

func (d *LogDoc) pushMap() {
	d.Maps = append(d.Maps, make(map[string]interface{}))
	d.Values = append(d.Values, d.Maps[len(d.Maps)-1])
}

func (d *LogDoc) pushList() {
	d.Lists = append(d.Lists, make([]interface{}, 0))
	d.Values = append(d.Values, d.Lists[len(d.Lists)-1])
}

func (d *LogDoc) pushValue(value interface{}) {
	d.Values = append(d.Values, value)
}

func (d *LogDoc) popValue() interface{} {
	value := d.Values[len(d.Values)-1]
	d.Values = d.Values[:len(d.Values)-1]
	return value
}

func (d *LogDoc) pushField(field string) {
	d.Fields = append(d.Fields, field)
}

func (d *LogDoc) popField() string {
	field := d.Fields[len(d.Fields)-1]
	d.Fields = d.Fields[:len(d.Fields)-1]
	return field
}

func (d *LogDoc) setMapValue() {
	field, value := d.popField(), d.popValue()
	d.Maps[len(d.Maps)-1][field] = value
}

func (d *LogDoc) popMap() {
	d.Maps = d.Maps[:len(d.Maps)-1]
}

func (d *LogDoc) popList() {
	d.Lists = d.Lists[:len(d.Lists)-1]
}

func (d *LogDoc) setListValue() {
	d.Lists[len(d.Lists)-1] = append(d.Lists[len(d.Lists)-1], d.popValue())
}
