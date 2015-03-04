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
	ruleOrderedDoc
	ruleOrderedDocElements
	ruleOrderedDocElem
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
	ruleAction46
	ruleAction47
	ruleAction48
	ruleAction49
	ruleAction50

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
	"OrderedDoc",
	"OrderedDocElements",
	"OrderedDocElem",
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
	"Action46",
	"Action47",
	"Action48",
	"Action49",
	"Action50",

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
	rules  [131]func() bool
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
			p.PushList()
		case ruleAction19:
			p.PopList()
		case ruleAction20:
			p.PushMap()
		case ruleAction21:
			p.SetMapValue()
			p.SetListValue()
		case ruleAction22:
			p.PopMap()
		case ruleAction23:
			p.StartField("exception")
		case ruleAction24:
			p.PushValue(buffer[begin:end])
			p.EndField()
		case ruleAction25:
			p.PushValue(buffer[begin:end])
		case ruleAction26:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction27:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction28:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction29:
			p.PushMap()
		case ruleAction30:
			p.PopMap()
		case ruleAction31:
			p.SetMapValue()
		case ruleAction32:
			p.PushList()
		case ruleAction33:
			p.PopList()
		case ruleAction34:
			p.SetListValue()
		case ruleAction35:
			p.PushField(buffer[begin:end])
		case ruleAction36:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction37:
			p.PushValue(buffer[begin:end])
		case ruleAction38:
			p.PushValue(nil)
		case ruleAction39:
			p.PushValue(true)
		case ruleAction40:
			p.PushValue(false)
		case ruleAction41:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction42:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction43:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction44:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction45:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction46:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction47:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction48:
			p.PushValue(p.Minkey())
		case ruleAction49:
			p.PushValue(p.Maxkey())
		case ruleAction50:
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
								add(ruleAction26, position)
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
								add(ruleAction27, position)
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
							add(ruleAction28, position)
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
							add(ruleAction23, position)
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
							add(ruleAction24, position)
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
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l176
								}
								position++
								break
							case '_':
								if buffer[position] != rune('_') {
									goto l176
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l176
								}
								position++
								break
							}
						}

					l180:
						{
							position181, tokenIndex181, depth181 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l181
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l181
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l181
									}
									position++
									break
								}
							}

							goto l180
						l181:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
						}
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
					position185 := position
					depth++
					{
						position186, tokenIndex186, depth186 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l187
						}
						position++
						{
							position188 := position
							depth++
							if buffer[position] != rune('{') {
								goto l187
							}
							position++
							{
								add(ruleAction18, position)
							}
							{
								position190, tokenIndex190, depth190 := position, tokenIndex, depth
								{
									position192 := position
									depth++
									if !_rules[ruleOrderedDocElem]() {
										goto l190
									}
								l193:
									{
										position194, tokenIndex194, depth194 := position, tokenIndex, depth
										if buffer[position] != rune(',') {
											goto l194
										}
										position++
										if !_rules[ruleOrderedDocElem]() {
											goto l194
										}
										goto l193
									l194:
										position, tokenIndex, depth = position194, tokenIndex194, depth194
									}
									depth--
									add(ruleOrderedDocElements, position192)
								}
								goto l191
							l190:
								position, tokenIndex, depth = position190, tokenIndex190, depth190
							}
						l191:
							if buffer[position] != rune('}') {
								goto l187
							}
							position++
							{
								add(ruleAction19, position)
							}
							depth--
							add(ruleOrderedDoc, position188)
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
				add(ruleplanSummaryElem, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 18 planSummaryStage <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))+> */
		nil,
		/* 19 planSummary <- <((' ' OrderedDoc Action16) / Action17)> */
		nil,
		/* 20 OrderedDoc <- <('{' Action18 OrderedDocElements? '}' Action19)> */
		nil,
		/* 21 OrderedDocElements <- <(OrderedDocElem (',' OrderedDocElem)*)> */
		nil,
		/* 22 OrderedDocElem <- <(S? Action20 Field S? Value Action21 S? Action22)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l204
					}
					goto l205
				l204:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
				}
			l205:
				{
					add(ruleAction20, position)
				}
				if !_rules[ruleField]() {
					goto l202
				}
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l207
					}
					goto l208
				l207:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
				}
			l208:
				if !_rules[ruleValue]() {
					goto l202
				}
				{
					add(ruleAction21, position)
				}
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
				{
					add(ruleAction22, position)
				}
				depth--
				add(ruleOrderedDocElem, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 23 exceptionField <- <('e' 'x' 'c' 'e' 'p' 't' 'i' 'o' 'n' ':' Action23 <(&(. !('c' 'o' 'd' 'e' ':')) .)+> S? Action24)> */
		nil,
		/* 24 LineValue <- <((Doc / Numeric / PartialDoc) S?)> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l217
					}
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if !_rules[ruleNumeric]() {
						goto l218
					}
					goto l216
				l218:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					{
						position219 := position
						depth++
						{
							position220 := position
							depth++
							{
								position221 := position
								depth++
								if buffer[position] != rune('{') {
									goto l214
								}
								position++
								{
									position224, tokenIndex224, depth224 := position, tokenIndex, depth
									if buffer[position] != rune('}') {
										goto l224
									}
									position++
									goto l214
								l224:
									position, tokenIndex, depth = position224, tokenIndex224, depth224
								}
								if !matchDot() {
									goto l214
								}
							l222:
								{
									position223, tokenIndex223, depth223 := position, tokenIndex, depth
									{
										position225, tokenIndex225, depth225 := position, tokenIndex, depth
										if buffer[position] != rune('}') {
											goto l225
										}
										position++
										goto l223
									l225:
										position, tokenIndex, depth = position225, tokenIndex225, depth225
									}
									if !matchDot() {
										goto l223
									}
									goto l222
								l223:
									position, tokenIndex, depth = position223, tokenIndex223, depth223
								}
								if buffer[position] != rune('}') {
									goto l214
								}
								position++
							l226:
								{
									position227, tokenIndex227, depth227 := position, tokenIndex, depth
									{
										position228 := position
										depth++
										{
											position229, tokenIndex229, depth229 := position, tokenIndex, depth
											if !matchDot() {
												goto l227
											}
											{
												position230, tokenIndex230, depth230 := position, tokenIndex, depth
												{
													position231 := position
													depth++
													{
														position232, tokenIndex232, depth232 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l233
														}
														position++
														if buffer[position] != rune('i') {
															goto l233
														}
														position++
														if buffer[position] != rune('n') {
															goto l233
														}
														position++
														if buffer[position] != rune('s') {
															goto l233
														}
														position++
														if buffer[position] != rune('e') {
															goto l233
														}
														position++
														if buffer[position] != rune('r') {
															goto l233
														}
														position++
														if buffer[position] != rune('t') {
															goto l233
														}
														position++
														if buffer[position] != rune('e') {
															goto l233
														}
														position++
														if buffer[position] != rune('d') {
															goto l233
														}
														position++
														goto l232
													l233:
														position, tokenIndex, depth = position232, tokenIndex232, depth232
														{
															switch buffer[position] {
															case 'n':
																if buffer[position] != rune('n') {
																	goto l230
																}
																position++
																if buffer[position] != rune('t') {
																	goto l230
																}
																position++
																if buffer[position] != rune('o') {
																	goto l230
																}
																position++
																if buffer[position] != rune('r') {
																	goto l230
																}
																position++
																if buffer[position] != rune('e') {
																	goto l230
																}
																position++
																if buffer[position] != rune('t') {
																	goto l230
																}
																position++
																if buffer[position] != rune('u') {
																	goto l230
																}
																position++
																if buffer[position] != rune('r') {
																	goto l230
																}
																position++
																if buffer[position] != rune('n') {
																	goto l230
																}
																position++
																break
															case 'c':
																if buffer[position] != rune('c') {
																	goto l230
																}
																position++
																if buffer[position] != rune('u') {
																	goto l230
																}
																position++
																if buffer[position] != rune('r') {
																	goto l230
																}
																position++
																if buffer[position] != rune('s') {
																	goto l230
																}
																position++
																if buffer[position] != rune('o') {
																	goto l230
																}
																position++
																if buffer[position] != rune('r') {
																	goto l230
																}
																position++
																if buffer[position] != rune('i') {
																	goto l230
																}
																position++
																if buffer[position] != rune('d') {
																	goto l230
																}
																position++
																break
															default:
																if buffer[position] != rune('p') {
																	goto l230
																}
																position++
																if buffer[position] != rune('l') {
																	goto l230
																}
																position++
																if buffer[position] != rune('a') {
																	goto l230
																}
																position++
																if buffer[position] != rune('n') {
																	goto l230
																}
																position++
																if buffer[position] != rune('S') {
																	goto l230
																}
																position++
																if buffer[position] != rune('u') {
																	goto l230
																}
																position++
																if buffer[position] != rune('m') {
																	goto l230
																}
																position++
																if buffer[position] != rune('m') {
																	goto l230
																}
																position++
																if buffer[position] != rune('a') {
																	goto l230
																}
																position++
																if buffer[position] != rune('r') {
																	goto l230
																}
																position++
																if buffer[position] != rune('y') {
																	goto l230
																}
																position++
																break
															}
														}

													}
												l232:
													depth--
													add(ruleknownField, position231)
												}
												goto l227
											l230:
												position, tokenIndex, depth = position230, tokenIndex230, depth230
											}
											position, tokenIndex, depth = position229, tokenIndex229, depth229
										}
										if !matchDot() {
											goto l227
										}
										depth--
										add(rulepartialDocExtra, position228)
									}
									goto l226
								l227:
									position, tokenIndex, depth = position227, tokenIndex227, depth227
								}
								depth--
								add(rulepartialDoc, position221)
							}
							depth--
							add(rulePegText, position220)
						}
						{
							add(ruleAction25, position)
						}
						depth--
						add(rulePartialDoc, position219)
					}
				}
			l216:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l236
					}
					goto l237
				l236:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
				}
			l237:
				depth--
				add(ruleLineValue, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
			return false
		},
		/* 25 PartialDoc <- <(<partialDoc> Action25)> */
		nil,
		/* 26 partialDoc <- <('{' (!'}' .)+ '}' partialDocExtra*)> */
		nil,
		/* 27 partialDocExtra <- <(&(. !knownField) .)> */
		nil,
		/* 28 knownField <- <(('n' 'i' 'n' 's' 'e' 'r' 't' 'e' 'd') / ((&('n') ('n' 't' 'o' 'r' 'e' 't' 'u' 'r' 'n')) | (&('c') ('c' 'u' 'r' 's' 'o' 'r' 'i' 'd')) | (&('p') ('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y'))))> */
		nil,
		/* 29 timestamp24 <- <(<(date ' ' time)> Action26)> */
		nil,
		/* 30 timestamp26 <- <(<datetime26> Action27)> */
		nil,
		/* 31 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 32 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 33 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position246, tokenIndex246, depth246 := position, tokenIndex, depth
			{
				position247 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l246
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l246
				}
				position++
				depth--
				add(ruledigit2, position247)
			}
			return true
		l246:
			position, tokenIndex, depth = position246, tokenIndex246, depth246
			return false
		},
		/* 34 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 35 tz <- <('+' [0-9]+)> */
		nil,
		/* 36 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l250
					}
					depth--
					add(rulehour, position252)
				}
				if buffer[position] != rune(':') {
					goto l250
				}
				position++
				{
					position253 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l250
					}
					depth--
					add(ruleminute, position253)
				}
				if buffer[position] != rune(':') {
					goto l250
				}
				position++
				{
					position254 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l250
					}
					depth--
					add(rulesecond, position254)
				}
				if buffer[position] != rune('.') {
					goto l250
				}
				position++
				{
					position255 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l250
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l250
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l250
					}
					position++
					depth--
					add(rulemillisecond, position255)
				}
				depth--
				add(ruletime, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 37 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 38 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 39 dayNum <- <([0-9] [0-9]?)> */
		nil,
		/* 40 hour <- <digit2> */
		nil,
		/* 41 minute <- <digit2> */
		nil,
		/* 42 second <- <digit2> */
		nil,
		/* 43 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 44 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 45 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 46 extra <- <(<.+> Action28)> */
		nil,
		/* 47 S <- <' '+> */
		func() bool {
			position266, tokenIndex266, depth266 := position, tokenIndex, depth
			{
				position267 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l266
				}
				position++
			l268:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
				}
				depth--
				add(ruleS, position267)
			}
			return true
		l266:
			position, tokenIndex, depth = position266, tokenIndex266, depth266
			return false
		},
		/* 48 Doc <- <('{' Action29 DocElements? '}' Action30)> */
		func() bool {
			position270, tokenIndex270, depth270 := position, tokenIndex, depth
			{
				position271 := position
				depth++
				if buffer[position] != rune('{') {
					goto l270
				}
				position++
				{
					add(ruleAction29, position)
				}
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					{
						position275 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l273
						}
					l276:
						{
							position277, tokenIndex277, depth277 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l277
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l277
							}
							goto l276
						l277:
							position, tokenIndex, depth = position277, tokenIndex277, depth277
						}
						depth--
						add(ruleDocElements, position275)
					}
					goto l274
				l273:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
				}
			l274:
				if buffer[position] != rune('}') {
					goto l270
				}
				position++
				{
					add(ruleAction30, position)
				}
				depth--
				add(ruleDoc, position271)
			}
			return true
		l270:
			position, tokenIndex, depth = position270, tokenIndex270, depth270
			return false
		},
		/* 49 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 50 DocElem <- <(S? Field S? Value S? Action31)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l282
					}
					goto l283
				l282:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
				}
			l283:
				if !_rules[ruleField]() {
					goto l280
				}
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l284
					}
					goto l285
				l284:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
				}
			l285:
				if !_rules[ruleValue]() {
					goto l280
				}
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l286
					}
					goto l287
				l286:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
				}
			l287:
				{
					add(ruleAction31, position)
				}
				depth--
				add(ruleDocElem, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 51 List <- <('[' Action32 ListElements? ']' Action33)> */
		nil,
		/* 52 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 53 ListElem <- <(S? Value S? Action34)> */
		func() bool {
			position291, tokenIndex291, depth291 := position, tokenIndex, depth
			{
				position292 := position
				depth++
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l293
					}
					goto l294
				l293:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
				}
			l294:
				if !_rules[ruleValue]() {
					goto l291
				}
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l295
					}
					goto l296
				l295:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
				}
			l296:
				{
					add(ruleAction34, position)
				}
				depth--
				add(ruleListElem, position292)
			}
			return true
		l291:
			position, tokenIndex, depth = position291, tokenIndex291, depth291
			return false
		},
		/* 54 Field <- <(<fieldChar+> ':' Action35)> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300 := position
					depth++
					if !_rules[rulefieldChar]() {
						goto l298
					}
				l301:
					{
						position302, tokenIndex302, depth302 := position, tokenIndex, depth
						if !_rules[rulefieldChar]() {
							goto l302
						}
						goto l301
					l302:
						position, tokenIndex, depth = position302, tokenIndex302, depth302
					}
					depth--
					add(rulePegText, position300)
				}
				if buffer[position] != rune(':') {
					goto l298
				}
				position++
				{
					add(ruleAction35, position)
				}
				depth--
				add(ruleField, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 55 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position304, tokenIndex304, depth304 := position, tokenIndex, depth
			{
				position305 := position
				depth++
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					{
						position308 := position
						depth++
						if buffer[position] != rune('n') {
							goto l307
						}
						position++
						if buffer[position] != rune('u') {
							goto l307
						}
						position++
						if buffer[position] != rune('l') {
							goto l307
						}
						position++
						if buffer[position] != rune('l') {
							goto l307
						}
						position++
						{
							add(ruleAction38, position)
						}
						depth--
						add(ruleNull, position308)
					}
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					{
						position311 := position
						depth++
						if buffer[position] != rune('M') {
							goto l310
						}
						position++
						if buffer[position] != rune('i') {
							goto l310
						}
						position++
						if buffer[position] != rune('n') {
							goto l310
						}
						position++
						if buffer[position] != rune('K') {
							goto l310
						}
						position++
						if buffer[position] != rune('e') {
							goto l310
						}
						position++
						if buffer[position] != rune('y') {
							goto l310
						}
						position++
						{
							add(ruleAction48, position)
						}
						depth--
						add(ruleMinKey, position311)
					}
					goto l306
				l310:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					{
						switch buffer[position] {
						case 'M':
							{
								position314 := position
								depth++
								if buffer[position] != rune('M') {
									goto l304
								}
								position++
								if buffer[position] != rune('a') {
									goto l304
								}
								position++
								if buffer[position] != rune('x') {
									goto l304
								}
								position++
								if buffer[position] != rune('K') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('y') {
									goto l304
								}
								position++
								{
									add(ruleAction49, position)
								}
								depth--
								add(ruleMaxKey, position314)
							}
							break
						case 'u':
							{
								position316 := position
								depth++
								if buffer[position] != rune('u') {
									goto l304
								}
								position++
								if buffer[position] != rune('n') {
									goto l304
								}
								position++
								if buffer[position] != rune('d') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('f') {
									goto l304
								}
								position++
								if buffer[position] != rune('i') {
									goto l304
								}
								position++
								if buffer[position] != rune('n') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('d') {
									goto l304
								}
								position++
								{
									add(ruleAction50, position)
								}
								depth--
								add(ruleUndefined, position316)
							}
							break
						case 'N':
							{
								position318 := position
								depth++
								if buffer[position] != rune('N') {
									goto l304
								}
								position++
								if buffer[position] != rune('u') {
									goto l304
								}
								position++
								if buffer[position] != rune('m') {
									goto l304
								}
								position++
								if buffer[position] != rune('b') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('r') {
									goto l304
								}
								position++
								if buffer[position] != rune('L') {
									goto l304
								}
								position++
								if buffer[position] != rune('o') {
									goto l304
								}
								position++
								if buffer[position] != rune('n') {
									goto l304
								}
								position++
								if buffer[position] != rune('g') {
									goto l304
								}
								position++
								if buffer[position] != rune('(') {
									goto l304
								}
								position++
								{
									position319 := position
									depth++
									{
										position322, tokenIndex322, depth322 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l322
										}
										position++
										goto l304
									l322:
										position, tokenIndex, depth = position322, tokenIndex322, depth322
									}
									if !matchDot() {
										goto l304
									}
								l320:
									{
										position321, tokenIndex321, depth321 := position, tokenIndex, depth
										{
											position323, tokenIndex323, depth323 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l323
											}
											position++
											goto l321
										l323:
											position, tokenIndex, depth = position323, tokenIndex323, depth323
										}
										if !matchDot() {
											goto l321
										}
										goto l320
									l321:
										position, tokenIndex, depth = position321, tokenIndex321, depth321
									}
									depth--
									add(rulePegText, position319)
								}
								if buffer[position] != rune(')') {
									goto l304
								}
								position++
								{
									add(ruleAction47, position)
								}
								depth--
								add(ruleNumberLong, position318)
							}
							break
						case '/':
							{
								position325 := position
								depth++
								if buffer[position] != rune('/') {
									goto l304
								}
								position++
								{
									position326 := position
									depth++
									{
										position327 := position
										depth++
										{
											position330 := position
											depth++
											{
												position331, tokenIndex331, depth331 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l331
												}
												position++
												goto l304
											l331:
												position, tokenIndex, depth = position331, tokenIndex331, depth331
											}
											if !matchDot() {
												goto l304
											}
											depth--
											add(ruleregexChar, position330)
										}
									l328:
										{
											position329, tokenIndex329, depth329 := position, tokenIndex, depth
											{
												position332 := position
												depth++
												{
													position333, tokenIndex333, depth333 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l333
													}
													position++
													goto l329
												l333:
													position, tokenIndex, depth = position333, tokenIndex333, depth333
												}
												if !matchDot() {
													goto l329
												}
												depth--
												add(ruleregexChar, position332)
											}
											goto l328
										l329:
											position, tokenIndex, depth = position329, tokenIndex329, depth329
										}
										if buffer[position] != rune('/') {
											goto l304
										}
										position++
									l334:
										{
											position335, tokenIndex335, depth335 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l335
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l335
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l335
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l335
													}
													position++
													break
												}
											}

											goto l334
										l335:
											position, tokenIndex, depth = position335, tokenIndex335, depth335
										}
										depth--
										add(ruleregexBody, position327)
									}
									depth--
									add(rulePegText, position326)
								}
								{
									add(ruleAction44, position)
								}
								depth--
								add(ruleRegex, position325)
							}
							break
						case 'T':
							{
								position338 := position
								depth++
								{
									position339, tokenIndex339, depth339 := position, tokenIndex, depth
									{
										position341 := position
										depth++
										if buffer[position] != rune('T') {
											goto l340
										}
										position++
										if buffer[position] != rune('i') {
											goto l340
										}
										position++
										if buffer[position] != rune('m') {
											goto l340
										}
										position++
										if buffer[position] != rune('e') {
											goto l340
										}
										position++
										if buffer[position] != rune('s') {
											goto l340
										}
										position++
										if buffer[position] != rune('t') {
											goto l340
										}
										position++
										if buffer[position] != rune('a') {
											goto l340
										}
										position++
										if buffer[position] != rune('m') {
											goto l340
										}
										position++
										if buffer[position] != rune('p') {
											goto l340
										}
										position++
										if buffer[position] != rune('(') {
											goto l340
										}
										position++
										{
											position342 := position
											depth++
											{
												position345, tokenIndex345, depth345 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l345
												}
												position++
												goto l340
											l345:
												position, tokenIndex, depth = position345, tokenIndex345, depth345
											}
											if !matchDot() {
												goto l340
											}
										l343:
											{
												position344, tokenIndex344, depth344 := position, tokenIndex, depth
												{
													position346, tokenIndex346, depth346 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l346
													}
													position++
													goto l344
												l346:
													position, tokenIndex, depth = position346, tokenIndex346, depth346
												}
												if !matchDot() {
													goto l344
												}
												goto l343
											l344:
												position, tokenIndex, depth = position344, tokenIndex344, depth344
											}
											depth--
											add(rulePegText, position342)
										}
										if buffer[position] != rune(')') {
											goto l340
										}
										position++
										{
											add(ruleAction45, position)
										}
										depth--
										add(ruletimestampParen, position341)
									}
									goto l339
								l340:
									position, tokenIndex, depth = position339, tokenIndex339, depth339
									{
										position348 := position
										depth++
										if buffer[position] != rune('T') {
											goto l304
										}
										position++
										if buffer[position] != rune('i') {
											goto l304
										}
										position++
										if buffer[position] != rune('m') {
											goto l304
										}
										position++
										if buffer[position] != rune('e') {
											goto l304
										}
										position++
										if buffer[position] != rune('s') {
											goto l304
										}
										position++
										if buffer[position] != rune('t') {
											goto l304
										}
										position++
										if buffer[position] != rune('a') {
											goto l304
										}
										position++
										if buffer[position] != rune('m') {
											goto l304
										}
										position++
										if buffer[position] != rune('p') {
											goto l304
										}
										position++
										if buffer[position] != rune(' ') {
											goto l304
										}
										position++
										{
											position349 := position
											depth++
											{
												position352, tokenIndex352, depth352 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l353
												}
												position++
												goto l352
											l353:
												position, tokenIndex, depth = position352, tokenIndex352, depth352
												if buffer[position] != rune('|') {
													goto l304
												}
												position++
											}
										l352:
										l350:
											{
												position351, tokenIndex351, depth351 := position, tokenIndex, depth
												{
													position354, tokenIndex354, depth354 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l355
													}
													position++
													goto l354
												l355:
													position, tokenIndex, depth = position354, tokenIndex354, depth354
													if buffer[position] != rune('|') {
														goto l351
													}
													position++
												}
											l354:
												goto l350
											l351:
												position, tokenIndex, depth = position351, tokenIndex351, depth351
											}
											depth--
											add(rulePegText, position349)
										}
										{
											add(ruleAction46, position)
										}
										depth--
										add(ruletimestampPipe, position348)
									}
								}
							l339:
								depth--
								add(ruleTimestampVal, position338)
							}
							break
						case 'B':
							{
								position357 := position
								depth++
								if buffer[position] != rune('B') {
									goto l304
								}
								position++
								if buffer[position] != rune('i') {
									goto l304
								}
								position++
								if buffer[position] != rune('n') {
									goto l304
								}
								position++
								if buffer[position] != rune('D') {
									goto l304
								}
								position++
								if buffer[position] != rune('a') {
									goto l304
								}
								position++
								if buffer[position] != rune('t') {
									goto l304
								}
								position++
								if buffer[position] != rune('a') {
									goto l304
								}
								position++
								if buffer[position] != rune('(') {
									goto l304
								}
								position++
								{
									position358 := position
									depth++
									{
										position361, tokenIndex361, depth361 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l361
										}
										position++
										goto l304
									l361:
										position, tokenIndex, depth = position361, tokenIndex361, depth361
									}
									if !matchDot() {
										goto l304
									}
								l359:
									{
										position360, tokenIndex360, depth360 := position, tokenIndex, depth
										{
											position362, tokenIndex362, depth362 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l362
											}
											position++
											goto l360
										l362:
											position, tokenIndex, depth = position362, tokenIndex362, depth362
										}
										if !matchDot() {
											goto l360
										}
										goto l359
									l360:
										position, tokenIndex, depth = position360, tokenIndex360, depth360
									}
									depth--
									add(rulePegText, position358)
								}
								if buffer[position] != rune(')') {
									goto l304
								}
								position++
								{
									add(ruleAction43, position)
								}
								depth--
								add(ruleBinData, position357)
							}
							break
						case 'D', 'n':
							{
								position364 := position
								depth++
								{
									position365, tokenIndex365, depth365 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l365
									}
									position++
									if buffer[position] != rune('e') {
										goto l365
									}
									position++
									if buffer[position] != rune('w') {
										goto l365
									}
									position++
									if buffer[position] != rune(' ') {
										goto l365
									}
									position++
									goto l366
								l365:
									position, tokenIndex, depth = position365, tokenIndex365, depth365
								}
							l366:
								if buffer[position] != rune('D') {
									goto l304
								}
								position++
								if buffer[position] != rune('a') {
									goto l304
								}
								position++
								if buffer[position] != rune('t') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('(') {
									goto l304
								}
								position++
								{
									position367, tokenIndex367, depth367 := position, tokenIndex, depth
									if buffer[position] != rune('-') {
										goto l367
									}
									position++
									goto l368
								l367:
									position, tokenIndex, depth = position367, tokenIndex367, depth367
								}
							l368:
								{
									position369 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l304
									}
									position++
								l370:
									{
										position371, tokenIndex371, depth371 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l371
										}
										position++
										goto l370
									l371:
										position, tokenIndex, depth = position371, tokenIndex371, depth371
									}
									depth--
									add(rulePegText, position369)
								}
								if buffer[position] != rune(')') {
									goto l304
								}
								position++
								{
									add(ruleAction41, position)
								}
								depth--
								add(ruleDate, position364)
							}
							break
						case 'O':
							{
								position373 := position
								depth++
								if buffer[position] != rune('O') {
									goto l304
								}
								position++
								if buffer[position] != rune('b') {
									goto l304
								}
								position++
								if buffer[position] != rune('j') {
									goto l304
								}
								position++
								if buffer[position] != rune('e') {
									goto l304
								}
								position++
								if buffer[position] != rune('c') {
									goto l304
								}
								position++
								if buffer[position] != rune('t') {
									goto l304
								}
								position++
								if buffer[position] != rune('I') {
									goto l304
								}
								position++
								if buffer[position] != rune('d') {
									goto l304
								}
								position++
								if buffer[position] != rune('(') {
									goto l304
								}
								position++
								{
									position374, tokenIndex374, depth374 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l375
									}
									position++
									goto l374
								l375:
									position, tokenIndex, depth = position374, tokenIndex374, depth374
									if buffer[position] != rune('"') {
										goto l304
									}
									position++
								}
							l374:
								{
									position376 := position
									depth++
								l377:
									{
										position378, tokenIndex378, depth378 := position, tokenIndex, depth
										{
											position379 := position
											depth++
											{
												position380, tokenIndex380, depth380 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l381
												}
												position++
												goto l380
											l381:
												position, tokenIndex, depth = position380, tokenIndex380, depth380
												{
													position382, tokenIndex382, depth382 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l383
													}
													position++
													goto l382
												l383:
													position, tokenIndex, depth = position382, tokenIndex382, depth382
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l378
													}
													position++
												}
											l382:
											}
										l380:
											depth--
											add(rulehexChar, position379)
										}
										goto l377
									l378:
										position, tokenIndex, depth = position378, tokenIndex378, depth378
									}
									depth--
									add(rulePegText, position376)
								}
								{
									position384, tokenIndex384, depth384 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l385
									}
									position++
									goto l384
								l385:
									position, tokenIndex, depth = position384, tokenIndex384, depth384
									if buffer[position] != rune('"') {
										goto l304
									}
									position++
								}
							l384:
								if buffer[position] != rune(')') {
									goto l304
								}
								position++
								{
									add(ruleAction42, position)
								}
								depth--
								add(ruleObjectID, position373)
							}
							break
						case '"':
							{
								position387 := position
								depth++
								if buffer[position] != rune('"') {
									goto l304
								}
								position++
								{
									position388 := position
									depth++
								l389:
									{
										position390, tokenIndex390, depth390 := position, tokenIndex, depth
										{
											position391 := position
											depth++
											{
												position392, tokenIndex392, depth392 := position, tokenIndex, depth
												{
													position394, tokenIndex394, depth394 := position, tokenIndex, depth
													{
														position395, tokenIndex395, depth395 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l396
														}
														position++
														goto l395
													l396:
														position, tokenIndex, depth = position395, tokenIndex395, depth395
														if buffer[position] != rune('\\') {
															goto l394
														}
														position++
													}
												l395:
													goto l393
												l394:
													position, tokenIndex, depth = position394, tokenIndex394, depth394
												}
												if !matchDot() {
													goto l393
												}
												goto l392
											l393:
												position, tokenIndex, depth = position392, tokenIndex392, depth392
												if buffer[position] != rune('\\') {
													goto l390
												}
												position++
												if !matchDot() {
													goto l390
												}
											}
										l392:
											depth--
											add(rulestringChar, position391)
										}
										goto l389
									l390:
										position, tokenIndex, depth = position390, tokenIndex390, depth390
									}
									depth--
									add(rulePegText, position388)
								}
								if buffer[position] != rune('"') {
									goto l304
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleString, position387)
							}
							break
						case 'f', 't':
							{
								position398 := position
								depth++
								{
									position399, tokenIndex399, depth399 := position, tokenIndex, depth
									{
										position401 := position
										depth++
										if buffer[position] != rune('t') {
											goto l400
										}
										position++
										if buffer[position] != rune('r') {
											goto l400
										}
										position++
										if buffer[position] != rune('u') {
											goto l400
										}
										position++
										if buffer[position] != rune('e') {
											goto l400
										}
										position++
										{
											add(ruleAction39, position)
										}
										depth--
										add(ruleTrue, position401)
									}
									goto l399
								l400:
									position, tokenIndex, depth = position399, tokenIndex399, depth399
									{
										position403 := position
										depth++
										if buffer[position] != rune('f') {
											goto l304
										}
										position++
										if buffer[position] != rune('a') {
											goto l304
										}
										position++
										if buffer[position] != rune('l') {
											goto l304
										}
										position++
										if buffer[position] != rune('s') {
											goto l304
										}
										position++
										if buffer[position] != rune('e') {
											goto l304
										}
										position++
										{
											add(ruleAction40, position)
										}
										depth--
										add(ruleFalse, position403)
									}
								}
							l399:
								depth--
								add(ruleBoolean, position398)
							}
							break
						case '[':
							{
								position405 := position
								depth++
								if buffer[position] != rune('[') {
									goto l304
								}
								position++
								{
									add(ruleAction32, position)
								}
								{
									position407, tokenIndex407, depth407 := position, tokenIndex, depth
									{
										position409 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l407
										}
									l410:
										{
											position411, tokenIndex411, depth411 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l411
											}
											position++
											if !_rules[ruleListElem]() {
												goto l411
											}
											goto l410
										l411:
											position, tokenIndex, depth = position411, tokenIndex411, depth411
										}
										depth--
										add(ruleListElements, position409)
									}
									goto l408
								l407:
									position, tokenIndex, depth = position407, tokenIndex407, depth407
								}
							l408:
								if buffer[position] != rune(']') {
									goto l304
								}
								position++
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleList, position405)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l304
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l304
							}
							break
						}
					}

				}
			l306:
				depth--
				add(ruleValue, position305)
			}
			return true
		l304:
			position, tokenIndex, depth = position304, tokenIndex304, depth304
			return false
		},
		/* 56 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action36)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415 := position
					depth++
					{
						position416, tokenIndex416, depth416 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l416
						}
						position++
						goto l417
					l416:
						position, tokenIndex, depth = position416, tokenIndex416, depth416
					}
				l417:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l413
					}
					position++
				l418:
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l419
						}
						position++
						goto l418
					l419:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
					}
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l420
						}
						position++
						goto l421
					l420:
						position, tokenIndex, depth = position420, tokenIndex420, depth420
					}
				l421:
				l422:
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l423
						}
						position++
						goto l422
					l423:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
					}
					depth--
					add(rulePegText, position415)
				}
				{
					add(ruleAction36, position)
				}
				depth--
				add(ruleNumeric, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 57 Boolean <- <(True / False)> */
		nil,
		/* 58 String <- <('"' <stringChar*> '"' Action37)> */
		nil,
		/* 59 Null <- <('n' 'u' 'l' 'l' Action38)> */
		nil,
		/* 60 True <- <('t' 'r' 'u' 'e' Action39)> */
		nil,
		/* 61 False <- <('f' 'a' 'l' 's' 'e' Action40)> */
		nil,
		/* 62 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') '-'? <[0-9]+> ')' Action41)> */
		nil,
		/* 63 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' ('\'' / '"') <hexChar*> ('\'' / '"') ')' Action42)> */
		nil,
		/* 64 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action43)> */
		nil,
		/* 65 Regex <- <('/' <regexBody> Action44)> */
		nil,
		/* 66 TimestampVal <- <(timestampParen / timestampPipe)> */
		nil,
		/* 67 timestampParen <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action45)> */
		nil,
		/* 68 timestampPipe <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' ' ' <([0-9] / '|')+> Action46)> */
		nil,
		/* 69 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action47)> */
		nil,
		/* 70 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action48)> */
		nil,
		/* 71 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action49)> */
		nil,
		/* 72 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action50)> */
		nil,
		/* 73 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 74 regexChar <- <(!'/' .)> */
		nil,
		/* 75 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 76 stringChar <- <((!('"' / '\\') .) / ('\\' .))> */
		nil,
		/* 77 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l445
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l445
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l445
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l445
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l445
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l445
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l445
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		nil,
		/* 80 Action0 <- <{ p.SetField("severity", buffer[begin:end]) }> */
		nil,
		/* 81 Action1 <- <{ p.SetField("component", buffer[begin:end]) }> */
		nil,
		/* 82 Action2 <- <{ p.SetField("context", buffer[begin:end]) }> */
		nil,
		/* 83 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 84 Action4 <- <{ p.SetField("warning", buffer[begin:end]) }> */
		nil,
		/* 85 Action5 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 86 Action6 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 87 Action7 <- <{ p.EndField() }> */
		nil,
		/* 88 Action8 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 89 Action9 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 90 Action10 <- <{ p.EndField() }> */
		nil,
		/* 91 Action11 <- <{ p.SetField("command_type", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 92 Action12 <- <{ p.EndField() }> */
		nil,
		/* 93 Action13 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 94 Action14 <- <{ p.EndField()}> */
		nil,
		/* 95 Action15 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 96 Action16 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 97 Action17 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 98 Action18 <- <{ p.PushList() }> */
		nil,
		/* 99 Action19 <- <{ p.PopList() }> */
		nil,
		/* 100 Action20 <- <{ p.PushMap() }> */
		nil,
		/* 101 Action21 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 102 Action22 <- <{ p.PopMap() }> */
		nil,
		/* 103 Action23 <- <{ p.StartField("exception") }> */
		nil,
		/* 104 Action24 <- <{ p.PushValue(buffer[begin:end]); p.EndField() }> */
		nil,
		/* 105 Action25 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 106 Action26 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 107 Action27 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 108 Action28 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 109 Action29 <- <{ p.PushMap() }> */
		nil,
		/* 110 Action30 <- <{ p.PopMap() }> */
		nil,
		/* 111 Action31 <- <{ p.SetMapValue() }> */
		nil,
		/* 112 Action32 <- <{ p.PushList() }> */
		nil,
		/* 113 Action33 <- <{ p.PopList() }> */
		nil,
		/* 114 Action34 <- <{ p.SetListValue() }> */
		nil,
		/* 115 Action35 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 116 Action36 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 117 Action37 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 118 Action38 <- <{ p.PushValue(nil) }> */
		nil,
		/* 119 Action39 <- <{ p.PushValue(true) }> */
		nil,
		/* 120 Action40 <- <{ p.PushValue(false) }> */
		nil,
		/* 121 Action41 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 122 Action42 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 123 Action43 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 124 Action44 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 125 Action45 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 126 Action46 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 127 Action47 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 128 Action48 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 129 Action49 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 130 Action50 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
