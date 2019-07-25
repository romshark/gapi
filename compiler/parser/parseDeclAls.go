package parser

import (
	"fmt"
)

func (pr *Parser) parseDeclAls(lex *Lexer) *TypeAlias {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordAlias,
		FragTkKwdAls,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type ID
	fTypeID, err := readWord(
		lex,
		"alias type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Read '='
	fSymEq, err := readToken(lex, FragTkSymEq, "equals sign")
	if pr.err(err) {
		return nil
	}

	// Read aliased type identifier
	fAliasedTypeID, err := readWord(
		lex,
		"aliased type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Instantiate type
	newType := &TypeAlias{
		terminalType: terminalType{
			TypeName: fTypeID.src,
		},
	}
	newType.Src = NewConstruct(lex, FragDeclAls,
		fDeclKeyword,
		fTypeID,
		fSymEq,
		fAliasedTypeID,
	)

	// Define the type
	pr.defineType(newType)

	pr.deferJob(func() {
		// Ensure the aliased type exists after all types have been defined
		aliasedType := pr.findTypeByName(fAliasedTypeID.src)
		if aliasedType != nil {
			// Set the aliased type
			newType.AliasedType = aliasedType
			return
		}
		pr.err(&pErr{
			at:   fDeclKeyword.begin,
			code: ErrTypeUndef,
			message: fmt.Sprintf(
				"undefined type %s aliased by %s",
				fAliasedTypeID.src,
				fTypeID.src,
			),
		})
	})

	return newType
}
