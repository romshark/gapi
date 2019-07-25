package parser

import "fmt"

// parseStrField parses a struct field
func (pr *Parser) parseStrField(
	lex *Lexer,
	structType *TypeStruct,
) *StructField {
	// Read field name
	fName, err := readWord(
		lex,
		"field identifier",
		FragTkIdnFld,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new field instance
	newField := &StructField{
		Struct: structType,
		Name:   fName.src,
	}

	// Read type and set it when it's determined
	fType := pr.parseType(lex, func(t Type) {
		// Make sure the type of the field is pure
		if !t.IsPure() {
			pr.err(&pErr{
				at:      fName.begin,
				code:    ErrStructFieldImpure,
				message: fmt.Sprintf("struct field of impure type %s", t),
			})
		}

		newField.Type = t
	})
	if fType == nil {
		return nil
	}

	newField.Src = NewConstruct(lex, FragStrField,
		fName,
		fType,
	)

	// Define the graph node
	pr.defineGraphNode(newField)

	return newField
}
