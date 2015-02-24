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
									{
										position10 := position
										depth++
										{
											position11, tokenIndex11, depth11 := position, tokenIndex, depth
											if !_rules[ruledigit2]() {
												goto l11
											}
											goto l12
										l11:
											position, tokenIndex, depth = position11, tokenIndex11, depth11
										}
									l12:
										depth--
										add(ruledayNum, position10)
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
							position14 := position
							depth++
							{
								position15 := position
								depth++
								{
									position16 := position
									depth++
									{
										position17 := position
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
										add(ruledigit4, position17)
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
										position18, tokenIndex18, depth18 := position, tokenIndex, depth
										{
											position20 := position
											depth++
											if buffer[position] != rune('+') {
												goto l18
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l18
											}
											position++
										l21:
											{
												position22, tokenIndex22, depth22 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l22
												}
												position++
												goto l21
											l22:
												position, tokenIndex, depth = position22, tokenIndex22, depth22
											}
											depth--
											add(ruletz, position20)
										}
										goto l19
									l18:
										position, tokenIndex, depth = position18, tokenIndex18, depth18
									}
								l19:
									depth--
									add(ruledatetime26, position16)
								}
								depth--
								add(rulePegText, position15)
							}
							{
								add(ruleAction18, position)
							}
							depth--
							add(ruletimestamp26, position14)
						}
					}
				l3:
					{
						position24, tokenIndex24, depth24 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l24
						}
						goto l25
					l24:
						position, tokenIndex, depth = position24, tokenIndex24, depth24
					}
				l25:
					depth--
					add(ruleTimestamp, position2)
				}
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					{
						position28 := position
						depth++
						{
							position29 := position
							depth++
							{
								position30, tokenIndex30, depth30 := position, tokenIndex, depth
								if buffer[position] != rune('I') {
									goto l31
								}
								position++
								goto l30
							l31:
								position, tokenIndex, depth = position30, tokenIndex30, depth30
								if buffer[position] != rune('D') {
									goto l26
								}
								position++
							}
						l30:
							depth--
							add(rulePegText, position29)
						}
						if buffer[position] != rune(' ') {
							goto l26
						}
						position++
						{
							add(ruleAction0, position)
						}
						depth--
						add(ruleLogLevel, position28)
					}
					goto l27
				l26:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
				}
			l27:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					{
						position35 := position
						depth++
						{
							position36 := position
							depth++
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l33
							}
							position++
						l37:
							{
								position38, tokenIndex38, depth38 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l38
								}
								position++
								goto l37
							l38:
								position, tokenIndex, depth = position38, tokenIndex38, depth38
							}
							depth--
							add(rulePegText, position36)
						}
						if buffer[position] != rune(' ') {
							goto l33
						}
						position++
					l39:
						{
							position40, tokenIndex40, depth40 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l40
							}
							position++
							goto l39
						l40:
							position, tokenIndex, depth = position40, tokenIndex40, depth40
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleComponent, position35)
					}
					goto l34
				l33:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
				}
			l34:
				{
					position42 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position43 := position
						depth++
						{
							position46 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position48, tokenIndex48, depth48 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l49
										}
										position++
										goto l48
									l49:
										position, tokenIndex, depth = position48, tokenIndex48, depth48
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l48:
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
							add(ruleletterOrDigit, position46)
						}
					l44:
						{
							position45, tokenIndex45, depth45 := position, tokenIndex, depth
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
												goto l45
											}
											position++
										}
									l52:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l45
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l45
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l45
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position50)
							}
							goto l44
						l45:
							position, tokenIndex, depth = position45, tokenIndex45, depth45
						}
						depth--
						add(rulePegText, position43)
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
					add(ruleContext, position42)
				}
				{
					position55 := position
					depth++
					{
						position56 := position
						depth++
						{
							switch buffer[position] {
							case 'c':
								if buffer[position] != rune('c') {
									goto l0
								}
								position++
								if buffer[position] != rune('o') {
									goto l0
								}
								position++
								if buffer[position] != rune('m') {
									goto l0
								}
								position++
								if buffer[position] != rune('m') {
									goto l0
								}
								position++
								if buffer[position] != rune('a') {
									goto l0
								}
								position++
								if buffer[position] != rune('n') {
									goto l0
								}
								position++
								if buffer[position] != rune('d') {
									goto l0
								}
								position++
								break
							case 'g':
								if buffer[position] != rune('g') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								if buffer[position] != rune('t') {
									goto l0
								}
								position++
								if buffer[position] != rune('m') {
									goto l0
								}
								position++
								if buffer[position] != rune('o') {
									goto l0
								}
								position++
								if buffer[position] != rune('r') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								break
							case 'r':
								if buffer[position] != rune('r') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								if buffer[position] != rune('m') {
									goto l0
								}
								position++
								if buffer[position] != rune('o') {
									goto l0
								}
								position++
								if buffer[position] != rune('v') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								break
							case 'u':
								if buffer[position] != rune('u') {
									goto l0
								}
								position++
								if buffer[position] != rune('p') {
									goto l0
								}
								position++
								if buffer[position] != rune('d') {
									goto l0
								}
								position++
								if buffer[position] != rune('a') {
									goto l0
								}
								position++
								if buffer[position] != rune('t') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								break
							case 'i':
								if buffer[position] != rune('i') {
									goto l0
								}
								position++
								if buffer[position] != rune('n') {
									goto l0
								}
								position++
								if buffer[position] != rune('s') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								if buffer[position] != rune('r') {
									goto l0
								}
								position++
								if buffer[position] != rune('t') {
									goto l0
								}
								position++
								break
							default:
								if buffer[position] != rune('q') {
									goto l0
								}
								position++
								if buffer[position] != rune('u') {
									goto l0
								}
								position++
								if buffer[position] != rune('e') {
									goto l0
								}
								position++
								if buffer[position] != rune('r') {
									goto l0
								}
								position++
								if buffer[position] != rune('y') {
									goto l0
								}
								position++
								break
							}
						}

						depth--
						add(rulePegText, position56)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleOp, position55)
				}
				{
					position59 := position
					depth++
					{
						position60 := position
						depth++
						{
							position63 := position
							depth++
							{
								switch buffer[position] {
								case '$':
									if buffer[position] != rune('$') {
										goto l0
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l0
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l0
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l0
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l0
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('A') || c > rune('z') {
										goto l0
									}
									position++
									break
								}
							}

							depth--
							add(rulensChar, position63)
						}
					l61:
						{
							position62, tokenIndex62, depth62 := position, tokenIndex, depth
							{
								position65 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l62
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l62
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l62
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l62
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l62
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l62
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position65)
							}
							goto l61
						l62:
							position, tokenIndex, depth = position62, tokenIndex62, depth62
						}
						depth--
						add(rulePegText, position60)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction4, position)
					}
					depth--
					add(ruleNS, position59)
				}
			l68:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l69
					}
					goto l68
				l69:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
				}
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					{
						position72 := position
						depth++
						if buffer[position] != rune('l') {
							goto l70
						}
						position++
						if buffer[position] != rune('o') {
							goto l70
						}
						position++
						if buffer[position] != rune('c') {
							goto l70
						}
						position++
						if buffer[position] != rune('k') {
							goto l70
						}
						position++
						if buffer[position] != rune('s') {
							goto l70
						}
						position++
						if buffer[position] != rune('(') {
							goto l70
						}
						position++
						if buffer[position] != rune('m') {
							goto l70
						}
						position++
						if buffer[position] != rune('i') {
							goto l70
						}
						position++
						if buffer[position] != rune('c') {
							goto l70
						}
						position++
						if buffer[position] != rune('r') {
							goto l70
						}
						position++
						if buffer[position] != rune('o') {
							goto l70
						}
						position++
						if buffer[position] != rune('s') {
							goto l70
						}
						position++
						if buffer[position] != rune(')') {
							goto l70
						}
						position++
						{
							position73, tokenIndex73, depth73 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l73
							}
							goto l74
						l73:
							position, tokenIndex, depth = position73, tokenIndex73, depth73
						}
					l74:
					l75:
						{
							position76, tokenIndex76, depth76 := position, tokenIndex, depth
							{
								position77 := position
								depth++
								{
									position78 := position
									depth++
									{
										switch buffer[position] {
										case 'R':
											if buffer[position] != rune('R') {
												goto l76
											}
											position++
											break
										case 'r':
											if buffer[position] != rune('r') {
												goto l76
											}
											position++
											break
										default:
											{
												position80, tokenIndex80, depth80 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l81
												}
												position++
												goto l80
											l81:
												position, tokenIndex, depth = position80, tokenIndex80, depth80
												if buffer[position] != rune('W') {
													goto l76
												}
												position++
											}
										l80:
											break
										}
									}

									depth--
									add(rulePegText, position78)
								}
								{
									add(ruleAction5, position)
								}
								if buffer[position] != rune(':') {
									goto l76
								}
								position++
								if !_rules[ruleNumeric]() {
									goto l76
								}
								{
									position83, tokenIndex83, depth83 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l83
									}
									goto l84
								l83:
									position, tokenIndex, depth = position83, tokenIndex83, depth83
								}
							l84:
								{
									add(ruleAction6, position)
								}
								depth--
								add(rulelock, position77)
							}
							goto l75
						l76:
							position, tokenIndex, depth = position76, tokenIndex76, depth76
						}
						depth--
						add(ruleLocks, position72)
					}
					goto l71
				l70:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
				}
			l71:
			l86:
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l87
					}
					goto l86
				l87:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
				}
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					{
						position90 := position
						depth++
						{
							position91 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l88
							}
							position++
						l92:
							{
								position93, tokenIndex93, depth93 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l93
								}
								position++
								goto l92
							l93:
								position, tokenIndex, depth = position93, tokenIndex93, depth93
							}
							depth--
							add(rulePegText, position91)
						}
						if buffer[position] != rune('m') {
							goto l88
						}
						position++
						if buffer[position] != rune('s') {
							goto l88
						}
						position++
						{
							add(ruleAction7, position)
						}
						depth--
						add(ruleDuration, position90)
					}
					goto l89
				l88:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
				}
			l89:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					{
						position97 := position
						depth++
						{
							position98 := position
							depth++
							if !matchDot() {
								goto l95
							}
						l99:
							{
								position100, tokenIndex100, depth100 := position, tokenIndex, depth
								if !matchDot() {
									goto l100
								}
								goto l99
							l100:
								position, tokenIndex, depth = position100, tokenIndex100, depth100
							}
							depth--
							add(rulePegText, position98)
						}
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleextra, position97)
					}
					goto l96
				l95:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
				}
			l96:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if !matchDot() {
						goto l102
					}
					goto l0
				l102:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
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
		/* 5 Op <- <(<((&('c') ('c' 'o' 'm' 'm' 'a' 'n' 'd')) | (&('g') ('g' 'e' 't' 'm' 'o' 'r' 'e')) | (&('r') ('r' 'e' 'm' 'o' 'v' 'e')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('i') ('i' 'n' 's' 'e' 'r' 't')) | (&('q') ('q' 'u' 'e' 'r' 'y')))> ' ' Action3)> */
		nil,
		/* 6 LineField <- <((commandField / planSummaryField / plainField) S?)> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					{
						position112 := position
						depth++
						if buffer[position] != rune('c') {
							goto l111
						}
						position++
						if buffer[position] != rune('o') {
							goto l111
						}
						position++
						if buffer[position] != rune('m') {
							goto l111
						}
						position++
						if buffer[position] != rune('m') {
							goto l111
						}
						position++
						if buffer[position] != rune('a') {
							goto l111
						}
						position++
						if buffer[position] != rune('n') {
							goto l111
						}
						position++
						if buffer[position] != rune('d') {
							goto l111
						}
						position++
						if buffer[position] != rune(':') {
							goto l111
						}
						position++
						if buffer[position] != rune(' ') {
							goto l111
						}
						position++
						{
							position113 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l111
							}
						l114:
							{
								position115, tokenIndex115, depth115 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l115
								}
								goto l114
							l115:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
							}
							depth--
							add(rulePegText, position113)
						}
						{
							add(ruleAction10, position)
						}
						if !_rules[ruleLineValue]() {
							goto l111
						}
						{
							add(ruleAction11, position)
						}
						depth--
						add(rulecommandField, position112)
					}
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					{
						position119 := position
						depth++
						if buffer[position] != rune('p') {
							goto l118
						}
						position++
						if buffer[position] != rune('l') {
							goto l118
						}
						position++
						if buffer[position] != rune('a') {
							goto l118
						}
						position++
						if buffer[position] != rune('n') {
							goto l118
						}
						position++
						if buffer[position] != rune('S') {
							goto l118
						}
						position++
						if buffer[position] != rune('u') {
							goto l118
						}
						position++
						if buffer[position] != rune('m') {
							goto l118
						}
						position++
						if buffer[position] != rune('m') {
							goto l118
						}
						position++
						if buffer[position] != rune('a') {
							goto l118
						}
						position++
						if buffer[position] != rune('r') {
							goto l118
						}
						position++
						if buffer[position] != rune('y') {
							goto l118
						}
						position++
						if buffer[position] != rune(':') {
							goto l118
						}
						position++
						if buffer[position] != rune(' ') {
							goto l118
						}
						position++
						{
							add(ruleAction12, position)
						}
						{
							position121 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l118
							}
						l122:
							{
								position123, tokenIndex123, depth123 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l123
								}
								position++
								if buffer[position] != rune(' ') {
									goto l123
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l123
								}
								goto l122
							l123:
								position, tokenIndex, depth = position123, tokenIndex123, depth123
							}
							depth--
							add(ruleplanSummaryElements, position121)
						}
						{
							add(ruleAction13, position)
						}
						depth--
						add(ruleplanSummaryField, position119)
					}
					goto l110
				l118:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					{
						position125 := position
						depth++
						{
							position126 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l108
							}
						l127:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l128
								}
								goto l127
							l128:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
							}
							depth--
							add(rulePegText, position126)
						}
						if buffer[position] != rune(':') {
							goto l108
						}
						position++
						{
							position129, tokenIndex129, depth129 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l129
							}
							goto l130
						l129:
							position, tokenIndex, depth = position129, tokenIndex129, depth129
						}
					l130:
						{
							add(ruleAction8, position)
						}
						if !_rules[ruleLineValue]() {
							goto l108
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleplainField, position125)
					}
				}
			l110:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l133
					}
					goto l134
				l133:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
				}
			l134:
				depth--
				add(ruleLineField, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 7 NS <- <(<nsChar+> ' ' Action4)> */
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
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				{
					position145 := position
					depth++
					{
						position146 := position
						depth++
						{
							position147, tokenIndex147, depth147 := position, tokenIndex, depth
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
							if buffer[position] != rune('H') {
								goto l148
							}
							position++
							if buffer[position] != rune('A') {
								goto l148
							}
							position++
							if buffer[position] != rune('S') {
								goto l148
							}
							position++
							if buffer[position] != rune('H') {
								goto l148
							}
							position++
							goto l147
						l148:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('C') {
								goto l149
							}
							position++
							if buffer[position] != rune('A') {
								goto l149
							}
							position++
							if buffer[position] != rune('C') {
								goto l149
							}
							position++
							if buffer[position] != rune('H') {
								goto l149
							}
							position++
							if buffer[position] != rune('E') {
								goto l149
							}
							position++
							if buffer[position] != rune('D') {
								goto l149
							}
							position++
							if buffer[position] != rune('_') {
								goto l149
							}
							position++
							if buffer[position] != rune('P') {
								goto l149
							}
							position++
							if buffer[position] != rune('L') {
								goto l149
							}
							position++
							if buffer[position] != rune('A') {
								goto l149
							}
							position++
							if buffer[position] != rune('N') {
								goto l149
							}
							position++
							goto l147
						l149:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('C') {
								goto l150
							}
							position++
							if buffer[position] != rune('O') {
								goto l150
							}
							position++
							if buffer[position] != rune('L') {
								goto l150
							}
							position++
							if buffer[position] != rune('L') {
								goto l150
							}
							position++
							if buffer[position] != rune('S') {
								goto l150
							}
							position++
							if buffer[position] != rune('C') {
								goto l150
							}
							position++
							if buffer[position] != rune('A') {
								goto l150
							}
							position++
							if buffer[position] != rune('N') {
								goto l150
							}
							position++
							goto l147
						l150:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('C') {
								goto l151
							}
							position++
							if buffer[position] != rune('O') {
								goto l151
							}
							position++
							if buffer[position] != rune('U') {
								goto l151
							}
							position++
							if buffer[position] != rune('N') {
								goto l151
							}
							position++
							if buffer[position] != rune('T') {
								goto l151
							}
							position++
							goto l147
						l151:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('D') {
								goto l152
							}
							position++
							if buffer[position] != rune('E') {
								goto l152
							}
							position++
							if buffer[position] != rune('L') {
								goto l152
							}
							position++
							if buffer[position] != rune('E') {
								goto l152
							}
							position++
							if buffer[position] != rune('T') {
								goto l152
							}
							position++
							if buffer[position] != rune('E') {
								goto l152
							}
							position++
							goto l147
						l152:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('G') {
								goto l153
							}
							position++
							if buffer[position] != rune('E') {
								goto l153
							}
							position++
							if buffer[position] != rune('O') {
								goto l153
							}
							position++
							if buffer[position] != rune('_') {
								goto l153
							}
							position++
							if buffer[position] != rune('N') {
								goto l153
							}
							position++
							if buffer[position] != rune('E') {
								goto l153
							}
							position++
							if buffer[position] != rune('A') {
								goto l153
							}
							position++
							if buffer[position] != rune('R') {
								goto l153
							}
							position++
							if buffer[position] != rune('_') {
								goto l153
							}
							position++
							if buffer[position] != rune('2') {
								goto l153
							}
							position++
							if buffer[position] != rune('D') {
								goto l153
							}
							position++
							goto l147
						l153:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('G') {
								goto l154
							}
							position++
							if buffer[position] != rune('E') {
								goto l154
							}
							position++
							if buffer[position] != rune('O') {
								goto l154
							}
							position++
							if buffer[position] != rune('_') {
								goto l154
							}
							position++
							if buffer[position] != rune('N') {
								goto l154
							}
							position++
							if buffer[position] != rune('E') {
								goto l154
							}
							position++
							if buffer[position] != rune('A') {
								goto l154
							}
							position++
							if buffer[position] != rune('R') {
								goto l154
							}
							position++
							if buffer[position] != rune('_') {
								goto l154
							}
							position++
							if buffer[position] != rune('2') {
								goto l154
							}
							position++
							if buffer[position] != rune('D') {
								goto l154
							}
							position++
							if buffer[position] != rune('S') {
								goto l154
							}
							position++
							if buffer[position] != rune('P') {
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
							if buffer[position] != rune('R') {
								goto l154
							}
							position++
							if buffer[position] != rune('E') {
								goto l154
							}
							position++
							goto l147
						l154:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('I') {
								goto l155
							}
							position++
							if buffer[position] != rune('D') {
								goto l155
							}
							position++
							if buffer[position] != rune('H') {
								goto l155
							}
							position++
							if buffer[position] != rune('A') {
								goto l155
							}
							position++
							if buffer[position] != rune('C') {
								goto l155
							}
							position++
							if buffer[position] != rune('K') {
								goto l155
							}
							position++
							goto l147
						l155:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('S') {
								goto l156
							}
							position++
							if buffer[position] != rune('O') {
								goto l156
							}
							position++
							if buffer[position] != rune('R') {
								goto l156
							}
							position++
							if buffer[position] != rune('T') {
								goto l156
							}
							position++
							if buffer[position] != rune('_') {
								goto l156
							}
							position++
							if buffer[position] != rune('M') {
								goto l156
							}
							position++
							if buffer[position] != rune('E') {
								goto l156
							}
							position++
							if buffer[position] != rune('R') {
								goto l156
							}
							position++
							if buffer[position] != rune('G') {
								goto l156
							}
							position++
							if buffer[position] != rune('E') {
								goto l156
							}
							position++
							goto l147
						l156:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('S') {
								goto l157
							}
							position++
							if buffer[position] != rune('H') {
								goto l157
							}
							position++
							if buffer[position] != rune('A') {
								goto l157
							}
							position++
							if buffer[position] != rune('R') {
								goto l157
							}
							position++
							if buffer[position] != rune('D') {
								goto l157
							}
							position++
							if buffer[position] != rune('I') {
								goto l157
							}
							position++
							if buffer[position] != rune('N') {
								goto l157
							}
							position++
							if buffer[position] != rune('G') {
								goto l157
							}
							position++
							if buffer[position] != rune('_') {
								goto l157
							}
							position++
							if buffer[position] != rune('F') {
								goto l157
							}
							position++
							if buffer[position] != rune('I') {
								goto l157
							}
							position++
							if buffer[position] != rune('L') {
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
							if buffer[position] != rune('R') {
								goto l157
							}
							position++
							goto l147
						l157:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('S') {
								goto l158
							}
							position++
							if buffer[position] != rune('K') {
								goto l158
							}
							position++
							if buffer[position] != rune('I') {
								goto l158
							}
							position++
							if buffer[position] != rune('P') {
								goto l158
							}
							position++
							goto l147
						l158:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('S') {
								goto l159
							}
							position++
							if buffer[position] != rune('O') {
								goto l159
							}
							position++
							if buffer[position] != rune('R') {
								goto l159
							}
							position++
							if buffer[position] != rune('T') {
								goto l159
							}
							position++
							goto l147
						l159:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('X') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('B') {
										goto l143
									}
									position++
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									if buffer[position] != rune('L') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									if buffer[position] != rune('_') {
										goto l143
									}
									position++
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									if buffer[position] != rune('R') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('J') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('R') {
										goto l143
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('L') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('_') {
										goto l143
									}
									position++
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									if buffer[position] != rune('L') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('M') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									if buffer[position] != rune('_') {
										goto l143
									}
									position++
									if buffer[position] != rune('M') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('X') {
										goto l143
									}
									position++
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l143
									}
									position++
									if buffer[position] != rune('R') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('P') {
										goto l143
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('H') {
										goto l143
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('F') {
										goto l143
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('I') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('U') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('_') {
										goto l143
									}
									position++
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									if buffer[position] != rune('C') {
										goto l143
									}
									position++
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l143
									}
									position++
									if buffer[position] != rune('N') {
										goto l143
									}
									position++
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									if buffer[position] != rune('_') {
										goto l143
									}
									position++
									if buffer[position] != rune('S') {
										goto l143
									}
									position++
									if buffer[position] != rune('O') {
										goto l143
									}
									position++
									if buffer[position] != rune('R') {
										goto l143
									}
									position++
									if buffer[position] != rune('T') {
										goto l143
									}
									position++
									if buffer[position] != rune('E') {
										goto l143
									}
									position++
									if buffer[position] != rune('D') {
										goto l143
									}
									position++
									break
								}
							}

						}
					l147:
						depth--
						add(ruleplanSummaryStage, position146)
					}
					depth--
					add(rulePegText, position145)
				}
				{
					add(ruleAction14, position)
				}
				{
					position162 := position
					depth++
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l164
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l164
						}
						{
							add(ruleAction15, position)
						}
						goto l163
					l164:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
						{
							add(ruleAction16, position)
						}
					}
				l163:
					depth--
					add(ruleplanSummary, position162)
				}
				depth--
				add(ruleplanSummaryElem, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 16 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 17 planSummary <- <((' ' LineValue Action15) / Action16)> */
		nil,
		/* 18 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l172
					}
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if !_rules[ruleNumeric]() {
						goto l169
					}
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l173
					}
					goto l174
				l173:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
				}
			l174:
				depth--
				add(ruleLineValue, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
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
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l179
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l179
				}
				position++
				depth--
				add(ruledigit2, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 24 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 25 tz <- <('+' [0-9]+)> */
		nil,
		/* 26 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position183, tokenIndex183, depth183 := position, tokenIndex, depth
			{
				position184 := position
				depth++
				{
					position185 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l183
					}
					depth--
					add(rulehour, position185)
				}
				if buffer[position] != rune(':') {
					goto l183
				}
				position++
				{
					position186 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l183
					}
					depth--
					add(ruleminute, position186)
				}
				if buffer[position] != rune(':') {
					goto l183
				}
				position++
				{
					position187 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l183
					}
					depth--
					add(rulesecond, position187)
				}
				if buffer[position] != rune('.') {
					goto l183
				}
				position++
				{
					position188 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l183
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l183
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l183
					}
					position++
					depth--
					add(rulemillisecond, position188)
				}
				depth--
				add(ruletime, position184)
			}
			return true
		l183:
			position, tokenIndex, depth = position183, tokenIndex183, depth183
			return false
		},
		/* 27 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 28 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 29 dayNum <- <digit2?> */
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
			position199, tokenIndex199, depth199 := position, tokenIndex, depth
			{
				position200 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l199
				}
				position++
			l201:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
				}
				depth--
				add(ruleS, position200)
			}
			return true
		l199:
			position, tokenIndex, depth = position199, tokenIndex199, depth199
			return false
		},
		/* 38 Doc <- <('{' Action20 DocElements? '}' Action21)> */
		func() bool {
			position203, tokenIndex203, depth203 := position, tokenIndex, depth
			{
				position204 := position
				depth++
				if buffer[position] != rune('{') {
					goto l203
				}
				position++
				{
					add(ruleAction20, position)
				}
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					{
						position208 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l206
						}
					l209:
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l210
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l210
							}
							goto l209
						l210:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
						}
						depth--
						add(ruleDocElements, position208)
					}
					goto l207
				l206:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
				}
			l207:
				if buffer[position] != rune('}') {
					goto l203
				}
				position++
				{
					add(ruleAction21, position)
				}
				depth--
				add(ruleDoc, position204)
			}
			return true
		l203:
			position, tokenIndex, depth = position203, tokenIndex203, depth203
			return false
		},
		/* 39 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 40 DocElem <- <(S? Field S? Value S? Action22)> */
		func() bool {
			position213, tokenIndex213, depth213 := position, tokenIndex, depth
			{
				position214 := position
				depth++
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l215
					}
					goto l216
				l215:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
				}
			l216:
				{
					position217 := position
					depth++
					{
						position218 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l213
						}
					l219:
						{
							position220, tokenIndex220, depth220 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l220
							}
							goto l219
						l220:
							position, tokenIndex, depth = position220, tokenIndex220, depth220
						}
						depth--
						add(rulePegText, position218)
					}
					if buffer[position] != rune(':') {
						goto l213
					}
					position++
					{
						add(ruleAction26, position)
					}
					depth--
					add(ruleField, position217)
				}
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l222
					}
					goto l223
				l222:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
				}
			l223:
				if !_rules[ruleValue]() {
					goto l213
				}
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l224
					}
					goto l225
				l224:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
				}
			l225:
				{
					add(ruleAction22, position)
				}
				depth--
				add(ruleDocElem, position214)
			}
			return true
		l213:
			position, tokenIndex, depth = position213, tokenIndex213, depth213
			return false
		},
		/* 41 List <- <('[' Action23 ListElements? ']' Action24)> */
		nil,
		/* 42 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 43 ListElem <- <(S? Value S? Action25)> */
		func() bool {
			position229, tokenIndex229, depth229 := position, tokenIndex, depth
			{
				position230 := position
				depth++
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l231
					}
					goto l232
				l231:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
				}
			l232:
				if !_rules[ruleValue]() {
					goto l229
				}
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l233
					}
					goto l234
				l233:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
				}
			l234:
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleListElem, position230)
			}
			return true
		l229:
			position, tokenIndex, depth = position229, tokenIndex229, depth229
			return false
		},
		/* 44 Field <- <(<fieldChar+> ':' Action26)> */
		nil,
		/* 45 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position237, tokenIndex237, depth237 := position, tokenIndex, depth
			{
				position238 := position
				depth++
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					{
						position241 := position
						depth++
						if buffer[position] != rune('n') {
							goto l240
						}
						position++
						if buffer[position] != rune('u') {
							goto l240
						}
						position++
						if buffer[position] != rune('l') {
							goto l240
						}
						position++
						if buffer[position] != rune('l') {
							goto l240
						}
						position++
						{
							add(ruleAction29, position)
						}
						depth--
						add(ruleNull, position241)
					}
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					{
						position244 := position
						depth++
						if buffer[position] != rune('M') {
							goto l243
						}
						position++
						if buffer[position] != rune('i') {
							goto l243
						}
						position++
						if buffer[position] != rune('n') {
							goto l243
						}
						position++
						if buffer[position] != rune('K') {
							goto l243
						}
						position++
						if buffer[position] != rune('e') {
							goto l243
						}
						position++
						if buffer[position] != rune('y') {
							goto l243
						}
						position++
						{
							add(ruleAction38, position)
						}
						depth--
						add(ruleMinKey, position244)
					}
					goto l239
				l243:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					{
						switch buffer[position] {
						case 'M':
							{
								position247 := position
								depth++
								if buffer[position] != rune('M') {
									goto l237
								}
								position++
								if buffer[position] != rune('a') {
									goto l237
								}
								position++
								if buffer[position] != rune('x') {
									goto l237
								}
								position++
								if buffer[position] != rune('K') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('y') {
									goto l237
								}
								position++
								{
									add(ruleAction39, position)
								}
								depth--
								add(ruleMaxKey, position247)
							}
							break
						case 'u':
							{
								position249 := position
								depth++
								if buffer[position] != rune('u') {
									goto l237
								}
								position++
								if buffer[position] != rune('n') {
									goto l237
								}
								position++
								if buffer[position] != rune('d') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('f') {
									goto l237
								}
								position++
								if buffer[position] != rune('i') {
									goto l237
								}
								position++
								if buffer[position] != rune('n') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('d') {
									goto l237
								}
								position++
								{
									add(ruleAction40, position)
								}
								depth--
								add(ruleUndefined, position249)
							}
							break
						case 'N':
							{
								position251 := position
								depth++
								if buffer[position] != rune('N') {
									goto l237
								}
								position++
								if buffer[position] != rune('u') {
									goto l237
								}
								position++
								if buffer[position] != rune('m') {
									goto l237
								}
								position++
								if buffer[position] != rune('b') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('r') {
									goto l237
								}
								position++
								if buffer[position] != rune('L') {
									goto l237
								}
								position++
								if buffer[position] != rune('o') {
									goto l237
								}
								position++
								if buffer[position] != rune('n') {
									goto l237
								}
								position++
								if buffer[position] != rune('g') {
									goto l237
								}
								position++
								if buffer[position] != rune('(') {
									goto l237
								}
								position++
								{
									position252 := position
									depth++
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l255
										}
										position++
										goto l237
									l255:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
									}
									if !matchDot() {
										goto l237
									}
								l253:
									{
										position254, tokenIndex254, depth254 := position, tokenIndex, depth
										{
											position256, tokenIndex256, depth256 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l256
											}
											position++
											goto l254
										l256:
											position, tokenIndex, depth = position256, tokenIndex256, depth256
										}
										if !matchDot() {
											goto l254
										}
										goto l253
									l254:
										position, tokenIndex, depth = position254, tokenIndex254, depth254
									}
									depth--
									add(rulePegText, position252)
								}
								if buffer[position] != rune(')') {
									goto l237
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleNumberLong, position251)
							}
							break
						case '/':
							{
								position258 := position
								depth++
								if buffer[position] != rune('/') {
									goto l237
								}
								position++
								{
									position259 := position
									depth++
									{
										position260 := position
										depth++
										{
											position263 := position
											depth++
											{
												position264, tokenIndex264, depth264 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l264
												}
												position++
												goto l237
											l264:
												position, tokenIndex, depth = position264, tokenIndex264, depth264
											}
											if !matchDot() {
												goto l237
											}
											depth--
											add(ruleregexChar, position263)
										}
									l261:
										{
											position262, tokenIndex262, depth262 := position, tokenIndex, depth
											{
												position265 := position
												depth++
												{
													position266, tokenIndex266, depth266 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l266
													}
													position++
													goto l262
												l266:
													position, tokenIndex, depth = position266, tokenIndex266, depth266
												}
												if !matchDot() {
													goto l262
												}
												depth--
												add(ruleregexChar, position265)
											}
											goto l261
										l262:
											position, tokenIndex, depth = position262, tokenIndex262, depth262
										}
										if buffer[position] != rune('/') {
											goto l237
										}
										position++
									l267:
										{
											position268, tokenIndex268, depth268 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l268
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l268
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l268
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l268
													}
													position++
													break
												}
											}

											goto l267
										l268:
											position, tokenIndex, depth = position268, tokenIndex268, depth268
										}
										depth--
										add(ruleregexBody, position260)
									}
									depth--
									add(rulePegText, position259)
								}
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleRegex, position258)
							}
							break
						case 'T':
							{
								position271 := position
								depth++
								if buffer[position] != rune('T') {
									goto l237
								}
								position++
								if buffer[position] != rune('i') {
									goto l237
								}
								position++
								if buffer[position] != rune('m') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('s') {
									goto l237
								}
								position++
								if buffer[position] != rune('t') {
									goto l237
								}
								position++
								if buffer[position] != rune('a') {
									goto l237
								}
								position++
								if buffer[position] != rune('m') {
									goto l237
								}
								position++
								if buffer[position] != rune('p') {
									goto l237
								}
								position++
								if buffer[position] != rune('(') {
									goto l237
								}
								position++
								{
									position272 := position
									depth++
									{
										position275, tokenIndex275, depth275 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l275
										}
										position++
										goto l237
									l275:
										position, tokenIndex, depth = position275, tokenIndex275, depth275
									}
									if !matchDot() {
										goto l237
									}
								l273:
									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										{
											position276, tokenIndex276, depth276 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l276
											}
											position++
											goto l274
										l276:
											position, tokenIndex, depth = position276, tokenIndex276, depth276
										}
										if !matchDot() {
											goto l274
										}
										goto l273
									l274:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
									}
									depth--
									add(rulePegText, position272)
								}
								if buffer[position] != rune(')') {
									goto l237
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleTimestampVal, position271)
							}
							break
						case 'B':
							{
								position278 := position
								depth++
								if buffer[position] != rune('B') {
									goto l237
								}
								position++
								if buffer[position] != rune('i') {
									goto l237
								}
								position++
								if buffer[position] != rune('n') {
									goto l237
								}
								position++
								if buffer[position] != rune('D') {
									goto l237
								}
								position++
								if buffer[position] != rune('a') {
									goto l237
								}
								position++
								if buffer[position] != rune('t') {
									goto l237
								}
								position++
								if buffer[position] != rune('a') {
									goto l237
								}
								position++
								if buffer[position] != rune('(') {
									goto l237
								}
								position++
								{
									position279 := position
									depth++
									{
										position282, tokenIndex282, depth282 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l282
										}
										position++
										goto l237
									l282:
										position, tokenIndex, depth = position282, tokenIndex282, depth282
									}
									if !matchDot() {
										goto l237
									}
								l280:
									{
										position281, tokenIndex281, depth281 := position, tokenIndex, depth
										{
											position283, tokenIndex283, depth283 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l283
											}
											position++
											goto l281
										l283:
											position, tokenIndex, depth = position283, tokenIndex283, depth283
										}
										if !matchDot() {
											goto l281
										}
										goto l280
									l281:
										position, tokenIndex, depth = position281, tokenIndex281, depth281
									}
									depth--
									add(rulePegText, position279)
								}
								if buffer[position] != rune(')') {
									goto l237
								}
								position++
								{
									add(ruleAction34, position)
								}
								depth--
								add(ruleBinData, position278)
							}
							break
						case 'D', 'n':
							{
								position285 := position
								depth++
								{
									position286, tokenIndex286, depth286 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l286
									}
									position++
									if buffer[position] != rune('e') {
										goto l286
									}
									position++
									if buffer[position] != rune('w') {
										goto l286
									}
									position++
									if buffer[position] != rune(' ') {
										goto l286
									}
									position++
									goto l287
								l286:
									position, tokenIndex, depth = position286, tokenIndex286, depth286
								}
							l287:
								if buffer[position] != rune('D') {
									goto l237
								}
								position++
								if buffer[position] != rune('a') {
									goto l237
								}
								position++
								if buffer[position] != rune('t') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('(') {
									goto l237
								}
								position++
								{
									position288 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l237
									}
									position++
								l289:
									{
										position290, tokenIndex290, depth290 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l290
										}
										position++
										goto l289
									l290:
										position, tokenIndex, depth = position290, tokenIndex290, depth290
									}
									depth--
									add(rulePegText, position288)
								}
								if buffer[position] != rune(')') {
									goto l237
								}
								position++
								{
									add(ruleAction32, position)
								}
								depth--
								add(ruleDate, position285)
							}
							break
						case 'O':
							{
								position292 := position
								depth++
								if buffer[position] != rune('O') {
									goto l237
								}
								position++
								if buffer[position] != rune('b') {
									goto l237
								}
								position++
								if buffer[position] != rune('j') {
									goto l237
								}
								position++
								if buffer[position] != rune('e') {
									goto l237
								}
								position++
								if buffer[position] != rune('c') {
									goto l237
								}
								position++
								if buffer[position] != rune('t') {
									goto l237
								}
								position++
								if buffer[position] != rune('I') {
									goto l237
								}
								position++
								if buffer[position] != rune('d') {
									goto l237
								}
								position++
								if buffer[position] != rune('(') {
									goto l237
								}
								position++
								if buffer[position] != rune('"') {
									goto l237
								}
								position++
								{
									position293 := position
									depth++
								l294:
									{
										position295, tokenIndex295, depth295 := position, tokenIndex, depth
										{
											position296 := position
											depth++
											{
												position297, tokenIndex297, depth297 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l298
												}
												position++
												goto l297
											l298:
												position, tokenIndex, depth = position297, tokenIndex297, depth297
												{
													position299, tokenIndex299, depth299 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l300
													}
													position++
													goto l299
												l300:
													position, tokenIndex, depth = position299, tokenIndex299, depth299
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l295
													}
													position++
												}
											l299:
											}
										l297:
											depth--
											add(rulehexChar, position296)
										}
										goto l294
									l295:
										position, tokenIndex, depth = position295, tokenIndex295, depth295
									}
									depth--
									add(rulePegText, position293)
								}
								if buffer[position] != rune('"') {
									goto l237
								}
								position++
								if buffer[position] != rune(')') {
									goto l237
								}
								position++
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleObjectID, position292)
							}
							break
						case '"':
							{
								position302 := position
								depth++
								if buffer[position] != rune('"') {
									goto l237
								}
								position++
								{
									position303 := position
									depth++
								l304:
									{
										position305, tokenIndex305, depth305 := position, tokenIndex, depth
										{
											position306 := position
											depth++
											{
												position307, tokenIndex307, depth307 := position, tokenIndex, depth
												{
													position309, tokenIndex309, depth309 := position, tokenIndex, depth
													{
														position310, tokenIndex310, depth310 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l311
														}
														position++
														goto l310
													l311:
														position, tokenIndex, depth = position310, tokenIndex310, depth310
														if buffer[position] != rune('\\') {
															goto l309
														}
														position++
													}
												l310:
													goto l308
												l309:
													position, tokenIndex, depth = position309, tokenIndex309, depth309
												}
												if !matchDot() {
													goto l308
												}
												goto l307
											l308:
												position, tokenIndex, depth = position307, tokenIndex307, depth307
												if buffer[position] != rune('\\') {
													goto l305
												}
												position++
												{
													position312, tokenIndex312, depth312 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l313
													}
													position++
													goto l312
												l313:
													position, tokenIndex, depth = position312, tokenIndex312, depth312
													if buffer[position] != rune('\\') {
														goto l305
													}
													position++
												}
											l312:
											}
										l307:
											depth--
											add(rulestringChar, position306)
										}
										goto l304
									l305:
										position, tokenIndex, depth = position305, tokenIndex305, depth305
									}
									depth--
									add(rulePegText, position303)
								}
								if buffer[position] != rune('"') {
									goto l237
								}
								position++
								{
									add(ruleAction28, position)
								}
								depth--
								add(ruleString, position302)
							}
							break
						case 'f', 't':
							{
								position315 := position
								depth++
								{
									position316, tokenIndex316, depth316 := position, tokenIndex, depth
									{
										position318 := position
										depth++
										if buffer[position] != rune('t') {
											goto l317
										}
										position++
										if buffer[position] != rune('r') {
											goto l317
										}
										position++
										if buffer[position] != rune('u') {
											goto l317
										}
										position++
										if buffer[position] != rune('e') {
											goto l317
										}
										position++
										{
											add(ruleAction30, position)
										}
										depth--
										add(ruleTrue, position318)
									}
									goto l316
								l317:
									position, tokenIndex, depth = position316, tokenIndex316, depth316
									{
										position320 := position
										depth++
										if buffer[position] != rune('f') {
											goto l237
										}
										position++
										if buffer[position] != rune('a') {
											goto l237
										}
										position++
										if buffer[position] != rune('l') {
											goto l237
										}
										position++
										if buffer[position] != rune('s') {
											goto l237
										}
										position++
										if buffer[position] != rune('e') {
											goto l237
										}
										position++
										{
											add(ruleAction31, position)
										}
										depth--
										add(ruleFalse, position320)
									}
								}
							l316:
								depth--
								add(ruleBoolean, position315)
							}
							break
						case '[':
							{
								position322 := position
								depth++
								if buffer[position] != rune('[') {
									goto l237
								}
								position++
								{
									add(ruleAction23, position)
								}
								{
									position324, tokenIndex324, depth324 := position, tokenIndex, depth
									{
										position326 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l324
										}
									l327:
										{
											position328, tokenIndex328, depth328 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l328
											}
											position++
											if !_rules[ruleListElem]() {
												goto l328
											}
											goto l327
										l328:
											position, tokenIndex, depth = position328, tokenIndex328, depth328
										}
										depth--
										add(ruleListElements, position326)
									}
									goto l325
								l324:
									position, tokenIndex, depth = position324, tokenIndex324, depth324
								}
							l325:
								if buffer[position] != rune(']') {
									goto l237
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleList, position322)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l237
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l237
							}
							break
						}
					}

				}
			l239:
				depth--
				add(ruleValue, position238)
			}
			return true
		l237:
			position, tokenIndex, depth = position237, tokenIndex237, depth237
			return false
		},
		/* 46 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action27)> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				{
					position332 := position
					depth++
					{
						position333, tokenIndex333, depth333 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l333
						}
						position++
						goto l334
					l333:
						position, tokenIndex, depth = position333, tokenIndex333, depth333
					}
				l334:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l330
					}
					position++
				l335:
					{
						position336, tokenIndex336, depth336 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l336
						}
						position++
						goto l335
					l336:
						position, tokenIndex, depth = position336, tokenIndex336, depth336
					}
					{
						position337, tokenIndex337, depth337 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l337
						}
						position++
						goto l338
					l337:
						position, tokenIndex, depth = position337, tokenIndex337, depth337
					}
				l338:
				l339:
					{
						position340, tokenIndex340, depth340 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l340
						}
						position++
						goto l339
					l340:
						position, tokenIndex, depth = position340, tokenIndex340, depth340
					}
					depth--
					add(rulePegText, position332)
				}
				{
					add(ruleAction27, position)
				}
				depth--
				add(ruleNumeric, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
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
			position360, tokenIndex360, depth360 := position, tokenIndex, depth
			{
				position361 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l360
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l360
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l360
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l360
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l360
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l360
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l360
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position361)
			}
			return true
		l360:
			position, tokenIndex, depth = position360, tokenIndex360, depth360
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
