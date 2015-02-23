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
		/* 0 MongoLogLine <- <(Timestamp ' ' Thread ' ' Op ' ' NS ' ' LineField* Locks? LineField* Duration? extra !.)> */
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
					position76 := position
					depth++
					{
						position77 := position
						depth++
					l78:
						{
							position79, tokenIndex79, depth79 := position, tokenIndex, depth
							if !matchDot() {
								goto l79
							}
							goto l78
						l79:
							position, tokenIndex, depth = position79, tokenIndex79, depth79
						}
						depth--
						add(rulePegText, position77)
					}
					{
						add(ruleAction9, position)
					}
					depth--
					add(ruleextra, position76)
				}
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if !matchDot() {
						goto l81
					}
					goto l0
				l81:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
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
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					{
						position89 := position
						depth++
						if buffer[position] != rune('c') {
							goto l88
						}
						position++
						if buffer[position] != rune('o') {
							goto l88
						}
						position++
						if buffer[position] != rune('m') {
							goto l88
						}
						position++
						if buffer[position] != rune('m') {
							goto l88
						}
						position++
						if buffer[position] != rune('a') {
							goto l88
						}
						position++
						if buffer[position] != rune('n') {
							goto l88
						}
						position++
						if buffer[position] != rune('d') {
							goto l88
						}
						position++
						if buffer[position] != rune(':') {
							goto l88
						}
						position++
						if buffer[position] != rune(' ') {
							goto l88
						}
						position++
						{
							position90 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l88
							}
						l91:
							{
								position92, tokenIndex92, depth92 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l92
								}
								goto l91
							l92:
								position, tokenIndex, depth = position92, tokenIndex92, depth92
							}
							depth--
							add(rulePegText, position90)
						}
						{
							add(ruleAction6, position)
						}
						depth--
						add(rulecommandFieldName, position89)
					}
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					{
						position94 := position
						depth++
						{
							position95 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l85
							}
						l96:
							{
								position97, tokenIndex97, depth97 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l97
								}
								goto l96
							l97:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
							}
							depth--
							add(rulePegText, position95)
						}
						if buffer[position] != rune(':') {
							goto l85
						}
						position++
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleplainFieldName, position94)
					}
				}
			l87:
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
					position101 := position
					depth++
					{
						position102, tokenIndex102, depth102 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l102
						}
						goto l103
					l102:
						position, tokenIndex, depth = position102, tokenIndex102, depth102
					}
				l103:
					{
						position104, tokenIndex104, depth104 := position, tokenIndex, depth
						if !_rules[ruleDoc]() {
							goto l105
						}
						goto l104
					l105:
						position, tokenIndex, depth = position104, tokenIndex104, depth104
						if !_rules[ruleNumeric]() {
							goto l85
						}
					}
				l104:
					depth--
					add(ruleLineValue, position101)
				}
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
					add(ruleAction2, position)
				}
				depth--
				add(ruleLineField, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
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
		/* 11 LineValue <- <(S? (Doc / Numeric))> */
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
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l120
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l120
				}
				position++
				depth--
				add(ruledigit2, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 17 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 18 tz <- <('+' [0-9]+)> */
		nil,
		/* 19 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position124, tokenIndex124, depth124 := position, tokenIndex, depth
			{
				position125 := position
				depth++
				{
					position126 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l124
					}
					depth--
					add(rulehour, position126)
				}
				if buffer[position] != rune(':') {
					goto l124
				}
				position++
				{
					position127 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l124
					}
					depth--
					add(ruleminute, position127)
				}
				if buffer[position] != rune(':') {
					goto l124
				}
				position++
				{
					position128 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l124
					}
					depth--
					add(rulesecond, position128)
				}
				if buffer[position] != rune('.') {
					goto l124
				}
				position++
				{
					position129 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l124
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l124
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l124
					}
					position++
					depth--
					add(rulemillisecond, position129)
				}
				depth--
				add(ruletime, position125)
			}
			return true
		l124:
			position, tokenIndex, depth = position124, tokenIndex124, depth124
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
		/* 29 extra <- <(<.*> Action9)> */
		nil,
		/* 30 S <- <' '+> */
		func() bool {
			position140, tokenIndex140, depth140 := position, tokenIndex, depth
			{
				position141 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l140
				}
				position++
			l142:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
				}
				depth--
				add(ruleS, position141)
			}
			return true
		l140:
			position, tokenIndex, depth = position140, tokenIndex140, depth140
			return false
		},
		/* 31 Doc <- <('{' Action10 DocElements? '}' Action11)> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				if buffer[position] != rune('{') {
					goto l144
				}
				position++
				{
					add(ruleAction10, position)
				}
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					{
						position149 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l147
						}
					l150:
						{
							position151, tokenIndex151, depth151 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l151
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l151
							}
							goto l150
						l151:
							position, tokenIndex, depth = position151, tokenIndex151, depth151
						}
						depth--
						add(ruleDocElements, position149)
					}
					goto l148
				l147:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
				}
			l148:
				if buffer[position] != rune('}') {
					goto l144
				}
				position++
				{
					add(ruleAction11, position)
				}
				depth--
				add(ruleDoc, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 32 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 33 DocElem <- <(S? Field S? Value S? Action12)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
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
				{
					position158 := position
					depth++
					{
						position159 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l154
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
						goto l154
					}
					position++
					{
						add(ruleAction16, position)
					}
					depth--
					add(ruleField, position158)
				}
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
				if !_rules[ruleValue]() {
					goto l154
				}
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
					add(ruleAction12, position)
				}
				depth--
				add(ruleDocElem, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 34 List <- <('[' Action13 ListElements? ']' Action14)> */
		nil,
		/* 35 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 36 ListElem <- <(S? Value S? Action15)> */
		func() bool {
			position170, tokenIndex170, depth170 := position, tokenIndex, depth
			{
				position171 := position
				depth++
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l172
					}
					goto l173
				l172:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
				}
			l173:
				if !_rules[ruleValue]() {
					goto l170
				}
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l174
					}
					goto l175
				l174:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
				}
			l175:
				{
					add(ruleAction15, position)
				}
				depth--
				add(ruleListElem, position171)
			}
			return true
		l170:
			position, tokenIndex, depth = position170, tokenIndex170, depth170
			return false
		},
		/* 37 Field <- <(<fieldChar+> ':' Action16)> */
		nil,
		/* 38 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					{
						position182 := position
						depth++
						if buffer[position] != rune('n') {
							goto l181
						}
						position++
						if buffer[position] != rune('u') {
							goto l181
						}
						position++
						if buffer[position] != rune('l') {
							goto l181
						}
						position++
						if buffer[position] != rune('l') {
							goto l181
						}
						position++
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleNull, position182)
					}
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					{
						position185 := position
						depth++
						if buffer[position] != rune('M') {
							goto l184
						}
						position++
						if buffer[position] != rune('i') {
							goto l184
						}
						position++
						if buffer[position] != rune('n') {
							goto l184
						}
						position++
						if buffer[position] != rune('K') {
							goto l184
						}
						position++
						if buffer[position] != rune('e') {
							goto l184
						}
						position++
						if buffer[position] != rune('y') {
							goto l184
						}
						position++
						{
							add(ruleAction28, position)
						}
						depth--
						add(ruleMinKey, position185)
					}
					goto l180
				l184:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					{
						switch buffer[position] {
						case 'M':
							{
								position188 := position
								depth++
								if buffer[position] != rune('M') {
									goto l178
								}
								position++
								if buffer[position] != rune('a') {
									goto l178
								}
								position++
								if buffer[position] != rune('x') {
									goto l178
								}
								position++
								if buffer[position] != rune('K') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('y') {
									goto l178
								}
								position++
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleMaxKey, position188)
							}
							break
						case 'u':
							{
								position190 := position
								depth++
								if buffer[position] != rune('u') {
									goto l178
								}
								position++
								if buffer[position] != rune('n') {
									goto l178
								}
								position++
								if buffer[position] != rune('d') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('f') {
									goto l178
								}
								position++
								if buffer[position] != rune('i') {
									goto l178
								}
								position++
								if buffer[position] != rune('n') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('d') {
									goto l178
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleUndefined, position190)
							}
							break
						case 'N':
							{
								position192 := position
								depth++
								if buffer[position] != rune('N') {
									goto l178
								}
								position++
								if buffer[position] != rune('u') {
									goto l178
								}
								position++
								if buffer[position] != rune('m') {
									goto l178
								}
								position++
								if buffer[position] != rune('b') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('r') {
									goto l178
								}
								position++
								if buffer[position] != rune('L') {
									goto l178
								}
								position++
								if buffer[position] != rune('o') {
									goto l178
								}
								position++
								if buffer[position] != rune('n') {
									goto l178
								}
								position++
								if buffer[position] != rune('g') {
									goto l178
								}
								position++
								if buffer[position] != rune('(') {
									goto l178
								}
								position++
								{
									position193 := position
									depth++
									{
										position196, tokenIndex196, depth196 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l196
										}
										position++
										goto l178
									l196:
										position, tokenIndex, depth = position196, tokenIndex196, depth196
									}
									if !matchDot() {
										goto l178
									}
								l194:
									{
										position195, tokenIndex195, depth195 := position, tokenIndex, depth
										{
											position197, tokenIndex197, depth197 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l197
											}
											position++
											goto l195
										l197:
											position, tokenIndex, depth = position197, tokenIndex197, depth197
										}
										if !matchDot() {
											goto l195
										}
										goto l194
									l195:
										position, tokenIndex, depth = position195, tokenIndex195, depth195
									}
									depth--
									add(rulePegText, position193)
								}
								if buffer[position] != rune(')') {
									goto l178
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleNumberLong, position192)
							}
							break
						case '/':
							{
								position199 := position
								depth++
								if buffer[position] != rune('/') {
									goto l178
								}
								position++
								{
									position200 := position
									depth++
									{
										position201 := position
										depth++
										{
											position204 := position
											depth++
											{
												position205, tokenIndex205, depth205 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l205
												}
												position++
												goto l178
											l205:
												position, tokenIndex, depth = position205, tokenIndex205, depth205
											}
											if !matchDot() {
												goto l178
											}
											depth--
											add(ruleregexChar, position204)
										}
									l202:
										{
											position203, tokenIndex203, depth203 := position, tokenIndex, depth
											{
												position206 := position
												depth++
												{
													position207, tokenIndex207, depth207 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l207
													}
													position++
													goto l203
												l207:
													position, tokenIndex, depth = position207, tokenIndex207, depth207
												}
												if !matchDot() {
													goto l203
												}
												depth--
												add(ruleregexChar, position206)
											}
											goto l202
										l203:
											position, tokenIndex, depth = position203, tokenIndex203, depth203
										}
										if buffer[position] != rune('/') {
											goto l178
										}
										position++
									l208:
										{
											position209, tokenIndex209, depth209 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l209
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l209
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l209
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l209
													}
													position++
													break
												}
											}

											goto l208
										l209:
											position, tokenIndex, depth = position209, tokenIndex209, depth209
										}
										depth--
										add(ruleregexBody, position201)
									}
									depth--
									add(rulePegText, position200)
								}
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleRegex, position199)
							}
							break
						case 'T':
							{
								position212 := position
								depth++
								if buffer[position] != rune('T') {
									goto l178
								}
								position++
								if buffer[position] != rune('i') {
									goto l178
								}
								position++
								if buffer[position] != rune('m') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('s') {
									goto l178
								}
								position++
								if buffer[position] != rune('t') {
									goto l178
								}
								position++
								if buffer[position] != rune('a') {
									goto l178
								}
								position++
								if buffer[position] != rune('m') {
									goto l178
								}
								position++
								if buffer[position] != rune('p') {
									goto l178
								}
								position++
								if buffer[position] != rune('(') {
									goto l178
								}
								position++
								{
									position213 := position
									depth++
									{
										position216, tokenIndex216, depth216 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l216
										}
										position++
										goto l178
									l216:
										position, tokenIndex, depth = position216, tokenIndex216, depth216
									}
									if !matchDot() {
										goto l178
									}
								l214:
									{
										position215, tokenIndex215, depth215 := position, tokenIndex, depth
										{
											position217, tokenIndex217, depth217 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l217
											}
											position++
											goto l215
										l217:
											position, tokenIndex, depth = position217, tokenIndex217, depth217
										}
										if !matchDot() {
											goto l215
										}
										goto l214
									l215:
										position, tokenIndex, depth = position215, tokenIndex215, depth215
									}
									depth--
									add(rulePegText, position213)
								}
								if buffer[position] != rune(')') {
									goto l178
								}
								position++
								{
									add(ruleAction26, position)
								}
								depth--
								add(ruleTimestampVal, position212)
							}
							break
						case 'B':
							{
								position219 := position
								depth++
								if buffer[position] != rune('B') {
									goto l178
								}
								position++
								if buffer[position] != rune('i') {
									goto l178
								}
								position++
								if buffer[position] != rune('n') {
									goto l178
								}
								position++
								if buffer[position] != rune('D') {
									goto l178
								}
								position++
								if buffer[position] != rune('a') {
									goto l178
								}
								position++
								if buffer[position] != rune('t') {
									goto l178
								}
								position++
								if buffer[position] != rune('a') {
									goto l178
								}
								position++
								if buffer[position] != rune('(') {
									goto l178
								}
								position++
								{
									position220 := position
									depth++
									{
										position223, tokenIndex223, depth223 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l223
										}
										position++
										goto l178
									l223:
										position, tokenIndex, depth = position223, tokenIndex223, depth223
									}
									if !matchDot() {
										goto l178
									}
								l221:
									{
										position222, tokenIndex222, depth222 := position, tokenIndex, depth
										{
											position224, tokenIndex224, depth224 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l224
											}
											position++
											goto l222
										l224:
											position, tokenIndex, depth = position224, tokenIndex224, depth224
										}
										if !matchDot() {
											goto l222
										}
										goto l221
									l222:
										position, tokenIndex, depth = position222, tokenIndex222, depth222
									}
									depth--
									add(rulePegText, position220)
								}
								if buffer[position] != rune(')') {
									goto l178
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleBinData, position219)
							}
							break
						case 'D', 'n':
							{
								position226 := position
								depth++
								{
									position227, tokenIndex227, depth227 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l227
									}
									position++
									if buffer[position] != rune('e') {
										goto l227
									}
									position++
									if buffer[position] != rune('w') {
										goto l227
									}
									position++
									if buffer[position] != rune(' ') {
										goto l227
									}
									position++
									goto l228
								l227:
									position, tokenIndex, depth = position227, tokenIndex227, depth227
								}
							l228:
								if buffer[position] != rune('D') {
									goto l178
								}
								position++
								if buffer[position] != rune('a') {
									goto l178
								}
								position++
								if buffer[position] != rune('t') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('(') {
									goto l178
								}
								position++
								{
									position229 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l178
									}
									position++
								l230:
									{
										position231, tokenIndex231, depth231 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l231
										}
										position++
										goto l230
									l231:
										position, tokenIndex, depth = position231, tokenIndex231, depth231
									}
									depth--
									add(rulePegText, position229)
								}
								if buffer[position] != rune(')') {
									goto l178
								}
								position++
								{
									add(ruleAction22, position)
								}
								depth--
								add(ruleDate, position226)
							}
							break
						case 'O':
							{
								position233 := position
								depth++
								if buffer[position] != rune('O') {
									goto l178
								}
								position++
								if buffer[position] != rune('b') {
									goto l178
								}
								position++
								if buffer[position] != rune('j') {
									goto l178
								}
								position++
								if buffer[position] != rune('e') {
									goto l178
								}
								position++
								if buffer[position] != rune('c') {
									goto l178
								}
								position++
								if buffer[position] != rune('t') {
									goto l178
								}
								position++
								if buffer[position] != rune('I') {
									goto l178
								}
								position++
								if buffer[position] != rune('d') {
									goto l178
								}
								position++
								if buffer[position] != rune('(') {
									goto l178
								}
								position++
								if buffer[position] != rune('"') {
									goto l178
								}
								position++
								{
									position234 := position
									depth++
								l235:
									{
										position236, tokenIndex236, depth236 := position, tokenIndex, depth
										{
											position237 := position
											depth++
											{
												position238, tokenIndex238, depth238 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l239
												}
												position++
												goto l238
											l239:
												position, tokenIndex, depth = position238, tokenIndex238, depth238
												{
													position240, tokenIndex240, depth240 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l241
													}
													position++
													goto l240
												l241:
													position, tokenIndex, depth = position240, tokenIndex240, depth240
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l236
													}
													position++
												}
											l240:
											}
										l238:
											depth--
											add(rulehexChar, position237)
										}
										goto l235
									l236:
										position, tokenIndex, depth = position236, tokenIndex236, depth236
									}
									depth--
									add(rulePegText, position234)
								}
								if buffer[position] != rune('"') {
									goto l178
								}
								position++
								if buffer[position] != rune(')') {
									goto l178
								}
								position++
								{
									add(ruleAction23, position)
								}
								depth--
								add(ruleObjectID, position233)
							}
							break
						case '"':
							{
								position243 := position
								depth++
								if buffer[position] != rune('"') {
									goto l178
								}
								position++
								{
									position244 := position
									depth++
								l245:
									{
										position246, tokenIndex246, depth246 := position, tokenIndex, depth
										{
											position247 := position
											depth++
											{
												position248, tokenIndex248, depth248 := position, tokenIndex, depth
												{
													position250, tokenIndex250, depth250 := position, tokenIndex, depth
													{
														position251, tokenIndex251, depth251 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l252
														}
														position++
														goto l251
													l252:
														position, tokenIndex, depth = position251, tokenIndex251, depth251
														if buffer[position] != rune('\\') {
															goto l250
														}
														position++
													}
												l251:
													goto l249
												l250:
													position, tokenIndex, depth = position250, tokenIndex250, depth250
												}
												if !matchDot() {
													goto l249
												}
												goto l248
											l249:
												position, tokenIndex, depth = position248, tokenIndex248, depth248
												if buffer[position] != rune('\\') {
													goto l246
												}
												position++
												{
													position253, tokenIndex253, depth253 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l254
													}
													position++
													goto l253
												l254:
													position, tokenIndex, depth = position253, tokenIndex253, depth253
													if buffer[position] != rune('\\') {
														goto l246
													}
													position++
												}
											l253:
											}
										l248:
											depth--
											add(rulestringChar, position247)
										}
										goto l245
									l246:
										position, tokenIndex, depth = position246, tokenIndex246, depth246
									}
									depth--
									add(rulePegText, position244)
								}
								if buffer[position] != rune('"') {
									goto l178
								}
								position++
								{
									add(ruleAction18, position)
								}
								depth--
								add(ruleString, position243)
							}
							break
						case 'f', 't':
							{
								position256 := position
								depth++
								{
									position257, tokenIndex257, depth257 := position, tokenIndex, depth
									{
										position259 := position
										depth++
										if buffer[position] != rune('t') {
											goto l258
										}
										position++
										if buffer[position] != rune('r') {
											goto l258
										}
										position++
										if buffer[position] != rune('u') {
											goto l258
										}
										position++
										if buffer[position] != rune('e') {
											goto l258
										}
										position++
										{
											add(ruleAction20, position)
										}
										depth--
										add(ruleTrue, position259)
									}
									goto l257
								l258:
									position, tokenIndex, depth = position257, tokenIndex257, depth257
									{
										position261 := position
										depth++
										if buffer[position] != rune('f') {
											goto l178
										}
										position++
										if buffer[position] != rune('a') {
											goto l178
										}
										position++
										if buffer[position] != rune('l') {
											goto l178
										}
										position++
										if buffer[position] != rune('s') {
											goto l178
										}
										position++
										if buffer[position] != rune('e') {
											goto l178
										}
										position++
										{
											add(ruleAction21, position)
										}
										depth--
										add(ruleFalse, position261)
									}
								}
							l257:
								depth--
								add(ruleBoolean, position256)
							}
							break
						case '[':
							{
								position263 := position
								depth++
								if buffer[position] != rune('[') {
									goto l178
								}
								position++
								{
									add(ruleAction13, position)
								}
								{
									position265, tokenIndex265, depth265 := position, tokenIndex, depth
									{
										position267 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l265
										}
									l268:
										{
											position269, tokenIndex269, depth269 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l269
											}
											position++
											if !_rules[ruleListElem]() {
												goto l269
											}
											goto l268
										l269:
											position, tokenIndex, depth = position269, tokenIndex269, depth269
										}
										depth--
										add(ruleListElements, position267)
									}
									goto l266
								l265:
									position, tokenIndex, depth = position265, tokenIndex265, depth265
								}
							l266:
								if buffer[position] != rune(']') {
									goto l178
								}
								position++
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleList, position263)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l178
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l178
							}
							break
						}
					}

				}
			l180:
				depth--
				add(ruleValue, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 39 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action17)> */
		func() bool {
			position271, tokenIndex271, depth271 := position, tokenIndex, depth
			{
				position272 := position
				depth++
				{
					position273 := position
					depth++
					{
						position274, tokenIndex274, depth274 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l274
						}
						position++
						goto l275
					l274:
						position, tokenIndex, depth = position274, tokenIndex274, depth274
					}
				l275:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l271
					}
					position++
				l276:
					{
						position277, tokenIndex277, depth277 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l277
						}
						position++
						goto l276
					l277:
						position, tokenIndex, depth = position277, tokenIndex277, depth277
					}
					{
						position278, tokenIndex278, depth278 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l278
						}
						position++
						goto l279
					l278:
						position, tokenIndex, depth = position278, tokenIndex278, depth278
					}
				l279:
				l280:
					{
						position281, tokenIndex281, depth281 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l281
						}
						position++
						goto l280
					l281:
						position, tokenIndex, depth = position281, tokenIndex281, depth281
					}
					depth--
					add(rulePegText, position273)
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleNumeric, position272)
			}
			return true
		l271:
			position, tokenIndex, depth = position271, tokenIndex271, depth271
			return false
		},
		/* 40 Boolean <- <(True / False)> */
		nil,
		/* 41 String <- <('"' <stringChar*> '"' Action18)> */
		nil,
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
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				{
					switch buffer[position] {
					case '$', '_':
						{
							position304, tokenIndex304, depth304 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l305
							}
							position++
							goto l304
						l305:
							position, tokenIndex, depth = position304, tokenIndex304, depth304
							if buffer[position] != rune('$') {
								goto l301
							}
							position++
						}
					l304:
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l301
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l301
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l301
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
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
