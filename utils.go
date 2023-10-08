package main

import (
	"path/filepath"
)

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
func getGoGeoserverPackageDir() string {
	dir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	return dir
}
