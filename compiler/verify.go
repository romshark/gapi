package compiler

import "github.com/pkg/errors"

func isUpLetter(r byte) bool {
	return !(r < 0x41 || r > 0x5A)
}

func isLowLetter(r byte) bool {
	return !(r < 0x61 || r > 0x7A)
}

func isDigit(r byte) bool {
	return !(r < 0x30 || r > 0x39)
}

func isLetter(r byte) bool {
	return isLowLetter(r) || isUpLetter(r)
}

func verifyCapitalizedCamelCase(name string) error {
	if len(name) < 1 {
		return errors.New("empty")
	}

	// [A-Z]
	if !isUpLetter(name[0]) {
		// Non-capitalized first letter
		return errors.New(
			"must begin with a capitalized latin character (A-Z)",
		)
	}

	for i := 1; i < len(name); i++ {
		r := name[i]
		// [a-zA-Z0-9]
		if !isLetter(r) && !isDigit(r) {
			return errors.Errorf(
				"contains illegal character '%s'",
				string(r),
			)
		}
	}

	return nil

}

func verifyLowerCamelCase(name string) error {
	if len(name) < 1 {
		return errors.New("empty")
	}

	// [a-z]
	if !isLowLetter(name[0]) {
		// Non-capitalized first letter
		return errors.New(
			"must begin with a lower case latin character (a-z))",
		)
	}

	for i := 1; i < len(name); i++ {
		r := name[i]
		// [a-zA-Z0-9]
		if !isLetter(r) && !isDigit(r) {
			return errors.Errorf(
				"contains illegal character '%s'",
				string(r),
			)
		}
	}

	return nil
}
