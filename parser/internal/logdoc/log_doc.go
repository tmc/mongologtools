//go:generate peg -inline -switch log_doc.peg

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
	Maps   []int
	Lists  []int
	Fields []string
	Values []interface{}
}

func (d *LogDoc) Init() {
	d.Maps = make([]int, 0)
	d.Lists = make([]int, 0)
}

func (d *LogDoc) PushMap() {
	d.Values = append(d.Values, make(map[string]interface{}))
	d.Maps = append(d.Maps, len(d.Values)-1)
}

func (d *LogDoc) PushList() {
	d.Values = append(d.Values, make([]interface{}, 0))
	d.Lists = append(d.Lists, len(d.Values)-1)
}

func (d *LogDoc) PushValue(value interface{}) {
	d.Values = append(d.Values, value)
}

func (d *LogDoc) PopValue() interface{} {
	value := d.Values[len(d.Values)-1]
	d.Values = d.Values[:len(d.Values)-1]
	return value
}

func (d *LogDoc) PushField(field string) {
	d.Fields = append(d.Fields, field)
}

func (d *LogDoc) PopField() string {
	field := d.Fields[len(d.Fields)-1]
	d.Fields = d.Fields[:len(d.Fields)-1]
	return field
}

func (d *LogDoc) SetMapValue() {
	field, value := d.PopField(), d.PopValue()
	i := d.Maps[len(d.Maps)-1]
	d.Values[i].(map[string]interface{})[field] = value
}

func (d *LogDoc) PopMap() {
	d.Maps = d.Maps[:len(d.Maps)-1]
}

func (d *LogDoc) PopList() {
	d.Lists = d.Lists[:len(d.Lists)-1]
}

func (d *LogDoc) SetListValue() {
	i := d.Lists[len(d.Lists)-1]
	d.Values[i] = append(d.Values[i].([]interface{}), d.PopValue())
}
