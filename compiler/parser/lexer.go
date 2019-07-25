package parser

import (
	"fmt"
)

// Lexer represents the source code lexer
type Lexer struct {
	tail Cursor
	src  *SourceFile
}

// NewLexer creates a new lexer instance
func NewLexer(sourceFile SourceFile) *Lexer {
	return &Lexer{
		src: &sourceFile,
		tail: Cursor{
			Index:  0,
			Line:   1,
			Column: 1,
			File:   &sourceFile.File,
		},
	}
}

// Cursor returns the current cursor position
func (lex *Lexer) Cursor() Cursor { return lex.tail }

// New creates a new lexer branching off the original one
func (lex *Lexer) New() *Lexer {
	return &Lexer{
		tail: lex.tail,
		src:  lex.src,
	}
}

// Peek returns true if expected is in front of the current position,
// otherwise returns false
func (lex *Lexer) Peek(expected string) bool {
	if lex.tail.Index+uint32(len(expected)) >= uint32(len(lex.src.Src)) {
		// The source is smaller than the expected target
		return false
	}
	e := 0
	i := lex.tail.Index + 1
	for i < uint32(len(expected)) {
		if lex.src.Src[i] != expected[e] {
			// The source doesn't match the expected target
			return false
		}
		i++
		e++
	}
	return true
}

func (lex *Lexer) readSpace() (*Token, Error) {
	begin := lex.tail
LOOP:
	for {
		if lex.tail.Index >= uint32(len(lex.src.Src)) {
			return lex.newToken(begin, FragTkSpace), nil
		}
		switch lex.src.Src[lex.tail.Index] {
		case ' ':
		case '\t':
		case '\n':
			// Unix line-break
			lex.tail.Index++
			lex.tail.Line++
			lex.tail.Column = 1
			continue LOOP
		case '\r':
			if lex.Peek("\n") {
				// \r\n line-break
				lex.tail.Index += 2
				lex.tail.Line++
				lex.tail.Column = 1
				continue LOOP
			}
			// Return carriage character without a following line-feed
			return nil, &pErr{
				at:   begin,
				code: ErrSyntax,
				message: fmt.Sprintf(
					"unexpected character '%d'",
					lex.src.Src[begin.Index],
				),
			}
		default:
			return lex.newToken(begin, FragTkSpace), nil
		}
		lex.tail.Index++
		lex.tail.Column++
	}
}

func (lex *Lexer) readLatinAlphanum() *Token {
	begin := lex.tail
	for {
		if lex.tail.Index >= uint32(len(lex.src.Src)) {
			return lex.newToken(begin, FragTkLatinAlphanum)
		}
		if !isLatinAlphanum(lex.src.Src, lex.tail.Index) {
			return lex.newToken(begin, FragTkLatinAlphanum)
		}
		lex.tail.Index++
		lex.tail.Column++
	}
}

func (lex *Lexer) tryReadSymList() (*Token, Error) {
	begin := lex.tail
	if lex.Peek("]") {
		lex.tail.Index += 2
		lex.tail.Column += 2
		return lex.newToken(begin, FragTkSymList), nil
	}
	return nil, &pErr{
		at:   begin,
		code: ErrSyntax,
		message: fmt.Sprintf(
			"unexpected character '%s'",
			string(lex.src.Src[begin.Index]),
		),
	}
}

func (lex *Lexer) newToken(begin Cursor, id FragID) *Token {
	if lex.tail.Index == begin.Index {
		// EOF
		return nil
	}
	newToken := &Token{
		id:    id,
		src:   lex.src.Src[begin.Index:lex.tail.Index],
		begin: begin,
		end:   lex.tail,
	}
	return newToken
}

