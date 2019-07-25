package parser

// Construct represents a non-terminal fragment (construct)
type Construct struct {
	id       FragID
	src      string
	begin    Cursor
	end      Cursor
	elements []Fragment
}

// NewConstruct creates a new construct
func NewConstruct(
	lexer *Lexer,
	id FragID,
	elements ...Fragment,
) *Construct {
	begin := elements[0].Begin()
	end := elements[len(elements)-1].End()
	return &Construct{
		id:       id,
		src:      lexer.src.Src[begin.Index:end.Index],
		begin:    begin,
		end:      end,
		elements: elements,
	}
}

// FragID returns the token's fragment identifier
func (con *Construct) FragID() FragID { return con.id }

// Begin returns the token's begin cursor
func (con *Construct) Begin() Cursor { return con.begin }

// End returns the token's end cursor
func (con *Construct) End() Cursor { return con.end }

// Src returns the token's raw source code
func (con *Construct) Src() string { return con.src }

// Elements always returns all sub-fragments
func (con *Construct) Elements() []Fragment { return con.elements }
