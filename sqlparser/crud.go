package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	insertFormat      = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CglyZXR1cm4gZm10LlNwcmludGYoYElOU0VSVCBJTlRPICV2KGtleXNGaWVsZCkgVkFMVUVTKHZhbHVlc0ZpZWxkKWAsIHRhYmxlLCBlbGVtZW50c0ZpZWxkKQp9"
	queryFormat       = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKCWRhdGEsIGxlbmd0aCA6PSBteXNxbC5RdWVyeShjb21tYW5kKQoJaWYgZGF0YSA9PSBuaWwgfHwgbGVuZ3RoIDw9IDAgewoJCXJldHVybgoJfQoJYiA6PSAqZGF0YQoJZm9yIGkgOj0gMDsgaSA8IGxlbmd0aDsgaSsrIHsKCQllbGVtZW50IDo9IHBhcnNlcihiW2ldKQoJCWVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCBlbGVtZW50KQoJfQoJcmV0dXJuCn0"
	parserFormat      = "ZnVuYyBmdW5jTmFtZSh2YWx1ZXNGaWVsZCBtYXBbc3RyaW5nXXN0cmluZykgKnN0cnVjdE5hbWUgewoJcmV0dXJuICZzdHJ1Y3ROYW1lewoJCWNvbnRlbnRGaWVsZAoJfQp9"
	selectFormat   = "ZnVuYyBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CglyZXR1cm4gZm10LlNwcmludGYoYFNFTEVDVCAqIEZST00gJXYgV0hFUkUgZmllbGROYW1lPXZhbHVlRmllbGRgLCB0YWJsZSwgZmllbGROYW1lKQp9"
	publicFormat  = "ZnVuYyBpbnNlcnRGdW5jKGVsZW1lbnQgKnN0cnVjdE5hbWUpIChzcWwuUmVzdWx0LCBlcnJvcikgewogICAgcmV0dXJuIG15c3FsLkV4ZWMoaW5zZXJ0KGVsZW1lbnQsIHRhYmxlTmFtZSkpCn0KCmZ1bmMgc2VsZWN0RnVuYygpIFtdKnN0cnVjdE5hbWUgewogICAgcmV0dXJuIHF1ZXJ5KGRhdGFiYXNlLlNlbGVjdFRhYmxlKHRhYmxlTmFtZSkpCn0KCmZ1bmMgdXBkYXRlRnVuYyhjb21tYW5kIHN0cmluZywgaWQgaW50KSAoc3FsLlJlc3VsdCwgZXJyb3IpIHsKICAgIHJldHVybiBteXNxbC5FeGVjKGRhdGFiYXNlLlVwZGF0ZShjb21tYW5kLCBpZCwgdGFibGVOYW1lKSkKfQ"
	removeCrudFormat = "ZnVuYyBmdW5jTmFtZShpZCBpbnQpIChzcWwuUmVzdWx0LCBlcnJvcikgewoJcmV0dXJuIG15c3FsLkV4ZWMoc3ViRnVuYyhpZCwgc3RydWN0TmFtZSkpCn0"
	publicSubFormat = "ZnVuYyBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSBbXSpzdHJ1Y3ROYW1lIHsKCXJldHVybiBxdWVyeShzdWJGdW5jKGZpZWxkTmFtZSwgdGFibGVOYW1lKSkKfQ"
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
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf(`%v:database.ParseInt(%v["%v"]),`, toFieldUpperFormat(m.Fields[i].Name), valuesField, m.Fields[i].Name))
		case "BIGINT":
			elements = append(elements, fmt.Sprintf(`%v:database.ParseInt64(%v["%v"]),`, toFieldUpperFormat(m.Fields[i].Name), valuesField, m.Fields[i].Name))
		default:
			elements = append(elements, fmt.Sprintf(`%v:%v["%v"],`, toFieldUpperFormat(m.Fields[i].Name), valuesField, m.Fields[i].Name))
		}
	}

	contentField := strings.Join(elements, "\n\t\t")
	return toParserFormat(contentField, valuesField, structName, funcName)
}

func (m *MetadataTable)ToSelectFuncFormat(funcName string) (b string) {
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].Name != "id" && m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}
		b += "\n\n"
		b += toSelectFuncFormat(m.Fields[i].Name, m.Fields[i].DataType, funcName)
	}
	return
}

func toSelectFuncFormat(name, dateType, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName + toFieldUpperFormat(name), -1)
	b = strings.Replace(b, "fieldName", name, -1)

	switch dateType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
		b = strings.Replace(b, "fieldType", "int", -1)
		b = strings.Replace(b, "valueField", `%v`, -1)
	case "BIGINT":
		b = strings.Replace(b, "fieldType", "int64", -1)
		b = strings.Replace(b, "valueField", `%v`, -1)
	default:
		b = strings.Replace(b, "fieldType", "", -1)
		b = strings.Replace(b, "valueField", `"%v"`, -1)
	}
	return
}

func (m *MetadataTable)ToPublicCrudFormat(insertFunc, selectFunc, updateFunc, structPrefix, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(publicFormat)
	b = strings.Replace(string(fieldFormat), "insertFunc", insertFunc, -1)
	b = strings.Replace(b, "selectFunc", selectFunc, -1)
	b = strings.Replace(b, "updateFunc", updateFunc, -1)
	b = strings.Replace(b, "structName", structPrefix + toFieldUpperFormat(m.Name), -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	return
}

func (m *MetadataTable)ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName string)(b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].Name != "id" && m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}
		subFunc := subPrefix + toFieldUpperFormat(m.Fields[i].Name)
		funcName := funcPrefix + toFieldUpperFormat(m.Fields[i].Name)
		b += "\n\n"
		b += toPublicSubCrudFormat(funcName, m.Fields[i].Name, m.Fields[i].DataType, subFunc, structName, tableName)
	}

	return
}

func toPublicSubCrudFormat(funcName, fieldName, fieldType, subFunc, structName, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(publicSubFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	switch fieldType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
		b = strings.Replace(b, "fieldType", `int`, -1)
	case "BIGINT":
		b = strings.Replace(b, "fieldType", `int64`, -1)
	default:
		b = strings.Replace(b, "fieldType", `string`, -1)
	}
	return
}

func toRemoveCrudFormat(funcName, subFunc, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(removeCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable)ToRemoveCrudFormat(funcName, structPrefix string)(b string) {
	structName := structPrefix + m.ToUpperCase()
	subFunc := structPrefix + funcName
	return toRemoveCrudFormat(funcName, subFunc, structName)
}