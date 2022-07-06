package sqlparser

import (
	"fmt"
	"strings"
)

func (m *MetadataTable) ToSelectStorageFuncFormat(fromPrefix, selectPrefix, databasePrefix, storePrefix, StorePrefix, currentPrefix string) (b string) {
	structName := fmt.Sprintf("%s.%s", databasePrefix, m.ToUpperCase())
	var idx []string
	var ids []string
	var args []string
	keys := m.PrimaryKey()
	switch len(keys) {
	case 1:
		funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, keys[0].ToUpperCase())
		b = fmt.Sprintf("func (%s *%s) %s(%s %s) *%s {\n\treturn %s.%s.%s%s(%s)\n}", storePrefix, StorePrefix, funcName, keys[0].ToLowerCase(), keys[0].TypeOf(), structName, storePrefix, m.ToUpperCase(), selectPrefix, keys[0].ToUpperCase(), keys[0].ToLowerCase())
		b += fmt.Sprintf("\n\nfunc %s(%s %s) *%s {\n\treturn %s.%s(%s)\n}", funcName, keys[0].ToLowerCase(), keys[0].TypeOf(), structName, currentPrefix, funcName, keys[0].ToLowerCase())
	default:
		for _, v := range keys {
			idx = append(idx, fmt.Sprintf("%s", v.ToUpperCase()))
			ids = append(ids, fmt.Sprintf("%s %s", v.ToLowerCase(), v.TypeOf()))
			args = append(args, v.ToLowerCase())
		}
		funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, strings.Join(idx, "And"))
		b = fmt.Sprintf("func (%s *%s) %s(%s) *%s {\n\treturn %s.%s.%s%s(%s)\n}", storePrefix, StorePrefix, funcName, strings.Join(ids, ", "), structName, storePrefix, m.ToUpperCase(), selectPrefix, strings.Join(idx, "And"), strings.Join(args, ", "))
		b += fmt.Sprintf("\n\nfunc %s(%s) *%s {\n\treturn %s.%s(%s)\n}", funcName, strings.Join(ids, ", "), structName, currentPrefix, funcName, strings.Join(args, ", "))
	}

	for i := range m.Fields {
		if !m.Fields[i].PrimaryKey && m.Fields[i].Unique {
			funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, m.Fields[i].ToUpperCase())
			b += fmt.Sprintf("\n\nfunc (%s *%s) %s(%s %s) *%s {\n\treturn %s.%s.%s%s(%s)\n}", storePrefix, StorePrefix, funcName, m.Fields[i].ToLowerCase(), m.Fields[i].TypeOf(), structName, storePrefix, m.ToUpperCase(), selectPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].ToLowerCase())
			b += fmt.Sprintf("\n\nfunc %s(%s %s) *%s {\n\treturn %s.%s(%s)\n}", funcName, m.Fields[i].ToLowerCase(), m.Fields[i].TypeOf(), structName, currentPrefix, funcName, m.Fields[i].ToLowerCase())
		}
	}

	return
}
