package sqlparser

import (
	"strings"
)

func findField(s string) (element *Field) {
	element = &Field{}
	parts := SplitBySpaceOutsideQuotes(s)
	for _, part := range parts {
		switch {
		case element.HasComment:
			element.Comment = Trim(part)
			element.HasComment = false
		case strings.EqualFold(part, "UNIQUE"):
			element.Unique = true
		case strings.EqualFold(part, "AUTO_INCREMENT"):
			element.AutoIncrment = true
		case strings.EqualFold(part, "PRIMARY"):
			element.PrimaryKey = true
		case strings.EqualFold(part, "COMMENT"):
			element.HasComment = true
		case element.DataType == "" && findDataTypeString(part) != "":
			element.DataType = findDataTypeString(part)
		case element.Name == "" && findKeywordString(part) == "":
			element.Name = Trim(part)
		default:
		}
		if element.Name != "" && element.Comment == "" {
			element.Comment = element.Name
		}
	}
	return
}

func matchTableName(s string) string {
	if idx := strings.IndexByte(s, '.'); idx != -1 {
		return s[idx+1:]
	}
	return s
}

func findTableName(s string) string {
	options := Split(s, " ")
	for i := range options {
		v := options[i]
		if findKeywordString(v) == "" {
			return matchTableName(v)
		}
	}

	return ""
}

func findPrimaryKey(s string) (elements []string) {
	options := SplitAndTrimSpecialChars(s, " ")
	for i := range options {
		v := options[i]
		if findKeywordString(v) == "" {
			elements = append(elements, v)
		}
	}
	return
}

func TrimFunc(s string, f func(rune) bool) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		if !f(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func Trim(s string) string {
	return TrimFunc(s, func(r rune) bool {
		return r == ',' || r == '`' || r == '"' || r == '\''
	})
}

func TrimArray(parts ...string) []string {
	p := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		trimmed := Trim(part)
		if trimmed != "" {
			p = append(p, trimmed)
		}
	}
	return p
}

func TrimSpecialChars(s string) string {
	return TrimFunc(s, func(r rune) bool {
		return r == ',' || r == '`' || r == '"' || r == '\'' || r == '(' || r == ')'
	})
}

func SplitAndTrimSpecialChars(s, sep string) (elements []string) {
	r := strings.Split(s, sep)
	for i := range r {
		v := TrimSpecialChars(r[i])
		if v == "" {
			continue
		}
		elements = append(elements, v)
	}
	return
}

func GenerateFunctionName(prefix string, keywords ...string) string {
	var builder strings.Builder
	builder.WriteString(prefix)

	for i, keyword := range keywords {
		if i != 0 {
			builder.WriteString("And")
		}
		builder.WriteString(ToUpperCamel(keyword))
	}
	return builder.String()
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

func SplitBySpaceOutsideQuotes(s string) []string {
	var result []string
	var part strings.Builder
	inSingleQuote, inDoubleQuote, inBacktick := false, false, false
	for _, c := range s {
		switch c {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			}
			part.WriteRune(c)
		case '"':
			if !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			}
			part.WriteRune(c)
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
			part.WriteRune(c)
		case ' ':
			if !inSingleQuote && !inDoubleQuote && !inBacktick {
				if part.Len() > 0 {
					result = append(result, part.String())
					part.Reset()
				}
			} else {
				part.WriteRune(c)
			}
		default:
			part.WriteRune(c)
		}
	}

	if part.Len() > 0 {
		result = append(result, part.String())
	}

	return result
}

func FromFile(filename string) (source *Database) {
	source = &Database{}
	table := &MetadataTable{}
	isValid := false
	readFile(func(s string) {
		if s == "" || strings.HasPrefix(s, "--") {
			return
		}
		parts := SplitBySpaceOutsideQuotes(s)
		v := findKeywordString(parts[0])
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

	source.EnableQueryFields(queryKeywords...)
	source.EnableRequiredUpdateFields(requiredUpdateKeywords...)
	source.EnableRequiredCreatedFields(requiredCreatedKeywords...)
	source.EnableRequiredDatetimeFields()

	return
}
