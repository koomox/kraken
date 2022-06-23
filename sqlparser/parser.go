package sqlparser

import (
	"strings"
)

func findField(s string) (element *Field) {
	element = &Field{}
	b := Split(s, " ")
	for i := 0; i < len(b); i++ {
		if b[i] == "" {
			continue
		}
		if element.Name != "" && strings.EqualFold(b[i], "UNIQUE") {
			element.Unique = true
		}
		if element.Name != "" && strings.EqualFold(b[i], "AUTO_INCREMENT") {
			element.AutoIncrment = true
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

func matchTableName(s string) string {
	var ch []byte
	isValid := false
	for i := range s {
		switch s[i] {
		case '.':
			isValid = true
		default:
			if isValid {
				ch = append(ch, s[i])
			}
		}
	}
	return string(ch)
}

func findTableName(s string) string {
	options := Split(s, " ")
	for i := range options {
		v := options[i]
		if findKeywordString(v) == "" && strings.Contains(v, ".") {
			return matchTableName(v)
		}
	}

	return ""
}

func findPrimaryKey(s string) (elements []string) {
	options := Split(s, " ")
	for i := range options {
		v := options[i]
		if findKeywordString(v) == "" {
			elements = append(elements, v)
		}
	}
	return
}

func Trim(s string) string {
	var ch []byte
	for i := range s {
		switch s[i] {
		case ',', '(', ')', '`', '"':
		default:
			ch = append(ch, s[i])
		}
	}
	return string(ch)
}

func Split(s, sep string) (elements []string) {
	r := strings.Split(s, sep)
	for i := range r {
		v := Trim(r[i])
		if v == "" {
			continue
		}
		elements = append(elements, v)
	}
	return
}

func FromFile(filename string) (source *Database) {
	source = &Database{}
	table := &MetadataTable{}
	isValid := false
	readFile(func(s string) {
		if s == "" || strings.HasPrefix(s, "--") {
			return
		}
		options := Split(s, " ")
		v := findKeywordString(options[0])
		switch v {
		case "PRIMARY":
			if strings.HasPrefix(s, "PRIMARY KEY") {
				if keys := findPrimaryKey(strings.TrimPrefix(s, "PRIMARY KEY")); keys != nil && len(keys) > 0 {
					for i := range keys {
						table.SetPrimaryKey(keys[i])
					}
				}
			}
			return
		case "UNIQUE":
			if strings.HasPrefix(s, "UNIQUE INDEX") {
				return
			}
		case "CREATE":
			if strings.HasPrefix(s, "CREATE TABLE") && strings.HasSuffix(s, "(") {
				isValid = true
				table = &MetadataTable{Name: findTableName(s)}
			}
			return
		default:
			if strings.HasPrefix(s, ")") && isValid {
				isValid = false
				source.Tables = append(source.Tables, table)
				table = &MetadataTable{}
				return
			}
		}
		if v != "" {
			return
		}
		table.Fields = append(table.Fields, findField(s))
	}, filename)

	return
}
