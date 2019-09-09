package parser

import (
	"fmt"
	"sort"

	parser "github.com/romshark/llparser"
)

func isLineBreak(source string, index uint) int {
	switch source[index] {
	case '\n':
		return 1
	case '\r':
		next := index + 1
		if next < uint(len(source)) && source[next] == '\n' {
			return 2
		}
	}
	return -1
}

func isSpecialChar(bt byte) bool {
	if bt >= 0x21 && bt <= 0x2F {
		// ! " # $ % & ' ( ) * + , - . /
		return true
	}
	if bt >= 0x3A && bt <= 0x40 {
		// : ; < = > ? @
		return true
	}
	if bt >= 0x5B && bt <= 0x60 {
		// [ \ ] ^ _ `
		return true
	}
	if bt >= 0x7B && bt <= 0x7E {
		// { | } ~
		return true
	}
	return false
}

func isSpace(bt byte) bool {
	if bt == ' ' || bt == '\t' {
		// whitespace or tab
		return true
	}
	return false
}

// isLowLatinLetter returns true if r is a lower-case latin character,
// otherwise returns false
func isLowLatinLetter(r rune) bool {
	return !(r < 0x61 || r > 0x7A)
}

// isUpLatinLetter returns true if r is an upper-case latin character,
// otherwise returns false
func isUpLatinLetter(r rune) bool {
	return !(r < 0x41 || r > 0x5A)
}

// isDigit returns true if r is a digit character, otherwise returns false
func isDigit(r rune) bool {
	return !(r < 0x30 || r > 0x39)
}

// isLatinLetter returns true if r is either a lower-case or an upper-case
// latin character, otherwise returns false
func isLatinLetter(r rune) bool {
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

func capitalizedCamelCase(str []rune) bool {
	if len(str) < 1 {
		return false
	}

	// [A-Z]
	if !isUpLatinLetter(str[0]) {
		// Non-capitalized first letter
		return false
	}

	for i := 1; i < len(str); i++ {
		r := str[i]
		// [a-zA-Z0-9]
		if !isLatinLetter(r) && !isDigit(r) {
			return false
		}
	}

	return true
}

func lowerCamelCase(str []rune) bool {
	if len(str) < 1 {
		return false
	}

	// [a-z]
	if !isLowLatinLetter(str[0]) {
		// Non-capitalized first letter
		return false
	}

	for i := 1; i < len(str); i++ {
		r := str[i]
		// [a-zA-Z0-9]
		if !isLatinLetter(r) && !isDigit(r) {
			return false
		}
	}

	return true
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

func findElement(
	frags []parser.Fragment,
	kind parser.FragmentKind,
	offset uint,
) (parser.Fragment, uint) {
	for ix := offset; ix < uint(len(frags)); ix++ {
		elem := frags[ix]
		if elem.Kind() == kind {
			return elem, ix
		}
	}
	return nil, 0
}
