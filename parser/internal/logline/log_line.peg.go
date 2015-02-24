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
						add(ruleAction3, position)
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
							add(ruleAction4, position)
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
							add(ruleAction9, position)
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
		/* 4 LineField <- <((commandFieldName / plainFieldName) S? LineValue S? Action2)> */
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
						depth--
						add(rulecommandFieldName, position93)
					}
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					{
						position98 := position
						depth++
						{
							position99 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l89
							}
						l100:
							{
								position101, tokenIndex101, depth101 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l101
								}
								goto l100
							l101:
								position, tokenIndex, depth = position101, tokenIndex101, depth101
							}
							depth--
							add(rulePegText, position99)
						}
						if buffer[position] != rune(':') {
							goto l89
						}
						position++
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleplainFieldName, position98)
					}
				}
			l91:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l103
					}
					goto l104
				l103:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
				}
			l104:
				{
					position105 := position
					depth++
					{
						position106, tokenIndex106, depth106 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l106
						}
						goto l107
					l106:
						position, tokenIndex, depth = position106, tokenIndex106, depth106
					}
				l107:
					{
						switch buffer[position] {
						case '"':
							if !_rules[ruleString]() {
								goto l89
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l89
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l89
							}
							break
						}
					}

					depth--
					add(ruleLineValue, position105)
				}
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l109
					}
					goto l110
				l109:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
				}
			l110:
				{
					add(ruleAction2, position)
				}
				depth--
				add(ruleLineField, position90)
			}
			return true
		l89:
			position, tokenIndex, depth = position89, tokenIndex89, depth89
			return false
		},
		/* 5 NS <- <(<nsChar+> ' ' Action3)> */
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
			position123, tokenIndex123, depth123 := position, tokenIndex, depth
			{
				position124 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l123
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l123
				}
				position++
				depth--
				add(ruledigit2, position124)
			}
			return true
		l123:
			position, tokenIndex, depth = position123, tokenIndex123, depth123
			return false
		},
		/* 17 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 18 tz <- <('+' [0-9]+)> */
		nil,
		/* 19 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				{
					position129 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l127
					}
					depth--
					add(rulehour, position129)
				}
				if buffer[position] != rune(':') {
					goto l127
				}
				position++
				{
					position130 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l127
					}
					depth--
					add(ruleminute, position130)
				}
				if buffer[position] != rune(':') {
					goto l127
				}
				position++
				{
					position131 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l127
					}
					depth--
					add(rulesecond, position131)
				}
				if buffer[position] != rune('.') {
					goto l127
				}
				position++
				{
					position132 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l127
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l127
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l127
					}
					position++
					depth--
					add(rulemillisecond, position132)
				}
				depth--
				add(ruletime, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
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
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l143
				}
				position++
			l145:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
				}
				depth--
				add(ruleS, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 31 Doc <- <('{' Action10 DocElements? '}' Action11)> */
		func() bool {
			position147, tokenIndex147, depth147 := position, tokenIndex, depth
			{
				position148 := position
				depth++
				if buffer[position] != rune('{') {
					goto l147
				}
				position++
				{
					add(ruleAction10, position)
				}
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					{
						position152 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l150
						}
					l153:
						{
							position154, tokenIndex154, depth154 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l154
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l154
							}
							goto l153
						l154:
							position, tokenIndex, depth = position154, tokenIndex154, depth154
						}
						depth--
						add(ruleDocElements, position152)
					}
					goto l151
				l150:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
				}
			l151:
				if buffer[position] != rune('}') {
					goto l147
				}
				position++
				{
					add(ruleAction11, position)
				}
				depth--
				add(ruleDoc, position148)
			}
			return true
		l147:
			position, tokenIndex, depth = position147, tokenIndex147, depth147
			return false
		},
		/* 32 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 33 DocElem <- <(S? Field S? Value S? Action12)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l159
					}
					goto l160
				l159:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
				}
			l160:
				{
					position161 := position
					depth++
					{
						position162 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l157
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
						goto l157
					}
					position++
					{
						add(ruleAction16, position)
					}
					depth--
					add(ruleField, position161)
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
				if !_rules[ruleValue]() {
					goto l157
				}
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l168
					}
					goto l169
				l168:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
				}
			l169:
				{
					add(ruleAction12, position)
				}
				depth--
				add(ruleDocElem, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 34 List <- <('[' Action13 ListElements? ']' Action14)> */
		nil,
		/* 35 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 36 ListElem <- <(S? Value S? Action15)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
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
				if !_rules[ruleValue]() {
					goto l173
				}
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l177
					}
					goto l178
				l177:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
				}
			l178:
				{
					add(ruleAction15, position)
				}
				depth--
				add(ruleListElem, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 37 Field <- <(<fieldChar+> ':' Action16)> */
		nil,
		/* 38 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position181, tokenIndex181, depth181 := position, tokenIndex, depth
			{
				position182 := position
				depth++
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					{
						position185 := position
						depth++
						if buffer[position] != rune('n') {
							goto l184
						}
						position++
						if buffer[position] != rune('u') {
							goto l184
						}
						position++
						if buffer[position] != rune('l') {
							goto l184
						}
						position++
						if buffer[position] != rune('l') {
							goto l184
						}
						position++
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleNull, position185)
					}
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					{
						position188 := position
						depth++
						if buffer[position] != rune('M') {
							goto l187
						}
						position++
						if buffer[position] != rune('i') {
							goto l187
						}
						position++
						if buffer[position] != rune('n') {
							goto l187
						}
						position++
						if buffer[position] != rune('K') {
							goto l187
						}
						position++
						if buffer[position] != rune('e') {
							goto l187
						}
						position++
						if buffer[position] != rune('y') {
							goto l187
						}
						position++
						{
							add(ruleAction28, position)
						}
						depth--
						add(ruleMinKey, position188)
					}
					goto l183
				l187:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					{
						switch buffer[position] {
						case 'M':
							{
								position191 := position
								depth++
								if buffer[position] != rune('M') {
									goto l181
								}
								position++
								if buffer[position] != rune('a') {
									goto l181
								}
								position++
								if buffer[position] != rune('x') {
									goto l181
								}
								position++
								if buffer[position] != rune('K') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('y') {
									goto l181
								}
								position++
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleMaxKey, position191)
							}
							break
						case 'u':
							{
								position193 := position
								depth++
								if buffer[position] != rune('u') {
									goto l181
								}
								position++
								if buffer[position] != rune('n') {
									goto l181
								}
								position++
								if buffer[position] != rune('d') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('f') {
									goto l181
								}
								position++
								if buffer[position] != rune('i') {
									goto l181
								}
								position++
								if buffer[position] != rune('n') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('d') {
									goto l181
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleUndefined, position193)
							}
							break
						case 'N':
							{
								position195 := position
								depth++
								if buffer[position] != rune('N') {
									goto l181
								}
								position++
								if buffer[position] != rune('u') {
									goto l181
								}
								position++
								if buffer[position] != rune('m') {
									goto l181
								}
								position++
								if buffer[position] != rune('b') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('r') {
									goto l181
								}
								position++
								if buffer[position] != rune('L') {
									goto l181
								}
								position++
								if buffer[position] != rune('o') {
									goto l181
								}
								position++
								if buffer[position] != rune('n') {
									goto l181
								}
								position++
								if buffer[position] != rune('g') {
									goto l181
								}
								position++
								if buffer[position] != rune('(') {
									goto l181
								}
								position++
								{
									position196 := position
									depth++
									{
										position199, tokenIndex199, depth199 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l199
										}
										position++
										goto l181
									l199:
										position, tokenIndex, depth = position199, tokenIndex199, depth199
									}
									if !matchDot() {
										goto l181
									}
								l197:
									{
										position198, tokenIndex198, depth198 := position, tokenIndex, depth
										{
											position200, tokenIndex200, depth200 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l200
											}
											position++
											goto l198
										l200:
											position, tokenIndex, depth = position200, tokenIndex200, depth200
										}
										if !matchDot() {
											goto l198
										}
										goto l197
									l198:
										position, tokenIndex, depth = position198, tokenIndex198, depth198
									}
									depth--
									add(rulePegText, position196)
								}
								if buffer[position] != rune(')') {
									goto l181
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleNumberLong, position195)
							}
							break
						case '/':
							{
								position202 := position
								depth++
								if buffer[position] != rune('/') {
									goto l181
								}
								position++
								{
									position203 := position
									depth++
									{
										position204 := position
										depth++
										{
											position207 := position
											depth++
											{
												position208, tokenIndex208, depth208 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l208
												}
												position++
												goto l181
											l208:
												position, tokenIndex, depth = position208, tokenIndex208, depth208
											}
											if !matchDot() {
												goto l181
											}
											depth--
											add(ruleregexChar, position207)
										}
									l205:
										{
											position206, tokenIndex206, depth206 := position, tokenIndex, depth
											{
												position209 := position
												depth++
												{
													position210, tokenIndex210, depth210 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l210
													}
													position++
													goto l206
												l210:
													position, tokenIndex, depth = position210, tokenIndex210, depth210
												}
												if !matchDot() {
													goto l206
												}
												depth--
												add(ruleregexChar, position209)
											}
											goto l205
										l206:
											position, tokenIndex, depth = position206, tokenIndex206, depth206
										}
										if buffer[position] != rune('/') {
											goto l181
										}
										position++
									l211:
										{
											position212, tokenIndex212, depth212 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l212
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l212
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l212
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l212
													}
													position++
													break
												}
											}

											goto l211
										l212:
											position, tokenIndex, depth = position212, tokenIndex212, depth212
										}
										depth--
										add(ruleregexBody, position204)
									}
									depth--
									add(rulePegText, position203)
								}
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleRegex, position202)
							}
							break
						case 'T':
							{
								position215 := position
								depth++
								if buffer[position] != rune('T') {
									goto l181
								}
								position++
								if buffer[position] != rune('i') {
									goto l181
								}
								position++
								if buffer[position] != rune('m') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('s') {
									goto l181
								}
								position++
								if buffer[position] != rune('t') {
									goto l181
								}
								position++
								if buffer[position] != rune('a') {
									goto l181
								}
								position++
								if buffer[position] != rune('m') {
									goto l181
								}
								position++
								if buffer[position] != rune('p') {
									goto l181
								}
								position++
								if buffer[position] != rune('(') {
									goto l181
								}
								position++
								{
									position216 := position
									depth++
									{
										position219, tokenIndex219, depth219 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l219
										}
										position++
										goto l181
									l219:
										position, tokenIndex, depth = position219, tokenIndex219, depth219
									}
									if !matchDot() {
										goto l181
									}
								l217:
									{
										position218, tokenIndex218, depth218 := position, tokenIndex, depth
										{
											position220, tokenIndex220, depth220 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l220
											}
											position++
											goto l218
										l220:
											position, tokenIndex, depth = position220, tokenIndex220, depth220
										}
										if !matchDot() {
											goto l218
										}
										goto l217
									l218:
										position, tokenIndex, depth = position218, tokenIndex218, depth218
									}
									depth--
									add(rulePegText, position216)
								}
								if buffer[position] != rune(')') {
									goto l181
								}
								position++
								{
									add(ruleAction26, position)
								}
								depth--
								add(ruleTimestampVal, position215)
							}
							break
						case 'B':
							{
								position222 := position
								depth++
								if buffer[position] != rune('B') {
									goto l181
								}
								position++
								if buffer[position] != rune('i') {
									goto l181
								}
								position++
								if buffer[position] != rune('n') {
									goto l181
								}
								position++
								if buffer[position] != rune('D') {
									goto l181
								}
								position++
								if buffer[position] != rune('a') {
									goto l181
								}
								position++
								if buffer[position] != rune('t') {
									goto l181
								}
								position++
								if buffer[position] != rune('a') {
									goto l181
								}
								position++
								if buffer[position] != rune('(') {
									goto l181
								}
								position++
								{
									position223 := position
									depth++
									{
										position226, tokenIndex226, depth226 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l226
										}
										position++
										goto l181
									l226:
										position, tokenIndex, depth = position226, tokenIndex226, depth226
									}
									if !matchDot() {
										goto l181
									}
								l224:
									{
										position225, tokenIndex225, depth225 := position, tokenIndex, depth
										{
											position227, tokenIndex227, depth227 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l227
											}
											position++
											goto l225
										l227:
											position, tokenIndex, depth = position227, tokenIndex227, depth227
										}
										if !matchDot() {
											goto l225
										}
										goto l224
									l225:
										position, tokenIndex, depth = position225, tokenIndex225, depth225
									}
									depth--
									add(rulePegText, position223)
								}
								if buffer[position] != rune(')') {
									goto l181
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleBinData, position222)
							}
							break
						case 'D', 'n':
							{
								position229 := position
								depth++
								{
									position230, tokenIndex230, depth230 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l230
									}
									position++
									if buffer[position] != rune('e') {
										goto l230
									}
									position++
									if buffer[position] != rune('w') {
										goto l230
									}
									position++
									if buffer[position] != rune(' ') {
										goto l230
									}
									position++
									goto l231
								l230:
									position, tokenIndex, depth = position230, tokenIndex230, depth230
								}
							l231:
								if buffer[position] != rune('D') {
									goto l181
								}
								position++
								if buffer[position] != rune('a') {
									goto l181
								}
								position++
								if buffer[position] != rune('t') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('(') {
									goto l181
								}
								position++
								{
									position232 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l181
									}
									position++
								l233:
									{
										position234, tokenIndex234, depth234 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l234
										}
										position++
										goto l233
									l234:
										position, tokenIndex, depth = position234, tokenIndex234, depth234
									}
									depth--
									add(rulePegText, position232)
								}
								if buffer[position] != rune(')') {
									goto l181
								}
								position++
								{
									add(ruleAction22, position)
								}
								depth--
								add(ruleDate, position229)
							}
							break
						case 'O':
							{
								position236 := position
								depth++
								if buffer[position] != rune('O') {
									goto l181
								}
								position++
								if buffer[position] != rune('b') {
									goto l181
								}
								position++
								if buffer[position] != rune('j') {
									goto l181
								}
								position++
								if buffer[position] != rune('e') {
									goto l181
								}
								position++
								if buffer[position] != rune('c') {
									goto l181
								}
								position++
								if buffer[position] != rune('t') {
									goto l181
								}
								position++
								if buffer[position] != rune('I') {
									goto l181
								}
								position++
								if buffer[position] != rune('d') {
									goto l181
								}
								position++
								if buffer[position] != rune('(') {
									goto l181
								}
								position++
								if buffer[position] != rune('"') {
									goto l181
								}
								position++
								{
									position237 := position
									depth++
								l238:
									{
										position239, tokenIndex239, depth239 := position, tokenIndex, depth
										{
											position240 := position
											depth++
											{
												position241, tokenIndex241, depth241 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l242
												}
												position++
												goto l241
											l242:
												position, tokenIndex, depth = position241, tokenIndex241, depth241
												{
													position243, tokenIndex243, depth243 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l244
													}
													position++
													goto l243
												l244:
													position, tokenIndex, depth = position243, tokenIndex243, depth243
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l239
													}
													position++
												}
											l243:
											}
										l241:
											depth--
											add(rulehexChar, position240)
										}
										goto l238
									l239:
										position, tokenIndex, depth = position239, tokenIndex239, depth239
									}
									depth--
									add(rulePegText, position237)
								}
								if buffer[position] != rune('"') {
									goto l181
								}
								position++
								if buffer[position] != rune(')') {
									goto l181
								}
								position++
								{
									add(ruleAction23, position)
								}
								depth--
								add(ruleObjectID, position236)
							}
							break
						case '"':
							if !_rules[ruleString]() {
								goto l181
							}
							break
						case 'f', 't':
							{
								position246 := position
								depth++
								{
									position247, tokenIndex247, depth247 := position, tokenIndex, depth
									{
										position249 := position
										depth++
										if buffer[position] != rune('t') {
											goto l248
										}
										position++
										if buffer[position] != rune('r') {
											goto l248
										}
										position++
										if buffer[position] != rune('u') {
											goto l248
										}
										position++
										if buffer[position] != rune('e') {
											goto l248
										}
										position++
										{
											add(ruleAction20, position)
										}
										depth--
										add(ruleTrue, position249)
									}
									goto l247
								l248:
									position, tokenIndex, depth = position247, tokenIndex247, depth247
									{
										position251 := position
										depth++
										if buffer[position] != rune('f') {
											goto l181
										}
										position++
										if buffer[position] != rune('a') {
											goto l181
										}
										position++
										if buffer[position] != rune('l') {
											goto l181
										}
										position++
										if buffer[position] != rune('s') {
											goto l181
										}
										position++
										if buffer[position] != rune('e') {
											goto l181
										}
										position++
										{
											add(ruleAction21, position)
										}
										depth--
										add(ruleFalse, position251)
									}
								}
							l247:
								depth--
								add(ruleBoolean, position246)
							}
							break
						case '[':
							{
								position253 := position
								depth++
								if buffer[position] != rune('[') {
									goto l181
								}
								position++
								{
									add(ruleAction13, position)
								}
								{
									position255, tokenIndex255, depth255 := position, tokenIndex, depth
									{
										position257 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l255
										}
									l258:
										{
											position259, tokenIndex259, depth259 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l259
											}
											position++
											if !_rules[ruleListElem]() {
												goto l259
											}
											goto l258
										l259:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
										}
										depth--
										add(ruleListElements, position257)
									}
									goto l256
								l255:
									position, tokenIndex, depth = position255, tokenIndex255, depth255
								}
							l256:
								if buffer[position] != rune(']') {
									goto l181
								}
								position++
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleList, position253)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l181
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l181
							}
							break
						}
					}

				}
			l183:
				depth--
				add(ruleValue, position182)
			}
			return true
		l181:
			position, tokenIndex, depth = position181, tokenIndex181, depth181
			return false
		},
		/* 39 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action17)> */
		func() bool {
			position261, tokenIndex261, depth261 := position, tokenIndex, depth
			{
				position262 := position
				depth++
				{
					position263 := position
					depth++
					{
						position264, tokenIndex264, depth264 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l264
						}
						position++
						goto l265
					l264:
						position, tokenIndex, depth = position264, tokenIndex264, depth264
					}
				l265:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l261
					}
					position++
				l266:
					{
						position267, tokenIndex267, depth267 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l267
						}
						position++
						goto l266
					l267:
						position, tokenIndex, depth = position267, tokenIndex267, depth267
					}
					{
						position268, tokenIndex268, depth268 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l268
						}
						position++
						goto l269
					l268:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
					}
				l269:
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
					add(rulePegText, position263)
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleNumeric, position262)
			}
			return true
		l261:
			position, tokenIndex, depth = position261, tokenIndex261, depth261
			return false
		},
		/* 40 Boolean <- <(True / False)> */
		nil,
		/* 41 String <- <('"' <stringChar*> '"' Action18)> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				if buffer[position] != rune('"') {
					goto l274
				}
				position++
				{
					position276 := position
					depth++
				l277:
					{
						position278, tokenIndex278, depth278 := position, tokenIndex, depth
						{
							position279 := position
							depth++
							{
								position280, tokenIndex280, depth280 := position, tokenIndex, depth
								{
									position282, tokenIndex282, depth282 := position, tokenIndex, depth
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
											goto l282
										}
										position++
									}
								l283:
									goto l281
								l282:
									position, tokenIndex, depth = position282, tokenIndex282, depth282
								}
								if !matchDot() {
									goto l281
								}
								goto l280
							l281:
								position, tokenIndex, depth = position280, tokenIndex280, depth280
								if buffer[position] != rune('\\') {
									goto l278
								}
								position++
								{
									position285, tokenIndex285, depth285 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l286
									}
									position++
									goto l285
								l286:
									position, tokenIndex, depth = position285, tokenIndex285, depth285
									if buffer[position] != rune('\\') {
										goto l278
									}
									position++
								}
							l285:
							}
						l280:
							depth--
							add(rulestringChar, position279)
						}
						goto l277
					l278:
						position, tokenIndex, depth = position278, tokenIndex278, depth278
					}
					depth--
					add(rulePegText, position276)
				}
				if buffer[position] != rune('"') {
					goto l274
				}
				position++
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleString, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
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
		/* 58 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position304, tokenIndex304, depth304 := position, tokenIndex, depth
			{
				position305 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l304
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l304
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l304
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l304
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l304
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l304
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l304
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position305)
			}
			return true
		l304:
			position, tokenIndex, depth = position304, tokenIndex304, depth304
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
