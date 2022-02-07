package sqlparser

import (
	"encoding/base64"
	"strings"
	"unicode"
)

const (
	structFormat      = "type structName struct {\ncontentStr\n}"
	structFieldFormat = "CWZpZWxkTmFtZSBmaWVsZERhdGFUeXBlIGB0YWdGaWVsZDoidGFnTmFtZSJg"
)

func toFieldUpperFormat(s string) string {
	isSymbol := true
	var ch []rune
	for _, c := range s {
		if c == '_' {
			isSymbol = true
			continue
		}
		if isSymbol && c != '_' {
			ch = append(ch, unicode.ToUpper(c))
			isSymbol = false
			continue
		}
		ch = append(ch, c)
	}

	return string(ch)
}

func (m *MetadataTable) toStructFieldFormat(tagField string) (elements []string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(structFieldFormat)
	var element string
	for i := 0; i < len(m.Fields); i++ {
		element = strings.Replace(string(fieldFormat), "fieldName", toFieldUpperFormat(m.Fields[i].Name), -1)
		element = strings.Replace(element, "tagField", tagField, -1)
		element = strings.Replace(element, "tagName", m.Fields[i].Name, -1)
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			element = strings.Replace(element, "fieldDataType", "int", -1)
		case "BIGINT":
			element = strings.Replace(element, "fieldDataType", "int64", -1)
		default:
			element = strings.Replace(element, "fieldDataType", "string", -1)
		}
		elements = append(elements, element)
	}
	return
}

func (m *MetadataTable) ToStructFormat(tagField string) (b string) {
	elements := m.toStructFieldFormat(tagField)
	b = strings.Replace(structFormat, "structName", toFieldUpperFormat(m.Name), -1)
	b = strings.Replace(b, "contentStr", strings.Join(elements, "\n"), -1)
	return
}
