package sqlparser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type frontendField struct {
	Label          string `json:"label"`
	Field          string `json:"field"`
	RenderType     string `json:"renderType"`
	Hidden         bool   `json:"hidden"`
	Visible        bool   `json:"visible"`
	Writable       bool   `json:"writable"`
	UpdateWritable bool   `json:"updateWritable"`
	UpdateVisible  bool   `json:"updateVisible"`
}

func (m *MetadataTable) ToFrontendColumnsFormat(columnsName string) string {
	fieldsLen := len(m.Fields)
	if columnsName == "" {
		columnsName = m.ToUpperCase()
	}
	var elements []*frontendField
	for i := 0; i < fieldsLen; i++ {
		element := &frontendField{Label: m.Fields[i].Comment, Field: m.Fields[i].Name, RenderType: "Input"}
		switch m.Fields[i].Name {
		case "id", "uid", "username":
			element.Writable = true
			element.UpdateVisible = true
		case "password":
			element.Hidden = true
			element.Writable = true
		case "status", "deleted", "created_by", "updated_by", "created_at", "updated_at":
			element.Hidden = true
			element.Visible = true
			element.UpdateVisible = true
		default:
			element.Visible = true
			element.Writable = true
			element.UpdateWritable = true
			element.UpdateVisible = true
		}
		elements = append(elements, element)
	}

	b, _ := json.MarshalIndent(elements, "", "    ")
	return fmt.Sprintf("export const %s = %s;", columnsName, string(b))
}

func (m *MetadataTable) ToForntendUnmarshalJSONFormat(funcPrefix, structName, elementName string) (b string) {
	fieldsLen := len(m.Fields)
	funcName := GenerateFunctionName(funcPrefix, m.Name)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].DataType {
		case "TINYINT", "SMALLINT", "MEDIUMINT":
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, elementName, m.Fields[i].ToUpperCase(), "toInt(val)"))
		case "INT", "BIGINT":
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, elementName, m.Fields[i].ToUpperCase(), "toInt64(val)"))
		case "FLOAT", "DOUBLE", "DECIMAL":
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, elementName, m.Fields[i].ToUpperCase(), "toFloat(val)"))
		default:
			elements = append(elements, fmt.Sprintf("\t\tcase \"%s\":\n\t\t\t%s.%s = %s", m.Fields[i].Name, elementName, m.Fields[i].ToUpperCase(), "val"))
		}
	}

	return fmt.Sprintf("func %s(m map[string]interface{}) *%s {\n\t%s := &%s{}\n\tfor k, v := range m {\n\t\tval := strings.TrimSpace(fmt.Sprintf(\"%%v\", v))\n\t\tswitch k {\n%s\n\t\t}\n\t}\n\treturn %s\n}", funcName, structName, elementName, structName, strings.Join(elements, "\n"), elementName)
}