package sqlparser

import (
	"fmt"
	"strings"
)

func (m *MetadataTable) ToSelectSQLFormat(funcName string) (b string) {
	return fmt.Sprintf("func %s(table string) string {\n\treturn fmt.Sprintf(`SELECT * FROM %%s`, table)\n}", funcName)
}

func (m *MetadataTable) ToInsertSQLFormat(funcName, structPrefix, structName string) (b string) {
	var keys []string
	var values []string
	var elements []string
	for i := range m.Fields {
		if m.Fields[i].AutoIncrment {
			continue
		}
		values = append(values, m.Fields[i].ValueOf())
		keys = append(keys, m.Fields[i].Name)
		elements = append(elements, fmt.Sprintf("%s.%s", structPrefix, m.Fields[i].ToUpperCase()))
	}
	b = fmt.Sprintf("func %s(%s *%s, table string) string {\n", funcName, structPrefix, structName)
	b += fmt.Sprintf("\treturn fmt.Sprintf(`INSERT INTO %%s(%s) VALUES(%s)`, table, %s)\n}", strings.Join(keys, ", "), strings.Join(values, ", "), strings.Join(elements, ", "))
	return
}

func (m *MetadataTable) ToUpdateSQLFormat(funcName string) (b string) {
	names, types, formats := m.ExtractPrimaryFieldFormat()
	b = fmt.Sprintf("func %s(command string, %s, table string) string {\n", funcName, strings.Join(types, ", "))
	b += fmt.Sprintf("\treturn fmt.Sprintf(`UPDATE %%s SET %%s WHERE %s`, table, command, %s)\n}", strings.Join(formats, " AND "), strings.Join(names, ", "))
	return
}

func (m *MetadataTable) ToRemoveSQLFormat(funcName string) (b string) {
	names, types, formats := m.ExtractPrimaryFieldFormat()
	b = fmt.Sprintf("func %s(%s, table string) string {\n", funcName, strings.Join(types, ", "))
	b += fmt.Sprintf("\treturn fmt.Sprintf(`DELETE FROM %%s WHERE %s`, table, %s)\n}", strings.Join(formats, " AND "), strings.Join(names, ", "))
	return
}

