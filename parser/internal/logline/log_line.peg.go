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
	ruleplanSummaryField
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
	"planSummaryField",
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
	rules  [94]func() bool
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
			p.SetField("planSummaryType", buffer[begin:end])
			p.StartField("planSummary")
		case ruleAction8:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction9:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction10:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction11:
			p.PushMap()
		case ruleAction12:
			p.PopMap()
		case ruleAction13:
			p.SetMapValue()
		case ruleAction14:
			p.PushList()
		case ruleAction15:
			p.PopList()
		case ruleAction16:
			p.SetListValue()
		case ruleAction17:
			p.PushField(buffer[begin:end])
		case ruleAction18:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction19:
			p.PushValue(buffer[begin:end])
		case ruleAction20:
			p.PushValue(nil)
		case ruleAction21:
			p.PushValue(true)
		case ruleAction22:
			p.PushValue(false)
		case ruleAction23:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction24:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction25:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction26:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction27:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction28:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction29:
			p.PushValue(p.Minkey())
		case ruleAction30:
			p.PushValue(p.Maxkey())
		case ruleAction31:
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
								add(ruleAction8, position)
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
								add(ruleAction9, position)
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
							add(ruleAction10, position)
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
		/* 4 LineField <- <((commandFieldName / planSummaryField / plainFieldName) S? LineValue S? Action2)> */
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
						position99 := position
						depth++
						if buffer[position] != rune('p') {
							goto l98
						}
						position++
						if buffer[position] != rune('l') {
							goto l98
						}
						position++
						if buffer[position] != rune('a') {
							goto l98
						}
						position++
						if buffer[position] != rune('n') {
							goto l98
						}
						position++
						if buffer[position] != rune('S') {
							goto l98
						}
						position++
						if buffer[position] != rune('u') {
							goto l98
						}
						position++
						if buffer[position] != rune('m') {
							goto l98
						}
						position++
						if buffer[position] != rune('m') {
							goto l98
						}
						position++
						if buffer[position] != rune('a') {
							goto l98
						}
						position++
						if buffer[position] != rune('r') {
							goto l98
						}
						position++
						if buffer[position] != rune('y') {
							goto l98
						}
						position++
						if buffer[position] != rune(':') {
							goto l98
						}
						position++
						if buffer[position] != rune(' ') {
							goto l98
						}
						position++
						{
							position100 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l98
							}
						l101:
							{
								position102, tokenIndex102, depth102 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l102
								}
								goto l101
							l102:
								position, tokenIndex, depth = position102, tokenIndex102, depth102
							}
							depth--
							add(rulePegText, position100)
						}
						{
							add(ruleAction7, position)
						}
						depth--
						add(ruleplanSummaryField, position99)
					}
					goto l91
				l98:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					{
						position104 := position
						depth++
						{
							position105 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l89
							}
						l106:
							{
								position107, tokenIndex107, depth107 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l107
								}
								goto l106
							l107:
								position, tokenIndex, depth = position107, tokenIndex107, depth107
							}
							depth--
							add(rulePegText, position105)
						}
						if buffer[position] != rune(':') {
							goto l89
						}
						position++
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleplainFieldName, position104)
					}
				}
			l91:
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
					position111 := position
					depth++
					{
						position112, tokenIndex112, depth112 := position, tokenIndex, depth
						if !_rules[ruleS]() {
							goto l112
						}
						goto l113
					l112:
						position, tokenIndex, depth = position112, tokenIndex112, depth112
					}
				l113:
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
					add(ruleLineValue, position111)
				}
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l115
					}
					goto l116
				l115:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
				}
			l116:
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
		/* 11 planSummaryField <- <('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y' ':' ' ' <fieldChar+> Action7)> */
		nil,
		/* 12 LineValue <- <(S? ((&('"') String) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		nil,
		/* 13 timestamp24 <- <(<(date ' ' time)> Action8)> */
		nil,
		/* 14 timestamp26 <- <(<datetime26> Action9)> */
		nil,
		/* 15 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 16 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 17 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position130, tokenIndex130, depth130 := position, tokenIndex, depth
			{
				position131 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l130
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l130
				}
				position++
				depth--
				add(ruledigit2, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 18 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 19 tz <- <('+' [0-9]+)> */
		nil,
		/* 20 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position134, tokenIndex134, depth134 := position, tokenIndex, depth
			{
				position135 := position
				depth++
				{
					position136 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l134
					}
					depth--
					add(rulehour, position136)
				}
				if buffer[position] != rune(':') {
					goto l134
				}
				position++
				{
					position137 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l134
					}
					depth--
					add(ruleminute, position137)
				}
				if buffer[position] != rune(':') {
					goto l134
				}
				position++
				{
					position138 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l134
					}
					depth--
					add(rulesecond, position138)
				}
				if buffer[position] != rune('.') {
					goto l134
				}
				position++
				{
					position139 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l134
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l134
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l134
					}
					position++
					depth--
					add(rulemillisecond, position139)
				}
				depth--
				add(ruletime, position135)
			}
			return true
		l134:
			position, tokenIndex, depth = position134, tokenIndex134, depth134
			return false
		},
		/* 21 day <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 22 month <- <([A-Z] [a-z] [a-z])> */
		nil,
		/* 23 dayNum <- <digit2?> */
		nil,
		/* 24 hour <- <digit2> */
		nil,
		/* 25 minute <- <digit2> */
		nil,
		/* 26 second <- <digit2> */
		nil,
		/* 27 millisecond <- <([0-9] [0-9] [0-9])> */
		nil,
		/* 28 letterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 29 nsChar <- <((&('$') '$') | (&(':') ':') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '[' | '\\' | ']' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [A-z]))> */
		nil,
		/* 30 extra <- <(<.+> Action10)> */
		nil,
		/* 31 S <- <' '+> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l150
				}
				position++
			l152:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
				}
				depth--
				add(ruleS, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
			return false
		},
		/* 32 Doc <- <('{' Action11 DocElements? '}' Action12)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				if buffer[position] != rune('{') {
					goto l154
				}
				position++
				{
					add(ruleAction11, position)
				}
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					{
						position159 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l157
						}
					l160:
						{
							position161, tokenIndex161, depth161 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l161
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l161
							}
							goto l160
						l161:
							position, tokenIndex, depth = position161, tokenIndex161, depth161
						}
						depth--
						add(ruleDocElements, position159)
					}
					goto l158
				l157:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
				}
			l158:
				if buffer[position] != rune('}') {
					goto l154
				}
				position++
				{
					add(ruleAction12, position)
				}
				depth--
				add(ruleDoc, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 33 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 34 DocElem <- <(S? Field S? Value S? Action13)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
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
					position168 := position
					depth++
					{
						position169 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l164
						}
					l170:
						{
							position171, tokenIndex171, depth171 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l171
							}
							goto l170
						l171:
							position, tokenIndex, depth = position171, tokenIndex171, depth171
						}
						depth--
						add(rulePegText, position169)
					}
					if buffer[position] != rune(':') {
						goto l164
					}
					position++
					{
						add(ruleAction17, position)
					}
					depth--
					add(ruleField, position168)
				}
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
					goto l164
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
					add(ruleAction13, position)
				}
				depth--
				add(ruleDocElem, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 35 List <- <('[' Action14 ListElements? ']' Action15)> */
		nil,
		/* 36 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 37 ListElem <- <(S? Value S? Action16)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l182
					}
					goto l183
				l182:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
				}
			l183:
				if !_rules[ruleValue]() {
					goto l180
				}
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l184
					}
					goto l185
				l184:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
				}
			l185:
				{
					add(ruleAction16, position)
				}
				depth--
				add(ruleListElem, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 38 Field <- <(<fieldChar+> ':' Action17)> */
		nil,
		/* 39 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					{
						position192 := position
						depth++
						if buffer[position] != rune('n') {
							goto l191
						}
						position++
						if buffer[position] != rune('u') {
							goto l191
						}
						position++
						if buffer[position] != rune('l') {
							goto l191
						}
						position++
						if buffer[position] != rune('l') {
							goto l191
						}
						position++
						{
							add(ruleAction20, position)
						}
						depth--
						add(ruleNull, position192)
					}
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					{
						position195 := position
						depth++
						if buffer[position] != rune('M') {
							goto l194
						}
						position++
						if buffer[position] != rune('i') {
							goto l194
						}
						position++
						if buffer[position] != rune('n') {
							goto l194
						}
						position++
						if buffer[position] != rune('K') {
							goto l194
						}
						position++
						if buffer[position] != rune('e') {
							goto l194
						}
						position++
						if buffer[position] != rune('y') {
							goto l194
						}
						position++
						{
							add(ruleAction29, position)
						}
						depth--
						add(ruleMinKey, position195)
					}
					goto l190
				l194:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					{
						switch buffer[position] {
						case 'M':
							{
								position198 := position
								depth++
								if buffer[position] != rune('M') {
									goto l188
								}
								position++
								if buffer[position] != rune('a') {
									goto l188
								}
								position++
								if buffer[position] != rune('x') {
									goto l188
								}
								position++
								if buffer[position] != rune('K') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('y') {
									goto l188
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleMaxKey, position198)
							}
							break
						case 'u':
							{
								position200 := position
								depth++
								if buffer[position] != rune('u') {
									goto l188
								}
								position++
								if buffer[position] != rune('n') {
									goto l188
								}
								position++
								if buffer[position] != rune('d') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('f') {
									goto l188
								}
								position++
								if buffer[position] != rune('i') {
									goto l188
								}
								position++
								if buffer[position] != rune('n') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('d') {
									goto l188
								}
								position++
								{
									add(ruleAction31, position)
								}
								depth--
								add(ruleUndefined, position200)
							}
							break
						case 'N':
							{
								position202 := position
								depth++
								if buffer[position] != rune('N') {
									goto l188
								}
								position++
								if buffer[position] != rune('u') {
									goto l188
								}
								position++
								if buffer[position] != rune('m') {
									goto l188
								}
								position++
								if buffer[position] != rune('b') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('r') {
									goto l188
								}
								position++
								if buffer[position] != rune('L') {
									goto l188
								}
								position++
								if buffer[position] != rune('o') {
									goto l188
								}
								position++
								if buffer[position] != rune('n') {
									goto l188
								}
								position++
								if buffer[position] != rune('g') {
									goto l188
								}
								position++
								if buffer[position] != rune('(') {
									goto l188
								}
								position++
								{
									position203 := position
									depth++
									{
										position206, tokenIndex206, depth206 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l206
										}
										position++
										goto l188
									l206:
										position, tokenIndex, depth = position206, tokenIndex206, depth206
									}
									if !matchDot() {
										goto l188
									}
								l204:
									{
										position205, tokenIndex205, depth205 := position, tokenIndex, depth
										{
											position207, tokenIndex207, depth207 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l207
											}
											position++
											goto l205
										l207:
											position, tokenIndex, depth = position207, tokenIndex207, depth207
										}
										if !matchDot() {
											goto l205
										}
										goto l204
									l205:
										position, tokenIndex, depth = position205, tokenIndex205, depth205
									}
									depth--
									add(rulePegText, position203)
								}
								if buffer[position] != rune(')') {
									goto l188
								}
								position++
								{
									add(ruleAction28, position)
								}
								depth--
								add(ruleNumberLong, position202)
							}
							break
						case '/':
							{
								position209 := position
								depth++
								if buffer[position] != rune('/') {
									goto l188
								}
								position++
								{
									position210 := position
									depth++
									{
										position211 := position
										depth++
										{
											position214 := position
											depth++
											{
												position215, tokenIndex215, depth215 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l215
												}
												position++
												goto l188
											l215:
												position, tokenIndex, depth = position215, tokenIndex215, depth215
											}
											if !matchDot() {
												goto l188
											}
											depth--
											add(ruleregexChar, position214)
										}
									l212:
										{
											position213, tokenIndex213, depth213 := position, tokenIndex, depth
											{
												position216 := position
												depth++
												{
													position217, tokenIndex217, depth217 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l217
													}
													position++
													goto l213
												l217:
													position, tokenIndex, depth = position217, tokenIndex217, depth217
												}
												if !matchDot() {
													goto l213
												}
												depth--
												add(ruleregexChar, position216)
											}
											goto l212
										l213:
											position, tokenIndex, depth = position213, tokenIndex213, depth213
										}
										if buffer[position] != rune('/') {
											goto l188
										}
										position++
									l218:
										{
											position219, tokenIndex219, depth219 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l219
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l219
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l219
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l219
													}
													position++
													break
												}
											}

											goto l218
										l219:
											position, tokenIndex, depth = position219, tokenIndex219, depth219
										}
										depth--
										add(ruleregexBody, position211)
									}
									depth--
									add(rulePegText, position210)
								}
								{
									add(ruleAction26, position)
								}
								depth--
								add(ruleRegex, position209)
							}
							break
						case 'T':
							{
								position222 := position
								depth++
								if buffer[position] != rune('T') {
									goto l188
								}
								position++
								if buffer[position] != rune('i') {
									goto l188
								}
								position++
								if buffer[position] != rune('m') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('s') {
									goto l188
								}
								position++
								if buffer[position] != rune('t') {
									goto l188
								}
								position++
								if buffer[position] != rune('a') {
									goto l188
								}
								position++
								if buffer[position] != rune('m') {
									goto l188
								}
								position++
								if buffer[position] != rune('p') {
									goto l188
								}
								position++
								if buffer[position] != rune('(') {
									goto l188
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
										goto l188
									l226:
										position, tokenIndex, depth = position226, tokenIndex226, depth226
									}
									if !matchDot() {
										goto l188
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
									goto l188
								}
								position++
								{
									add(ruleAction27, position)
								}
								depth--
								add(ruleTimestampVal, position222)
							}
							break
						case 'B':
							{
								position229 := position
								depth++
								if buffer[position] != rune('B') {
									goto l188
								}
								position++
								if buffer[position] != rune('i') {
									goto l188
								}
								position++
								if buffer[position] != rune('n') {
									goto l188
								}
								position++
								if buffer[position] != rune('D') {
									goto l188
								}
								position++
								if buffer[position] != rune('a') {
									goto l188
								}
								position++
								if buffer[position] != rune('t') {
									goto l188
								}
								position++
								if buffer[position] != rune('a') {
									goto l188
								}
								position++
								if buffer[position] != rune('(') {
									goto l188
								}
								position++
								{
									position230 := position
									depth++
									{
										position233, tokenIndex233, depth233 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l233
										}
										position++
										goto l188
									l233:
										position, tokenIndex, depth = position233, tokenIndex233, depth233
									}
									if !matchDot() {
										goto l188
									}
								l231:
									{
										position232, tokenIndex232, depth232 := position, tokenIndex, depth
										{
											position234, tokenIndex234, depth234 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l234
											}
											position++
											goto l232
										l234:
											position, tokenIndex, depth = position234, tokenIndex234, depth234
										}
										if !matchDot() {
											goto l232
										}
										goto l231
									l232:
										position, tokenIndex, depth = position232, tokenIndex232, depth232
									}
									depth--
									add(rulePegText, position230)
								}
								if buffer[position] != rune(')') {
									goto l188
								}
								position++
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleBinData, position229)
							}
							break
						case 'D', 'n':
							{
								position236 := position
								depth++
								{
									position237, tokenIndex237, depth237 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l237
									}
									position++
									if buffer[position] != rune('e') {
										goto l237
									}
									position++
									if buffer[position] != rune('w') {
										goto l237
									}
									position++
									if buffer[position] != rune(' ') {
										goto l237
									}
									position++
									goto l238
								l237:
									position, tokenIndex, depth = position237, tokenIndex237, depth237
								}
							l238:
								if buffer[position] != rune('D') {
									goto l188
								}
								position++
								if buffer[position] != rune('a') {
									goto l188
								}
								position++
								if buffer[position] != rune('t') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('(') {
									goto l188
								}
								position++
								{
									position239 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l188
									}
									position++
								l240:
									{
										position241, tokenIndex241, depth241 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l241
										}
										position++
										goto l240
									l241:
										position, tokenIndex, depth = position241, tokenIndex241, depth241
									}
									depth--
									add(rulePegText, position239)
								}
								if buffer[position] != rune(')') {
									goto l188
								}
								position++
								{
									add(ruleAction23, position)
								}
								depth--
								add(ruleDate, position236)
							}
							break
						case 'O':
							{
								position243 := position
								depth++
								if buffer[position] != rune('O') {
									goto l188
								}
								position++
								if buffer[position] != rune('b') {
									goto l188
								}
								position++
								if buffer[position] != rune('j') {
									goto l188
								}
								position++
								if buffer[position] != rune('e') {
									goto l188
								}
								position++
								if buffer[position] != rune('c') {
									goto l188
								}
								position++
								if buffer[position] != rune('t') {
									goto l188
								}
								position++
								if buffer[position] != rune('I') {
									goto l188
								}
								position++
								if buffer[position] != rune('d') {
									goto l188
								}
								position++
								if buffer[position] != rune('(') {
									goto l188
								}
								position++
								if buffer[position] != rune('"') {
									goto l188
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
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l249
												}
												position++
												goto l248
											l249:
												position, tokenIndex, depth = position248, tokenIndex248, depth248
												{
													position250, tokenIndex250, depth250 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l251
													}
													position++
													goto l250
												l251:
													position, tokenIndex, depth = position250, tokenIndex250, depth250
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l246
													}
													position++
												}
											l250:
											}
										l248:
											depth--
											add(rulehexChar, position247)
										}
										goto l245
									l246:
										position, tokenIndex, depth = position246, tokenIndex246, depth246
									}
									depth--
									add(rulePegText, position244)
								}
								if buffer[position] != rune('"') {
									goto l188
								}
								position++
								if buffer[position] != rune(')') {
									goto l188
								}
								position++
								{
									add(ruleAction24, position)
								}
								depth--
								add(ruleObjectID, position243)
							}
							break
						case '"':
							if !_rules[ruleString]() {
								goto l188
							}
							break
						case 'f', 't':
							{
								position253 := position
								depth++
								{
									position254, tokenIndex254, depth254 := position, tokenIndex, depth
									{
										position256 := position
										depth++
										if buffer[position] != rune('t') {
											goto l255
										}
										position++
										if buffer[position] != rune('r') {
											goto l255
										}
										position++
										if buffer[position] != rune('u') {
											goto l255
										}
										position++
										if buffer[position] != rune('e') {
											goto l255
										}
										position++
										{
											add(ruleAction21, position)
										}
										depth--
										add(ruleTrue, position256)
									}
									goto l254
								l255:
									position, tokenIndex, depth = position254, tokenIndex254, depth254
									{
										position258 := position
										depth++
										if buffer[position] != rune('f') {
											goto l188
										}
										position++
										if buffer[position] != rune('a') {
											goto l188
										}
										position++
										if buffer[position] != rune('l') {
											goto l188
										}
										position++
										if buffer[position] != rune('s') {
											goto l188
										}
										position++
										if buffer[position] != rune('e') {
											goto l188
										}
										position++
										{
											add(ruleAction22, position)
										}
										depth--
										add(ruleFalse, position258)
									}
								}
							l254:
								depth--
								add(ruleBoolean, position253)
							}
							break
						case '[':
							{
								position260 := position
								depth++
								if buffer[position] != rune('[') {
									goto l188
								}
								position++
								{
									add(ruleAction14, position)
								}
								{
									position262, tokenIndex262, depth262 := position, tokenIndex, depth
									{
										position264 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l262
										}
									l265:
										{
											position266, tokenIndex266, depth266 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l266
											}
											position++
											if !_rules[ruleListElem]() {
												goto l266
											}
											goto l265
										l266:
											position, tokenIndex, depth = position266, tokenIndex266, depth266
										}
										depth--
										add(ruleListElements, position264)
									}
									goto l263
								l262:
									position, tokenIndex, depth = position262, tokenIndex262, depth262
								}
							l263:
								if buffer[position] != rune(']') {
									goto l188
								}
								position++
								{
									add(ruleAction15, position)
								}
								depth--
								add(ruleList, position260)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l188
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l188
							}
							break
						}
					}

				}
			l190:
				depth--
				add(ruleValue, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 40 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action18)> */
		func() bool {
			position268, tokenIndex268, depth268 := position, tokenIndex, depth
			{
				position269 := position
				depth++
				{
					position270 := position
					depth++
					{
						position271, tokenIndex271, depth271 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l271
						}
						position++
						goto l272
					l271:
						position, tokenIndex, depth = position271, tokenIndex271, depth271
					}
				l272:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l268
					}
					position++
				l273:
					{
						position274, tokenIndex274, depth274 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l274
						}
						position++
						goto l273
					l274:
						position, tokenIndex, depth = position274, tokenIndex274, depth274
					}
					{
						position275, tokenIndex275, depth275 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l275
						}
						position++
						goto l276
					l275:
						position, tokenIndex, depth = position275, tokenIndex275, depth275
					}
				l276:
				l277:
					{
						position278, tokenIndex278, depth278 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l278
						}
						position++
						goto l277
					l278:
						position, tokenIndex, depth = position278, tokenIndex278, depth278
					}
					depth--
					add(rulePegText, position270)
				}
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleNumeric, position269)
			}
			return true
		l268:
			position, tokenIndex, depth = position268, tokenIndex268, depth268
			return false
		},
		/* 41 Boolean <- <(True / False)> */
		nil,
		/* 42 String <- <('"' <stringChar*> '"' Action19)> */
		func() bool {
			position281, tokenIndex281, depth281 := position, tokenIndex, depth
			{
				position282 := position
				depth++
				if buffer[position] != rune('"') {
					goto l281
				}
				position++
				{
					position283 := position
					depth++
				l284:
					{
						position285, tokenIndex285, depth285 := position, tokenIndex, depth
						{
							position286 := position
							depth++
							{
								position287, tokenIndex287, depth287 := position, tokenIndex, depth
								{
									position289, tokenIndex289, depth289 := position, tokenIndex, depth
									{
										position290, tokenIndex290, depth290 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l291
										}
										position++
										goto l290
									l291:
										position, tokenIndex, depth = position290, tokenIndex290, depth290
										if buffer[position] != rune('\\') {
											goto l289
										}
										position++
									}
								l290:
									goto l288
								l289:
									position, tokenIndex, depth = position289, tokenIndex289, depth289
								}
								if !matchDot() {
									goto l288
								}
								goto l287
							l288:
								position, tokenIndex, depth = position287, tokenIndex287, depth287
								if buffer[position] != rune('\\') {
									goto l285
								}
								position++
								{
									position292, tokenIndex292, depth292 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l293
									}
									position++
									goto l292
								l293:
									position, tokenIndex, depth = position292, tokenIndex292, depth292
									if buffer[position] != rune('\\') {
										goto l285
									}
									position++
								}
							l292:
							}
						l287:
							depth--
							add(rulestringChar, position286)
						}
						goto l284
					l285:
						position, tokenIndex, depth = position285, tokenIndex285, depth285
					}
					depth--
					add(rulePegText, position283)
				}
				if buffer[position] != rune('"') {
					goto l281
				}
				position++
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleString, position282)
			}
			return true
		l281:
			position, tokenIndex, depth = position281, tokenIndex281, depth281
			return false
		},
		/* 43 Null <- <('n' 'u' 'l' 'l' Action20)> */
		nil,
		/* 44 True <- <('t' 'r' 'u' 'e' Action21)> */
		nil,
		/* 45 False <- <('f' 'a' 'l' 's' 'e' Action22)> */
		nil,
		/* 46 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action23)> */
		nil,
		/* 47 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action24)> */
		nil,
		/* 48 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action25)> */
		nil,
		/* 49 Regex <- <('/' <regexBody> Action26)> */
		nil,
		/* 50 TimestampVal <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action27)> */
		nil,
		/* 51 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action28)> */
		nil,
		/* 52 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action29)> */
		nil,
		/* 53 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action30)> */
		nil,
		/* 54 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action31)> */
		nil,
		/* 55 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 56 regexChar <- <(!'/' .)> */
		nil,
		/* 57 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 58 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 59 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l311
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l311
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l311
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l311
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l311
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l311
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l311
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		nil,
		/* 62 Action0 <- <{ p.SetField("thread", buffer[begin:end]) }> */
		nil,
		/* 63 Action1 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 64 Action2 <- <{ p.EndField() }> */
		nil,
		/* 65 Action3 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 66 Action4 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 67 Action5 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 68 Action6 <- <{ p.SetField("commandType", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 69 Action7 <- <{ p.SetField("planSummaryType", buffer[begin:end]); p.StartField("planSummary") }> */
		nil,
		/* 70 Action8 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 71 Action9 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 72 Action10 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 73 Action11 <- <{ p.PushMap() }> */
		nil,
		/* 74 Action12 <- <{ p.PopMap() }> */
		nil,
		/* 75 Action13 <- <{ p.SetMapValue() }> */
		nil,
		/* 76 Action14 <- <{ p.PushList() }> */
		nil,
		/* 77 Action15 <- <{ p.PopList() }> */
		nil,
		/* 78 Action16 <- <{ p.SetListValue() }> */
		nil,
		/* 79 Action17 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 80 Action18 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 81 Action19 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 82 Action20 <- <{ p.PushValue(nil) }> */
		nil,
		/* 83 Action21 <- <{ p.PushValue(true) }> */
		nil,
		/* 84 Action22 <- <{ p.PushValue(false) }> */
		nil,
		/* 85 Action23 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 86 Action24 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 87 Action25 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 88 Action26 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 89 Action27 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 90 Action28 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 91 Action29 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 92 Action30 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 93 Action31 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
