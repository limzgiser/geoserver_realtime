package main

func GetGsFieldType(gsType string) string {
	switch gsType {
	case "Boolean":
		return "Boolean"
	case "Byte", "Short", "Integer", "Long", "Float", "Double":
		return "Number"
	default:
		return "String"
	}
}
