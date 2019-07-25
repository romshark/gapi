package parser

import (
	"fmt"
)

func (pr *Parser) parseDeclMut(lex *Lexer) *Mutation {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordMutation,
		FragTkKwdMut,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read endpoint name
	fName, err := readWord(
		lex,
		"endpoint name",
		FragTkIdnProp,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new query endpoint instance
	newMutation := &Mutation{
		Name: fName.src,
	}

	// Parse parameters
	fParams, params, parsed := pr.parseOptParams(lex, newMutation)
	if !parsed {
		return nil
	}
	newMutation.Parameters = params

	// Read type ID
	fType, err := readWord(
		lex,
		"mutation type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	pr.deferJob(func() {
		// Ensure the type of the query endpoint exists
		resultType := pr.findTypeByName(fType.src)
		if resultType != nil {
			// Set the type
			newMutation.Type = resultType
			return
		}
		pr.err(&pErr{
			at:   fType.begin,
			code: ErrTypeUndef,
			message: fmt.Sprintf(
				"undefined type %s referenced by query endpoint %s",
				fType.src,
				fName.src,
			),
		})
	})

	newMutation.Src = NewConstruct(lex, FragDeclMut,
		fDeclKeyword,
		fName,
		fParams,
		fType,
	)

	// Define the endpoint
	pr.defineGraphNode(newMutation)

	return newMutation
}
