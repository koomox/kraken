package sqlparser

import (
	"fmt"
	"strings"
)

func toLabelFormat(label, field, hidden, visible, writable, updateWritable, updateVisible string) (b string) {
	return fmt.Sprintf("  {\n    label: '%s',\n    field: '%s',\n    renderType: 'Input',\n    hidden: %s,\n    visible: %s,\n    writable: %s,\n    updateWritable: %s,\n    updateVisible: %s,\n  },", label, field, hidden, visible, writable, updateWritable, updateVisible)
}

func (m *MetadataTable) ToFrontendColumnsFormat(columnsName string) (b string) {
	fieldsLen := len(m.Fields)
	if columnsName == "" {
		columnsName = m.ToUpperCase()
	}
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "id", "uid", "username":
			elements = append(elements, toLabelFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "false", "false", "true", "false", "true"))
		case "password":
			elements = append(elements, toLabelFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "true", "false", "true", "false", "false"))
		case "status", "deleted", "created_by", "updated_by", "created_at", "updated_at":
			elements = append(elements, toLabelFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "true", "true", "false", "false", "true"))
		default:
			elements = append(elements, toLabelFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "false", "true", "true", "true", "true"))
		}
	}

	return fmt.Sprintf("const %s = [\n%s\n];", columnsName, strings.Join(elements, "\n"))
}

func (m *MetadataTable) ToForntendParseFormat(funcName, structName, element string) (b string) {
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, element, m.Fields[i].ToUpperCase(), "database.ParseInt(val)"))
		case "BIGINT":
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, element, m.Fields[i].ToUpperCase(), "database.ParseInt64(val)"))
		default:
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, element, m.Fields[i].ToUpperCase(), "val"))
		}
	}

	return fmt.Sprintf("func %s(m map[string]interface{}) (%s *%s) {\n\t%s = &%s{}\n\tfor k, v := range m {\n\t\tval := strings.TrimSpace(fmt.Sprintf(\"%%v\", v))\n\t\tswitch k {\n%s\n\t\t}\n\t}\nreturn\n}", funcName, element, structName, element, structName, strings.Join(elements, "\n"))
}