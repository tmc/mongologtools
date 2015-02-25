package logdoc

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
	ruleLogDoc
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
	ruleS
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	rulePegText
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

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"LogDoc",
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
	"S",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"PegText",
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

type LogDocParser struct {
	LogDoc

	Buffer string
	buffer []rune
	rules  [56]func() bool
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
	p *LogDocParser
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

func (p *LogDocParser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *LogDocParser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *LogDocParser) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {
		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruleAction0:
			p.PushMap()
		case ruleAction1:
			p.PopMap()
		case ruleAction2:
			p.SetMapValue()
		case ruleAction3:
			p.PushList()
		case ruleAction4:
			p.PopList()
		case ruleAction5:
			p.SetListValue()
		case ruleAction6:
			p.PushField(buffer[begin:end])
		case ruleAction7:
			p.PushValue(p.Numeric(buffer[begin:end]))
		case ruleAction8:
			p.PushValue(buffer[begin:end])
		case ruleAction9:
			p.PushValue(nil)
		case ruleAction10:
			p.PushValue(true)
		case ruleAction11:
			p.PushValue(false)
		case ruleAction12:
			p.PushValue(p.Date(buffer[begin:end]))
		case ruleAction13:
			p.PushValue(p.ObjectId(buffer[begin:end]))
		case ruleAction14:
			p.PushValue(p.Bindata(buffer[begin:end]))
		case ruleAction15:
			p.PushValue(p.Regex(buffer[begin:end]))
		case ruleAction16:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction17:
			p.PushValue(p.Timestamp(buffer[begin:end]))
		case ruleAction18:
			p.PushValue(p.Numberlong(buffer[begin:end]))
		case ruleAction19:
			p.PushValue(p.Minkey())
		case ruleAction20:
			p.PushValue(p.Maxkey())
		case ruleAction21:
			p.PushValue(p.Undefined())

		}
	}
}