// Next returns the next token or nil if there is EOF is reached
func (lex *Lexer) Next() (*Token, Error) {
	begin := lex.tail
	if lex.tail.Index >= uint32(len(lex.src.Src)) {
		// EOF
		return nil, nil
	}
	newSingleRuneTk := func(id FragID) *Token {
		lex.tail.Index++
		lex.tail.Column++
		newTk := lex.newToken(begin, id)
		return newTk
	}
	start := lex.src.Src[lex.tail.Index]
	switch start {
	case ' ':
		fallthrough
	case '\t':
		fallthrough
	case '\n':
		fallthrough
	case '\r':
		return lex.readSpace()
	case '{':
		return newSingleRuneTk(FragTkBlk), nil
	case '}':
		return newSingleRuneTk(FragTkBlkEnd), nil
	case '(':
		return newSingleRuneTk(FragTkPar), nil
	case ')':
		return newSingleRuneTk(FragTkParEnd), nil
	case ',':
		return newSingleRuneTk(FragTkSymSep), nil
	case '.':
		return newSingleRuneTk(FragTkMemAcc), nil
	case '#':
		return newSingleRuneTk(FragTkDocLineInit), nil
	case '=':
		return newSingleRuneTk(FragTkSymEq), nil
	case '?':
		return newSingleRuneTk(FragTkSymOpt), nil
	case '[':
		return lex.tryReadSymList()
	}
	if isLatinAlphanum(lex.src.Src, lex.tail.Index) {
		return lex.readLatinAlphanum(), nil
	}
	return nil, &pErr{
		at:      lex.tail,
		code:    ErrSyntax,
		message: fmt.Sprintf("unexpected symbol %d", start),
	}
}

// NextExpect returns an error if the next token isn't the expected one,
// otherwise returns the next token
func (lex *Lexer) NextExpect(
	expected FragID,
	msgFmt string,
	msgVars ...interface{},
) (*Token, Error) {
	tk, err := lex.Next()
	if err != nil {
		return nil, err
	}
	if tk == nil {
		return tk, &pErr{
			at:   lex.tail,
			code: ErrSyntax,
			message: fmt.Sprintf(msgFmt, msgVars...) +
				", reached end of file",
		}
	}
	if tk.id != expected {
		return tk, &pErr{
			at:   tk.begin,
			code: ErrSyntax,
			message: fmt.Sprintf(msgFmt, msgVars...) +
				fmt.Sprintf(", got: '%s'", tk.src),
		}
	}
	return tk, nil
}

// Skip helps defining skip-lists
type Skip []FragID

// NextExpectSkip returns an error if the next token isn't the expected one
// skipping all expected ignored fragments otherwise returns the next token
func (lex *Lexer) NextExpectSkip(
	expected FragID,
	ignore Skip,
	msgFmt string,
	msgVars ...interface{},
) (*Token, Error) {
SCAN_LOOP:
	for {
		tk, err := lex.Next()
		if err != nil {
			// Abort search on error
			return nil, err
		}
		if tk == nil {
			// Abort search on EOF
			return tk, &pErr{
				at:   lex.tail,
				code: ErrSyntax,
				message: fmt.Sprintf(msgFmt, msgVars...) +
					", reached end of file",
			}
		}
		if tk.id == expected {
			return tk, nil
		}
		for _, ignored := range ignore {
			if ignored == tk.id {
				// Ignore token and continue
				continue SCAN_LOOP
			}
		}
		return tk, &pErr{
			at:   tk.begin,
			code: ErrSyntax,
			message: fmt.Sprintf(msgFmt, msgVars...) +
				fmt.Sprintf(", got: '%s'", tk.src),
		}
	}
}

// NextSkip returns the next token skipping ignored fragments
func (lex *Lexer) NextSkip(ignore Skip) (*Token, Error) {
SCAN_LOOP:
	for {
		tk, err := lex.Next()
		if err != nil {
			// Abort search on error
			return nil, err
		}
		if tk == nil {
			// Abort search on EOF
			return nil, nil
		}
		for _, ignored := range ignore {
			if ignored == tk.id {
				// Ignore token and continue
				continue SCAN_LOOP
			}
		}
		return tk, nil
	}
}
