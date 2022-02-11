package sqlparser

import (
	"encoding/base64"
	"strings"
)

const (
	columnsFieldFormat = "Y29uc3QgY29sdW1uc0ZpZWxkID0gWwpjb250ZW50RmllbGQKXTs"
	labelFieldFormat   = "ICB7CiAgICBsYWJlbDogJ2xhYmVsRmllbGQnLAogICAgZmllbGQ6ICdmaWVsZE5hbWUnLAogICAgcmVuZGVyVHlwZTogJ0lucHV0JywKICAgIGhpZGRlbjogaGlkZGVuRmllbGQsCiAgICB2aXNpYmxlOiB2aXNpYmxlRmllbGQsCiAgICB3cml0YWJsZTogd3JpdGFibGVGaWVsZCwKICAgIHVwZGF0ZVdyaXRhYmxlOiB1cGRhdGVXcml0YWJsZUZpZWxkLAogICAgdXBkYXRlVmlzaWJsZTogdXBkYXRlVmlzaWJsZUZpZWxkLAogIH0s"
	parseFuncFormat    = "ZnVuYyBmdW5jTmFtZShtIG1hcFtzdHJpbmddaW50ZXJmYWNle30pIChlbGVtZW50ICpzdHJ1Y3ROYW1lKSB7CgllbGVtZW50ID0gJnN0cnVjdE5hbWV7fQoJZm9yIGssIHYgOj0gcmFuZ2UgbSB7CgkJdmFsIDo9IHN0cmluZ3MuVHJpbVNwYWNlKGZtdC5TcHJpbnRmKCIldiIsIHYpKQoJCXN3aXRjaCBrIHsKY29udGVudEZpZWxkCgkJfQoJfQoKCXJldHVybgp9"
	parseSubFuncFormat = "CQljYXNlICJuYW1lRmllbGQiOgoJCQllbGVtZW50LmZpZWxkTmFtZSA9IHZhbHVlRmllbGQ"
)

func toLabelFormat(labelField, fieldName, hiddenField, visibleField, writableField, updateWritableField, updateVisibleField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(labelFieldFormat)
	b = strings.Replace(string(fieldFormat), "labelField", labelField, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "hiddenField", hiddenField, -1)
	b = strings.Replace(b, "visibleField", visibleField, -1)
	b = strings.Replace(b, "writableField", writableField, -1)
	b = strings.Replace(b, "updateWritableField", updateWritableField, -1)
	return strings.Replace(b, "updateVisibleField", updateVisibleField, -1)
}

func toColumnsFormat(columnsField, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(columnsFieldFormat)
	b = strings.Replace(string(fieldFormat), "columnsField", columnsField, -1)
	return strings.Replace(b, "contentField", contentField, -1)
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

	return toColumnsFormat(columnsName, strings.Join(elements, "\n"))
}

func toParseFuncFormat(funcName, structName, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parseFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func toParseSubFuncFormat(nameField, fieldName, valueField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parseSubFuncFormat)
	b = strings.Replace(string(fieldFormat), "nameField", nameField, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	return strings.Replace(b, "valueField", valueField, -1)
}

func (m *MetadataTable) ToForntendParseFormat(funcPrefix string) (b string) {
	fieldsLen := len(m.Fields)
	funcName := funcPrefix + m.ToUpperCase()
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, toParseSubFuncFormat(m.Fields[i].Name, m.Fields[i].ToUpperCase(), "database.ParseInt(val)"))
		case "BIGINT":
			elements = append(elements, toParseSubFuncFormat(m.Fields[i].Name, m.Fields[i].ToUpperCase(), "database.ParseInt64(val)"))
		default:
			elements = append(elements, toParseSubFuncFormat(m.Fields[i].Name, m.Fields[i].ToUpperCase(), "val"))
		}
	}

	return toParseFuncFormat(funcName, m.ToUpperCase(), strings.Join(elements, "\n"))
}
