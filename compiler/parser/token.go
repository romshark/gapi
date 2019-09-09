package parser

import (
	"fmt"

	parser "github.com/romshark/llparser"
)

// Token represents a typed source code token
type Token struct {
	kind  parser.FragmentKind
	src   string
	begin parser.Cursor
	end   parser.Cursor
}

// FragmentKind returns the token's fragment identifier
func (tok *Token) FragmentKind() parser.FragmentKind { return tok.kind }

// Begin returns the token's begin cursor
func (tok *Token) Begin() parser.Cursor { return tok.begin }

// End returns the token's end cursor
func (tok *Token) End() parser.Cursor { return tok.end }

// Src returns the token's raw source code
func (tok *Token) Src() string { return tok.src }

// Elements always returns nil for token fragments
func (tok *Token) Elements() []parser.Fragment { return nil }

// String returns the stringified token
func (tok *Token) String() string {
	return fmt.Sprintf(
		"%d(%d:%d-%d:%d)",
		tok.kind,
		tok.begin.Line,
		tok.begin.Column,
		tok.end.Line,
		tok.end.Column,
	)
}
