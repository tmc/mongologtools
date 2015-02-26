package logline

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleMongoLogLine
	ruleTimestamp
	ruleSeverity
	ruleComponent
	ruleContext
	ruleOp
	ruleWarning
	ruleloglineSizeWarning
	ruleLineField
	ruleNS
	ruleLocks
	rulelock
	ruleDuration
	ruleplainField
	rulecommandField
	ruleplanSummaryField
	ruleplanSummaryElements
	ruleplanSummaryElem
	ruleplanSummaryStage
	ruleplanSummary
	ruleexceptionField
	ruleLineValue
	ruletimestamp24
	ruletimestamp26
	ruledatetime26
	ruledigit4
	ruledigit2
	ruledate
	ruletz
	ruletime
	ruleday
	rulemonth
	ruledayNum
	rulehour
	ruleminute
	rulesecond
	rulemillisecond
	ruleletterOrDigit
	rulensChar
	ruleextra
	ruleS
	ruleDoc
	ruleDocElements
	ruleDocElem
	ruleList
	ruleListElements
	ruleListElem
	ruleField
	ruleValue
	ruleNumeric
	ruleBoolean
	ruleString
	ruleNull
	ruleTrue
	ruleFalse
	ruleDate
	ruleObjectID
	ruleBinData
	ruleRegex
	ruleTimestampVal
	ruletimestampParen
	ruletimestampPipe
	ruleNumberLong
	ruleMinKey
	ruleMaxKey
	ruleUndefined
	rulehexChar
	ruleregexChar
	ruleregexBody
	rulestringChar
	rulefieldChar
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39
	ruleAction40
	ruleAction41
	ruleAction42
	ruleAction43
	ruleAction44

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"MongoLogLine",
	"Timestamp",
	"Severity",
	"Component",
	"Context",
	"Op",
	"Warning",
	"loglineSizeWarning",
	"LineField",
	"NS",
	"Locks",
	"lock",
	"Duration",
	"plainField",
	"commandField",
	"planSummaryField",
	"planSummaryElements",
	"planSummaryElem",
	"planSummaryStage",
	"planSummary",
	"exceptionField",
	"LineValue",
	"timestamp24",
	"timestamp26",
	"datetime26",
	"digit4",
	"digit2",
	"date",
	"tz",
	"time",
	"day",
	"month",
	"dayNum",
	"hour",
	"minute",
	"second",
	"millisecond",
	"letterOrDigit",
	"nsChar",
	"extra",
	"S",
	"Doc",
	"DocElements",
	"DocElem",
	"List",
	"ListElements",
	"ListElem",
	"Field",
	"Value",
	"Numeric",
	"Boolean",
	"String",
	"Null",
	"True",
	"False",
	"Date",
	"ObjectID",
	"BinData",
	"Regex",
	"TimestampVal",
	"timestampParen",
	"timestampPipe",
	"NumberLong",
	"MinKey",
	"MaxKey",
	"Undefined",
	"hexChar",
	"regexChar",
	"regexBody",
	"stringChar",
	"fieldChar",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(buffer[node.begin:node.end]))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens16) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type logLineParser struct {
	logLine

	Buffer string
	buffer []rune
	rules  [118]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *logLineParser
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *logLineParser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *logLineParser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *logLineParser) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {
		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruleAction0:
			p.SetField("severity", buffer[begin:end])
		case ruleAction1:
			p.SetField("component", buffer[begin:end])
		case ruleAction2:
			p.SetField("context", buffer[begin:end])
		case ruleAction3:
			p.SetField("op", buffer[begin:end])
		case ruleAction4:
			p.SetField("warning", buffer[begin:end])
		case ruleAction5:
			p.SetField("ns", buffer[begin:end])
		case ruleAction6:
			p.StartField(buffer[begin:end])
		case ruleAction7:
			p.EndField()
		case ruleAction8:
			p.SetField("duration_ms", buffer[begin:end])
		case ruleAction9:
			p.StartField(buffer[begin:end])
		case ruleAction10:
			p.EndField()
		case ruleAction11:
			p.SetField("command_type", buffer[begin:end])
			p.StartField("command")
		case ruleAction12:
			p.EndField()
		case ruleAction13:
			p.StartField("planSummary")
			p.PushList()
		case ruleAction14:
			p.EndField()
		case ruleAction15:
			p.PushMap()
			p.PushField(buffer[begin:end])
		case ruleAction16:
			p.SetMapValue()
			p.SetListValue()
		case ruleAction17:
			p.PushValue(1)
			p.SetMapValue()
			p.SetListValue()
		case ruleAction18:
			p.StartField("exception")
		case ruleAction19:
			p.PushValue(buffer[begin:end])
			p.EndField()
		case ruleAction20:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction21:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction22:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction23:
			p.PushMap()
		case ruleAction24:
			p.PopMap()
		case ruleAction25:
			p.SetMapValue()
		case ruleAction26:
			p.PushList()
		case ruleAction27:
			p.PopList()
		case ruleAction28:
			p.SetListValue()
		case ruleAction29:
			p.PushField(buffer[begin:end])
		case ruleAction30:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction31:
			p.PushValue(buffer[begin:end])
		case ruleAction32:
			p.PushValue(nil)
		case ruleAction33:
			p.PushValue(true)
		case ruleAction34:
			p.PushValue(false)
		case ruleAction35:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction36:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction37:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction38:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction39:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction40:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction41:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction42:
			p.PushValue(p.Minkey())
		case ruleAction43:
			p.PushValue(p.Maxkey())
		case ruleAction44:
			p.PushValue(p.Undefined())

		}
	}
}

