package compiler

import "github.com/pkg/errors"

func isUpLatin(r byte) bool {
	return !(r < 0x41 || r > 0x5A)
}

func isLowLatinAlpha(r byte) bool {
	return !(r < 0x61 || r > 0x7A)
}

func isDigit(r byte) bool {
	return !(r < 0x30 || r > 0x39)
}

func isLowLatinAlphanum(r byte) bool {
	return isLowLatinAlpha(r) || isDigit(r)
}

func verifyTypeName(name string) error {
	if len(name) < 1 {
		return errors.New("missing type name")
	}
	if !isUpLatin(name[0]) {
		// Non-capitalized first letter
		return errors.New(
			"not a type identifier. Type names must begin with " +
				"a capitalized latin character (A-Z)",
		)
	}

	for i := 1; i < len(name); i++ {
		if !isLowLatinAlphanum(name[i]) {
			return errors.Errorf(
				"not a type identifier. "+
					"Type name contains illegal character '%s'",
				string(name[i]),
			)
		}
	}

	return nil
}
