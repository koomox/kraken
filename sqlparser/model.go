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
		if m.Fields[i].AutoIncrment {
			continue
		}
		key := m.Fields[i].Name
		if m.Fields[i].Name == m.Name {
			key = fmt.Sprintf("%sed", m.Fields[i].Name)
		}
		args = append(args, fmt.Sprintf("%s %s", key, m.Fields[i].TypeOf()))
		params += fmt.Sprintf("\t\t%s: %s,\n", m.Fields[i].ToUpperCase(), key)
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
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment || m.Fields[i].RequiredUpdate || m.Fields[i].RequiredCreated {
			continue
		}
		key := m.Fields[i].Name
		if m.Fields[i].Name == m.Name {
			key = fmt.Sprintf("%sed", m.Fields[i].Name)
		}
		args = append(args, fmt.Sprintf("%s %s", key, m.Fields[i].TypeOf()))
		switch m.Fields[i].TypeOf() {
		case "string":
			params = append(params, fmt.Sprintf("\tif %s != %s.%s {\n\t\tcommand = append(command, fmt.Sprintf(`%s=\"%%v\"`, %s))\n\t}", key, structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].Name, key))
		default:
			params = append(params, fmt.Sprintf("\tif %s != %s.%s {\n\t\tcommand = append(command, fmt.Sprintf(`%s=%%v`, %s))\n\t}", key, structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].Name, key))
		}

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
		if !m.Fields[i].PrimaryKey && !m.Fields[i].RequiredUpdate {
			continue
		}

		if m.Fields[i].RequiredUpdate {
			if m.Fields[i].TypeOf() == "string" {
				params = append(params, fmt.Sprintf("\tcommand += fmt.Sprintf(`, %s=\"%%v\"`, %s)", m.Fields[i].Name, m.Fields[i].Name))
			} else {
				params = append(params, fmt.Sprintf("\tcommand += fmt.Sprintf(`, %s=%%v`, %s)", m.Fields[i].Name, m.Fields[i].Name))
			}
		}

		if m.Fields[i].PrimaryKey {
			keys = append(keys, m.Fields[i].Name)
		}

		key := m.Fields[i].Name
		if m.Fields[i].Name == m.Name {
			key = fmt.Sprintf("%sed", m.Fields[i].Name)
		}
		args = append(args, fmt.Sprintf("%s %s", key, m.Fields[i].TypeOf()))
	}

	b = fmt.Sprintf("func %s(%s, command string) (sql.Result, error) {\n", funcName, strings.Join(args, ", "))
	b += strings.Join(params, "\n")
	b += fmt.Sprintf("\n\treturn %s.%s(command, %s)\n}", m.ToLowerCase(), updateFunc, strings.Join(keys, ", "))

	return
}

func (m *MetadataTable) ToSelectTableModelFuncFormat(selectFunc, tablePrefix, databasePrefix string) string {
	funcName := fmt.Sprintf("%s%s%s", selectFunc, tablePrefix, m.ToUpperCase())

	return fmt.Sprintf("func %s() []*%s.%s {\n\treturn %s.%s()\n}", funcName, databasePrefix, m.ToUpperCase(), m.ToLowerCase(), selectFunc)
}

func (m *MetadataTable) ToRemoveModelFuncFormat(removeFunc string) string {
	funcName := GenerateFunctionName(removeFunc, m.Name)
	names, types, _ := m.ExtractPrimaryFieldFormat()

	return fmt.Sprintf("func %s(%s) (sql.Result, error) {\n\treturn %s.%s(%s)\n}", funcName, strings.Join(types, ", "), m.ToLowerCase(), removeFunc, strings.Join(names, ", "))
}

func (m *MetadataTable) ToWhereModelFuncFormat(whereFunc, databasePrefix string) (b string) {
	b = fmt.Sprintf("func %s%s(command string) []*%s.%s {\n", whereFunc, m.ToUpperCase(), databasePrefix, m.ToUpperCase())
	b += fmt.Sprintf("\treturn %s.%s(command)\n}", m.ToLowerCase(), whereFunc)
	return
}

func (m *MetadataTable) ToSelectModelFuncFormat(fromPrefix, selectPrefix, databasePrefix string) (b string) {
	var params []string
	structName := fmt.Sprintf("%s.%s", databasePrefix, m.ToUpperCase())
	prefixFunc := fmt.Sprintf("%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix)
	names, types, _ := m.ExtractPrimaryFieldFormat()
	mainFunc := GenerateFunctionName(prefixFunc, names...)
	subFunc := GenerateFunctionName(selectPrefix, names...)
	params = append(params, fmt.Sprintf("func %s(%s) []*%s {\n\treturn %s.%s(%s)\n}", mainFunc, strings.Join(types, ", "), structName, m.ToLowerCase(), subFunc, strings.Join(names, ", ")))

	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, "Limit")
			params = append(params, fmt.Sprintf("func %s(%s int) []*%s {\n\treturn %s.%s%s(%s)\n}", funcName, "limit", structName, m.ToLowerCase(), selectPrefix, "Limit", "limit"))
			continue
		}
		if m.Fields[i].HasQuery || m.Fields[i].Unique || m.Fields[i].RequiredCreated || m.Fields[i].RequiredUpdate {
			key := m.Fields[i].Name
			if m.Fields[i].Name == m.Name {
				key = fmt.Sprintf("%sed", m.Fields[i].Name)
			}
			funcName := fmt.Sprintf("%s%s%s%s", fromPrefix, m.ToUpperCase(), selectPrefix, m.Fields[i].ToUpperCase())
			if m.Fields[i].DataType == "DATETIME" || m.Fields[i].DataType == "DATE" {
				params = append(params, fmt.Sprintf("func %s(offset time.Duration, timezone, format string) []*%s {\n\t%s, err := FormatTimeWithOffset(offset, timezone, format)\n\tif err != nil {\n\t\treturn nil\n\t}\n\n\treturn %s.%s%s(%s)\n}", funcName, structName, key, m.ToLowerCase(), selectPrefix, m.Fields[i].ToUpperCase(), key))
			} else {
				params = append(params, fmt.Sprintf("func %s(%s %s) []*%s {\n\treturn %s.%s%s(%s)\n}", funcName, key, m.Fields[i].TypeOf(), structName, m.ToLowerCase(), selectPrefix, m.Fields[i].ToUpperCase(), key))
			}

		}
	}

	return strings.Join(params, "\n\n")
}

func (m *MetadataTable) ToSetModelFuncFormat(funcPrefix, setPrefix string) (b string) {
	names, types, _ := m.ExtractPrimaryAndUpdateFieldFormat()

	for i := range m.Fields {
		if m.Fields[i].PrimaryKey || m.Fields[i].RequiredUpdate || m.Fields[i].RequiredCreated {
			continue
		}
		funcName := fmt.Sprintf("%s%s%s%s", funcPrefix, m.ToUpperCase(), setPrefix, m.Fields[i].ToUpperCase())
		setFunc := fmt.Sprintf("%s.%s%s", m.ToLowerCase(), setPrefix, m.Fields[i].ToUpperCase())
		b += "\n"
		key := m.Fields[i].Name
		if m.Fields[i].Name == m.Name {
			key = fmt.Sprintf("%sed", m.Fields[i].Name)
		}
		b += fmt.Sprintf("func %s(%s %s, %s) (sql.Result, error) {\n", funcName, key, m.Fields[i].TypeOf(), strings.Join(types, ", "))
		b += fmt.Sprintf("\treturn %s(%s, %s)\n}", setFunc, key, strings.Join(names, ", "))
		b += "\n"
	}

	return
}
