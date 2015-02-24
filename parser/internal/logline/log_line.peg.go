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
	ruleThread
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

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"MongoLogLine",
	"Timestamp",
	"Thread",
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
	rules  [103]func() bool
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
			p.SetField("thread", buffer[begin:end])
		case ruleAction1:
			p.SetField("op", buffer[begin:end])
		case ruleAction2:
			p.SetField("ns", buffer[begin:end])
		case ruleAction3:
			p.SetField("duration_ms", buffer[begin:end])
		case ruleAction4:
			p.StartField(buffer[begin:end])
		case ruleAction5:
			p.EndField()
		case ruleAction6:
			p.SetField("commandType", buffer[begin:end])
			p.StartField("command")
		case ruleAction7:
			p.EndField()
		case ruleAction8:
			p.StartField("planSummary")
			p.PushList()
		case ruleAction9:
			p.EndField()
		case ruleAction10:
			p.PushMap()
			p.PushField(buffer[begin:end])
		case ruleAction11:
			p.SetMapValue()
			p.SetListValue()
		case ruleAction12:
			p.PushValue(1)
			p.SetMapValue()
			p.SetListValue()
		case ruleAction13:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction14:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction15:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction16:
			p.PushMap()
		case ruleAction17:
			p.PopMap()
		case ruleAction18:
			p.SetMapValue()
		case ruleAction19:
			p.PushList()
		case ruleAction20:
			p.PopList()
		case ruleAction21:
			p.SetListValue()
		case ruleAction22:
			p.PushField(buffer[begin:end])
		case ruleAction23:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction24:
			p.PushValue(buffer[begin:end])
		case ruleAction25:
			p.PushValue(nil)
		case ruleAction26:
			p.PushValue(true)
		case ruleAction27:
			p.PushValue(false)
		case ruleAction28:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction29:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction30:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction31:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction32:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction33:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction34:
			p.PushValue(p.Minkey())
		case ruleAction35:
			p.PushValue(p.Maxkey())
		case ruleAction36:
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
		/* 0 MongoLogLine <- <(Timestamp Thread Op NS LineField* Locks? LineField* Duration? extra? !.)> */
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
								add(ruleAction13, position)
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
								add(ruleAction14, position)
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
					position26 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position27 := position
						depth++
						{
							position30 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position32, tokenIndex32, depth32 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l33
										}
										position++
										goto l32
									l33:
										position, tokenIndex, depth = position32, tokenIndex32, depth32
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l32:
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
							add(ruleletterOrDigit, position30)
						}
					l28:
						{
							position29, tokenIndex29, depth29 := position, tokenIndex, depth
							{
								position34 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position36, tokenIndex36, depth36 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l37
											}
											position++
											goto l36
										l37:
											position, tokenIndex, depth = position36, tokenIndex36, depth36
											if buffer[position] != rune('$') {
												goto l29
											}
											position++
										}
									l36:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l29
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l29
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l29
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position34)
							}
							goto l28
						l29:
							position, tokenIndex, depth = position29, tokenIndex29, depth29
						}
						depth--
						add(rulePegText, position27)
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
						add(ruleAction0, position)
					}
					depth--
					add(ruleThread, position26)
				}
				{
					position39 := position
					depth++
					{
						position40 := position
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
						add(rulePegText, position40)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction1, position)
					}
					depth--
					add(ruleOp, position39)
				}
				{
					position43 := position
					depth++
					{
						position44 := position
						depth++
						{
							position47 := position
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
							add(rulensChar, position47)
						}
					l45:
						{
							position46, tokenIndex46, depth46 := position, tokenIndex, depth
							{
								position49 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l46
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l46
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l46
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l46
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l46
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l46
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position49)
							}
							goto l45
						l46:
							position, tokenIndex, depth = position46, tokenIndex46, depth46
						}
						depth--
						add(rulePegText, position44)
					}
					if buffer[position] != rune(' ') {
						goto l0
					}
					position++
					{
						add(ruleAction2, position)
					}
					depth--
					add(ruleNS, position43)
				}
			l52:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l53
					}
					goto l52
				l53:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
				}
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					{
						position56 := position
						depth++
						if buffer[position] != rune('l') {
							goto l54
						}
						position++
						if buffer[position] != rune('o') {
							goto l54
						}
						position++
						if buffer[position] != rune('c') {
							goto l54
						}
						position++
						if buffer[position] != rune('k') {
							goto l54
						}
						position++
						if buffer[position] != rune('s') {
							goto l54
						}
						position++
						if buffer[position] != rune('(') {
							goto l54
						}
						position++
						if buffer[position] != rune('m') {
							goto l54
						}
						position++
						if buffer[position] != rune('i') {
							goto l54
						}
						position++
						if buffer[position] != rune('c') {
							goto l54
						}
						position++
						if buffer[position] != rune('r') {
							goto l54
						}
						position++
						if buffer[position] != rune('o') {
							goto l54
						}
						position++
						if buffer[position] != rune('s') {
							goto l54
						}
						position++
						if buffer[position] != rune(')') {
							goto l54
						}
						position++
						{
							position57, tokenIndex57, depth57 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l57
							}
							goto l58
						l57:
							position, tokenIndex, depth = position57, tokenIndex57, depth57
						}
					l58:
					l59:
						{
							position60, tokenIndex60, depth60 := position, tokenIndex, depth
							{
								position61 := position
								depth++
								{
									switch buffer[position] {
									case 'R':
										if buffer[position] != rune('R') {
											goto l60
										}
										position++
										break
									case 'r':
										if buffer[position] != rune('r') {
											goto l60
										}
										position++
										break
									default:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('w') {
												goto l64
											}
											position++
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('W') {
												goto l60
											}
											position++
										}
									l63:
										break
									}
								}

								if buffer[position] != rune(':') {
									goto l60
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l60
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
								{
									position67, tokenIndex67, depth67 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l67
									}
									goto l68
								l67:
									position, tokenIndex, depth = position67, tokenIndex67, depth67
								}
							l68:
								depth--
								add(rulelock, position61)
							}
							goto l59
						l60:
							position, tokenIndex, depth = position60, tokenIndex60, depth60
						}
						depth--
						add(ruleLocks, position56)
					}
					goto l55
				l54:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
				}
			l55:
			l69:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l70
					}
					goto l69
				l70:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
				}
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					{
						position73 := position
						depth++
						{
							position74 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l71
							}
							position++
						l75:
							{
								position76, tokenIndex76, depth76 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l76
								}
								position++
								goto l75
							l76:
								position, tokenIndex, depth = position76, tokenIndex76, depth76
							}
							depth--
							add(rulePegText, position74)
						}
						if buffer[position] != rune('m') {
							goto l71
						}
						position++
						if buffer[position] != rune('s') {
							goto l71
						}
						position++
						{
							add(ruleAction3, position)
						}
						depth--
						add(ruleDuration, position73)
					}
					goto l72
				l71:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
				}
			l72:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					{
						position80 := position
						depth++
						{
							position81 := position
							depth++
							if !matchDot() {
								goto l78
							}
						l82:
							{
								position83, tokenIndex83, depth83 := position, tokenIndex, depth
								if !matchDot() {
									goto l83
								}
								goto l82
							l83:
								position, tokenIndex, depth = position83, tokenIndex83, depth83
							}
							depth--
							add(rulePegText, position81)
						}
						{
							add(ruleAction15, position)
						}
						depth--
						add(ruleextra, position80)
					}
					goto l79
				l78:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
				}
			l79:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if !matchDot() {
						goto l85
					}
					goto l0
				l85:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
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
		/* 2 Thread <- <('[' <letterOrDigit+> ']' ' ' Action0)> */
		nil,
		/* 3 Op <- <(<((&('c') ('c' 'o' 'm' 'm' 'a' 'n' 'd')) | (&('g') ('g' 'e' 't' 'm' 'o' 'r' 'e')) | (&('r') ('r' 'e' 'm' 'o' 'v' 'e')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('i') ('i' 'n' 's' 'e' 'r' 't')) | (&('q') ('q' 'u' 'e' 'r' 'y')))> ' ' Action1)> */
		nil,
		/* 4 LineField <- <((commandField / planSummaryField / plainField) S?)> */
		func() bool {
			position89, tokenIndex89, depth89 := position, tokenIndex, depth
			{
				position90 := position
				depth++
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					{
						position93 := position
						depth++
						if buffer[position] != rune('c') {
							goto l92
						}
						position++
						if buffer[position] != rune('o') {
							goto l92
						}
						position++
						if buffer[position] != rune('m') {
							goto l92
						}
						position++
						if buffer[position] != rune('m') {
							goto l92
						}
						position++
						if buffer[position] != rune('a') {
							goto l92
						}
						position++
						if buffer[position] != rune('n') {
							goto l92
						}
						position++
						if buffer[position] != rune('d') {
							goto l92
						}
						position++
						if buffer[position] != rune(':') {
							goto l92
						}
						position++
						if buffer[position] != rune(' ') {
							goto l92
						}
						position++
						{
							position94 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l92
							}
						l95:
							{
								position96, tokenIndex96, depth96 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l96
								}
								goto l95
							l96:
								position, tokenIndex, depth = position96, tokenIndex96, depth96
							}
							depth--
							add(rulePegText, position94)
						}
						{
							add(ruleAction6, position)
						}
						if !_rules[ruleLineValue]() {
							goto l92
						}
						{
							add(ruleAction7, position)
						}
						depth--
						add(rulecommandField, position93)
					}
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					{
						position100 := position
						depth++
						if buffer[position] != rune('p') {
							goto l99
						}
						position++
						if buffer[position] != rune('l') {
							goto l99
						}
						position++
						if buffer[position] != rune('a') {
							goto l99
						}
						position++
						if buffer[position] != rune('n') {
							goto l99
						}
						position++
						if buffer[position] != rune('S') {
							goto l99
						}
						position++
						if buffer[position] != rune('u') {
							goto l99
						}
						position++
						if buffer[position] != rune('m') {
							goto l99
						}
						position++
						if buffer[position] != rune('m') {
							goto l99
						}
						position++
						if buffer[position] != rune('a') {
							goto l99
						}
						position++
						if buffer[position] != rune('r') {
							goto l99
						}
						position++
						if buffer[position] != rune('y') {
							goto l99
						}
						position++
						if buffer[position] != rune(':') {
							goto l99
						}
						position++
						if buffer[position] != rune(' ') {
							goto l99
						}
						position++
						{
							add(ruleAction8, position)
						}
						{
							position102 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l99
							}
						l103:
							{
								position104, tokenIndex104, depth104 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l104
								}
								position++
								if buffer[position] != rune(' ') {
									goto l104
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l104
								}
								goto l103
							l104:
								position, tokenIndex, depth = position104, tokenIndex104, depth104
							}
							depth--
							add(ruleplanSummaryElements, position102)
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleplanSummaryField, position100)
					}
					goto l91
				l99:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					{
						position106 := position
						depth++
						{
							position107 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l89
							}
						l108:
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l109
								}
								goto l108
							l109:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
							}
							depth--
							add(rulePegText, position107)
						}
						if buffer[position] != rune(':') {
							goto l89
						}
						position++
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
						{
							add(ruleAction4, position)
						}
						if !_rules[ruleLineValue]() {
							goto l89
						}
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleplainField, position106)
					}
				}
			l91:
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l114
					}
					goto l115
				l114:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
				}
			l115:
				depth--
				add(ruleLineField, position90)
			}
			return true
		l89:
			position, tokenIndex, depth = position89, tokenIndex89, depth89
			return false
		},
		/* 5 NS <- <(<nsChar+> ' ' Action2)> */
		nil,
		/* 6 Locks <- <('l' 'o' 'c' 'k' 's' '(' 'm' 'i' 'c' 'r' 'o' 's' ')' S? lock*)> */
		nil,
		/* 7 lock <- <(((&('R') 'R') | (&('r') 'r') | (&('W' | 'w') ('w' / 'W'))) ':' [0-9]+ S?)> */
		nil,
		/* 8 Duration <- <(<[0-9]+> ('m' 's') Action3)> */
		nil,
		/* 9 plainField <- <(<fieldChar+> ':' S? Action4 LineValue Action5)> */
		nil,
		/* 10 commandField <- <('c' 'o' 'm' 'm' 'a' 'n' 'd' ':' ' ' <fieldChar+> Action6 LineValue Action7)> */
		nil,
		/* 11 planSummaryField <- <('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y' ':' ' ' Action8 planSummaryElements Action9)> */
		nil,
		/* 12 planSummaryElements <- <(planSummaryElem (',' ' ' planSummaryElem)*)> */
		nil,
		/* 13 planSummaryElem <- <(<planSummaryStage> Action10 planSummary)> */
		func() bool {
			position124, tokenIndex124, depth124 := position, tokenIndex, depth
			{
				position125 := position
				depth++
				{
					position126 := position
					depth++
					{
						position127 := position
						depth++
						{
							position128, tokenIndex128, depth128 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l129
							}
							position++
							if buffer[position] != rune('N') {
								goto l129
							}
							position++
							if buffer[position] != rune('D') {
								goto l129
							}
							position++
							if buffer[position] != rune('_') {
								goto l129
							}
							position++
							if buffer[position] != rune('H') {
								goto l129
							}
							position++
							if buffer[position] != rune('A') {
								goto l129
							}
							position++
							if buffer[position] != rune('S') {
								goto l129
							}
							position++
							if buffer[position] != rune('H') {
								goto l129
							}
							position++
							goto l128
						l129:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('C') {
								goto l130
							}
							position++
							if buffer[position] != rune('A') {
								goto l130
							}
							position++
							if buffer[position] != rune('C') {
								goto l130
							}
							position++
							if buffer[position] != rune('H') {
								goto l130
							}
							position++
							if buffer[position] != rune('E') {
								goto l130
							}
							position++
							if buffer[position] != rune('D') {
								goto l130
							}
							position++
							if buffer[position] != rune('_') {
								goto l130
							}
							position++
							if buffer[position] != rune('P') {
								goto l130
							}
							position++
							if buffer[position] != rune('L') {
								goto l130
							}
							position++
							if buffer[position] != rune('A') {
								goto l130
							}
							position++
							if buffer[position] != rune('N') {
								goto l130
							}
							position++
							goto l128
						l130:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('C') {
								goto l131
							}
							position++
							if buffer[position] != rune('O') {
								goto l131
							}
							position++
							if buffer[position] != rune('L') {
								goto l131
							}
							position++
							if buffer[position] != rune('L') {
								goto l131
							}
							position++
							if buffer[position] != rune('S') {
								goto l131
							}
							position++
							if buffer[position] != rune('C') {
								goto l131
							}
							position++
							if buffer[position] != rune('A') {
								goto l131
							}
							position++
							if buffer[position] != rune('N') {
								goto l131
							}
							position++
							goto l128
						l131:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('C') {
								goto l132
							}
							position++
							if buffer[position] != rune('O') {
								goto l132
							}
							position++
							if buffer[position] != rune('U') {
								goto l132
							}
							position++
							if buffer[position] != rune('N') {
								goto l132
							}
							position++
							if buffer[position] != rune('T') {
								goto l132
							}
							position++
							goto l128
						l132:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('D') {
								goto l133
							}
							position++
							if buffer[position] != rune('E') {
								goto l133
							}
							position++
							if buffer[position] != rune('L') {
								goto l133
							}
							position++
							if buffer[position] != rune('E') {
								goto l133
							}
							position++
							if buffer[position] != rune('T') {
								goto l133
							}
							position++
							if buffer[position] != rune('E') {
								goto l133
							}
							position++
							goto l128
						l133:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('G') {
								goto l134
							}
							position++
							if buffer[position] != rune('E') {
								goto l134
							}
							position++
							if buffer[position] != rune('O') {
								goto l134
							}
							position++
							if buffer[position] != rune('_') {
								goto l134
							}
							position++
							if buffer[position] != rune('N') {
								goto l134
							}
							position++
							if buffer[position] != rune('E') {
								goto l134
							}
							position++
							if buffer[position] != rune('A') {
								goto l134
							}
							position++
							if buffer[position] != rune('R') {
								goto l134
							}
							position++
							if buffer[position] != rune('_') {
								goto l134
							}
							position++
							if buffer[position] != rune('2') {
								goto l134
							}
							position++
							if buffer[position] != rune('D') {
								goto l134
							}
							position++
							goto l128
						l134:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('G') {
								goto l135
							}
							position++
							if buffer[position] != rune('E') {
								goto l135
							}
							position++
							if buffer[position] != rune('O') {
								goto l135
							}
							position++
							if buffer[position] != rune('_') {
								goto l135
							}
							position++
							if buffer[position] != rune('N') {
								goto l135
							}
							position++
							if buffer[position] != rune('E') {
								goto l135
							}
							position++
							if buffer[position] != rune('A') {
								goto l135
							}
							position++
							if buffer[position] != rune('R') {
								goto l135
							}
							position++
							if buffer[position] != rune('_') {
								goto l135
							}
							position++
							if buffer[position] != rune('2') {
								goto l135
							}
							position++
							if buffer[position] != rune('D') {
								goto l135
							}
							position++
							if buffer[position] != rune('S') {
								goto l135
							}
							position++
							if buffer[position] != rune('P') {
								goto l135
							}
							position++
							if buffer[position] != rune('H') {
								goto l135
							}
							position++
							if buffer[position] != rune('E') {
								goto l135
							}
							position++
							if buffer[position] != rune('R') {
								goto l135
							}
							position++
							if buffer[position] != rune('E') {
								goto l135
							}
							position++
							goto l128
						l135:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('I') {
								goto l136
							}
							position++
							if buffer[position] != rune('D') {
								goto l136
							}
							position++
							if buffer[position] != rune('H') {
								goto l136
							}
							position++
							if buffer[position] != rune('A') {
								goto l136
							}
							position++
							if buffer[position] != rune('C') {
								goto l136
							}
							position++
							if buffer[position] != rune('K') {
								goto l136
							}
							position++
							goto l128
						l136:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('S') {
								goto l137
							}
							position++
							if buffer[position] != rune('O') {
								goto l137
							}
							position++
							if buffer[position] != rune('R') {
								goto l137
							}
							position++
							if buffer[position] != rune('T') {
								goto l137
							}
							position++
							if buffer[position] != rune('_') {
								goto l137
							}
							position++
							if buffer[position] != rune('M') {
								goto l137
							}
							position++
							if buffer[position] != rune('E') {
								goto l137
							}
							position++
							if buffer[position] != rune('R') {
								goto l137
							}
							position++
							if buffer[position] != rune('G') {
								goto l137
							}
							position++
							if buffer[position] != rune('E') {
								goto l137
							}
							position++
							goto l128
						l137:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('S') {
								goto l138
							}
							position++
							if buffer[position] != rune('H') {
								goto l138
							}
							position++
							if buffer[position] != rune('A') {
								goto l138
							}
							position++
							if buffer[position] != rune('R') {
								goto l138
							}
							position++
							if buffer[position] != rune('D') {
								goto l138
							}
							position++
							if buffer[position] != rune('I') {
								goto l138
							}
							position++
							if buffer[position] != rune('N') {
								goto l138
							}
							position++
							if buffer[position] != rune('G') {
								goto l138
							}
							position++
							if buffer[position] != rune('_') {
								goto l138
							}
							position++
							if buffer[position] != rune('F') {
								goto l138
							}
							position++
							if buffer[position] != rune('I') {
								goto l138
							}
							position++
							if buffer[position] != rune('L') {
								goto l138
							}
							position++
							if buffer[position] != rune('T') {
								goto l138
							}
							position++
							if buffer[position] != rune('E') {
								goto l138
							}
							position++
							if buffer[position] != rune('R') {
								goto l138
							}
							position++
							goto l128
						l138:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('S') {
								goto l139
							}
							position++
							if buffer[position] != rune('K') {
								goto l139
							}
							position++
							if buffer[position] != rune('I') {
								goto l139
							}
							position++
							if buffer[position] != rune('P') {
								goto l139
							}
							position++
							goto l128
						l139:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('S') {
								goto l140
							}
							position++
							if buffer[position] != rune('O') {
								goto l140
							}
							position++
							if buffer[position] != rune('R') {
								goto l140
							}
							position++
							if buffer[position] != rune('T') {
								goto l140
							}
							position++
							goto l128
						l140:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('X') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('B') {
										goto l124
									}
									position++
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									if buffer[position] != rune('L') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									if buffer[position] != rune('_') {
										goto l124
									}
									position++
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									if buffer[position] != rune('R') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('J') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('R') {
										goto l124
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('L') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('_') {
										goto l124
									}
									position++
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									if buffer[position] != rune('L') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('M') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									if buffer[position] != rune('_') {
										goto l124
									}
									position++
									if buffer[position] != rune('M') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('X') {
										goto l124
									}
									position++
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l124
									}
									position++
									if buffer[position] != rune('R') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('P') {
										goto l124
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('H') {
										goto l124
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('F') {
										goto l124
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('I') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('U') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('_') {
										goto l124
									}
									position++
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									if buffer[position] != rune('C') {
										goto l124
									}
									position++
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l124
									}
									position++
									if buffer[position] != rune('N') {
										goto l124
									}
									position++
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									if buffer[position] != rune('_') {
										goto l124
									}
									position++
									if buffer[position] != rune('S') {
										goto l124
									}
									position++
									if buffer[position] != rune('O') {
										goto l124
									}
									position++
									if buffer[position] != rune('R') {
										goto l124
									}
									position++
									if buffer[position] != rune('T') {
										goto l124
									}
									position++
									if buffer[position] != rune('E') {
										goto l124
									}
									position++
									if buffer[position] != rune('D') {
										goto l124
									}
									position++
									break
								}
							}

						}
					l128:
						depth--
						add(ruleplanSummaryStage, position127)
					}
					depth--
					add(rulePegText, position126)
				}
				{
					add(ruleAction10, position)
				}
				{
					position143 := position
					depth++
					{
						position144, tokenIndex144, depth144 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l145
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l145
						}
						{
							add(ruleAction11, position)
						}
						goto l144
					l145:
						position, tokenIndex, depth = position144, tokenIndex144, depth144
						{
							add(ruleAction12, position)
						}
					}
				l144:
					depth--
					add(ruleplanSummary, position143)
				}
				depth--
				add(ruleplanSummaryElem, position125)
			}
			return true
		l124:
			position, tokenIndex, depth = position124, tokenIndex124, depth124
			return false
		},
		/* 14 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 15 planSummary <- <((' ' LineValue Action11) / Action12)> */
		nil,
		/* 16 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l153
					}
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if !_rules[ruleNumeric]() {
						goto l150
					}
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l154
					}
					goto l155
				l154:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
				}
			l155:
				depth--
				add(ruleLineValue, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
			return false
		},
		/* 17 timestamp24 <- <(<(date ' ' time)> Action13)> */
		nil,
		/* 18 timestamp26 <- <(<datetime26> Action14)> */
		nil,
		/* 19 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 20 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 21 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l160
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l160
				}
				position++
				depth--
				add(ruledigit2, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 22 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 23 tz <- <('+' [0-9]+)> */
		nil,
		/* 24 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				{
					position166 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l164
					}
					depth--
					add(rulehour, position166)
				}
				if buffer[position] != rune(':') {
					goto l164
				}
				position++
				{
					position167 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l164
					}
					depth--
					add(ruleminute, position167)
				}
				if buffer[position] != rune(':') {
					goto l164
				}
				position++
				{
					position168 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l164
					}
					depth--
					add(rulesecond, position168)
				}
				if buffer[position] != rune('.') {
					goto l164
				}
				position++
				{
					position169 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l164
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l164
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l164
					}
					position++
					depth--
					add(rulemillisecond, position169)
				}
				depth--
				add(ruletime, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 25 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 26 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 27 dayNum <- <digit2?> */
		nil,
		/* 28 hour <- <digit2> */
		nil,
		/* 29 minute <- <digit2> */
		nil,
		/* 30 second <- <digit2> */
		nil,
		/* 31 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 32 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 33 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 34 extra <- <(<.+> Action15)> */
		nil,
		/* 35 S <- <' '+> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l180
				}
				position++
			l182:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
				}
				depth--
				add(ruleS, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 36 Doc <- <('{' Action16 DocElements? '}' Action17)> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				if buffer[position] != rune('{') {
					goto l184
				}
				position++
				{
					add(ruleAction16, position)
				}
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					{
						position189 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l187
						}
					l190:
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l191
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l191
							}
							goto l190
						l191:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
						}
						depth--
						add(ruleDocElements, position189)
					}
					goto l188
				l187:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
				}
			l188:
				if buffer[position] != rune('}') {
					goto l184
				}
				position++
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleDoc, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 37 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 38 DocElem <- <(S? Field S? Value S? Action18)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
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
				{
					position198 := position
					depth++
					{
						position199 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l194
						}
					l200:
						{
							position201, tokenIndex201, depth201 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l201
							}
							goto l200
						l201:
							position, tokenIndex, depth = position201, tokenIndex201, depth201
						}
						depth--
						add(rulePegText, position199)
					}
					if buffer[position] != rune(':') {
						goto l194
					}
					position++
					{
						add(ruleAction22, position)
					}
					depth--
					add(ruleField, position198)
				}
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l203
					}
					goto l204
				l203:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
				}
			l204:
				if !_rules[ruleValue]() {
					goto l194
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l205
					}
					goto l206
				l205:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
				}
			l206:
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleDocElem, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 39 List <- <('[' Action19 ListElements? ']' Action20)> */
		nil,
		/* 40 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 41 ListElem <- <(S? Value S? Action21)> */
		func() bool {
			position210, tokenIndex210, depth210 := position, tokenIndex, depth
			{
				position211 := position
				depth++
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l212
					}
					goto l213
				l212:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
				}
			l213:
				if !_rules[ruleValue]() {
					goto l210
				}
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l214
					}
					goto l215
				l214:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
				}
			l215:
				{
					add(ruleAction21, position)
				}
				depth--
				add(ruleListElem, position211)
			}
			return true
		l210:
			position, tokenIndex, depth = position210, tokenIndex210, depth210
			return false
		},
		/* 42 Field <- <(<fieldChar+> ':' Action22)> */
		nil,
		/* 43 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					{
						position222 := position
						depth++
						if buffer[position] != rune('n') {
							goto l221
						}
						position++
						if buffer[position] != rune('u') {
							goto l221
						}
						position++
						if buffer[position] != rune('l') {
							goto l221
						}
						position++
						if buffer[position] != rune('l') {
							goto l221
						}
						position++
						{
							add(ruleAction25, position)
						}
						depth--
						add(ruleNull, position222)
					}
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					{
						position225 := position
						depth++
						if buffer[position] != rune('M') {
							goto l224
						}
						position++
						if buffer[position] != rune('i') {
							goto l224
						}
						position++
						if buffer[position] != rune('n') {
							goto l224
						}
						position++
						if buffer[position] != rune('K') {
							goto l224
						}
						position++
						if buffer[position] != rune('e') {
							goto l224
						}
						position++
						if buffer[position] != rune('y') {
							goto l224
						}
						position++
						{
							add(ruleAction34, position)
						}
						depth--
						add(ruleMinKey, position225)
					}
					goto l220
				l224:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					{
						switch buffer[position] {
						case 'M':
							{
								position228 := position
								depth++
								if buffer[position] != rune('M') {
									goto l218
								}
								position++
								if buffer[position] != rune('a') {
									goto l218
								}
								position++
								if buffer[position] != rune('x') {
									goto l218
								}
								position++
								if buffer[position] != rune('K') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('y') {
									goto l218
								}
								position++
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleMaxKey, position228)
							}
							break
						case 'u':
							{
								position230 := position
								depth++
								if buffer[position] != rune('u') {
									goto l218
								}
								position++
								if buffer[position] != rune('n') {
									goto l218
								}
								position++
								if buffer[position] != rune('d') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('f') {
									goto l218
								}
								position++
								if buffer[position] != rune('i') {
									goto l218
								}
								position++
								if buffer[position] != rune('n') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('d') {
									goto l218
								}
								position++
								{
									add(ruleAction36, position)
								}
								depth--
								add(ruleUndefined, position230)
							}
							break
						case 'N':
							{
								position232 := position
								depth++
								if buffer[position] != rune('N') {
									goto l218
								}
								position++
								if buffer[position] != rune('u') {
									goto l218
								}
								position++
								if buffer[position] != rune('m') {
									goto l218
								}
								position++
								if buffer[position] != rune('b') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('r') {
									goto l218
								}
								position++
								if buffer[position] != rune('L') {
									goto l218
								}
								position++
								if buffer[position] != rune('o') {
									goto l218
								}
								position++
								if buffer[position] != rune('n') {
									goto l218
								}
								position++
								if buffer[position] != rune('g') {
									goto l218
								}
								position++
								if buffer[position] != rune('(') {
									goto l218
								}
								position++
								{
									position233 := position
									depth++
									{
										position236, tokenIndex236, depth236 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l236
										}
										position++
										goto l218
									l236:
										position, tokenIndex, depth = position236, tokenIndex236, depth236
									}
									if !matchDot() {
										goto l218
									}
								l234:
									{
										position235, tokenIndex235, depth235 := position, tokenIndex, depth
										{
											position237, tokenIndex237, depth237 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l237
											}
											position++
											goto l235
										l237:
											position, tokenIndex, depth = position237, tokenIndex237, depth237
										}
										if !matchDot() {
											goto l235
										}
										goto l234
									l235:
										position, tokenIndex, depth = position235, tokenIndex235, depth235
									}
									depth--
									add(rulePegText, position233)
								}
								if buffer[position] != rune(')') {
									goto l218
								}
								position++
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleNumberLong, position232)
							}
							break
						case '/':
							{
								position239 := position
								depth++
								if buffer[position] != rune('/') {
									goto l218
								}
								position++
								{
									position240 := position
									depth++
									{
										position241 := position
										depth++
										{
											position244 := position
											depth++
											{
												position245, tokenIndex245, depth245 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l245
												}
												position++
												goto l218
											l245:
												position, tokenIndex, depth = position245, tokenIndex245, depth245
											}
											if !matchDot() {
												goto l218
											}
											depth--
											add(ruleregexChar, position244)
										}
									l242:
										{
											position243, tokenIndex243, depth243 := position, tokenIndex, depth
											{
												position246 := position
												depth++
												{
													position247, tokenIndex247, depth247 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l247
													}
													position++
													goto l243
												l247:
													position, tokenIndex, depth = position247, tokenIndex247, depth247
												}
												if !matchDot() {
													goto l243
												}
												depth--
												add(ruleregexChar, position246)
											}
											goto l242
										l243:
											position, tokenIndex, depth = position243, tokenIndex243, depth243
										}
										if buffer[position] != rune('/') {
											goto l218
										}
										position++
									l248:
										{
											position249, tokenIndex249, depth249 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l249
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l249
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l249
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l249
													}
													position++
													break
												}
											}

											goto l248
										l249:
											position, tokenIndex, depth = position249, tokenIndex249, depth249
										}
										depth--
										add(ruleregexBody, position241)
									}
									depth--
									add(rulePegText, position240)
								}
								{
									add(ruleAction31, position)
								}
								depth--
								add(ruleRegex, position239)
							}
							break
						case 'T':
							{
								position252 := position
								depth++
								if buffer[position] != rune('T') {
									goto l218
								}
								position++
								if buffer[position] != rune('i') {
									goto l218
								}
								position++
								if buffer[position] != rune('m') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('s') {
									goto l218
								}
								position++
								if buffer[position] != rune('t') {
									goto l218
								}
								position++
								if buffer[position] != rune('a') {
									goto l218
								}
								position++
								if buffer[position] != rune('m') {
									goto l218
								}
								position++
								if buffer[position] != rune('p') {
									goto l218
								}
								position++
								if buffer[position] != rune('(') {
									goto l218
								}
								position++
								{
									position253 := position
									depth++
									{
										position256, tokenIndex256, depth256 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l256
										}
										position++
										goto l218
									l256:
										position, tokenIndex, depth = position256, tokenIndex256, depth256
									}
									if !matchDot() {
										goto l218
									}
								l254:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										{
											position257, tokenIndex257, depth257 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l257
											}
											position++
											goto l255
										l257:
											position, tokenIndex, depth = position257, tokenIndex257, depth257
										}
										if !matchDot() {
											goto l255
										}
										goto l254
									l255:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
									}
									depth--
									add(rulePegText, position253)
								}
								if buffer[position] != rune(')') {
									goto l218
								}
								position++
								{
									add(ruleAction32, position)
								}
								depth--
								add(ruleTimestampVal, position252)
							}
							break
						case 'B':
							{
								position259 := position
								depth++
								if buffer[position] != rune('B') {
									goto l218
								}
								position++
								if buffer[position] != rune('i') {
									goto l218
								}
								position++
								if buffer[position] != rune('n') {
									goto l218
								}
								position++
								if buffer[position] != rune('D') {
									goto l218
								}
								position++
								if buffer[position] != rune('a') {
									goto l218
								}
								position++
								if buffer[position] != rune('t') {
									goto l218
								}
								position++
								if buffer[position] != rune('a') {
									goto l218
								}
								position++
								if buffer[position] != rune('(') {
									goto l218
								}
								position++
								{
									position260 := position
									depth++
									{
										position263, tokenIndex263, depth263 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l263
										}
										position++
										goto l218
									l263:
										position, tokenIndex, depth = position263, tokenIndex263, depth263
									}
									if !matchDot() {
										goto l218
									}
								l261:
									{
										position262, tokenIndex262, depth262 := position, tokenIndex, depth
										{
											position264, tokenIndex264, depth264 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l264
											}
											position++
											goto l262
										l264:
											position, tokenIndex, depth = position264, tokenIndex264, depth264
										}
										if !matchDot() {
											goto l262
										}
										goto l261
									l262:
										position, tokenIndex, depth = position262, tokenIndex262, depth262
									}
									depth--
									add(rulePegText, position260)
								}
								if buffer[position] != rune(')') {
									goto l218
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleBinData, position259)
							}
							break
						case 'D', 'n':
							{
								position266 := position
								depth++
								{
									position267, tokenIndex267, depth267 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l267
									}
									position++
									if buffer[position] != rune('e') {
										goto l267
									}
									position++
									if buffer[position] != rune('w') {
										goto l267
									}
									position++
									if buffer[position] != rune(' ') {
										goto l267
									}
									position++
									goto l268
								l267:
									position, tokenIndex, depth = position267, tokenIndex267, depth267
								}
							l268:
								if buffer[position] != rune('D') {
									goto l218
								}
								position++
								if buffer[position] != rune('a') {
									goto l218
								}
								position++
								if buffer[position] != rune('t') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('(') {
									goto l218
								}
								position++
								{
									position269 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l218
									}
									position++
								l270:
									{
										position271, tokenIndex271, depth271 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l271
										}
										position++
										goto l270
									l271:
										position, tokenIndex, depth = position271, tokenIndex271, depth271
									}
									depth--
									add(rulePegText, position269)
								}
								if buffer[position] != rune(')') {
									goto l218
								}
								position++
								{
									add(ruleAction28, position)
								}
								depth--
								add(ruleDate, position266)
							}
							break
						case 'O':
							{
								position273 := position
								depth++
								if buffer[position] != rune('O') {
									goto l218
								}
								position++
								if buffer[position] != rune('b') {
									goto l218
								}
								position++
								if buffer[position] != rune('j') {
									goto l218
								}
								position++
								if buffer[position] != rune('e') {
									goto l218
								}
								position++
								if buffer[position] != rune('c') {
									goto l218
								}
								position++
								if buffer[position] != rune('t') {
									goto l218
								}
								position++
								if buffer[position] != rune('I') {
									goto l218
								}
								position++
								if buffer[position] != rune('d') {
									goto l218
								}
								position++
								if buffer[position] != rune('(') {
									goto l218
								}
								position++
								if buffer[position] != rune('"') {
									goto l218
								}
								position++
								{
									position274 := position
									depth++
								l275:
									{
										position276, tokenIndex276, depth276 := position, tokenIndex, depth
										{
											position277 := position
											depth++
											{
												position278, tokenIndex278, depth278 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l279
												}
												position++
												goto l278
											l279:
												position, tokenIndex, depth = position278, tokenIndex278, depth278
												{
													position280, tokenIndex280, depth280 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l281
													}
													position++
													goto l280
												l281:
													position, tokenIndex, depth = position280, tokenIndex280, depth280
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l276
													}
													position++
												}
											l280:
											}
										l278:
											depth--
											add(rulehexChar, position277)
										}
										goto l275
									l276:
										position, tokenIndex, depth = position276, tokenIndex276, depth276
									}
									depth--
									add(rulePegText, position274)
								}
								if buffer[position] != rune('"') {
									goto l218
								}
								position++
								if buffer[position] != rune(')') {
									goto l218
								}
								position++
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleObjectID, position273)
							}
							break
						case '"':
							{
								position283 := position
								depth++
								if buffer[position] != rune('"') {
									goto l218
								}
								position++
								{
									position284 := position
									depth++
								l285:
									{
										position286, tokenIndex286, depth286 := position, tokenIndex, depth
										{
											position287 := position
											depth++
											{
												position288, tokenIndex288, depth288 := position, tokenIndex, depth
												{
													position290, tokenIndex290, depth290 := position, tokenIndex, depth
													{
														position291, tokenIndex291, depth291 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l292
														}
														position++
														goto l291
													l292:
														position, tokenIndex, depth = position291, tokenIndex291, depth291
														if buffer[position] != rune('\\') {
															goto l290
														}
														position++
													}
												l291:
													goto l289
												l290:
													position, tokenIndex, depth = position290, tokenIndex290, depth290
												}
												if !matchDot() {
													goto l289
												}
												goto l288
											l289:
												position, tokenIndex, depth = position288, tokenIndex288, depth288
												if buffer[position] != rune('\\') {
													goto l286
												}
												position++
												{
													position293, tokenIndex293, depth293 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l294
													}
													position++
													goto l293
												l294:
													position, tokenIndex, depth = position293, tokenIndex293, depth293
													if buffer[position] != rune('\\') {
														goto l286
													}
													position++
												}
											l293:
											}
										l288:
											depth--
											add(rulestringChar, position287)
										}
										goto l285
									l286:
										position, tokenIndex, depth = position286, tokenIndex286, depth286
									}
									depth--
									add(rulePegText, position284)
								}
								if buffer[position] != rune('"') {
									goto l218
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleString, position283)
							}
							break
						case 'f', 't':
							{
								position296 := position
								depth++
								{
									position297, tokenIndex297, depth297 := position, tokenIndex, depth
									{
										position299 := position
										depth++
										if buffer[position] != rune('t') {
											goto l298
										}
										position++
										if buffer[position] != rune('r') {
											goto l298
										}
										position++
										if buffer[position] != rune('u') {
											goto l298
										}
										position++
										if buffer[position] != rune('e') {
											goto l298
										}
										position++
										{
											add(ruleAction26, position)
										}
										depth--
										add(ruleTrue, position299)
									}
									goto l297
								l298:
									position, tokenIndex, depth = position297, tokenIndex297, depth297
									{
										position301 := position
										depth++
										if buffer[position] != rune('f') {
											goto l218
										}
										position++
										if buffer[position] != rune('a') {
											goto l218
										}
										position++
										if buffer[position] != rune('l') {
											goto l218
										}
										position++
										if buffer[position] != rune('s') {
											goto l218
										}
										position++
										if buffer[position] != rune('e') {
											goto l218
										}
										position++
										{
											add(ruleAction27, position)
										}
										depth--
										add(ruleFalse, position301)
									}
								}
							l297:
								depth--
								add(ruleBoolean, position296)
							}
							break
						case '[':
							{
								position303 := position
								depth++
								if buffer[position] != rune('[') {
									goto l218
								}
								position++
								{
									add(ruleAction19, position)
								}
								{
									position305, tokenIndex305, depth305 := position, tokenIndex, depth
									{
										position307 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l305
										}
									l308:
										{
											position309, tokenIndex309, depth309 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l309
											}
											position++
											if !_rules[ruleListElem]() {
												goto l309
											}
											goto l308
										l309:
											position, tokenIndex, depth = position309, tokenIndex309, depth309
										}
										depth--
										add(ruleListElements, position307)
									}
									goto l306
								l305:
									position, tokenIndex, depth = position305, tokenIndex305, depth305
								}
							l306:
								if buffer[position] != rune(']') {
									goto l218
								}
								position++
								{
									add(ruleAction20, position)
								}
								depth--
								add(ruleList, position303)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l218
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l218
							}
							break
						}
					}

				}
			l220:
				depth--
				add(ruleValue, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 44 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action23)> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
				{
					position313 := position
					depth++
					{
						position314, tokenIndex314, depth314 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l314
						}
						position++
						goto l315
					l314:
						position, tokenIndex, depth = position314, tokenIndex314, depth314
					}
				l315:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l311
					}
					position++
				l316:
					{
						position317, tokenIndex317, depth317 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l317
						}
						position++
						goto l316
					l317:
						position, tokenIndex, depth = position317, tokenIndex317, depth317
					}
					{
						position318, tokenIndex318, depth318 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l318
						}
						position++
						goto l319
					l318:
						position, tokenIndex, depth = position318, tokenIndex318, depth318
					}
				l319:
				l320:
					{
						position321, tokenIndex321, depth321 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l321
						}
						position++
						goto l320
					l321:
						position, tokenIndex, depth = position321, tokenIndex321, depth321
					}
					depth--
					add(rulePegText, position313)
				}
				{
					add(ruleAction23, position)
				}
				depth--
				add(ruleNumeric, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		/* 45 Boolean <- <(True / False)> */
		nil,
		/* 46 String <- <('"' <stringChar*> '"' Action24)> */
		nil,
		/* 47 Null <- <('n' 'u' 'l' 'l' Action25)> */
		nil,
		/* 48 True <- <('t' 'r' 'u' 'e' Action26)> */
		nil,
		/* 49 False <- <('f' 'a' 'l' 's' 'e' Action27)> */
		nil,
		/* 50 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action28)> */
		nil,
		/* 51 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action29)> */
		nil,
		/* 52 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action30)> */
		nil,
		/* 53 Regex <- <('/' <regexBody> Action31)> */
		nil,
		/* 54 TimestampVal <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action32)> */
		nil,
		/* 55 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action33)> */
		nil,
		/* 56 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action34)> */
		nil,
		/* 57 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action35)> */
		nil,
		/* 58 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action36)> */
		nil,
		/* 59 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 60 regexChar <- <(!'/' .)> */
		nil,
		/* 61 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 62 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 63 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l341
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l341
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l341
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l341
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l341
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l341
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l341
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		nil,
		/* 66 Action0 <- <{ p.SetField("thread", buffer[begin:end]) }> */
		nil,
		/* 67 Action1 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 68 Action2 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 69 Action3 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 70 Action4 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 71 Action5 <- <{ p.EndField() }> */
		nil,
		/* 72 Action6 <- <{ p.SetField("commandType", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 73 Action7 <- <{ p.EndField() }> */
		nil,
		/* 74 Action8 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 75 Action9 <- <{ p.EndField()}> */
		nil,
		/* 76 Action10 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 77 Action11 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 78 Action12 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 79 Action13 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 80 Action14 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 81 Action15 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 82 Action16 <- <{ p.PushMap() }> */
		nil,
		/* 83 Action17 <- <{ p.PopMap() }> */
		nil,
		/* 84 Action18 <- <{ p.SetMapValue() }> */
		nil,
		/* 85 Action19 <- <{ p.PushList() }> */
		nil,
		/* 86 Action20 <- <{ p.PopList() }> */
		nil,
		/* 87 Action21 <- <{ p.SetListValue() }> */
		nil,
		/* 88 Action22 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 89 Action23 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 90 Action24 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 91 Action25 <- <{ p.PushValue(nil) }> */
		nil,
		/* 92 Action26 <- <{ p.PushValue(true) }> */
		nil,
		/* 93 Action27 <- <{ p.PushValue(false) }> */
		nil,
		/* 94 Action28 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 95 Action29 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 96 Action30 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 97 Action31 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 98 Action32 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 99 Action33 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 100 Action34 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 101 Action35 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 102 Action36 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
