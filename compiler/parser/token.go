package parser

import "fmt"

// Token represents a typed source code token
type Token struct {
	id    FragID
	src   string
	begin Cursor
	end   Cursor
}

// FragID returns the token's fragment identifier
func (tok *Token) FragID() FragID { return tok.id }

// Begin returns the token's begin cursor
func (tok *Token) Begin() Cursor { return tok.begin }

// End returns the token's end cursor
func (tok *Token) End() Cursor { return tok.end }

// Src returns the token's raw source code
func (tok *Token) Src() string { return tok.src }

// Elements always returns nil for token fragments
func (tok *Token) Elements() []Fragment { return nil }

// String returns the stringified token
func (tok *Token) String() string {
	return fmt.Sprintf(
		"%s(%d:%d-%d:%d)",
		tok.id.String(),
		tok.begin.Line,
		tok.begin.Column,
		tok.end.Line,
		tok.end.Column,
	)
}
