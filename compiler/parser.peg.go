package compiler

// Code generated by R:\go-workspace\bin\peg.exe ./parser.peg DO NOT EDIT

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleSchema
	ruleDclSc
	ruleDclAl
	ruleDclEn
	ruleDclSt
	ruleDclRv
	ruleDclTr
	ruleDclUn
	ruleDclQr
	ruleDclMt
	ruleDclSb
	ruleBlkEn
	ruleBlkSt
	ruleBlkRv
	ruleBlkUn
	rulePropSt
	rulePropRv
	ruleArgs
	ruleTp
	ruleTpNullable
	ruleTpList
	ruleTpName
	ruleWdLowCase
	ruleSpLn
	ruleSpOpt
	ruleSpLnOpt
	ruleSpVr
	ruleOEQ
	ruleKSC
	ruleKAL
	ruleKEN
	ruleKST
	ruleKRV
	ruleKTR
	ruleKUN
	ruleKQR
	ruleKMT
	ruleKSB
	rulePRN
	rulePRNE
	ruleBLK
	ruleBLKE
)

var rul3s = [...]string{
	"Unknown",
	"Schema",
	"DclSc",
	"DclAl",
	"DclEn",
	"DclSt",
	"DclRv",
	"DclTr",
	"DclUn",
	"DclQr",
	"DclMt",
	"DclSb",
	"BlkEn",
	"BlkSt",
	"BlkRv",
	"BlkUn",
	"PropSt",
	"PropRv",
	"Args",
	"Tp",
	"TpNullable",
	"TpList",
	"TpName",
	"WdLowCase",
	"SpLn",
	"SpOpt",
	"SpLnOpt",
	"SpVr",
	"OEQ",
	"KSC",
	"KAL",
	"KEN",
	"KST",
	"KRV",
	"KTR",
	"KUN",
	"KQR",
	"KMT",
	"KSB",
	"PRN",
	"PRNE",
	"BLK",
	"BLKE",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
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
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type GAPIParser struct {
	Buffer string
	buffer []rune
	rules  [43]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *GAPIParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *GAPIParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
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
	p   *GAPIParser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *GAPIParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *GAPIParser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func Pretty(pretty bool) func(*GAPIParser) error {
	return func(p *GAPIParser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*GAPIParser) error {
	return func(p *GAPIParser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *GAPIParser) Init(options ...func(*GAPIParser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
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
		/* 0 Schema <- <(SpOpt DclSc (SpOpt (DclAl / DclEn / DclSt / DclTr / DclRv / DclUn / DclQr / DclMt / DclSb))+ SpOpt)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleSpOpt]() {
					goto l0
				}
				if !_rules[ruleDclSc]() {
					goto l0
				}
				if !_rules[ruleSpOpt]() {
					goto l0
				}
				{
					position4, tokenIndex4 := position, tokenIndex
					if !_rules[ruleDclAl]() {
						goto l5
					}
					goto l4
				l5:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclEn]() {
						goto l6
					}
					goto l4
				l6:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclSt]() {
						goto l7
					}
					goto l4
				l7:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclTr]() {
						goto l8
					}
					goto l4
				l8:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclRv]() {
						goto l9
					}
					goto l4
				l9:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclUn]() {
						goto l10
					}
					goto l4
				l10:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclQr]() {
						goto l11
					}
					goto l4
				l11:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclMt]() {
						goto l12
					}
					goto l4
				l12:
					position, tokenIndex = position4, tokenIndex4
					if !_rules[ruleDclSb]() {
						goto l0
					}
				}
			l4:
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[ruleSpOpt]() {
						goto l3
					}
					{
						position13, tokenIndex13 := position, tokenIndex
						if !_rules[ruleDclAl]() {
							goto l14
						}
						goto l13
					l14:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclEn]() {
							goto l15
						}
						goto l13
					l15:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclSt]() {
							goto l16
						}
						goto l13
					l16:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclTr]() {
							goto l17
						}
						goto l13
					l17:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclRv]() {
							goto l18
						}
						goto l13
					l18:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclUn]() {
							goto l19
						}
						goto l13
					l19:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclQr]() {
							goto l20
						}
						goto l13
					l20:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclMt]() {
							goto l21
						}
						goto l13
					l21:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruleDclSb]() {
							goto l3
						}
					}
				l13:
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleSpOpt]() {
					goto l0
				}
				add(ruleSchema, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 DclSc <- <(KSC SpLn WdLowCase)> */
		func() bool {
			position22, tokenIndex22 := position, tokenIndex
			{
				position23 := position
				if !_rules[ruleKSC]() {
					goto l22
				}
				if !_rules[ruleSpLn]() {
					goto l22
				}
				if !_rules[ruleWdLowCase]() {
					goto l22
				}
				add(ruleDclSc, position23)
			}
			return true
		l22:
			position, tokenIndex = position22, tokenIndex22
			return false
		},
		/* 2 DclAl <- <(KAL SpLn TpName SpLn OEQ SpLn TpName)> */
		func() bool {
			position24, tokenIndex24 := position, tokenIndex
			{
				position25 := position
				if !_rules[ruleKAL]() {
					goto l24
				}
				if !_rules[ruleSpLn]() {
					goto l24
				}
				if !_rules[ruleTpName]() {
					goto l24
				}
				if !_rules[ruleSpLn]() {
					goto l24
				}
				if !_rules[ruleOEQ]() {
					goto l24
				}
				if !_rules[ruleSpLn]() {
					goto l24
				}
				if !_rules[ruleTpName]() {
					goto l24
				}
				add(ruleDclAl, position25)
			}
			return true
		l24:
			position, tokenIndex = position24, tokenIndex24
			return false
		},
		/* 3 DclEn <- <(KEN SpLn TpName SpLnOpt BlkEn)> */
		func() bool {
			position26, tokenIndex26 := position, tokenIndex
			{
				position27 := position
				if !_rules[ruleKEN]() {
					goto l26
				}
				if !_rules[ruleSpLn]() {
					goto l26
				}
				if !_rules[ruleTpName]() {
					goto l26
				}
				if !_rules[ruleSpLnOpt]() {
					goto l26
				}
				if !_rules[ruleBlkEn]() {
					goto l26
				}
				add(ruleDclEn, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 4 DclSt <- <(KST SpLn TpName SpLnOpt BlkSt)> */
		func() bool {
			position28, tokenIndex28 := position, tokenIndex
			{
				position29 := position
				if !_rules[ruleKST]() {
					goto l28
				}
				if !_rules[ruleSpLn]() {
					goto l28
				}
				if !_rules[ruleTpName]() {
					goto l28
				}
				if !_rules[ruleSpLnOpt]() {
					goto l28
				}
				if !_rules[ruleBlkSt]() {
					goto l28
				}
				add(ruleDclSt, position29)
			}
			return true
		l28:
			position, tokenIndex = position28, tokenIndex28
			return false
		},
		/* 5 DclRv <- <(KRV SpLn TpName SpLnOpt BlkRv)> */
		func() bool {
			position30, tokenIndex30 := position, tokenIndex
			{
				position31 := position
				if !_rules[ruleKRV]() {
					goto l30
				}
				if !_rules[ruleSpLn]() {
					goto l30
				}
				if !_rules[ruleTpName]() {
					goto l30
				}
				if !_rules[ruleSpLnOpt]() {
					goto l30
				}
				if !_rules[ruleBlkRv]() {
					goto l30
				}
				add(ruleDclRv, position31)
			}
			return true
		l30:
			position, tokenIndex = position30, tokenIndex30
			return false
		},
		/* 6 DclTr <- <(KTR SpLn TpName SpLnOpt BlkRv)> */
		func() bool {
			position32, tokenIndex32 := position, tokenIndex
			{
				position33 := position
				if !_rules[ruleKTR]() {
					goto l32
				}
				if !_rules[ruleSpLn]() {
					goto l32
				}
				if !_rules[ruleTpName]() {
					goto l32
				}
				if !_rules[ruleSpLnOpt]() {
					goto l32
				}
				if !_rules[ruleBlkRv]() {
					goto l32
				}
				add(ruleDclTr, position33)
			}
			return true
		l32:
			position, tokenIndex = position32, tokenIndex32
			return false
		},
		/* 7 DclUn <- <(KUN SpLn TpName SpLnOpt BlkUn)> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				if !_rules[ruleKUN]() {
					goto l34
				}
				if !_rules[ruleSpLn]() {
					goto l34
				}
				if !_rules[ruleTpName]() {
					goto l34
				}
				if !_rules[ruleSpLnOpt]() {
					goto l34
				}
				if !_rules[ruleBlkUn]() {
					goto l34
				}
				add(ruleDclUn, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 8 DclQr <- <(KQR SpLn WdLowCase SpLnOpt Args? SpLnOpt Tp)> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				if !_rules[ruleKQR]() {
					goto l36
				}
				if !_rules[ruleSpLn]() {
					goto l36
				}
				if !_rules[ruleWdLowCase]() {
					goto l36
				}
				if !_rules[ruleSpLnOpt]() {
					goto l36
				}
				{
					position38, tokenIndex38 := position, tokenIndex
					if !_rules[ruleArgs]() {
						goto l38
					}
					goto l39
				l38:
					position, tokenIndex = position38, tokenIndex38
				}
			l39:
				if !_rules[ruleSpLnOpt]() {
					goto l36
				}
				if !_rules[ruleTp]() {
					goto l36
				}
				add(ruleDclQr, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 9 DclMt <- <(KMT SpLn WdLowCase SpLnOpt Args? SpLnOpt Tp)> */
		func() bool {
			position40, tokenIndex40 := position, tokenIndex
			{
				position41 := position
				if !_rules[ruleKMT]() {
					goto l40
				}
				if !_rules[ruleSpLn]() {
					goto l40
				}
				if !_rules[ruleWdLowCase]() {
					goto l40
				}
				if !_rules[ruleSpLnOpt]() {
					goto l40
				}
				{
					position42, tokenIndex42 := position, tokenIndex
					if !_rules[ruleArgs]() {
						goto l42
					}
					goto l43
				l42:
					position, tokenIndex = position42, tokenIndex42
				}
			l43:
				if !_rules[ruleSpLnOpt]() {
					goto l40
				}
				if !_rules[ruleTp]() {
					goto l40
				}
				add(ruleDclMt, position41)
			}
			return true
		l40:
			position, tokenIndex = position40, tokenIndex40
			return false
		},
		/* 10 DclSb <- <(KSB SpLn WdLowCase SpLnOpt Args? SpLnOpt Tp)> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				if !_rules[ruleKSB]() {
					goto l44
				}
				if !_rules[ruleSpLn]() {
					goto l44
				}
				if !_rules[ruleWdLowCase]() {
					goto l44
				}
				if !_rules[ruleSpLnOpt]() {
					goto l44
				}
				{
					position46, tokenIndex46 := position, tokenIndex
					if !_rules[ruleArgs]() {
						goto l46
					}
					goto l47
				l46:
					position, tokenIndex = position46, tokenIndex46
				}
			l47:
				if !_rules[ruleSpLnOpt]() {
					goto l44
				}
				if !_rules[ruleTp]() {
					goto l44
				}
				add(ruleDclSb, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 11 BlkEn <- <(BLK (SpOpt WdLowCase)+ SpOpt BLKE)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if !_rules[ruleBLK]() {
					goto l48
				}
				if !_rules[ruleSpOpt]() {
					goto l48
				}
				if !_rules[ruleWdLowCase]() {
					goto l48
				}
			l50:
				{
					position51, tokenIndex51 := position, tokenIndex
					if !_rules[ruleSpOpt]() {
						goto l51
					}
					if !_rules[ruleWdLowCase]() {
						goto l51
					}
					goto l50
				l51:
					position, tokenIndex = position51, tokenIndex51
				}
				if !_rules[ruleSpOpt]() {
					goto l48
				}
				if !_rules[ruleBLKE]() {
					goto l48
				}
				add(ruleBlkEn, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 12 BlkSt <- <(BLK SpOpt PropSt+ SpOpt BLKE)> */
		func() bool {
			position52, tokenIndex52 := position, tokenIndex
			{
				position53 := position
				if !_rules[ruleBLK]() {
					goto l52
				}
				if !_rules[ruleSpOpt]() {
					goto l52
				}
				if !_rules[rulePropSt]() {
					goto l52
				}
			l54:
				{
					position55, tokenIndex55 := position, tokenIndex
					if !_rules[rulePropSt]() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex = position55, tokenIndex55
				}
				if !_rules[ruleSpOpt]() {
					goto l52
				}
				if !_rules[ruleBLKE]() {
					goto l52
				}
				add(ruleBlkSt, position53)
			}
			return true
		l52:
			position, tokenIndex = position52, tokenIndex52
			return false
		},
		/* 13 BlkRv <- <(BLK SpOpt PropRv+ SpOpt BLKE)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				if !_rules[ruleBLK]() {
					goto l56
				}
				if !_rules[ruleSpOpt]() {
					goto l56
				}
				if !_rules[rulePropRv]() {
					goto l56
				}
			l58:
				{
					position59, tokenIndex59 := position, tokenIndex
					if !_rules[rulePropRv]() {
						goto l59
					}
					goto l58
				l59:
					position, tokenIndex = position59, tokenIndex59
				}
				if !_rules[ruleSpOpt]() {
					goto l56
				}
				if !_rules[ruleBLKE]() {
					goto l56
				}
				add(ruleBlkRv, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 14 BlkUn <- <(BLK (SpOpt TpName)+ SpOpt BLKE)> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				if !_rules[ruleBLK]() {
					goto l60
				}
				if !_rules[ruleSpOpt]() {
					goto l60
				}
				if !_rules[ruleTpName]() {
					goto l60
				}
			l62:
				{
					position63, tokenIndex63 := position, tokenIndex
					if !_rules[ruleSpOpt]() {
						goto l63
					}
					if !_rules[ruleTpName]() {
						goto l63
					}
					goto l62
				l63:
					position, tokenIndex = position63, tokenIndex63
				}
				if !_rules[ruleSpOpt]() {
					goto l60
				}
				if !_rules[ruleBLKE]() {
					goto l60
				}
				add(ruleBlkUn, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 15 PropSt <- <(WdLowCase SpLn Tp)> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				if !_rules[ruleWdLowCase]() {
					goto l64
				}
				if !_rules[ruleSpLn]() {
					goto l64
				}
				if !_rules[ruleTp]() {
					goto l64
				}
				add(rulePropSt, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 16 PropRv <- <(WdLowCase SpLn Args? Tp)> */
		func() bool {
			position66, tokenIndex66 := position, tokenIndex
			{
				position67 := position
				if !_rules[ruleWdLowCase]() {
					goto l66
				}
				if !_rules[ruleSpLn]() {
					goto l66
				}
				{
					position68, tokenIndex68 := position, tokenIndex
					if !_rules[ruleArgs]() {
						goto l68
					}
					goto l69
				l68:
					position, tokenIndex = position68, tokenIndex68
				}
			l69:
				if !_rules[ruleTp]() {
					goto l66
				}
				add(rulePropRv, position67)
			}
			return true
		l66:
			position, tokenIndex = position66, tokenIndex66
			return false
		},
		/* 17 Args <- <(PRN SpOpt (WdLowCase Tp ',')+ SpOpt PRNE)> */
		func() bool {
			position70, tokenIndex70 := position, tokenIndex
			{
				position71 := position
				if !_rules[rulePRN]() {
					goto l70
				}
				if !_rules[ruleSpOpt]() {
					goto l70
				}
				if !_rules[ruleWdLowCase]() {
					goto l70
				}
				if !_rules[ruleTp]() {
					goto l70
				}
				if buffer[position] != rune(',') {
					goto l70
				}
				position++
			l72:
				{
					position73, tokenIndex73 := position, tokenIndex
					if !_rules[ruleWdLowCase]() {
						goto l73
					}
					if !_rules[ruleTp]() {
						goto l73
					}
					if buffer[position] != rune(',') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex = position73, tokenIndex73
				}
				if !_rules[ruleSpOpt]() {
					goto l70
				}
				if !_rules[rulePRNE]() {
					goto l70
				}
				add(ruleArgs, position71)
			}
			return true
		l70:
			position, tokenIndex = position70, tokenIndex70
			return false
		},
		/* 18 Tp <- <(TpName / TpNullable / TpList)> */
		func() bool {
			position74, tokenIndex74 := position, tokenIndex
			{
				position75 := position
				{
					position76, tokenIndex76 := position, tokenIndex
					if !_rules[ruleTpName]() {
						goto l77
					}
					goto l76
				l77:
					position, tokenIndex = position76, tokenIndex76
					if !_rules[ruleTpNullable]() {
						goto l78
					}
					goto l76
				l78:
					position, tokenIndex = position76, tokenIndex76
					if !_rules[ruleTpList]() {
						goto l74
					}
				}
			l76:
				add(ruleTp, position75)
			}
			return true
		l74:
			position, tokenIndex = position74, tokenIndex74
			return false
		},
		/* 19 TpNullable <- <('?' Tp)> */
		func() bool {
			position79, tokenIndex79 := position, tokenIndex
			{
				position80 := position
				if buffer[position] != rune('?') {
					goto l79
				}
				position++
				if !_rules[ruleTp]() {
					goto l79
				}
				add(ruleTpNullable, position80)
			}
			return true
		l79:
			position, tokenIndex = position79, tokenIndex79
			return false
		},
		/* 20 TpList <- <('[' ']' Tp)> */
		func() bool {
			position81, tokenIndex81 := position, tokenIndex
			{
				position82 := position
				if buffer[position] != rune('[') {
					goto l81
				}
				position++
				if buffer[position] != rune(']') {
					goto l81
				}
				position++
				if !_rules[ruleTp]() {
					goto l81
				}
				add(ruleTpList, position82)
			}
			return true
		l81:
			position, tokenIndex = position81, tokenIndex81
			return false
		},
		/* 21 TpName <- <([A-Z] ([a-z] / [A-Z] / [0-9])*)> */
		func() bool {
			position83, tokenIndex83 := position, tokenIndex
			{
				position84 := position
				if c := buffer[position]; c < rune('A') || c > rune('Z') {
					goto l83
				}
				position++
			l85:
				{
					position86, tokenIndex86 := position, tokenIndex
					{
						position87, tokenIndex87 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l88
						}
						position++
						goto l87
					l88:
						position, tokenIndex = position87, tokenIndex87
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l89
						}
						position++
						goto l87
					l89:
						position, tokenIndex = position87, tokenIndex87
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l86
						}
						position++
					}
				l87:
					goto l85
				l86:
					position, tokenIndex = position86, tokenIndex86
				}
				add(ruleTpName, position84)
			}
			return true
		l83:
			position, tokenIndex = position83, tokenIndex83
			return false
		},
		/* 22 WdLowCase <- <([a-z] ([a-z] / [A-Z] / [0-9])*)> */
		func() bool {
			position90, tokenIndex90 := position, tokenIndex
			{
				position91 := position
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l90
				}
				position++
			l92:
				{
					position93, tokenIndex93 := position, tokenIndex
					{
						position94, tokenIndex94 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex = position94, tokenIndex94
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l96
						}
						position++
						goto l94
					l96:
						position, tokenIndex = position94, tokenIndex94
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l93
						}
						position++
					}
				l94:
					goto l92
				l93:
					position, tokenIndex = position93, tokenIndex93
				}
				add(ruleWdLowCase, position91)
			}
			return true
		l90:
			position, tokenIndex = position90, tokenIndex90
			return false
		},
		/* 23 SpLn <- <(' ' / '\t')+> */
		func() bool {
			position97, tokenIndex97 := position, tokenIndex
			{
				position98 := position
				{
					position101, tokenIndex101 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex = position101, tokenIndex101
					if buffer[position] != rune('\t') {
						goto l97
					}
					position++
				}
			l101:
			l99:
				{
					position100, tokenIndex100 := position, tokenIndex
					{
						position103, tokenIndex103 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l104
						}
						position++
						goto l103
					l104:
						position, tokenIndex = position103, tokenIndex103
						if buffer[position] != rune('\t') {
							goto l100
						}
						position++
					}
				l103:
					goto l99
				l100:
					position, tokenIndex = position100, tokenIndex100
				}
				add(ruleSpLn, position98)
			}
			return true
		l97:
			position, tokenIndex = position97, tokenIndex97
			return false
		},
		/* 24 SpOpt <- <(SpLn / SpVr)*> */
		func() bool {
			{
				position106 := position
			l107:
				{
					position108, tokenIndex108 := position, tokenIndex
					{
						position109, tokenIndex109 := position, tokenIndex
						if !_rules[ruleSpLn]() {
							goto l110
						}
						goto l109
					l110:
						position, tokenIndex = position109, tokenIndex109
						if !_rules[ruleSpVr]() {
							goto l108
						}
					}
				l109:
					goto l107
				l108:
					position, tokenIndex = position108, tokenIndex108
				}
				add(ruleSpOpt, position106)
			}
			return true
		},
		/* 25 SpLnOpt <- <(' ' / '\t')*> */
		func() bool {
			{
				position112 := position
			l113:
				{
					position114, tokenIndex114 := position, tokenIndex
					{
						position115, tokenIndex115 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l116
						}
						position++
						goto l115
					l116:
						position, tokenIndex = position115, tokenIndex115
						if buffer[position] != rune('\t') {
							goto l114
						}
						position++
					}
				l115:
					goto l113
				l114:
					position, tokenIndex = position114, tokenIndex114
				}
				add(ruleSpLnOpt, position112)
			}
			return true
		},
		/* 26 SpVr <- <('\n' / ('\r' '\n'))+> */
		func() bool {
			position117, tokenIndex117 := position, tokenIndex
			{
				position118 := position
				{
					position121, tokenIndex121 := position, tokenIndex
					if buffer[position] != rune('\n') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex = position121, tokenIndex121
					if buffer[position] != rune('\r') {
						goto l117
					}
					position++
					if buffer[position] != rune('\n') {
						goto l117
					}
					position++
				}
			l121:
			l119:
				{
					position120, tokenIndex120 := position, tokenIndex
					{
						position123, tokenIndex123 := position, tokenIndex
						if buffer[position] != rune('\n') {
							goto l124
						}
						position++
						goto l123
					l124:
						position, tokenIndex = position123, tokenIndex123
						if buffer[position] != rune('\r') {
							goto l120
						}
						position++
						if buffer[position] != rune('\n') {
							goto l120
						}
						position++
					}
				l123:
					goto l119
				l120:
					position, tokenIndex = position120, tokenIndex120
				}
				add(ruleSpVr, position118)
			}
			return true
		l117:
			position, tokenIndex = position117, tokenIndex117
			return false
		},
		/* 27 OEQ <- <'='> */
		func() bool {
			position125, tokenIndex125 := position, tokenIndex
			{
				position126 := position
				if buffer[position] != rune('=') {
					goto l125
				}
				position++
				add(ruleOEQ, position126)
			}
			return true
		l125:
			position, tokenIndex = position125, tokenIndex125
			return false
		},
		/* 28 KSC <- <('s' 'c' 'h' 'e' 'm' 'a')> */
		func() bool {
			position127, tokenIndex127 := position, tokenIndex
			{
				position128 := position
				if buffer[position] != rune('s') {
					goto l127
				}
				position++
				if buffer[position] != rune('c') {
					goto l127
				}
				position++
				if buffer[position] != rune('h') {
					goto l127
				}
				position++
				if buffer[position] != rune('e') {
					goto l127
				}
				position++
				if buffer[position] != rune('m') {
					goto l127
				}
				position++
				if buffer[position] != rune('a') {
					goto l127
				}
				position++
				add(ruleKSC, position128)
			}
			return true
		l127:
			position, tokenIndex = position127, tokenIndex127
			return false
		},
		/* 29 KAL <- <('a' 'l' 'i' 'a' 's')> */
		func() bool {
			position129, tokenIndex129 := position, tokenIndex
			{
				position130 := position
				if buffer[position] != rune('a') {
					goto l129
				}
				position++
				if buffer[position] != rune('l') {
					goto l129
				}
				position++
				if buffer[position] != rune('i') {
					goto l129
				}
				position++
				if buffer[position] != rune('a') {
					goto l129
				}
				position++
				if buffer[position] != rune('s') {
					goto l129
				}
				position++
				add(ruleKAL, position130)
			}
			return true
		l129:
			position, tokenIndex = position129, tokenIndex129
			return false
		},
		/* 30 KEN <- <('e' 'n' 'u' 'm')> */
		func() bool {
			position131, tokenIndex131 := position, tokenIndex
			{
				position132 := position
				if buffer[position] != rune('e') {
					goto l131
				}
				position++
				if buffer[position] != rune('n') {
					goto l131
				}
				position++
				if buffer[position] != rune('u') {
					goto l131
				}
				position++
				if buffer[position] != rune('m') {
					goto l131
				}
				position++
				add(ruleKEN, position132)
			}
			return true
		l131:
			position, tokenIndex = position131, tokenIndex131
			return false
		},
		/* 31 KST <- <('s' 't' 'r' 'u' 'c' 't')> */
		func() bool {
			position133, tokenIndex133 := position, tokenIndex
			{
				position134 := position
				if buffer[position] != rune('s') {
					goto l133
				}
				position++
				if buffer[position] != rune('t') {
					goto l133
				}
				position++
				if buffer[position] != rune('r') {
					goto l133
				}
				position++
				if buffer[position] != rune('u') {
					goto l133
				}
				position++
				if buffer[position] != rune('c') {
					goto l133
				}
				position++
				if buffer[position] != rune('t') {
					goto l133
				}
				position++
				add(ruleKST, position134)
			}
			return true
		l133:
			position, tokenIndex = position133, tokenIndex133
			return false
		},
		/* 32 KRV <- <('r' 'e' 's' 'o' 'l' 'v' 'e' 'r')> */
		func() bool {
			position135, tokenIndex135 := position, tokenIndex
			{
				position136 := position
				if buffer[position] != rune('r') {
					goto l135
				}
				position++
				if buffer[position] != rune('e') {
					goto l135
				}
				position++
				if buffer[position] != rune('s') {
					goto l135
				}
				position++
				if buffer[position] != rune('o') {
					goto l135
				}
				position++
				if buffer[position] != rune('l') {
					goto l135
				}
				position++
				if buffer[position] != rune('v') {
					goto l135
				}
				position++
				if buffer[position] != rune('e') {
					goto l135
				}
				position++
				if buffer[position] != rune('r') {
					goto l135
				}
				position++
				add(ruleKRV, position136)
			}
			return true
		l135:
			position, tokenIndex = position135, tokenIndex135
			return false
		},
		/* 33 KTR <- <('t' 'r' 'a' 'i' 't')> */
		func() bool {
			position137, tokenIndex137 := position, tokenIndex
			{
				position138 := position
				if buffer[position] != rune('t') {
					goto l137
				}
				position++
				if buffer[position] != rune('r') {
					goto l137
				}
				position++
				if buffer[position] != rune('a') {
					goto l137
				}
				position++
				if buffer[position] != rune('i') {
					goto l137
				}
				position++
				if buffer[position] != rune('t') {
					goto l137
				}
				position++
				add(ruleKTR, position138)
			}
			return true
		l137:
			position, tokenIndex = position137, tokenIndex137
			return false
		},
		/* 34 KUN <- <('u' 'n' 'i' 'o' 'n')> */
		func() bool {
			position139, tokenIndex139 := position, tokenIndex
			{
				position140 := position
				if buffer[position] != rune('u') {
					goto l139
				}
				position++
				if buffer[position] != rune('n') {
					goto l139
				}
				position++
				if buffer[position] != rune('i') {
					goto l139
				}
				position++
				if buffer[position] != rune('o') {
					goto l139
				}
				position++
				if buffer[position] != rune('n') {
					goto l139
				}
				position++
				add(ruleKUN, position140)
			}
			return true
		l139:
			position, tokenIndex = position139, tokenIndex139
			return false
		},
		/* 35 KQR <- <('q' 'u' 'e' 'r' 'y')> */
		func() bool {
			position141, tokenIndex141 := position, tokenIndex
			{
				position142 := position
				if buffer[position] != rune('q') {
					goto l141
				}
				position++
				if buffer[position] != rune('u') {
					goto l141
				}
				position++
				if buffer[position] != rune('e') {
					goto l141
				}
				position++
				if buffer[position] != rune('r') {
					goto l141
				}
				position++
				if buffer[position] != rune('y') {
					goto l141
				}
				position++
				add(ruleKQR, position142)
			}
			return true
		l141:
			position, tokenIndex = position141, tokenIndex141
			return false
		},
		/* 36 KMT <- <('m' 'u' 't' 'a' 't' 'i' 'o' 'n')> */
		func() bool {
			position143, tokenIndex143 := position, tokenIndex
			{
				position144 := position
				if buffer[position] != rune('m') {
					goto l143
				}
				position++
				if buffer[position] != rune('u') {
					goto l143
				}
				position++
				if buffer[position] != rune('t') {
					goto l143
				}
				position++
				if buffer[position] != rune('a') {
					goto l143
				}
				position++
				if buffer[position] != rune('t') {
					goto l143
				}
				position++
				if buffer[position] != rune('i') {
					goto l143
				}
				position++
				if buffer[position] != rune('o') {
					goto l143
				}
				position++
				if buffer[position] != rune('n') {
					goto l143
				}
				position++
				add(ruleKMT, position144)
			}
			return true
		l143:
			position, tokenIndex = position143, tokenIndex143
			return false
		},
		/* 37 KSB <- <('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n')> */
		func() bool {
			position145, tokenIndex145 := position, tokenIndex
			{
				position146 := position
				if buffer[position] != rune('s') {
					goto l145
				}
				position++
				if buffer[position] != rune('u') {
					goto l145
				}
				position++
				if buffer[position] != rune('b') {
					goto l145
				}
				position++
				if buffer[position] != rune('s') {
					goto l145
				}
				position++
				if buffer[position] != rune('c') {
					goto l145
				}
				position++
				if buffer[position] != rune('r') {
					goto l145
				}
				position++
				if buffer[position] != rune('i') {
					goto l145
				}
				position++
				if buffer[position] != rune('p') {
					goto l145
				}
				position++
				if buffer[position] != rune('t') {
					goto l145
				}
				position++
				if buffer[position] != rune('i') {
					goto l145
				}
				position++
				if buffer[position] != rune('o') {
					goto l145
				}
				position++
				if buffer[position] != rune('n') {
					goto l145
				}
				position++
				add(ruleKSB, position146)
			}
			return true
		l145:
			position, tokenIndex = position145, tokenIndex145
			return false
		},
		/* 38 PRN <- <'('> */
		func() bool {
			position147, tokenIndex147 := position, tokenIndex
			{
				position148 := position
				if buffer[position] != rune('(') {
					goto l147
				}
				position++
				add(rulePRN, position148)
			}
			return true
		l147:
			position, tokenIndex = position147, tokenIndex147
			return false
		},
		/* 39 PRNE <- <')'> */
		func() bool {
			position149, tokenIndex149 := position, tokenIndex
			{
				position150 := position
				if buffer[position] != rune(')') {
					goto l149
				}
				position++
				add(rulePRNE, position150)
			}
			return true
		l149:
			position, tokenIndex = position149, tokenIndex149
			return false
		},
		/* 40 BLK <- <'{'> */
		func() bool {
			position151, tokenIndex151 := position, tokenIndex
			{
				position152 := position
				if buffer[position] != rune('{') {
					goto l151
				}
				position++
				add(ruleBLK, position152)
			}
			return true
		l151:
			position, tokenIndex = position151, tokenIndex151
			return false
		},
		/* 41 BLKE <- <'}'> */
		func() bool {
			position153, tokenIndex153 := position, tokenIndex
			{
				position154 := position
				if buffer[position] != rune('}') {
					goto l153
				}
				position++
				add(ruleBLKE, position154)
			}
			return true
		l153:
			position, tokenIndex = position153, tokenIndex153
			return false
		},
	}
	p.rules = _rules
	return nil
}