func (m *MetadataTable) ToWhereSQLFormat(funcName string) (b string) {
	b = fmt.Sprintf("func %s(command, table string) string {\n", funcName)
	b += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %%s WHERE %%s`, table, command)\n}")
	return
}

func (m *MetadataTable) ToQuerySQLFormat(funcName, structPrefix, structName string) (b string) {
	b = fmt.Sprintf("func %s(command string) (%s []*%s) {\n", funcName, structPrefix, structName)
	b += "\tdata, length := mysql.Query(command)\n"
	b += "\tif data == nil || length <= 0 {\n\t\treturn\n\t}\n"
	b += "\tb := *data\n"
	b += fmt.Sprintf("\tfor i := 0; i < length; i++ {\n\t\telement := parser(b[i])\n\t\t%s = append(%s, element)\n\t}\n", structPrefix, structPrefix)
	b += "\treturn\n}"
	return
}

func (m *MetadataTable) ToParserSQLFormat(funcName, prefixName, structName, databasePrefix string) (b string) {
	b = fmt.Sprintf("func %s(%s map[string]string) *%s {\n", funcName, prefixName, structName)
	b += fmt.Sprintf("\treturn &%s{\n", structName)
	for i := range m.Fields {
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			b += fmt.Sprintf("\t\t%s:%s.ParseInt(%s[\"%s\"]),\n", m.Fields[i].ToUpperCase(), databasePrefix, prefixName, m.Fields[i].Name)
		case "BIGINT":
			b += fmt.Sprintf("\t\t%s:%s.ParseInt64(%s[\"%s\"]),\n", m.Fields[i].ToUpperCase(), databasePrefix, prefixName, m.Fields[i].Name)
		default:
			b += fmt.Sprintf("\t\t%s:%s[\"%s\"],\n", m.Fields[i].ToUpperCase(), prefixName, m.Fields[i].Name)
		}
	}
	b += "\t}\n}"
	return
}

func (m *MetadataTable) ToSubSelectSQLFormat(prefixFunc string) (b string) {
	names, types, formats := m.ExtractPrimaryFieldFormat()
	mainFunc := GenerateFunctionName(prefixFunc, names...)
	b += fmt.Sprintf("func %s(%s, table string) string {\n", mainFunc, strings.Join(types, ", "))
	b += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %%s WHERE %s`, table, %s)\n}\n\n", strings.Join(formats, " AND "), strings.Join(names, ", "))

	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			continue
		}
		if strings.EqualFold(m.Fields[i].Name, "created_by") || m.Fields[i].HasQuery || m.Fields[i].Unique {
			b += fmt.Sprintf("func %s(%s %s, table string) string {\n", GenerateFunctionName(prefixFunc, m.Fields[i].Name), m.Fields[i].Name, m.Fields[i].TypeOf())
			b += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %%s WHERE %s=%s`, table, %s)\n}\n\n", m.Fields[i].Name, m.Fields[i].ValueOf(), m.Fields[i].Name)
		}
	}
	return
}

func (m *MetadataTable) ToSetSQLFormat(funcPrefix string) (b string) {
	names, types, format := m.ExtractPrimaryFieldFormat()
	keys, args, upFormat := m.ExtractUpdateFieldFormat()
	keys = append(keys, names...)
	types = append(types, args...)

	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at", "created_by", "created_at":
			continue
		default:
			if m.Fields[i].PrimaryKey {
				continue
			}
		}
		funcName := funcPrefix + m.Fields[i].ToUpperCase()
		formats := []string{fmt.Sprintf("%v=%v", m.Fields[i].Name, m.Fields[i].ValueOf())}
		formats = append(formats, upFormat...)
		keywords := []string{m.Fields[i].Name}
		keywords = append(keywords, keys...)
		b += fmt.Sprintf("\nfunc %s(%s %s, %s, table string) string {\n", funcName, m.Fields[i].Name, m.Fields[i].TypeOf(), strings.Join(types, ", "))
		b += fmt.Sprintf("\treturn fmt.Sprintf(`UPDATE %%s SET %s WHERE %s`, table, %s)\n}\n", strings.Join(formats, ", "), strings.Join(format, " AND "), strings.Join(keywords, ", "))
	}

	return
}

func (m *MetadataTable) ToSubSelectCrudFormat(prefixFunc, queryFunc, subPrefixFunc, structName, tableName string) (b string) {
	names, types, _ := m.ExtractPrimaryFieldFormat()
	mainFunc := GenerateFunctionName(prefixFunc, names...)
	subFunc := GenerateFunctionName(subPrefixFunc, names...)
	b = fmt.Sprintf("func %s(%s) []*%s{\n", mainFunc, strings.Join(types, ", "), structName)
	b += fmt.Sprintf("\treturn %s(%s(%s, %s))\n}\n\n", queryFunc, subFunc, strings.Join(names, ", "), tableName)

	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			continue
		}
		if strings.EqualFold(m.Fields[i].Name, "created_by") || m.Fields[i].HasQuery || m.Fields[i].Unique {
			b += fmt.Sprintf("func %s(%s %s) []*%s {\n", GenerateFunctionName(prefixFunc, m.Fields[i].ToUpperCase()), m.Fields[i].Name, m.Fields[i].TypeOf(), structName)
			b += fmt.Sprintf("\treturn %s(%s(%s, %s))\n}\n\n", queryFunc, GenerateFunctionName(subPrefixFunc, m.Fields[i].ToUpperCase()), m.Fields[i].Name, tableName)
		}
	}
	return
}

func (m *MetadataTable) ToRemoveCrudFormat(funcName, removeFunc, tableName string) (b string) {
	names, types, _ := m.ExtractPrimaryFieldFormat()
	b = fmt.Sprintf("func %s(%s) (sql.Result, error) {\n", funcName, strings.Join(types, ", "))
	b += fmt.Sprintf("\treturn mysql.Exec(%s(%s, %s))\n}", removeFunc, strings.Join(names, ", "), tableName)
	return
}

func (m *MetadataTable) ToUpdateCrudFormat(funcName, updateFunc, tableName string) (b string) {
	names, types, _ := m.ExtractPrimaryFieldFormat()
	b = fmt.Sprintf("func %s(command string, %s) (sql.Result, error) {\n", funcName, strings.Join(types, ", "))
	b += fmt.Sprintf("\treturn mysql.Exec(%s(command, %s, %s))\n}", updateFunc, strings.Join(names, ", "), tableName)
	return
}

func (m *MetadataTable) ToSelectCrudFormat(funcName, queryFunc, selectFunc, structName, tableName string) (b string) {
	b = fmt.Sprintf("func %s() []*%s {\n", funcName, structName)
	b += fmt.Sprintf("\treturn %s(%s(%s))\n}", queryFunc, selectFunc, tableName)
	return
}

func (m *MetadataTable) ToInsertCrudFormat(funcName, insertFunc, structPrefix, structName, tableName string) (b string) {
	b = fmt.Sprintf("func %s(%s *%s) (sql.Result, error) {\n", funcName, structPrefix, structName)
	b += fmt.Sprintf("\treturn mysql.Exec(%s(%s, %s))\n}", insertFunc, structPrefix, tableName)
	return
}

func (m *MetadataTable) ToWhereCrudFormat(funcName, queryFunc, whereFunc, structName, tableName string) (b string) {
	b = fmt.Sprintf("func %s(command string) []*%s {\n", funcName, structName)
	b += fmt.Sprintf("\treturn %s(%s(command, %s))\n}", queryFunc, whereFunc, tableName)
	return
}

func (m *MetadataTable) ToUpdateTickerCrudFormat(funcName, queryFunc, updateTickerFunc, structName, tableName string) (b string) {
	b = fmt.Sprintf("func %s(updated_at string) []*%s {\n", funcName, structName)
	b += fmt.Sprintf("\treturn %s(%s(updated_at, %s))\n}", queryFunc, updateTickerFunc, tableName)
	return
}

func (m *MetadataTable) ToSetCrudFormat(funcPrefix, setPrefix, tableName string) (b string) {
	names, types, _ := m.ExtractPrimaryAndUpdateFieldFormat()

	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at", "created_by", "created_at":
			continue
		default:
			if m.Fields[i].PrimaryKey {
				continue
			}
		}
		funcName := funcPrefix + m.Fields[i].ToUpperCase()
		setFunc := setPrefix + m.Fields[i].ToUpperCase()
		b += fmt.Sprintf("\nfunc %s(%s %s, %s) (sql.Result, error) {\n", funcName, m.Fields[i].Name, m.Fields[i].TypeOf(), strings.Join(types, ", "))
		b += fmt.Sprintf("\treturn mysql.Exec(%s(%s, %s, %s))\n}\n", setFunc, m.Fields[i].Name, strings.Join(names, ", "), tableName)
	}

	return
}
