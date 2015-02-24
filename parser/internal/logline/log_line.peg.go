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
	rules  [109]func() bool
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
			p.SetField("ns", buffer[begin:end])
		case ruleAction5:
			p.StartField(buffer[begin:end])
		case ruleAction6:
			p.EndField()
		case ruleAction7:
			p.SetField("duration_ms", buffer[begin:end])
		case ruleAction8:
			p.StartField(buffer[begin:end])
		case ruleAction9:
			p.EndField()
		case ruleAction10:
			p.SetField("commandType", buffer[begin:end])
			p.StartField("command")
		case ruleAction11:
			p.EndField()
		case ruleAction12:
			p.StartField("planSummary")
			p.PushList()
		case ruleAction13:
			p.EndField()
		case ruleAction14:
			p.PushMap()
			p.PushField(buffer[begin:end])
		case ruleAction15:
			p.SetMapValue()
			p.SetListValue()
		case ruleAction16:
			p.PushValue(1)
			p.SetMapValue()
			p.SetListValue()
		case ruleAction17:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction18:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction19:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction20:
			p.PushMap()
		case ruleAction21:
			p.PopMap()
		case ruleAction22:
			p.SetMapValue()
		case ruleAction23:
			p.PushList()
		case ruleAction24:
			p.PopList()
		case ruleAction25:
			p.SetListValue()
		case ruleAction26:
			p.PushField(buffer[begin:end])
		case ruleAction27:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction28:
			p.PushValue(buffer[begin:end])
		case ruleAction29:
			p.PushValue(nil)
		case ruleAction30:
			p.PushValue(true)
		case ruleAction31:
			p.PushValue(false)
		case ruleAction32:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction33:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction34:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction35:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction36:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction37:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction38:
			p.PushValue(p.Minkey())
		case ruleAction39:
			p.PushValue(p.Maxkey())
		case ruleAction40:
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
		/* 0 MongoLogLine <- <(Timestamp LogLevel? Component? Context Op NS LineField* Locks? LineField* Duration? extra? !.)> */
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
								add(ruleAction17, position)
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
								add(ruleAction18, position)
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
					position57 := position
					depth++
					{
						position58 := position
						depth++
						{
							position61, tokenIndex61, depth61 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l62
							}
							position++
							goto l61
						l62:
							position, tokenIndex, depth = position61, tokenIndex61, depth61
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l0
							}
							position++
						}
					l61:
					l59:
						{
							position60, tokenIndex60, depth60 := position, tokenIndex, depth
							{
								position63, tokenIndex63, depth63 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l64
								}
								position++
								goto l63
							l64:
								position, tokenIndex, depth = position63, tokenIndex63, depth63
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l60
								}
								position++
							}
						l63:
							goto l59
						l60:
							position, tokenIndex, depth = position60, tokenIndex60, depth60
						}
						depth--
						add(rulePegText, position58)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleOp, position57)
				}
				{
					position66 := position
					depth++
					{
						position67 := position
						depth++
					l68:
						{
							position69, tokenIndex69, depth69 := position, tokenIndex, depth
							{
								position70 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l69
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l69
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l69
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l69
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l69
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l69
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position70)
							}
							goto l68
						l69:
							position, tokenIndex, depth = position69, tokenIndex69, depth69
						}
						depth--
						add(rulePegText, position67)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction4, position)
					}
					depth--
					add(ruleNS, position66)
				}
			l73:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l74
					}
					goto l73
				l74:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
				}
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					{
						position77 := position
						depth++
						if buffer[position] != rune('l') {
							goto l75
						}
						position++
						if buffer[position] != rune('o') {
							goto l75
						}
						position++
						if buffer[position] != rune('c') {
							goto l75
						}
						position++
						if buffer[position] != rune('k') {
							goto l75
						}
						position++
						if buffer[position] != rune('s') {
							goto l75
						}
						position++
						if buffer[position] != rune('(') {
							goto l75
						}
						position++
						if buffer[position] != rune('m') {
							goto l75
						}
						position++
						if buffer[position] != rune('i') {
							goto l75
						}
						position++
						if buffer[position] != rune('c') {
							goto l75
						}
						position++
						if buffer[position] != rune('r') {
							goto l75
						}
						position++
						if buffer[position] != rune('o') {
							goto l75
						}
						position++
						if buffer[position] != rune('s') {
							goto l75
						}
						position++
						if buffer[position] != rune(')') {
							goto l75
						}
						position++
						{
							position78, tokenIndex78, depth78 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l78
							}
							goto l79
						l78:
							position, tokenIndex, depth = position78, tokenIndex78, depth78
						}
					l79:
					l80:
						{
							position81, tokenIndex81, depth81 := position, tokenIndex, depth
							{
								position82 := position
								depth++
								{
									position83 := position
									depth++
									{
										switch buffer[position] {
										case 'R':
											if buffer[position] != rune('R') {
												goto l81
											}
											position++
											break
										case 'r':
											if buffer[position] != rune('r') {
												goto l81
											}
											position++
											break
										default:
											{
												position85, tokenIndex85, depth85 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l86
												}
												position++
												goto l85
											l86:
												position, tokenIndex, depth = position85, tokenIndex85, depth85
												if buffer[position] != rune('W') {
													goto l81
												}
												position++
											}
										l85:
											break
										}
									}

									depth--
									add(rulePegText, position83)
								}
								{
									add(ruleAction5, position)
								}
								if buffer[position] != rune(':') {
									goto l81
								}
								position++
								if !_rules[ruleNumeric]() {
									goto l81
								}
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
								{
									add(ruleAction6, position)
								}
								depth--
								add(rulelock, position82)
							}
							goto l80
						l81:
							position, tokenIndex, depth = position81, tokenIndex81, depth81
						}
						depth--
						add(ruleLocks, position77)
					}
					goto l76
				l75:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
				}
			l76:
			l91:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l92
					}
					goto l91
				l92:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
				}
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					{
						position95 := position
						depth++
						{
							position96 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l93
							}
							position++
						l97:
							{
								position98, tokenIndex98, depth98 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l98
								}
								position++
								goto l97
							l98:
								position, tokenIndex, depth = position98, tokenIndex98, depth98
							}
							depth--
							add(rulePegText, position96)
						}
						if buffer[position] != rune('m') {
							goto l93
						}
						position++
						if buffer[position] != rune('s') {
							goto l93
						}
						position++
						{
							add(ruleAction7, position)
						}
						depth--
						add(ruleDuration, position95)
					}
					goto l94
				l93:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
				}
			l94:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					{
						position102 := position
						depth++
						{
							position103 := position
							depth++
							if !matchDot() {
								goto l100
							}
						l104:
							{
								position105, tokenIndex105, depth105 := position, tokenIndex, depth
								if !matchDot() {
									goto l105
								}
								goto l104
							l105:
								position, tokenIndex, depth = position105, tokenIndex105, depth105
							}
							depth--
							add(rulePegText, position103)
						}
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleextra, position102)
					}
					goto l101
				l100:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
				}
			l101:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if !matchDot() {
						goto l107
					}
					goto l0
				l107:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
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
		/* 6 LineField <- <((commandField / planSummaryField / plainField) S?)> */
		func() bool {
			position113, tokenIndex113, depth113 := position, tokenIndex, depth
			{
				position114 := position
				depth++
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					{
						position117 := position
						depth++
						if buffer[position] != rune('c') {
							goto l116
						}
						position++
						if buffer[position] != rune('o') {
							goto l116
						}
						position++
						if buffer[position] != rune('m') {
							goto l116
						}
						position++
						if buffer[position] != rune('m') {
							goto l116
						}
						position++
						if buffer[position] != rune('a') {
							goto l116
						}
						position++
						if buffer[position] != rune('n') {
							goto l116
						}
						position++
						if buffer[position] != rune('d') {
							goto l116
						}
						position++
						if buffer[position] != rune(':') {
							goto l116
						}
						position++
						if buffer[position] != rune(' ') {
							goto l116
						}
						position++
						{
							position118 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l116
							}
						l119:
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l120
								}
								goto l119
							l120:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
							}
							depth--
							add(rulePegText, position118)
						}
						{
							add(ruleAction10, position)
						}
						if !_rules[ruleLineValue]() {
							goto l116
						}
						{
							add(ruleAction11, position)
						}
						depth--
						add(rulecommandField, position117)
					}
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					{
						position124 := position
						depth++
						if buffer[position] != rune('p') {
							goto l123
						}
						position++
						if buffer[position] != rune('l') {
							goto l123
						}
						position++
						if buffer[position] != rune('a') {
							goto l123
						}
						position++
						if buffer[position] != rune('n') {
							goto l123
						}
						position++
						if buffer[position] != rune('S') {
							goto l123
						}
						position++
						if buffer[position] != rune('u') {
							goto l123
						}
						position++
						if buffer[position] != rune('m') {
							goto l123
						}
						position++
						if buffer[position] != rune('m') {
							goto l123
						}
						position++
						if buffer[position] != rune('a') {
							goto l123
						}
						position++
						if buffer[position] != rune('r') {
							goto l123
						}
						position++
						if buffer[position] != rune('y') {
							goto l123
						}
						position++
						if buffer[position] != rune(':') {
							goto l123
						}
						position++
						if buffer[position] != rune(' ') {
							goto l123
						}
						position++
						{
							add(ruleAction12, position)
						}
						{
							position126 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l123
							}
						l127:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l128
								}
								position++
								if buffer[position] != rune(' ') {
									goto l128
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l128
								}
								goto l127
							l128:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
							}
							depth--
							add(ruleplanSummaryElements, position126)
						}
						{
							add(ruleAction13, position)
						}
						depth--
						add(ruleplanSummaryField, position124)
					}
					goto l115
				l123:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					{
						position130 := position
						depth++
						{
							position131 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l113
							}
						l132:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l133
								}
								goto l132
							l133:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
							}
							depth--
							add(rulePegText, position131)
						}
						if buffer[position] != rune(':') {
							goto l113
						}
						position++
						{
							position134, tokenIndex134, depth134 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l134
							}
							goto l135
						l134:
							position, tokenIndex, depth = position134, tokenIndex134, depth134
						}
					l135:
						{
							add(ruleAction8, position)
						}
						if !_rules[ruleLineValue]() {
							goto l113
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleplainField, position130)
					}
				}
			l115:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l138
					}
					goto l139
				l138:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
				}
			l139:
				depth--
				add(ruleLineField, position114)
			}
			return true
		l113:
			position, tokenIndex, depth = position113, tokenIndex113, depth113
			return false
		},
		/* 7 NS <- <(<nsChar*> ' ' Action4)> */
		nil,
		/* 8 Locks <- <('l' 'o' 'c' 'k' 's' '(' 'm' 'i' 'c' 'r' 'o' 's' ')' S? lock*)> */
		nil,
		/* 9 lock <- <(<((&('R') 'R') | (&('r') 'r') | (&('W' | 'w') ('w' / 'W')))> Action5 ':' Numeric S? Action6)> */
		nil,
		/* 10 Duration <- <(<[0-9]+> ('m' 's') Action7)> */
		nil,
		/* 11 plainField <- <(<fieldChar+> ':' S? Action8 LineValue Action9)> */
		nil,
		/* 12 commandField <- <('c' 'o' 'm' 'm' 'a' 'n' 'd' ':' ' ' <fieldChar+> Action10 LineValue Action11)> */
		nil,
		/* 13 planSummaryField <- <('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y' ':' ' ' Action12 planSummaryElements Action13)> */
		nil,
		/* 14 planSummaryElements <- <(planSummaryElem (',' ' ' planSummaryElem)*)> */
		nil,
		/* 15 planSummaryElem <- <(<planSummaryStage> Action14 planSummary)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				{
					position150 := position
					depth++
					{
						position151 := position
						depth++
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l153
							}
							position++
							if buffer[position] != rune('N') {
								goto l153
							}
							position++
							if buffer[position] != rune('D') {
								goto l153
							}
							position++
							if buffer[position] != rune('_') {
								goto l153
							}
							position++
							if buffer[position] != rune('H') {
								goto l153
							}
							position++
							if buffer[position] != rune('A') {
								goto l153
							}
							position++
							if buffer[position] != rune('S') {
								goto l153
							}
							position++
							if buffer[position] != rune('H') {
								goto l153
							}
							position++
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('C') {
								goto l154
							}
							position++
							if buffer[position] != rune('A') {
								goto l154
							}
							position++
							if buffer[position] != rune('C') {
								goto l154
							}
							position++
							if buffer[position] != rune('H') {
								goto l154
							}
							position++
							if buffer[position] != rune('E') {
								goto l154
							}
							position++
							if buffer[position] != rune('D') {
								goto l154
							}
							position++
							if buffer[position] != rune('_') {
								goto l154
							}
							position++
							if buffer[position] != rune('P') {
								goto l154
							}
							position++
							if buffer[position] != rune('L') {
								goto l154
							}
							position++
							if buffer[position] != rune('A') {
								goto l154
							}
							position++
							if buffer[position] != rune('N') {
								goto l154
							}
							position++
							goto l152
						l154:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('C') {
								goto l155
							}
							position++
							if buffer[position] != rune('O') {
								goto l155
							}
							position++
							if buffer[position] != rune('L') {
								goto l155
							}
							position++
							if buffer[position] != rune('L') {
								goto l155
							}
							position++
							if buffer[position] != rune('S') {
								goto l155
							}
							position++
							if buffer[position] != rune('C') {
								goto l155
							}
							position++
							if buffer[position] != rune('A') {
								goto l155
							}
							position++
							if buffer[position] != rune('N') {
								goto l155
							}
							position++
							goto l152
						l155:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('C') {
								goto l156
							}
							position++
							if buffer[position] != rune('O') {
								goto l156
							}
							position++
							if buffer[position] != rune('U') {
								goto l156
							}
							position++
							if buffer[position] != rune('N') {
								goto l156
							}
							position++
							if buffer[position] != rune('T') {
								goto l156
							}
							position++
							goto l152
						l156:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('D') {
								goto l157
							}
							position++
							if buffer[position] != rune('E') {
								goto l157
							}
							position++
							if buffer[position] != rune('L') {
								goto l157
							}
							position++
							if buffer[position] != rune('E') {
								goto l157
							}
							position++
							if buffer[position] != rune('T') {
								goto l157
							}
							position++
							if buffer[position] != rune('E') {
								goto l157
							}
							position++
							goto l152
						l157:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('G') {
								goto l158
							}
							position++
							if buffer[position] != rune('E') {
								goto l158
							}
							position++
							if buffer[position] != rune('O') {
								goto l158
							}
							position++
							if buffer[position] != rune('_') {
								goto l158
							}
							position++
							if buffer[position] != rune('N') {
								goto l158
							}
							position++
							if buffer[position] != rune('E') {
								goto l158
							}
							position++
							if buffer[position] != rune('A') {
								goto l158
							}
							position++
							if buffer[position] != rune('R') {
								goto l158
							}
							position++
							if buffer[position] != rune('_') {
								goto l158
							}
							position++
							if buffer[position] != rune('2') {
								goto l158
							}
							position++
							if buffer[position] != rune('D') {
								goto l158
							}
							position++
							goto l152
						l158:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('G') {
								goto l159
							}
							position++
							if buffer[position] != rune('E') {
								goto l159
							}
							position++
							if buffer[position] != rune('O') {
								goto l159
							}
							position++
							if buffer[position] != rune('_') {
								goto l159
							}
							position++
							if buffer[position] != rune('N') {
								goto l159
							}
							position++
							if buffer[position] != rune('E') {
								goto l159
							}
							position++
							if buffer[position] != rune('A') {
								goto l159
							}
							position++
							if buffer[position] != rune('R') {
								goto l159
							}
							position++
							if buffer[position] != rune('_') {
								goto l159
							}
							position++
							if buffer[position] != rune('2') {
								goto l159
							}
							position++
							if buffer[position] != rune('D') {
								goto l159
							}
							position++
							if buffer[position] != rune('S') {
								goto l159
							}
							position++
							if buffer[position] != rune('P') {
								goto l159
							}
							position++
							if buffer[position] != rune('H') {
								goto l159
							}
							position++
							if buffer[position] != rune('E') {
								goto l159
							}
							position++
							if buffer[position] != rune('R') {
								goto l159
							}
							position++
							if buffer[position] != rune('E') {
								goto l159
							}
							position++
							goto l152
						l159:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('I') {
								goto l160
							}
							position++
							if buffer[position] != rune('D') {
								goto l160
							}
							position++
							if buffer[position] != rune('H') {
								goto l160
							}
							position++
							if buffer[position] != rune('A') {
								goto l160
							}
							position++
							if buffer[position] != rune('C') {
								goto l160
							}
							position++
							if buffer[position] != rune('K') {
								goto l160
							}
							position++
							goto l152
						l160:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('S') {
								goto l161
							}
							position++
							if buffer[position] != rune('O') {
								goto l161
							}
							position++
							if buffer[position] != rune('R') {
								goto l161
							}
							position++
							if buffer[position] != rune('T') {
								goto l161
							}
							position++
							if buffer[position] != rune('_') {
								goto l161
							}
							position++
							if buffer[position] != rune('M') {
								goto l161
							}
							position++
							if buffer[position] != rune('E') {
								goto l161
							}
							position++
							if buffer[position] != rune('R') {
								goto l161
							}
							position++
							if buffer[position] != rune('G') {
								goto l161
							}
							position++
							if buffer[position] != rune('E') {
								goto l161
							}
							position++
							goto l152
						l161:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('S') {
								goto l162
							}
							position++
							if buffer[position] != rune('H') {
								goto l162
							}
							position++
							if buffer[position] != rune('A') {
								goto l162
							}
							position++
							if buffer[position] != rune('R') {
								goto l162
							}
							position++
							if buffer[position] != rune('D') {
								goto l162
							}
							position++
							if buffer[position] != rune('I') {
								goto l162
							}
							position++
							if buffer[position] != rune('N') {
								goto l162
							}
							position++
							if buffer[position] != rune('G') {
								goto l162
							}
							position++
							if buffer[position] != rune('_') {
								goto l162
							}
							position++
							if buffer[position] != rune('F') {
								goto l162
							}
							position++
							if buffer[position] != rune('I') {
								goto l162
							}
							position++
							if buffer[position] != rune('L') {
								goto l162
							}
							position++
							if buffer[position] != rune('T') {
								goto l162
							}
							position++
							if buffer[position] != rune('E') {
								goto l162
							}
							position++
							if buffer[position] != rune('R') {
								goto l162
							}
							position++
							goto l152
						l162:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('S') {
								goto l163
							}
							position++
							if buffer[position] != rune('K') {
								goto l163
							}
							position++
							if buffer[position] != rune('I') {
								goto l163
							}
							position++
							if buffer[position] != rune('P') {
								goto l163
							}
							position++
							goto l152
						l163:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('S') {
								goto l164
							}
							position++
							if buffer[position] != rune('O') {
								goto l164
							}
							position++
							if buffer[position] != rune('R') {
								goto l164
							}
							position++
							if buffer[position] != rune('T') {
								goto l164
							}
							position++
							goto l152
						l164:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('X') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('B') {
										goto l148
									}
									position++
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									if buffer[position] != rune('L') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									if buffer[position] != rune('_') {
										goto l148
									}
									position++
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									if buffer[position] != rune('R') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('J') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('R') {
										goto l148
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('L') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('_') {
										goto l148
									}
									position++
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									if buffer[position] != rune('L') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('M') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									if buffer[position] != rune('_') {
										goto l148
									}
									position++
									if buffer[position] != rune('M') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('X') {
										goto l148
									}
									position++
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l148
									}
									position++
									if buffer[position] != rune('R') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('P') {
										goto l148
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('H') {
										goto l148
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('F') {
										goto l148
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('I') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('U') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('_') {
										goto l148
									}
									position++
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									if buffer[position] != rune('C') {
										goto l148
									}
									position++
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l148
									}
									position++
									if buffer[position] != rune('N') {
										goto l148
									}
									position++
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									if buffer[position] != rune('_') {
										goto l148
									}
									position++
									if buffer[position] != rune('S') {
										goto l148
									}
									position++
									if buffer[position] != rune('O') {
										goto l148
									}
									position++
									if buffer[position] != rune('R') {
										goto l148
									}
									position++
									if buffer[position] != rune('T') {
										goto l148
									}
									position++
									if buffer[position] != rune('E') {
										goto l148
									}
									position++
									if buffer[position] != rune('D') {
										goto l148
									}
									position++
									break
								}
							}

						}
					l152:
						depth--
						add(ruleplanSummaryStage, position151)
					}
					depth--
					add(rulePegText, position150)
				}
				{
					add(ruleAction14, position)
				}
				{
					position167 := position
					depth++
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l169
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l169
						}
						{
							add(ruleAction15, position)
						}
						goto l168
					l169:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
						{
							add(ruleAction16, position)
						}
					}
				l168:
					depth--
					add(ruleplanSummary, position167)
				}
				depth--
				add(ruleplanSummaryElem, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 16 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 17 planSummary <- <((' ' LineValue Action15) / Action16)> */
		nil,
		/* 18 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l177
					}
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if !_rules[ruleNumeric]() {
						goto l174
					}
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l178
					}
					goto l179
				l178:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
				}
			l179:
				depth--
				add(ruleLineValue, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 19 timestamp24 <- <(<(date ' ' time)> Action17)> */
		nil,
		/* 20 timestamp26 <- <(<datetime26> Action18)> */
		nil,
		/* 21 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 22 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 23 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l184
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l184
				}
				position++
				depth--
				add(ruledigit2, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 24 date <- <(day ' ' month ' '+ dayNum)> */
		nil,
		/* 25 tz <- <('+' [0-9]+)> */
		nil,
		/* 26 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				{
					position190 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l188
					}
					depth--
					add(rulehour, position190)
				}
				if buffer[position] != rune(':') {
					goto l188
				}
				position++
				{
					position191 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l188
					}
					depth--
					add(ruleminute, position191)
				}
				if buffer[position] != rune(':') {
					goto l188
				}
				position++
				{
					position192 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l188
					}
					depth--
					add(rulesecond, position192)
				}
				if buffer[position] != rune('.') {
					goto l188
				}
				position++
				{
					position193 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l188
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l188
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l188
					}
					position++
					depth--
					add(rulemillisecond, position193)
				}
				depth--
				add(ruletime, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 27 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 28 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 29 dayNum <- <([0-9] [0-9]?)> */
		nil,
		/* 30 hour <- <digit2> */
		nil,
		/* 31 minute <- <digit2> */
		nil,
		/* 32 second <- <digit2> */
		nil,
		/* 33 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 34 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 35 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 36 extra <- <(<.+> Action19)> */
		nil,
		/* 37 S <- <' '+> */
		func() bool {
			position204, tokenIndex204, depth204 := position, tokenIndex, depth
			{
				position205 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l204
				}
				position++
			l206:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
				}
				depth--
				add(ruleS, position205)
			}
			return true
		l204:
			position, tokenIndex, depth = position204, tokenIndex204, depth204
			return false
		},
		/* 38 Doc <- <('{' Action20 DocElements? '}' Action21)> */
		func() bool {
			position208, tokenIndex208, depth208 := position, tokenIndex, depth
			{
				position209 := position
				depth++
				if buffer[position] != rune('{') {
					goto l208
				}
				position++
				{
					add(ruleAction20, position)
				}
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					{
						position213 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l211
						}
					l214:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l215
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l215
							}
							goto l214
						l215:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
						}
						depth--
						add(ruleDocElements, position213)
					}
					goto l212
				l211:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
				}
			l212:
				if buffer[position] != rune('}') {
					goto l208
				}
				position++
				{
					add(ruleAction21, position)
				}
				depth--
				add(ruleDoc, position209)
			}
			return true
		l208:
			position, tokenIndex, depth = position208, tokenIndex208, depth208
			return false
		},
		/* 39 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 40 DocElem <- <(S? Field S? Value S? Action22)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l220
					}
					goto l221
				l220:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
				}
			l221:
				{
					position222 := position
					depth++
					{
						position223 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l218
						}
					l224:
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l225
							}
							goto l224
						l225:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
						}
						depth--
						add(rulePegText, position223)
					}
					if buffer[position] != rune(':') {
						goto l218
					}
					position++
					{
						add(ruleAction26, position)
					}
					depth--
					add(ruleField, position222)
				}
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l227
					}
					goto l228
				l227:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
				}
			l228:
				if !_rules[ruleValue]() {
					goto l218
				}
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l229
					}
					goto l230
				l229:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
				}
			l230:
				{
					add(ruleAction22, position)
				}
				depth--
				add(ruleDocElem, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 41 List <- <('[' Action23 ListElements? ']' Action24)> */
		nil,
		/* 42 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 43 ListElem <- <(S? Value S? Action25)> */
		func() bool {
			position234, tokenIndex234, depth234 := position, tokenIndex, depth
			{
				position235 := position
				depth++
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
				if !_rules[ruleValue]() {
					goto l234
				}
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
					add(ruleAction25, position)
				}
				depth--
				add(ruleListElem, position235)
			}
			return true
		l234:
			position, tokenIndex, depth = position234, tokenIndex234, depth234
			return false
		},
		/* 44 Field <- <(<fieldChar+> ':' Action26)> */
		nil,
		/* 45 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position242, tokenIndex242, depth242 := position, tokenIndex, depth
			{
				position243 := position
				depth++
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					{
						position246 := position
						depth++
						if buffer[position] != rune('n') {
							goto l245
						}
						position++
						if buffer[position] != rune('u') {
							goto l245
						}
						position++
						if buffer[position] != rune('l') {
							goto l245
						}
						position++
						if buffer[position] != rune('l') {
							goto l245
						}
						position++
						{
							add(ruleAction29, position)
						}
						depth--
						add(ruleNull, position246)
					}
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					{
						position249 := position
						depth++
						if buffer[position] != rune('M') {
							goto l248
						}
						position++
						if buffer[position] != rune('i') {
							goto l248
						}
						position++
						if buffer[position] != rune('n') {
							goto l248
						}
						position++
						if buffer[position] != rune('K') {
							goto l248
						}
						position++
						if buffer[position] != rune('e') {
							goto l248
						}
						position++
						if buffer[position] != rune('y') {
							goto l248
						}
						position++
						{
							add(ruleAction38, position)
						}
						depth--
						add(ruleMinKey, position249)
					}
					goto l244
				l248:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					{
						switch buffer[position] {
						case 'M':
							{
								position252 := position
								depth++
								if buffer[position] != rune('M') {
									goto l242
								}
								position++
								if buffer[position] != rune('a') {
									goto l242
								}
								position++
								if buffer[position] != rune('x') {
									goto l242
								}
								position++
								if buffer[position] != rune('K') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('y') {
									goto l242
								}
								position++
								{
									add(ruleAction39, position)
								}
								depth--
								add(ruleMaxKey, position252)
							}
							break
						case 'u':
							{
								position254 := position
								depth++
								if buffer[position] != rune('u') {
									goto l242
								}
								position++
								if buffer[position] != rune('n') {
									goto l242
								}
								position++
								if buffer[position] != rune('d') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('f') {
									goto l242
								}
								position++
								if buffer[position] != rune('i') {
									goto l242
								}
								position++
								if buffer[position] != rune('n') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('d') {
									goto l242
								}
								position++
								{
									add(ruleAction40, position)
								}
								depth--
								add(ruleUndefined, position254)
							}
							break
						case 'N':
							{
								position256 := position
								depth++
								if buffer[position] != rune('N') {
									goto l242
								}
								position++
								if buffer[position] != rune('u') {
									goto l242
								}
								position++
								if buffer[position] != rune('m') {
									goto l242
								}
								position++
								if buffer[position] != rune('b') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('r') {
									goto l242
								}
								position++
								if buffer[position] != rune('L') {
									goto l242
								}
								position++
								if buffer[position] != rune('o') {
									goto l242
								}
								position++
								if buffer[position] != rune('n') {
									goto l242
								}
								position++
								if buffer[position] != rune('g') {
									goto l242
								}
								position++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								{
									position257 := position
									depth++
									{
										position260, tokenIndex260, depth260 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l260
										}
										position++
										goto l242
									l260:
										position, tokenIndex, depth = position260, tokenIndex260, depth260
									}
									if !matchDot() {
										goto l242
									}
								l258:
									{
										position259, tokenIndex259, depth259 := position, tokenIndex, depth
										{
											position261, tokenIndex261, depth261 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l261
											}
											position++
											goto l259
										l261:
											position, tokenIndex, depth = position261, tokenIndex261, depth261
										}
										if !matchDot() {
											goto l259
										}
										goto l258
									l259:
										position, tokenIndex, depth = position259, tokenIndex259, depth259
									}
									depth--
									add(rulePegText, position257)
								}
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleNumberLong, position256)
							}
							break
						case '/':
							{
								position263 := position
								depth++
								if buffer[position] != rune('/') {
									goto l242
								}
								position++
								{
									position264 := position
									depth++
									{
										position265 := position
										depth++
										{
											position268 := position
											depth++
											{
												position269, tokenIndex269, depth269 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l269
												}
												position++
												goto l242
											l269:
												position, tokenIndex, depth = position269, tokenIndex269, depth269
											}
											if !matchDot() {
												goto l242
											}
											depth--
											add(ruleregexChar, position268)
										}
									l266:
										{
											position267, tokenIndex267, depth267 := position, tokenIndex, depth
											{
												position270 := position
												depth++
												{
													position271, tokenIndex271, depth271 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l271
													}
													position++
													goto l267
												l271:
													position, tokenIndex, depth = position271, tokenIndex271, depth271
												}
												if !matchDot() {
													goto l267
												}
												depth--
												add(ruleregexChar, position270)
											}
											goto l266
										l267:
											position, tokenIndex, depth = position267, tokenIndex267, depth267
										}
										if buffer[position] != rune('/') {
											goto l242
										}
										position++
									l272:
										{
											position273, tokenIndex273, depth273 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l273
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l273
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l273
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l273
													}
													position++
													break
												}
											}

											goto l272
										l273:
											position, tokenIndex, depth = position273, tokenIndex273, depth273
										}
										depth--
										add(ruleregexBody, position265)
									}
									depth--
									add(rulePegText, position264)
								}
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleRegex, position263)
							}
							break
						case 'T':
							{
								position276 := position
								depth++
								if buffer[position] != rune('T') {
									goto l242
								}
								position++
								if buffer[position] != rune('i') {
									goto l242
								}
								position++
								if buffer[position] != rune('m') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('s') {
									goto l242
								}
								position++
								if buffer[position] != rune('t') {
									goto l242
								}
								position++
								if buffer[position] != rune('a') {
									goto l242
								}
								position++
								if buffer[position] != rune('m') {
									goto l242
								}
								position++
								if buffer[position] != rune('p') {
									goto l242
								}
								position++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								{
									position277 := position
									depth++
									{
										position280, tokenIndex280, depth280 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l280
										}
										position++
										goto l242
									l280:
										position, tokenIndex, depth = position280, tokenIndex280, depth280
									}
									if !matchDot() {
										goto l242
									}
								l278:
									{
										position279, tokenIndex279, depth279 := position, tokenIndex, depth
										{
											position281, tokenIndex281, depth281 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l281
											}
											position++
											goto l279
										l281:
											position, tokenIndex, depth = position281, tokenIndex281, depth281
										}
										if !matchDot() {
											goto l279
										}
										goto l278
									l279:
										position, tokenIndex, depth = position279, tokenIndex279, depth279
									}
									depth--
									add(rulePegText, position277)
								}
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleTimestampVal, position276)
							}
							break
						case 'B':
							{
								position283 := position
								depth++
								if buffer[position] != rune('B') {
									goto l242
								}
								position++
								if buffer[position] != rune('i') {
									goto l242
								}
								position++
								if buffer[position] != rune('n') {
									goto l242
								}
								position++
								if buffer[position] != rune('D') {
									goto l242
								}
								position++
								if buffer[position] != rune('a') {
									goto l242
								}
								position++
								if buffer[position] != rune('t') {
									goto l242
								}
								position++
								if buffer[position] != rune('a') {
									goto l242
								}
								position++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								{
									position284 := position
									depth++
									{
										position287, tokenIndex287, depth287 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l287
										}
										position++
										goto l242
									l287:
										position, tokenIndex, depth = position287, tokenIndex287, depth287
									}
									if !matchDot() {
										goto l242
									}
								l285:
									{
										position286, tokenIndex286, depth286 := position, tokenIndex, depth
										{
											position288, tokenIndex288, depth288 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l288
											}
											position++
											goto l286
										l288:
											position, tokenIndex, depth = position288, tokenIndex288, depth288
										}
										if !matchDot() {
											goto l286
										}
										goto l285
									l286:
										position, tokenIndex, depth = position286, tokenIndex286, depth286
									}
									depth--
									add(rulePegText, position284)
								}
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								{
									add(ruleAction34, position)
								}
								depth--
								add(ruleBinData, position283)
							}
							break
						case 'D', 'n':
							{
								position290 := position
								depth++
								{
									position291, tokenIndex291, depth291 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l291
									}
									position++
									if buffer[position] != rune('e') {
										goto l291
									}
									position++
									if buffer[position] != rune('w') {
										goto l291
									}
									position++
									if buffer[position] != rune(' ') {
										goto l291
									}
									position++
									goto l292
								l291:
									position, tokenIndex, depth = position291, tokenIndex291, depth291
								}
							l292:
								if buffer[position] != rune('D') {
									goto l242
								}
								position++
								if buffer[position] != rune('a') {
									goto l242
								}
								position++
								if buffer[position] != rune('t') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								{
									position293 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l242
									}
									position++
								l294:
									{
										position295, tokenIndex295, depth295 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l295
										}
										position++
										goto l294
									l295:
										position, tokenIndex, depth = position295, tokenIndex295, depth295
									}
									depth--
									add(rulePegText, position293)
								}
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								{
									add(ruleAction32, position)
								}
								depth--
								add(ruleDate, position290)
							}
							break
						case 'O':
							{
								position297 := position
								depth++
								if buffer[position] != rune('O') {
									goto l242
								}
								position++
								if buffer[position] != rune('b') {
									goto l242
								}
								position++
								if buffer[position] != rune('j') {
									goto l242
								}
								position++
								if buffer[position] != rune('e') {
									goto l242
								}
								position++
								if buffer[position] != rune('c') {
									goto l242
								}
								position++
								if buffer[position] != rune('t') {
									goto l242
								}
								position++
								if buffer[position] != rune('I') {
									goto l242
								}
								position++
								if buffer[position] != rune('d') {
									goto l242
								}
								position++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								if buffer[position] != rune('"') {
									goto l242
								}
								position++
								{
									position298 := position
									depth++
								l299:
									{
										position300, tokenIndex300, depth300 := position, tokenIndex, depth
										{
											position301 := position
											depth++
											{
												position302, tokenIndex302, depth302 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l303
												}
												position++
												goto l302
											l303:
												position, tokenIndex, depth = position302, tokenIndex302, depth302
												{
													position304, tokenIndex304, depth304 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l305
													}
													position++
													goto l304
												l305:
													position, tokenIndex, depth = position304, tokenIndex304, depth304
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l300
													}
													position++
												}
											l304:
											}
										l302:
											depth--
											add(rulehexChar, position301)
										}
										goto l299
									l300:
										position, tokenIndex, depth = position300, tokenIndex300, depth300
									}
									depth--
									add(rulePegText, position298)
								}
								if buffer[position] != rune('"') {
									goto l242
								}
								position++
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleObjectID, position297)
							}
							break
						case '"':
							{
								position307 := position
								depth++
								if buffer[position] != rune('"') {
									goto l242
								}
								position++
								{
									position308 := position
									depth++
								l309:
									{
										position310, tokenIndex310, depth310 := position, tokenIndex, depth
										{
											position311 := position
											depth++
											{
												position312, tokenIndex312, depth312 := position, tokenIndex, depth
												{
													position314, tokenIndex314, depth314 := position, tokenIndex, depth
													{
														position315, tokenIndex315, depth315 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l316
														}
														position++
														goto l315
													l316:
														position, tokenIndex, depth = position315, tokenIndex315, depth315
														if buffer[position] != rune('\\') {
															goto l314
														}
														position++
													}
												l315:
													goto l313
												l314:
													position, tokenIndex, depth = position314, tokenIndex314, depth314
												}
												if !matchDot() {
													goto l313
												}
												goto l312
											l313:
												position, tokenIndex, depth = position312, tokenIndex312, depth312
												if buffer[position] != rune('\\') {
													goto l310
												}
												position++
												{
													position317, tokenIndex317, depth317 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l318
													}
													position++
													goto l317
												l318:
													position, tokenIndex, depth = position317, tokenIndex317, depth317
													if buffer[position] != rune('\\') {
														goto l310
													}
													position++
												}
											l317:
											}
										l312:
											depth--
											add(rulestringChar, position311)
										}
										goto l309
									l310:
										position, tokenIndex, depth = position310, tokenIndex310, depth310
									}
									depth--
									add(rulePegText, position308)
								}
								if buffer[position] != rune('"') {
									goto l242
								}
								position++
								{
									add(ruleAction28, position)
								}
								depth--
								add(ruleString, position307)
							}
							break
						case 'f', 't':
							{
								position320 := position
								depth++
								{
									position321, tokenIndex321, depth321 := position, tokenIndex, depth
									{
										position323 := position
										depth++
										if buffer[position] != rune('t') {
											goto l322
										}
										position++
										if buffer[position] != rune('r') {
											goto l322
										}
										position++
										if buffer[position] != rune('u') {
											goto l322
										}
										position++
										if buffer[position] != rune('e') {
											goto l322
										}
										position++
										{
											add(ruleAction30, position)
										}
										depth--
										add(ruleTrue, position323)
									}
									goto l321
								l322:
									position, tokenIndex, depth = position321, tokenIndex321, depth321
									{
										position325 := position
										depth++
										if buffer[position] != rune('f') {
											goto l242
										}
										position++
										if buffer[position] != rune('a') {
											goto l242
										}
										position++
										if buffer[position] != rune('l') {
											goto l242
										}
										position++
										if buffer[position] != rune('s') {
											goto l242
										}
										position++
										if buffer[position] != rune('e') {
											goto l242
										}
										position++
										{
											add(ruleAction31, position)
										}
										depth--
										add(ruleFalse, position325)
									}
								}
							l321:
								depth--
								add(ruleBoolean, position320)
							}
							break
						case '[':
							{
								position327 := position
								depth++
								if buffer[position] != rune('[') {
									goto l242
								}
								position++
								{
									add(ruleAction23, position)
								}
								{
									position329, tokenIndex329, depth329 := position, tokenIndex, depth
									{
										position331 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l329
										}
									l332:
										{
											position333, tokenIndex333, depth333 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l333
											}
											position++
											if !_rules[ruleListElem]() {
												goto l333
											}
											goto l332
										l333:
											position, tokenIndex, depth = position333, tokenIndex333, depth333
										}
										depth--
										add(ruleListElements, position331)
									}
									goto l330
								l329:
									position, tokenIndex, depth = position329, tokenIndex329, depth329
								}
							l330:
								if buffer[position] != rune(']') {
									goto l242
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleList, position327)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l242
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l242
							}
							break
						}
					}

				}
			l244:
				depth--
				add(ruleValue, position243)
			}
			return true
		l242:
			position, tokenIndex, depth = position242, tokenIndex242, depth242
			return false
		},
		/* 46 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action27)> */
		func() bool {
			position335, tokenIndex335, depth335 := position, tokenIndex, depth
			{
				position336 := position
				depth++
				{
					position337 := position
					depth++
					{
						position338, tokenIndex338, depth338 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l338
						}
						position++
						goto l339
					l338:
						position, tokenIndex, depth = position338, tokenIndex338, depth338
					}
				l339:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l335
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
					{
						position342, tokenIndex342, depth342 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l342
						}
						position++
						goto l343
					l342:
						position, tokenIndex, depth = position342, tokenIndex342, depth342
					}
				l343:
				l344:
					{
						position345, tokenIndex345, depth345 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l345
						}
						position++
						goto l344
					l345:
						position, tokenIndex, depth = position345, tokenIndex345, depth345
					}
					depth--
					add(rulePegText, position337)
				}
				{
					add(ruleAction27, position)
				}
				depth--
				add(ruleNumeric, position336)
			}
			return true
		l335:
			position, tokenIndex, depth = position335, tokenIndex335, depth335
			return false
		},
		/* 47 Boolean <- <(True / False)> */
		nil,
		/* 48 String <- <('"' <stringChar*> '"' Action28)> */
		nil,
		/* 49 Null <- <('n' 'u' 'l' 'l' Action29)> */
		nil,
		/* 50 True <- <('t' 'r' 'u' 'e' Action30)> */
		nil,
		/* 51 False <- <('f' 'a' 'l' 's' 'e' Action31)> */
		nil,
		/* 52 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action32)> */
		nil,
		/* 53 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action33)> */
		nil,
		/* 54 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action34)> */
		nil,
		/* 55 Regex <- <('/' <regexBody> Action35)> */
		nil,
		/* 56 TimestampVal <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action36)> */
		nil,
		/* 57 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action37)> */
		nil,
		/* 58 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action38)> */
		nil,
		/* 59 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action39)> */
		nil,
		/* 60 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action40)> */
		nil,
		/* 61 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 62 regexChar <- <(!'/' .)> */
		nil,
		/* 63 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 64 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 65 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l365
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l365
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l365
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l365
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l365
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l365
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l365
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		nil,
		/* 68 Action0 <- <{ p.SetField("log_level", buffer[begin:end]) }> */
		nil,
		/* 69 Action1 <- <{ p.SetField("component", buffer[begin:end]) }> */
		nil,
		/* 70 Action2 <- <{ p.SetField("context", buffer[begin:end]) }> */
		nil,
		/* 71 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 72 Action4 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 73 Action5 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 74 Action6 <- <{ p.EndField() }> */
		nil,
		/* 75 Action7 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 76 Action8 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 77 Action9 <- <{ p.EndField() }> */
		nil,
		/* 78 Action10 <- <{ p.SetField("commandType", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 79 Action11 <- <{ p.EndField() }> */
		nil,
		/* 80 Action12 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 81 Action13 <- <{ p.EndField()}> */
		nil,
		/* 82 Action14 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 83 Action15 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 84 Action16 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 85 Action17 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 86 Action18 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 87 Action19 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 88 Action20 <- <{ p.PushMap() }> */
		nil,
		/* 89 Action21 <- <{ p.PopMap() }> */
		nil,
		/* 90 Action22 <- <{ p.SetMapValue() }> */
		nil,
		/* 91 Action23 <- <{ p.PushList() }> */
		nil,
		/* 92 Action24 <- <{ p.PopList() }> */
		nil,
		/* 93 Action25 <- <{ p.SetListValue() }> */
		nil,
		/* 94 Action26 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 95 Action27 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 96 Action28 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 97 Action29 <- <{ p.PushValue(nil) }> */
		nil,
		/* 98 Action30 <- <{ p.PushValue(true) }> */
		nil,
		/* 99 Action31 <- <{ p.PushValue(false) }> */
		nil,
		/* 100 Action32 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 101 Action33 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 102 Action34 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 103 Action35 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 104 Action36 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 105 Action37 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 106 Action38 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 107 Action39 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 108 Action40 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
