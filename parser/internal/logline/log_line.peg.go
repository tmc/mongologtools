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
	ruleSubsystem
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
	ruleAction37
	ruleAction38

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"MongoLogLine",
	"Timestamp",
	"LogLevel",
	"Subsystem",
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
	"Action37",
	"Action38",

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
	rules  [107]func() bool
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
			p.SetField("subsystem", buffer[begin:end])
		case ruleAction2:
			p.SetField("thread", buffer[begin:end])
		case ruleAction3:
			p.SetField("op", buffer[begin:end])
		case ruleAction4:
			p.SetField("ns", buffer[begin:end])
		case ruleAction5:
			p.SetField("duration_ms", buffer[begin:end])
		case ruleAction6:
			p.StartField(buffer[begin:end])
		case ruleAction7:
			p.EndField()
		case ruleAction8:
			p.SetField("commandType", buffer[begin:end])
			p.StartField("command")
		case ruleAction9:
			p.EndField()
		case ruleAction10:
			p.StartField("planSummary")
			p.PushList()
		case ruleAction11:
			p.EndField()
		case ruleAction12:
			p.PushMap()
			p.PushField(buffer[begin:end])
		case ruleAction13:
			p.SetMapValue()
			p.SetListValue()
		case ruleAction14:
			p.PushValue(1)
			p.SetMapValue()
			p.SetListValue()
		case ruleAction15:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction16:
			p.SetField("timestamp", buffer[begin:end])
		case ruleAction17:
			p.SetField("xextra", buffer[begin:end])
		case ruleAction18:
			p.PushMap()
		case ruleAction19:
			p.PopMap()
		case ruleAction20:
			p.SetMapValue()
		case ruleAction21:
			p.PushList()
		case ruleAction22:
			p.PopList()
		case ruleAction23:
			p.SetListValue()
		case ruleAction24:
			p.PushField(buffer[begin:end])
		case ruleAction25:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction26:
			p.PushValue(buffer[begin:end])
		case ruleAction27:
			p.PushValue(nil)
		case ruleAction28:
			p.PushValue(true)
		case ruleAction29:
			p.PushValue(false)
		case ruleAction30:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction31:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction32:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction33:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction34:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction35:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction36:
			p.PushValue(p.Minkey())
		case ruleAction37:
			p.PushValue(p.Maxkey())
		case ruleAction38:
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
		/* 0 MongoLogLine <- <(Timestamp LogLevel? Subsystem? Thread Op NS LineField* Locks? LineField* Duration? extra? !.)> */
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
								add(ruleAction15, position)
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
								add(ruleAction16, position)
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
						add(ruleSubsystem, position35)
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
					add(ruleThread, position42)
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
											position79, tokenIndex79, depth79 := position, tokenIndex, depth
											if buffer[position] != rune('w') {
												goto l80
											}
											position++
											goto l79
										l80:
											position, tokenIndex, depth = position79, tokenIndex79, depth79
											if buffer[position] != rune('W') {
												goto l76
											}
											position++
										}
									l79:
										break
									}
								}

								if buffer[position] != rune(':') {
									goto l76
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l76
								}
								position++
							l81:
								{
									position82, tokenIndex82, depth82 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l82
									}
									position++
									goto l81
								l82:
									position, tokenIndex, depth = position82, tokenIndex82, depth82
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
			l85:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if !_rules[ruleLineField]() {
						goto l86
					}
					goto l85
				l86:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
				}
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					{
						position89 := position
						depth++
						{
							position90 := position
							depth++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l87
							}
							position++
						l91:
							{
								position92, tokenIndex92, depth92 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l92
								}
								position++
								goto l91
							l92:
								position, tokenIndex, depth = position92, tokenIndex92, depth92
							}
							depth--
							add(rulePegText, position90)
						}
						if buffer[position] != rune('m') {
							goto l87
						}
						position++
						if buffer[position] != rune('s') {
							goto l87
						}
						position++
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleDuration, position89)
					}
					goto l88
				l87:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
				}
			l88:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					{
						position96 := position
						depth++
						{
							position97 := position
							depth++
							if !matchDot() {
								goto l94
							}
						l98:
							{
								position99, tokenIndex99, depth99 := position, tokenIndex, depth
								if !matchDot() {
									goto l99
								}
								goto l98
							l99:
								position, tokenIndex, depth = position99, tokenIndex99, depth99
							}
							depth--
							add(rulePegText, position97)
						}
						{
							add(ruleAction17, position)
						}
						depth--
						add(ruleextra, position96)
					}
					goto l95
				l94:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
				}
			l95:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if !matchDot() {
						goto l101
					}
					goto l0
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
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
		/* 3 Subsystem <- <(<[A-Z]+> ' '+ Action1)> */
		nil,
		/* 4 Thread <- <('[' <letterOrDigit+> ']' ' ' Action2)> */
		nil,
		/* 5 Op <- <(<((&('c') ('c' 'o' 'm' 'm' 'a' 'n' 'd')) | (&('g') ('g' 'e' 't' 'm' 'o' 'r' 'e')) | (&('r') ('r' 'e' 'm' 'o' 'v' 'e')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('i') ('i' 'n' 's' 'e' 'r' 't')) | (&('q') ('q' 'u' 'e' 'r' 'y')))> ' ' Action3)> */
		nil,
		/* 6 LineField <- <((commandField / planSummaryField / plainField) S?)> */
		func() bool {
			position107, tokenIndex107, depth107 := position, tokenIndex, depth
			{
				position108 := position
				depth++
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					{
						position111 := position
						depth++
						if buffer[position] != rune('c') {
							goto l110
						}
						position++
						if buffer[position] != rune('o') {
							goto l110
						}
						position++
						if buffer[position] != rune('m') {
							goto l110
						}
						position++
						if buffer[position] != rune('m') {
							goto l110
						}
						position++
						if buffer[position] != rune('a') {
							goto l110
						}
						position++
						if buffer[position] != rune('n') {
							goto l110
						}
						position++
						if buffer[position] != rune('d') {
							goto l110
						}
						position++
						if buffer[position] != rune(':') {
							goto l110
						}
						position++
						if buffer[position] != rune(' ') {
							goto l110
						}
						position++
						{
							position112 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l110
							}
						l113:
							{
								position114, tokenIndex114, depth114 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l114
								}
								goto l113
							l114:
								position, tokenIndex, depth = position114, tokenIndex114, depth114
							}
							depth--
							add(rulePegText, position112)
						}
						{
							add(ruleAction8, position)
						}
						if !_rules[ruleLineValue]() {
							goto l110
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(rulecommandField, position111)
					}
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					{
						position118 := position
						depth++
						if buffer[position] != rune('p') {
							goto l117
						}
						position++
						if buffer[position] != rune('l') {
							goto l117
						}
						position++
						if buffer[position] != rune('a') {
							goto l117
						}
						position++
						if buffer[position] != rune('n') {
							goto l117
						}
						position++
						if buffer[position] != rune('S') {
							goto l117
						}
						position++
						if buffer[position] != rune('u') {
							goto l117
						}
						position++
						if buffer[position] != rune('m') {
							goto l117
						}
						position++
						if buffer[position] != rune('m') {
							goto l117
						}
						position++
						if buffer[position] != rune('a') {
							goto l117
						}
						position++
						if buffer[position] != rune('r') {
							goto l117
						}
						position++
						if buffer[position] != rune('y') {
							goto l117
						}
						position++
						if buffer[position] != rune(':') {
							goto l117
						}
						position++
						if buffer[position] != rune(' ') {
							goto l117
						}
						position++
						{
							add(ruleAction10, position)
						}
						{
							position120 := position
							depth++
							if !_rules[ruleplanSummaryElem]() {
								goto l117
							}
						l121:
							{
								position122, tokenIndex122, depth122 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l122
								}
								position++
								if buffer[position] != rune(' ') {
									goto l122
								}
								position++
								if !_rules[ruleplanSummaryElem]() {
									goto l122
								}
								goto l121
							l122:
								position, tokenIndex, depth = position122, tokenIndex122, depth122
							}
							depth--
							add(ruleplanSummaryElements, position120)
						}
						{
							add(ruleAction11, position)
						}
						depth--
						add(ruleplanSummaryField, position118)
					}
					goto l109
				l117:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					{
						position124 := position
						depth++
						{
							position125 := position
							depth++
							if !_rules[rulefieldChar]() {
								goto l107
							}
						l126:
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if !_rules[rulefieldChar]() {
									goto l127
								}
								goto l126
							l127:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
							}
							depth--
							add(rulePegText, position125)
						}
						if buffer[position] != rune(':') {
							goto l107
						}
						position++
						{
							position128, tokenIndex128, depth128 := position, tokenIndex, depth
							if !_rules[ruleS]() {
								goto l128
							}
							goto l129
						l128:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
						}
					l129:
						{
							add(ruleAction6, position)
						}
						if !_rules[ruleLineValue]() {
							goto l107
						}
						{
							add(ruleAction7, position)
						}
						depth--
						add(ruleplainField, position124)
					}
				}
			l109:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l132
					}
					goto l133
				l132:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
				}
			l133:
				depth--
				add(ruleLineField, position108)
			}
			return true
		l107:
			position, tokenIndex, depth = position107, tokenIndex107, depth107
			return false
		},
		/* 7 NS <- <(<nsChar+> ' ' Action4)> */
		nil,
		/* 8 Locks <- <('l' 'o' 'c' 'k' 's' '(' 'm' 'i' 'c' 'r' 'o' 's' ')' S? lock*)> */
		nil,
		/* 9 lock <- <(((&('R') 'R') | (&('r') 'r') | (&('W' | 'w') ('w' / 'W'))) ':' [0-9]+ S?)> */
		nil,
		/* 10 Duration <- <(<[0-9]+> ('m' 's') Action5)> */
		nil,
		/* 11 plainField <- <(<fieldChar+> ':' S? Action6 LineValue Action7)> */
		nil,
		/* 12 commandField <- <('c' 'o' 'm' 'm' 'a' 'n' 'd' ':' ' ' <fieldChar+> Action8 LineValue Action9)> */
		nil,
		/* 13 planSummaryField <- <('p' 'l' 'a' 'n' 'S' 'u' 'm' 'm' 'a' 'r' 'y' ':' ' ' Action10 planSummaryElements Action11)> */
		nil,
		/* 14 planSummaryElements <- <(planSummaryElem (',' ' ' planSummaryElem)*)> */
		nil,
		/* 15 planSummaryElem <- <(<planSummaryStage> Action12 planSummary)> */
		func() bool {
			position142, tokenIndex142, depth142 := position, tokenIndex, depth
			{
				position143 := position
				depth++
				{
					position144 := position
					depth++
					{
						position145 := position
						depth++
						{
							position146, tokenIndex146, depth146 := position, tokenIndex, depth
							if buffer[position] != rune('A') {
								goto l147
							}
							position++
							if buffer[position] != rune('N') {
								goto l147
							}
							position++
							if buffer[position] != rune('D') {
								goto l147
							}
							position++
							if buffer[position] != rune('_') {
								goto l147
							}
							position++
							if buffer[position] != rune('H') {
								goto l147
							}
							position++
							if buffer[position] != rune('A') {
								goto l147
							}
							position++
							if buffer[position] != rune('S') {
								goto l147
							}
							position++
							if buffer[position] != rune('H') {
								goto l147
							}
							position++
							goto l146
						l147:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('C') {
								goto l148
							}
							position++
							if buffer[position] != rune('A') {
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
							goto l146
						l148:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('C') {
								goto l149
							}
							position++
							if buffer[position] != rune('O') {
								goto l149
							}
							position++
							if buffer[position] != rune('L') {
								goto l149
							}
							position++
							if buffer[position] != rune('L') {
								goto l149
							}
							position++
							if buffer[position] != rune('S') {
								goto l149
							}
							position++
							if buffer[position] != rune('C') {
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
							goto l146
						l149:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('C') {
								goto l150
							}
							position++
							if buffer[position] != rune('O') {
								goto l150
							}
							position++
							if buffer[position] != rune('U') {
								goto l150
							}
							position++
							if buffer[position] != rune('N') {
								goto l150
							}
							position++
							if buffer[position] != rune('T') {
								goto l150
							}
							position++
							goto l146
						l150:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('D') {
								goto l151
							}
							position++
							if buffer[position] != rune('E') {
								goto l151
							}
							position++
							if buffer[position] != rune('L') {
								goto l151
							}
							position++
							if buffer[position] != rune('E') {
								goto l151
							}
							position++
							if buffer[position] != rune('T') {
								goto l151
							}
							position++
							if buffer[position] != rune('E') {
								goto l151
							}
							position++
							goto l146
						l151:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('G') {
								goto l152
							}
							position++
							if buffer[position] != rune('E') {
								goto l152
							}
							position++
							if buffer[position] != rune('O') {
								goto l152
							}
							position++
							if buffer[position] != rune('_') {
								goto l152
							}
							position++
							if buffer[position] != rune('N') {
								goto l152
							}
							position++
							if buffer[position] != rune('E') {
								goto l152
							}
							position++
							if buffer[position] != rune('A') {
								goto l152
							}
							position++
							if buffer[position] != rune('R') {
								goto l152
							}
							position++
							if buffer[position] != rune('_') {
								goto l152
							}
							position++
							if buffer[position] != rune('2') {
								goto l152
							}
							position++
							if buffer[position] != rune('D') {
								goto l152
							}
							position++
							goto l146
						l152:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
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
							if buffer[position] != rune('S') {
								goto l153
							}
							position++
							if buffer[position] != rune('P') {
								goto l153
							}
							position++
							if buffer[position] != rune('H') {
								goto l153
							}
							position++
							if buffer[position] != rune('E') {
								goto l153
							}
							position++
							if buffer[position] != rune('R') {
								goto l153
							}
							position++
							if buffer[position] != rune('E') {
								goto l153
							}
							position++
							goto l146
						l153:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('I') {
								goto l154
							}
							position++
							if buffer[position] != rune('D') {
								goto l154
							}
							position++
							if buffer[position] != rune('H') {
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
							if buffer[position] != rune('K') {
								goto l154
							}
							position++
							goto l146
						l154:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('S') {
								goto l155
							}
							position++
							if buffer[position] != rune('O') {
								goto l155
							}
							position++
							if buffer[position] != rune('R') {
								goto l155
							}
							position++
							if buffer[position] != rune('T') {
								goto l155
							}
							position++
							if buffer[position] != rune('_') {
								goto l155
							}
							position++
							if buffer[position] != rune('M') {
								goto l155
							}
							position++
							if buffer[position] != rune('E') {
								goto l155
							}
							position++
							if buffer[position] != rune('R') {
								goto l155
							}
							position++
							if buffer[position] != rune('G') {
								goto l155
							}
							position++
							if buffer[position] != rune('E') {
								goto l155
							}
							position++
							goto l146
						l155:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('S') {
								goto l156
							}
							position++
							if buffer[position] != rune('H') {
								goto l156
							}
							position++
							if buffer[position] != rune('A') {
								goto l156
							}
							position++
							if buffer[position] != rune('R') {
								goto l156
							}
							position++
							if buffer[position] != rune('D') {
								goto l156
							}
							position++
							if buffer[position] != rune('I') {
								goto l156
							}
							position++
							if buffer[position] != rune('N') {
								goto l156
							}
							position++
							if buffer[position] != rune('G') {
								goto l156
							}
							position++
							if buffer[position] != rune('_') {
								goto l156
							}
							position++
							if buffer[position] != rune('F') {
								goto l156
							}
							position++
							if buffer[position] != rune('I') {
								goto l156
							}
							position++
							if buffer[position] != rune('L') {
								goto l156
							}
							position++
							if buffer[position] != rune('T') {
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
							goto l146
						l156:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('S') {
								goto l157
							}
							position++
							if buffer[position] != rune('K') {
								goto l157
							}
							position++
							if buffer[position] != rune('I') {
								goto l157
							}
							position++
							if buffer[position] != rune('P') {
								goto l157
							}
							position++
							goto l146
						l157:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							if buffer[position] != rune('S') {
								goto l158
							}
							position++
							if buffer[position] != rune('O') {
								goto l158
							}
							position++
							if buffer[position] != rune('R') {
								goto l158
							}
							position++
							if buffer[position] != rune('T') {
								goto l158
							}
							position++
							goto l146
						l158:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
							{
								switch buffer[position] {
								case 'U':
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									break
								case 'T':
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('X') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									break
								case 'S':
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('B') {
										goto l142
									}
									position++
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									if buffer[position] != rune('L') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									break
								case 'Q':
									if buffer[position] != rune('Q') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									if buffer[position] != rune('_') {
										goto l142
									}
									position++
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									break
								case 'P':
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									if buffer[position] != rune('R') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('J') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									break
								case 'O':
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('R') {
										goto l142
									}
									position++
									break
								case 'M':
									if buffer[position] != rune('M') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('L') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('_') {
										goto l142
									}
									position++
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									if buffer[position] != rune('L') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									break
								case 'L':
									if buffer[position] != rune('L') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('M') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									break
								case 'K':
									if buffer[position] != rune('K') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									if buffer[position] != rune('_') {
										goto l142
									}
									position++
									if buffer[position] != rune('M') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									break
								case 'I':
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('X') {
										goto l142
									}
									position++
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									break
								case 'G':
									if buffer[position] != rune('G') {
										goto l142
									}
									position++
									if buffer[position] != rune('R') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('P') {
										goto l142
									}
									position++
									break
								case 'F':
									if buffer[position] != rune('F') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('H') {
										goto l142
									}
									position++
									break
								case 'E':
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('F') {
										goto l142
									}
									position++
									break
								case 'D':
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('I') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									break
								case 'C':
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('U') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('_') {
										goto l142
									}
									position++
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									if buffer[position] != rune('C') {
										goto l142
									}
									position++
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									break
								default:
									if buffer[position] != rune('A') {
										goto l142
									}
									position++
									if buffer[position] != rune('N') {
										goto l142
									}
									position++
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									if buffer[position] != rune('_') {
										goto l142
									}
									position++
									if buffer[position] != rune('S') {
										goto l142
									}
									position++
									if buffer[position] != rune('O') {
										goto l142
									}
									position++
									if buffer[position] != rune('R') {
										goto l142
									}
									position++
									if buffer[position] != rune('T') {
										goto l142
									}
									position++
									if buffer[position] != rune('E') {
										goto l142
									}
									position++
									if buffer[position] != rune('D') {
										goto l142
									}
									position++
									break
								}
							}

						}
					l146:
						depth--
						add(ruleplanSummaryStage, position145)
					}
					depth--
					add(rulePegText, position144)
				}
				{
					add(ruleAction12, position)
				}
				{
					position161 := position
					depth++
					{
						position162, tokenIndex162, depth162 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l163
						}
						position++
						if !_rules[ruleLineValue]() {
							goto l163
						}
						{
							add(ruleAction13, position)
						}
						goto l162
					l163:
						position, tokenIndex, depth = position162, tokenIndex162, depth162
						{
							add(ruleAction14, position)
						}
					}
				l162:
					depth--
					add(ruleplanSummary, position161)
				}
				depth--
				add(ruleplanSummaryElem, position143)
			}
			return true
		l142:
			position, tokenIndex, depth = position142, tokenIndex142, depth142
			return false
		},
		/* 16 planSummaryStage <- <(('A' 'N' 'D' '_' 'H' 'A' 'S' 'H') / ('C' 'A' 'C' 'H' 'E' 'D' '_' 'P' 'L' 'A' 'N') / ('C' 'O' 'L' 'L' 'S' 'C' 'A' 'N') / ('C' 'O' 'U' 'N' 'T') / ('D' 'E' 'L' 'E' 'T' 'E') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D') / ('G' 'E' 'O' '_' 'N' 'E' 'A' 'R' '_' '2' 'D' 'S' 'P' 'H' 'E' 'R' 'E') / ('I' 'D' 'H' 'A' 'C' 'K') / ('S' 'O' 'R' 'T' '_' 'M' 'E' 'R' 'G' 'E') / ('S' 'H' 'A' 'R' 'D' 'I' 'N' 'G' '_' 'F' 'I' 'L' 'T' 'E' 'R') / ('S' 'K' 'I' 'P') / ('S' 'O' 'R' 'T') / ((&('U') ('U' 'P' 'D' 'A' 'T' 'E')) | (&('T') ('T' 'E' 'X' 'T')) | (&('S') ('S' 'U' 'B' 'P' 'L' 'A' 'N')) | (&('Q') ('Q' 'U' 'E' 'U' 'E' 'D' '_' 'D' 'A' 'T' 'A')) | (&('P') ('P' 'R' 'O' 'J' 'E' 'C' 'T' 'I' 'O' 'N')) | (&('O') ('O' 'R')) | (&('M') ('M' 'U' 'L' 'T' 'I' '_' 'P' 'L' 'A' 'N')) | (&('L') ('L' 'I' 'M' 'I' 'T')) | (&('K') ('K' 'E' 'E' 'P' '_' 'M' 'U' 'T' 'A' 'T' 'I' 'O' 'N' 'S')) | (&('I') ('I' 'X' 'S' 'C' 'A' 'N')) | (&('G') ('G' 'R' 'O' 'U' 'P')) | (&('F') ('F' 'E' 'T' 'C' 'H')) | (&('E') ('E' 'O' 'F')) | (&('D') ('D' 'I' 'S' 'T' 'I' 'N' 'C' 'T')) | (&('C') ('C' 'O' 'U' 'N' 'T' '_' 'S' 'C' 'A' 'N')) | (&('A') ('A' 'N' 'D' '_' 'S' 'O' 'R' 'T' 'E' 'D'))))> */
		nil,
		/* 17 planSummary <- <((' ' LineValue Action13) / Action14)> */
		nil,
		/* 18 LineValue <- <((Doc / Numeric) S?)> */
		func() bool {
			position168, tokenIndex168, depth168 := position, tokenIndex, depth
			{
				position169 := position
				depth++
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if !_rules[ruleDoc]() {
						goto l171
					}
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if !_rules[ruleNumeric]() {
						goto l168
					}
				}
			l170:
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
				depth--
				add(ruleLineValue, position169)
			}
			return true
		l168:
			position, tokenIndex, depth = position168, tokenIndex168, depth168
			return false
		},
		/* 19 timestamp24 <- <(<(date ' ' time)> Action15)> */
		nil,
		/* 20 timestamp26 <- <(<datetime26> Action16)> */
		nil,
		/* 21 datetime26 <- <(digit4 '-' digit2 '-' digit2 'T' time tz?)> */
		nil,
		/* 22 digit4 <- <([0-9] [0-9] [0-9] [0-9])> */
		nil,
		/* 23 digit2 <- <([0-9] [0-9])> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l178
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l178
				}
				position++
				depth--
				add(ruledigit2, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 24 date <- <(day ' ' month ' ' dayNum)> */
		nil,
		/* 25 tz <- <('+' [0-9]+)> */
		nil,
		/* 26 time <- <(hour ':' minute ':' second '.' millisecond)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				{
					position184 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l182
					}
					depth--
					add(rulehour, position184)
				}
				if buffer[position] != rune(':') {
					goto l182
				}
				position++
				{
					position185 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l182
					}
					depth--
					add(ruleminute, position185)
				}
				if buffer[position] != rune(':') {
					goto l182
				}
				position++
				{
					position186 := position
					depth++
					if !_rules[ruledigit2]() {
						goto l182
					}
					depth--
					add(rulesecond, position186)
				}
				if buffer[position] != rune('.') {
					goto l182
				}
				position++
				{
					position187 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l182
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l182
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l182
					}
					position++
					depth--
					add(rulemillisecond, position187)
				}
				depth--
				add(ruletime, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
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
		/* 36 extra <- <(<.+> Action17)> */
		nil,
		/* 37 S <- <' '+> */
		func() bool {
			position198, tokenIndex198, depth198 := position, tokenIndex, depth
			{
				position199 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l198
				}
				position++
			l200:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
				}
				depth--
				add(ruleS, position199)
			}
			return true
		l198:
			position, tokenIndex, depth = position198, tokenIndex198, depth198
			return false
		},
		/* 38 Doc <- <('{' Action18 DocElements? '}' Action19)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				if buffer[position] != rune('{') {
					goto l202
				}
				position++
				{
					add(ruleAction18, position)
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					{
						position207 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l205
						}
					l208:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l209
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l209
							}
							goto l208
						l209:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
						}
						depth--
						add(ruleDocElements, position207)
					}
					goto l206
				l205:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
				}
			l206:
				if buffer[position] != rune('}') {
					goto l202
				}
				position++
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleDoc, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 39 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 40 DocElem <- <(S? Field S? Value S? Action20)> */
		func() bool {
			position212, tokenIndex212, depth212 := position, tokenIndex, depth
			{
				position213 := position
				depth++
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
					position216 := position
					depth++
					{
						position217 := position
						depth++
						if !_rules[rulefieldChar]() {
							goto l212
						}
					l218:
						{
							position219, tokenIndex219, depth219 := position, tokenIndex, depth
							if !_rules[rulefieldChar]() {
								goto l219
							}
							goto l218
						l219:
							position, tokenIndex, depth = position219, tokenIndex219, depth219
						}
						depth--
						add(rulePegText, position217)
					}
					if buffer[position] != rune(':') {
						goto l212
					}
					position++
					{
						add(ruleAction24, position)
					}
					depth--
					add(ruleField, position216)
				}
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l221
					}
					goto l222
				l221:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
				}
			l222:
				if !_rules[ruleValue]() {
					goto l212
				}
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l223
					}
					goto l224
				l223:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
				}
			l224:
				{
					add(ruleAction20, position)
				}
				depth--
				add(ruleDocElem, position213)
			}
			return true
		l212:
			position, tokenIndex, depth = position212, tokenIndex212, depth212
			return false
		},
		/* 41 List <- <('[' Action21 ListElements? ']' Action22)> */
		nil,
		/* 42 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 43 ListElem <- <(S? Value S? Action23)> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l230
					}
					goto l231
				l230:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
				}
			l231:
				if !_rules[ruleValue]() {
					goto l228
				}
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l232
					}
					goto l233
				l232:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
				}
			l233:
				{
					add(ruleAction23, position)
				}
				depth--
				add(ruleListElem, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 44 Field <- <(<fieldChar+> ':' Action24)> */
		nil,
		/* 45 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position236, tokenIndex236, depth236 := position, tokenIndex, depth
			{
				position237 := position
				depth++
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					{
						position240 := position
						depth++
						if buffer[position] != rune('n') {
							goto l239
						}
						position++
						if buffer[position] != rune('u') {
							goto l239
						}
						position++
						if buffer[position] != rune('l') {
							goto l239
						}
						position++
						if buffer[position] != rune('l') {
							goto l239
						}
						position++
						{
							add(ruleAction27, position)
						}
						depth--
						add(ruleNull, position240)
					}
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					{
						position243 := position
						depth++
						if buffer[position] != rune('M') {
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
							add(ruleAction36, position)
						}
						depth--
						add(ruleMinKey, position243)
					}
					goto l238
				l242:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					{
						switch buffer[position] {
						case 'M':
							{
								position246 := position
								depth++
								if buffer[position] != rune('M') {
									goto l236
								}
								position++
								if buffer[position] != rune('a') {
									goto l236
								}
								position++
								if buffer[position] != rune('x') {
									goto l236
								}
								position++
								if buffer[position] != rune('K') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('y') {
									goto l236
								}
								position++
								{
									add(ruleAction37, position)
								}
								depth--
								add(ruleMaxKey, position246)
							}
							break
						case 'u':
							{
								position248 := position
								depth++
								if buffer[position] != rune('u') {
									goto l236
								}
								position++
								if buffer[position] != rune('n') {
									goto l236
								}
								position++
								if buffer[position] != rune('d') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('f') {
									goto l236
								}
								position++
								if buffer[position] != rune('i') {
									goto l236
								}
								position++
								if buffer[position] != rune('n') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('d') {
									goto l236
								}
								position++
								{
									add(ruleAction38, position)
								}
								depth--
								add(ruleUndefined, position248)
							}
							break
						case 'N':
							{
								position250 := position
								depth++
								if buffer[position] != rune('N') {
									goto l236
								}
								position++
								if buffer[position] != rune('u') {
									goto l236
								}
								position++
								if buffer[position] != rune('m') {
									goto l236
								}
								position++
								if buffer[position] != rune('b') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('r') {
									goto l236
								}
								position++
								if buffer[position] != rune('L') {
									goto l236
								}
								position++
								if buffer[position] != rune('o') {
									goto l236
								}
								position++
								if buffer[position] != rune('n') {
									goto l236
								}
								position++
								if buffer[position] != rune('g') {
									goto l236
								}
								position++
								if buffer[position] != rune('(') {
									goto l236
								}
								position++
								{
									position251 := position
									depth++
									{
										position254, tokenIndex254, depth254 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l254
										}
										position++
										goto l236
									l254:
										position, tokenIndex, depth = position254, tokenIndex254, depth254
									}
									if !matchDot() {
										goto l236
									}
								l252:
									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										{
											position255, tokenIndex255, depth255 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l255
											}
											position++
											goto l253
										l255:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
										}
										if !matchDot() {
											goto l253
										}
										goto l252
									l253:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
									}
									depth--
									add(rulePegText, position251)
								}
								if buffer[position] != rune(')') {
									goto l236
								}
								position++
								{
									add(ruleAction35, position)
								}
								depth--
								add(ruleNumberLong, position250)
							}
							break
						case '/':
							{
								position257 := position
								depth++
								if buffer[position] != rune('/') {
									goto l236
								}
								position++
								{
									position258 := position
									depth++
									{
										position259 := position
										depth++
										{
											position262 := position
											depth++
											{
												position263, tokenIndex263, depth263 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l263
												}
												position++
												goto l236
											l263:
												position, tokenIndex, depth = position263, tokenIndex263, depth263
											}
											if !matchDot() {
												goto l236
											}
											depth--
											add(ruleregexChar, position262)
										}
									l260:
										{
											position261, tokenIndex261, depth261 := position, tokenIndex, depth
											{
												position264 := position
												depth++
												{
													position265, tokenIndex265, depth265 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l265
													}
													position++
													goto l261
												l265:
													position, tokenIndex, depth = position265, tokenIndex265, depth265
												}
												if !matchDot() {
													goto l261
												}
												depth--
												add(ruleregexChar, position264)
											}
											goto l260
										l261:
											position, tokenIndex, depth = position261, tokenIndex261, depth261
										}
										if buffer[position] != rune('/') {
											goto l236
										}
										position++
									l266:
										{
											position267, tokenIndex267, depth267 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l267
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l267
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l267
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l267
													}
													position++
													break
												}
											}

											goto l266
										l267:
											position, tokenIndex, depth = position267, tokenIndex267, depth267
										}
										depth--
										add(ruleregexBody, position259)
									}
									depth--
									add(rulePegText, position258)
								}
								{
									add(ruleAction33, position)
								}
								depth--
								add(ruleRegex, position257)
							}
							break
						case 'T':
							{
								position270 := position
								depth++
								if buffer[position] != rune('T') {
									goto l236
								}
								position++
								if buffer[position] != rune('i') {
									goto l236
								}
								position++
								if buffer[position] != rune('m') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('s') {
									goto l236
								}
								position++
								if buffer[position] != rune('t') {
									goto l236
								}
								position++
								if buffer[position] != rune('a') {
									goto l236
								}
								position++
								if buffer[position] != rune('m') {
									goto l236
								}
								position++
								if buffer[position] != rune('p') {
									goto l236
								}
								position++
								if buffer[position] != rune('(') {
									goto l236
								}
								position++
								{
									position271 := position
									depth++
									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l274
										}
										position++
										goto l236
									l274:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
									}
									if !matchDot() {
										goto l236
									}
								l272:
									{
										position273, tokenIndex273, depth273 := position, tokenIndex, depth
										{
											position275, tokenIndex275, depth275 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l275
											}
											position++
											goto l273
										l275:
											position, tokenIndex, depth = position275, tokenIndex275, depth275
										}
										if !matchDot() {
											goto l273
										}
										goto l272
									l273:
										position, tokenIndex, depth = position273, tokenIndex273, depth273
									}
									depth--
									add(rulePegText, position271)
								}
								if buffer[position] != rune(')') {
									goto l236
								}
								position++
								{
									add(ruleAction34, position)
								}
								depth--
								add(ruleTimestampVal, position270)
							}
							break
						case 'B':
							{
								position277 := position
								depth++
								if buffer[position] != rune('B') {
									goto l236
								}
								position++
								if buffer[position] != rune('i') {
									goto l236
								}
								position++
								if buffer[position] != rune('n') {
									goto l236
								}
								position++
								if buffer[position] != rune('D') {
									goto l236
								}
								position++
								if buffer[position] != rune('a') {
									goto l236
								}
								position++
								if buffer[position] != rune('t') {
									goto l236
								}
								position++
								if buffer[position] != rune('a') {
									goto l236
								}
								position++
								if buffer[position] != rune('(') {
									goto l236
								}
								position++
								{
									position278 := position
									depth++
									{
										position281, tokenIndex281, depth281 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l281
										}
										position++
										goto l236
									l281:
										position, tokenIndex, depth = position281, tokenIndex281, depth281
									}
									if !matchDot() {
										goto l236
									}
								l279:
									{
										position280, tokenIndex280, depth280 := position, tokenIndex, depth
										{
											position282, tokenIndex282, depth282 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l282
											}
											position++
											goto l280
										l282:
											position, tokenIndex, depth = position282, tokenIndex282, depth282
										}
										if !matchDot() {
											goto l280
										}
										goto l279
									l280:
										position, tokenIndex, depth = position280, tokenIndex280, depth280
									}
									depth--
									add(rulePegText, position278)
								}
								if buffer[position] != rune(')') {
									goto l236
								}
								position++
								{
									add(ruleAction32, position)
								}
								depth--
								add(ruleBinData, position277)
							}
							break
						case 'D', 'n':
							{
								position284 := position
								depth++
								{
									position285, tokenIndex285, depth285 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l285
									}
									position++
									if buffer[position] != rune('e') {
										goto l285
									}
									position++
									if buffer[position] != rune('w') {
										goto l285
									}
									position++
									if buffer[position] != rune(' ') {
										goto l285
									}
									position++
									goto l286
								l285:
									position, tokenIndex, depth = position285, tokenIndex285, depth285
								}
							l286:
								if buffer[position] != rune('D') {
									goto l236
								}
								position++
								if buffer[position] != rune('a') {
									goto l236
								}
								position++
								if buffer[position] != rune('t') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('(') {
									goto l236
								}
								position++
								{
									position287 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l236
									}
									position++
								l288:
									{
										position289, tokenIndex289, depth289 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l289
										}
										position++
										goto l288
									l289:
										position, tokenIndex, depth = position289, tokenIndex289, depth289
									}
									depth--
									add(rulePegText, position287)
								}
								if buffer[position] != rune(')') {
									goto l236
								}
								position++
								{
									add(ruleAction30, position)
								}
								depth--
								add(ruleDate, position284)
							}
							break
						case 'O':
							{
								position291 := position
								depth++
								if buffer[position] != rune('O') {
									goto l236
								}
								position++
								if buffer[position] != rune('b') {
									goto l236
								}
								position++
								if buffer[position] != rune('j') {
									goto l236
								}
								position++
								if buffer[position] != rune('e') {
									goto l236
								}
								position++
								if buffer[position] != rune('c') {
									goto l236
								}
								position++
								if buffer[position] != rune('t') {
									goto l236
								}
								position++
								if buffer[position] != rune('I') {
									goto l236
								}
								position++
								if buffer[position] != rune('d') {
									goto l236
								}
								position++
								if buffer[position] != rune('(') {
									goto l236
								}
								position++
								if buffer[position] != rune('"') {
									goto l236
								}
								position++
								{
									position292 := position
									depth++
								l293:
									{
										position294, tokenIndex294, depth294 := position, tokenIndex, depth
										{
											position295 := position
											depth++
											{
												position296, tokenIndex296, depth296 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l297
												}
												position++
												goto l296
											l297:
												position, tokenIndex, depth = position296, tokenIndex296, depth296
												{
													position298, tokenIndex298, depth298 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l299
													}
													position++
													goto l298
												l299:
													position, tokenIndex, depth = position298, tokenIndex298, depth298
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l294
													}
													position++
												}
											l298:
											}
										l296:
											depth--
											add(rulehexChar, position295)
										}
										goto l293
									l294:
										position, tokenIndex, depth = position294, tokenIndex294, depth294
									}
									depth--
									add(rulePegText, position292)
								}
								if buffer[position] != rune('"') {
									goto l236
								}
								position++
								if buffer[position] != rune(')') {
									goto l236
								}
								position++
								{
									add(ruleAction31, position)
								}
								depth--
								add(ruleObjectID, position291)
							}
							break
						case '"':
							{
								position301 := position
								depth++
								if buffer[position] != rune('"') {
									goto l236
								}
								position++
								{
									position302 := position
									depth++
								l303:
									{
										position304, tokenIndex304, depth304 := position, tokenIndex, depth
										{
											position305 := position
											depth++
											{
												position306, tokenIndex306, depth306 := position, tokenIndex, depth
												{
													position308, tokenIndex308, depth308 := position, tokenIndex, depth
													{
														position309, tokenIndex309, depth309 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l310
														}
														position++
														goto l309
													l310:
														position, tokenIndex, depth = position309, tokenIndex309, depth309
														if buffer[position] != rune('\\') {
															goto l308
														}
														position++
													}
												l309:
													goto l307
												l308:
													position, tokenIndex, depth = position308, tokenIndex308, depth308
												}
												if !matchDot() {
													goto l307
												}
												goto l306
											l307:
												position, tokenIndex, depth = position306, tokenIndex306, depth306
												if buffer[position] != rune('\\') {
													goto l304
												}
												position++
												{
													position311, tokenIndex311, depth311 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l312
													}
													position++
													goto l311
												l312:
													position, tokenIndex, depth = position311, tokenIndex311, depth311
													if buffer[position] != rune('\\') {
														goto l304
													}
													position++
												}
											l311:
											}
										l306:
											depth--
											add(rulestringChar, position305)
										}
										goto l303
									l304:
										position, tokenIndex, depth = position304, tokenIndex304, depth304
									}
									depth--
									add(rulePegText, position302)
								}
								if buffer[position] != rune('"') {
									goto l236
								}
								position++
								{
									add(ruleAction26, position)
								}
								depth--
								add(ruleString, position301)
							}
							break
						case 'f', 't':
							{
								position314 := position
								depth++
								{
									position315, tokenIndex315, depth315 := position, tokenIndex, depth
									{
										position317 := position
										depth++
										if buffer[position] != rune('t') {
											goto l316
										}
										position++
										if buffer[position] != rune('r') {
											goto l316
										}
										position++
										if buffer[position] != rune('u') {
											goto l316
										}
										position++
										if buffer[position] != rune('e') {
											goto l316
										}
										position++
										{
											add(ruleAction28, position)
										}
										depth--
										add(ruleTrue, position317)
									}
									goto l315
								l316:
									position, tokenIndex, depth = position315, tokenIndex315, depth315
									{
										position319 := position
										depth++
										if buffer[position] != rune('f') {
											goto l236
										}
										position++
										if buffer[position] != rune('a') {
											goto l236
										}
										position++
										if buffer[position] != rune('l') {
											goto l236
										}
										position++
										if buffer[position] != rune('s') {
											goto l236
										}
										position++
										if buffer[position] != rune('e') {
											goto l236
										}
										position++
										{
											add(ruleAction29, position)
										}
										depth--
										add(ruleFalse, position319)
									}
								}
							l315:
								depth--
								add(ruleBoolean, position314)
							}
							break
						case '[':
							{
								position321 := position
								depth++
								if buffer[position] != rune('[') {
									goto l236
								}
								position++
								{
									add(ruleAction21, position)
								}
								{
									position323, tokenIndex323, depth323 := position, tokenIndex, depth
									{
										position325 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l323
										}
									l326:
										{
											position327, tokenIndex327, depth327 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l327
											}
											position++
											if !_rules[ruleListElem]() {
												goto l327
											}
											goto l326
										l327:
											position, tokenIndex, depth = position327, tokenIndex327, depth327
										}
										depth--
										add(ruleListElements, position325)
									}
									goto l324
								l323:
									position, tokenIndex, depth = position323, tokenIndex323, depth323
								}
							l324:
								if buffer[position] != rune(']') {
									goto l236
								}
								position++
								{
									add(ruleAction22, position)
								}
								depth--
								add(ruleList, position321)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l236
							}
							break
						default:
							if !_rules[ruleNumeric]() {
								goto l236
							}
							break
						}
					}

				}
			l238:
				depth--
				add(ruleValue, position237)
			}
			return true
		l236:
			position, tokenIndex, depth = position236, tokenIndex236, depth236
			return false
		},
		/* 46 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action25)> */
		func() bool {
			position329, tokenIndex329, depth329 := position, tokenIndex, depth
			{
				position330 := position
				depth++
				{
					position331 := position
					depth++
					{
						position332, tokenIndex332, depth332 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l332
						}
						position++
						goto l333
					l332:
						position, tokenIndex, depth = position332, tokenIndex332, depth332
					}
				l333:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l329
					}
					position++
				l334:
					{
						position335, tokenIndex335, depth335 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l335
						}
						position++
						goto l334
					l335:
						position, tokenIndex, depth = position335, tokenIndex335, depth335
					}
					{
						position336, tokenIndex336, depth336 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l336
						}
						position++
						goto l337
					l336:
						position, tokenIndex, depth = position336, tokenIndex336, depth336
					}
				l337:
				l338:
					{
						position339, tokenIndex339, depth339 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l339
						}
						position++
						goto l338
					l339:
						position, tokenIndex, depth = position339, tokenIndex339, depth339
					}
					depth--
					add(rulePegText, position331)
				}
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleNumeric, position330)
			}
			return true
		l329:
			position, tokenIndex, depth = position329, tokenIndex329, depth329
			return false
		},
		/* 47 Boolean <- <(True / False)> */
		nil,
		/* 48 String <- <('"' <stringChar*> '"' Action26)> */
		nil,
		/* 49 Null <- <('n' 'u' 'l' 'l' Action27)> */
		nil,
		/* 50 True <- <('t' 'r' 'u' 'e' Action28)> */
		nil,
		/* 51 False <- <('f' 'a' 'l' 's' 'e' Action29)> */
		nil,
		/* 52 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action30)> */
		nil,
		/* 53 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action31)> */
		nil,
		/* 54 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action32)> */
		nil,
		/* 55 Regex <- <('/' <regexBody> Action33)> */
		nil,
		/* 56 TimestampVal <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action34)> */
		nil,
		/* 57 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action35)> */
		nil,
		/* 58 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action36)> */
		nil,
		/* 59 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action37)> */
		nil,
		/* 60 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action38)> */
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
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					switch buffer[position] {
					case '$', '*', '.', '_':
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l359
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l359
								}
								position++
								break
							case '$':
								if buffer[position] != rune('$') {
									goto l359
								}
								position++
								break
							default:
								if buffer[position] != rune('_') {
									goto l359
								}
								position++
								break
							}
						}

						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l359
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l359
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l359
						}
						position++
						break
					}
				}

				depth--
				add(rulefieldChar, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		nil,
		/* 68 Action0 <- <{ p.SetField("log_level", buffer[begin:end]) }> */
		nil,
		/* 69 Action1 <- <{ p.SetField("subsystem", buffer[begin:end]) }> */
		nil,
		/* 70 Action2 <- <{ p.SetField("thread", buffer[begin:end]) }> */
		nil,
		/* 71 Action3 <- <{ p.SetField("op", buffer[begin:end]) }> */
		nil,
		/* 72 Action4 <- <{ p.SetField("ns", buffer[begin:end]) }> */
		nil,
		/* 73 Action5 <- <{ p.SetField("duration_ms", buffer[begin:end]) }> */
		nil,
		/* 74 Action6 <- <{ p.StartField(buffer[begin:end]) }> */
		nil,
		/* 75 Action7 <- <{ p.EndField() }> */
		nil,
		/* 76 Action8 <- <{ p.SetField("commandType", buffer[begin:end]); p.StartField("command") }> */
		nil,
		/* 77 Action9 <- <{ p.EndField() }> */
		nil,
		/* 78 Action10 <- <{ p.StartField("planSummary"); p.PushList() }> */
		nil,
		/* 79 Action11 <- <{ p.EndField()}> */
		nil,
		/* 80 Action12 <- <{ p.PushMap(); p.PushField(buffer[begin:end]) }> */
		nil,
		/* 81 Action13 <- <{ p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 82 Action14 <- <{ p.PushValue(1); p.SetMapValue(); p.SetListValue() }> */
		nil,
		/* 83 Action15 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 84 Action16 <- <{ p.SetField("timestamp", buffer[begin:end]) }> */
		nil,
		/* 85 Action17 <- <{ p.SetField("xextra", buffer[begin:end]) }> */
		nil,
		/* 86 Action18 <- <{ p.PushMap() }> */
		nil,
		/* 87 Action19 <- <{ p.PopMap() }> */
		nil,
		/* 88 Action20 <- <{ p.SetMapValue() }> */
		nil,
		/* 89 Action21 <- <{ p.PushList() }> */
		nil,
		/* 90 Action22 <- <{ p.PopList() }> */
		nil,
		/* 91 Action23 <- <{ p.SetListValue() }> */
		nil,
		/* 92 Action24 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 93 Action25 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 94 Action26 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 95 Action27 <- <{ p.PushValue(nil) }> */
		nil,
		/* 96 Action28 <- <{ p.PushValue(true) }> */
		nil,
		/* 97 Action29 <- <{ p.PushValue(false) }> */
		nil,
		/* 98 Action30 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 99 Action31 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 100 Action32 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 101 Action33 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 102 Action34 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 103 Action35 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 104 Action36 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 105 Action37 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 106 Action38 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