func (p *LogDocParser) Init() {
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
		/* 0 LogDoc <- <(Doc !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[ruleDoc]() {
					goto l0
				}
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
				}
				depth--
				add(ruleLogDoc, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Doc <- <('{' Action0 DocElements? '}' Action1)> */
		func() bool {
			position3, tokenIndex3, depth3 := position, tokenIndex, depth
			{
				position4 := position
				depth++
				if buffer[position] != rune('{') {
					goto l3
				}
				position++
				{
					add(ruleAction0, position)
				}
				{
					position6, tokenIndex6, depth6 := position, tokenIndex, depth
					{
						position8 := position
						depth++
						if !_rules[ruleDocElem]() {
							goto l6
						}
					l9:
						{
							position10, tokenIndex10, depth10 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l10
							}
							position++
							if !_rules[ruleDocElem]() {
								goto l10
							}
							goto l9
						l10:
							position, tokenIndex, depth = position10, tokenIndex10, depth10
						}
						depth--
						add(ruleDocElements, position8)
					}
					goto l7
				l6:
					position, tokenIndex, depth = position6, tokenIndex6, depth6
				}
			l7:
				if buffer[position] != rune('}') {
					goto l3
				}
				position++
				{
					add(ruleAction1, position)
				}
				depth--
				add(ruleDoc, position4)
			}
			return true
		l3:
			position, tokenIndex, depth = position3, tokenIndex3, depth3
			return false
		},
		/* 2 DocElements <- <(DocElem (',' DocElem)*)> */
		nil,
		/* 3 DocElem <- <(S? Field S? Value S? Action2)> */
		func() bool {
			position13, tokenIndex13, depth13 := position, tokenIndex, depth
			{
				position14 := position
				depth++
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l15
					}
					goto l16
				l15:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
				}
			l16:
				{
					position17 := position
					depth++
					{
						position18 := position
						depth++
						{
							position21 := position
							depth++
							{
								switch buffer[position] {
								case '$', '*', '.', '_':
									{
										switch buffer[position] {
										case '*':
											if buffer[position] != rune('*') {
												goto l13
											}
											position++
											break
										case '.':
											if buffer[position] != rune('.') {
												goto l13
											}
											position++
											break
										case '$':
											if buffer[position] != rune('$') {
												goto l13
											}
											position++
											break
										default:
											if buffer[position] != rune('_') {
												goto l13
											}
											position++
											break
										}
									}

									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l13
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l13
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l13
									}
									position++
									break
								}
							}

							depth--
							add(rulefieldChar, position21)
						}
					l19:
						{
							position20, tokenIndex20, depth20 := position, tokenIndex, depth
							{
								position24 := position
								depth++
								{
									switch buffer[position] {
									case '$', '*', '.', '_':
										{
											switch buffer[position] {
											case '*':
												if buffer[position] != rune('*') {
													goto l20
												}
												position++
												break
											case '.':
												if buffer[position] != rune('.') {
													goto l20
												}
												position++
												break
											case '$':
												if buffer[position] != rune('$') {
													goto l20
												}
												position++
												break
											default:
												if buffer[position] != rune('_') {
													goto l20
												}
												position++
												break
											}
										}

										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l20
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l20
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l20
										}
										position++
										break
									}
								}

								depth--
								add(rulefieldChar, position24)
							}
							goto l19
						l20:
							position, tokenIndex, depth = position20, tokenIndex20, depth20
						}
						depth--
						add(rulePegText, position18)
					}
					if buffer[position] != rune(':') {
						goto l13
					}
					position++
					{
						add(ruleAction6, position)
					}
					depth--
					add(ruleField, position17)
				}
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
				if !_rules[ruleValue]() {
					goto l13
				}
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l30
					}
					goto l31
				l30:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
				}
			l31:
				{
					add(ruleAction2, position)
				}
				depth--
				add(ruleDocElem, position14)
			}
			return true
		l13:
			position, tokenIndex, depth = position13, tokenIndex13, depth13
			return false
		},
		/* 4 List <- <('[' Action3 ListElements? ']' Action4)> */
		nil,
		/* 5 ListElements <- <(ListElem (',' ListElem)*)> */
		nil,
		/* 6 ListElem <- <(S? Value S? Action5)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l37
					}
					goto l38
				l37:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
				}
			l38:
				if !_rules[ruleValue]() {
					goto l35
				}
				{
					position39, tokenIndex39, depth39 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l39
					}
					goto l40
				l39:
					position, tokenIndex, depth = position39, tokenIndex39, depth39
				}
			l40:
				{
					add(ruleAction5, position)
				}
				depth--
				add(ruleListElem, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 7 Field <- <(<fieldChar+> ':' Action6)> */
		nil,
		/* 8 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') TimestampVal) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					{
						position47 := position
						depth++
						if buffer[position] != rune('n') {
							goto l46
						}
						position++
						if buffer[position] != rune('u') {
							goto l46
						}
						position++
						if buffer[position] != rune('l') {
							goto l46
						}
						position++
						if buffer[position] != rune('l') {
							goto l46
						}
						position++
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleNull, position47)
					}
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					{
						position50 := position
						depth++
						if buffer[position] != rune('M') {
							goto l49
						}
						position++
						if buffer[position] != rune('i') {
							goto l49
						}
						position++
						if buffer[position] != rune('n') {
							goto l49
						}
						position++
						if buffer[position] != rune('K') {
							goto l49
						}
						position++
						if buffer[position] != rune('e') {
							goto l49
						}
						position++
						if buffer[position] != rune('y') {
							goto l49
						}
						position++
						{
							add(ruleAction19, position)
						}
						depth--
						add(ruleMinKey, position50)
					}
					goto l45
				l49:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					{
						switch buffer[position] {
						case 'M':
							{
								position53 := position
								depth++
								if buffer[position] != rune('M') {
									goto l43
								}
								position++
								if buffer[position] != rune('a') {
									goto l43
								}
								position++
								if buffer[position] != rune('x') {
									goto l43
								}
								position++
								if buffer[position] != rune('K') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('y') {
									goto l43
								}
								position++
								{
									add(ruleAction20, position)
								}
								depth--
								add(ruleMaxKey, position53)
							}
							break
						case 'u':
							{
								position55 := position
								depth++
								if buffer[position] != rune('u') {
									goto l43
								}
								position++
								if buffer[position] != rune('n') {
									goto l43
								}
								position++
								if buffer[position] != rune('d') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('f') {
									goto l43
								}
								position++
								if buffer[position] != rune('i') {
									goto l43
								}
								position++
								if buffer[position] != rune('n') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('d') {
									goto l43
								}
								position++
								{
									add(ruleAction21, position)
								}
								depth--
								add(ruleUndefined, position55)
							}
							break
						case 'N':
							{
								position57 := position
								depth++
								if buffer[position] != rune('N') {
									goto l43
								}
								position++
								if buffer[position] != rune('u') {
									goto l43
								}
								position++
								if buffer[position] != rune('m') {
									goto l43
								}
								position++
								if buffer[position] != rune('b') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('r') {
									goto l43
								}
								position++
								if buffer[position] != rune('L') {
									goto l43
								}
								position++
								if buffer[position] != rune('o') {
									goto l43
								}
								position++
								if buffer[position] != rune('n') {
									goto l43
								}
								position++
								if buffer[position] != rune('g') {
									goto l43
								}
								position++
								if buffer[position] != rune('(') {
									goto l43
								}
								position++
								{
									position58 := position
									depth++
									{
										position61, tokenIndex61, depth61 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l61
										}
										position++
										goto l43
									l61:
										position, tokenIndex, depth = position61, tokenIndex61, depth61
									}
									if !matchDot() {
										goto l43
									}
								l59:
									{
										position60, tokenIndex60, depth60 := position, tokenIndex, depth
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l62
											}
											position++
											goto l60
										l62:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
										}
										if !matchDot() {
											goto l60
										}
										goto l59
									l60:
										position, tokenIndex, depth = position60, tokenIndex60, depth60
									}
									depth--
									add(rulePegText, position58)
								}
								if buffer[position] != rune(')') {
									goto l43
								}
								position++
								{
									add(ruleAction18, position)
								}
								depth--
								add(ruleNumberLong, position57)
							}
							break
						case '/':
							{
								position64 := position
								depth++
								if buffer[position] != rune('/') {
									goto l43
								}
								position++
								{
									position65 := position
									depth++
									{
										position66 := position
										depth++
										{
											position69 := position
											depth++
											{
												position70, tokenIndex70, depth70 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l70
												}
												position++
												goto l43
											l70:
												position, tokenIndex, depth = position70, tokenIndex70, depth70
											}
											if !matchDot() {
												goto l43
											}
											depth--
											add(ruleregexChar, position69)
										}
									l67:
										{
											position68, tokenIndex68, depth68 := position, tokenIndex, depth
											{
												position71 := position
												depth++
												{
													position72, tokenIndex72, depth72 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l72
													}
													position++
													goto l68
												l72:
													position, tokenIndex, depth = position72, tokenIndex72, depth72
												}
												if !matchDot() {
													goto l68
												}
												depth--
												add(ruleregexChar, position71)
											}
											goto l67
										l68:
											position, tokenIndex, depth = position68, tokenIndex68, depth68
										}
										if buffer[position] != rune('/') {
											goto l43
										}
										position++
									l73:
										{
											position74, tokenIndex74, depth74 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l74
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l74
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l74
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l74
													}
													position++
													break
												}
											}

											goto l73
										l74:
											position, tokenIndex, depth = position74, tokenIndex74, depth74
										}
										depth--
										add(ruleregexBody, position66)
									}
									depth--
									add(rulePegText, position65)
								}
								{
									add(ruleAction15, position)
								}
								depth--
								add(ruleRegex, position64)
							}
							break
						case 'T':
							{
								position77 := position
								depth++
								{
									position78, tokenIndex78, depth78 := position, tokenIndex, depth
									{
										position80 := position
										depth++
										if buffer[position] != rune('T') {
											goto l79
										}
										position++
										if buffer[position] != rune('i') {
											goto l79
										}
										position++
										if buffer[position] != rune('m') {
											goto l79
										}
										position++
										if buffer[position] != rune('e') {
											goto l79
										}
										position++
										if buffer[position] != rune('s') {
											goto l79
										}
										position++
										if buffer[position] != rune('t') {
											goto l79
										}
										position++
										if buffer[position] != rune('a') {
											goto l79
										}
										position++
										if buffer[position] != rune('m') {
											goto l79
										}
										position++
										if buffer[position] != rune('p') {
											goto l79
										}
										position++
										if buffer[position] != rune('(') {
											goto l79
										}
										position++
										{
											position81 := position
											depth++
											{
												position84, tokenIndex84, depth84 := position, tokenIndex, depth
												if buffer[position] != rune(')') {
													goto l84
												}
												position++
												goto l79
											l84:
												position, tokenIndex, depth = position84, tokenIndex84, depth84
											}
											if !matchDot() {
												goto l79
											}
										l82:
											{
												position83, tokenIndex83, depth83 := position, tokenIndex, depth
												{
													position85, tokenIndex85, depth85 := position, tokenIndex, depth
													if buffer[position] != rune(')') {
														goto l85
													}
													position++
													goto l83
												l85:
													position, tokenIndex, depth = position85, tokenIndex85, depth85
												}
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
										if buffer[position] != rune(')') {
											goto l79
										}
										position++
										{
											add(ruleAction16, position)
										}
										depth--
										add(ruletimestampParen, position80)
									}
									goto l78
								l79:
									position, tokenIndex, depth = position78, tokenIndex78, depth78
									{
										position87 := position
										depth++
										if buffer[position] != rune('T') {
											goto l43
										}
										position++
										if buffer[position] != rune('i') {
											goto l43
										}
										position++
										if buffer[position] != rune('m') {
											goto l43
										}
										position++
										if buffer[position] != rune('e') {
											goto l43
										}
										position++
										if buffer[position] != rune('s') {
											goto l43
										}
										position++
										if buffer[position] != rune('t') {
											goto l43
										}
										position++
										if buffer[position] != rune('a') {
											goto l43
										}
										position++
										if buffer[position] != rune('m') {
											goto l43
										}
										position++
										if buffer[position] != rune('p') {
											goto l43
										}
										position++
										if buffer[position] != rune(' ') {
											goto l43
										}
										position++
										{
											position88 := position
											depth++
											{
												position91, tokenIndex91, depth91 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l92
												}
												position++
												goto l91
											l92:
												position, tokenIndex, depth = position91, tokenIndex91, depth91
												if buffer[position] != rune('|') {
													goto l43
												}
												position++
											}
										l91:
										l89:
											{
												position90, tokenIndex90, depth90 := position, tokenIndex, depth
												{
													position93, tokenIndex93, depth93 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l94
													}
													position++
													goto l93
												l94:
													position, tokenIndex, depth = position93, tokenIndex93, depth93
													if buffer[position] != rune('|') {
														goto l90
													}
													position++
												}
											l93:
												goto l89
											l90:
												position, tokenIndex, depth = position90, tokenIndex90, depth90
											}
											depth--
											add(rulePegText, position88)
										}
										{
											add(ruleAction17, position)
										}
										depth--
										add(ruletimestampPipe, position87)
									}
								}
							l78:
								depth--
								add(ruleTimestampVal, position77)
							}
							break
						case 'B':
							{
								position96 := position
								depth++
								if buffer[position] != rune('B') {
									goto l43
								}
								position++
								if buffer[position] != rune('i') {
									goto l43
								}
								position++
								if buffer[position] != rune('n') {
									goto l43
								}
								position++
								if buffer[position] != rune('D') {
									goto l43
								}
								position++
								if buffer[position] != rune('a') {
									goto l43
								}
								position++
								if buffer[position] != rune('t') {
									goto l43
								}
								position++
								if buffer[position] != rune('a') {
									goto l43
								}
								position++
								if buffer[position] != rune('(') {
									goto l43
								}
								position++
								{
									position97 := position
									depth++
									{
										position100, tokenIndex100, depth100 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l100
										}
										position++
										goto l43
									l100:
										position, tokenIndex, depth = position100, tokenIndex100, depth100
									}
									if !matchDot() {
										goto l43
									}
								l98:
									{
										position99, tokenIndex99, depth99 := position, tokenIndex, depth
										{
											position101, tokenIndex101, depth101 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l101
											}
											position++
											goto l99
										l101:
											position, tokenIndex, depth = position101, tokenIndex101, depth101
										}
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
								if buffer[position] != rune(')') {
									goto l43
								}
								position++
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleBinData, position96)
							}
							break
						case 'D', 'n':
							{
								position103 := position
								depth++
								{
									position104, tokenIndex104, depth104 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l104
									}
									position++
									if buffer[position] != rune('e') {
										goto l104
									}
									position++
									if buffer[position] != rune('w') {
										goto l104
									}
									position++
									if buffer[position] != rune(' ') {
										goto l104
									}
									position++
									goto l105
								l104:
									position, tokenIndex, depth = position104, tokenIndex104, depth104
								}
							l105:
								if buffer[position] != rune('D') {
									goto l43
								}
								position++
								if buffer[position] != rune('a') {
									goto l43
								}
								position++
								if buffer[position] != rune('t') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('(') {
									goto l43
								}
								position++
								{
									position106 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l43
									}
									position++
								l107:
									{
										position108, tokenIndex108, depth108 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l108
										}
										position++
										goto l107
									l108:
										position, tokenIndex, depth = position108, tokenIndex108, depth108
									}
									depth--
									add(rulePegText, position106)
								}
								if buffer[position] != rune(')') {
									goto l43
								}
								position++
								{
									add(ruleAction12, position)
								}
								depth--
								add(ruleDate, position103)
							}
							break
						case 'O':
							{
								position110 := position
								depth++
								if buffer[position] != rune('O') {
									goto l43
								}
								position++
								if buffer[position] != rune('b') {
									goto l43
								}
								position++
								if buffer[position] != rune('j') {
									goto l43
								}
								position++
								if buffer[position] != rune('e') {
									goto l43
								}
								position++
								if buffer[position] != rune('c') {
									goto l43
								}
								position++
								if buffer[position] != rune('t') {
									goto l43
								}
								position++
								if buffer[position] != rune('I') {
									goto l43
								}
								position++
								if buffer[position] != rune('d') {
									goto l43
								}
								position++
								if buffer[position] != rune('(') {
									goto l43
								}
								position++
								{
									position111, tokenIndex111, depth111 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l112
									}
									position++
									goto l111
								l112:
									position, tokenIndex, depth = position111, tokenIndex111, depth111
									if buffer[position] != rune('"') {
										goto l43
									}
									position++
								}
							l111:
								{
									position113 := position
									depth++
								l114:
									{
										position115, tokenIndex115, depth115 := position, tokenIndex, depth
										{
											position116 := position
											depth++
											{
												position117, tokenIndex117, depth117 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
												goto l117
											l118:
												position, tokenIndex, depth = position117, tokenIndex117, depth117
												{
													position119, tokenIndex119, depth119 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l120
													}
													position++
													goto l119
												l120:
													position, tokenIndex, depth = position119, tokenIndex119, depth119
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l115
													}
													position++
												}
											l119:
											}
										l117:
											depth--
											add(rulehexChar, position116)
										}
										goto l114
									l115:
										position, tokenIndex, depth = position115, tokenIndex115, depth115
									}
									depth--
									add(rulePegText, position113)
								}
								{
									position121, tokenIndex121, depth121 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l122
									}
									position++
									goto l121
								l122:
									position, tokenIndex, depth = position121, tokenIndex121, depth121
									if buffer[position] != rune('"') {
										goto l43
									}
									position++
								}
							l121:
								if buffer[position] != rune(')') {
									goto l43
								}
								position++
								{
									add(ruleAction13, position)
								}
								depth--
								add(ruleObjectID, position110)
							}
							break
						case '"':
							{
								position124 := position
								depth++
								if buffer[position] != rune('"') {
									goto l43
								}
								position++
								{
									position125 := position
									depth++
								l126:
									{
										position127, tokenIndex127, depth127 := position, tokenIndex, depth
										{
											position128 := position
											depth++
											{
												position129, tokenIndex129, depth129 := position, tokenIndex, depth
												{
													position131, tokenIndex131, depth131 := position, tokenIndex, depth
													{
														position132, tokenIndex132, depth132 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l133
														}
														position++
														goto l132
													l133:
														position, tokenIndex, depth = position132, tokenIndex132, depth132
														if buffer[position] != rune('\\') {
															goto l131
														}
														position++
													}
												l132:
													goto l130
												l131:
													position, tokenIndex, depth = position131, tokenIndex131, depth131
												}
												if !matchDot() {
													goto l130
												}
												goto l129
											l130:
												position, tokenIndex, depth = position129, tokenIndex129, depth129
												if buffer[position] != rune('\\') {
													goto l127
												}
												position++
												{
													position134, tokenIndex134, depth134 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l135
													}
													position++
													goto l134
												l135:
													position, tokenIndex, depth = position134, tokenIndex134, depth134
													if buffer[position] != rune('\\') {
														goto l127
													}
													position++
												}
											l134:
											}
										l129:
											depth--
											add(rulestringChar, position128)
										}
										goto l126
									l127:
										position, tokenIndex, depth = position127, tokenIndex127, depth127
									}
									depth--
									add(rulePegText, position125)
								}
								if buffer[position] != rune('"') {
									goto l43
								}
								position++
								{
									add(ruleAction8, position)
								}
								depth--
								add(ruleString, position124)
							}
							break
						case 'f', 't':
							{
								position137 := position
								depth++
								{
									position138, tokenIndex138, depth138 := position, tokenIndex, depth
									{
										position140 := position
										depth++
										if buffer[position] != rune('t') {
											goto l139
										}
										position++
										if buffer[position] != rune('r') {
											goto l139
										}
										position++
										if buffer[position] != rune('u') {
											goto l139
										}
										position++
										if buffer[position] != rune('e') {
											goto l139
										}
										position++
										{
											add(ruleAction10, position)
										}
										depth--
										add(ruleTrue, position140)
									}
									goto l138
								l139:
									position, tokenIndex, depth = position138, tokenIndex138, depth138
									{
										position142 := position
										depth++
										if buffer[position] != rune('f') {
											goto l43
										}
										position++
										if buffer[position] != rune('a') {
											goto l43
										}
										position++
										if buffer[position] != rune('l') {
											goto l43
										}
										position++
										if buffer[position] != rune('s') {
											goto l43
										}
										position++
										if buffer[position] != rune('e') {
											goto l43
										}
										position++
										{
											add(ruleAction11, position)
										}
										depth--
										add(ruleFalse, position142)
									}
								}
							l138:
								depth--
								add(ruleBoolean, position137)
							}
							break
						case '[':
							{
								position144 := position
								depth++
								if buffer[position] != rune('[') {
									goto l43
								}
								position++
								{
									add(ruleAction3, position)
								}
								{
									position146, tokenIndex146, depth146 := position, tokenIndex, depth
									{
										position148 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l146
										}
									l149:
										{
											position150, tokenIndex150, depth150 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l150
											}
											position++
											if !_rules[ruleListElem]() {
												goto l150
											}
											goto l149
										l150:
											position, tokenIndex, depth = position150, tokenIndex150, depth150
										}
										depth--
										add(ruleListElements, position148)
									}
									goto l147
								l146:
									position, tokenIndex, depth = position146, tokenIndex146, depth146
								}
							l147:
								if buffer[position] != rune(']') {
									goto l43
								}
								position++
								{
									add(ruleAction4, position)
								}
								depth--
								add(ruleList, position144)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l43
							}
							break
						default:
							{
								position152 := position
								depth++
								{
									position153 := position
									depth++
									{
										position154, tokenIndex154, depth154 := position, tokenIndex, depth
										if buffer[position] != rune('-') {
											goto l154
										}
										position++
										goto l155
									l154:
										position, tokenIndex, depth = position154, tokenIndex154, depth154
									}
								l155:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l43
									}
									position++
								l156:
									{
										position157, tokenIndex157, depth157 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l157
										}
										position++
										goto l156
									l157:
										position, tokenIndex, depth = position157, tokenIndex157, depth157
									}
									{
										position158, tokenIndex158, depth158 := position, tokenIndex, depth
										if buffer[position] != rune('.') {
											goto l158
										}
										position++
										goto l159
									l158:
										position, tokenIndex, depth = position158, tokenIndex158, depth158
									}
								l159:
								l160:
									{
										position161, tokenIndex161, depth161 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l161
										}
										position++
										goto l160
									l161:
										position, tokenIndex, depth = position161, tokenIndex161, depth161
									}
									depth--
									add(rulePegText, position153)
								}
								{
									add(ruleAction7, position)
								}
								depth--
								add(ruleNumeric, position152)
							}
							break
						}
					}

				}
			l45:
				depth--
				add(ruleValue, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 9 Numeric <- <(<('-'? [0-9]+ '.'? [0-9]*)> Action7)> */
		nil,
		/* 10 Boolean <- <(True / False)> */
		nil,
		/* 11 String <- <('"' <stringChar*> '"' Action8)> */
		nil,
		/* 12 Null <- <('n' 'u' 'l' 'l' Action9)> */
		nil,
		/* 13 True <- <('t' 'r' 'u' 'e' Action10)> */
		nil,
		/* 14 False <- <('f' 'a' 'l' 's' 'e' Action11)> */
		nil,
		/* 15 Date <- <(('n' 'e' 'w' ' ')? ('D' 'a' 't' 'e' '(') <[0-9]+> ')' Action12)> */
		nil,
		/* 16 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' ('\'' / '"') <hexChar*> ('\'' / '"') ')' Action13)> */
		nil,
		/* 17 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action14)> */
		nil,
		/* 18 Regex <- <('/' <regexBody> Action15)> */
		nil,
		/* 19 TimestampVal <- <(timestampParen / timestampPipe)> */
		nil,
		/* 20 timestampParen <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action16)> */
		nil,
		/* 21 timestampPipe <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' ' ' <([0-9] / '|')+> Action17)> */
		nil,
		/* 22 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action18)> */
		nil,
		/* 23 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action19)> */
		nil,
		/* 24 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action20)> */
		nil,
		/* 25 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action21)> */
		nil,
		/* 26 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 27 regexChar <- <(!'/' .)> */
		nil,
		/* 28 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 29 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 30 fieldChar <- <((&('$' | '*' | '.' | '_') ((&('*') '*') | (&('.') '.') | (&('$') '$') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 31 S <- <' '> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l185
				}
				position++
				depth--
				add(ruleS, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 33 Action0 <- <{ p.PushMap() }> */
		nil,
		/* 34 Action1 <- <{ p.PopMap() }> */
		nil,
		/* 35 Action2 <- <{ p.SetMapValue() }> */
		nil,
		/* 36 Action3 <- <{ p.PushList() }> */
		nil,
		/* 37 Action4 <- <{ p.PopList() }> */
		nil,
		/* 38 Action5 <- <{ p.SetListValue() }> */
		nil,
		nil,
		/* 40 Action6 <- <{ p.PushField(buffer[begin:end]) }> */
		nil,
		/* 41 Action7 <- <{ p.PushValue(p.Numeric(buffer[begin:end])) }> */
		nil,
		/* 42 Action8 <- <{ p.PushValue(buffer[begin:end]) }> */
		nil,
		/* 43 Action9 <- <{ p.PushValue(nil) }> */
		nil,
		/* 44 Action10 <- <{ p.PushValue(true) }> */
		nil,
		/* 45 Action11 <- <{ p.PushValue(false) }> */
		nil,
		/* 46 Action12 <- <{ p.PushValue(p.Date(buffer[begin:end])) }> */
		nil,
		/* 47 Action13 <- <{ p.PushValue(p.ObjectId(buffer[begin:end])) }> */
		nil,
		/* 48 Action14 <- <{ p.PushValue(p.Bindata(buffer[begin:end])) }> */
		nil,
		/* 49 Action15 <- <{ p.PushValue(p.Regex(buffer[begin:end])) }> */
		nil,
		/* 50 Action16 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 51 Action17 <- <{ p.PushValue(p.Timestamp(buffer[begin:end])) }> */
		nil,
		/* 52 Action18 <- <{ p.PushValue(p.Numberlong(buffer[begin:end])) }> */
		nil,
		/* 53 Action19 <- <{ p.PushValue(p.Minkey()) }> */
		nil,
		/* 54 Action20 <- <{ p.PushValue(p.Maxkey()) }> */
		nil,
		/* 55 Action21 <- <{ p.PushValue(p.Undefined()) }> */
		nil,
	}
	p.rules = _rules
}
