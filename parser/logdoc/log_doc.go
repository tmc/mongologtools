package logdoc

import "fmt"

type logDoc struct {
	maps   []map[string]interface{}
	lists  [][]interface{}
	fields []string
	values []interface{}
}

// ConvertLogToExtended converts MongoDB log line formatted documents to an extended JSON representation
func ConvertLogToExtended(input []byte) (map[string]interface{}, error) {
	p := &logDocParser{Buffer: string(input)}
	p.Init()
	p.logDoc.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	if len(p.values) == 0 {
		return nil, fmt.Errorf("log_doc: no values present after parsing")
	}
	if doc, ok := p.values[0].(map[string]interface{}); ok {
		return doc, nil
	}
	return nil, fmt.Errorf("log_doc: got unexpected type %T", p.values[0])
}

func (d *logDoc) Init() {
	d.maps = make([]map[string]interface{}, 0)
	d.lists = make([][]interface{}, 0)
}

func (d *logDoc) setRemainder(remainder string) {
	d.pushField("xx extra")
	d.pushValue(remainder)
	d.setMapValue()
}

func (d *logDoc) pushMap() {
	d.maps = append(d.maps, make(map[string]interface{}))
	d.values = append(d.values, d.maps[len(d.maps)-1])
}

func (d *logDoc) pushList() {
	d.lists = append(d.lists, make([]interface{}, 0))
	d.values = append(d.values, d.lists[len(d.lists)-1])
}

func (d *logDoc) pushValue(value interface{}) {
	d.values = append(d.values, value)
}

func (d *logDoc) popValue() interface{} {
	value := d.values[len(d.values)-1]
	d.values = d.values[:len(d.values)-1]
	return value
}

func (d *logDoc) pushField(field string) {
	d.fields = append(d.fields, field)
}

func (d *logDoc) popField() string {
	field := d.fields[len(d.fields)-1]
	d.fields = d.fields[:len(d.fields)-1]
	return field
}

func (d *logDoc) setMapValue() {
	field, value := d.popField(), d.popValue()
	d.maps[len(d.maps)-1][field] = value
}

func (d *logDoc) popMap() {
	d.maps = d.maps[:len(d.maps)-1]
}

func (d *logDoc) popList() {
	d.lists = d.lists[:len(d.lists)-1]
}

func (d *logDoc) setListValue() {
	d.lists[len(d.lists)-1] = append(d.lists[len(d.lists)-1], d.popValue())
}
