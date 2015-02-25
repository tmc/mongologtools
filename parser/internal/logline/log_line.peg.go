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
	rules  [115]func() bool
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
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction19:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction20:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction21:
			p.PushMap()
		case ruleAction22:
			p.PopMap()
		case ruleAction23:
			p.SetMapValue()
		case ruleAction24:
			p.PushList()
		case ruleAction25:
			p.PopList()
		case ruleAction26:
			p.SetListValue()
		case ruleAction27:
			p.PushField(buffer[begin:end])
		case ruleAction28:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction29:
			p.PushValue(buffer[begin:end])
		case ruleAction30:
			p.PushValue(nil)
		case ruleAction31:
			p.PushValue(true)
		case ruleAction32:
			p.PushValue(false)
		case ruleAction33:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction34:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction35:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction36:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction37:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction38:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction39:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction40:
			p.PushValue(p.Minkey())
		case ruleAction41:
			p.PushValue(p.Maxkey())
		case ruleAction42:
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
								add(ruleAction18, position)
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
								add(ruleAction19, position)
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
							add(ruleAction20, position)
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
		/* 8 LineField <- <((commandField / planSummaryField / plainField) S?)> */
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
						if buffer[position] != rune('c') {
							goto l132
						}
						position++
						if buffer[position] != rune('o') {
							goto l132
						}
						position++
						if buffer[position] != rune('m') {
							goto l132
						}
						position++
						if buffer[position] != rune('m') {
							goto l132
						}
						position++
						if buffer[position] != rune('a') {
							goto l132
						}
						position++
						if buffer[position] != rune('n') {
							goto l132
						}
						position++
						if buffer[position] != rune('d') {
							goto l132
						}
						position++
						if buffer[position] != rune(':') {
							goto l132
						}
						position++
						if buffer[position] != rune(' ') {
							goto l132
						}
						position++
						{
							position134 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l132
							}
						l135:
							{
								position136, tokenIndex136, depth136 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l136
								}
								goto l135
							l136:
								position, tokenIndex, depth = position136, tokenIndex136, depth136
							}
							depth--
							add(rulePegText, position134)
						}
						{
							position137, tokenIndex137, depth137 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l137
							}
							goto l138
						l137:
							position, tokenIndex, depth = position137, tokenIndex137, depth137
						}
					l138:
						{
							add(ruleAction11, position)
						}
						if !_rules[ruleLineValue]() {
							goto l132
						}
						{
							add(ruleAction12, position)
						}
						depth--
						add(rulecommandField, position133)
					}
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position142 := position
						depth++
						if buffer[position] != rune('p') {
							goto l141
						}
						position++
						if buffer[position] != rune('l') {
							goto l141
						}
						position++
						if buffer[position] != rune('a') {
							goto l141
						}
						position++
						if buffer[position] != rune('n') {
							goto l141
						}
						position++
						if buffer[position] != rune('S') {
							goto l141
						}
						position++
						if buffer[position] != rune('u') {
							goto l141
						}
						position++
						if buffer[position] != rune('m') {
							goto l141
						}
						position++
						if buffer[position] != rune('m') {
							goto l141
						}
						position++
						if buffer[position] != rune('a') {
							goto l141
						}
						position++
						if buffer[position] != rune('r') {
							goto l141
						}
						position++
						if buffer[position] != rune('y') {
							goto l141
						}
						position++
						if buffer[position] != rune(':') {
							goto l141
						}
						position++
						if buffer[position] != rune(' ') {
							goto l141
						}
						position++
						{
							add(ruleAction13, position)
						}
						{
							position144 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l141
							}
						l145:
							{
								position146, tokenIndex146, depth146 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l146
								}
								position++
								if buffer[position] != rune(' ') {
									goto l146
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l146
								}
								goto l145
							l146:
								position, tokenIndex, depth = position146, tokenIndex146, depth146
							}
							depth--
							add(ruleplanSummaryElements, position144)
						}
						{
							add(ruleAction14, position)
						}
						depth--
						add(ruleplanSummaryField, position142)
					}
					goto l131
				l141:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position148 := position
						depth++
						{
							position149 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l129
							}
						l150:
							{
								position151, tokenIndex151, depth151 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l151
								}
								goto l150
							l151:
								position, tokenIndex, depth = position151, tokenIndex151, depth151
							}
							depth--
							add(rulePegText, position149)
						}
						if buffer[position] != rune(':') {
							goto l129
						}
						position++
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l152
							}
							goto l153
						l152:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
						}
					l153:
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
						add(ruleplainField, position148)
					}
				}
			l131:
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l156
					}
					goto l157
				l156:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
				}
			l157:
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
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				{
					position168 := position
					depth++
					{
						position169 := position
						depth++
						{
							position170, tokenIndex170, depth170 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l171
							}
							position++
							if buffer[position] != rune('N') {
								goto l171
							}
							position++
							if buffer[position] != rune('D') {
								goto l171
							}
							position++
							if buffer[position] != rune('_') {
								goto l171
							}
							position++
							if buffer[position] != rune('H') {
								goto l171
							}
							position++
							if buffer[position] != rune('A') {
								goto l171
							}
							position++
							if buffer[position] != rune('S') {
								goto l171
							}
							position++
							if buffer[position] != rune('H') {
								goto l171
							}
							position++
							goto l170
						l171:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('C') {
								goto l172
							}
							position++
							if buffer[position] != rune('A') {
								goto l172
							}
							position++
							if buffer[position] != rune('C') {
								goto l172
							}
							position++
							if buffer[position] != rune('H') {
								goto l172
							}
							position++
							if buffer[position] != rune('E') {
								goto l172
							}
							position++
							if buffer[position] != rune('D') {
								goto l172
							}
							position++
							if buffer[position] != rune('_') {
								goto l172
							}
							position++
							if buffer[position] != rune('P') {
								goto l172
							}
							position++
							if buffer[position] != rune('L') {
								goto l172
							}
							position++
							if buffer[position] != rune('A') {
								goto l172
							}
							position++
							if buffer[position] != rune('N') {
								goto l172
							}
							position++
							goto l170
						l172:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('C') {
								goto l173
							}
							position++
							if buffer[position] != rune('O') {
								goto l173
							}
							position++
							if buffer[position] != rune('L') {
								goto l173
							}
							position++
							if buffer[position] != rune('L') {
								goto l173
							}
							position++
							if buffer[position] != rune('S') {
								goto l173
							}
							position++
							if buffer[position] != rune('C') {
								goto l173
							}
							position++
							if buffer[position] != rune('A') {
								goto l173
							}
							position++
							if buffer[position] != rune('N') {
								goto l173
							}
							position++
							goto l170
						l173:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('C') {
								goto l174
							}
							position++
							if buffer[position] != rune('O') {
								goto l174
							}
							position++
							if buffer[position] != rune('U') {
								goto l174
							}
							position++
							if buffer[position] != rune('N') {
								goto l174
							}
							position++
							if buffer[position] != rune('T') {
								goto l174
							}
							position++
							goto l170
						l174:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('D') {
								goto l175
							}
							position++
							if buffer[position] != rune('E') {
								goto l175
							}
							position++
							if buffer[position] != rune('L') {
								goto l175
							}
							position++
							if buffer[position] != rune('E') {
								goto l175
							}
							position++
							if buffer[position] != rune('T') {
								goto l175
							}
							position++
							if buffer[position] != rune('E') {
								goto l175
							}
							position++
							goto l170
						l175:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('G') {
								goto l176
							}
							position++
							if buffer[position] != rune('E') {
								goto l176
							}
							position++
							if buffer[position] != rune('O') {
								goto l176
							}
							position++
							if buffer[position] != rune('_') {
								goto l176
							}
							position++
							if buffer[position] != rune('N') {
								goto l176
							}
							position++
							if buffer[position] != rune('E') {
								goto l176
							}
							position++
							if buffer[position] != rune('A') {
								goto l176
							}
							position++
							if buffer[position] != rune('R') {
								goto l176
							}
							position++
							if buffer[position] != rune('_') {
								goto l176
							}
							position++
							if buffer[position] != rune('2') {
								goto l176
							}
							position++
							if buffer[position] != rune('D') {
								goto l176
							}
							position++
							goto l170
						l176:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('G') {
								goto l177
							}
							position++
							if buffer[position] != rune('E') {
								goto l177
							}
							position++
							if buffer[position] != rune('O') {
								goto l177
							}
							position++
							if buffer[position] != rune('_') {
								goto l177
							}
							position++
							if buffer[position] != rune('N') {
								goto l177
							}
							position++
							if buffer[position] != rune('E') {
								goto l177
							}
							position++
							if buffer[position] != rune('A') {
								goto l177
							}
							position++
							if buffer[position] != rune('R') {
								goto l177
							}
							position++
							if buffer[position] != rune('_') {
								goto l177
							}
							position++
							if buffer[position] != rune('2') {
								goto l177
							}
							position++
							if buffer[position] != rune('D') {
								goto l177
							}
							position++
							if buffer[position] != rune('S') {
								goto l177
							}
							position++
							if buffer[position] != rune('P') {
								goto l177
							}
							position++
							if buffer[position] != rune('H') {
								goto l177
							}
							position++
							if buffer[position] != rune('E') {
								goto l177
							}
							position++
							if buffer[position] != rune('R') {
								goto l177
							}
							position++
							if buffer[position] != rune('E') {
								goto l177
							}
							position++
							goto l170
						l177:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('I') {
								goto l178
							}
							position++
							if buffer[position] != rune('D') {
								goto l178
							}
							position++
							if buffer[position] != rune('H') {
								goto l178
							}
							position++
							if buffer[position] != rune('A') {
								goto l178
							}
							position++
							if buffer[position] != rune('C') {
								goto l178
							}
							position++
							if buffer[position] != rune('K') {
								goto l178
							}
							position++
							goto l170
						l178:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
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
							if buffer[position] != rune('_') {
								goto l179
							}
							position++
							if buffer[position] != rune('M') {
								goto l179
							}
							position++
							if buffer[position] != rune('E') {
								goto l179
							}
							position++
							if buffer[position] != rune('R') {
								goto l179
							}
							position++
							if buffer[position] != rune('G') {
								goto l179
							}
							position++
							if buffer[position] != rune('E') {
								goto l179
							}
							position++
							goto l170
						l179:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('S') {
								goto l180
							}
							position++
							if buffer[position] != rune('H') {
								goto l180
							}
							position++
							if buffer[position] != rune('A') {
								goto l180
							}
							position++
							if buffer[position] != rune('R') {
								goto l180
							}
							position++
							if buffer[position] != rune('D') {
								goto l180
							}
							position++
							if buffer[position] != rune('I') {
								goto l180
							}
							position++
							if buffer[position] != rune('N') {
								goto l180
							}
							position++
							if buffer[position] != rune('G') {
								goto l180
							}
							position++
							if buffer[position] != rune('_') {
								goto l180
							}
							position++
							if buffer[position] != rune('F') {
								goto l180
							}
							position++
							if buffer[position] != rune('I') {
								goto l180
							}
							position++
							if buffer[position] != rune('L') {
								goto l180
							}
							position++
							if buffer[position] != rune('T') {
								goto l180
							}
							position++
							if buffer[position] != rune('E') {
								goto l180
							}
							position++
							if buffer[position] != rune('R') {
								goto l180
							}
							position++
							goto l170
						l180:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('S') {
								goto l181
							}
							position++
							if buffer[position] != rune('K') {
								goto l181
							}
							position++
							if buffer[position] != rune('I') {
								goto l181
							}
							position++
							if buffer[position] != rune('P') {
								goto l181
							}
							position++
							goto l170
						l181:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('S') {
								goto l182
							}
							position++
							if buffer[position] != rune('O') {
								goto l182
							}
							position++
							if buffer[position] != rune('R') {
								goto l182
							}
							position++
							if buffer[position] != rune('T') {
								goto l182
							}
							position++
							goto l170
						l182:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('X') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('B') {
										goto l166
									}
									position++
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									if buffer[position] != rune('L') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									if buffer[position] != rune('_') {
										goto l166
									}
									position++
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									if buffer[position] != rune('R') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('J') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('R') {
										goto l166
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('L') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('_') {
										goto l166
									}
									position++
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									if buffer[position] != rune('L') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('M') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									if buffer[position] != rune('_') {
										goto l166
									}
									position++
									if buffer[position] != rune('M') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('X') {
										goto l166
									}
									position++
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l166
									}
									position++
									if buffer[position] != rune('R') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('P') {
										goto l166
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('H') {
										goto l166
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('F') {
										goto l166
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('I') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('U') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('_') {
										goto l166
									}
									position++
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									if buffer[position] != rune('C') {
										goto l166
									}
									position++
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l166
									}
									position++
									if buffer[position] != rune('N') {
										goto l166
									}
									position++
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									if buffer[position] != rune('_') {
										goto l166
									}
									position++
									if buffer[position] != rune('S') {
										goto l166
									}
									position++
									if buffer[position] != rune('O') {
										goto l166
									}
									position++
									if buffer[position] != rune('R') {
										goto l166
									}
									position++
									if buffer[position] != rune('T') {
										goto l166
									}
									position++
									if buffer[position] != rune('E') {
										goto l166
									}
									position++
									if buffer[position] != rune('D') {
										goto l166
									}
									position++
									break
								}
							}

						}
					l170:
						depth--
						add(ruleplanSummaryStage, position169)
					}
					depth--
					add(rulePegText, position168)
				}
				{
					add(ruleAction15, position)
				}
				{
					position185 := position
					depth++
					{
						position186, tokenIndex186, depth186 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l187
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l187
						}
						{
							add(ruleAction16, position)
						}
						goto l186
					l187:
						position, tokenIndex, depth = position186, tokenIndex186, depth186
						{
							add(ruleAction17, position)
						}
					}
				l186:
					depth--
					add(ruleplanSummary, position185)
				}
				depth--
				add(ruleplanSummaryElem, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 18 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 19 planSummary <- <((' ' LineValue Action16) / Action17)> */
		nil,
		/* 20 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position192, tokenIndex192, depth192 := position, tokenIndex, depth
			{
				position193 := position
				depth++
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l195
					}
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if !_rules[ruleNumeric]() {
						goto l192
					}
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l196
					}
					goto l197
				l196:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
				}
			l197:
				depth--
				add(ruleLineValue, position193)
			}
			return true
		l192:
			position, tokenIndex, depth = position192, tokenIndex192, depth192
			return false
		},
		/* 21 timestamp24 <- <(<(date ' ' time)> Action18)> */
		nil,
		/* 22 timestamp26 <- <(<datetime26> Action19)> */
		nil,
		/* 23 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 24 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 25 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l202
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l202
				}
				position++
				depth--
				add(ruledigit2, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 26 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 27 tz <- <('+' [0-9]+)> */
		nil,
		/* 28 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position206, tokenIndex206, depth206 := position, tokenIndex, depth
			{
				position207 := position
				depth++
				{
					position208 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l206
					}
					depth--
					add(rulehour, position208)
				}
				if buffer[position] != rune(':') {
					goto l206
				}
				position++
				{
					position209 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l206
					}
					depth--
					add(ruleminute, position209)
				}
				if buffer[position] != rune(':') {
					goto l206
				}
				position++
				{
					position210 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l206
					}
					depth--
					add(rulesecond, position210)
				}
				if buffer[position] != rune('.') {
					goto l206
				}
				position++
				{
					position211 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l206
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l206
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l206
					}
					position++
					depth--
					add(rulemillisecond, position211)
				}
				depth--
				add(ruletime, position207)
			}
			return true
		l206:
			position, tokenIndex, depth = position206, tokenIndex206, depth206
			return false
		},
		/* 29 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 30 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 31 dayNum <- <([0-9] [0-9]?)> */
		nil,
		/* 32 hour <- <digit2> */
		nil,
		/* 33 minute <- <digit2> */
		nil,
		/* 34 second <- <digit2> */
		nil,
		/* 35 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 36 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 37 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 38 extra <- <(<.+> Action20)> */
		nil,
		/* 39 S <- <' '+> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l222
				}
				position++
			l224:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
				}
				depth--
				add(ruleS, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 40 Doc <- <('{' Action21 DocElements? '}' Action22)> */
		func() bool {
			position226, tokenIndex226, depth226 := position, tokenIndex, depth
			{
				position227 := position
				depth++
				if buffer[position] != rune('{') {
					goto l226
				}
				position++
				{
					add(ruleAction21, position)
				}
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					{
						position231 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l229
						}
					l232:
						{
							position233, tokenIndex233, depth233 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l233
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l233
							}
							goto l232
						l233:
							position, tokenIndex, depth = position233, tokenIndex233, depth233
						}
						depth--
						add(ruleDocElements, position231)
					}
					goto l230
				l229:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
				}
			l230:
				if buffer[position] != rune('}') {
					goto l226
				}
				position++
				{
					add(ruleAction22, position)
				}
				depth--
				add(ruleDoc, position227)
			}
			return true
		l226:
			position, tokenIndex, depth = position226, tokenIndex226, depth226
			return false
		},
		/* 41 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 42 DocElem <- <(S? Field S? Value S? Action23)> */
		func() bool {
			position236, tokenIndex236, depth236 := position, tokenIndex, depth
			{
				position237 := position
				depth++
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l238
					}
					goto l239
				l238:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
				}
			l239:
				{
					position240 := position
					depth++
					{
						position241 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l236
						}
					l242:
						{
							position243, tokenIndex243, depth243 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l243
							}
							goto l242
						l243:
							position, tokenIndex, depth = position243, tokenIndex243, depth243
						}
						depth--
						add(rulePegText, position241)
					}
					if buffer[position] != rune(':') {
						goto l236
					}
					position++
					{
						add(ruleAction27, position)
					}
					depth--
					add(ruleField, position240)
				}
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l245
					}
					goto l246
				l245:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
				}
			l246:
				if !_rules[ruleValue]() {
					goto l236
				}
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l247
					}
					goto l248
				l247:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
				}
			l248:
				{
					add(ruleAction23, position)
				}
				depth--
				add(ruleDocElem, position237)
			}
			return true
		l236:
			position, tokenIndex, depth = position236, tokenIndex236, depth236
			return false
		},
		/* 43 List <- <('[' Action24 ListElements? ']' Action25)> */
		nil,
		/* 44 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 45 ListElem <- <(S? Value S? Action26)> */
		func() bool {
			position252, tokenIndex252, depth252 := position, tokenIndex, depth
			{
				position253 := position
				depth++
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l254
					}
					goto l255
				l254:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
				}
			l255:
				if !_rules[ruleValue]() {
					goto l252
				}
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l256
					}
					goto l257
				l256:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
				}
			l257:
				{
					add(ruleAction26, position)
				}
				depth--
				add(ruleListElem, position253)
			}
			return true
		l252:
			position, tokenIndex, depth = position252, tokenIndex252, depth252
			return false
		},
		/* 46 Field <- <(<fieldChar+> ':' Action27)> */
		nil,
		/* 47 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position260, tokenIndex260, depth260 := position, tokenIndex, depth
			{
				position261 := position
				depth++
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					{
						position264 := position
						depth++
						if buffer[position] != rune('n') {
							goto l263
						}
						position++
						if buffer[position] != rune('u') {
							goto l263
						}
						position++
						if buffer[position] != rune('l') {
							goto l263
						}
						position++
						if buffer[position] != rune('l') {
							goto l263
						}
						position++
						{
							add(ruleAction30, position)
						}
						depth--
						add(ruleNull, position264)
					}
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					{
						position267 := position
						depth++
						if buffer[position] != rune('M') {
							goto l266
						}
						position++
						if buffer[position] != rune('i') {
							goto l266
						}
						position++
						if buffer[position] != rune('n') {
							goto l266
						}
						position++
						if buffer[position] != rune('K') {
							goto l266
						}
						position++
						if buffer[position] != rune('e') {
							goto l266
						}
						position++
						if buffer[position] != rune('y') {
							goto l266
						}
						position++
						{
							add(ruleAction40, position)
						}
						depth--
						add(ruleMinKey, position267)
					}
					goto l262
				l266:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					{
						switch buffer[position] {
						case 'M':
							{
								position270 := position
								depth++
								if buffer[position] != rune('M') {
									goto l260
								}
								position++
								if buffer[position] != rune('a') {
									goto l260
								}
								position++
								if buffer[position] != rune('x') {
									goto l260
								}
								position++
								if buffer[position] != rune('K') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('y') {
									goto l260
								}
								position++
								{
									add(ruleAction41, position)
								}
								depth--
								add(ruleMaxKey, position270)
							}
							break
						case 'u':
							{
								position272 := position
								depth++
								if buffer[position] != rune('u') {
									goto l260
								}
								position++
								if buffer[position] != rune('n') {
									goto l260
								}
								position++
								if buffer[position] != rune('d') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('f') {
									goto l260
								}
								position++
								if buffer[position] != rune('i') {
									goto l260
								}
								position++
								if buffer[position] != rune('n') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('d') {
									goto l260
								}
								position++
								{
									add(ruleAction42, position)
								}
								depth--
								add(ruleUndefined, position272)
							}
							break
						case 'N':
							{
								position274 := position
								depth++
								if buffer[position] != rune('N') {
									goto l260
								}
								position++
								if buffer[position] != rune('u') {
									goto l260
								}
								position++
								if buffer[position] != rune('m') {
									goto l260
								}
								position++
								if buffer[position] != rune('b') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('r') {
									goto l260
								}
								position++
								if buffer[position] != rune('L') {
									goto l260
								}
								position++
								if buffer[position] != rune('o') {
									goto l260
								}
								position++
								if buffer[position] != rune('n') {
									goto l260
								}
								position++
								if buffer[position] != rune('g') {
									goto l260
								}
								position++
								if buffer[position] != rune('(') {
									goto l260
								}
								position++
								{
									position275 := position
									depth++
									{
										position278, tokenIndex278, depth278 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l278
										}
										position++
										goto l260
									l278:
										position, tokenIndex, depth = position278, tokenIndex278, depth278
									}
									if !matchDot() {
										goto l260
									}
								l276:
									{
										position277, tokenIndex277, depth277 := position, tokenIndex, depth
										{
											position279, tokenIndex279, depth279 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l279
											}
											position++
											goto l277
										l279:
											position, tokenIndex, depth = position279, tokenIndex279, depth279
										}
										if !matchDot() {
											goto l277
										}
										goto l276
									l277:
										position, tokenIndex, depth = position277, tokenIndex277, depth277
									}
									depth--
									add(rulePegText, position275)
								}
								if buffer[position] != rune(')') {
									goto l260
								}
								position++
								{
									add(ruleAction39, position)
								}
								depth--
								add(ruleNumberLong, position274)
							}
							break
						case '/':
							{
								position281 := position
								depth++
								if buffer[position] != rune('/') {
									goto l260
								}
								position++
								{
									position282 := position
									depth++
									{
										position283 := position
										depth++
										{
											position286 := position
											depth++
											{
												position287, tokenIndex287, depth287 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l287
												}
												position++
												goto l260
											l287:
												position, tokenIndex, depth = position287, tokenIndex287, depth287
											}
											if !matchDot() {
												goto l260
											}
											depth--
											add(ruleregexChar, position286)
										}
									l284:
										{
											position285, tokenIndex285, depth285 := position, tokenIndex, depth
											{
												position288 := position
												depth++
												{
													position289, tokenIndex289, depth289 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l289
													}
													position++
													goto l285
												l289:
													position, tokenIndex, depth = position289, tokenIndex289, depth289
												}
												if !matchDot() {
													goto l285
												}
												depth--
												add(ruleregexChar, position288)
											}
											goto l284
										l285:
											position, tokenIndex, depth = position285, tokenIndex285, depth285
										}
										if buffer[position] != rune('/') {
											goto l260
										}
										position++
									l290:
										{
											position291, tokenIndex291, depth291 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l291
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l291
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l291
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l291
													}
													position++
													break
												}
											}

											goto l290
										l291:
											position, tokenIndex, depth = position291, tokenIndex291, depth291
										}
										depth--
										add(ruleregexBody, position283)
									}
									depth--
									add(rulePegText, position282)
								}
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleRegex, position281)
							}
							break
						case 'T':
							{
								position294 := position
								depth++
								{
									position295, tokenIndex295, depth295 := position, tokenIndex, depth
									{
										position297 := position
										depth++
										if buffer[position] != rune('T') {
											goto l296
										}
										position++
										if buffer[position] != rune('i') {
											goto l296
										}
										position++
										if buffer[position] != rune('m') {
											goto l296
										}
										position++
										if buffer[position] != rune('e') {
											goto l296
										}
										position++
										if buffer[position] != rune('s') {
											goto l296
										}
										position++
										if buffer[position] != rune('t') {
											goto l296
										}
										position++
										if buffer[position] != rune('a') {
											goto l296
										}
										position++
										if buffer[position] != rune('m') {
											goto l296
										}
										position++
										if buffer[position] != rune('p') {
											goto l296
										}
										position++
										if buffer[position] != rune('(') {
											goto l296
										}
										position++
										{
											position298 := position
											depth++
											{
												position301, tokenIndex301, depth301 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l301
												}
												position++
												goto l296
											l301:
												position, tokenIndex, depth = position301, tokenIndex301, depth301
											}
											if !matchDot() {
												goto l296
											}
										l299:
											{
												position300, tokenIndex300, depth300 := position, tokenIndex, depth
												{
													position302, tokenIndex302, depth302 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l302
													}
													position++
													goto l300
												l302:
													position, tokenIndex, depth = position302, tokenIndex302, depth302
												}
												if !matchDot() {
													goto l300
												}
												goto l299
											l300:
												position, tokenIndex, depth = position300, tokenIndex300, depth300
											}
											depth--
											add(rulePegText, position298)
										}
										if buffer[position] != rune(')') {
											goto l296
										}
										position++
										{
											add(ruleAction37, position)
										}
										depth--
										add(ruletimestampParen, position297)
									}
									goto l295
								l296:
									position, tokenIndex, depth = position295, tokenIndex295, depth295
									{
										position304 := position
										depth++
										if buffer[position] != rune('T') {
											goto l260
										}
										position++
										if buffer[position] != rune('i') {
											goto l260
										}
										position++
										if buffer[position] != rune('m') {
											goto l260
										}
										position++
										if buffer[position] != rune('e') {
											goto l260
										}
										position++
										if buffer[position] != rune('s') {
											goto l260
										}
										position++
										if buffer[position] != rune('t') {
											goto l260
										}
										position++
										if buffer[position] != rune('a') {
											goto l260
										}
										position++
										if buffer[position] != rune('m') {
											goto l260
										}
										position++
										if buffer[position] != rune('p') {
											goto l260
										}
										position++
										if buffer[position] != rune(' ') {
											goto l260
										}
										position++
										{
											position305 := position
											depth++
											{
												position308, tokenIndex308, depth308 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l309
												}
												position++
												goto l308
											l309:
												position, tokenIndex, depth = position308, tokenIndex308, depth308
												if buffer[position] != rune('|') {
													goto l260
												}
												position++
											}
										l308:
										l306:
											{
												position307, tokenIndex307, depth307 := position, tokenIndex, depth
												{
													position310, tokenIndex310, depth310 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l311
													}
													position++
													goto l310
												l311:
													position, tokenIndex, depth = position310, tokenIndex310, depth310
													if buffer[position] != rune('|') {
														goto l307
													}
													position++
												}
											l310:
												goto l306
											l307:
												position, tokenIndex, depth = position307, tokenIndex307, depth307
											}
											depth--
											add(rulePegText, position305)
										}
										{
											add(ruleAction38, position)
										}
										depth--
										add(ruletimestampPipe, position304)
									}
								}
							l295:
								depth--
								add(ruleTimestampVal, position294)
							}
							break
						case 'B':
							{
								position313 := position
								depth++
								if buffer[position] != rune('B') {
									goto l260
								}
								position++
								if buffer[position] != rune('i') {
									goto l260
								}
								position++
								if buffer[position] != rune('n') {
									goto l260
								}
								position++
								if buffer[position] != rune('D') {
									goto l260
								}
								position++
								if buffer[position] != rune('a') {
									goto l260
								}
								position++
								if buffer[position] != rune('t') {
									goto l260
								}
								position++
								if buffer[position] != rune('a') {
									goto l260
								}
								position++
								if buffer[position] != rune('(') {
									goto l260
								}
								position++
								{
									position314 := position
									depth++
									{
										position317, tokenIndex317, depth317 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l317
										}
										position++
										goto l260
									l317:
										position, tokenIndex, depth = position317, tokenIndex317, depth317
									}
									if !matchDot() {
										goto l260
									}
								l315:
									{
										position316, tokenIndex316, depth316 := position, tokenIndex, depth
										{
											position318, tokenIndex318, depth318 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l318
											}
											position++
											goto l316
										l318:
											position, tokenIndex, depth = position318, tokenIndex318, depth318
										}
										if !matchDot() {
											goto l316
										}
										goto l315
									l316:
										position, tokenIndex, depth = position316, tokenIndex316, depth316
									}
									depth--
									add(rulePegText, position314)
								}
								if buffer[position] != rune(')') {
									goto l260
								}
								position++
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleBinData, position313)
							}
							break
						case 'D', 'n':
							{
								position320 := position
								depth++
								{
									position321, tokenIndex321, depth321 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l321
									}
									position++
									if buffer[position] != rune('e') {
										goto l321
									}
									position++
									if buffer[position] != rune('w') {
										goto l321
									}
									position++
									if buffer[position] != rune(' ') {
										goto l321
									}
									position++
									goto l322
								l321:
									position, tokenIndex, depth = position321, tokenIndex321, depth321
								}
							l322:
								if buffer[position] != rune('D') {
									goto l260
								}
								position++
								if buffer[position] != rune('a') {
									goto l260
								}
								position++
								if buffer[position] != rune('t') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('(') {
									goto l260
								}
								position++
								{
									position323 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l260
									}
									position++
								l324:
									{
										position325, tokenIndex325, depth325 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l325
										}
										position++
										goto l324
									l325:
										position, tokenIndex, depth = position325, tokenIndex325, depth325
									}
									depth--
									add(rulePegText, position323)
								}
								if buffer[position] != rune(')') {
									goto l260
								}
								position++
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleDate, position320)
							}
							break
						case 'O':
							{
								position327 := position
								depth++
								if buffer[position] != rune('O') {
									goto l260
								}
								position++
								if buffer[position] != rune('b') {
									goto l260
								}
								position++
								if buffer[position] != rune('j') {
									goto l260
								}
								position++
								if buffer[position] != rune('e') {
									goto l260
								}
								position++
								if buffer[position] != rune('c') {
									goto l260
								}
								position++
								if buffer[position] != rune('t') {
									goto l260
								}
								position++
								if buffer[position] != rune('I') {
									goto l260
								}
								position++
								if buffer[position] != rune('d') {
									goto l260
								}
								position++
								if buffer[position] != rune('(') {
									goto l260
								}
								position++
								{
									position328, tokenIndex328, depth328 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l329
									}
									position++
									goto l328
								l329:
									position, tokenIndex, depth = position328, tokenIndex328, depth328
									if buffer[position] != rune('"') {
										goto l260
									}
									position++
								}
							l328:
								{
									position330 := position
									depth++
								l331:
									{
										position332, tokenIndex332, depth332 := position, tokenIndex, depth
										{
											position333 := position
											depth++
											{
												position334, tokenIndex334, depth334 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l335
												}
												position++
												goto l334
											l335:
												position, tokenIndex, depth = position334, tokenIndex334, depth334
												{
													position336, tokenIndex336, depth336 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l337
													}
													position++
													goto l336
												l337:
													position, tokenIndex, depth = position336, tokenIndex336, depth336
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l332
													}
													position++
												}
											l336:
											}
										l334:
											depth--
											add(rulehexChar, position333)
										}
										goto l331
									l332:
										position, tokenIndex, depth = position332, tokenIndex332, depth332
									}
									depth--
									add(rulePegText, position330)
								}
								{
									position338, tokenIndex338, depth338 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l339
									}
									position++
									goto l338
								l339:
									position, tokenIndex, depth = position338, tokenIndex338, depth338
									if buffer[position] != rune('"') {
										goto l260
									}
									position++
								}
							l338:
								if buffer[position] != rune(')') {
									goto l260
								}
								position++
								{
									add(ruleAction34, position)
								}
								depth--
								add(ruleObjectID, position327)
							}
							break
						case '"':
							{
								position341 := position
								depth++
								if buffer[position] != rune('"') {
									goto l260
								}
								position++
								{
									position342 := position
									depth++
								l343:
									{
										position344, tokenIndex344, depth344 := position, tokenIndex, depth
										{
											position345 := position
											depth++
											{
												position346, tokenIndex346, depth346 := position, tokenIndex, depth
												{
													position348, tokenIndex348, depth348 := position, tokenIndex, depth
													{
														position349, tokenIndex349, depth349 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l350
														}
														position++
														goto l349
													l350:
														position, tokenIndex, depth = position349, tokenIndex349, depth349
														if buffer[position] != rune('\\') {
															goto l348
														}
														position++
													}
												l349:
													goto l347
												l348:
													position, tokenIndex, depth = position348, tokenIndex348, depth348
												}
												if !matchDot() {
													goto l347
												}
												goto l346
											l347:
												position, tokenIndex, depth = position346, tokenIndex346, depth346
												if buffer[position] != rune('\\') {
													goto l344
												}
												position++
												if !matchDot() {
													goto l344
												}
											}
										l346:
											depth--
											add(rulestringChar, position345)
										}
										goto l343
									l344:
										position, tokenIndex, depth = position344, tokenIndex344, depth344
									}
									depth--
									add(rulePegText, position342)
								}
								if buffer[position] != rune('"') {
									goto l260
								}
								position++
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleString, position341)
							}
							break
						case 'f', 't':
							{
								position352 := position
								depth++
								{
									position353, tokenIndex353, depth353 := position, tokenIndex, depth
									{
										position355 := position
										depth++
										if buffer[position] != rune('t') {
											goto l354
										}
										position++
										if buffer[position] != rune('r') {
											goto l354
										}
										position++
										if buffer[position] != rune('u') {
											goto l354
										}
										position++
										if buffer[position] != rune('e') {
											goto l354
										}
										position++
										{
											add(ruleAction31, position)
										}
										depth--
										add(ruleTrue, position355)
									}
									goto l353
								l354:
									position, tokenIndex, depth = position353, tokenIndex353, depth353
									{
										position357 := position
										depth++
										if buffer[position] != rune('f') {
											goto l260
										}
										position++
										if buffer[position] != rune('a') {
											goto l260
										}
										position++
										if buffer[position] != rune('l') {
											goto l260
										}
										position++
										if buffer[position] != rune('s') {
											goto l260
										}
										position++
										if buffer[position] != rune('e') {
											goto l260
										}
										position++
										{
											add(ruleAction32, position)
										}
										depth--
										add(ruleFalse, position357)
									}
								}
							l353:
								depth--
								add(ruleBoolean, position352)
							}
							break
						case '[':
							{
								position359 := position
								depth++
								if buffer[position] != rune('[') {
									goto l260
								}
								position++
								{
									add(ruleAction24, position)
								}
								{
									position361, tokenIndex361, depth361 := position, tokenIndex, depth
									{
										position363 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l361
										}
									l364:
										{
											position365, tokenIndex365, depth365 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l365
											}
											position++
											if !_rules[ruleListElem]() {
												goto l365
											}
											goto l364
										l365:
											position, tokenIndex, depth = position365, tokenIndex365, depth365
										}
										depth--
										add(ruleListElements, position363)
									}
									goto l362
								l361:
									position, tokenIndex, depth = position361, tokenIndex361, depth361
								}
							l362:
								if buffer[position] != rune(']') {
									goto l260
								}
								position++
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleList, position359)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l260
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l260
							}
							break
						}
					}

				}
			l262:
				depth--
				add(ruleValue, position261)
			}
			return true
		l260:
			position, tokenIndex, depth = position260, tokenIndex260, depth260
			return false
		},
		/* 48 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action28)> */
		func() bool {
			position367, tokenIndex367, depth367 := position, tokenIndex, depth
			{
				position368 := position
				depth++
				{
					position369 := position
					depth++
					{
						position370, tokenIndex370, depth370 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l370
						}
						position++
						goto l371
					l370:
						position, tokenIndex, depth = position370, tokenIndex370, depth370
					}
				l371:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l367
					}
					position++
				l372:
					{
						position373, tokenIndex373, depth373 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l373
						}
						position++
						goto l372
					l373:
						position, tokenIndex, depth = position373, tokenIndex373, depth373
					}
					{
						position374, tokenIndex374, depth374 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l374
						}
						position++
						goto l375
					l374:
						position, tokenIndex, depth = position374, tokenIndex374, depth374
					}
				l375:
				l376:
					{
						position377, tokenIndex377, depth377 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l377
						}
						position++
						goto l376
					l377:
						position, tokenIndex, depth = position377, tokenIndex377, depth377
					}
					depth--
					add(rulePegText, position369)
				}
				{
					add(ruleAction28, position)
				}
				depth--
				add(ruleNumeric, position368)
			}
			return true
		l367:
			position, tokenIndex, depth = position367, tokenIndex367, depth367
			return false
		},
		/* 49 Boolean <- <(True / False)> */
		nil,
		/* 50 String <- <('"' <stringChar*> '"' Action29)> */
		nil,
		/* 51 Null <- <('n' 'u' 'l' 'l' Action30)> */
		nil,
		/* 52 True <- <('t' 'r' 'u' 'e' Action31)> */
		nil,
		/* 53 False <- <('f' 'a' 'l' 's' 'e' Action32)> */
		nil,
		/* 54 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action33)> */
		nil,
		/* 55 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' ('\'' / '"') <hexChar*> ('\'' / '"') ')' Action34)> */
		nil,
		/* 56 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action35)> */
		nil,
		/* 57 Regex <- <('/' <regexBody> Action36)> */
		nil,
		/* 58 TimestampVal <- <(timestampParen / timestampPipe)> */
		nil,
		/* 59 timestampParen <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action37)> */
		nil,
		/* 60 timestampPipe <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' ' ' <([0-9] / '|')+> Action38)> */
		nil,
		/* 61 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action39)> */
		nil,
		/* 62 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action40)> */
		nil,
		/* 63 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action41)> */
		nil,
		/* 64 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action42)> */
		nil,
		/* 65 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 66 regexChar <- <(!'/' .)> */
		nil,
		/* 67 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 68 stringChar <- <((!('"' / '\\') .) / ('\\' .))> */
		nil,
		/* 69 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position399, tokenIndex399, depth399 := position, tokenIndex, depth
			{
				position400 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l399
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l399
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l399
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l399
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l399
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l399
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l399
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position400)
			}
			return true
		l399:
			position, tokenIndex, depth = position399, tokenIndex399, depth399
			return false
		},
		nil,
		/* 72 Action0 <- <{ p.SetField("log_level", buffer[begin:end]) }> */
		nil,
		/* 73 Action1 <- <{ p.SetField("component", buffer[begin:end]) }> */
		nil,
		/* 74 Action2 <- <{ p.SetField("context", buffer[begin:end]) }> */
		nil,
		/* 75 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 76 Action4 <- <{ p.SetField("warning", buffer[begin:end]) }> */
		nil,
		/* 77 Action5 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 78 Action6 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 79 Action7 <- <{ p.EndField() }> */
		nil,
		/* 80 Action8 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 81 Action9 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 82 Action10 <- <{ p.EndField() }> */
		nil,
		/* 83 Action11 <- <{ p.SetField("command_type", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 84 Action12 <- <{ p.EndField() }> */
		nil,
		/* 85 Action13 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 86 Action14 <- <{ p.EndField()}> */
		nil,
		/* 87 Action15 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 88 Action16 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 89 Action17 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 90 Action18 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 91 Action19 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 92 Action20 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 93 Action21 <- <{ p.PushMap() }> */
		nil,
		/* 94 Action22 <- <{ p.PopMap() }> */
		nil,
		/* 95 Action23 <- <{ p.SetMapValue() }> */
		nil,
		/* 96 Action24 <- <{ p.PushList() }> */
		nil,
		/* 97 Action25 <- <{ p.PopList() }> */
		nil,
		/* 98 Action26 <- <{ p.SetListValue() }> */
		nil,
		/* 99 Action27 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 100 Action28 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 101 Action29 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 102 Action30 <- <{ p.PushValue(nil) }> */
		nil,
		/* 103 Action31 <- <{ p.PushValue(true) }> */
		nil,
		/* 104 Action32 <- <{ p.PushValue(false) }> */
		nil,
		/* 105 Action33 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 106 Action34 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 107 Action35 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 108 Action36 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 109 Action37 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 110 Action38 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 111 Action39 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 112 Action40 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 113 Action41 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 114 Action42 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
