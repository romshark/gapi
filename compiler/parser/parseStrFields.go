package parser

import "fmt"

// parseStrFields parses the fields block of a struct declaration
func (pr *Parser) parseStrFields(
	lex *Lexer,
	structType *TypeStruct,
) (Fragment, []*StructField) {
	// Read '{'
	fBlockBegin, err := readToken(
		lex,
		FragTkBlk,
		"struct fields block opening '{'",
	)
	if pr.err(err) {
		return nil, nil
	}

	frags := []Fragment{fBlockBegin}
	byName := map[string]*Token{}
	fields := []*StructField{}

	// Parse fields
SCAN_LOOP:
	for {
		peeker := lex.New()
		// Peek for 1 token to find out whether
		// the block ended or a new field began
		tk, err := peeker.NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil, nil
		}
		if tk == nil {
			// Unexpected EOF
			pr.err(&pErr{
				at:      peeker.Cursor(),
				code:    ErrSyntax,
				message: "unexpected end of file",
			})
			return nil, nil
		}

		switch tk.id {
		case FragTkLatinAlphanum:
			// A field
			newField := pr.parseStrField(lex, structType)
			if newField == nil {
				return nil, nil
			}
			frags = append(frags, newField.Src)
			fields = append(fields, newField)
		case FragTkBlkEnd:
			// End of the block
			frags = append(frags, tk)
			_, _ = lex.NextSkip(Skip{FragTkSpace})
			break SCAN_LOOP
		default:
			// Unexpected token
			pr.err(&pErr{
				at:      tk.begin,
				code:    ErrSyntax,
				message: fmt.Sprintf("unexpected token '%s'", tk.src),
			})
			return nil, nil
		}

		fieldName := tk.src

		// Check for redeclarations
		if defined, isDefined := byName[fieldName]; isDefined {
			pr.err(&pErr{
				at:   tk.begin,
				code: ErrStructFieldRedecl,
				message: fmt.Sprintf(
					"Redeclaration of struct field %s "+
						"(previously declared at %s)",
					fieldName,
					defined.begin,
				),
			})
			return nil, nil
		}

		byName[fieldName] = tk
	}

	// Make sure there's at least 1 field
	if len(fields) < 1 {
		pr.err(&pErr{
			at:   fBlockBegin.begin,
			code: ErrStructNoFields,
			message: fmt.Sprintf(
				"struct %s is missing fields",
				structType.TypeName,
			),
		})
		return nil, nil
	}

	return NewConstruct(lex, FragStrFields, frags...), fields
}

/*
	valid = true
	nd := first
	for nd != nil {
		var newField *StructField
		var ndFldTp *Lexer
		ndField := nd

		// Read field name
		nd = ndField.up
		fieldName := c.getSrc(nd)

		// Verify field identifier
		if err := verifyLowerCamelCase(fieldName); err != nil {
			c.err(pErr{
				ErrStructFieldIllegalIdent,
				fmt.Sprintf(
					"invalid struct field identifier at %d:%d: %s",
					nd.begin,
					nd.end,
					err,
				),
			})
			valid = false
			goto NEXT
		}

		// Check for redeclared fields
		if field := structType.FieldByName(fieldName); field != nil {
			c.err(pErr{
				ErrStructFieldRedecl,
				fmt.Sprintf(
					"Redeclaration of struct field %s at %d:%d "+
						"(previously declared at %d:%d)",
					fieldName,
					nd.begin,
					nd.end,
					field.Begin,
					field.End,
				),
			})
			valid = false
			goto NEXT
		}

		// Add field
		newField = &StructField{
			Src: Src{
				Begin: nd.begin,
				End:   nd.end,
			},
			Struct:  structType,
			Name:    fieldName,
			GraphID: 0,   // Set during definition
			Type:    nil, // Deferred
		}
		c.defineGraphNode(newField)
		structType.Fields = append(structType.Fields, newField)

		nd = skipUntil(nd.next, ruleTp)
		ndFldTp = nd

		// Parse the field type in deferred mode
		c.deferJob(func() error {
			fieldType, err := c.parseType(ndFldTp)
			if err != nil {
				c.err(err)
				return nil
			}

			// Ensure all struct fields are of a pure type
			if !fieldType.IsPure() {
				c.err(pErr{
					ErrStructFieldImpure,
					fmt.Sprintf(
						"Struct field %s has impure type %s at %d:%d",
						newField.GraphNodeName(),
						fieldType,
						ndField.begin,
						ndField.end,
					),
				})
			}

			// Set the field type
			newField.Type = fieldType

			return nil
		})

	NEXT:
		nd = skipUntil(ndField.next, ruleStFld)
	}
	return
*/
