package parser

import (
	"fmt"
)

func (pr *Parser) parseDeclQry(lex *Lexer) *Query {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordQuery,
		FragTkKwdQry,
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
	newQuery := &Query{
		Name: fName.src,
	}

	// Parse parameters
	fParams, params, parsed := pr.parseOptParams(lex, newQuery)
	if !parsed {
		return nil
	}
	newQuery.Parameters = params

	// Read type ID
	fType, err := readWord(
		lex,
		"query type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	pr.deferJob(func() {
		// Ensure the type of the query endpoint exists
		resultType := pr.findTypeByDesignation(fType.src)
		if resultType != nil {
			// Set the type
			newQuery.Type = resultType
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

	newQuery.Src = NewConstruct(lex, FragDeclQry,
		fDeclKeyword,
		fName,
		fParams,
		fType,
	)

	// Define the endpoint
	pr.defineGraphNode(newQuery)

	return newQuery
}