func (p *logLineParser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := 0, 0, 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin int) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 MongoLogLine <- <(LineField? Timestamp Severity? Component? Context Warning? Op NS LineField* Locks? LineField* Duration? extra? !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l2
					}
					goto l3
				l2:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
				}
			l3:
				{
					position4 := position
					depth++
					{
						position5, tokenIndex5, depth5 := position, tokenIndex, depth
						{
							position7 := position
							depth++
							{
								position8 := position
								depth++
								{
									position9 := position
									depth++
									{
										position10 := position
										depth++
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l6
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l6
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l6
										}
										position++
										depth--
										add(ruleday, position10)
									}
									if buffer[position] != rune(' ') {
										goto l6
									}
									position++
									{
										position11 := position
										depth++
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l6
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l6
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l6
										}
										position++
										depth--
										add(rulemonth, position11)
									}
									if buffer[position] != rune(' ') {
										goto l6
									}
									position++
								l12:
									{
										position13, tokenIndex13, depth13 := position, tokenIndex, depth
										if buffer[position] != rune(' ') {
											goto l13
										}
										position++
										goto l12
									l13:
										position, tokenIndex, depth = position13, tokenIndex13, depth13
									}
									{
										position14 := position
										depth++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l6
										}
										position++
										{
											position15, tokenIndex15, depth15 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l15
											}
											position++
											goto l16
										l15:
											position, tokenIndex, depth = position15, tokenIndex15, depth15
										}
									l16:
										depth--
										add(ruledayNum, position14)
									}
									depth--
									add(ruledate, position9)
								}
								if buffer[position] != rune(' ') {
									goto l6
								}
								position++
								if !_rules[ruletime]() {
									goto l6
								}
								depth--
								add(rulePegText, position8)
							}
							{
								add(ruleAction20, position)
							}
							depth--
							add(ruletimestamp24, position7)
						}
						goto l5
					l6:
						position, tokenIndex, depth = position5, tokenIndex5, depth5
						{
							position18 := position
							depth++
							{
								position19 := position
								depth++
								{
									position20 := position
									depth++
									{
										position21 := position
										depth++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l0
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l0
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l0
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l0
										}
										position++
										depth--
										add(ruledigit4, position21)
									}
									if buffer[position] != rune('-') {
										goto l0
									}
									position++
									if !_rules[ruledigit2]() {
										goto l0
									}
									if buffer[position] != rune('-') {
										goto l0
									}
									position++
									if !_rules[ruledigit2]() {
										goto l0
									}
									if buffer[position] != rune('T') {
										goto l0
									}
									position++
									if !_rules[ruletime]() {
										goto l0
									}
									{
										position22, tokenIndex22, depth22 := position, tokenIndex, depth
										{
											position24 := position
											depth++
											if buffer[position] != rune('+') {
												goto l22
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l22
											}
											position++
										l25:
											{
												position26, tokenIndex26, depth26 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l26
												}
												position++
												goto l25
											l26:
												position, tokenIndex, depth = position26, tokenIndex26, depth26
											}
											depth--
											add(ruletz, position24)
										}
										goto l23
									l22:
										position, tokenIndex, depth = position22, tokenIndex22, depth22
									}
								l23:
									depth--
									add(ruledatetime26, position20)
								}
								depth--
								add(rulePegText, position19)
							}
							{
								add(ruleAction21, position)
							}
							depth--
							add(ruletimestamp26, position18)
						}
					}
				l5:
					{
						position28, tokenIndex28, depth28 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l28
						}
						goto l29
					l28:
						position, tokenIndex, depth = position28, tokenIndex28, depth28
					}
				l29:
					depth--
					add(ruleTimestamp, position4)
				}
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					{
						position32 := position
						depth++
						{
							position33 := position
							depth++
							{
								position34, tokenIndex34, depth34 := position, tokenIndex, depth
								if buffer[position] != rune('I') {
									goto l35
								}
								position++
								goto l34
							l35:
								position, tokenIndex, depth = position34, tokenIndex34, depth34
								if buffer[position] != rune('D') {
									goto l30
								}
								position++
							}
						l34:
							depth--
							add(rulePegText, position33)
						}
						if buffer[position] != rune(' ') {
							goto l30
						}
						position++
						{
							add(ruleAction0, position)
						}
						depth--
						add(ruleSeverity, position32)
					}
					goto l31
				l30:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
				}
			l31:
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					{
						position39 := position
						depth++
						{
							position40 := position
							depth++
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l37
							}
							position++
						l41:
							{
								position42, tokenIndex42, depth42 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l42
								}
								position++
								goto l41
							l42:
								position, tokenIndex, depth = position42, tokenIndex42, depth42
							}
							depth--
							add(rulePegText, position40)
						}
						if buffer[position] != rune(' ') {
							goto l37
						}
						position++
					l43:
						{
							position44, tokenIndex44, depth44 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l44
							}
							position++
							goto l43
						l44:
							position, tokenIndex, depth = position44, tokenIndex44, depth44
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleComponent, position39)
					}
					goto l38
				l37:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
				}
			l38:
				{
					position46 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position47 := position
						depth++
						{
							position50 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position52, tokenIndex52, depth52 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l53
										}
										position++
										goto l52
									l53:
										position, tokenIndex, depth = position52, tokenIndex52, depth52
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l52:
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l0
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l0
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l0
									}
									position++
									break
								}
							}

							depth--
							add(ruleletterOrDigit, position50)
						}
					l48:
						{
							position49, tokenIndex49, depth49 := position, tokenIndex, depth
							{
								position54 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position56, tokenIndex56, depth56 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l57
											}
											position++
											goto l56
										l57:
											position, tokenIndex, depth = position56, tokenIndex56, depth56
											if buffer[position] != rune('$') {
												goto l49
											}
											position++
										}
									l56:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l49
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l49
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l49
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position54)
							}
							goto l48
						l49:
							position, tokenIndex, depth = position49, tokenIndex49, depth49
						}
						depth--
						add(rulePegText, position47)
					}
					if buffer[position] != rune(']') {
						goto l0
					}
					position++
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction2, position)
					}
					depth--
					add(ruleContext, position46)
				}
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					{
						position61 := position
						depth++
						{
							position62 := position
							depth++
							{
								position63 := position
								depth++
								if buffer[position] != rune('w') {
									goto l59
								}
								position++
								if buffer[position] != rune('a') {
									goto l59
								}
								position++
								if buffer[position] != rune('r') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('g') {
									goto l59
								}
								position++
								if buffer[position] != rune(':') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('l') {
									goto l59
								}
								position++
								if buffer[position] != rune('o') {
									goto l59
								}
								position++
								if buffer[position] != rune('g') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('l') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('a') {
									goto l59
								}
								position++
								if buffer[position] != rune('t') {
									goto l59
								}
								position++
								if buffer[position] != rune('t') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune('m') {
									goto l59
								}
								position++
								if buffer[position] != rune('p') {
									goto l59
								}
								position++
								if buffer[position] != rune('t') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune('d') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('(') {
									goto l59
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l59
								}
								position++
							l64:
								{
									position65, tokenIndex65, depth65 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l65
									}
									position++
									goto l64
								l65:
									position, tokenIndex, depth = position65, tokenIndex65, depth65
								}
								if buffer[position] != rune('k') {
									goto l59
								}
								position++
								if buffer[position] != rune(')') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('o') {
									goto l59
								}
								position++
								if buffer[position] != rune('v') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune('r') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('m') {
									goto l59
								}
								position++
								if buffer[position] != rune('a') {
									goto l59
								}
								position++
								if buffer[position] != rune('x') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('s') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('z') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('(') {
									goto l59
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l59
								}
								position++
							l66:
								{
									position67, tokenIndex67, depth67 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l67
									}
									position++
									goto l66
								l67:
									position, tokenIndex, depth = position67, tokenIndex67, depth67
								}
								if buffer[position] != rune('k') {
									goto l59
								}
								position++
								if buffer[position] != rune(')') {
									goto l59
								}
								position++
								if buffer[position] != rune(',') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('p') {
									goto l59
								}
								position++
								if buffer[position] != rune('r') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('t') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('g') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('b') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune('g') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('i') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('g') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('a') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('d') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('e') {
									goto l59
								}
								position++
								if buffer[position] != rune('n') {
									goto l59
								}
								position++
								if buffer[position] != rune('d') {
									goto l59
								}
								position++
								if buffer[position] != rune(' ') {
									goto l59
								}
								position++
								if buffer[position] != rune('.') {
									goto l59
								}
								position++
								if buffer[position] != rune('.') {
									goto l59
								}
								position++
								if buffer[position] != rune('.') {
									goto l59
								}
								position++
								depth--
								add(ruleloglineSizeWarning, position63)
							}
							depth--
							add(rulePegText, position62)
						}
						if buffer[position] != rune(' ') {
							goto l59
						}
						position++
						{
							add(ruleAction4, position)
						}
						depth--
						add(ruleWarning, position61)
					}
					goto l60
				l59:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
				}
			l60:
				{
					position69 := position
					depth++
					{
						position70 := position
						depth++
						{
							position73, tokenIndex73, depth73 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l74
							}
							position++
							goto l73
						l74:
							position, tokenIndex, depth = position73, tokenIndex73, depth73
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l0
							}
							position++
						}
					l73:
					l71:
						{
							position72, tokenIndex72, depth72 := position, tokenIndex, depth
							{
								position75, tokenIndex75, depth75 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l76
								}
								position++
								goto l75
							l76:
								position, tokenIndex, depth = position75, tokenIndex75, depth75
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l72
								}
								position++
							}
						l75:
							goto l71
						l72:
							position, tokenIndex, depth = position72, tokenIndex72, depth72
						}
						depth--
						add(rulePegText, position70)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleOp, position69)
				}
				{
					position78 := position
					depth++
					{
						position79 := position
						depth++
					l80:
						{
							position81, tokenIndex81, depth81 := position, tokenIndex, depth
							{
								position82 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l81
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l81
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l81
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l81
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l81
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l81
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position82)
							}
							goto l80
						l81:
							position, tokenIndex, depth = position81, tokenIndex81, depth81
						}
						depth--
						add(rulePegText, position79)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction5, position)
					}
					depth--
					add(ruleNS, position78)
				}
			l85:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l86
					}
					goto l85
				l86:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
				}
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					{
						position89 := position
						depth++
						if buffer[position] != rune('l') {
							goto l87
						}
						position++
						if buffer[position] != rune('o') {
							goto l87
						}
						position++
						if buffer[position] != rune('c') {
							goto l87
						}
						position++
						if buffer[position] != rune('k') {
							goto l87
						}
						position++
						if buffer[position] != rune('s') {
							goto l87
						}
						position++
						if buffer[position] != rune('(') {
							goto l87
						}
						position++
						if buffer[position] != rune('m') {
							goto l87
						}
						position++
						if buffer[position] != rune('i') {
							goto l87
						}
						position++
						if buffer[position] != rune('c') {
							goto l87
						}
						position++
						if buffer[position] != rune('r') {
							goto l87
						}
						position++
						if buffer[position] != rune('o') {
							goto l87
						}
						position++
						if buffer[position] != rune('s') {
							goto l87
						}
						position++
						if buffer[position] != rune(')') {
							goto l87
						}
						position++
						{
							position90, tokenIndex90, depth90 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l90
							}
							goto l91
						l90:
							position, tokenIndex, depth = position90, tokenIndex90, depth90
						}
					l91:
					l92:
						{
							position93, tokenIndex93, depth93 := position, tokenIndex, depth
							{
								position94 := position
								depth++
								{
									position95 := position
									depth++
									{
										switch buffer[position] {
										case 'R':
											if buffer[position] != rune('R') {
												goto l93
											}
											position++
											break
										case 'r':
											if buffer[position] != rune('r') {
												goto l93
											}
											position++
											break
										default:
											{
												position97, tokenIndex97, depth97 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l98
												}
												position++
												goto l97
											l98:
												position, tokenIndex, depth = position97, tokenIndex97, depth97
												if buffer[position] != rune('W') {
													goto l93
												}
												position++
											}
										l97:
											break
										}
									}

									depth--
									add(rulePegText, position95)
								}
								{
									add(ruleAction6, position)
								}
								if buffer[position] != rune(':') {
									goto l93
								}
								position++
								if !_rules[ruleNumeric]() {
									goto l93
								}
								{
									position100, tokenIndex100, depth100 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l100
									}
									goto l101
								l100:
									position, tokenIndex, depth = position100, tokenIndex100, depth100
								}
							l101:
								{
									add(ruleAction7, position)
								}
								depth--
								add(rulelock, position94)
							}
							goto l92
						l93:
							position, tokenIndex, depth = position93, tokenIndex93, depth93
						}
						depth--
						add(ruleLocks, position89)
					}
					goto l88
				l87:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
				}
			l88:
			l103:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l104
					}
					goto l103
				l104:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
				}
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					{
						position107 := position
						depth++
						{
							position108 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l105
							}
							position++
						l109:
							{
								position110, tokenIndex110, depth110 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l110
								}
								position++
								goto l109
							l110:
								position, tokenIndex, depth = position110, tokenIndex110, depth110
							}
							depth--
							add(rulePegText, position108)
						}
						if buffer[position] != rune('m') {
							goto l105
						}
						position++
						if buffer[position] != rune('s') {
							goto l105
						}
						position++
						{
							add(ruleAction8, position)
						}
						depth--
						add(ruleDuration, position107)
					}
					goto l106
				l105:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
				}
			l106:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					{
						position114 := position
						depth++
						{
							position115 := position
							depth++
							if !matchDot() {
								goto l112
							}
						l116:
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if !matchDot() {
									goto l117
								}
								goto l116
							l117:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
							}
							depth--
							add(rulePegText, position115)
						}
						{
							add(ruleAction22, position)
						}
						depth--
						add(ruleextra, position114)
					}
					goto l113
				l112:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
				}
			l113:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if !matchDot() {
						goto l119
					}
					goto l0
				l119:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
				}
				depth--
				add(ruleMongoLogLine, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Timestamp <- <((timestamp24 / timestamp26) S?)> */
		nil,
		/* 2 Severity <- <(<('I' / 'D')> ' ' Action0)> */
		nil,
		/* 3 Component <- <(<[A-Z]+> ' '+ Action1)> */
		nil,
		/* 4 Context <- <('[' <letterOrDigit+> ']' ' ' Action2)> */
		nil,
		/* 5 Op <- <(<([a-z] / [A-Z])+> ' ' Action3)> */
		nil,
		/* 6 Warning <- <(<loglineSizeWarning> ' ' Action4)> */
		nil,
		/* 7 loglineSizeWarning <- <('w' 'a' 'r' 'n' 'i' 'n' 'g' ':' ' ' 'l' 'o' 'g' ' ' 'l' 'i' 'n' 'e' ' ' 'a' 't' 't' 'e' 'm' 'p' 't' 'e' 'd' ' ' '(' [0-9]+ ('k' ')' ' ' 'o' 'v' 'e' 'r' ' ' 'm' 'a' 'x' ' ' 's' 'i' 'z' 'e' ' ' '(') [0-9]+ ('k' ')' ',' ' ' 'p' 'r' 'i' 'n' 't' 'i' 'n' 'g' ' ' 'b' 'e' 'g' 'i' 'n' 'n' 'i' 'n' 'g' ' ' 'a' 'n' 'd' ' ' 'e' 'n' 'd' ' ' '.' '.' '.'))> */
		nil,
		/* 8 LineField <- <((exceptionField / commandField / planSummaryField / plainField) S?)> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					{
						position131 := position
						depth++
						if buffer[position] != rune('e') {
							goto l130
						}
						position++
						if buffer[position] != rune('x') {
							goto l130
						}
						position++
						if buffer[position] != rune('c') {
							goto l130
						}
						position++
						if buffer[position] != rune('e') {
							goto l130
						}
						position++
						if buffer[position] != rune('p') {
							goto l130
						}
						position++
						if buffer[position] != rune('t') {
							goto l130
						}
						position++
						if buffer[position] != rune('i') {
							goto l130
						}
						position++
						if buffer[position] != rune('o') {
							goto l130
						}
						position++
						if buffer[position] != rune('n') {
							goto l130
						}
						position++
						if buffer[position] != rune(':') {
							goto l130
						}
						position++
						{
							add(ruleAction18, position)
						}
						{
							position133 := position
							depth++
							{
								position136, tokenIndex136, depth136 := position, tokenIndex, depth
								if !matchDot() {
									goto l130
								}
								{
									position137, tokenIndex137, depth137 := position, tokenIndex, depth
									if buffer[position] != rune('c') {
										goto l137
									}
									position++
									if buffer[position] != rune('o') {
										goto l137
									}
									position++
									if buffer[position] != rune('d') {
										goto l137
									}
									position++
									if buffer[position] != rune('e') {
										goto l137
									}
									position++
									if buffer[position] != rune(':') {
										goto l137
									}
									position++
									goto l130
								l137:
									position, tokenIndex, depth = position137, tokenIndex137, depth137
								}
								position, tokenIndex, depth = position136, tokenIndex136, depth136
							}
							if !matchDot() {
								goto l130
							}
						l134:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								{
									position138, tokenIndex138, depth138 := position, tokenIndex, depth
									if !matchDot() {
										goto l135
									}
									{
										position139, tokenIndex139, depth139 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l139
										}
										position++
										if buffer[position] != rune('o') {
											goto l139
										}
										position++
										if buffer[position] != rune('d') {
											goto l139
										}
										position++
										if buffer[position] != rune('e') {
											goto l139
										}
										position++
										if buffer[position] != rune(':') {
											goto l139
										}
										position++
										goto l135
									l139:
										position, tokenIndex, depth = position139, tokenIndex139, depth139
									}
									position, tokenIndex, depth = position138, tokenIndex138, depth138
								}
								if !matchDot() {
									goto l135
								}
								goto l134
							l135:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
							}
							depth--
							add(rulePegText, position133)
						}
						{
							position140, tokenIndex140, depth140 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l140
							}
							goto l141
						l140:
							position, tokenIndex, depth = position140, tokenIndex140, depth140
						}
					l141:
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleexceptionField, position131)
					}
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					{
						position144 := position
						depth++
						if buffer[position] != rune('c') {
							goto l143
						}
						position++
						if buffer[position] != rune('o') {
							goto l143
						}
						position++
						if buffer[position] != rune('m') {
							goto l143
						}
						position++
						if buffer[position] != rune('m') {
							goto l143
						}
						position++
						if buffer[position] != rune('a') {
							goto l143
						}
						position++
						if buffer[position] != rune('n') {
							goto l143
						}
						position++
						if buffer[position] != rune('d') {
							goto l143
						}
						position++
						if buffer[position] != rune(':') {
							goto l143
						}
						position++
						if buffer[position] != rune(' ') {
							goto l143
						}
						position++
						{
							position145 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l143
							}
						l146:
							{
								position147, tokenIndex147, depth147 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l147
								}
								goto l146
							l147:
								position, tokenIndex, depth = position147, tokenIndex147, depth147
							}
							depth--
							add(rulePegText, position145)
						}
						{
							position148, tokenIndex148, depth148 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l148
							}
							goto l149
						l148:
							position, tokenIndex, depth = position148, tokenIndex148, depth148
						}
					l149:
						{
							add(ruleAction11, position)
						}
						if !_rules[ruleLineValue]() {
							goto l143
						}
						{
							add(ruleAction12, position)
						}
						depth--
						add(rulecommandField, position144)
					}
					goto l129
				l143:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					{
						position153 := position
						depth++
						if buffer[position] != rune('p') {
							goto l152
						}
						position++
						if buffer[position] != rune('l') {
							goto l152
						}
						position++
						if buffer[position] != rune('a') {
							goto l152
						}
						position++
						if buffer[position] != rune('n') {
							goto l152
						}
						position++
						if buffer[position] != rune('S') {
							goto l152
						}
						position++
						if buffer[position] != rune('u') {
							goto l152
						}
						position++
						if buffer[position] != rune('m') {
							goto l152
						}
						position++
						if buffer[position] != rune('m') {
							goto l152
						}
						position++
						if buffer[position] != rune('a') {
							goto l152
						}
						position++
						if buffer[position] != rune('r') {
							goto l152
						}
						position++
						if buffer[position] != rune('y') {
							goto l152
						}
						position++
						if buffer[position] != rune(':') {
							goto l152
						}
						position++
						if buffer[position] != rune(' ') {
							goto l152
						}
						position++
						{
							add(ruleAction13, position)
						}
						{
							position155 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l152
							}
						l156:
							{
								position157, tokenIndex157, depth157 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l157
								}
								position++
								if buffer[position] != rune(' ') {
									goto l157
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l157
								}
								goto l156
							l157:
								position, tokenIndex, depth = position157, tokenIndex157, depth157
							}
							depth--
							add(ruleplanSummaryElements, position155)
						}
						{
							add(ruleAction14, position)
						}
						depth--
						add(ruleplanSummaryField, position153)
					}
					goto l129
				l152:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					{
						position159 := position
						depth++
						{
							position160 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l127
							}
						l161:
							{
								position162, tokenIndex162, depth162 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l162
								}
								goto l161
							l162:
								position, tokenIndex, depth = position162, tokenIndex162, depth162
							}
							depth--
							add(rulePegText, position160)
						}
						if buffer[position] != rune(':') {
							goto l127
						}
						position++
						{
							position163, tokenIndex163, depth163 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l163
							}
							goto l164
						l163:
							position, tokenIndex, depth = position163, tokenIndex163, depth163
						}
					l164:
						{
							add(ruleAction9, position)
						}
						if !_rules[ruleLineValue]() {
							goto l127
						}
						{
							add(ruleAction10, position)
						}
						depth--
						add(ruleplainField, position159)
					}
				}
			l129:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l167
					}
					goto l168
				l167:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
				}
			l168:
				depth--
				add(ruleLineField, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
			return false
		},
		/* 9 NS <- <(<nsChar*> ' ' Action5)> */
		nil,
		/* 10 Locks <- <('l' 'o' 'c' 'k' 's' '(' 'm' 'i' 'c' 'r' 'o' 's' ')' S? lock*)> */
		nil,
		/* 11 lock <- <(<((&('R') 'R') | (&('r') 'r') | (&('W' | 'w') ('w' / 'W')))> Action6 ':' Numeric S? Action7)> */
		nil,
		/* 12 Duration <- <(<[0-9]+> ('m' 's') Action8)> */
		nil,
		/* 13 plainField <- <(<fieldChar+> ':' S? Action9 LineValue Action10)> */
		nil,
		/* 14 commandField <- <('c' 'o' 'm' 'm' 'a' 'n' 'd' ':' ' ' <fieldChar+> S? Action11 LineValue Action12)> */
		nil,
		/* 15 planSummaryField <- <('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y' ':' ' ' Action13 planSummaryElements Action14)> */
		nil,
		/* 16 planSummaryElements <- <(planSummaryElem (',' ' ' planSummaryElem)*)> */
		nil,
		/* 17 planSummaryElem <- <(<planSummaryStage> Action15 planSummary)> */
		func() bool {
			position177, tokenIndex177, depth177 := position, tokenIndex, depth
			{
				position178 := position
				depth++
				{
					position179 := position
					depth++
					{
						position180 := position
						depth++
						{
							position181, tokenIndex181, depth181 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l182
							}
							position++
							if buffer[position] != rune('N') {
								goto l182
							}
							position++
							if buffer[position] != rune('D') {
								goto l182
							}
							position++
							if buffer[position] != rune('_') {
								goto l182
							}
							position++
							if buffer[position] != rune('H') {
								goto l182
							}
							position++
							if buffer[position] != rune('A') {
								goto l182
							}
							position++
							if buffer[position] != rune('S') {
								goto l182
							}
							position++
							if buffer[position] != rune('H') {
								goto l182
							}
							position++
							goto l181
						l182:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('C') {
								goto l183
							}
							position++
							if buffer[position] != rune('A') {
								goto l183
							}
							position++
							if buffer[position] != rune('C') {
								goto l183
							}
							position++
							if buffer[position] != rune('H') {
								goto l183
							}
							position++
							if buffer[position] != rune('E') {
								goto l183
							}
							position++
							if buffer[position] != rune('D') {
								goto l183
							}
							position++
							if buffer[position] != rune('_') {
								goto l183
							}
							position++
							if buffer[position] != rune('P') {
								goto l183
							}
							position++
							if buffer[position] != rune('L') {
								goto l183
							}
							position++
							if buffer[position] != rune('A') {
								goto l183
							}
							position++
							if buffer[position] != rune('N') {
								goto l183
							}
							position++
							goto l181
						l183:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('C') {
								goto l184
							}
							position++
							if buffer[position] != rune('O') {
								goto l184
							}
							position++
							if buffer[position] != rune('L') {
								goto l184
							}
							position++
							if buffer[position] != rune('L') {
								goto l184
							}
							position++
							if buffer[position] != rune('S') {
								goto l184
							}
							position++
							if buffer[position] != rune('C') {
								goto l184
							}
							position++
							if buffer[position] != rune('A') {
								goto l184
							}
							position++
							if buffer[position] != rune('N') {
								goto l184
							}
							position++
							goto l181
						l184:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('C') {
								goto l185
							}
							position++
							if buffer[position] != rune('O') {
								goto l185
							}
							position++
							if buffer[position] != rune('U') {
								goto l185
							}
							position++
							if buffer[position] != rune('N') {
								goto l185
							}
							position++
							if buffer[position] != rune('T') {
								goto l185
							}
							position++
							goto l181
						l185:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('D') {
								goto l186
							}
							position++
							if buffer[position] != rune('E') {
								goto l186
							}
							position++
							if buffer[position] != rune('L') {
								goto l186
							}
							position++
							if buffer[position] != rune('E') {
								goto l186
							}
							position++
							if buffer[position] != rune('T') {
								goto l186
							}
							position++
							if buffer[position] != rune('E') {
								goto l186
							}
							position++
							goto l181
						l186:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('G') {
								goto l187
							}
							position++
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
							if buffer[position] != rune('O') {
								goto l187
							}
							position++
							if buffer[position] != rune('_') {
								goto l187
							}
							position++
							if buffer[position] != rune('N') {
								goto l187
							}
							position++
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
							if buffer[position] != rune('A') {
								goto l187
							}
							position++
							if buffer[position] != rune('R') {
								goto l187
							}
							position++
							if buffer[position] != rune('_') {
								goto l187
							}
							position++
							if buffer[position] != rune('2') {
								goto l187
							}
							position++
							if buffer[position] != rune('D') {
								goto l187
							}
							position++
							goto l181
						l187:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('G') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							if buffer[position] != rune('O') {
								goto l188
							}
							position++
							if buffer[position] != rune('_') {
								goto l188
							}
							position++
							if buffer[position] != rune('N') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							if buffer[position] != rune('A') {
								goto l188
							}
							position++
							if buffer[position] != rune('R') {
								goto l188
							}
							position++
							if buffer[position] != rune('_') {
								goto l188
							}
							position++
							if buffer[position] != rune('2') {
								goto l188
							}
							position++
							if buffer[position] != rune('D') {
								goto l188
							}
							position++
							if buffer[position] != rune('S') {
								goto l188
							}
							position++
							if buffer[position] != rune('P') {
								goto l188
							}
							position++
							if buffer[position] != rune('H') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							if buffer[position] != rune('R') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							goto l181
						l188:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('I') {
								goto l189
							}
							position++
							if buffer[position] != rune('D') {
								goto l189
							}
							position++
							if buffer[position] != rune('H') {
								goto l189
							}
							position++
							if buffer[position] != rune('A') {
								goto l189
							}
							position++
							if buffer[position] != rune('C') {
								goto l189
							}
							position++
							if buffer[position] != rune('K') {
								goto l189
							}
							position++
							goto l181
						l189:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('S') {
								goto l190
							}
							position++
							if buffer[position] != rune('O') {
								goto l190
							}
							position++
							if buffer[position] != rune('R') {
								goto l190
							}
							position++
							if buffer[position] != rune('T') {
								goto l190
							}
							position++
							if buffer[position] != rune('_') {
								goto l190
							}
							position++
							if buffer[position] != rune('M') {
								goto l190
							}
							position++
							if buffer[position] != rune('E') {
								goto l190
							}
							position++
							if buffer[position] != rune('R') {
								goto l190
							}
							position++
							if buffer[position] != rune('G') {
								goto l190
							}
							position++
							if buffer[position] != rune('E') {
								goto l190
							}
							position++
							goto l181
						l190:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('S') {
								goto l191
							}
							position++
							if buffer[position] != rune('H') {
								goto l191
							}
							position++
							if buffer[position] != rune('A') {
								goto l191
							}
							position++
							if buffer[position] != rune('R') {
								goto l191
							}
							position++
							if buffer[position] != rune('D') {
								goto l191
							}
							position++
							if buffer[position] != rune('I') {
								goto l191
							}
							position++
							if buffer[position] != rune('N') {
								goto l191
							}
							position++
							if buffer[position] != rune('G') {
								goto l191
							}
							position++
							if buffer[position] != rune('_') {
								goto l191
							}
							position++
							if buffer[position] != rune('F') {
								goto l191
							}
							position++
							if buffer[position] != rune('I') {
								goto l191
							}
							position++
							if buffer[position] != rune('L') {
								goto l191
							}
							position++
							if buffer[position] != rune('T') {
								goto l191
							}
							position++
							if buffer[position] != rune('E') {
								goto l191
							}
							position++
							if buffer[position] != rune('R') {
								goto l191
							}
							position++
							goto l181
						l191:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('S') {
								goto l192
							}
							position++
							if buffer[position] != rune('K') {
								goto l192
							}
							position++
							if buffer[position] != rune('I') {
								goto l192
							}
							position++
							if buffer[position] != rune('P') {
								goto l192
							}
							position++
							goto l181
						l192:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('S') {
								goto l193
							}
							position++
							if buffer[position] != rune('O') {
								goto l193
							}
							position++
							if buffer[position] != rune('R') {
								goto l193
							}
							position++
							if buffer[position] != rune('T') {
								goto l193
							}
							position++
							goto l181
						l193:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('X') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('B') {
										goto l177
									}
									position++
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									if buffer[position] != rune('L') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									if buffer[position] != rune('_') {
										goto l177
									}
									position++
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									if buffer[position] != rune('R') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('J') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('R') {
										goto l177
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('L') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('_') {
										goto l177
									}
									position++
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									if buffer[position] != rune('L') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('M') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									if buffer[position] != rune('_') {
										goto l177
									}
									position++
									if buffer[position] != rune('M') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('X') {
										goto l177
									}
									position++
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l177
									}
									position++
									if buffer[position] != rune('R') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('P') {
										goto l177
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('H') {
										goto l177
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('F') {
										goto l177
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('I') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('U') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('_') {
										goto l177
									}
									position++
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									if buffer[position] != rune('C') {
										goto l177
									}
									position++
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l177
									}
									position++
									if buffer[position] != rune('N') {
										goto l177
									}
									position++
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									if buffer[position] != rune('_') {
										goto l177
									}
									position++
									if buffer[position] != rune('S') {
										goto l177
									}
									position++
									if buffer[position] != rune('O') {
										goto l177
									}
									position++
									if buffer[position] != rune('R') {
										goto l177
									}
									position++
									if buffer[position] != rune('T') {
										goto l177
									}
									position++
									if buffer[position] != rune('E') {
										goto l177
									}
									position++
									if buffer[position] != rune('D') {
										goto l177
									}
									position++
									break
								}
							}

						}
					l181:
						depth--
						add(ruleplanSummaryStage, position180)
					}
					depth--
					add(rulePegText, position179)
				}
				{
					add(ruleAction15, position)
				}
				{
					position196 := position
					depth++
					{
						position197, tokenIndex197, depth197 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l198
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l198
						}
						{
							add(ruleAction16, position)
						}
						goto l197
					l198:
						position, tokenIndex, depth = position197, tokenIndex197, depth197
						{
							add(ruleAction17, position)
						}
					}
				l197:
					depth--
					add(ruleplanSummary, position196)
				}
				depth--
				add(ruleplanSummaryElem, position178)
			}
			return true
		l177:
			position, tokenIndex, depth = position177, tokenIndex177, depth177
			return false
		},
		/* 18 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 19 planSummary <- <((' ' LineValue Action16) / Action17)> */
		nil,
		/* 20 exceptionField <- <('e' 'x' 'c' 'e' 'p' 't' 'i' 'o' 'n' ':' Action18 <(&(. !('c' 'o' 'd' 'e' ':')) .)+> S? Action19)> */
		nil,
		/* 21 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position204, tokenIndex204, depth204 := position, tokenIndex, depth
			{
				position205 := position
				depth++
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l207
					}
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if !_rules[ruleNumeric]() {
						goto l204
					}
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l208
					}
					goto l209
				l208:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
				}
			l209:
				depth--
				add(ruleLineValue, position205)
			}
			return true
		l204:
			position, tokenIndex, depth = position204, tokenIndex204, depth204
			return false
		},
		/* 22 timestamp24 <- <(<(date ' ' time)> Action20)> */
		nil,
		/* 23 timestamp26 <- <(<datetime26> Action21)> */
		nil,
		/* 24 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 25 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 26 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l214
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l214
				}
				position++
				depth--
				add(ruledigit2, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
			return false
		},
		/* 27 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 28 tz <- <('+' [0-9]+)> */
		nil,
		/* 29 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l218
					}
					depth--
					add(rulehour, position220)
				}
				if buffer[position] != rune(':') {
					goto l218
				}
				position++
				{
					position221 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l218
					}
					depth--
					add(ruleminute, position221)
				}
				if buffer[position] != rune(':') {
					goto l218
				}
				position++
				{
					position222 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l218
					}
					depth--
					add(rulesecond, position222)
				}
				if buffer[position] != rune('.') {
					goto l218
				}
				position++
				{
					position223 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l218
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l218
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l218
					}
					position++
					depth--
					add(rulemillisecond, position223)
				}
				depth--
				add(ruletime, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 30 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 31 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 32 dayNum <- <([0-9] [0-9]?)> */
		nil,
		/* 33 hour <- <digit2> */
		nil,
		/* 34 minute <- <digit2> */
		nil,
		/* 35 second <- <digit2> */
		nil,
		/* 36 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 37 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 38 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 39 extra <- <(<.+> Action22)> */
		nil,
		/* 40 S <- <' '+> */
		func() bool {
			position234, tokenIndex234, depth234 := position, tokenIndex, depth
			{
				position235 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l234
				}
				position++
			l236:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
				}
				depth--
				add(ruleS, position235)
			}
			return true
		l234:
			position, tokenIndex, depth = position234, tokenIndex234, depth234
			return false
		},
		/* 41 Doc <- <('{' Action23 DocElements? '}' Action24)> */
		func() bool {
			position238, tokenIndex238, depth238 := position, tokenIndex, depth
			{
				position239 := position
				depth++
				if buffer[position] != rune('{') {
					goto l238
				}
				position++
				{
					add(ruleAction23, position)
				}
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					{
						position243 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l241
						}
					l244:
						{
							position245, tokenIndex245, depth245 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l245
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l245
							}
							goto l244
						l245:
							position, tokenIndex, depth = position245, tokenIndex245, depth245
						}
						depth--
						add(ruleDocElements, position243)
					}
					goto l242
				l241:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
				}
			l242:
				if buffer[position] != rune('}') {
					goto l238
				}
				position++
				{
					add(ruleAction24, position)
				}
				depth--
				add(ruleDoc, position239)
			}
			return true
		l238:
			position, tokenIndex, depth = position238, tokenIndex238, depth238
			return false
		},
		/* 42 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 43 DocElem <- <(S? Field S? Value S? Action25)> */
		func() bool {
			position248, tokenIndex248, depth248 := position, tokenIndex, depth
			{
				position249 := position
				depth++
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l250
					}
					goto l251
				l250:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
				}
			l251:
				{
					position252 := position
					depth++
					{
						position253 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l248
						}
					l254:
						{
							position255, tokenIndex255, depth255 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l255
							}
							goto l254
						l255:
							position, tokenIndex, depth = position255, tokenIndex255, depth255
						}
						depth--
						add(rulePegText, position253)
					}
					if buffer[position] != rune(':') {
						goto l248
					}
					position++
					{
						add(ruleAction29, position)
					}
					depth--
					add(ruleField, position252)
				}
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l257
					}
					goto l258
				l257:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
				}
			l258:
				if !_rules[ruleValue]() {
					goto l248
				}
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l259
					}
					goto l260
				l259:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
				}
			l260:
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleDocElem, position249)
			}
			return true
		l248:
			position, tokenIndex, depth = position248, tokenIndex248, depth248
			return false
		},
		/* 44 List <- <('[' Action26 ListElements? ']' Action27)> */
		nil,
		/* 45 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 46 ListElem <- <(S? Value S? Action28)> */
		func() bool {
			position264, tokenIndex264, depth264 := position, tokenIndex, depth
			{
				position265 := position
				depth++
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l266
					}
					goto l267
				l266:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
				}
			l267:
				if !_rules[ruleValue]() {
					goto l264
				}
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l268
					}
					goto l269
				l268:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
				}
			l269:
				{
					add(ruleAction28, position)
				}
				depth--
				add(ruleListElem, position265)
			}
			return true
		l264:
			position, tokenIndex, depth = position264, tokenIndex264, depth264
			return false
		},
		/* 47 Field <- <(<fieldChar+> ':' Action29)> */
		nil,
		/* 48 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					{
						position276 := position
						depth++
						if buffer[position] != rune('n') {
							goto l275
						}
						position++
						if buffer[position] != rune('u') {
							goto l275
						}
						position++
						if buffer[position] != rune('l') {
							goto l275
						}
						position++
						if buffer[position] != rune('l') {
							goto l275
						}
						position++
						{
							add(ruleAction32, position)
						}
						depth--
						add(ruleNull, position276)
					}
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					{
						position279 := position
						depth++
						if buffer[position] != rune('M') {
							goto l278
						}
						position++
						if buffer[position] != rune('i') {
							goto l278
						}
						position++
						if buffer[position] != rune('n') {
							goto l278
						}
						position++
						if buffer[position] != rune('K') {
							goto l278
						}
						position++
						if buffer[position] != rune('e') {
							goto l278
						}
						position++
						if buffer[position] != rune('y') {
							goto l278
						}
						position++
						{
							add(ruleAction42, position)
						}
						depth--
						add(ruleMinKey, position279)
					}
					goto l274
				l278:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					{
						switch buffer[position] {
						case 'M':
							{
								position282 := position
								depth++
								if buffer[position] != rune('M') {
									goto l272
								}
								position++
								if buffer[position] != rune('a') {
									goto l272
								}
								position++
								if buffer[position] != rune('x') {
									goto l272
								}
								position++
								if buffer[position] != rune('K') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('y') {
									goto l272
								}
								position++
								{
									add(ruleAction43, position)
								}
								depth--
								add(ruleMaxKey, position282)
							}
							break
						case 'u':
							{
								position284 := position
								depth++
								if buffer[position] != rune('u') {
									goto l272
								}
								position++
								if buffer[position] != rune('n') {
									goto l272
								}
								position++
								if buffer[position] != rune('d') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('f') {
									goto l272
								}
								position++
								if buffer[position] != rune('i') {
									goto l272
								}
								position++
								if buffer[position] != rune('n') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('d') {
									goto l272
								}
								position++
								{
									add(ruleAction44, position)
								}
								depth--
								add(ruleUndefined, position284)
							}
							break
						case 'N':
							{
								position286 := position
								depth++
								if buffer[position] != rune('N') {
									goto l272
								}
								position++
								if buffer[position] != rune('u') {
									goto l272
								}
								position++
								if buffer[position] != rune('m') {
									goto l272
								}
								position++
								if buffer[position] != rune('b') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('r') {
									goto l272
								}
								position++
								if buffer[position] != rune('L') {
									goto l272
								}
								position++
								if buffer[position] != rune('o') {
									goto l272
								}
								position++
								if buffer[position] != rune('n') {
									goto l272
								}
								position++
								if buffer[position] != rune('g') {
									goto l272
								}
								position++
								if buffer[position] != rune('(') {
									goto l272
								}
								position++
								{
									position287 := position
									depth++
									{
										position290, tokenIndex290, depth290 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l290
										}
										position++
										goto l272
									l290:
										position, tokenIndex, depth = position290, tokenIndex290, depth290
									}
									if !matchDot() {
										goto l272
									}
								l288:
									{
										position289, tokenIndex289, depth289 := position, tokenIndex, depth
										{
											position291, tokenIndex291, depth291 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l291
											}
											position++
											goto l289
										l291:
											position, tokenIndex, depth = position291, tokenIndex291, depth291
										}
										if !matchDot() {
											goto l289
										}
										goto l288
									l289:
										position, tokenIndex, depth = position289, tokenIndex289, depth289
									}
									depth--
									add(rulePegText, position287)
								}
								if buffer[position] != rune(')') {
									goto l272
								}
								position++
								{
									add(ruleAction41, position)
								}
								depth--
								add(ruleNumberLong, position286)
							}
							break
						case '/':
							{
								position293 := position
								depth++
								if buffer[position] != rune('/') {
									goto l272
								}
								position++
								{
									position294 := position
									depth++
									{
										position295 := position
										depth++
										{
											position298 := position
											depth++
											{
												position299, tokenIndex299, depth299 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l299
												}
												position++
												goto l272
											l299:
												position, tokenIndex, depth = position299, tokenIndex299, depth299
											}
											if !matchDot() {
												goto l272
											}
											depth--
											add(ruleregexChar, position298)
										}
									l296:
										{
											position297, tokenIndex297, depth297 := position, tokenIndex, depth
											{
												position300 := position
												depth++
												{
													position301, tokenIndex301, depth301 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l301
													}
													position++
													goto l297
												l301:
													position, tokenIndex, depth = position301, tokenIndex301, depth301
												}
												if !matchDot() {
													goto l297
												}
												depth--
												add(ruleregexChar, position300)
											}
											goto l296
										l297:
											position, tokenIndex, depth = position297, tokenIndex297, depth297
										}
										if buffer[position] != rune('/') {
											goto l272
										}
										position++
									l302:
										{
											position303, tokenIndex303, depth303 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l303
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l303
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l303
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l303
													}
													position++
													break
												}
											}

											goto l302
										l303:
											position, tokenIndex, depth = position303, tokenIndex303, depth303
										}
										depth--
										add(ruleregexBody, position295)
									}
									depth--
									add(rulePegText, position294)
								}
								{
									add(ruleAction38, position)
								}
								depth--
								add(ruleRegex, position293)
							}
							break
						case 'T':
							{
								position306 := position
								depth++
								{
									position307, tokenIndex307, depth307 := position, tokenIndex, depth
									{
										position309 := position
										depth++
										if buffer[position] != rune('T') {
											goto l308
										}
										position++
										if buffer[position] != rune('i') {
											goto l308
										}
										position++
										if buffer[position] != rune('m') {
											goto l308
										}
										position++
										if buffer[position] != rune('e') {
											goto l308
										}
										position++
										if buffer[position] != rune('s') {
											goto l308
										}
										position++
										if buffer[position] != rune('t') {
											goto l308
										}
										position++
										if buffer[position] != rune('a') {
											goto l308
										}
										position++
										if buffer[position] != rune('m') {
											goto l308
										}
										position++
										if buffer[position] != rune('p') {
											goto l308
										}
										position++
										if buffer[position] != rune('(') {
											goto l308
										}
										position++
										{
											position310 := position
											depth++
											{
												position313, tokenIndex313, depth313 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l313
												}
												position++
												goto l308
											l313:
												position, tokenIndex, depth = position313, tokenIndex313, depth313
											}
											if !matchDot() {
												goto l308
											}
										l311:
											{
												position312, tokenIndex312, depth312 := position, tokenIndex, depth
												{
													position314, tokenIndex314, depth314 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l314
													}
													position++
													goto l312
												l314:
													position, tokenIndex, depth = position314, tokenIndex314, depth314
												}
												if !matchDot() {
													goto l312
												}
												goto l311
											l312:
												position, tokenIndex, depth = position312, tokenIndex312, depth312
											}
											depth--
											add(rulePegText, position310)
										}
										if buffer[position] != rune(')') {
											goto l308
										}
										position++
										{
											add(ruleAction39, position)
										}
										depth--
										add(ruletimestampParen, position309)
									}
									goto l307
								l308:
									position, tokenIndex, depth = position307, tokenIndex307, depth307
									{
										position316 := position
										depth++
										if buffer[position] != rune('T') {
											goto l272
										}
										position++
										if buffer[position] != rune('i') {
											goto l272
										}
										position++
										if buffer[position] != rune('m') {
											goto l272
										}
										position++
										if buffer[position] != rune('e') {
											goto l272
										}
										position++
										if buffer[position] != rune('s') {
											goto l272
										}
										position++
										if buffer[position] != rune('t') {
											goto l272
										}
										position++
										if buffer[position] != rune('a') {
											goto l272
										}
										position++
										if buffer[position] != rune('m') {
											goto l272
										}
										position++
										if buffer[position] != rune('p') {
											goto l272
										}
										position++
										if buffer[position] != rune(' ') {
											goto l272
										}
										position++
										{
											position317 := position
											depth++
											{
												position320, tokenIndex320, depth320 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l321
												}
												position++
												goto l320
											l321:
												position, tokenIndex, depth = position320, tokenIndex320, depth320
												if buffer[position] != rune('|') {
													goto l272
												}
												position++
											}
										l320:
										l318:
											{
												position319, tokenIndex319, depth319 := position, tokenIndex, depth
												{
													position322, tokenIndex322, depth322 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l323
													}
													position++
													goto l322
												l323:
													position, tokenIndex, depth = position322, tokenIndex322, depth322
													if buffer[position] != rune('|') {
														goto l319
													}
													position++
												}
											l322:
												goto l318
											l319:
												position, tokenIndex, depth = position319, tokenIndex319, depth319
											}
											depth--
											add(rulePegText, position317)
										}
										{
											add(ruleAction40, position)
										}
										depth--
										add(ruletimestampPipe, position316)
									}
								}
							l307:
								depth--
								add(ruleTimestampVal, position306)
							}
							break
						case 'B':
							{
								position325 := position
								depth++
								if buffer[position] != rune('B') {
									goto l272
								}
								position++
								if buffer[position] != rune('i') {
									goto l272
								}
								position++
								if buffer[position] != rune('n') {
									goto l272
								}
								position++
								if buffer[position] != rune('D') {
									goto l272
								}
								position++
								if buffer[position] != rune('a') {
									goto l272
								}
								position++
								if buffer[position] != rune('t') {
									goto l272
								}
								position++
								if buffer[position] != rune('a') {
									goto l272
								}
								position++
								if buffer[position] != rune('(') {
									goto l272
								}
								position++
								{
									position326 := position
									depth++
									{
										position329, tokenIndex329, depth329 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l329
										}
										position++
										goto l272
									l329:
										position, tokenIndex, depth = position329, tokenIndex329, depth329
									}
									if !matchDot() {
										goto l272
									}
								l327:
									{
										position328, tokenIndex328, depth328 := position, tokenIndex, depth
										{
											position330, tokenIndex330, depth330 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l330
											}
											position++
											goto l328
										l330:
											position, tokenIndex, depth = position330, tokenIndex330, depth330
										}
										if !matchDot() {
											goto l328
										}
										goto l327
									l328:
										position, tokenIndex, depth = position328, tokenIndex328, depth328
									}
									depth--
									add(rulePegText, position326)
								}
								if buffer[position] != rune(')') {
									goto l272
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleBinData, position325)
							}
							break
						case 'D', 'n':
							{
								position332 := position
								depth++
								{
									position333, tokenIndex333, depth333 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l333
									}
									position++
									if buffer[position] != rune('e') {
										goto l333
									}
									position++
									if buffer[position] != rune('w') {
										goto l333
									}
									position++
									if buffer[position] != rune(' ') {
										goto l333
									}
									position++
									goto l334
								l333:
									position, tokenIndex, depth = position333, tokenIndex333, depth333
								}
							l334:
								if buffer[position] != rune('D') {
									goto l272
								}
								position++
								if buffer[position] != rune('a') {
									goto l272
								}
								position++
								if buffer[position] != rune('t') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('(') {
									goto l272
								}
								position++
								{
									position335, tokenIndex335, depth335 := position, tokenIndex, depth
									if buffer[position] != rune('-') {
										goto l335
									}
									position++
									goto l336
								l335:
									position, tokenIndex, depth = position335, tokenIndex335, depth335
								}
							l336:
								{
									position337 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l272
									}
									position++
								l338:
									{
										position339, tokenIndex339, depth339 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l339
										}
										position++
										goto l338
									l339:
										position, tokenIndex, depth = position339, tokenIndex339, depth339
									}
									depth--
									add(rulePegText, position337)
								}
								if buffer[position] != rune(')') {
									goto l272
								}
								position++
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleDate, position332)
							}
							break
						case 'O':
							{
								position341 := position
								depth++
								if buffer[position] != rune('O') {
									goto l272
								}
								position++
								if buffer[position] != rune('b') {
									goto l272
								}
								position++
								if buffer[position] != rune('j') {
									goto l272
								}
								position++
								if buffer[position] != rune('e') {
									goto l272
								}
								position++
								if buffer[position] != rune('c') {
									goto l272
								}
								position++
								if buffer[position] != rune('t') {
									goto l272
								}
								position++
								if buffer[position] != rune('I') {
									goto l272
								}
								position++
								if buffer[position] != rune('d') {
									goto l272
								}
								position++
								if buffer[position] != rune('(') {
									goto l272
								}
								position++
								{
									position342, tokenIndex342, depth342 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l343
									}
									position++
									goto l342
								l343:
									position, tokenIndex, depth = position342, tokenIndex342, depth342
									if buffer[position] != rune('"') {
										goto l272
									}
									position++
								}
							l342:
								{
									position344 := position
									depth++
								l345:
									{
										position346, tokenIndex346, depth346 := position, tokenIndex, depth
										{
											position347 := position
											depth++
											{
												position348, tokenIndex348, depth348 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l349
												}
												position++
												goto l348
											l349:
												position, tokenIndex, depth = position348, tokenIndex348, depth348
												{
													position350, tokenIndex350, depth350 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l351
													}
													position++
													goto l350
												l351:
													position, tokenIndex, depth = position350, tokenIndex350, depth350
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l346
													}
													position++
												}
											l350:
											}
										l348:
											depth--
											add(rulehexChar, position347)
										}
										goto l345
									l346:
										position, tokenIndex, depth = position346, tokenIndex346, depth346
									}
									depth--
									add(rulePegText, position344)
								}
								{
									position352, tokenIndex352, depth352 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l353
									}
									position++
									goto l352
								l353:
									position, tokenIndex, depth = position352, tokenIndex352, depth352
									if buffer[position] != rune('"') {
										goto l272
									}
									position++
								}
							l352:
								if buffer[position] != rune(')') {
									goto l272
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleObjectID, position341)
							}
							break
						case '"':
							{
								position355 := position
								depth++
								if buffer[position] != rune('"') {
									goto l272
								}
								position++
								{
									position356 := position
									depth++
								l357:
									{
										position358, tokenIndex358, depth358 := position, tokenIndex, depth
										{
											position359 := position
											depth++
											{
												position360, tokenIndex360, depth360 := position, tokenIndex, depth
												{
													position362, tokenIndex362, depth362 := position, tokenIndex, depth
													{
														position363, tokenIndex363, depth363 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l364
														}
														position++
														goto l363
													l364:
														position, tokenIndex, depth = position363, tokenIndex363, depth363
														if buffer[position] != rune('\\') {
															goto l362
														}
														position++
													}
												l363:
													goto l361
												l362:
													position, tokenIndex, depth = position362, tokenIndex362, depth362
												}
												if !matchDot() {
													goto l361
												}
												goto l360
											l361:
												position, tokenIndex, depth = position360, tokenIndex360, depth360
												if buffer[position] != rune('\\') {
													goto l358
												}
												position++
												if !matchDot() {
													goto l358
												}
											}
										l360:
											depth--
											add(rulestringChar, position359)
										}
										goto l357
									l358:
										position, tokenIndex, depth = position358, tokenIndex358, depth358
									}
									depth--
									add(rulePegText, position356)
								}
								if buffer[position] != rune('"') {
									goto l272
								}
								position++
								{
									add(ruleAction31, position)
								}
								depth--
								add(ruleString, position355)
							}
							break
						case 'f', 't':
							{
								position366 := position
								depth++
								{
									position367, tokenIndex367, depth367 := position, tokenIndex, depth
									{
										position369 := position
										depth++
										if buffer[position] != rune('t') {
											goto l368
										}
										position++
										if buffer[position] != rune('r') {
											goto l368
										}
										position++
										if buffer[position] != rune('u') {
											goto l368
										}
										position++
										if buffer[position] != rune('e') {
											goto l368
										}
										position++
										{
											add(ruleAction33, position)
										}
										depth--
										add(ruleTrue, position369)
									}
									goto l367
								l368:
									position, tokenIndex, depth = position367, tokenIndex367, depth367
									{
										position371 := position
										depth++
										if buffer[position] != rune('f') {
											goto l272
										}
										position++
										if buffer[position] != rune('a') {
											goto l272
										}
										position++
										if buffer[position] != rune('l') {
											goto l272
										}
										position++
										if buffer[position] != rune('s') {
											goto l272
										}
										position++
										if buffer[position] != rune('e') {
											goto l272
										}
										position++
										{
											add(ruleAction34, position)
										}
										depth--
										add(ruleFalse, position371)
									}
								}
							l367:
								depth--
								add(ruleBoolean, position366)
							}
							break
						case '[':
							{
								position373 := position
								depth++
								if buffer[position] != rune('[') {
									goto l272
								}
								position++
								{
									add(ruleAction26, position)
								}
								{
									position375, tokenIndex375, depth375 := position, tokenIndex, depth
									{
										position377 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l375
										}
									l378:
										{
											position379, tokenIndex379, depth379 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l379
											}
											position++
											if !_rules[ruleListElem]() {
												goto l379
											}
											goto l378
										l379:
											position, tokenIndex, depth = position379, tokenIndex379, depth379
										}
										depth--
										add(ruleListElements, position377)
									}
									goto l376
								l375:
									position, tokenIndex, depth = position375, tokenIndex375, depth375
								}
							l376:
								if buffer[position] != rune(']') {
									goto l272
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleList, position373)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l272
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l272
							}
							break
						}
					}

				}
			l274:
				depth--
				add(ruleValue, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 49 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action30)> */
		func() bool {
			position381, tokenIndex381, depth381 := position, tokenIndex, depth
			{
				position382 := position
				depth++
				{
					position383 := position
					depth++
					{
						position384, tokenIndex384, depth384 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l384
						}
						position++
						goto l385
					l384:
						position, tokenIndex, depth = position384, tokenIndex384, depth384
					}
				l385:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l381
					}
					position++
				l386:
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l387
						}
						position++
						goto l386
					l387:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
					}
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l388
						}
						position++
						goto l389
					l388:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
					}
				l389:
				l390:
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l391
						}
						position++
						goto l390
					l391:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
					}
					depth--
					add(rulePegText, position383)
				}
				{
					add(ruleAction30, position)
				}
				depth--
				add(ruleNumeric, position382)
			}
			return true
		l381:
			position, tokenIndex, depth = position381, tokenIndex381, depth381
			return false
		},
		/* 50 Boolean <- <(True / False)> */
		nil,
		/* 51 String <- <('"' <stringChar*> '"' Action31)> */
		nil,
		/* 52 Null <- <('n' 'u' 'l' 'l' Action32)> */
		nil,
		/* 53 True <- <('t' 'r' 'u' 'e' Action33)> */
		nil,
		/* 54 False <- <('f' 'a' 'l' 's' 'e' Action34)> */
		nil,
		/* 55 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') '-'? <[0-9]+> ')' Action35)> */
		nil,
		/* 56 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' ('\'' / '"') <hexChar*> ('\'' / '"') ')' Action36)> */
		nil,
		/* 57 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action37)> */
		nil,
		/* 58 Regex <- <('/' <regexBody> Action38)> */
		nil,
		/* 59 TimestampVal <- <(timestampParen / timestampPipe)> */
		nil,
		/* 60 timestampParen <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action39)> */
		nil,
		/* 61 timestampPipe <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' ' ' <([0-9] / '|')+> Action40)> */
		nil,
		/* 62 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action41)> */
		nil,
		/* 63 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action42)> */
		nil,
		/* 64 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action43)> */
		nil,
		/* 65 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action44)> */
		nil,
		/* 66 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 67 regexChar <- <(!'/' .)> */
		nil,
		/* 68 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 69 stringChar <- <((!('"' / '\\') .) / ('\\' .))> */
		nil,
		/* 70 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l413
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l413
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l413
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l413
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l413
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l413
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l413
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		nil,
		/* 73 Action0 <- <{ p.SetField("severity", buffer[begin:end]) }> */
		nil,
		/* 74 Action1 <- <{ p.SetField("component", buffer[begin:end]) }> */
		nil,
		/* 75 Action2 <- <{ p.SetField("context", buffer[begin:end]) }> */
		nil,
		/* 76 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 77 Action4 <- <{ p.SetField("warning", buffer[begin:end]) }> */
		nil,
		/* 78 Action5 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 79 Action6 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 80 Action7 <- <{ p.EndField() }> */
		nil,
		/* 81 Action8 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 82 Action9 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 83 Action10 <- <{ p.EndField() }> */
		nil,
		/* 84 Action11 <- <{ p.SetField("command_type", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 85 Action12 <- <{ p.EndField() }> */
		nil,
		/* 86 Action13 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 87 Action14 <- <{ p.EndField()}> */
		nil,
		/* 88 Action15 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 89 Action16 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 90 Action17 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 91 Action18 <- <{ p.StartField("exception") }> */
		nil,
		/* 92 Action19 <- <{ p.PushValue(buffer[begin:end]); p.EndField() }> */
		nil,
		/* 93 Action20 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 94 Action21 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 95 Action22 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 96 Action23 <- <{ p.PushMap() }> */
		nil,
		/* 97 Action24 <- <{ p.PopMap() }> */
		nil,
		/* 98 Action25 <- <{ p.SetMapValue() }> */
		nil,
		/* 99 Action26 <- <{ p.PushList() }> */
		nil,
		/* 100 Action27 <- <{ p.PopList() }> */
		nil,
		/* 101 Action28 <- <{ p.SetListValue() }> */
		nil,
		/* 102 Action29 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 103 Action30 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 104 Action31 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 105 Action32 <- <{ p.PushValue(nil) }> */
		nil,
		/* 106 Action33 <- <{ p.PushValue(true) }> */
		nil,
		/* 107 Action34 <- <{ p.PushValue(false) }> */
		nil,
		/* 108 Action35 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 109 Action36 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 110 Action37 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 111 Action38 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 112 Action39 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 113 Action40 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 114 Action41 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 115 Action42 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 116 Action43 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 117 Action44 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
