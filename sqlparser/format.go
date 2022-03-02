package sqlparser

import (
	"fmt"
	"unicode"
)

func ImportFormatString(args ...string) string {
	values := "import(\n"
	for i := range args {
		values += fmt.Sprintf("\t%v\n", args[i])
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

func (m *MetadataTable) ToStructFormat(tagField string) (b string) {
	b = fmt.Sprintf("type %v struct {\n", toFieldUpperFormat(m.Name))
	for i := range m.Fields {
		b += fmt.Sprintf("\t%v %v ", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf())
		b += fmt.Sprintf("`")
		b += fmt.Sprintf(`%v:"%v"`, tagField, m.Fields[i].Name)
		b += fmt.Sprintf("`\n")
	}
	b += "}\n"
	return
}
