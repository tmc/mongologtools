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
	ruleTimestamp
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
	"Timestamp",
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
	rules  [53]func() bool
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
			p.pushMap()
		case ruleAction1:
			p.popMap()
		case ruleAction2:
			p.setMapValue()
		case ruleAction3:
			p.pushList()
		case ruleAction4:
			p.popList()
		case ruleAction5:
			p.setListValue()
		case ruleAction6:
			p.pushField(buffer[begin:end])
		case ruleAction7:
			p.pushValue(numeric(buffer[begin:end]))
		case ruleAction8:
			p.pushValue(buffer[begin:end])
		case ruleAction9:
			p.pushValue(nil)
		case ruleAction10:
			p.pushValue(true)
		case ruleAction11:
			p.pushValue(false)
		case ruleAction12:
			p.pushValue(date(buffer[begin:end]))
		case ruleAction13:
			p.pushValue(objectid(buffer[begin:end]))
		case ruleAction14:
			p.pushValue(bindata(buffer[begin:end]))
		case ruleAction15:
			p.pushValue(regex(buffer[begin:end]))
		case ruleAction16:
			p.pushValue(timestamp(buffer[begin:end]))
		case ruleAction17:
			p.pushValue(numberlong(buffer[begin:end]))
		case ruleAction18:
			p.pushValue(minkey())
		case ruleAction19:
			p.pushValue(maxkey())
		case ruleAction20:
			p.pushValue(undefined())

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
							{
								position11, tokenIndex11, depth11 := position, tokenIndex, depth
								if !_rules[ruleS]() {
									goto l11
								}
								goto l12
							l11:
								position, tokenIndex, depth = position11, tokenIndex11, depth11
							}
						l12:
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
		/* 2 DocElements <- <(DocElem (',' S? DocElem)*)> */
		nil,
		/* 3 DocElem <- <(S? Field S? Value S? Action2)> */
		func() bool {
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				{
					position17, tokenIndex17, depth17 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l17
					}
					goto l18
				l17:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
				}
			l18:
				{
					position19 := position
					depth++
					{
						position20 := position
						depth++
						{
							position23 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position25, tokenIndex25, depth25 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l26
										}
										position++
										goto l25
									l26:
										position, tokenIndex, depth = position25, tokenIndex25, depth25
										if buffer[position] != rune('$') {
											goto l15
										}
										position++
									}
								l25:
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l15
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l15
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l15
									}
									position++
									break
								}
							}

							depth--
							add(rulefieldChar, position23)
						}
					l21:
						{
							position22, tokenIndex22, depth22 := position, tokenIndex, depth
							{
								position27 := position
								depth++
								{
									switch buffer[position] {
									case '$', '_':
										{
											position29, tokenIndex29, depth29 := position, tokenIndex, depth
											if buffer[position] != rune('_') {
												goto l30
											}
											position++
											goto l29
										l30:
											position, tokenIndex, depth = position29, tokenIndex29, depth29
											if buffer[position] != rune('$') {
												goto l22
											}
											position++
										}
									l29:
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l22
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l22
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l22
										}
										position++
										break
									}
								}

								depth--
								add(rulefieldChar, position27)
							}
							goto l21
						l22:
							position, tokenIndex, depth = position22, tokenIndex22, depth22
						}
						depth--
						add(rulePegText, position20)
					}
					if buffer[position] != rune(':') {
						goto l15
					}
					position++
					{
						add(ruleAction6, position)
					}
					depth--
					add(ruleField, position19)
				}
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l32
					}
					goto l33
				l32:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
				}
			l33:
				if !_rules[ruleValue]() {
					goto l15
				}
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l34
					}
					goto l35
				l34:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
				}
			l35:
				{
					add(ruleAction2, position)
				}
				depth--
				add(ruleDocElem, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 4 List <- <('[' Action3 ListElements? ']' Action4)> */
		nil,
		/* 5 ListElements <- <(ListElem (',' S? ListElem)*)> */
		nil,
		/* 6 ListElem <- <(Value S? Action5)> */
		func() bool {
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				if !_rules[ruleValue]() {
					goto l39
				}
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if !_rules[ruleS]() {
						goto l41
					}
					goto l42
				l41:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
				}
			l42:
				{
					add(ruleAction5, position)
				}
				depth--
				add(ruleListElem, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 7 Field <- <(<fieldChar+> ':' Action6)> */
		nil,
		/* 8 Value <- <(Null / MinKey / ((&('M') MaxKey) | (&('u') Undefined) | (&('N') NumberLong) | (&('/') Regex) | (&('T') Timestamp) | (&('B') BinData) | (&('D' | 'n') Date) | (&('O') ObjectID) | (&('"') String) | (&('f' | 't') Boolean) | (&('[') List) | (&('{') Doc) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		func() bool {
			position45, tokenIndex45, depth45 := position, tokenIndex, depth
			{
				position46 := position
				depth++
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					{
						position49 := position
						depth++
						if buffer[position] != rune('n') {
							goto l48
						}
						position++
						if buffer[position] != rune('u') {
							goto l48
						}
						position++
						if buffer[position] != rune('l') {
							goto l48
						}
						position++
						if buffer[position] != rune('l') {
							goto l48
						}
						position++
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleNull, position49)
					}
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					{
						position52 := position
						depth++
						if buffer[position] != rune('M') {
							goto l51
						}
						position++
						if buffer[position] != rune('i') {
							goto l51
						}
						position++
						if buffer[position] != rune('n') {
							goto l51
						}
						position++
						if buffer[position] != rune('K') {
							goto l51
						}
						position++
						if buffer[position] != rune('e') {
							goto l51
						}
						position++
						if buffer[position] != rune('y') {
							goto l51
						}
						position++
						{
							add(ruleAction18, position)
						}
						depth--
						add(ruleMinKey, position52)
					}
					goto l47
				l51:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					{
						switch buffer[position] {
						case 'M':
							{
								position55 := position
								depth++
								if buffer[position] != rune('M') {
									goto l45
								}
								position++
								if buffer[position] != rune('a') {
									goto l45
								}
								position++
								if buffer[position] != rune('x') {
									goto l45
								}
								position++
								if buffer[position] != rune('K') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('y') {
									goto l45
								}
								position++
								{
									add(ruleAction19, position)
								}
								depth--
								add(ruleMaxKey, position55)
							}
							break
						case 'u':
							{
								position57 := position
								depth++
								if buffer[position] != rune('u') {
									goto l45
								}
								position++
								if buffer[position] != rune('n') {
									goto l45
								}
								position++
								if buffer[position] != rune('d') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('f') {
									goto l45
								}
								position++
								if buffer[position] != rune('i') {
									goto l45
								}
								position++
								if buffer[position] != rune('n') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('d') {
									goto l45
								}
								position++
								{
									add(ruleAction20, position)
								}
								depth--
								add(ruleUndefined, position57)
							}
							break
						case 'N':
							{
								position59 := position
								depth++
								if buffer[position] != rune('N') {
									goto l45
								}
								position++
								if buffer[position] != rune('u') {
									goto l45
								}
								position++
								if buffer[position] != rune('m') {
									goto l45
								}
								position++
								if buffer[position] != rune('b') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('r') {
									goto l45
								}
								position++
								if buffer[position] != rune('L') {
									goto l45
								}
								position++
								if buffer[position] != rune('o') {
									goto l45
								}
								position++
								if buffer[position] != rune('n') {
									goto l45
								}
								position++
								if buffer[position] != rune('g') {
									goto l45
								}
								position++
								if buffer[position] != rune('(') {
									goto l45
								}
								position++
								{
									position60 := position
									depth++
									{
										position63, tokenIndex63, depth63 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l63
										}
										position++
										goto l45
									l63:
										position, tokenIndex, depth = position63, tokenIndex63, depth63
									}
									if !matchDot() {
										goto l45
									}
								l61:
									{
										position62, tokenIndex62, depth62 := position, tokenIndex, depth
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l64
											}
											position++
											goto l62
										l64:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
										}
										if !matchDot() {
											goto l62
										}
										goto l61
									l62:
										position, tokenIndex, depth = position62, tokenIndex62, depth62
									}
									depth--
									add(rulePegText, position60)
								}
								if buffer[position] != rune(')') {
									goto l45
								}
								position++
								{
									add(ruleAction17, position)
								}
								depth--
								add(ruleNumberLong, position59)
							}
							break
						case '/':
							{
								position66 := position
								depth++
								if buffer[position] != rune('/') {
									goto l45
								}
								position++
								{
									position67 := position
									depth++
									{
										position68 := position
										depth++
										{
											position71 := position
											depth++
											{
												position72, tokenIndex72, depth72 := position, tokenIndex, depth
												if buffer[position] != rune('/') {
													goto l72
												}
												position++
												goto l45
											l72:
												position, tokenIndex, depth = position72, tokenIndex72, depth72
											}
											if !matchDot() {
												goto l45
											}
											depth--
											add(ruleregexChar, position71)
										}
									l69:
										{
											position70, tokenIndex70, depth70 := position, tokenIndex, depth
											{
												position73 := position
												depth++
												{
													position74, tokenIndex74, depth74 := position, tokenIndex, depth
													if buffer[position] != rune('/') {
														goto l74
													}
													position++
													goto l70
												l74:
													position, tokenIndex, depth = position74, tokenIndex74, depth74
												}
												if !matchDot() {
													goto l70
												}
												depth--
												add(ruleregexChar, position73)
											}
											goto l69
										l70:
											position, tokenIndex, depth = position70, tokenIndex70, depth70
										}
										if buffer[position] != rune('/') {
											goto l45
										}
										position++
									l75:
										{
											position76, tokenIndex76, depth76 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 's':
													if buffer[position] != rune('s') {
														goto l76
													}
													position++
													break
												case 'm':
													if buffer[position] != rune('m') {
														goto l76
													}
													position++
													break
												case 'i':
													if buffer[position] != rune('i') {
														goto l76
													}
													position++
													break
												default:
													if buffer[position] != rune('g') {
														goto l76
													}
													position++
													break
												}
											}

											goto l75
										l76:
											position, tokenIndex, depth = position76, tokenIndex76, depth76
										}
										depth--
										add(ruleregexBody, position68)
									}
									depth--
									add(rulePegText, position67)
								}
								{
									add(ruleAction15, position)
								}
								depth--
								add(ruleRegex, position66)
							}
							break
						case 'T':
							{
								position79 := position
								depth++
								if buffer[position] != rune('T') {
									goto l45
								}
								position++
								if buffer[position] != rune('i') {
									goto l45
								}
								position++
								if buffer[position] != rune('m') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('s') {
									goto l45
								}
								position++
								if buffer[position] != rune('t') {
									goto l45
								}
								position++
								if buffer[position] != rune('a') {
									goto l45
								}
								position++
								if buffer[position] != rune('m') {
									goto l45
								}
								position++
								if buffer[position] != rune('p') {
									goto l45
								}
								position++
								if buffer[position] != rune('(') {
									goto l45
								}
								position++
								{
									position80 := position
									depth++
									{
										position83, tokenIndex83, depth83 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l83
										}
										position++
										goto l45
									l83:
										position, tokenIndex, depth = position83, tokenIndex83, depth83
									}
									if !matchDot() {
										goto l45
									}
								l81:
									{
										position82, tokenIndex82, depth82 := position, tokenIndex, depth
										{
											position84, tokenIndex84, depth84 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l84
											}
											position++
											goto l82
										l84:
											position, tokenIndex, depth = position84, tokenIndex84, depth84
										}
										if !matchDot() {
											goto l82
										}
										goto l81
									l82:
										position, tokenIndex, depth = position82, tokenIndex82, depth82
									}
									depth--
									add(rulePegText, position80)
								}
								if buffer[position] != rune(')') {
									goto l45
								}
								position++
								{
									add(ruleAction16, position)
								}
								depth--
								add(ruleTimestamp, position79)
							}
							break
						case 'B':
							{
								position86 := position
								depth++
								if buffer[position] != rune('B') {
									goto l45
								}
								position++
								if buffer[position] != rune('i') {
									goto l45
								}
								position++
								if buffer[position] != rune('n') {
									goto l45
								}
								position++
								if buffer[position] != rune('D') {
									goto l45
								}
								position++
								if buffer[position] != rune('a') {
									goto l45
								}
								position++
								if buffer[position] != rune('t') {
									goto l45
								}
								position++
								if buffer[position] != rune('a') {
									goto l45
								}
								position++
								if buffer[position] != rune('(') {
									goto l45
								}
								position++
								{
									position87 := position
									depth++
									{
										position90, tokenIndex90, depth90 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l90
										}
										position++
										goto l45
									l90:
										position, tokenIndex, depth = position90, tokenIndex90, depth90
									}
									if !matchDot() {
										goto l45
									}
								l88:
									{
										position89, tokenIndex89, depth89 := position, tokenIndex, depth
										{
											position91, tokenIndex91, depth91 := position, tokenIndex, depth
											if buffer[position] != rune(')') {
												goto l91
											}
											position++
											goto l89
										l91:
											position, tokenIndex, depth = position91, tokenIndex91, depth91
										}
										if !matchDot() {
											goto l89
										}
										goto l88
									l89:
										position, tokenIndex, depth = position89, tokenIndex89, depth89
									}
									depth--
									add(rulePegText, position87)
								}
								if buffer[position] != rune(')') {
									goto l45
								}
								position++
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleBinData, position86)
							}
							break
						case 'D', 'n':
							{
								position93 := position
								depth++
								{
									position94, tokenIndex94, depth94 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l94
									}
									position++
									if buffer[position] != rune('e') {
										goto l94
									}
									position++
									if buffer[position] != rune('w') {
										goto l94
									}
									position++
									if buffer[position] != rune(' ') {
										goto l94
									}
									position++
									goto l95
								l94:
									position, tokenIndex, depth = position94, tokenIndex94, depth94
								}
							l95:
								if buffer[position] != rune('D') {
									goto l45
								}
								position++
								if buffer[position] != rune('a') {
									goto l45
								}
								position++
								if buffer[position] != rune('t') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('(') {
									goto l45
								}
								position++
								{
									position96 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l45
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
								if buffer[position] != rune(')') {
									goto l45
								}
								position++
								{
									add(ruleAction12, position)
								}
								depth--
								add(ruleDate, position93)
							}
							break
						case 'O':
							{
								position100 := position
								depth++
								if buffer[position] != rune('O') {
									goto l45
								}
								position++
								if buffer[position] != rune('b') {
									goto l45
								}
								position++
								if buffer[position] != rune('j') {
									goto l45
								}
								position++
								if buffer[position] != rune('e') {
									goto l45
								}
								position++
								if buffer[position] != rune('c') {
									goto l45
								}
								position++
								if buffer[position] != rune('t') {
									goto l45
								}
								position++
								if buffer[position] != rune('I') {
									goto l45
								}
								position++
								if buffer[position] != rune('d') {
									goto l45
								}
								position++
								if buffer[position] != rune('(') {
									goto l45
								}
								position++
								if buffer[position] != rune('"') {
									goto l45
								}
								position++
								{
									position101 := position
									depth++
								l102:
									{
										position103, tokenIndex103, depth103 := position, tokenIndex, depth
										{
											position104 := position
											depth++
											{
												position105, tokenIndex105, depth105 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l106
												}
												position++
												goto l105
											l106:
												position, tokenIndex, depth = position105, tokenIndex105, depth105
												{
													position107, tokenIndex107, depth107 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('f') {
														goto l108
													}
													position++
													goto l107
												l108:
													position, tokenIndex, depth = position107, tokenIndex107, depth107
													if c := buffer[position]; c < rune('A') || c > rune('F') {
														goto l103
													}
													position++
												}
											l107:
											}
										l105:
											depth--
											add(rulehexChar, position104)
										}
										goto l102
									l103:
										position, tokenIndex, depth = position103, tokenIndex103, depth103
									}
									depth--
									add(rulePegText, position101)
								}
								if buffer[position] != rune('"') {
									goto l45
								}
								position++
								if buffer[position] != rune(')') {
									goto l45
								}
								position++
								{
									add(ruleAction13, position)
								}
								depth--
								add(ruleObjectID, position100)
							}
							break
						case '"':
							{
								position110 := position
								depth++
								if buffer[position] != rune('"') {
									goto l45
								}
								position++
								{
									position111 := position
									depth++
								l112:
									{
										position113, tokenIndex113, depth113 := position, tokenIndex, depth
										{
											position114 := position
											depth++
											{
												position115, tokenIndex115, depth115 := position, tokenIndex, depth
												{
													position117, tokenIndex117, depth117 := position, tokenIndex, depth
													{
														position118, tokenIndex118, depth118 := position, tokenIndex, depth
														if buffer[position] != rune('"') {
															goto l119
														}
														position++
														goto l118
													l119:
														position, tokenIndex, depth = position118, tokenIndex118, depth118
														if buffer[position] != rune('\\') {
															goto l117
														}
														position++
													}
												l118:
													goto l116
												l117:
													position, tokenIndex, depth = position117, tokenIndex117, depth117
												}
												if !matchDot() {
													goto l116
												}
												goto l115
											l116:
												position, tokenIndex, depth = position115, tokenIndex115, depth115
												if buffer[position] != rune('\\') {
													goto l113
												}
												position++
												{
													position120, tokenIndex120, depth120 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l121
													}
													position++
													goto l120
												l121:
													position, tokenIndex, depth = position120, tokenIndex120, depth120
													if buffer[position] != rune('\\') {
														goto l113
													}
													position++
												}
											l120:
											}
										l115:
											depth--
											add(rulestringChar, position114)
										}
										goto l112
									l113:
										position, tokenIndex, depth = position113, tokenIndex113, depth113
									}
									depth--
									add(rulePegText, position111)
								}
								if buffer[position] != rune('"') {
									goto l45
								}
								position++
								{
									add(ruleAction8, position)
								}
								depth--
								add(ruleString, position110)
							}
							break
						case 'f', 't':
							{
								position123 := position
								depth++
								{
									position124, tokenIndex124, depth124 := position, tokenIndex, depth
									{
										position126 := position
										depth++
										if buffer[position] != rune('t') {
											goto l125
										}
										position++
										if buffer[position] != rune('r') {
											goto l125
										}
										position++
										if buffer[position] != rune('u') {
											goto l125
										}
										position++
										if buffer[position] != rune('e') {
											goto l125
										}
										position++
										{
											add(ruleAction10, position)
										}
										depth--
										add(ruleTrue, position126)
									}
									goto l124
								l125:
									position, tokenIndex, depth = position124, tokenIndex124, depth124
									{
										position128 := position
										depth++
										if buffer[position] != rune('f') {
											goto l45
										}
										position++
										if buffer[position] != rune('a') {
											goto l45
										}
										position++
										if buffer[position] != rune('l') {
											goto l45
										}
										position++
										if buffer[position] != rune('s') {
											goto l45
										}
										position++
										if buffer[position] != rune('e') {
											goto l45
										}
										position++
										{
											add(ruleAction11, position)
										}
										depth--
										add(ruleFalse, position128)
									}
								}
							l124:
								depth--
								add(ruleBoolean, position123)
							}
							break
						case '[':
							{
								position130 := position
								depth++
								if buffer[position] != rune('[') {
									goto l45
								}
								position++
								{
									add(ruleAction3, position)
								}
								{
									position132, tokenIndex132, depth132 := position, tokenIndex, depth
									{
										position134 := position
										depth++
										if !_rules[ruleListElem]() {
											goto l132
										}
									l135:
										{
											position136, tokenIndex136, depth136 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l136
											}
											position++
											{
												position137, tokenIndex137, depth137 := position, tokenIndex, depth
												if !_rules[ruleS]() {
													goto l137
												}
												goto l138
											l137:
												position, tokenIndex, depth = position137, tokenIndex137, depth137
											}
										l138:
											if !_rules[ruleListElem]() {
												goto l136
											}
											goto l135
										l136:
											position, tokenIndex, depth = position136, tokenIndex136, depth136
										}
										depth--
										add(ruleListElements, position134)
									}
									goto l133
								l132:
									position, tokenIndex, depth = position132, tokenIndex132, depth132
								}
							l133:
								if buffer[position] != rune(']') {
									goto l45
								}
								position++
								{
									add(ruleAction4, position)
								}
								depth--
								add(ruleList, position130)
							}
							break
						case '{':
							if !_rules[ruleDoc]() {
								goto l45
							}
							break
						default:
							{
								position140 := position
								depth++
								{
									position141 := position
									depth++
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l45
									}
									position++
								l142:
									{
										position143, tokenIndex143, depth143 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l143
										}
										position++
										goto l142
									l143:
										position, tokenIndex, depth = position143, tokenIndex143, depth143
									}
									{
										position144, tokenIndex144, depth144 := position, tokenIndex, depth
										if buffer[position] != rune('.') {
											goto l144
										}
										position++
										goto l145
									l144:
										position, tokenIndex, depth = position144, tokenIndex144, depth144
									}
								l145:
								l146:
									{
										position147, tokenIndex147, depth147 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l147
										}
										position++
										goto l146
									l147:
										position, tokenIndex, depth = position147, tokenIndex147, depth147
									}
									depth--
									add(rulePegText, position141)
								}
								{
									add(ruleAction7, position)
								}
								depth--
								add(ruleNumeric, position140)
							}
							break
						}
					}

				}
			l47:
				depth--
				add(ruleValue, position46)
			}
			return true
		l45:
			position, tokenIndex, depth = position45, tokenIndex45, depth45
			return false
		},
		/* 9 Numeric <- <(<([0-9]+ '.'? [0-9]*)> Action7)> */
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
		/* 16 ObjectID <- <('O' 'b' 'j' 'e' 'c' 't' 'I' 'd' '(' '"' <hexChar*> ('"' ')') Action13)> */
		nil,
		/* 17 BinData <- <('B' 'i' 'n' 'D' 'a' 't' 'a' '(' <(!')' .)+> ')' Action14)> */
		nil,
		/* 18 Regex <- <('/' <regexBody> Action15)> */
		nil,
		/* 19 Timestamp <- <('T' 'i' 'm' 'e' 's' 't' 'a' 'm' 'p' '(' <(!')' .)+> ')' Action16)> */
		nil,
		/* 20 NumberLong <- <('N' 'u' 'm' 'b' 'e' 'r' 'L' 'o' 'n' 'g' '(' <(!')' .)+> ')' Action17)> */
		nil,
		/* 21 MinKey <- <('M' 'i' 'n' 'K' 'e' 'y' Action18)> */
		nil,
		/* 22 MaxKey <- <('M' 'a' 'x' 'K' 'e' 'y' Action19)> */
		nil,
		/* 23 Undefined <- <('u' 'n' 'd' 'e' 'f' 'i' 'n' 'e' 'd' Action20)> */
		nil,
		/* 24 hexChar <- <([0-9] / ([a-f] / [A-F]))> */
		nil,
		/* 25 regexChar <- <(!'/' .)> */
		nil,
		/* 26 regexBody <- <(regexChar+ '/' ((&('s') 's') | (&('m') 'm') | (&('i') 'i') | (&('g') 'g'))*)> */
		nil,
		/* 27 stringChar <- <((!('"' / '\\') .) / ('\\' ('"' / '\\')))> */
		nil,
		/* 28 fieldChar <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 29 S <- <' '> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l169
				}
				position++
				depth--
				add(ruleS, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 31 Action0 <- <{ p.pushMap() }> */
		nil,
		/* 32 Action1 <- <{ p.popMap() }> */
		nil,
		/* 33 Action2 <- <{ p.setMapValue() }> */
		nil,
		/* 34 Action3 <- <{ p.pushList() }> */
		nil,
		/* 35 Action4 <- <{ p.popList() }> */
		nil,
		/* 36 Action5 <- <{ p.setListValue() }> */
		nil,
		nil,
		/* 38 Action6 <- <{ p.pushField(buffer[begin:end]) }> */
		nil,
		/* 39 Action7 <- <{ p.pushValue(numeric(buffer[begin:end])) }> */
		nil,
		/* 40 Action8 <- <{ p.pushValue(buffer[begin:end]) }> */
		nil,
		/* 41 Action9 <- <{ p.pushValue(nil) }> */
		nil,
		/* 42 Action10 <- <{ p.pushValue(true) }> */
		nil,
		/* 43 Action11 <- <{ p.pushValue(false) }> */
		nil,
		/* 44 Action12 <- <{ p.pushValue(date(buffer[begin:end])) }> */
		nil,
		/* 45 Action13 <- <{ p.pushValue(objectid(buffer[begin:end])) }> */
		nil,
		/* 46 Action14 <- <{ p.pushValue(bindata(buffer[begin:end])) }> */
		nil,
		/* 47 Action15 <- <{ p.pushValue(regex(buffer[begin:end])) }> */
		nil,
		/* 48 Action16 <- <{ p.pushValue(timestamp(buffer[begin:end])) }> */
		nil,
		/* 49 Action17 <- <{ p.pushValue(numberlong(buffer[begin:end])) }> */
		nil,
		/* 50 Action18 <- <{ p.pushValue(minkey()) }> */
		nil,
		/* 51 Action19 <- <{ p.pushValue(maxkey()) }> */
		nil,
		/* 52 Action20 <- <{ p.pushValue(undefined()) }> */
		nil,
	}
	p.rules = _rules
}
