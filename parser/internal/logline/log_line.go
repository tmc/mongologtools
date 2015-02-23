//go:generate peg -switch -inline log_line.peg

package logline

import "github.com/tmc/mongologtools/parser/internal/logdoc"

func ParseLogLine(input string) (map[string]interface{}, error) {
	p := logLineParser{Buffer: input}
	p.Init()
	p.logLine.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()
	return p.Fields, nil

}

type logLine struct {
	logdoc.LogDoc

	Fields     map[string]interface{}
	fieldNames []string
}

func (m *logLine) Init() {
	m.LogDoc.Init()
	m.Fields = make(map[string]interface{})
	m.fieldNames = make([]string, 0)
}

func (m *logLine) SetField(key string, value string) {
	m.Fields[key] = value
}

func (m *logLine) StartField(fieldName string) {
	m.fieldNames = append(m.fieldNames, fieldName)
}

func (m *logLine) EndField() {
	lastValue := m.LogDoc.PopValue()
	targetField := m.fieldNames[len(m.fieldNames)-1]
	m.Fields[targetField] = lastValue
	m.fieldNames = m.fieldNames[:len(m.fieldNames)-1]
}
