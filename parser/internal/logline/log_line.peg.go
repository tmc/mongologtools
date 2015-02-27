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
	rulePartialDoc
	rulepartialDoc
	rulepartialDocExtra
	ruleknownField
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
	ruleAction45

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
	"PartialDoc",
	"partialDoc",
	"partialDocExtra",
	"knownField",
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
	"Action45",

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
	rules  [123]func() bool
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
			p.PushValue(buffer[begin:end])
		case ruleAction21:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction22:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction23:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction24:
			p.PushMap()
		case ruleAction25:
			p.PopMap()
		case ruleAction26:
			p.SetMapValue()
		case ruleAction27:
			p.PushList()
		case ruleAction28:
			p.PopList()
		case ruleAction29:
			p.SetListValue()
		case ruleAction30:
			p.PushField(buffer[begin:end])
		case ruleAction31:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction32:
			p.PushValue(buffer[begin:end])
		case ruleAction33:
			p.PushValue(nil)
		case ruleAction34:
			p.PushValue(true)
		case ruleAction35:
			p.PushValue(false)
		case ruleAction36:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction37:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction38:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction39:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction40:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction41:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction42:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction43:
			p.PushValue(p.Minkey())
		case ruleAction44:
			p.PushValue(p.Maxkey())
		case ruleAction45:
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
								add(ruleAction21, position)
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
								add(ruleAction22, position)
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
								switch buffer[position] {
								case 'F':
									if buffer[position] != rune('F') {
										goto l30
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l30
									}
									position++
									break
								case 'W':
									if buffer[position] != rune('W') {
										goto l30
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l30
									}
									position++
									break
								default:
									if buffer[position] != rune('D') {
										goto l30
									}
									position++
									break
								}
							}

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
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					{
						position38 := position
						depth++
						{
							position39 := position
							depth++
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l36
							}
							position++
						l40:
							{
								position41, tokenIndex41, depth41 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l41
								}
								position++
								goto l40
							l41:
								position, tokenIndex, depth = position41, tokenIndex41, depth41
							}
							depth--
							add(rulePegText, position39)
						}
						if buffer[position] != rune(' ') {
							goto l36
						}
						position++
					l42:
						{
							position43, tokenIndex43, depth43 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l43
							}
							position++
							goto l42
						l43:
							position, tokenIndex, depth = position43, tokenIndex43, depth43
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleComponent, position38)
					}
					goto l37
				l36:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
				}
			l37:
				{
					position45 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position46 := position
						depth++
						{
							position49 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position51, tokenIndex51, depth51 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l52
										}
										position++
										goto l51
									l52:
										position, tokenIndex, depth = position51, tokenIndex51, depth51
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l51:
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
							add(ruleletterOrDigit, position49)
						}
					l47:
						{
							position48, tokenIndex48, depth48 := position, tokenIndex, depth
							{
								position53 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position55, tokenIndex55, depth55 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l56
											}
											position++
											goto l55
										l56:
											position, tokenIndex, depth = position55, tokenIndex55, depth55
											if buffer[position] != rune('$') {
												goto l48
											}
											position++
										}
									l55:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l48
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l48
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l48
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position53)
							}
							goto l47
						l48:
							position, tokenIndex, depth = position48, tokenIndex48, depth48
						}
						depth--
						add(rulePegText, position46)
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
					add(ruleContext, position45)
				}
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					{
						position60 := position
						depth++
						{
							position61 := position
							depth++
							{
								position62 := position
								depth++
								if buffer[position] != rune('w') {
									goto l58
								}
								position++
								if buffer[position] != rune('a') {
									goto l58
								}
								position++
								if buffer[position] != rune('r') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('g') {
									goto l58
								}
								position++
								if buffer[position] != rune(':') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('l') {
									goto l58
								}
								position++
								if buffer[position] != rune('o') {
									goto l58
								}
								position++
								if buffer[position] != rune('g') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('l') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('a') {
									goto l58
								}
								position++
								if buffer[position] != rune('t') {
									goto l58
								}
								position++
								if buffer[position] != rune('t') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune('m') {
									goto l58
								}
								position++
								if buffer[position] != rune('p') {
									goto l58
								}
								position++
								if buffer[position] != rune('t') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune('d') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('(') {
									goto l58
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l58
								}
								position++
							l63:
								{
									position64, tokenIndex64, depth64 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l64
									}
									position++
									goto l63
								l64:
									position, tokenIndex, depth = position64, tokenIndex64, depth64
								}
								if buffer[position] != rune('k') {
									goto l58
								}
								position++
								if buffer[position] != rune(')') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('o') {
									goto l58
								}
								position++
								if buffer[position] != rune('v') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune('r') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('m') {
									goto l58
								}
								position++
								if buffer[position] != rune('a') {
									goto l58
								}
								position++
								if buffer[position] != rune('x') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('s') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('z') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('(') {
									goto l58
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l58
								}
								position++
							l65:
								{
									position66, tokenIndex66, depth66 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l66
									}
									position++
									goto l65
								l66:
									position, tokenIndex, depth = position66, tokenIndex66, depth66
								}
								if buffer[position] != rune('k') {
									goto l58
								}
								position++
								if buffer[position] != rune(')') {
									goto l58
								}
								position++
								if buffer[position] != rune(',') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('p') {
									goto l58
								}
								position++
								if buffer[position] != rune('r') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('t') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('g') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('b') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune('g') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('i') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('g') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('a') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('d') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('e') {
									goto l58
								}
								position++
								if buffer[position] != rune('n') {
									goto l58
								}
								position++
								if buffer[position] != rune('d') {
									goto l58
								}
								position++
								if buffer[position] != rune(' ') {
									goto l58
								}
								position++
								if buffer[position] != rune('.') {
									goto l58
								}
								position++
								if buffer[position] != rune('.') {
									goto l58
								}
								position++
								if buffer[position] != rune('.') {
									goto l58
								}
								position++
								depth--
								add(ruleloglineSizeWarning, position62)
							}
							depth--
							add(rulePegText, position61)
						}
						if buffer[position] != rune(' ') {
							goto l58
						}
						position++
						{
							add(ruleAction4, position)
						}
						depth--
						add(ruleWarning, position60)
					}
					goto l59
				l58:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
				}
			l59:
				{
					position68 := position
					depth++
					{
						position69 := position
						depth++
						{
							position72, tokenIndex72, depth72 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l73
							}
							position++
							goto l72
						l73:
							position, tokenIndex, depth = position72, tokenIndex72, depth72
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l0
							}
							position++
						}
					l72:
					l70:
						{
							position71, tokenIndex71, depth71 := position, tokenIndex, depth
							{
								position74, tokenIndex74, depth74 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l75
								}
								position++
								goto l74
							l75:
								position, tokenIndex, depth = position74, tokenIndex74, depth74
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l71
								}
								position++
							}
						l74:
							goto l70
						l71:
							position, tokenIndex, depth = position71, tokenIndex71, depth71
						}
						depth--
						add(rulePegText, position69)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleOp, position68)
				}
				{
					position77 := position
					depth++
					{
						position78 := position
						depth++
					l79:
						{
							position80, tokenIndex80, depth80 := position, tokenIndex, depth
							{
								position81 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l80
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l80
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l80
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l80
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l80
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l80
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position81)
							}
							goto l79
						l80:
							position, tokenIndex, depth = position80, tokenIndex80, depth80
						}
						depth--
						add(rulePegText, position78)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction5, position)
					}
					depth--
					add(ruleNS, position77)
				}
			l84:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l85
					}
					goto l84
				l85:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
				}
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					{
						position88 := position
						depth++
						if buffer[position] != rune('l') {
							goto l86
						}
						position++
						if buffer[position] != rune('o') {
							goto l86
						}
						position++
						if buffer[position] != rune('c') {
							goto l86
						}
						position++
						if buffer[position] != rune('k') {
							goto l86
						}
						position++
						if buffer[position] != rune('s') {
							goto l86
						}
						position++
						if buffer[position] != rune('(') {
							goto l86
						}
						position++
						if buffer[position] != rune('m') {
							goto l86
						}
						position++
						if buffer[position] != rune('i') {
							goto l86
						}
						position++
						if buffer[position] != rune('c') {
							goto l86
						}
						position++
						if buffer[position] != rune('r') {
							goto l86
						}
						position++
						if buffer[position] != rune('o') {
							goto l86
						}
						position++
						if buffer[position] != rune('s') {
							goto l86
						}
						position++
						if buffer[position] != rune(')') {
							goto l86
						}
						position++
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l89
							}
							goto l90
						l89:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
						}
					l90:
					l91:
						{
							position92, tokenIndex92, depth92 := position, tokenIndex, depth
							{
								position93 := position
								depth++
								{
									position94 := position
									depth++
									{
										switch buffer[position] {
										case 'R':
											if buffer[position] != rune('R') {
												goto l92
											}
											position++
											break
										case 'r':
											if buffer[position] != rune('r') {
												goto l92
											}
											position++
											break
										default:
											{
												position96, tokenIndex96, depth96 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l97
												}
												position++
												goto l96
											l97:
												position, tokenIndex, depth = position96, tokenIndex96, depth96
												if buffer[position] != rune('W') {
													goto l92
												}
												position++
											}
										l96:
											break
										}
									}

									depth--
									add(rulePegText, position94)
								}
								{
									add(ruleAction6, position)
								}
								if buffer[position] != rune(':') {
									goto l92
								}
								position++
								if !_rules[ruleNumeric]() {
									goto l92
								}
								{
									position99, tokenIndex99, depth99 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l99
									}
									goto l100
								l99:
									position, tokenIndex, depth = position99, tokenIndex99, depth99
								}
							l100:
								{
									add(ruleAction7, position)
								}
								depth--
								add(rulelock, position93)
							}
							goto l91
						l92:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
						}
						depth--
						add(ruleLocks, position88)
					}
					goto l87
				l86:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
				}
			l87:
			l102:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l103
					}
					goto l102
				l103:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
				}
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					{
						position106 := position
						depth++
						{
							position107 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l104
							}
							position++
						l108:
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l109
								}
								position++
								goto l108
							l109:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
							}
							depth--
							add(rulePegText, position107)
						}
						if buffer[position] != rune('m') {
							goto l104
						}
						position++
						if buffer[position] != rune('s') {
							goto l104
						}
						position++
						{
							add(ruleAction8, position)
						}
						depth--
						add(ruleDuration, position106)
					}
					goto l105
				l104:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
				}
			l105:
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					{
						position113 := position
						depth++
						{
							position114 := position
							depth++
							if !matchDot() {
								goto l111
							}
						l115:
							{
								position116, tokenIndex116, depth116 := position, tokenIndex, depth
								if !matchDot() {
									goto l116
								}
								goto l115
							l116:
								position, tokenIndex, depth = position116, tokenIndex116, depth116
							}
							depth--
							add(rulePegText, position114)
						}
						{
							add(ruleAction23, position)
						}
						depth--
						add(ruleextra, position113)
					}
					goto l112
				l111:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
				}
			l112:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if !matchDot() {
						goto l118
					}
					goto l0
				l118:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
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
		/* 2 Severity <- <(<((&('F') 'F') | (&('E') 'E') | (&('W') 'W') | (&('I') 'I') | (&('D') 'D'))> ' ' Action0)> */
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
			position126, tokenIndex126, depth126 := position, tokenIndex, depth
			{
				position127 := position
				depth++
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					{
						position130 := position
						depth++
						if buffer[position] != rune('e') {
							goto l129
						}
						position++
						if buffer[position] != rune('x') {
							goto l129
						}
						position++
						if buffer[position] != rune('c') {
							goto l129
						}
						position++
						if buffer[position] != rune('e') {
							goto l129
						}
						position++
						if buffer[position] != rune('p') {
							goto l129
						}
						position++
						if buffer[position] != rune('t') {
							goto l129
						}
						position++
						if buffer[position] != rune('i') {
							goto l129
						}
						position++
						if buffer[position] != rune('o') {
							goto l129
						}
						position++
						if buffer[position] != rune('n') {
							goto l129
						}
						position++
						if buffer[position] != rune(':') {
							goto l129
						}
						position++
						{
							add(ruleAction18, position)
						}
						{
							position132 := position
							depth++
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if !matchDot() {
									goto l129
								}
								{
									position136, tokenIndex136, depth136 := position, tokenIndex, depth
									if buffer[position] != rune('c') {
										goto l136
									}
									position++
									if buffer[position] != rune('o') {
										goto l136
									}
									position++
									if buffer[position] != rune('d') {
										goto l136
									}
									position++
									if buffer[position] != rune('e') {
										goto l136
									}
									position++
									if buffer[position] != rune(':') {
										goto l136
									}
									position++
									goto l129
								l136:
									position, tokenIndex, depth = position136, tokenIndex136, depth136
								}
								position, tokenIndex, depth = position135, tokenIndex135, depth135
							}
							if !matchDot() {
								goto l129
							}
						l133:
							{
								position134, tokenIndex134, depth134 := position, tokenIndex, depth
								{
									position137, tokenIndex137, depth137 := position, tokenIndex, depth
									if !matchDot() {
										goto l134
									}
									{
										position138, tokenIndex138, depth138 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l138
										}
										position++
										if buffer[position] != rune('o') {
											goto l138
										}
										position++
										if buffer[position] != rune('d') {
											goto l138
										}
										position++
										if buffer[position] != rune('e') {
											goto l138
										}
										position++
										if buffer[position] != rune(':') {
											goto l138
										}
										position++
										goto l134
									l138:
										position, tokenIndex, depth = position138, tokenIndex138, depth138
									}
									position, tokenIndex, depth = position137, tokenIndex137, depth137
								}
								if !matchDot() {
									goto l134
								}
								goto l133
							l134:
								position, tokenIndex, depth = position134, tokenIndex134, depth134
							}
							depth--
							add(rulePegText, position132)
						}
						{
							position139, tokenIndex139, depth139 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l139
							}
							goto l140
						l139:
							position, tokenIndex, depth = position139, tokenIndex139, depth139
						}
					l140:
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleexceptionField, position130)
					}
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					{
						position143 := position
						depth++
						if buffer[position] != rune('c') {
							goto l142
						}
						position++
						if buffer[position] != rune('o') {
							goto l142
						}
						position++
						if buffer[position] != rune('m') {
							goto l142
						}
						position++
						if buffer[position] != rune('m') {
							goto l142
						}
						position++
						if buffer[position] != rune('a') {
							goto l142
						}
						position++
						if buffer[position] != rune('n') {
							goto l142
						}
						position++
						if buffer[position] != rune('d') {
							goto l142
						}
						position++
						if buffer[position] != rune(':') {
							goto l142
						}
						position++
						if buffer[position] != rune(' ') {
							goto l142
						}
						position++
						{
							position144 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l142
							}
						l145:
							{
								position146, tokenIndex146, depth146 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l146
								}
								goto l145
							l146:
								position, tokenIndex, depth = position146, tokenIndex146, depth146
							}
							depth--
							add(rulePegText, position144)
						}
						{
							position147, tokenIndex147, depth147 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l147
							}
							goto l148
						l147:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
						}
					l148:
						{
							add(ruleAction11, position)
						}
						if !_rules[ruleLineValue]() {
							goto l142
						}
						{
							add(ruleAction12, position)
						}
						depth--
						add(rulecommandField, position143)
					}
					goto l128
				l142:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					{
						position152 := position
						depth++
						if buffer[position] != rune('p') {
							goto l151
						}
						position++
						if buffer[position] != rune('l') {
							goto l151
						}
						position++
						if buffer[position] != rune('a') {
							goto l151
						}
						position++
						if buffer[position] != rune('n') {
							goto l151
						}
						position++
						if buffer[position] != rune('S') {
							goto l151
						}
						position++
						if buffer[position] != rune('u') {
							goto l151
						}
						position++
						if buffer[position] != rune('m') {
							goto l151
						}
						position++
						if buffer[position] != rune('m') {
							goto l151
						}
						position++
						if buffer[position] != rune('a') {
							goto l151
						}
						position++
						if buffer[position] != rune('r') {
							goto l151
						}
						position++
						if buffer[position] != rune('y') {
							goto l151
						}
						position++
						if buffer[position] != rune(':') {
							goto l151
						}
						position++
						if buffer[position] != rune(' ') {
							goto l151
						}
						position++
						{
							add(ruleAction13, position)
						}
						{
							position154 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l151
							}
						l155:
							{
								position156, tokenIndex156, depth156 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l156
								}
								position++
								if buffer[position] != rune(' ') {
									goto l156
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l156
								}
								goto l155
							l156:
								position, tokenIndex, depth = position156, tokenIndex156, depth156
							}
							depth--
							add(ruleplanSummaryElements, position154)
						}
						{
							add(ruleAction14, position)
						}
						depth--
						add(ruleplanSummaryField, position152)
					}
					goto l128
				l151:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					{
						position158 := position
						depth++
						{
							position159 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l126
							}
						l160:
							{
								position161, tokenIndex161, depth161 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l161
								}
								goto l160
							l161:
								position, tokenIndex, depth = position161, tokenIndex161, depth161
							}
							depth--
							add(rulePegText, position159)
						}
						if buffer[position] != rune(':') {
							goto l126
						}
						position++
						{
							position162, tokenIndex162, depth162 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l162
							}
							goto l163
						l162:
							position, tokenIndex, depth = position162, tokenIndex162, depth162
						}
					l163:
						{
							add(ruleAction9, position)
						}
						if !_rules[ruleLineValue]() {
							goto l126
						}
						{
							add(ruleAction10, position)
						}
						depth--
						add(ruleplainField, position158)
					}
				}
			l128:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l166
					}
					goto l167
				l166:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
				}
			l167:
				depth--
				add(ruleLineField, position127)
			}
			return true
		l126:
			position, tokenIndex, depth = position126, tokenIndex126, depth126
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
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				{
					position178 := position
					depth++
					{
						position179 := position
						depth++
						{
							position180, tokenIndex180, depth180 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l181
							}
							position++
							if buffer[position] != rune('N') {
								goto l181
							}
							position++
							if buffer[position] != rune('D') {
								goto l181
							}
							position++
							if buffer[position] != rune('_') {
								goto l181
							}
							position++
							if buffer[position] != rune('H') {
								goto l181
							}
							position++
							if buffer[position] != rune('A') {
								goto l181
							}
							position++
							if buffer[position] != rune('S') {
								goto l181
							}
							position++
							if buffer[position] != rune('H') {
								goto l181
							}
							position++
							goto l180
						l181:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('C') {
								goto l182
							}
							position++
							if buffer[position] != rune('A') {
								goto l182
							}
							position++
							if buffer[position] != rune('C') {
								goto l182
							}
							position++
							if buffer[position] != rune('H') {
								goto l182
							}
							position++
							if buffer[position] != rune('E') {
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
							if buffer[position] != rune('P') {
								goto l182
							}
							position++
							if buffer[position] != rune('L') {
								goto l182
							}
							position++
							if buffer[position] != rune('A') {
								goto l182
							}
							position++
							if buffer[position] != rune('N') {
								goto l182
							}
							position++
							goto l180
						l182:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('C') {
								goto l183
							}
							position++
							if buffer[position] != rune('O') {
								goto l183
							}
							position++
							if buffer[position] != rune('L') {
								goto l183
							}
							position++
							if buffer[position] != rune('L') {
								goto l183
							}
							position++
							if buffer[position] != rune('S') {
								goto l183
							}
							position++
							if buffer[position] != rune('C') {
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
							goto l180
						l183:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('C') {
								goto l184
							}
							position++
							if buffer[position] != rune('O') {
								goto l184
							}
							position++
							if buffer[position] != rune('U') {
								goto l184
							}
							position++
							if buffer[position] != rune('N') {
								goto l184
							}
							position++
							if buffer[position] != rune('T') {
								goto l184
							}
							position++
							goto l180
						l184:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('D') {
								goto l185
							}
							position++
							if buffer[position] != rune('E') {
								goto l185
							}
							position++
							if buffer[position] != rune('L') {
								goto l185
							}
							position++
							if buffer[position] != rune('E') {
								goto l185
							}
							position++
							if buffer[position] != rune('T') {
								goto l185
							}
							position++
							if buffer[position] != rune('E') {
								goto l185
							}
							position++
							goto l180
						l185:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('G') {
								goto l186
							}
							position++
							if buffer[position] != rune('E') {
								goto l186
							}
							position++
							if buffer[position] != rune('O') {
								goto l186
							}
							position++
							if buffer[position] != rune('_') {
								goto l186
							}
							position++
							if buffer[position] != rune('N') {
								goto l186
							}
							position++
							if buffer[position] != rune('E') {
								goto l186
							}
							position++
							if buffer[position] != rune('A') {
								goto l186
							}
							position++
							if buffer[position] != rune('R') {
								goto l186
							}
							position++
							if buffer[position] != rune('_') {
								goto l186
							}
							position++
							if buffer[position] != rune('2') {
								goto l186
							}
							position++
							if buffer[position] != rune('D') {
								goto l186
							}
							position++
							goto l180
						l186:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
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
							if buffer[position] != rune('S') {
								goto l187
							}
							position++
							if buffer[position] != rune('P') {
								goto l187
							}
							position++
							if buffer[position] != rune('H') {
								goto l187
							}
							position++
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
							if buffer[position] != rune('R') {
								goto l187
							}
							position++
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
							goto l180
						l187:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('I') {
								goto l188
							}
							position++
							if buffer[position] != rune('D') {
								goto l188
							}
							position++
							if buffer[position] != rune('H') {
								goto l188
							}
							position++
							if buffer[position] != rune('A') {
								goto l188
							}
							position++
							if buffer[position] != rune('C') {
								goto l188
							}
							position++
							if buffer[position] != rune('K') {
								goto l188
							}
							position++
							goto l180
						l188:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('S') {
								goto l189
							}
							position++
							if buffer[position] != rune('O') {
								goto l189
							}
							position++
							if buffer[position] != rune('R') {
								goto l189
							}
							position++
							if buffer[position] != rune('T') {
								goto l189
							}
							position++
							if buffer[position] != rune('_') {
								goto l189
							}
							position++
							if buffer[position] != rune('M') {
								goto l189
							}
							position++
							if buffer[position] != rune('E') {
								goto l189
							}
							position++
							if buffer[position] != rune('R') {
								goto l189
							}
							position++
							if buffer[position] != rune('G') {
								goto l189
							}
							position++
							if buffer[position] != rune('E') {
								goto l189
							}
							position++
							goto l180
						l189:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('S') {
								goto l190
							}
							position++
							if buffer[position] != rune('H') {
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
							if buffer[position] != rune('D') {
								goto l190
							}
							position++
							if buffer[position] != rune('I') {
								goto l190
							}
							position++
							if buffer[position] != rune('N') {
								goto l190
							}
							position++
							if buffer[position] != rune('G') {
								goto l190
							}
							position++
							if buffer[position] != rune('_') {
								goto l190
							}
							position++
							if buffer[position] != rune('F') {
								goto l190
							}
							position++
							if buffer[position] != rune('I') {
								goto l190
							}
							position++
							if buffer[position] != rune('L') {
								goto l190
							}
							position++
							if buffer[position] != rune('T') {
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
							goto l180
						l190:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('S') {
								goto l191
							}
							position++
							if buffer[position] != rune('K') {
								goto l191
							}
							position++
							if buffer[position] != rune('I') {
								goto l191
							}
							position++
							if buffer[position] != rune('P') {
								goto l191
							}
							position++
							goto l180
						l191:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
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
							goto l180
						l192:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('X') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('B') {
										goto l176
									}
									position++
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									if buffer[position] != rune('L') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									if buffer[position] != rune('_') {
										goto l176
									}
									position++
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									if buffer[position] != rune('R') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('J') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('R') {
										goto l176
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('L') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('_') {
										goto l176
									}
									position++
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									if buffer[position] != rune('L') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('M') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									if buffer[position] != rune('_') {
										goto l176
									}
									position++
									if buffer[position] != rune('M') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('X') {
										goto l176
									}
									position++
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l176
									}
									position++
									if buffer[position] != rune('R') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('P') {
										goto l176
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('H') {
										goto l176
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('F') {
										goto l176
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('I') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('U') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('_') {
										goto l176
									}
									position++
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									if buffer[position] != rune('C') {
										goto l176
									}
									position++
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l176
									}
									position++
									if buffer[position] != rune('N') {
										goto l176
									}
									position++
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									if buffer[position] != rune('_') {
										goto l176
									}
									position++
									if buffer[position] != rune('S') {
										goto l176
									}
									position++
									if buffer[position] != rune('O') {
										goto l176
									}
									position++
									if buffer[position] != rune('R') {
										goto l176
									}
									position++
									if buffer[position] != rune('T') {
										goto l176
									}
									position++
									if buffer[position] != rune('E') {
										goto l176
									}
									position++
									if buffer[position] != rune('D') {
										goto l176
									}
									position++
									break
								}
							}

						}
					l180:
						depth--
						add(ruleplanSummaryStage, position179)
					}
					depth--
					add(rulePegText, position178)
				}
				{
					add(ruleAction15, position)
				}
				{
					position195 := position
					depth++
					{
						position196, tokenIndex196, depth196 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l197
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l197
						}
						{
							add(ruleAction16, position)
						}
						goto l196
					l197:
						position, tokenIndex, depth = position196, tokenIndex196, depth196
						{
							add(ruleAction17, position)
						}
					}
				l196:
					depth--
					add(ruleplanSummary, position195)
				}
				depth--
				add(ruleplanSummaryElem, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 18 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 19 planSummary <- <((' ' LineValue Action16) / Action17)> */
		nil,
		/* 20 exceptionField <- <('e' 'x' 'c' 'e' 'p' 't' 'i' 'o' 'n' ':' Action18 <(&(. !('c' 'o' 'd' 'e' ':')) .)+> S? Action19)> */
		nil,
		/* 21 LineValue <- <((Doc / Numeric / PartialDoc) S?)> */
		func() bool {
			position203, tokenIndex203, depth203 := position, tokenIndex, depth
			{
				position204 := position
				depth++
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l206
					}
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if !_rules[ruleNumeric]() {
						goto l207
					}
					goto l205
				l207:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					{
						position208 := position
						depth++
						{
							position209 := position
							depth++
							{
								position210 := position
								depth++
								if buffer[position] != rune('{') {
									goto l203
								}
								position++
								{
									position213, tokenIndex213, depth213 := position, tokenIndex, depth
									if buffer[position] != rune('}') {
										goto l213
									}
									position++
									goto l203
								l213:
									position, tokenIndex, depth = position213, tokenIndex213, depth213
								}
								if !matchDot() {
									goto l203
								}
							l211:
								{
									position212, tokenIndex212, depth212 := position, tokenIndex, depth
									{
										position214, tokenIndex214, depth214 := position, tokenIndex, depth
										if buffer[position] != rune('}') {
											goto l214
										}
										position++
										goto l212
									l214:
										position, tokenIndex, depth = position214, tokenIndex214, depth214
									}
									if !matchDot() {
										goto l212
									}
									goto l211
								l212:
									position, tokenIndex, depth = position212, tokenIndex212, depth212
								}
								if buffer[position] != rune('}') {
									goto l203
								}
								position++
							l215:
								{
									position216, tokenIndex216, depth216 := position, tokenIndex, depth
									{
										position217 := position
										depth++
										{
											position218, tokenIndex218, depth218 := position, tokenIndex, depth
											if !matchDot() {
												goto l216
											}
											{
												position219, tokenIndex219, depth219 := position, tokenIndex, depth
												{
													position220 := position
													depth++
													{
														position221, tokenIndex221, depth221 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l222
														}
														position++
														if buffer[position] != rune('i') {
															goto l222
														}
														position++
														if buffer[position] != rune('n') {
															goto l222
														}
														position++
														if buffer[position] != rune('s') {
															goto l222
														}
														position++
														if buffer[position] != rune('e') {
															goto l222
														}
														position++
														if buffer[position] != rune('r') {
															goto l222
														}
														position++
														if buffer[position] != rune('t') {
															goto l222
														}
														position++
														if buffer[position] != rune('e') {
															goto l222
														}
														position++
														if buffer[position] != rune('d') {
															goto l222
														}
														position++
														goto l221
													l222:
														position, tokenIndex, depth = position221, tokenIndex221, depth221
														{
															switch buffer[position] {
															case 'n':
																if buffer[position] != rune('n') {
																	goto l219
																}
																position++
																if buffer[position] != rune('t') {
																	goto l219
																}
																position++
																if buffer[position] != rune('o') {
																	goto l219
																}
																position++
																if buffer[position] != rune('r') {
																	goto l219
																}
																position++
																if buffer[position] != rune('e') {
																	goto l219
																}
																position++
																if buffer[position] != rune('t') {
																	goto l219
																}
																position++
																if buffer[position] != rune('u') {
																	goto l219
																}
																position++
																if buffer[position] != rune('r') {
																	goto l219
																}
																position++
																if buffer[position] != rune('n') {
																	goto l219
																}
																position++
																break
															case 'c':
																if buffer[position] != rune('c') {
																	goto l219
																}
																position++
																if buffer[position] != rune('u') {
																	goto l219
																}
																position++
																if buffer[position] != rune('r') {
																	goto l219
																}
																position++
																if buffer[position] != rune('s') {
																	goto l219
																}
																position++
																if buffer[position] != rune('o') {
																	goto l219
																}
																position++
																if buffer[position] != rune('r') {
																	goto l219
																}
																position++
																if buffer[position] != rune('i') {
																	goto l219
																}
																position++
																if buffer[position] != rune('d') {
																	goto l219
																}
																position++
																break
															default:
																if buffer[position] != rune('p') {
																	goto l219
																}
																position++
																if buffer[position] != rune('l') {
																	goto l219
																}
																position++
																if buffer[position] != rune('a') {
																	goto l219
																}
																position++
																if buffer[position] != rune('n') {
																	goto l219
																}
																position++
																if buffer[position] != rune('S') {
																	goto l219
																}
																position++
																if buffer[position] != rune('u') {
																	goto l219
																}
																position++
																if buffer[position] != rune('m') {
																	goto l219
																}
																position++
																if buffer[position] != rune('m') {
																	goto l219
																}
																position++
																if buffer[position] != rune('a') {
																	goto l219
																}
																position++
																if buffer[position] != rune('r') {
																	goto l219
																}
																position++
																if buffer[position] != rune('y') {
																	goto l219
																}
																position++
																break
															}
														}

													}
												l221:
													depth--
													add(ruleknownField, position220)
												}
												goto l216
											l219:
												position, tokenIndex, depth = position219, tokenIndex219, depth219
											}
											position, tokenIndex, depth = position218, tokenIndex218, depth218
										}
										if !matchDot() {
											goto l216
										}
										depth--
										add(rulepartialDocExtra, position217)
									}
									goto l215
								l216:
									position, tokenIndex, depth = position216, tokenIndex216, depth216
								}
								depth--
								add(rulepartialDoc, position210)
							}
							depth--
							add(rulePegText, position209)
						}
						{
							add(ruleAction20, position)
						}
						depth--
						add(rulePartialDoc, position208)
					}
				}
			l205:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l225
					}
					goto l226
				l225:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
				}
			l226:
				depth--
				add(ruleLineValue, position204)
			}
			return true
		l203:
			position, tokenIndex, depth = position203, tokenIndex203, depth203
			return false
		},
		/* 22 PartialDoc <- <(<partialDoc> Action20)> */
		nil,
		/* 23 partialDoc <- <('{' (!'}' .)+ '}' partialDocExtra*)> */
		nil,
		/* 24 partialDocExtra <- <(&(. !knownField) .)> */
		nil,
		/* 25 knownField <- <(('n' 'i' 'n' 's' 'e' 'r' 't' 'e' 'd') / ((&('n') ('n' 't' 'o' 'r' 'e' 't' 'u' 'r' 'n')) | (&('c') ('c' 'u' 'r' 's' 'o' 'r' 'i' 'd')) | (&('p') ('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y'))))> */
		nil,
		/* 26 timestamp24 <- <(<(date ' ' time)> Action21)> */
		nil,
		/* 27 timestamp26 <- <(<datetime26> Action22)> */
		nil,
		/* 28 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 29 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 30 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position235, tokenIndex235, depth235 := position, tokenIndex, depth
			{
				position236 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l235
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l235
				}
				position++
				depth--
				add(ruledigit2, position236)
			}
			return true
		l235:
			position, tokenIndex, depth = position235, tokenIndex235, depth235
			return false
		},
		/* 31 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 32 tz <- <('+' [0-9]+)> */
		nil,
		/* 33 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position239, tokenIndex239, depth239 := position, tokenIndex, depth
			{
				position240 := position
				depth++
				{
					position241 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l239
					}
					depth--
					add(rulehour, position241)
				}
				if buffer[position] != rune(':') {
					goto l239
				}
				position++
				{
					position242 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l239
					}
					depth--
					add(ruleminute, position242)
				}
				if buffer[position] != rune(':') {
					goto l239
				}
				position++
				{
					position243 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l239
					}
					depth--
					add(rulesecond, position243)
				}
				if buffer[position] != rune('.') {
					goto l239
				}
				position++
				{
					position244 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l239
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l239
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l239
					}
					position++
					depth--
					add(rulemillisecond, position244)
				}
				depth--
				add(ruletime, position240)
			}
			return true
		l239:
			position, tokenIndex, depth = position239, tokenIndex239, depth239
			return false
		},
		/* 34 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 35 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 36 dayNum <- <([0-9] [0-9]?)> */
		nil,
		/* 37 hour <- <digit2> */
		nil,
		/* 38 minute <- <digit2> */
		nil,
		/* 39 second <- <digit2> */
		nil,
		/* 40 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 41 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 42 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 43 extra <- <(<.+> Action23)> */
		nil,
		/* 44 S <- <' '+> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l255
				}
				position++
			l257:
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
				}
				depth--
				add(ruleS, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 45 Doc <- <('{' Action24 DocElements? '}' Action25)> */
		func() bool {
			position259, tokenIndex259, depth259 := position, tokenIndex, depth
			{
				position260 := position
				depth++
				if buffer[position] != rune('{') {
					goto l259
				}
				position++
				{
					add(ruleAction24, position)
				}
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					{
						position264 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l262
						}
					l265:
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l266
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l266
							}
							goto l265
						l266:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
						}
						depth--
						add(ruleDocElements, position264)
					}
					goto l263
				l262:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
				}
			l263:
				if buffer[position] != rune('}') {
					goto l259
				}
				position++
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleDoc, position260)
			}
			return true
		l259:
			position, tokenIndex, depth = position259, tokenIndex259, depth259
			return false
		},
		/* 46 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 47 DocElem <- <(S? Field S? Value S? Action26)> */
		func() bool {
			position269, tokenIndex269, depth269 := position, tokenIndex, depth
			{
				position270 := position
				depth++
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l271
					}
					goto l272
				l271:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
				}
			l272:
				{
					position273 := position
					depth++
					{
						position274 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l269
						}
					l275:
						{
							position276, tokenIndex276, depth276 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l276
							}
							goto l275
						l276:
							position, tokenIndex, depth = position276, tokenIndex276, depth276
						}
						depth--
						add(rulePegText, position274)
					}
					if buffer[position] != rune(':') {
						goto l269
					}
					position++
					{
						add(ruleAction30, position)
					}
					depth--
					add(ruleField, position273)
				}
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l278
					}
					goto l279
				l278:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
				}
			l279:
				if !_rules[ruleValue]() {
					goto l269
				}
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l280
					}
					goto l281
				l280:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
				}
			l281:
				{
					add(ruleAction26, position)
				}
				depth--
				add(ruleDocElem, position270)
			}
			return true
		l269:
			position, tokenIndex, depth = position269, tokenIndex269, depth269
			return false
		},
		/* 48 List <- <('[' Action27 ListElements? ']' Action28)> */
		nil,
		/* 49 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 50 ListElem <- <(S? Value S? Action29)> */
		func() bool {
			position285, tokenIndex285, depth285 := position, tokenIndex, depth
			{
				position286 := position
				depth++
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l287
					}
					goto l288
				l287:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
				}
			l288:
				if !_rules[ruleValue]() {
					goto l285
				}
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l289
					}
					goto l290
				l289:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
				}
			l290:
				{
					add(ruleAction29, position)
				}
				depth--
				add(ruleListElem, position286)
			}
			return true
		l285:
			position, tokenIndex, depth = position285, tokenIndex285, depth285
			return false
		},
		/* 51 Field <- <(<fieldChar+> ':' Action30)> */
		nil,
		/* 52 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position293, tokenIndex293, depth293 := position, tokenIndex, depth
			{
				position294 := position
				depth++
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					{
						position297 := position
						depth++
						if buffer[position] != rune('n') {
							goto l296
						}
						position++
						if buffer[position] != rune('u') {
							goto l296
						}
						position++
						if buffer[position] != rune('l') {
							goto l296
						}
						position++
						if buffer[position] != rune('l') {
							goto l296
						}
						position++
						{
							add(ruleAction33, position)
						}
						depth--
						add(ruleNull, position297)
					}
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					{
						position300 := position
						depth++
						if buffer[position] != rune('M') {
							goto l299
						}
						position++
						if buffer[position] != rune('i') {
							goto l299
						}
						position++
						if buffer[position] != rune('n') {
							goto l299
						}
						position++
						if buffer[position] != rune('K') {
							goto l299
						}
						position++
						if buffer[position] != rune('e') {
							goto l299
						}
						position++
						if buffer[position] != rune('y') {
							goto l299
						}
						position++
						{
							add(ruleAction43, position)
						}
						depth--
						add(ruleMinKey, position300)
					}
					goto l295
				l299:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					{
						switch buffer[position] {
						case 'M':
							{
								position303 := position
								depth++
								if buffer[position] != rune('M') {
									goto l293
								}
								position++
								if buffer[position] != rune('a') {
									goto l293
								}
								position++
								if buffer[position] != rune('x') {
									goto l293
								}
								position++
								if buffer[position] != rune('K') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('y') {
									goto l293
								}
								position++
								{
									add(ruleAction44, position)
								}
								depth--
								add(ruleMaxKey, position303)
							}
							break
						case 'u':
							{
								position305 := position
								depth++
								if buffer[position] != rune('u') {
									goto l293
								}
								position++
								if buffer[position] != rune('n') {
									goto l293
								}
								position++
								if buffer[position] != rune('d') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('f') {
									goto l293
								}
								position++
								if buffer[position] != rune('i') {
									goto l293
								}
								position++
								if buffer[position] != rune('n') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('d') {
									goto l293
								}
								position++
								{
									add(ruleAction45, position)
								}
								depth--
								add(ruleUndefined, position305)
							}
							break
						case 'N':
							{
								position307 := position
								depth++
								if buffer[position] != rune('N') {
									goto l293
								}
								position++
								if buffer[position] != rune('u') {
									goto l293
								}
								position++
								if buffer[position] != rune('m') {
									goto l293
								}
								position++
								if buffer[position] != rune('b') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('r') {
									goto l293
								}
								position++
								if buffer[position] != rune('L') {
									goto l293
								}
								position++
								if buffer[position] != rune('o') {
									goto l293
								}
								position++
								if buffer[position] != rune('n') {
									goto l293
								}
								position++
								if buffer[position] != rune('g') {
									goto l293
								}
								position++
								if buffer[position] != rune('(') {
									goto l293
								}
								position++
								{
									position308 := position
									depth++
									{
										position311, tokenIndex311, depth311 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l311
										}
										position++
										goto l293
									l311:
										position, tokenIndex, depth = position311, tokenIndex311, depth311
									}
									if !matchDot() {
										goto l293
									}
								l309:
									{
										position310, tokenIndex310, depth310 := position, tokenIndex, depth
										{
											position312, tokenIndex312, depth312 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l312
											}
											position++
											goto l310
										l312:
											position, tokenIndex, depth = position312, tokenIndex312, depth312
										}
										if !matchDot() {
											goto l310
										}
										goto l309
									l310:
										position, tokenIndex, depth = position310, tokenIndex310, depth310
									}
									depth--
									add(rulePegText, position308)
								}
								if buffer[position] != rune(')') {
									goto l293
								}
								position++
								{
									add(ruleAction42, position)
								}
								depth--
								add(ruleNumberLong, position307)
							}
							break
						case '/':
							{
								position314 := position
								depth++
								if buffer[position] != rune('/') {
									goto l293
								}
								position++
								{
									position315 := position
									depth++
									{
										position316 := position
										depth++
										{
											position319 := position
											depth++
											{
												position320, tokenIndex320, depth320 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l320
												}
												position++
												goto l293
											l320:
												position, tokenIndex, depth = position320, tokenIndex320, depth320
											}
											if !matchDot() {
												goto l293
											}
											depth--
											add(ruleregexChar, position319)
										}
									l317:
										{
											position318, tokenIndex318, depth318 := position, tokenIndex, depth
											{
												position321 := position
												depth++
												{
													position322, tokenIndex322, depth322 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l322
													}
													position++
													goto l318
												l322:
													position, tokenIndex, depth = position322, tokenIndex322, depth322
												}
												if !matchDot() {
													goto l318
												}
												depth--
												add(ruleregexChar, position321)
											}
											goto l317
										l318:
											position, tokenIndex, depth = position318, tokenIndex318, depth318
										}
										if buffer[position] != rune('/') {
											goto l293
										}
										position++
									l323:
										{
											position324, tokenIndex324, depth324 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l324
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l324
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l324
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l324
													}
													position++
													break
												}
											}

											goto l323
										l324:
											position, tokenIndex, depth = position324, tokenIndex324, depth324
										}
										depth--
										add(ruleregexBody, position316)
									}
									depth--
									add(rulePegText, position315)
								}
								{
									add(ruleAction39, position)
								}
								depth--
								add(ruleRegex, position314)
							}
							break
						case 'T':
							{
								position327 := position
								depth++
								{
									position328, tokenIndex328, depth328 := position, tokenIndex, depth
									{
										position330 := position
										depth++
										if buffer[position] != rune('T') {
											goto l329
										}
										position++
										if buffer[position] != rune('i') {
											goto l329
										}
										position++
										if buffer[position] != rune('m') {
											goto l329
										}
										position++
										if buffer[position] != rune('e') {
											goto l329
										}
										position++
										if buffer[position] != rune('s') {
											goto l329
										}
										position++
										if buffer[position] != rune('t') {
											goto l329
										}
										position++
										if buffer[position] != rune('a') {
											goto l329
										}
										position++
										if buffer[position] != rune('m') {
											goto l329
										}
										position++
										if buffer[position] != rune('p') {
											goto l329
										}
										position++
										if buffer[position] != rune('(') {
											goto l329
										}
										position++
										{
											position331 := position
											depth++
											{
												position334, tokenIndex334, depth334 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l334
												}
												position++
												goto l329
											l334:
												position, tokenIndex, depth = position334, tokenIndex334, depth334
											}
											if !matchDot() {
												goto l329
											}
										l332:
											{
												position333, tokenIndex333, depth333 := position, tokenIndex, depth
												{
													position335, tokenIndex335, depth335 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l335
													}
													position++
													goto l333
												l335:
													position, tokenIndex, depth = position335, tokenIndex335, depth335
												}
												if !matchDot() {
													goto l333
												}
												goto l332
											l333:
												position, tokenIndex, depth = position333, tokenIndex333, depth333
											}
											depth--
											add(rulePegText, position331)
										}
										if buffer[position] != rune(')') {
											goto l329
										}
										position++
										{
											add(ruleAction40, position)
										}
										depth--
										add(ruletimestampParen, position330)
									}
									goto l328
								l329:
									position, tokenIndex, depth = position328, tokenIndex328, depth328
									{
										position337 := position
										depth++
										if buffer[position] != rune('T') {
											goto l293
										}
										position++
										if buffer[position] != rune('i') {
											goto l293
										}
										position++
										if buffer[position] != rune('m') {
											goto l293
										}
										position++
										if buffer[position] != rune('e') {
											goto l293
										}
										position++
										if buffer[position] != rune('s') {
											goto l293
										}
										position++
										if buffer[position] != rune('t') {
											goto l293
										}
										position++
										if buffer[position] != rune('a') {
											goto l293
										}
										position++
										if buffer[position] != rune('m') {
											goto l293
										}
										position++
										if buffer[position] != rune('p') {
											goto l293
										}
										position++
										if buffer[position] != rune(' ') {
											goto l293
										}
										position++
										{
											position338 := position
											depth++
											{
												position341, tokenIndex341, depth341 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l342
												}
												position++
												goto l341
											l342:
												position, tokenIndex, depth = position341, tokenIndex341, depth341
												if buffer[position] != rune('|') {
													goto l293
												}
												position++
											}
										l341:
										l339:
											{
												position340, tokenIndex340, depth340 := position, tokenIndex, depth
												{
													position343, tokenIndex343, depth343 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l344
													}
													position++
													goto l343
												l344:
													position, tokenIndex, depth = position343, tokenIndex343, depth343
													if buffer[position] != rune('|') {
														goto l340
													}
													position++
												}
											l343:
												goto l339
											l340:
												position, tokenIndex, depth = position340, tokenIndex340, depth340
											}
											depth--
											add(rulePegText, position338)
										}
										{
											add(ruleAction41, position)
										}
										depth--
										add(ruletimestampPipe, position337)
									}
								}
							l328:
								depth--
								add(ruleTimestampVal, position327)
							}
							break
						case 'B':
							{
								position346 := position
								depth++
								if buffer[position] != rune('B') {
									goto l293
								}
								position++
								if buffer[position] != rune('i') {
									goto l293
								}
								position++
								if buffer[position] != rune('n') {
									goto l293
								}
								position++
								if buffer[position] != rune('D') {
									goto l293
								}
								position++
								if buffer[position] != rune('a') {
									goto l293
								}
								position++
								if buffer[position] != rune('t') {
									goto l293
								}
								position++
								if buffer[position] != rune('a') {
									goto l293
								}
								position++
								if buffer[position] != rune('(') {
									goto l293
								}
								position++
								{
									position347 := position
									depth++
									{
										position350, tokenIndex350, depth350 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l350
										}
										position++
										goto l293
									l350:
										position, tokenIndex, depth = position350, tokenIndex350, depth350
									}
									if !matchDot() {
										goto l293
									}
								l348:
									{
										position349, tokenIndex349, depth349 := position, tokenIndex, depth
										{
											position351, tokenIndex351, depth351 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l351
											}
											position++
											goto l349
										l351:
											position, tokenIndex, depth = position351, tokenIndex351, depth351
										}
										if !matchDot() {
											goto l349
										}
										goto l348
									l349:
										position, tokenIndex, depth = position349, tokenIndex349, depth349
									}
									depth--
									add(rulePegText, position347)
								}
								if buffer[position] != rune(')') {
									goto l293
								}
								position++
								{
									add(ruleAction38, position)
								}
								depth--
								add(ruleBinData, position346)
							}
							break
						case 'D', 'n':
							{
								position353 := position
								depth++
								{
									position354, tokenIndex354, depth354 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l354
									}
									position++
									if buffer[position] != rune('e') {
										goto l354
									}
									position++
									if buffer[position] != rune('w') {
										goto l354
									}
									position++
									if buffer[position] != rune(' ') {
										goto l354
									}
									position++
									goto l355
								l354:
									position, tokenIndex, depth = position354, tokenIndex354, depth354
								}
							l355:
								if buffer[position] != rune('D') {
									goto l293
								}
								position++
								if buffer[position] != rune('a') {
									goto l293
								}
								position++
								if buffer[position] != rune('t') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('(') {
									goto l293
								}
								position++
								{
									position356, tokenIndex356, depth356 := position, tokenIndex, depth
									if buffer[position] != rune('-') {
										goto l356
									}
									position++
									goto l357
								l356:
									position, tokenIndex, depth = position356, tokenIndex356, depth356
								}
							l357:
								{
									position358 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l293
									}
									position++
								l359:
									{
										position360, tokenIndex360, depth360 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l360
										}
										position++
										goto l359
									l360:
										position, tokenIndex, depth = position360, tokenIndex360, depth360
									}
									depth--
									add(rulePegText, position358)
								}
								if buffer[position] != rune(')') {
									goto l293
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleDate, position353)
							}
							break
						case 'O':
							{
								position362 := position
								depth++
								if buffer[position] != rune('O') {
									goto l293
								}
								position++
								if buffer[position] != rune('b') {
									goto l293
								}
								position++
								if buffer[position] != rune('j') {
									goto l293
								}
								position++
								if buffer[position] != rune('e') {
									goto l293
								}
								position++
								if buffer[position] != rune('c') {
									goto l293
								}
								position++
								if buffer[position] != rune('t') {
									goto l293
								}
								position++
								if buffer[position] != rune('I') {
									goto l293
								}
								position++
								if buffer[position] != rune('d') {
									goto l293
								}
								position++
								if buffer[position] != rune('(') {
									goto l293
								}
								position++
								{
									position363, tokenIndex363, depth363 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l364
									}
									position++
									goto l363
								l364:
									position, tokenIndex, depth = position363, tokenIndex363, depth363
									if buffer[position] != rune('"') {
										goto l293
									}
									position++
								}
							l363:
								{
									position365 := position
									depth++
								l366:
									{
										position367, tokenIndex367, depth367 := position, tokenIndex, depth
										{
											position368 := position
											depth++
											{
												position369, tokenIndex369, depth369 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l370
												}
												position++
												goto l369
											l370:
												position, tokenIndex, depth = position369, tokenIndex369, depth369
												{
													position371, tokenIndex371, depth371 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l372
													}
													position++
													goto l371
												l372:
													position, tokenIndex, depth = position371, tokenIndex371, depth371
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l367
													}
													position++
												}
											l371:
											}
										l369:
											depth--
											add(rulehexChar, position368)
										}
										goto l366
									l367:
										position, tokenIndex, depth = position367, tokenIndex367, depth367
									}
									depth--
									add(rulePegText, position365)
								}
								{
									position373, tokenIndex373, depth373 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l374
									}
									position++
									goto l373
								l374:
									position, tokenIndex, depth = position373, tokenIndex373, depth373
									if buffer[position] != rune('"') {
										goto l293
									}
									position++
								}
							l373:
								if buffer[position] != rune(')') {
									goto l293
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleObjectID, position362)
							}
							break
						case '"':
							{
								position376 := position
								depth++
								if buffer[position] != rune('"') {
									goto l293
								}
								position++
								{
									position377 := position
									depth++
								l378:
									{
										position379, tokenIndex379, depth379 := position, tokenIndex, depth
										{
											position380 := position
											depth++
											{
												position381, tokenIndex381, depth381 := position, tokenIndex, depth
												{
													position383, tokenIndex383, depth383 := position, tokenIndex, depth
													{
														position384, tokenIndex384, depth384 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l385
														}
														position++
														goto l384
													l385:
														position, tokenIndex, depth = position384, tokenIndex384, depth384
														if buffer[position] != rune('\\') {
															goto l383
														}
														position++
													}
												l384:
													goto l382
												l383:
													position, tokenIndex, depth = position383, tokenIndex383, depth383
												}
												if !matchDot() {
													goto l382
												}
												goto l381
											l382:
												position, tokenIndex, depth = position381, tokenIndex381, depth381
												if buffer[position] != rune('\\') {
													goto l379
												}
												position++
												if !matchDot() {
													goto l379
												}
											}
										l381:
											depth--
											add(rulestringChar, position380)
										}
										goto l378
									l379:
										position, tokenIndex, depth = position379, tokenIndex379, depth379
									}
									depth--
									add(rulePegText, position377)
								}
								if buffer[position] != rune('"') {
									goto l293
								}
								position++
								{
									add(ruleAction32, position)
								}
								depth--
								add(ruleString, position376)
							}
							break
						case 'f', 't':
							{
								position387 := position
								depth++
								{
									position388, tokenIndex388, depth388 := position, tokenIndex, depth
									{
										position390 := position
										depth++
										if buffer[position] != rune('t') {
											goto l389
										}
										position++
										if buffer[position] != rune('r') {
											goto l389
										}
										position++
										if buffer[position] != rune('u') {
											goto l389
										}
										position++
										if buffer[position] != rune('e') {
											goto l389
										}
										position++
										{
											add(ruleAction34, position)
										}
										depth--
										add(ruleTrue, position390)
									}
									goto l388
								l389:
									position, tokenIndex, depth = position388, tokenIndex388, depth388
									{
										position392 := position
										depth++
										if buffer[position] != rune('f') {
											goto l293
										}
										position++
										if buffer[position] != rune('a') {
											goto l293
										}
										position++
										if buffer[position] != rune('l') {
											goto l293
										}
										position++
										if buffer[position] != rune('s') {
											goto l293
										}
										position++
										if buffer[position] != rune('e') {
											goto l293
										}
										position++
										{
											add(ruleAction35, position)
										}
										depth--
										add(ruleFalse, position392)
									}
								}
							l388:
								depth--
								add(ruleBoolean, position387)
							}
							break
						case '[':
							{
								position394 := position
								depth++
								if buffer[position] != rune('[') {
									goto l293
								}
								position++
								{
									add(ruleAction27, position)
								}
								{
									position396, tokenIndex396, depth396 := position, tokenIndex, depth
									{
										position398 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l396
										}
									l399:
										{
											position400, tokenIndex400, depth400 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l400
											}
											position++
											if !_rules[ruleListElem]() {
												goto l400
											}
											goto l399
										l400:
											position, tokenIndex, depth = position400, tokenIndex400, depth400
										}
										depth--
										add(ruleListElements, position398)
									}
									goto l397
								l396:
									position, tokenIndex, depth = position396, tokenIndex396, depth396
								}
							l397:
								if buffer[position] != rune(']') {
									goto l293
								}
								position++
								{
									add(ruleAction28, position)
								}
								depth--
								add(ruleList, position394)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l293
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l293
							}
							break
						}
					}

				}
			l295:
				depth--
				add(ruleValue, position294)
			}
			return true
		l293:
			position, tokenIndex, depth = position293, tokenIndex293, depth293
			return false
		},
		/* 53 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action31)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404 := position
					depth++
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l405
						}
						position++
						goto l406
					l405:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
					}
				l406:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l402
					}
					position++
				l407:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l408
						}
						position++
						goto l407
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
					{
						position409, tokenIndex409, depth409 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l409
						}
						position++
						goto l410
					l409:
						position, tokenIndex, depth = position409, tokenIndex409, depth409
					}
				l410:
				l411:
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l412
						}
						position++
						goto l411
					l412:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
					}
					depth--
					add(rulePegText, position404)
				}
				{
					add(ruleAction31, position)
				}
				depth--
				add(ruleNumeric, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 54 Boolean <- <(True / False)> */
		nil,
		/* 55 String <- <('"' <stringChar*> '"' Action32)> */
		nil,
		/* 56 Null <- <('n' 'u' 'l' 'l' Action33)> */
		nil,
		/* 57 True <- <('t' 'r' 'u' 'e' Action34)> */
		nil,
		/* 58 False <- <('f' 'a' 'l' 's' 'e' Action35)> */
		nil,
		/* 59 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') '-'? <[0-9]+> ')' Action36)> */
		nil,
		/* 60 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' ('\'' / '"') <hexChar*> ('\'' / '"') ')' Action37)> */
		nil,
		/* 61 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action38)> */
		nil,
		/* 62 Regex <- <('/' <regexBody> Action39)> */
		nil,
		/* 63 TimestampVal <- <(timestampParen / timestampPipe)> */
		nil,
		/* 64 timestampParen <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action40)> */
		nil,
		/* 65 timestampPipe <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' ' ' <([0-9] / '|')+> Action41)> */
		nil,
		/* 66 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action42)> */
		nil,
		/* 67 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action43)> */
		nil,
		/* 68 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action44)> */
		nil,
		/* 69 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action45)> */
		nil,
		/* 70 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 71 regexChar <- <(!'/' .)> */
		nil,
		/* 72 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 73 stringChar <- <((!('"' / '\\') .) / ('\\' .))> */
		nil,
		/* 74 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l434
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l434
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l434
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l434
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l434
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l434
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l434
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		nil,
		/* 77 Action0 <- <{ p.SetField("severity", buffer[begin:end]) }> */
		nil,
		/* 78 Action1 <- <{ p.SetField("component", buffer[begin:end]) }> */
		nil,
		/* 79 Action2 <- <{ p.SetField("context", buffer[begin:end]) }> */
		nil,
		/* 80 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 81 Action4 <- <{ p.SetField("warning", buffer[begin:end]) }> */
		nil,
		/* 82 Action5 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 83 Action6 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 84 Action7 <- <{ p.EndField() }> */
		nil,
		/* 85 Action8 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 86 Action9 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 87 Action10 <- <{ p.EndField() }> */
		nil,
		/* 88 Action11 <- <{ p.SetField("command_type", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 89 Action12 <- <{ p.EndField() }> */
		nil,
		/* 90 Action13 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 91 Action14 <- <{ p.EndField()}> */
		nil,
		/* 92 Action15 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 93 Action16 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 94 Action17 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 95 Action18 <- <{ p.StartField("exception") }> */
		nil,
		/* 96 Action19 <- <{ p.PushValue(buffer[begin:end]); p.EndField() }> */
		nil,
		/* 97 Action20 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 98 Action21 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 99 Action22 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 100 Action23 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 101 Action24 <- <{ p.PushMap() }> */
		nil,
		/* 102 Action25 <- <{ p.PopMap() }> */
		nil,
		/* 103 Action26 <- <{ p.SetMapValue() }> */
		nil,
		/* 104 Action27 <- <{ p.PushList() }> */
		nil,
		/* 105 Action28 <- <{ p.PopList() }> */
		nil,
		/* 106 Action29 <- <{ p.SetListValue() }> */
		nil,
		/* 107 Action30 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 108 Action31 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 109 Action32 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 110 Action33 <- <{ p.PushValue(nil) }> */
		nil,
		/* 111 Action34 <- <{ p.PushValue(true) }> */
		nil,
		/* 112 Action35 <- <{ p.PushValue(false) }> */
		nil,
		/* 113 Action36 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 114 Action37 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 115 Action38 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 116 Action39 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 117 Action40 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 118 Action41 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 119 Action42 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 120 Action43 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 121 Action44 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 122 Action45 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
