package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	insertFormat      = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CglyZXR1cm4gZm10LlNwcmludGYoYElOU0VSVCBJTlRPICV2KGtleXNGaWVsZCkgVkFMVUVTKHZhbHVlc0ZpZWxkKWAsIHRhYmxlLCBlbGVtZW50c0ZpZWxkKQp9"
	computedFormat    = "ZnVuYyBmdW5jTmFtZShkLCBzICpzdHJ1Y3ROYW1lKSAoY29tbWFuZCBzdHJpbmcpIHsKCWlmIHMuSWQgIT0gZC5JZCB7CgkJcmV0dXJuCgl9CmNvbnRlbnRGaWVsZAoJaWYgY29tbWFuZCA9PSAiIiB7CgkJcmV0dXJuCgl9CglpZiBzLkNyZWF0ZWRCeSAhPSBkLkNyZWF0ZWRCeSB7CgkJY29tbWFuZCArPSBmbXQuU3ByaW50ZihgY3JlYXRlZF9ieT0ldiwgYCwgZC5DcmVhdGVkQnkpCgl9CglpZiBzLlVwZGF0ZWRCeSAhPSBkLlVwZGF0ZWRCeSB7CgkJY29tbWFuZCArPSBmbXQuU3ByaW50ZihgdXBkYXRlZF9ieT0ldiwgYCwgZC5VcGRhdGVkQnkpCgl9Cgljb21tYW5kICs9IGZtdC5TcHJpbnRmKGB1cGRhdGVkX2F0PSIldiJgLCBleHQuTmV3RGF0ZVRpbWUoIiIpLlN0cmluZygpKQoJcmV0dXJuCn0"
	computedSubFormat = "CWlmIHMuZmllbGROYW1lICE9IGQuZmllbGROYW1lIHsKCQljb21tYW5kICs9IGZtdC5TcHJpbnRmKGB0YWdOYW1lPXZhbHVlRmllbGQsIGAsIGQuZmllbGROYW1lKQoJfQ"
	queryFormat       = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKCWRhdGEsIGxlbmd0aCA6PSBteXNxbC5RdWVyeShjb21tYW5kKQoJaWYgZGF0YSA9PSBuaWwgfHwgbGVuZ3RoIDw9IDAgewoJCXJldHVybgoJfQoJYiA6PSAqZGF0YQoJZm9yIGkgOj0gMDsgaSA8IGxlbmd0aDsgaSsrIHsKCQllbGVtZW50IDo9IHBhcnNlcihiW2ldKQoJCWVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCBlbGVtZW50KQoJfQoJcmV0dXJuCn0"
	parserFormat      = "ZnVuYyBmdW5jTmFtZSh2YWx1ZXNGaWVsZCBtYXBbc3RyaW5nXXN0cmluZykgKnN0cnVjdE5hbWUgewoJcmV0dXJuICZzdHJ1Y3ROYW1lewoJCWNvbnRlbnRGaWVsZAoJfQp9"
)

func (m *MetadataTable) ToInsertFormat(structPrefix, funcName string) (b string) {
	return m.toInsertFormat(structPrefix, funcName)
}

func (m *MetadataTable) ToQueryFormat(structPrefix, funcName string) (b string) {
	return m.toQueryFormat(structPrefix, funcName)
}

func (m *MetadataTable) ToParserFormat(valuesField, structPrefix, funcName string) (b string) {
	return m.toParserFormat(valuesField, structPrefix, funcName)
}

func (m *MetadataTable) ToComputedFormat(structPrefix, funcName string) (b string) {
	return m.toComputedFormat(structPrefix, funcName)
}

func toInsertFormat(keysField, valuesField, elementsField, structName, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(insertFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	b = strings.Replace(b, "elementsField", elementsField, -1)
	b = strings.Replace(b, "valuesField", valuesField, -1)
	return
}

func toQueryFormat(structName, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(queryFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return
}

func toParserFormat(contentField, valuesField, structName, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parserFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "valuesField", valuesField, -1)
	b = strings.Replace(b, "contentField", contentField, -1)
	return
}

func toComputedFormat(contentField, structName, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(computedFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "contentField", contentField, -1)
	return
}

func (m *MetadataTable) toInsertFormat(structPrefix, funcName string) (b string) {
	var keys []string
	var values []string
	var elements []string

	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	elementPrefix := "element."

	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].Name == "id" {
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			values = append(values, "%v")
		default:
			values = append(values, `"%v"`)
		}
		keys = append(keys, m.Fields[i].Name)
		elements = append(elements, elementPrefix+toFieldUpperFormat(m.Fields[i].Name))
	}

	keysField := strings.Join(keys, ", ")
	valuesField := strings.Join(values, ", ")
	elementsField := strings.Join(elements, ", ")

	return toInsertFormat(keysField, valuesField, elementsField, structName, funcName)
}

func (m *MetadataTable) toQueryFormat(structPrefix, funcName string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	return toQueryFormat(structName, funcName)
}

func (m *MetadataTable) toParserFormat(valuesField, structPrefix, funcName string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf(`%v:ext.Atoi(%v["%v"]),`, toFieldUpperFormat(m.Fields[i].Name), valuesField, m.Fields[i].Name))
		default:
			elements = append(elements, fmt.Sprintf(`%v:%v["%v"],`, toFieldUpperFormat(m.Fields[i].Name), valuesField, m.Fields[i].Name))
		}
	}

	contentField := strings.Join(elements, "\n\t\t")
	return toParserFormat(contentField, valuesField, structName, funcName)
}

func (m *MetadataTable) toComputedFormat(structPrefix, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(computedSubFormat)
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "id", "created_at", "updated_at":
			continue
		}
		element := strings.Replace(string(fieldFormat), "fieldName", toFieldUpperFormat(m.Fields[i].Name), -1)
		element = strings.Replace(element, "tagName", m.Fields[i].Name, -1)
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			element = strings.Replace(element, "valueField", `%v`, -1)
		default:
			element = strings.Replace(element, "valueField", `"%v"`, -1)
		}

		elements = append(elements, element)
	}

	contentField := strings.Join(elements, "\n")
	return toComputedFormat(contentField, structName, funcName)
}
