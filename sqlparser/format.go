package sqlparser

import (
	"fmt"
	"unicode"
)

func toImportFormat(args ...string) string {
	values := "import(\n"
	for i := range args {
		values += fmt.Sprintf("\t%v%v%v\n", `"`, args[i], `"`)
	}
	values += ")\n"
	return values
}

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

func toFieldLowerFormat(s string) string {
	var ch []rune
	for _, c := range s {
		if c == '_' {
			continue
		}
		ch = append(ch, unicode.ToLower(c))
	}

	return string(ch)
}

func ToUpperCamel(s string) string {
	isSymbol := false
	var ch []rune
	for i, c := range s {
		if c == '_' {
			isSymbol = true
			continue
		}
		if i == 0 {
			isSymbol = true
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

func ToLowerCamel(s string) string {
	isSymbol := false
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

func toLowerCamelFormat(s string) string {
	isSymbol := false
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

func (m *MetadataTable) ToStructFormat(tagField, labelField string) (b string) {
	b = fmt.Sprintf("type %v struct {\n", toFieldUpperFormat(m.Name))
	for i := range m.Fields {
		b += fmt.Sprintf("\t%v %v `%v:%v%v%v %v:%v%v%v`\n", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), tagField, `"`, m.Fields[i].Name, `"`, labelField, `"`, m.Fields[i].Comment, `"`)
	}
	b += "}"
	return
}

func (m *MetadataTable) ToStructSafeFormat(safeName, tagField, labelField string) (b string) {
	b = fmt.Sprintf("type %s%s struct {\n", toFieldUpperFormat(m.Name), safeName)
	for i := range m.Fields {
		b += fmt.Sprintf("\t%v %v `%v:%v%v%v %v:%v%v%v`\n", m.Fields[i].ToUpperCase(), m.Fields[i].TypeSafeOf(), tagField, `"`, m.Fields[i].Name, `"`, labelField, `"`, m.Fields[i].Comment, `"`)
	}
	b += "}"
	return
}

func (m *MetadataTable) ToStructCompareFormat(src, dst, funcName string) (b string) {
	b = fmt.Sprintf("func (%v *%v) %v(%v *%v) bool {\n", src, m.ToUpperCase(), funcName, dst, m.ToUpperCase())
	for i := range m.Fields {
		b += fmt.Sprintf("\tif %v.%v != %v.%v {\n\t\treturn false\n\t}\n", src, m.Fields[i].ToUpperCase(), dst, m.Fields[i].ToUpperCase())
	}
	b += "\treturn true\n}"
	return
}
