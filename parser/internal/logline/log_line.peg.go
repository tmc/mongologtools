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
	ruleplainFieldName
	rulecommandFieldName
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
	"plainFieldName",
	"commandFieldName",
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
	rules  [92]func() bool
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
			p.EndField()
		case ruleAction3:
			p.SetField("ns", buffer[begin:end])
		case ruleAction4:
			p.SetField("duration_ms", buffer[begin:end])
		case ruleAction5:
			p.StartField(buffer[begin:end])
		case ruleAction6:
			p.SetField("commandType", buffer[begin:end])
			p.StartField("command")
		case ruleAction7:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction8:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction9:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction10:
			p.PushMap()
		case ruleAction11:
			p.PopMap()
		case ruleAction12:
			p.SetMapValue()
		case ruleAction13:
			p.PushList()
		case ruleAction14:
			p.PopList()
		case ruleAction15:
			p.SetListValue()
		case ruleAction16:
			p.PushField(buffer[begin:end])
		case ruleAction17:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction18:
			p.PushValue(buffer[begin:end])
		case ruleAction19:
			p.PushValue(nil)
		case ruleAction20:
			p.PushValue(true)
		case ruleAction21:
			p.PushValue(false)
		case ruleAction22:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction23:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction24:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction25:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction26:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction27:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction28:
			p.PushValue(p.Minkey())
		case ruleAction29:
			p.PushValue(p.Maxkey())
		case ruleAction30:
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
		/* 0 MongoLogLine <- <(Timestamp ' ' Thread ' ' Op ' ' NS ' ' LineField* Locks? LineField* Duration? extra? !.)> */
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
								add(ruleAction7, position)
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
								add(ruleAction8, position)
							}
							depth--
							add(ruletimestamp26, position14)
						}
					}
				l3:
					depth--
					add(ruleTimestamp, position2)
				}
				if buffer[position] != rune(' ') {
					goto l0
				}
				position++
				{
					position24 := position
					depth++
					if buffer[position] != rune('[') {
						goto l0
					}
					position++
					{
						position25 := position
						depth++
						{
							position28 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position30, tokenIndex30, depth30 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l31
										}
										position++
										goto l30
									l31:
										position, tokenIndex, depth = position30, tokenIndex30, depth30
										if buffer[position] != rune('$') {
											goto l0
										}
										position++
									}
								l30:
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
							add(ruleletterOrDigit, position28)
						}
					l26:
						{
							position27, tokenIndex27, depth27 := position, tokenIndex, depth
							{
								position32 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position34, tokenIndex34, depth34 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l35
											}
											position++
											goto l34
										l35:
											position, tokenIndex, depth = position34, tokenIndex34, depth34
											if buffer[position] != rune('$') {
												goto l27
											}
											position++
										}
									l34:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l27
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l27
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l27
										}
										position++
										break
									}
								}

								depth--
								add(ruleletterOrDigit, position32)
							}
							goto l26
						l27:
							position, tokenIndex, depth = position27, tokenIndex27, depth27
						}
						depth--
						add(rulePegText, position25)
					}
					if buffer[position] != rune(']') {
						goto l0
					}
					position++
					{
						add(ruleAction0, position)
					}
					depth--
					add(ruleThread, position24)
				}
				if buffer[position] != rune(' ') {
					goto l0
				}
				position++
				{
					position37 := position
					depth++
					{
						position38 := position
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
						add(rulePegText, position38)
					}
					{
						add(ruleAction1, position)
					}
					depth--
					add(ruleOp, position37)
				}
				if buffer[position] != rune(' ') {
					goto l0
				}
				position++
				{
					position41 := position
					depth++
					{
						position42 := position
						depth++
						{
							position45 := position
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
							add(rulensChar, position45)
						}
					l43:
						{
							position44, tokenIndex44, depth44 := position, tokenIndex, depth
							{
								position47 := position
								depth++
								{
									switch buffer[position] {
									case '$':
										if buffer[position] != rune('$') {
											goto l44
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l44
										}
										position++
										break
									case '.':
										if buffer[position] != rune('.') {
											goto l44
										}
										position++
										break
									case '-':
										if buffer[position] != rune('-') {
											goto l44
										}
										position++
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l44
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('A') || c > rune('z') {
											goto l44
										}
										position++
										break
									}
								}

								depth--
								add(rulensChar, position47)
							}
							goto l43
						l44:
							position, tokenIndex, depth = position44, tokenIndex44, depth44
						}
						depth--
						add(rulePegText, position42)
					}
					{
						add(ruleAction3, position)
					}
					depth--
					add(ruleNS, position41)
				}
				if buffer[position] != rune(' ') {
					goto l0
				}
				position++
			l50:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l51
					}
					goto l50
				l51:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
				}
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					{
						position54 := position
						depth++
						if buffer[position] != rune('l') {
							goto l52
						}
						position++
						if buffer[position] != rune('o') {
							goto l52
						}
						position++
						if buffer[position] != rune('c') {
							goto l52
						}
						position++
						if buffer[position] != rune('k') {
							goto l52
						}
						position++
						if buffer[position] != rune('s') {
							goto l52
						}
						position++
						if buffer[position] != rune('(') {
							goto l52
						}
						position++
						if buffer[position] != rune('m') {
							goto l52
						}
						position++
						if buffer[position] != rune('i') {
							goto l52
						}
						position++
						if buffer[position] != rune('c') {
							goto l52
						}
						position++
						if buffer[position] != rune('r') {
							goto l52
						}
						position++
						if buffer[position] != rune('o') {
							goto l52
						}
						position++
						if buffer[position] != rune('s') {
							goto l52
						}
						position++
						if buffer[position] != rune(')') {
							goto l52
						}
						position++
						{
							position55, tokenIndex55, depth55 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l55
							}
							goto l56
						l55:
							position, tokenIndex, depth = position55, tokenIndex55, depth55
						}
					l56:
					l57:
						{
							position58, tokenIndex58, depth58 := position, tokenIndex, depth
							{
								position59 := position
								depth++
								{
									switch buffer[position] {
									case 'R':
										if buffer[position] != rune('R') {
											goto l58
										}
										position++
										break
									case 'r':
										if buffer[position] != rune('r') {
											goto l58
										}
										position++
										break
									default:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('w') {
												goto l62
											}
											position++
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('W') {
												goto l58
											}
											position++
										}
									l61:
										break
									}
								}

								if buffer[position] != rune(':') {
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
								{
									position65, tokenIndex65, depth65 := position, tokenIndex, depth
									if !_rules[ruleS]() {
										goto l65
									}
									goto l66
								l65:
									position, tokenIndex, depth = position65, tokenIndex65, depth65
								}
							l66:
								depth--
								add(rulelock, position59)
							}
							goto l57
						l58:
							position, tokenIndex, depth = position58, tokenIndex58, depth58
						}
						depth--
						add(ruleLocks, position54)
					}
					goto l53
				l52:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
				}
			l53:
			l67:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l68
					}
					goto l67
				l68:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
				}
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					{
						position71 := position
						depth++
						{
							position72 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l69
							}
							position++
						l73:
							{
								position74, tokenIndex74, depth74 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l74
								}
								position++
								goto l73
							l74:
								position, tokenIndex, depth = position74, tokenIndex74, depth74
							}
							depth--
							add(rulePegText, position72)
						}
						if buffer[position] != rune('m') {
							goto l69
						}
						position++
						if buffer[position] != rune('s') {
							goto l69
						}
						position++
						{
							add(ruleAction4, position)
						}
						depth--
						add(ruleDuration, position71)
					}
					goto l70
				l69:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
				}
			l70:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					{
						position78 := position
						depth++
						{
							position79 := position
							depth++
							if !matchDot() {
								goto l76
							}
						l80:
							{
								position81, tokenIndex81, depth81 := position, tokenIndex, depth
								if !matchDot() {
									goto l81
								}
								goto l80
							l81:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
							}
							depth--
							add(rulePegText, position79)
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleextra, position78)
					}
					goto l77
				l76:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
				}
			l77:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if !matchDot() {
						goto l83
					}
					goto l0
				l83:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
				}
				depth--
				add(ruleMongoLogLine, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Timestamp <- <(timestamp24 / timestamp26)> */
		nil,
		/* 2 Thread <- <('[' <letterOrDigit+> ']' Action0)> */
		nil,
		/* 3 Op <- <(<((&('c') ('c' 'o' 'm' 'm' 'a' 'n' 'd')) | (&('g') ('g' 'e' 't' 'm' 'o' 'r' 'e')) | (&('r') ('r' 'e' 'm' 'o' 'v' 'e')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('i') ('i' 'n' 's' 'e' 'r' 't')) | (&('q') ('q' 'u' 'e' 'r' 'y')))> Action1)> */
		nil,
		/* 4 LineField <- <((commandFieldName / plainFieldName) S? LineValue S? Action2)> */
		func() bool {
			position87, tokenIndex87, depth87 := position, tokenIndex, depth
			{
				position88 := position
				depth++
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					{
						position91 := position
						depth++
						if buffer[position] != rune('c') {
							goto l90
						}
						position++
						if buffer[position] != rune('o') {
							goto l90
						}
						position++
						if buffer[position] != rune('m') {
							goto l90
						}
						position++
						if buffer[position] != rune('m') {
							goto l90
						}
						position++
						if buffer[position] != rune('a') {
							goto l90
						}
						position++
						if buffer[position] != rune('n') {
							goto l90
						}
						position++
						if buffer[position] != rune('d') {
							goto l90
						}
						position++
						if buffer[position] != rune(':') {
							goto l90
						}
						position++
						if buffer[position] != rune(' ') {
							goto l90
						}
						position++
						{
							position92 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l90
							}
						l93:
							{
								position94, tokenIndex94, depth94 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l94
								}
								goto l93
							l94:
								position, tokenIndex, depth = position94, tokenIndex94, depth94
							}
							depth--
							add(rulePegText, position92)
						}
						{
							add(ruleAction6, position)
						}
						depth--
						add(rulecommandFieldName, position91)
					}
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					{
						position96 := position
						depth++
						{
							position97 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l87
							}
						l98:
							{
								position99, tokenIndex99, depth99 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l99
								}
								goto l98
							l99:
								position, tokenIndex, depth = position99, tokenIndex99, depth99
							}
							depth--
							add(rulePegText, position97)
						}
						if buffer[position] != rune(':') {
							goto l87
						}
						position++
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleplainFieldName, position96)
					}
				}
			l89:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l101
					}
					goto l102
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
				}
			l102:
				{
					position103 := position
					depth++
					{
						position104, tokenIndex104, depth104 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l104
						}
						goto l105
					l104:
						position, tokenIndex, depth = position104, tokenIndex104, depth104
					}
				l105:
					{
						switch buffer[position] {
						case '"':
							if !_rules[ruleString]() {
								goto l87
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l87
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l87
							}
							break
						}
					}

					depth--
					add(ruleLineValue, position103)
				}
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l107
					}
					goto l108
				l107:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
				}
			l108:
				{
					add(ruleAction2, position)
				}
				depth--
				add(ruleLineField, position88)
			}
			return true
		l87:
			position, tokenIndex, depth = position87, tokenIndex87, depth87
			return false
		},
		/* 5 NS <- <(<nsChar+> Action3)> */
		nil,
		/* 6 Locks <- <('l' 'o' 'c' 'k' 's' '(' 'm' 'i' 'c' 'r' 'o' 's' ')' S? lock*)> */
		nil,
		/* 7 lock <- <(((&('R') 'R') | (&('r') 'r') | (&('W' | 'w') ('w' / 'W'))) ':' [0-9]+ S?)> */
		nil,
		/* 8 Duration <- <(<[0-9]+> ('m' 's') Action4)> */
		nil,
		/* 9 plainFieldName <- <(<fieldChar+> ':' Action5)> */
		nil,
		/* 10 commandFieldName <- <('c' 'o' 'm' 'm' 'a' 'n' 'd' ':' ' ' <fieldChar+> Action6)> */
		nil,
		/* 11 LineValue <- <(S? ((&('"') String) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		nil,
		/* 12 timestamp24 <- <(<(date ' ' time)> Action7)> */
		nil,
		/* 13 timestamp26 <- <(<datetime26> Action8)> */
		nil,
		/* 14 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 15 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 16 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position121, tokenIndex121, depth121 := position, tokenIndex, depth
			{
				position122 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l121
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l121
				}
				position++
				depth--
				add(ruledigit2, position122)
			}
			return true
		l121:
			position, tokenIndex, depth = position121, tokenIndex121, depth121
			return false
		},
		/* 17 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 18 tz <- <('+' [0-9]+)> */
		nil,
		/* 19 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position125, tokenIndex125, depth125 := position, tokenIndex, depth
			{
				position126 := position
				depth++
				{
					position127 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l125
					}
					depth--
					add(rulehour, position127)
				}
				if buffer[position] != rune(':') {
					goto l125
				}
				position++
				{
					position128 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l125
					}
					depth--
					add(ruleminute, position128)
				}
				if buffer[position] != rune(':') {
					goto l125
				}
				position++
				{
					position129 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l125
					}
					depth--
					add(rulesecond, position129)
				}
				if buffer[position] != rune('.') {
					goto l125
				}
				position++
				{
					position130 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l125
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l125
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l125
					}
					position++
					depth--
					add(rulemillisecond, position130)
				}
				depth--
				add(ruletime, position126)
			}
			return true
		l125:
			position, tokenIndex, depth = position125, tokenIndex125, depth125
			return false
		},
		/* 20 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 21 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 22 dayNum <- <digit2?> */
		nil,
		/* 23 hour <- <digit2> */
		nil,
		/* 24 minute <- <digit2> */
		nil,
		/* 25 second <- <digit2> */
		nil,
		/* 26 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 27 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 28 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 29 extra <- <(<.+> Action9)> */
		nil,
		/* 30 S <- <' '+> */
		func() bool {
			position141, tokenIndex141, depth141 := position, tokenIndex, depth
			{
				position142 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l141
				}
				position++
			l143:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
				}
				depth--
				add(ruleS, position142)
			}
			return true
		l141:
			position, tokenIndex, depth = position141, tokenIndex141, depth141
			return false
		},
		/* 31 Doc <- <('{' Action10 DocElements? '}' Action11)> */
		func() bool {
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				if buffer[position] != rune('{') {
					goto l145
				}
				position++
				{
					add(ruleAction10, position)
				}
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					{
						position150 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l148
						}
					l151:
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l152
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l152
							}
							goto l151
						l152:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
						}
						depth--
						add(ruleDocElements, position150)
					}
					goto l149
				l148:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
				}
			l149:
				if buffer[position] != rune('}') {
					goto l145
				}
				position++
				{
					add(ruleAction11, position)
				}
				depth--
				add(ruleDoc, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
			return false
		},
		/* 32 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 33 DocElem <- <(S? Field S? Value S? Action12)> */
		func() bool {
			position155, tokenIndex155, depth155 := position, tokenIndex, depth
			{
				position156 := position
				depth++
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l157
					}
					goto l158
				l157:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
				}
			l158:
				{
					position159 := position
					depth++
					{
						position160 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l155
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
						goto l155
					}
					position++
					{
						add(ruleAction16, position)
					}
					depth--
					add(ruleField, position159)
				}
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l164
					}
					goto l165
				l164:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
				}
			l165:
				if !_rules[ruleValue]() {
					goto l155
				}
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
				{
					add(ruleAction12, position)
				}
				depth--
				add(ruleDocElem, position156)
			}
			return true
		l155:
			position, tokenIndex, depth = position155, tokenIndex155, depth155
			return false
		},
		/* 34 List <- <('[' Action13 ListElements? ']' Action14)> */
		nil,
		/* 35 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 36 ListElem <- <(S? Value S? Action15)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
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
				if !_rules[ruleValue]() {
					goto l171
				}
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l175
					}
					goto l176
				l175:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
				}
			l176:
				{
					add(ruleAction15, position)
				}
				depth--
				add(ruleListElem, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 37 Field <- <(<fieldChar+> ':' Action16)> */
		nil,
		/* 38 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					{
						position183 := position
						depth++
						if buffer[position] != rune('n') {
							goto l182
						}
						position++
						if buffer[position] != rune('u') {
							goto l182
						}
						position++
						if buffer[position] != rune('l') {
							goto l182
						}
						position++
						if buffer[position] != rune('l') {
							goto l182
						}
						position++
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleNull, position183)
					}
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					{
						position186 := position
						depth++
						if buffer[position] != rune('M') {
							goto l185
						}
						position++
						if buffer[position] != rune('i') {
							goto l185
						}
						position++
						if buffer[position] != rune('n') {
							goto l185
						}
						position++
						if buffer[position] != rune('K') {
							goto l185
						}
						position++
						if buffer[position] != rune('e') {
							goto l185
						}
						position++
						if buffer[position] != rune('y') {
							goto l185
						}
						position++
						{
							add(ruleAction28, position)
						}
						depth--
						add(ruleMinKey, position186)
					}
					goto l181
				l185:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					{
						switch buffer[position] {
						case 'M':
							{
								position189 := position
								depth++
								if buffer[position] != rune('M') {
									goto l179
								}
								position++
								if buffer[position] != rune('a') {
									goto l179
								}
								position++
								if buffer[position] != rune('x') {
									goto l179
								}
								position++
								if buffer[position] != rune('K') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('y') {
									goto l179
								}
								position++
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleMaxKey, position189)
							}
							break
						case 'u':
							{
								position191 := position
								depth++
								if buffer[position] != rune('u') {
									goto l179
								}
								position++
								if buffer[position] != rune('n') {
									goto l179
								}
								position++
								if buffer[position] != rune('d') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('f') {
									goto l179
								}
								position++
								if buffer[position] != rune('i') {
									goto l179
								}
								position++
								if buffer[position] != rune('n') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('d') {
									goto l179
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleUndefined, position191)
							}
							break
						case 'N':
							{
								position193 := position
								depth++
								if buffer[position] != rune('N') {
									goto l179
								}
								position++
								if buffer[position] != rune('u') {
									goto l179
								}
								position++
								if buffer[position] != rune('m') {
									goto l179
								}
								position++
								if buffer[position] != rune('b') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('r') {
									goto l179
								}
								position++
								if buffer[position] != rune('L') {
									goto l179
								}
								position++
								if buffer[position] != rune('o') {
									goto l179
								}
								position++
								if buffer[position] != rune('n') {
									goto l179
								}
								position++
								if buffer[position] != rune('g') {
									goto l179
								}
								position++
								if buffer[position] != rune('(') {
									goto l179
								}
								position++
								{
									position194 := position
									depth++
									{
										position197, tokenIndex197, depth197 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l197
										}
										position++
										goto l179
									l197:
										position, tokenIndex, depth = position197, tokenIndex197, depth197
									}
									if !matchDot() {
										goto l179
									}
								l195:
									{
										position196, tokenIndex196, depth196 := position, tokenIndex, depth
										{
											position198, tokenIndex198, depth198 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l198
											}
											position++
											goto l196
										l198:
											position, tokenIndex, depth = position198, tokenIndex198, depth198
										}
										if !matchDot() {
											goto l196
										}
										goto l195
									l196:
										position, tokenIndex, depth = position196, tokenIndex196, depth196
									}
									depth--
									add(rulePegText, position194)
								}
								if buffer[position] != rune(')') {
									goto l179
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleNumberLong, position193)
							}
							break
						case '/':
							{
								position200 := position
								depth++
								if buffer[position] != rune('/') {
									goto l179
								}
								position++
								{
									position201 := position
									depth++
									{
										position202 := position
										depth++
										{
											position205 := position
											depth++
											{
												position206, tokenIndex206, depth206 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l206
												}
												position++
												goto l179
											l206:
												position, tokenIndex, depth = position206, tokenIndex206, depth206
											}
											if !matchDot() {
												goto l179
											}
											depth--
											add(ruleregexChar, position205)
										}
									l203:
										{
											position204, tokenIndex204, depth204 := position, tokenIndex, depth
											{
												position207 := position
												depth++
												{
													position208, tokenIndex208, depth208 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l208
													}
													position++
													goto l204
												l208:
													position, tokenIndex, depth = position208, tokenIndex208, depth208
												}
												if !matchDot() {
													goto l204
												}
												depth--
												add(ruleregexChar, position207)
											}
											goto l203
										l204:
											position, tokenIndex, depth = position204, tokenIndex204, depth204
										}
										if buffer[position] != rune('/') {
											goto l179
										}
										position++
									l209:
										{
											position210, tokenIndex210, depth210 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l210
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l210
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l210
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l210
													}
													position++
													break
												}
											}

											goto l209
										l210:
											position, tokenIndex, depth = position210, tokenIndex210, depth210
										}
										depth--
										add(ruleregexBody, position202)
									}
									depth--
									add(rulePegText, position201)
								}
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleRegex, position200)
							}
							break
						case 'T':
							{
								position213 := position
								depth++
								if buffer[position] != rune('T') {
									goto l179
								}
								position++
								if buffer[position] != rune('i') {
									goto l179
								}
								position++
								if buffer[position] != rune('m') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('s') {
									goto l179
								}
								position++
								if buffer[position] != rune('t') {
									goto l179
								}
								position++
								if buffer[position] != rune('a') {
									goto l179
								}
								position++
								if buffer[position] != rune('m') {
									goto l179
								}
								position++
								if buffer[position] != rune('p') {
									goto l179
								}
								position++
								if buffer[position] != rune('(') {
									goto l179
								}
								position++
								{
									position214 := position
									depth++
									{
										position217, tokenIndex217, depth217 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l217
										}
										position++
										goto l179
									l217:
										position, tokenIndex, depth = position217, tokenIndex217, depth217
									}
									if !matchDot() {
										goto l179
									}
								l215:
									{
										position216, tokenIndex216, depth216 := position, tokenIndex, depth
										{
											position218, tokenIndex218, depth218 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l218
											}
											position++
											goto l216
										l218:
											position, tokenIndex, depth = position218, tokenIndex218, depth218
										}
										if !matchDot() {
											goto l216
										}
										goto l215
									l216:
										position, tokenIndex, depth = position216, tokenIndex216, depth216
									}
									depth--
									add(rulePegText, position214)
								}
								if buffer[position] != rune(')') {
									goto l179
								}
								position++
								{
									add(ruleAction26, position)
								}
								depth--
								add(ruleTimestampVal, position213)
							}
							break
						case 'B':
							{
								position220 := position
								depth++
								if buffer[position] != rune('B') {
									goto l179
								}
								position++
								if buffer[position] != rune('i') {
									goto l179
								}
								position++
								if buffer[position] != rune('n') {
									goto l179
								}
								position++
								if buffer[position] != rune('D') {
									goto l179
								}
								position++
								if buffer[position] != rune('a') {
									goto l179
								}
								position++
								if buffer[position] != rune('t') {
									goto l179
								}
								position++
								if buffer[position] != rune('a') {
									goto l179
								}
								position++
								if buffer[position] != rune('(') {
									goto l179
								}
								position++
								{
									position221 := position
									depth++
									{
										position224, tokenIndex224, depth224 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l224
										}
										position++
										goto l179
									l224:
										position, tokenIndex, depth = position224, tokenIndex224, depth224
									}
									if !matchDot() {
										goto l179
									}
								l222:
									{
										position223, tokenIndex223, depth223 := position, tokenIndex, depth
										{
											position225, tokenIndex225, depth225 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
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
									depth--
									add(rulePegText, position221)
								}
								if buffer[position] != rune(')') {
									goto l179
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleBinData, position220)
							}
							break
						case 'D', 'n':
							{
								position227 := position
								depth++
								{
									position228, tokenIndex228, depth228 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l228
									}
									position++
									if buffer[position] != rune('e') {
										goto l228
									}
									position++
									if buffer[position] != rune('w') {
										goto l228
									}
									position++
									if buffer[position] != rune(' ') {
										goto l228
									}
									position++
									goto l229
								l228:
									position, tokenIndex, depth = position228, tokenIndex228, depth228
								}
							l229:
								if buffer[position] != rune('D') {
									goto l179
								}
								position++
								if buffer[position] != rune('a') {
									goto l179
								}
								position++
								if buffer[position] != rune('t') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('(') {
									goto l179
								}
								position++
								{
									position230 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l179
									}
									position++
								l231:
									{
										position232, tokenIndex232, depth232 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l232
										}
										position++
										goto l231
									l232:
										position, tokenIndex, depth = position232, tokenIndex232, depth232
									}
									depth--
									add(rulePegText, position230)
								}
								if buffer[position] != rune(')') {
									goto l179
								}
								position++
								{
									add(ruleAction22, position)
								}
								depth--
								add(ruleDate, position227)
							}
							break
						case 'O':
							{
								position234 := position
								depth++
								if buffer[position] != rune('O') {
									goto l179
								}
								position++
								if buffer[position] != rune('b') {
									goto l179
								}
								position++
								if buffer[position] != rune('j') {
									goto l179
								}
								position++
								if buffer[position] != rune('e') {
									goto l179
								}
								position++
								if buffer[position] != rune('c') {
									goto l179
								}
								position++
								if buffer[position] != rune('t') {
									goto l179
								}
								position++
								if buffer[position] != rune('I') {
									goto l179
								}
								position++
								if buffer[position] != rune('d') {
									goto l179
								}
								position++
								if buffer[position] != rune('(') {
									goto l179
								}
								position++
								if buffer[position] != rune('"') {
									goto l179
								}
								position++
								{
									position235 := position
									depth++
								l236:
									{
										position237, tokenIndex237, depth237 := position, tokenIndex, depth
										{
											position238 := position
											depth++
											{
												position239, tokenIndex239, depth239 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l240
												}
												position++
												goto l239
											l240:
												position, tokenIndex, depth = position239, tokenIndex239, depth239
												{
													position241, tokenIndex241, depth241 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l242
													}
													position++
													goto l241
												l242:
													position, tokenIndex, depth = position241, tokenIndex241, depth241
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l237
													}
													position++
												}
											l241:
											}
										l239:
											depth--
											add(rulehexChar, position238)
										}
										goto l236
									l237:
										position, tokenIndex, depth = position237, tokenIndex237, depth237
									}
									depth--
									add(rulePegText, position235)
								}
								if buffer[position] != rune('"') {
									goto l179
								}
								position++
								if buffer[position] != rune(')') {
									goto l179
								}
								position++
								{
									add(ruleAction23, position)
								}
								depth--
								add(ruleObjectID, position234)
							}
							break
						case '"':
							if !_rules[ruleString]() {
								goto l179
							}
							break
						case 'f', 't':
							{
								position244 := position
								depth++
								{
									position245, tokenIndex245, depth245 := position, tokenIndex, depth
									{
										position247 := position
										depth++
										if buffer[position] != rune('t') {
											goto l246
										}
										position++
										if buffer[position] != rune('r') {
											goto l246
										}
										position++
										if buffer[position] != rune('u') {
											goto l246
										}
										position++
										if buffer[position] != rune('e') {
											goto l246
										}
										position++
										{
											add(ruleAction20, position)
										}
										depth--
										add(ruleTrue, position247)
									}
									goto l245
								l246:
									position, tokenIndex, depth = position245, tokenIndex245, depth245
									{
										position249 := position
										depth++
										if buffer[position] != rune('f') {
											goto l179
										}
										position++
										if buffer[position] != rune('a') {
											goto l179
										}
										position++
										if buffer[position] != rune('l') {
											goto l179
										}
										position++
										if buffer[position] != rune('s') {
											goto l179
										}
										position++
										if buffer[position] != rune('e') {
											goto l179
										}
										position++
										{
											add(ruleAction21, position)
										}
										depth--
										add(ruleFalse, position249)
									}
								}
							l245:
								depth--
								add(ruleBoolean, position244)
							}
							break
						case '[':
							{
								position251 := position
								depth++
								if buffer[position] != rune('[') {
									goto l179
								}
								position++
								{
									add(ruleAction13, position)
								}
								{
									position253, tokenIndex253, depth253 := position, tokenIndex, depth
									{
										position255 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l253
										}
									l256:
										{
											position257, tokenIndex257, depth257 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l257
											}
											position++
											if !_rules[ruleListElem]() {
												goto l257
											}
											goto l256
										l257:
											position, tokenIndex, depth = position257, tokenIndex257, depth257
										}
										depth--
										add(ruleListElements, position255)
									}
									goto l254
								l253:
									position, tokenIndex, depth = position253, tokenIndex253, depth253
								}
							l254:
								if buffer[position] != rune(']') {
									goto l179
								}
								position++
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleList, position251)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l179
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l179
							}
							break
						}
					}

				}
			l181:
				depth--
				add(ruleValue, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 39 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action17)> */
		func() bool {
			position259, tokenIndex259, depth259 := position, tokenIndex, depth
			{
				position260 := position
				depth++
				{
					position261 := position
					depth++
					{
						position262, tokenIndex262, depth262 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l262
						}
						position++
						goto l263
					l262:
						position, tokenIndex, depth = position262, tokenIndex262, depth262
					}
				l263:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l259
					}
					position++
				l264:
					{
						position265, tokenIndex265, depth265 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l265
						}
						position++
						goto l264
					l265:
						position, tokenIndex, depth = position265, tokenIndex265, depth265
					}
					{
						position266, tokenIndex266, depth266 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l266
						}
						position++
						goto l267
					l266:
						position, tokenIndex, depth = position266, tokenIndex266, depth266
					}
				l267:
				l268:
					{
						position269, tokenIndex269, depth269 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l269
						}
						position++
						goto l268
					l269:
						position, tokenIndex, depth = position269, tokenIndex269, depth269
					}
					depth--
					add(rulePegText, position261)
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleNumeric, position260)
			}
			return true
		l259:
			position, tokenIndex, depth = position259, tokenIndex259, depth259
			return false
		},
		/* 40 Boolean <- <(True / False)> */
		nil,
		/* 41 String <- <('"' <stringChar*> '"' Action18)> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				if buffer[position] != rune('"') {
					goto l272
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
								{
									position280, tokenIndex280, depth280 := position, tokenIndex, depth
									{
										position281, tokenIndex281, depth281 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l282
										}
										position++
										goto l281
									l282:
										position, tokenIndex, depth = position281, tokenIndex281, depth281
										if buffer[position] != rune('\\') {
											goto l280
										}
										position++
									}
								l281:
									goto l279
								l280:
									position, tokenIndex, depth = position280, tokenIndex280, depth280
								}
								if !matchDot() {
									goto l279
								}
								goto l278
							l279:
								position, tokenIndex, depth = position278, tokenIndex278, depth278
								if buffer[position] != rune('\\') {
									goto l276
								}
								position++
								{
									position283, tokenIndex283, depth283 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l284
									}
									position++
									goto l283
								l284:
									position, tokenIndex, depth = position283, tokenIndex283, depth283
									if buffer[position] != rune('\\') {
										goto l276
									}
									position++
								}
							l283:
							}
						l278:
							depth--
							add(rulestringChar, position277)
						}
						goto l275
					l276:
						position, tokenIndex, depth = position276, tokenIndex276, depth276
					}
					depth--
					add(rulePegText, position274)
				}
				if buffer[position] != rune('"') {
					goto l272
				}
				position++
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleString, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 42 Null <- <('n' 'u' 'l' 'l' Action19)> */
		nil,
		/* 43 True <- <('t' 'r' 'u' 'e' Action20)> */
		nil,
		/* 44 False <- <('f' 'a' 'l' 's' 'e' Action21)> */
		nil,
		/* 45 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action22)> */
		nil,
		/* 46 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action23)> */
		nil,
		/* 47 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action24)> */
		nil,
		/* 48 Regex <- <('/' <regexBody> Action25)> */
		nil,
		/* 49 TimestampVal <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action26)> */
		nil,
		/* 50 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action27)> */
		nil,
		/* 51 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action28)> */
		nil,
		/* 52 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action29)> */
		nil,
		/* 53 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action30)> */
		nil,
		/* 54 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 55 regexChar <- <(!'/' .)> */
		nil,
		/* 56 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 57 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 58 fieldChar <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position302, tokenIndex302, depth302 := position, tokenIndex, depth
			{
				position303 := position
				depth++
				{
					switch buffer[position] {
					case '$', '_':
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('$') {
								goto l302
							}
							position++
						}
					l305:
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l302
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l302
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l302
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position303)
			}
			return true
		l302:
			position, tokenIndex, depth = position302, tokenIndex302, depth302
			return false
		},
		nil,
		/* 61 Action0 <- <{ p.SetField("thread", buffer[begin:end]) }> */
		nil,
		/* 62 Action1 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 63 Action2 <- <{ p.EndField() }> */
		nil,
		/* 64 Action3 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 65 Action4 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 66 Action5 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 67 Action6 <- <{ p.SetField("commandType", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 68 Action7 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 69 Action8 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 70 Action9 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 71 Action10 <- <{ p.PushMap() }> */
		nil,
		/* 72 Action11 <- <{ p.PopMap() }> */
		nil,
		/* 73 Action12 <- <{ p.SetMapValue() }> */
		nil,
		/* 74 Action13 <- <{ p.PushList() }> */
		nil,
		/* 75 Action14 <- <{ p.PopList() }> */
		nil,
		/* 76 Action15 <- <{ p.SetListValue() }> */
		nil,
		/* 77 Action16 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 78 Action17 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 79 Action18 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 80 Action19 <- <{ p.PushValue(nil) }> */
		nil,
		/* 81 Action20 <- <{ p.PushValue(true) }> */
		nil,
		/* 82 Action21 <- <{ p.PushValue(false) }> */
		nil,
		/* 83 Action22 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 84 Action23 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 85 Action24 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 86 Action25 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 87 Action26 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 88 Action27 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 89 Action28 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 90 Action29 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 91 Action30 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
