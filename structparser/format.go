package structparser

import (
	"unicode"
)

const (
	columnsFieldFormat = "Y29uc3QgY29sdW1uc0ZpZWxkID0gWwpjb250ZW50RmllbGQKXTs"
	labelFieldFormat = "ICB7CiAgICBsYWJlbDogJ2xhYmVsRmllbGQnLAogICAgZmllbGQ6ICdmaWVsZE5hbWUnLAogICAgcmVuZGVyVHlwZTogJ0lucHV0JywKICAgIGhpZGRlbjogaGlkZGVuRmllbGQsCiAgICB2aXNpYmxlOiB2aXNpYmxlRmllbGQsCiAgICB3cml0YWJsZTogd3JpdGFibGVGaWVsZCwKICAgIHVwZGF0ZVdyaXRhYmxlOiB1cGRhdGVXcml0YWJsZUZpZWxkLAogICAgdXBkYXRlVmlzaWJsZTogdXBkYXRlVmlzaWJsZUZpZWxkLAogIH0s"
)

func toUpperCamelCase(s string) string {
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

func toLowerCase(s string) string {
	var ch []rune
	for _, c := range s {
		if c == '_' {
			continue
		}
		ch = append(ch, unicode.ToLower(c))
	}

	return string(ch)
}

func toLowerCamelCase(s string) string {
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