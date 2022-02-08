package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	columnsFieldFormat = "Y29uc3QgY29sdW1uc0ZpZWxkID0gWwpjb250ZW50RmllbGQKXTs"
	labelFieldFormat = "ICB7CiAgICBsYWJlbDogJ2xhYmVsRmllbGQnLAogICAgZmllbGQ6ICdmaWVsZE5hbWUnLAogICAgcmVuZGVyVHlwZTogJ0lucHV0JywKICAgIGhpZGRlbjogaGlkZGVuRmllbGQsCiAgICB2aXNpYmxlOiB2aXNpYmxlRmllbGQsCiAgICB3cml0YWJsZTogd3JpdGFibGVGaWVsZCwKICAgIHVwZGF0ZVdyaXRhYmxlOiB1cGRhdGVXcml0YWJsZUZpZWxkLAogICAgdXBkYXRlVmlzaWJsZTogdXBkYXRlVmlzaWJsZUZpZWxkLAogIH0s"
	parseFuncFormat = "ZnVuYyBmdW5jTmFtZShtIG1hcFtzdHJpbmddaW50ZXJmYWNle30pIChlbGVtZW50ICpzdHJ1Y3ROYW1lKSB7CgllbGVtZW50ID0gJnN0cnVjdE5hbWV7fQoJZm9yIGssIHYgOj0gcmFuZ2UgbSB7CgkJdmFsIDo9IHN0cmluZ3MuVHJpbVNwYWNlKGZtdC5TcHJpbnRmKCIldiIsIHYpKQoJCXN3aXRjaCBrIHsKY29udGVudEZpZWxkCgkJfQoJfQoKCXJldHVybgp9"
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

func (m *MetadataTable)ToFrontendColumnsFormat()(b string) {

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

func (m *MetadataTable)ToForntendParseFormat() (b string) {
	
}