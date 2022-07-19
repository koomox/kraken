package sqlparser

import (
	"fmt"
	"strings"
)

func (m *MetadataTable) ToCreateModelFuncFormat(createFunc, insertFunc, databasePrefix string) (b string) {
	var args []string
	var params string
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		args = append(args, fmt.Sprintf("%s %s", m.Fields[i].Name, m.Fields[i].TypeOf()))
		params += fmt.Sprintf("\t\t%s: %s,\n", m.Fields[i].ToUpperCase(), m.Fields[i].Name)
	}
	b = fmt.Sprintf("func %s%s(%s) (sql.Result, error) {\n", createFunc, m.ToUpperCase(), strings.Join(args, ", "))
	b += fmt.Sprintf("\treturn %s.%s(&%s.%s{\n", m.ToLowerCase(), insertFunc, databasePrefix, m.ToUpperCase())
	b += params
	b += "\t})\n}"
	return
}

func (m *MetadataTable) ToCompareModelFuncFormat(compareFunc, structPrefix, databasePrefix string) (b string) {
	var args []string
	var params []string
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		switch m.Fields[i].Name {
		case "created_by", "updated_by", "created_at", "updated_at":
			continue
		}
		args = append(args, fmt.Sprintf("%s %s", m.Fields[i].Name, m.Fields[i].TypeOf()))
		params = append(params, fmt.Sprintf("\tif %s != %s.%s {\n\t\tcommand = append(command, fmt.Sprintf(`%s=%%v`, %s))\n\t}", m.Fields[i].Name, structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].Name, m.Fields[i].Name))
	}

	b = fmt.Sprintf("func %s%s(%s, %s *%s.%s) string {\n", compareFunc, m.ToUpperCase(), strings.Join(args, ", "), structPrefix, databasePrefix, m.ToUpperCase())
	b += "\tvar command []string\n"
	b += strings.Join(params, "\n")
	b += "\n\tif command == nil || len(command) == 0 {\n\t\treturn \"\"\n\t}\n\n\treturn strings.Join(command, \", \")\n}"
	return
}

func (m *MetadataTable) ToUpdateModelFuncFormat(updateFunc string) (b string) {
	var args []string
	var params []string
	var keys []string
	funcName := fmt.Sprintf("%s%s", updateFunc, m.ToUpperCase())
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at":
			params = append(params, fmt.Sprintf("\tcommand += fmt.Sprintf(`, %s=%%v`, %s)", m.Fields[i].Name, m.Fields[i].Name))
		default:
			if !m.Fields[i].PrimaryKey {
				continue
			}
			keys = append(keys, m.Fields[i].Name)
		}
		args = append(args, fmt.Sprintf("%s %s", m.Fields[i].Name, m.Fields[i].TypeOf()))
	}

	b = fmt.Sprintf("func %s(%s, command string) (sql.Result, error) {\n", funcName, strings.Join(args, ", "))
	b += strings.Join(params, "\n")
	b += fmt.Sprintf("\n\treturn %s.%s(command, %s)\n}", m.ToLowerCase(), updateFunc, strings.Join(keys, ", "))

	return
}

func (m *MetadataTable) ToRemoveModelFuncFormat(removeFunc string) (b string) {
	var keys []string
	var args []string
	funcName := fmt.Sprintf("%s%s", removeFunc, m.ToUpperCase())
	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at":
		default:
			if !m.Fields[i].PrimaryKey {
				continue
			}
		}
		keys = append(keys, m.Fields[i].Name)
		args = append(args, fmt.Sprintf("%s %s", m.Fields[i].Name, m.Fields[i].TypeOf()))
	}
	b = fmt.Sprintf("func %s(%s) (sql.Result, error) {\n\treturn %s.%s(%s)\n}", funcName, strings.Join(args, ", "), m.ToLowerCase(), removeFunc, strings.Join(keys, ", "))

	return
}

func (m *MetadataTable) ToWhereModelFuncFormat(whereFunc, databasePrefix string) (b string) {
	b = fmt.Sprintf("func %s%s(command string) []*%s.%s {\n", whereFunc, m.ToUpperCase(), databasePrefix, m.ToUpperCase())
	b += fmt.Sprintf("\treturn %s.%s(command)\n}", m.ToLowerCase(), whereFunc)
	return
}

func (m *MetadataTable) ToSelectModelFuncFormat(fromPrefix, selectPrefix, databasePrefix string) (b string) {
	structName := fmt.Sprintf("%s.%s", databasePrefix, m.ToUpperCase())
	var idx []string
	var ids []string
	var args []string
	var params []string
	keys := m.PrimaryKey()
	switch len(keys) {
	case 1:
		funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, keys[0].ToUpperCase())
		params = append(params, fmt.Sprintf("func %s(%s %s) []*%s {\n\treturn %s.%s%s(%s)\n}", funcName, keys[0].Name, keys[0].TypeOf(), structName, m.ToLowerCase(), selectPrefix, keys[0].ToUpperCase(), keys[0].Name))
	default:
		for _, v := range keys {
			idx = append(idx, fmt.Sprintf("%s", v.ToUpperCase()))
			ids = append(ids, fmt.Sprintf("%s %s", v.ToLowerCase(), v.TypeOf()))
			args = append(args, v.ToLowerCase())
		}
		funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, strings.Join(idx, "And"))
		params = append(params, fmt.Sprintf("func %s(%s) []*%s {\n\treturn %s.%s%s(%s)\n}", funcName, strings.Join(ids, ", "), structName, m.ToLowerCase(), selectPrefix, strings.Join(idx, "And"), strings.Join(args, ", ")))
	}

	
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			continue
		}
		if m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}

		funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, m.Fields[i].ToUpperCase())
		params = append(params, fmt.Sprintf("func %s(%s %s) []*%s {\n\treturn %s.%s%s(%s)\n}", funcName, m.Fields[i].Name, m.Fields[i].TypeOf(), structName, m.ToLowerCase(), selectPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].Name))
	}

	return strings.Join(params, "\n\n")
}