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
	ruleLogLevel
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
	"LogLevel",
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
			p.SetField("log_level", buffer[begin:end])
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
		/* 0 MongoLogLine <- <(Timestamp LogLevel? Component? Context Warning? Op NS LineField* Locks? LineField* Duration? S? LineField* extra? !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2 := position
					depth++
					{
						position3, tokenIndex3, depth3 := position, tokenIndex, depth
						{
							position5 := position
							depth++
							{
								position6 := position
								depth++
								{
									position7 := position
									depth++
									{
										position8 := position
										depth++
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l4
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l4
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l4
										}
										position++
										depth--
										add(ruleday, position8)
									}
									if buffer[position] != rune(' ') {
										goto l4
									}
									position++
									{
										position9 := position
										depth++
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l4
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l4
										}
										position++
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l4
										}
										position++
										depth--
										add(rulemonth, position9)
									}
									if buffer[position] != rune(' ') {
										goto l4
									}
									position++
								l10:
									{
										position11, tokenIndex11, depth11 := position, tokenIndex, depth
										if buffer[position] != rune(' ') {
											goto l11
										}
										position++
										goto l10
									l11:
										position, tokenIndex, depth = position11, tokenIndex11, depth11
									}
									{
										position12 := position
										depth++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l4
										}
										position++
										{
											position13, tokenIndex13, depth13 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l13
											}
											position++
											goto l14
										l13:
											position, tokenIndex, depth = position13, tokenIndex13, depth13
										}
									l14:
										depth--
										add(ruledayNum, position12)
									}
									depth--
									add(ruledate, position7)
								}
								if buffer[position] != rune(' ') {
									goto l4
								}
								position++
								if !_rules[ruletime]() {
									goto l4
								}
								depth--
								add(rulePegText, position6)
							}
							{
								add(ruleAction20, position)
							}
							depth--
							add(ruletimestamp24, position5)
						}
						goto l3
					l4:
						position, tokenIndex, depth = position3, tokenIndex3, depth3
						{
							position16 := position
							depth++
							{
								position17 := position
								depth++
								{
									position18 := position
									depth++
									{
										position19 := position
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
										add(ruledigit4, position19)
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
										position20, tokenIndex20, depth20 := position, tokenIndex, depth
										{
											position22 := position
											depth++
											if buffer[position] != rune('+') {
												goto l20
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l20
											}
											position++
										l23:
											{
												position24, tokenIndex24, depth24 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l24
												}
												position++
												goto l23
											l24:
												position, tokenIndex, depth = position24, tokenIndex24, depth24
											}
											depth--
											add(ruletz, position22)
										}
										goto l21
									l20:
										position, tokenIndex, depth = position20, tokenIndex20, depth20
									}
								l21:
									depth--
									add(ruledatetime26, position18)
								}
								depth--
								add(rulePegText, position17)
							}
							{
								add(ruleAction21, position)
							}
							depth--
							add(ruletimestamp26, position16)
						}
					}
				l3:
					{
						position26, tokenIndex26, depth26 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l26
						}
						goto l27
					l26:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
					}
				l27:
					depth--
					add(ruleTimestamp, position2)
				}
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					{
						position30 := position
						depth++
						{
							position31 := position
							depth++
							{
								position32, tokenIndex32, depth32 := position, tokenIndex, depth
								if buffer[position] != rune('I') {
									goto l33
								}
								position++
								goto l32
							l33:
								position, tokenIndex, depth = position32, tokenIndex32, depth32
								if buffer[position] != rune('D') {
									goto l28
								}
								position++
							}
						l32:
							depth--
							add(rulePegText, position31)
						}
						if buffer[position] != rune(' ') {
							goto l28
						}
						position++
						{
							add(ruleAction0, position)
						}
						depth--
						add(ruleLogLevel, position30)
					}
					goto l29
				l28:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
				}
			l29:
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					{
						position37 := position
						depth++
						{
							position38 := position
							depth++
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l35
							}
							position++
						l39:
							{
								position40, tokenIndex40, depth40 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l40
								}
								position++
								goto l39
							l40:
								position, tokenIndex, depth = position40, tokenIndex40, depth40
							}
							depth--
							add(rulePegText, position38)
						}
						if buffer[position] != rune(' ') {
							goto l35
						}
						position++
					l41:
						{
							position42, tokenIndex42, depth42 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l42
							}
							position++
							goto l41
						l42:
							position, tokenIndex, depth = position42, tokenIndex42, depth42
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleComponent, position37)
					}
					goto l36
				l35:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
				}
			l36:
				{
					position44 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position45 := position
						depth++
						{
							position48 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position50, tokenIndex50, depth50 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l51
										}
										position++
										goto l50
									l51:
										position, tokenIndex, depth = position50, tokenIndex50, depth50
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l50:
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
							add(ruleletterOrDigit, position48)
						}
					l46:
						{
							position47, tokenIndex47, depth47 := position, tokenIndex, depth
							{
								position52 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position54, tokenIndex54, depth54 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l55
											}
											position++
											goto l54
										l55:
											position, tokenIndex, depth = position54, tokenIndex54, depth54
											if buffer[position] != rune('$') {
												goto l47
											}
											position++
										}
									l54:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l47
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l47
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l47
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position52)
							}
							goto l46
						l47:
							position, tokenIndex, depth = position47, tokenIndex47, depth47
						}
						depth--
						add(rulePegText, position45)
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
					add(ruleContext, position44)
				}
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					{
						position59 := position
						depth++
						{
							position60 := position
							depth++
							{
								position61 := position
								depth++
								if buffer[position] != rune('w') {
									goto l57
								}
								position++
								if buffer[position] != rune('a') {
									goto l57
								}
								position++
								if buffer[position] != rune('r') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('g') {
									goto l57
								}
								position++
								if buffer[position] != rune(':') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('l') {
									goto l57
								}
								position++
								if buffer[position] != rune('o') {
									goto l57
								}
								position++
								if buffer[position] != rune('g') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('l') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('a') {
									goto l57
								}
								position++
								if buffer[position] != rune('t') {
									goto l57
								}
								position++
								if buffer[position] != rune('t') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune('m') {
									goto l57
								}
								position++
								if buffer[position] != rune('p') {
									goto l57
								}
								position++
								if buffer[position] != rune('t') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune('d') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('(') {
									goto l57
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l57
								}
								position++
							l62:
								{
									position63, tokenIndex63, depth63 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l63
									}
									position++
									goto l62
								l63:
									position, tokenIndex, depth = position63, tokenIndex63, depth63
								}
								if buffer[position] != rune('k') {
									goto l57
								}
								position++
								if buffer[position] != rune(')') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('o') {
									goto l57
								}
								position++
								if buffer[position] != rune('v') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune('r') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('m') {
									goto l57
								}
								position++
								if buffer[position] != rune('a') {
									goto l57
								}
								position++
								if buffer[position] != rune('x') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('s') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('z') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('(') {
									goto l57
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l57
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
									goto l57
								}
								position++
								if buffer[position] != rune(')') {
									goto l57
								}
								position++
								if buffer[position] != rune(',') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('p') {
									goto l57
								}
								position++
								if buffer[position] != rune('r') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('t') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('g') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('b') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune('g') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('i') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('g') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('a') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('d') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('e') {
									goto l57
								}
								position++
								if buffer[position] != rune('n') {
									goto l57
								}
								position++
								if buffer[position] != rune('d') {
									goto l57
								}
								position++
								if buffer[position] != rune(' ') {
									goto l57
								}
								position++
								if buffer[position] != rune('.') {
									goto l57
								}
								position++
								if buffer[position] != rune('.') {
									goto l57
								}
								position++
								if buffer[position] != rune('.') {
									goto l57
								}
								position++
								depth--
								add(ruleloglineSizeWarning, position61)
							}
							depth--
							add(rulePegText, position60)
						}
						if buffer[position] != rune(' ') {
							goto l57
						}
						position++
						{
							add(ruleAction4, position)
						}
						depth--
						add(ruleWarning, position59)
					}
					goto l58
				l57:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
				}
			l58:
				{
					position67 := position
					depth++
					{
						position68 := position
						depth++
						{
							position71, tokenIndex71, depth71 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l72
							}
							position++
							goto l71
						l72:
							position, tokenIndex, depth = position71, tokenIndex71, depth71
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l0
							}
							position++
						}
					l71:
					l69:
						{
							position70, tokenIndex70, depth70 := position, tokenIndex, depth
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
									goto l70
								}
								position++
							}
						l73:
							goto l69
						l70:
							position, tokenIndex, depth = position70, tokenIndex70, depth70
						}
						depth--
						add(rulePegText, position68)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleOp, position67)
				}
				{
					position76 := position
					depth++
					{
						position77 := position
						depth++
					l78:
						{
							position79, tokenIndex79, depth79 := position, tokenIndex, depth
							{
								position80 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l79
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l79
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l79
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l79
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l79
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l79
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position80)
							}
							goto l78
						l79:
							position, tokenIndex, depth = position79, tokenIndex79, depth79
						}
						depth--
						add(rulePegText, position77)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction5, position)
					}
					depth--
					add(ruleNS, position76)
				}
			l83:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l84
					}
					goto l83
				l84:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
				}
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					{
						position87 := position
						depth++
						if buffer[position] != rune('l') {
							goto l85
						}
						position++
						if buffer[position] != rune('o') {
							goto l85
						}
						position++
						if buffer[position] != rune('c') {
							goto l85
						}
						position++
						if buffer[position] != rune('k') {
							goto l85
						}
						position++
						if buffer[position] != rune('s') {
							goto l85
						}
						position++
						if buffer[position] != rune('(') {
							goto l85
						}
						position++
						if buffer[position] != rune('m') {
							goto l85
						}
						position++
						if buffer[position] != rune('i') {
							goto l85
						}
						position++
						if buffer[position] != rune('c') {
							goto l85
						}
						position++
						if buffer[position] != rune('r') {
							goto l85
						}
						position++
						if buffer[position] != rune('o') {
							goto l85
						}
						position++
						if buffer[position] != rune('s') {
							goto l85
						}
						position++
						if buffer[position] != rune(')') {
							goto l85
						}
						position++
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l88
							}
							goto l89
						l88:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
						}
					l89:
					l90:
						{
							position91, tokenIndex91, depth91 := position, tokenIndex, depth
							{
								position92 := position
								depth++
								{
									position93 := position
									depth++
									{
										switch buffer[position] {
										case 'R':
											if buffer[position] != rune('R') {
												goto l91
											}
											position++
											break
										case 'r':
											if buffer[position] != rune('r') {
												goto l91
											}
											position++
											break
										default:
											{
												position95, tokenIndex95, depth95 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l96
												}
												position++
												goto l95
											l96:
												position, tokenIndex, depth = position95, tokenIndex95, depth95
												if buffer[position] != rune('W') {
													goto l91
												}
												position++
											}
										l95:
											break
										}
									}

									depth--
									add(rulePegText, position93)
								}
								{
									add(ruleAction6, position)
								}
								if buffer[position] != rune(':') {
									goto l91
								}
								position++
								if !_rules[ruleNumeric]() {
									goto l91
								}
								{
									position98, tokenIndex98, depth98 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l98
									}
									goto l99
								l98:
									position, tokenIndex, depth = position98, tokenIndex98, depth98
								}
							l99:
								{
									add(ruleAction7, position)
								}
								depth--
								add(rulelock, position92)
							}
							goto l90
						l91:
							position, tokenIndex, depth = position91, tokenIndex91, depth91
						}
						depth--
						add(ruleLocks, position87)
					}
					goto l86
				l85:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
				}
			l86:
			l101:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l102
					}
					goto l101
				l102:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
				}
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					{
						position105 := position
						depth++
						{
							position106 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l103
							}
							position++
						l107:
							{
								position108, tokenIndex108, depth108 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l108
								}
								position++
								goto l107
							l108:
								position, tokenIndex, depth = position108, tokenIndex108, depth108
							}
							depth--
							add(rulePegText, position106)
						}
						if buffer[position] != rune('m') {
							goto l103
						}
						position++
						if buffer[position] != rune('s') {
							goto l103
						}
						position++
						{
							add(ruleAction8, position)
						}
						depth--
						add(ruleDuration, position105)
					}
					goto l104
				l103:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
				}
			l104:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l110
					}
					goto l111
				l110:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
				}
			l111:
			l112:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l113
					}
					goto l112
				l113:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
				}
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					{
						position116 := position
						depth++
						{
							position117 := position
							depth++
							if !matchDot() {
								goto l114
							}
						l118:
							{
								position119, tokenIndex119, depth119 := position, tokenIndex, depth
								if !matchDot() {
									goto l119
								}
								goto l118
							l119:
								position, tokenIndex, depth = position119, tokenIndex119, depth119
							}
							depth--
							add(rulePegText, position117)
						}
						{
							add(ruleAction22, position)
						}
						depth--
						add(ruleextra, position116)
					}
					goto l115
				l114:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
				}
			l115:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if !matchDot() {
						goto l121
					}
					goto l0
				l121:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
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
		/* 2 LogLevel <- <(<('I' / 'D')> ' ' Action0)> */
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
			position129, tokenIndex129, depth129 := position, tokenIndex, depth
			{
				position130 := position
				depth++
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					{
						position133 := position
						depth++
						if buffer[position] != rune('e') {
							goto l132
						}
						position++
						if buffer[position] != rune('x') {
							goto l132
						}
						position++
						if buffer[position] != rune('c') {
							goto l132
						}
						position++
						if buffer[position] != rune('e') {
							goto l132
						}
						position++
						if buffer[position] != rune('p') {
							goto l132
						}
						position++
						if buffer[position] != rune('t') {
							goto l132
						}
						position++
						if buffer[position] != rune('i') {
							goto l132
						}
						position++
						if buffer[position] != rune('o') {
							goto l132
						}
						position++
						if buffer[position] != rune('n') {
							goto l132
						}
						position++
						if buffer[position] != rune(':') {
							goto l132
						}
						position++
						{
							add(ruleAction18, position)
						}
						{
							position135 := position
							depth++
							{
								position138, tokenIndex138, depth138 := position, tokenIndex, depth
								if !matchDot() {
									goto l132
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
									goto l132
								l139:
									position, tokenIndex, depth = position139, tokenIndex139, depth139
								}
								position, tokenIndex, depth = position138, tokenIndex138, depth138
							}
							if !matchDot() {
								goto l132
							}
						l136:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								{
									position140, tokenIndex140, depth140 := position, tokenIndex, depth
									if !matchDot() {
										goto l137
									}
									{
										position141, tokenIndex141, depth141 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l141
										}
										position++
										if buffer[position] != rune('o') {
											goto l141
										}
										position++
										if buffer[position] != rune('d') {
											goto l141
										}
										position++
										if buffer[position] != rune('e') {
											goto l141
										}
										position++
										if buffer[position] != rune(':') {
											goto l141
										}
										position++
										goto l137
									l141:
										position, tokenIndex, depth = position141, tokenIndex141, depth141
									}
									position, tokenIndex, depth = position140, tokenIndex140, depth140
								}
								if !matchDot() {
									goto l137
								}
								goto l136
							l137:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
							}
							depth--
							add(rulePegText, position135)
						}
						{
							position142, tokenIndex142, depth142 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l142
							}
							goto l143
						l142:
							position, tokenIndex, depth = position142, tokenIndex142, depth142
						}
					l143:
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleexceptionField, position133)
					}
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position146 := position
						depth++
						if buffer[position] != rune('c') {
							goto l145
						}
						position++
						if buffer[position] != rune('o') {
							goto l145
						}
						position++
						if buffer[position] != rune('m') {
							goto l145
						}
						position++
						if buffer[position] != rune('m') {
							goto l145
						}
						position++
						if buffer[position] != rune('a') {
							goto l145
						}
						position++
						if buffer[position] != rune('n') {
							goto l145
						}
						position++
						if buffer[position] != rune('d') {
							goto l145
						}
						position++
						if buffer[position] != rune(':') {
							goto l145
						}
						position++
						if buffer[position] != rune(' ') {
							goto l145
						}
						position++
						{
							position147 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l145
							}
						l148:
							{
								position149, tokenIndex149, depth149 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l149
								}
								goto l148
							l149:
								position, tokenIndex, depth = position149, tokenIndex149, depth149
							}
							depth--
							add(rulePegText, position147)
						}
						{
							position150, tokenIndex150, depth150 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l150
							}
							goto l151
						l150:
							position, tokenIndex, depth = position150, tokenIndex150, depth150
						}
					l151:
						{
							add(ruleAction11, position)
						}
						if !_rules[ruleLineValue]() {
							goto l145
						}
						{
							add(ruleAction12, position)
						}
						depth--
						add(rulecommandField, position146)
					}
					goto l131
				l145:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position155 := position
						depth++
						if buffer[position] != rune('p') {
							goto l154
						}
						position++
						if buffer[position] != rune('l') {
							goto l154
						}
						position++
						if buffer[position] != rune('a') {
							goto l154
						}
						position++
						if buffer[position] != rune('n') {
							goto l154
						}
						position++
						if buffer[position] != rune('S') {
							goto l154
						}
						position++
						if buffer[position] != rune('u') {
							goto l154
						}
						position++
						if buffer[position] != rune('m') {
							goto l154
						}
						position++
						if buffer[position] != rune('m') {
							goto l154
						}
						position++
						if buffer[position] != rune('a') {
							goto l154
						}
						position++
						if buffer[position] != rune('r') {
							goto l154
						}
						position++
						if buffer[position] != rune('y') {
							goto l154
						}
						position++
						if buffer[position] != rune(':') {
							goto l154
						}
						position++
						if buffer[position] != rune(' ') {
							goto l154
						}
						position++
						{
							add(ruleAction13, position)
						}
						{
							position157 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l154
							}
						l158:
							{
								position159, tokenIndex159, depth159 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l159
								}
								position++
								if buffer[position] != rune(' ') {
									goto l159
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l159
								}
								goto l158
							l159:
								position, tokenIndex, depth = position159, tokenIndex159, depth159
							}
							depth--
							add(ruleplanSummaryElements, position157)
						}
						{
							add(ruleAction14, position)
						}
						depth--
						add(ruleplanSummaryField, position155)
					}
					goto l131
				l154:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position161 := position
						depth++
						{
							position162 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l129
							}
						l163:
							{
								position164, tokenIndex164, depth164 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l164
								}
								goto l163
							l164:
								position, tokenIndex, depth = position164, tokenIndex164, depth164
							}
							depth--
							add(rulePegText, position162)
						}
						if buffer[position] != rune(':') {
							goto l129
						}
						position++
						{
							position165, tokenIndex165, depth165 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l165
							}
							goto l166
						l165:
							position, tokenIndex, depth = position165, tokenIndex165, depth165
						}
					l166:
						{
							add(ruleAction9, position)
						}
						if !_rules[ruleLineValue]() {
							goto l129
						}
						{
							add(ruleAction10, position)
						}
						depth--
						add(ruleplainField, position161)
					}
				}
			l131:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l169
					}
					goto l170
				l169:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
				}
			l170:
				depth--
				add(ruleLineField, position130)
			}
			return true
		l129:
			position, tokenIndex, depth = position129, tokenIndex129, depth129
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
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				{
					position181 := position
					depth++
					{
						position182 := position
						depth++
						{
							position183, tokenIndex183, depth183 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l184
							}
							position++
							if buffer[position] != rune('N') {
								goto l184
							}
							position++
							if buffer[position] != rune('D') {
								goto l184
							}
							position++
							if buffer[position] != rune('_') {
								goto l184
							}
							position++
							if buffer[position] != rune('H') {
								goto l184
							}
							position++
							if buffer[position] != rune('A') {
								goto l184
							}
							position++
							if buffer[position] != rune('S') {
								goto l184
							}
							position++
							if buffer[position] != rune('H') {
								goto l184
							}
							position++
							goto l183
						l184:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('C') {
								goto l185
							}
							position++
							if buffer[position] != rune('A') {
								goto l185
							}
							position++
							if buffer[position] != rune('C') {
								goto l185
							}
							position++
							if buffer[position] != rune('H') {
								goto l185
							}
							position++
							if buffer[position] != rune('E') {
								goto l185
							}
							position++
							if buffer[position] != rune('D') {
								goto l185
							}
							position++
							if buffer[position] != rune('_') {
								goto l185
							}
							position++
							if buffer[position] != rune('P') {
								goto l185
							}
							position++
							if buffer[position] != rune('L') {
								goto l185
							}
							position++
							if buffer[position] != rune('A') {
								goto l185
							}
							position++
							if buffer[position] != rune('N') {
								goto l185
							}
							position++
							goto l183
						l185:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('C') {
								goto l186
							}
							position++
							if buffer[position] != rune('O') {
								goto l186
							}
							position++
							if buffer[position] != rune('L') {
								goto l186
							}
							position++
							if buffer[position] != rune('L') {
								goto l186
							}
							position++
							if buffer[position] != rune('S') {
								goto l186
							}
							position++
							if buffer[position] != rune('C') {
								goto l186
							}
							position++
							if buffer[position] != rune('A') {
								goto l186
							}
							position++
							if buffer[position] != rune('N') {
								goto l186
							}
							position++
							goto l183
						l186:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('C') {
								goto l187
							}
							position++
							if buffer[position] != rune('O') {
								goto l187
							}
							position++
							if buffer[position] != rune('U') {
								goto l187
							}
							position++
							if buffer[position] != rune('N') {
								goto l187
							}
							position++
							if buffer[position] != rune('T') {
								goto l187
							}
							position++
							goto l183
						l187:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('D') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							if buffer[position] != rune('L') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							if buffer[position] != rune('T') {
								goto l188
							}
							position++
							if buffer[position] != rune('E') {
								goto l188
							}
							position++
							goto l183
						l188:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('G') {
								goto l189
							}
							position++
							if buffer[position] != rune('E') {
								goto l189
							}
							position++
							if buffer[position] != rune('O') {
								goto l189
							}
							position++
							if buffer[position] != rune('_') {
								goto l189
							}
							position++
							if buffer[position] != rune('N') {
								goto l189
							}
							position++
							if buffer[position] != rune('E') {
								goto l189
							}
							position++
							if buffer[position] != rune('A') {
								goto l189
							}
							position++
							if buffer[position] != rune('R') {
								goto l189
							}
							position++
							if buffer[position] != rune('_') {
								goto l189
							}
							position++
							if buffer[position] != rune('2') {
								goto l189
							}
							position++
							if buffer[position] != rune('D') {
								goto l189
							}
							position++
							goto l183
						l189:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('G') {
								goto l190
							}
							position++
							if buffer[position] != rune('E') {
								goto l190
							}
							position++
							if buffer[position] != rune('O') {
								goto l190
							}
							position++
							if buffer[position] != rune('_') {
								goto l190
							}
							position++
							if buffer[position] != rune('N') {
								goto l190
							}
							position++
							if buffer[position] != rune('E') {
								goto l190
							}
							position++
							if buffer[position] != rune('A') {
								goto l190
							}
							position++
							if buffer[position] != rune('R') {
								goto l190
							}
							position++
							if buffer[position] != rune('_') {
								goto l190
							}
							position++
							if buffer[position] != rune('2') {
								goto l190
							}
							position++
							if buffer[position] != rune('D') {
								goto l190
							}
							position++
							if buffer[position] != rune('S') {
								goto l190
							}
							position++
							if buffer[position] != rune('P') {
								goto l190
							}
							position++
							if buffer[position] != rune('H') {
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
							if buffer[position] != rune('E') {
								goto l190
							}
							position++
							goto l183
						l190:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('I') {
								goto l191
							}
							position++
							if buffer[position] != rune('D') {
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
							if buffer[position] != rune('C') {
								goto l191
							}
							position++
							if buffer[position] != rune('K') {
								goto l191
							}
							position++
							goto l183
						l191:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('S') {
								goto l192
							}
							position++
							if buffer[position] != rune('O') {
								goto l192
							}
							position++
							if buffer[position] != rune('R') {
								goto l192
							}
							position++
							if buffer[position] != rune('T') {
								goto l192
							}
							position++
							if buffer[position] != rune('_') {
								goto l192
							}
							position++
							if buffer[position] != rune('M') {
								goto l192
							}
							position++
							if buffer[position] != rune('E') {
								goto l192
							}
							position++
							if buffer[position] != rune('R') {
								goto l192
							}
							position++
							if buffer[position] != rune('G') {
								goto l192
							}
							position++
							if buffer[position] != rune('E') {
								goto l192
							}
							position++
							goto l183
						l192:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('S') {
								goto l193
							}
							position++
							if buffer[position] != rune('H') {
								goto l193
							}
							position++
							if buffer[position] != rune('A') {
								goto l193
							}
							position++
							if buffer[position] != rune('R') {
								goto l193
							}
							position++
							if buffer[position] != rune('D') {
								goto l193
							}
							position++
							if buffer[position] != rune('I') {
								goto l193
							}
							position++
							if buffer[position] != rune('N') {
								goto l193
							}
							position++
							if buffer[position] != rune('G') {
								goto l193
							}
							position++
							if buffer[position] != rune('_') {
								goto l193
							}
							position++
							if buffer[position] != rune('F') {
								goto l193
							}
							position++
							if buffer[position] != rune('I') {
								goto l193
							}
							position++
							if buffer[position] != rune('L') {
								goto l193
							}
							position++
							if buffer[position] != rune('T') {
								goto l193
							}
							position++
							if buffer[position] != rune('E') {
								goto l193
							}
							position++
							if buffer[position] != rune('R') {
								goto l193
							}
							position++
							goto l183
						l193:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('S') {
								goto l194
							}
							position++
							if buffer[position] != rune('K') {
								goto l194
							}
							position++
							if buffer[position] != rune('I') {
								goto l194
							}
							position++
							if buffer[position] != rune('P') {
								goto l194
							}
							position++
							goto l183
						l194:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('S') {
								goto l195
							}
							position++
							if buffer[position] != rune('O') {
								goto l195
							}
							position++
							if buffer[position] != rune('R') {
								goto l195
							}
							position++
							if buffer[position] != rune('T') {
								goto l195
							}
							position++
							goto l183
						l195:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('X') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('B') {
										goto l179
									}
									position++
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									if buffer[position] != rune('L') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									if buffer[position] != rune('_') {
										goto l179
									}
									position++
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									if buffer[position] != rune('R') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('J') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('R') {
										goto l179
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('L') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('_') {
										goto l179
									}
									position++
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									if buffer[position] != rune('L') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('M') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									if buffer[position] != rune('_') {
										goto l179
									}
									position++
									if buffer[position] != rune('M') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('X') {
										goto l179
									}
									position++
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l179
									}
									position++
									if buffer[position] != rune('R') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('P') {
										goto l179
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('H') {
										goto l179
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('F') {
										goto l179
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('I') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('U') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('_') {
										goto l179
									}
									position++
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									if buffer[position] != rune('C') {
										goto l179
									}
									position++
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l179
									}
									position++
									if buffer[position] != rune('N') {
										goto l179
									}
									position++
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									if buffer[position] != rune('_') {
										goto l179
									}
									position++
									if buffer[position] != rune('S') {
										goto l179
									}
									position++
									if buffer[position] != rune('O') {
										goto l179
									}
									position++
									if buffer[position] != rune('R') {
										goto l179
									}
									position++
									if buffer[position] != rune('T') {
										goto l179
									}
									position++
									if buffer[position] != rune('E') {
										goto l179
									}
									position++
									if buffer[position] != rune('D') {
										goto l179
									}
									position++
									break
								}
							}

						}
					l183:
						depth--
						add(ruleplanSummaryStage, position182)
					}
					depth--
					add(rulePegText, position181)
				}
				{
					add(ruleAction15, position)
				}
				{
					position198 := position
					depth++
					{
						position199, tokenIndex199, depth199 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l200
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l200
						}
						{
							add(ruleAction16, position)
						}
						goto l199
					l200:
						position, tokenIndex, depth = position199, tokenIndex199, depth199
						{
							add(ruleAction17, position)
						}
					}
				l199:
					depth--
					add(ruleplanSummary, position198)
				}
				depth--
				add(ruleplanSummaryElem, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
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
			position206, tokenIndex206, depth206 := position, tokenIndex, depth
			{
				position207 := position
				depth++
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l209
					}
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if !_rules[ruleNumeric]() {
						goto l206
					}
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l210
					}
					goto l211
				l210:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
				}
			l211:
				depth--
				add(ruleLineValue, position207)
			}
			return true
		l206:
			position, tokenIndex, depth = position206, tokenIndex206, depth206
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
			position216, tokenIndex216, depth216 := position, tokenIndex, depth
			{
				position217 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l216
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l216
				}
				position++
				depth--
				add(ruledigit2, position217)
			}
			return true
		l216:
			position, tokenIndex, depth = position216, tokenIndex216, depth216
			return false
		},
		/* 27 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 28 tz <- <('+' [0-9]+)> */
		nil,
		/* 29 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
				{
					position222 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l220
					}
					depth--
					add(rulehour, position222)
				}
				if buffer[position] != rune(':') {
					goto l220
				}
				position++
				{
					position223 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l220
					}
					depth--
					add(ruleminute, position223)
				}
				if buffer[position] != rune(':') {
					goto l220
				}
				position++
				{
					position224 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l220
					}
					depth--
					add(rulesecond, position224)
				}
				if buffer[position] != rune('.') {
					goto l220
				}
				position++
				{
					position225 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l220
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l220
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l220
					}
					position++
					depth--
					add(rulemillisecond, position225)
				}
				depth--
				add(ruletime, position221)
			}
			return true
		l220:
			position, tokenIndex, depth = position220, tokenIndex220, depth220
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
			position236, tokenIndex236, depth236 := position, tokenIndex, depth
			{
				position237 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l236
				}
				position++
			l238:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
				}
				depth--
				add(ruleS, position237)
			}
			return true
		l236:
			position, tokenIndex, depth = position236, tokenIndex236, depth236
			return false
		},
		/* 41 Doc <- <('{' Action23 DocElements? '}' Action24)> */
		func() bool {
			position240, tokenIndex240, depth240 := position, tokenIndex, depth
			{
				position241 := position
				depth++
				if buffer[position] != rune('{') {
					goto l240
				}
				position++
				{
					add(ruleAction23, position)
				}
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					{
						position245 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l243
						}
					l246:
						{
							position247, tokenIndex247, depth247 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l247
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l247
							}
							goto l246
						l247:
							position, tokenIndex, depth = position247, tokenIndex247, depth247
						}
						depth--
						add(ruleDocElements, position245)
					}
					goto l244
				l243:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
				}
			l244:
				if buffer[position] != rune('}') {
					goto l240
				}
				position++
				{
					add(ruleAction24, position)
				}
				depth--
				add(ruleDoc, position241)
			}
			return true
		l240:
			position, tokenIndex, depth = position240, tokenIndex240, depth240
			return false
		},
		/* 42 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 43 DocElem <- <(S? Field S? Value S? Action25)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l252
					}
					goto l253
				l252:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
				}
			l253:
				{
					position254 := position
					depth++
					{
						position255 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l250
						}
					l256:
						{
							position257, tokenIndex257, depth257 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l257
							}
							goto l256
						l257:
							position, tokenIndex, depth = position257, tokenIndex257, depth257
						}
						depth--
						add(rulePegText, position255)
					}
					if buffer[position] != rune(':') {
						goto l250
					}
					position++
					{
						add(ruleAction29, position)
					}
					depth--
					add(ruleField, position254)
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
				if !_rules[ruleValue]() {
					goto l250
				}
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l261
					}
					goto l262
				l261:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
				}
			l262:
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleDocElem, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 44 List <- <('[' Action26 ListElements? ']' Action27)> */
		nil,
		/* 45 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 46 ListElem <- <(S? Value S? Action28)> */
		func() bool {
			position266, tokenIndex266, depth266 := position, tokenIndex, depth
			{
				position267 := position
				depth++
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
				if !_rules[ruleValue]() {
					goto l266
				}
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l270
					}
					goto l271
				l270:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
				}
			l271:
				{
					add(ruleAction28, position)
				}
				depth--
				add(ruleListElem, position267)
			}
			return true
		l266:
			position, tokenIndex, depth = position266, tokenIndex266, depth266
			return false
		},
		/* 47 Field <- <(<fieldChar+> ':' Action29)> */
		nil,
		/* 48 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					{
						position278 := position
						depth++
						if buffer[position] != rune('n') {
							goto l277
						}
						position++
						if buffer[position] != rune('u') {
							goto l277
						}
						position++
						if buffer[position] != rune('l') {
							goto l277
						}
						position++
						if buffer[position] != rune('l') {
							goto l277
						}
						position++
						{
							add(ruleAction32, position)
						}
						depth--
						add(ruleNull, position278)
					}
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					{
						position281 := position
						depth++
						if buffer[position] != rune('M') {
							goto l280
						}
						position++
						if buffer[position] != rune('i') {
							goto l280
						}
						position++
						if buffer[position] != rune('n') {
							goto l280
						}
						position++
						if buffer[position] != rune('K') {
							goto l280
						}
						position++
						if buffer[position] != rune('e') {
							goto l280
						}
						position++
						if buffer[position] != rune('y') {
							goto l280
						}
						position++
						{
							add(ruleAction42, position)
						}
						depth--
						add(ruleMinKey, position281)
					}
					goto l276
				l280:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					{
						switch buffer[position] {
						case 'M':
							{
								position284 := position
								depth++
								if buffer[position] != rune('M') {
									goto l274
								}
								position++
								if buffer[position] != rune('a') {
									goto l274
								}
								position++
								if buffer[position] != rune('x') {
									goto l274
								}
								position++
								if buffer[position] != rune('K') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('y') {
									goto l274
								}
								position++
								{
									add(ruleAction43, position)
								}
								depth--
								add(ruleMaxKey, position284)
							}
							break
						case 'u':
							{
								position286 := position
								depth++
								if buffer[position] != rune('u') {
									goto l274
								}
								position++
								if buffer[position] != rune('n') {
									goto l274
								}
								position++
								if buffer[position] != rune('d') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('f') {
									goto l274
								}
								position++
								if buffer[position] != rune('i') {
									goto l274
								}
								position++
								if buffer[position] != rune('n') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('d') {
									goto l274
								}
								position++
								{
									add(ruleAction44, position)
								}
								depth--
								add(ruleUndefined, position286)
							}
							break
						case 'N':
							{
								position288 := position
								depth++
								if buffer[position] != rune('N') {
									goto l274
								}
								position++
								if buffer[position] != rune('u') {
									goto l274
								}
								position++
								if buffer[position] != rune('m') {
									goto l274
								}
								position++
								if buffer[position] != rune('b') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('r') {
									goto l274
								}
								position++
								if buffer[position] != rune('L') {
									goto l274
								}
								position++
								if buffer[position] != rune('o') {
									goto l274
								}
								position++
								if buffer[position] != rune('n') {
									goto l274
								}
								position++
								if buffer[position] != rune('g') {
									goto l274
								}
								position++
								if buffer[position] != rune('(') {
									goto l274
								}
								position++
								{
									position289 := position
									depth++
									{
										position292, tokenIndex292, depth292 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l292
										}
										position++
										goto l274
									l292:
										position, tokenIndex, depth = position292, tokenIndex292, depth292
									}
									if !matchDot() {
										goto l274
									}
								l290:
									{
										position291, tokenIndex291, depth291 := position, tokenIndex, depth
										{
											position293, tokenIndex293, depth293 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l293
											}
											position++
											goto l291
										l293:
											position, tokenIndex, depth = position293, tokenIndex293, depth293
										}
										if !matchDot() {
											goto l291
										}
										goto l290
									l291:
										position, tokenIndex, depth = position291, tokenIndex291, depth291
									}
									depth--
									add(rulePegText, position289)
								}
								if buffer[position] != rune(')') {
									goto l274
								}
								position++
								{
									add(ruleAction41, position)
								}
								depth--
								add(ruleNumberLong, position288)
							}
							break
						case '/':
							{
								position295 := position
								depth++
								if buffer[position] != rune('/') {
									goto l274
								}
								position++
								{
									position296 := position
									depth++
									{
										position297 := position
										depth++
										{
											position300 := position
											depth++
											{
												position301, tokenIndex301, depth301 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l301
												}
												position++
												goto l274
											l301:
												position, tokenIndex, depth = position301, tokenIndex301, depth301
											}
											if !matchDot() {
												goto l274
											}
											depth--
											add(ruleregexChar, position300)
										}
									l298:
										{
											position299, tokenIndex299, depth299 := position, tokenIndex, depth
											{
												position302 := position
												depth++
												{
													position303, tokenIndex303, depth303 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l303
													}
													position++
													goto l299
												l303:
													position, tokenIndex, depth = position303, tokenIndex303, depth303
												}
												if !matchDot() {
													goto l299
												}
												depth--
												add(ruleregexChar, position302)
											}
											goto l298
										l299:
											position, tokenIndex, depth = position299, tokenIndex299, depth299
										}
										if buffer[position] != rune('/') {
											goto l274
										}
										position++
									l304:
										{
											position305, tokenIndex305, depth305 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l305
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l305
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l305
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l305
													}
													position++
													break
												}
											}

											goto l304
										l305:
											position, tokenIndex, depth = position305, tokenIndex305, depth305
										}
										depth--
										add(ruleregexBody, position297)
									}
									depth--
									add(rulePegText, position296)
								}
								{
									add(ruleAction38, position)
								}
								depth--
								add(ruleRegex, position295)
							}
							break
						case 'T':
							{
								position308 := position
								depth++
								{
									position309, tokenIndex309, depth309 := position, tokenIndex, depth
									{
										position311 := position
										depth++
										if buffer[position] != rune('T') {
											goto l310
										}
										position++
										if buffer[position] != rune('i') {
											goto l310
										}
										position++
										if buffer[position] != rune('m') {
											goto l310
										}
										position++
										if buffer[position] != rune('e') {
											goto l310
										}
										position++
										if buffer[position] != rune('s') {
											goto l310
										}
										position++
										if buffer[position] != rune('t') {
											goto l310
										}
										position++
										if buffer[position] != rune('a') {
											goto l310
										}
										position++
										if buffer[position] != rune('m') {
											goto l310
										}
										position++
										if buffer[position] != rune('p') {
											goto l310
										}
										position++
										if buffer[position] != rune('(') {
											goto l310
										}
										position++
										{
											position312 := position
											depth++
											{
												position315, tokenIndex315, depth315 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l315
												}
												position++
												goto l310
											l315:
												position, tokenIndex, depth = position315, tokenIndex315, depth315
											}
											if !matchDot() {
												goto l310
											}
										l313:
											{
												position314, tokenIndex314, depth314 := position, tokenIndex, depth
												{
													position316, tokenIndex316, depth316 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l316
													}
													position++
													goto l314
												l316:
													position, tokenIndex, depth = position316, tokenIndex316, depth316
												}
												if !matchDot() {
													goto l314
												}
												goto l313
											l314:
												position, tokenIndex, depth = position314, tokenIndex314, depth314
											}
											depth--
											add(rulePegText, position312)
										}
										if buffer[position] != rune(')') {
											goto l310
										}
										position++
										{
											add(ruleAction39, position)
										}
										depth--
										add(ruletimestampParen, position311)
									}
									goto l309
								l310:
									position, tokenIndex, depth = position309, tokenIndex309, depth309
									{
										position318 := position
										depth++
										if buffer[position] != rune('T') {
											goto l274
										}
										position++
										if buffer[position] != rune('i') {
											goto l274
										}
										position++
										if buffer[position] != rune('m') {
											goto l274
										}
										position++
										if buffer[position] != rune('e') {
											goto l274
										}
										position++
										if buffer[position] != rune('s') {
											goto l274
										}
										position++
										if buffer[position] != rune('t') {
											goto l274
										}
										position++
										if buffer[position] != rune('a') {
											goto l274
										}
										position++
										if buffer[position] != rune('m') {
											goto l274
										}
										position++
										if buffer[position] != rune('p') {
											goto l274
										}
										position++
										if buffer[position] != rune(' ') {
											goto l274
										}
										position++
										{
											position319 := position
											depth++
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
													goto l274
												}
												position++
											}
										l322:
										l320:
											{
												position321, tokenIndex321, depth321 := position, tokenIndex, depth
												{
													position324, tokenIndex324, depth324 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l325
													}
													position++
													goto l324
												l325:
													position, tokenIndex, depth = position324, tokenIndex324, depth324
													if buffer[position] != rune('|') {
														goto l321
													}
													position++
												}
											l324:
												goto l320
											l321:
												position, tokenIndex, depth = position321, tokenIndex321, depth321
											}
											depth--
											add(rulePegText, position319)
										}
										{
											add(ruleAction40, position)
										}
										depth--
										add(ruletimestampPipe, position318)
									}
								}
							l309:
								depth--
								add(ruleTimestampVal, position308)
							}
							break
						case 'B':
							{
								position327 := position
								depth++
								if buffer[position] != rune('B') {
									goto l274
								}
								position++
								if buffer[position] != rune('i') {
									goto l274
								}
								position++
								if buffer[position] != rune('n') {
									goto l274
								}
								position++
								if buffer[position] != rune('D') {
									goto l274
								}
								position++
								if buffer[position] != rune('a') {
									goto l274
								}
								position++
								if buffer[position] != rune('t') {
									goto l274
								}
								position++
								if buffer[position] != rune('a') {
									goto l274
								}
								position++
								if buffer[position] != rune('(') {
									goto l274
								}
								position++
								{
									position328 := position
									depth++
									{
										position331, tokenIndex331, depth331 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l331
										}
										position++
										goto l274
									l331:
										position, tokenIndex, depth = position331, tokenIndex331, depth331
									}
									if !matchDot() {
										goto l274
									}
								l329:
									{
										position330, tokenIndex330, depth330 := position, tokenIndex, depth
										{
											position332, tokenIndex332, depth332 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l332
											}
											position++
											goto l330
										l332:
											position, tokenIndex, depth = position332, tokenIndex332, depth332
										}
										if !matchDot() {
											goto l330
										}
										goto l329
									l330:
										position, tokenIndex, depth = position330, tokenIndex330, depth330
									}
									depth--
									add(rulePegText, position328)
								}
								if buffer[position] != rune(')') {
									goto l274
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleBinData, position327)
							}
							break
						case 'D', 'n':
							{
								position334 := position
								depth++
								{
									position335, tokenIndex335, depth335 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l335
									}
									position++
									if buffer[position] != rune('e') {
										goto l335
									}
									position++
									if buffer[position] != rune('w') {
										goto l335
									}
									position++
									if buffer[position] != rune(' ') {
										goto l335
									}
									position++
									goto l336
								l335:
									position, tokenIndex, depth = position335, tokenIndex335, depth335
								}
							l336:
								if buffer[position] != rune('D') {
									goto l274
								}
								position++
								if buffer[position] != rune('a') {
									goto l274
								}
								position++
								if buffer[position] != rune('t') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('(') {
									goto l274
								}
								position++
								{
									position337, tokenIndex337, depth337 := position, tokenIndex, depth
									if buffer[position] != rune('-') {
										goto l337
									}
									position++
									goto l338
								l337:
									position, tokenIndex, depth = position337, tokenIndex337, depth337
								}
							l338:
								{
									position339 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l274
									}
									position++
								l340:
									{
										position341, tokenIndex341, depth341 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l341
										}
										position++
										goto l340
									l341:
										position, tokenIndex, depth = position341, tokenIndex341, depth341
									}
									depth--
									add(rulePegText, position339)
								}
								if buffer[position] != rune(')') {
									goto l274
								}
								position++
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleDate, position334)
							}
							break
						case 'O':
							{
								position343 := position
								depth++
								if buffer[position] != rune('O') {
									goto l274
								}
								position++
								if buffer[position] != rune('b') {
									goto l274
								}
								position++
								if buffer[position] != rune('j') {
									goto l274
								}
								position++
								if buffer[position] != rune('e') {
									goto l274
								}
								position++
								if buffer[position] != rune('c') {
									goto l274
								}
								position++
								if buffer[position] != rune('t') {
									goto l274
								}
								position++
								if buffer[position] != rune('I') {
									goto l274
								}
								position++
								if buffer[position] != rune('d') {
									goto l274
								}
								position++
								if buffer[position] != rune('(') {
									goto l274
								}
								position++
								{
									position344, tokenIndex344, depth344 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l345
									}
									position++
									goto l344
								l345:
									position, tokenIndex, depth = position344, tokenIndex344, depth344
									if buffer[position] != rune('"') {
										goto l274
									}
									position++
								}
							l344:
								{
									position346 := position
									depth++
								l347:
									{
										position348, tokenIndex348, depth348 := position, tokenIndex, depth
										{
											position349 := position
											depth++
											{
												position350, tokenIndex350, depth350 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l351
												}
												position++
												goto l350
											l351:
												position, tokenIndex, depth = position350, tokenIndex350, depth350
												{
													position352, tokenIndex352, depth352 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l353
													}
													position++
													goto l352
												l353:
													position, tokenIndex, depth = position352, tokenIndex352, depth352
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l348
													}
													position++
												}
											l352:
											}
										l350:
											depth--
											add(rulehexChar, position349)
										}
										goto l347
									l348:
										position, tokenIndex, depth = position348, tokenIndex348, depth348
									}
									depth--
									add(rulePegText, position346)
								}
								{
									position354, tokenIndex354, depth354 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l355
									}
									position++
									goto l354
								l355:
									position, tokenIndex, depth = position354, tokenIndex354, depth354
									if buffer[position] != rune('"') {
										goto l274
									}
									position++
								}
							l354:
								if buffer[position] != rune(')') {
									goto l274
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleObjectID, position343)
							}
							break
						case '"':
							{
								position357 := position
								depth++
								if buffer[position] != rune('"') {
									goto l274
								}
								position++
								{
									position358 := position
									depth++
								l359:
									{
										position360, tokenIndex360, depth360 := position, tokenIndex, depth
										{
											position361 := position
											depth++
											{
												position362, tokenIndex362, depth362 := position, tokenIndex, depth
												{
													position364, tokenIndex364, depth364 := position, tokenIndex, depth
													{
														position365, tokenIndex365, depth365 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l366
														}
														position++
														goto l365
													l366:
														position, tokenIndex, depth = position365, tokenIndex365, depth365
														if buffer[position] != rune('\\') {
															goto l364
														}
														position++
													}
												l365:
													goto l363
												l364:
													position, tokenIndex, depth = position364, tokenIndex364, depth364
												}
												if !matchDot() {
													goto l363
												}
												goto l362
											l363:
												position, tokenIndex, depth = position362, tokenIndex362, depth362
												if buffer[position] != rune('\\') {
													goto l360
												}
												position++
												if !matchDot() {
													goto l360
												}
											}
										l362:
											depth--
											add(rulestringChar, position361)
										}
										goto l359
									l360:
										position, tokenIndex, depth = position360, tokenIndex360, depth360
									}
									depth--
									add(rulePegText, position358)
								}
								if buffer[position] != rune('"') {
									goto l274
								}
								position++
								{
									add(ruleAction31, position)
								}
								depth--
								add(ruleString, position357)
							}
							break
						case 'f', 't':
							{
								position368 := position
								depth++
								{
									position369, tokenIndex369, depth369 := position, tokenIndex, depth
									{
										position371 := position
										depth++
										if buffer[position] != rune('t') {
											goto l370
										}
										position++
										if buffer[position] != rune('r') {
											goto l370
										}
										position++
										if buffer[position] != rune('u') {
											goto l370
										}
										position++
										if buffer[position] != rune('e') {
											goto l370
										}
										position++
										{
											add(ruleAction33, position)
										}
										depth--
										add(ruleTrue, position371)
									}
									goto l369
								l370:
									position, tokenIndex, depth = position369, tokenIndex369, depth369
									{
										position373 := position
										depth++
										if buffer[position] != rune('f') {
											goto l274
										}
										position++
										if buffer[position] != rune('a') {
											goto l274
										}
										position++
										if buffer[position] != rune('l') {
											goto l274
										}
										position++
										if buffer[position] != rune('s') {
											goto l274
										}
										position++
										if buffer[position] != rune('e') {
											goto l274
										}
										position++
										{
											add(ruleAction34, position)
										}
										depth--
										add(ruleFalse, position373)
									}
								}
							l369:
								depth--
								add(ruleBoolean, position368)
							}
							break
						case '[':
							{
								position375 := position
								depth++
								if buffer[position] != rune('[') {
									goto l274
								}
								position++
								{
									add(ruleAction26, position)
								}
								{
									position377, tokenIndex377, depth377 := position, tokenIndex, depth
									{
										position379 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l377
										}
									l380:
										{
											position381, tokenIndex381, depth381 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l381
											}
											position++
											if !_rules[ruleListElem]() {
												goto l381
											}
											goto l380
										l381:
											position, tokenIndex, depth = position381, tokenIndex381, depth381
										}
										depth--
										add(ruleListElements, position379)
									}
									goto l378
								l377:
									position, tokenIndex, depth = position377, tokenIndex377, depth377
								}
							l378:
								if buffer[position] != rune(']') {
									goto l274
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleList, position375)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l274
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l274
							}
							break
						}
					}

				}
			l276:
				depth--
				add(ruleValue, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
			return false
		},
		/* 49 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action30)> */
		func() bool {
			position383, tokenIndex383, depth383 := position, tokenIndex, depth
			{
				position384 := position
				depth++
				{
					position385 := position
					depth++
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l386
						}
						position++
						goto l387
					l386:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
					}
				l387:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l383
					}
					position++
				l388:
					{
						position389, tokenIndex389, depth389 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l389
						}
						position++
						goto l388
					l389:
						position, tokenIndex, depth = position389, tokenIndex389, depth389
					}
					{
						position390, tokenIndex390, depth390 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l390
						}
						position++
						goto l391
					l390:
						position, tokenIndex, depth = position390, tokenIndex390, depth390
					}
				l391:
				l392:
					{
						position393, tokenIndex393, depth393 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l393
						}
						position++
						goto l392
					l393:
						position, tokenIndex, depth = position393, tokenIndex393, depth393
					}
					depth--
					add(rulePegText, position385)
				}
				{
					add(ruleAction30, position)
				}
				depth--
				add(ruleNumeric, position384)
			}
			return true
		l383:
			position, tokenIndex, depth = position383, tokenIndex383, depth383
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
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l415
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l415
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l415
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l415
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l415
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l415
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l415
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		nil,
		/* 73 Action0 <- <{ p.SetField("log_level", buffer[begin:end]) }> */
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
