package parser

import "fmt"

// onAnonymousType is executed when an anonymous type is parsed.
// It checks whether this anonymous type was already defined before
// and returns either the defined one if it was or the newly registered one
// if it wasn't yet
func (pr *Parser) onAnonymousType(newType Type) Type {
	name := newType.String()

	switch newType.(type) {
	case *TypeOptional:
	case *TypeList:
	default:
		panic(fmt.Errorf("%s isn't an anonymous type", name))
	}

	// Check whether this anonymous type was already defined before
	// and return the defined one if it was
	if defined, isDefined := pr.typeByName[name]; isDefined {
		return defined
	}

	// Issue a new type ID
	pr.lastIssuedTypeID += TypeID(1)
	newID := pr.lastIssuedTypeID

	// Register the newly defined anonymous type
	pr.typeByID[newID] = newType
	pr.typeByName[name] = newType
	pr.ast.Types = append(pr.ast.Types, newType)
	pr.ast.AnonymousTypes = append(pr.ast.AnonymousTypes, newType)

	// Set ID
	switch t := newType.(type) {
	case *TypeOptional:
		t.ID = newID
	case *TypeList:
		t.ID = newID
	}

	return newType
}
