package parser

// stdTypeByName returns a standard primitive type instance by name
// or nil if name doesn't identify any built-in primitive type
func stdTypeByName(name string) Type {
	switch name {
	case "None":
		return TypeStdNone{}
	case "Bool":
		return TypeStdBool{}
	case "Byte":
		return TypeStdByte{}
	case "Int32":
		return TypeStdInt32{}
	case "Uint32":
		return TypeStdUint32{}
	case "Int64":
		return TypeStdInt64{}
	case "Uint64":
		return TypeStdUint64{}
	case "Float64":
		return TypeStdFloat64{}
	case "String":
		return TypeStdString{}
	case "Time":
		return TypeStdTime{}
	default:
		return nil
	}
}
