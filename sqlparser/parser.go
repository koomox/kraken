package sqlparser

import (
	"strings"
)

func findField(s string) (element Field) {
	b := strings.Split(s, " ")
	for i := 0; i < len(b); i++ {
		if b[i] == "" {
			continue
		}
		if element.Name != "" && strings.EqualFold(b[i], "UNIQUE") {
			element.Unique = true
		}
		if findKeywordString(b[i]) != "" {
			continue
		}
		if v := findDataTypeString(b[i]); v != "" {
			element.DataType = v
			continue
		}
		if element.Name == "" {
			element.Name = b[i]
		}
	}
	return
}

func matchTableName(s, sub string) (elements []string) {
	b := strings.Split(s, " ")
	for i := 0; i < len(b); i++ {
		if b[i] == "" || b[i] == sub {
			continue
		}
		if findKeywordString(b[i]) != "" {
			continue
		}
		elements = append(elements, b[i])
	}
	return
}

func findTableName(s, sub string) (element string) {
	elements := matchTableName(s, sub)
	element = strings.Join(elements, " ")
	element = strings.TrimSpace(element)
	if strings.Contains(element, ".") {
		elements = strings.Split(element, ".")
		element = elements[1]
	}
	return
}

func FromFile(filename string) (elements []MetadataTable) {
	element := &MetadataTable{}
	readFile(func(s string) {
		if s == "" || strings.HasPrefix(s, "--") || strings.HasPrefix(s, "CREATE DATABASE") {
			return
		}
		lines := strings.Split(s, " ")
		switch lines[0] {
		case "DESC", "USE", "SELECT", "INSERT", "GRANT", "PRIMARY", "UNIQUE", "DROP":
			return
		default:
		}
		if strings.HasPrefix(s, ")") {
			elements = append(elements, *element)
			element = &MetadataTable{}
			return
		}
		if strings.HasSuffix(s, "(") {
			element = &MetadataTable{Name: findTableName(s, "(")}
			return
		}
		element.Fields = append(element.Fields, findField(s))
	}, filename)
	return
}
