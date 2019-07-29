package parser

import (
	"errors"
	"fmt"
	"sort"
)

// isSpace returns true if src[i] is either a whitespace, a tab, a line-feed
// or a combination of the carriage-return and line-feed characters,
// otherwise returns false
func isSpace(src string, i uint32) bool {
	switch src[i] {
	case ' ':
		return true
	case '\t':
		return true
	case '\n':
		return true
	case '\r':
		if i+1 < uint32(len(src)) &&
			src[i+1] == '\n' {
			return false
		}
	}
	return false
}

// isLatinAlphanum returns true if src[i] is either a digit, a lower-case
// lating character or an upper-case latin character, otherwise returns false
func isLatinAlphanum(src string, i uint32) bool {
	if isSpace(src, i) {
		return false
	}
	c := src[i]
	if c >= 0x30 && c <= 0x39 {
		// Digit
		return true
	}
	if c >= 0x41 && c <= 0x5A {
		// Latin uppercase rune
		return true
	}
	if c >= 0x61 && c <= 0x7A {
		// Latin lowercase rune
		return true
	}
	return false
}

// isLowLatinLetter returns true if r is a lower-case latin character,
// otherwise returns false
func isLowLatinLetter(r byte) bool {
	return !(r < 0x61 || r > 0x7A)
}

// isUpLatinLetter returns true if r is an upper-case latin character,
// otherwise returns false
func isUpLatinLetter(r byte) bool {
	return !(r < 0x41 || r > 0x5A)
}

// isDigit returns true if r is a digit character, otherwise returns false
func isDigit(r byte) bool {
	return !(r < 0x30 || r > 0x39)
}

// isLatinLetter returns true if r is either a lower-case or an upper-case
// latin character, otherwise returns false
func isLatinLetter(r byte) bool {
	return isLowLatinLetter(r) || isUpLatinLetter(r)
}

func verify(
	token *Token,
	tokenType string,
	verificationMethod func(string) error,
) Error {
	if err := verificationMethod(token.src); err != nil {
		return &pErr{
			at:   token.begin,
			code: ErrSyntax,
			message: fmt.Sprintf(
				"illegal %s '%s' (%s)",
				tokenType,
				token.src,
				err,
			),
		}
	}
	return nil
}

func capitalizedCamelCase(ident string) error {
	if len(ident) < 1 {
		return errors.New("empty")
	}

	// [A-Z]
	if !isUpLatinLetter(ident[0]) {
		// Non-capitalized first letter
		return errors.New(
			"must begin with a capitalized latin character (A-Z)",
		)
	}

	for i := 1; i < len(ident); i++ {
		r := ident[i]
		// [a-zA-Z0-9]
		if !isLatinLetter(r) && !isDigit(r) {
			return fmt.Errorf(
				"contains illegal character '%s'",
				string(r),
			)
		}
	}

	return nil

}

func lowerCamelCase(ident string) error {
	if len(ident) < 1 {
		return errors.New("empty")
	}

	// [a-z]
	if !isLowLatinLetter(ident[0]) {
		// Non-capitalized first letter
		return errors.New(
			"must begin with a lower case latin character (a-z))",
		)
	}

	for i := 1; i < len(ident); i++ {
		r := ident[i]
		// [a-zA-Z0-9]
		if !isLatinLetter(r) && !isDigit(r) {
			return fmt.Errorf(
				"contains illegal character '%s'",
				string(r),
			)
		}
	}

	return nil
}

func sortTypesByName(types []Type) {
	sort.Slice(types, func(i, j int) bool {
		return types[i].String() < types[j].String()
	})
}

func sortQueryEndpointsByName(endpoints []*Query) {
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].GraphNodeName() < endpoints[j].GraphNodeName()
	})
}

func sortMutationsByName(mutations []*Mutation) {
	sort.Slice(mutations, func(i, j int) bool {
		return mutations[i].GraphNodeName() < mutations[j].GraphNodeName()
	})
}

func stringifyType(t Type) (name string) {
	if t == nil {
		return
	}
	for {
		if v, isOptional := t.(*TypeOptional); isOptional {
			if v.StoreType == nil {
				name += "?<unknown>"
				break
			}
			name += "?"
			t = v.StoreType
			continue
		}
		if v, isList := t.(*TypeList); isList {
			if v.StoreType == nil {
				name += "[]<unknown>"
				break
			}
			name += "[]"
			t = v.StoreType
			continue
		}
		name += t.String()
		break
	}
	return
}

func readWord(
	lex *Lexer,
	expectation string,
	fragID FragID,
	verificationMethod func(string) error,
) (*Token, Error) {
	tk, err := lex.NextExpectSkip(
		FragTkLatinAlphanum,
		Skip{FragTkSpace},
		"expected "+expectation,
	)
	if err != nil {
		return nil, err
	}
	tk.id = fragID
	if err := verify(
		tk,
		expectation,
		verificationMethod,
	); err != nil {
		return nil, err
	}
	return tk, nil
}

func readWordExact(
	lex *Lexer,
	expectedWord string,
	fragID FragID,
	expectation string,
) (*Token, Error) {
	tk, err := lex.NextExpectSkip(
		FragTkLatinAlphanum,
		Skip{FragTkSpace},
		"expected "+expectation,
	)
	if err != nil {
		return nil, err
	}
	tk.id = fragID
	if tk.Src() != expectedWord {
		return nil, &pErr{
			at:   tk.Begin(),
			code: ErrSyntax,
			message: fmt.Sprintf(
				"expected %s, got: '%s'",
				expectation,
				tk.Src(),
			),
		}
	}
	return tk, nil
}

func readToken(
	lex *Lexer,
	expectedFragID FragID,
	expectation string,
) (*Token, Error) {
	return lex.NextExpectSkip(
		expectedFragID,
		Skip{FragTkSpace},
		"expected "+expectation,
	)
}
