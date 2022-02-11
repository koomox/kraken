package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	insertFormat     = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CglyZXR1cm4gZm10LlNwcmludGYoYElOU0VSVCBJTlRPICV2KGtleXNGaWVsZCkgVkFMVUVTKHZhbHVlc0ZpZWxkKWAsIHRhYmxlLCBlbGVtZW50c0ZpZWxkKQp9"
	queryFormat      = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKCWRhdGEsIGxlbmd0aCA6PSBteXNxbC5RdWVyeShjb21tYW5kKQoJaWYgZGF0YSA9PSBuaWwgfHwgbGVuZ3RoIDw9IDAgewoJCXJldHVybgoJfQoJYiA6PSAqZGF0YQoJZm9yIGkgOj0gMDsgaSA8IGxlbmd0aDsgaSsrIHsKCQllbGVtZW50IDo9IHBhcnNlcihiW2ldKQoJCWVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCBlbGVtZW50KQoJfQoJcmV0dXJuCn0"
	parserFormat     = "ZnVuYyBmdW5jTmFtZSh2YWx1ZXNGaWVsZCBtYXBbc3RyaW5nXXN0cmluZykgKnN0cnVjdE5hbWUgewoJcmV0dXJuICZzdHJ1Y3ROYW1lewoJCWNvbnRlbnRGaWVsZAoJfQp9"
	selectFormat     = "ZnVuYyBmdW5jTmFtZShrZXlOYW1lIGtleVR5cGUsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgU0VMRUNUICogRlJPTSAldiBXSEVSRSBrZXlOYW1lPXZhbHVlRmllbGRgLCB0YWJsZSwga2V5TmFtZSkKfQ"
	insertCrudFormat = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lKSAoc3FsLlJlc3VsdCwgZXJyb3IpIHsKCXJldHVybiBteXNxbC5FeGVjKGluc2VydChlbGVtZW50LCB0YWJsZU5hbWUpKQp9"
	selectCrudFormat = "ZnVuYyBmdW5jTmFtZSgpIHN0cnVjdE5hbWUgewoJcmV0dXJuIHF1ZXJ5KHN1YkZ1bmModGFibGVOYW1lKSkKfQ"
	updateCrudFormat = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZywgaWQgZmllbGRUeXBlKSAoc3FsLlJlc3VsdCwgZXJyb3IpIHsKICAgIHJldHVybiBteXNxbC5FeGVjKHN1YkZ1bmMoY29tbWFuZCwgaWQsIHRhYmxlTmFtZSkpCn0"
	removeCrudFormat = "ZnVuYyBmdW5jTmFtZShwYXJhbXNGaWVsZCkgKHNxbC5SZXN1bHQsIGVycm9yKSB7CiAgICByZXR1cm4gbXlzcWwuRXhlYyhzdWJGdW5jKHZhbHVlc0ZpZWxkLCB0YWJsZU5hbWUpKQp9"
	whereCrudFormat  = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgc3RydWN0TmFtZSB7CglyZXR1cm4gcXVlcnkoc3ViRnVuYyhjb21tYW5kLCB0YWJsZU5hbWUpKQp9"
	publicSubFormat  = "ZnVuYyBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSBbXSpzdHJ1Y3ROYW1lIHsKCXJldHVybiBxdWVyeShzdWJGdW5jKGZpZWxkTmFtZSwgdGFibGVOYW1lKSkKfQ"
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
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
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

func (m *MetadataTable) ToSelectFuncFormat(funcName string) (b string) {
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		if m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}
		b += "\n\n"
		b += toSelectFuncFormat(m.Fields[i].Name, m.Fields[i].DataType, funcName)
	}
	return
}

func toSelectFuncFormat(name, dateType, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName+toFieldUpperFormat(name), -1)
	b = strings.Replace(b, "keyName", name, -1)

	switch dateType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
		b = strings.Replace(b, "keyType", "int", -1)
		b = strings.Replace(b, "valueField", `%v`, -1)
	case "BIGINT":
		b = strings.Replace(b, "keyType", "int64", -1)
		b = strings.Replace(b, "valueField", `%v`, -1)
	default:
		b = strings.Replace(b, "keyType", "", -1)
		b = strings.Replace(b, "valueField", `"%v"`, -1)
	}
	return
}

func (m *MetadataTable) ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		if m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
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

func toRemoveCrudFormat(funcName, paramsField, subFunc, valuesField, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(removeCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	b = strings.Replace(b, "valuesField", valuesField, -1)
	return strings.Replace(b, "paramsField", paramsField, -1)
}

func (m *MetadataTable) ToRemoveCrudFormat(funcName, structPrefix, tableName string) (b string) {
	subFunc := structPrefix + funcName
	fieldsLen := len(m.Fields)
	var params []string
	var values []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at":
		default:
			if m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
				continue
			}
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int"))
		case "BIGINT":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int64"))
		default:
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "string"))
		}
		values = append(values, m.Fields[i].Name)
	}

	return toRemoveCrudFormat(funcName, strings.Join(params, ", "), subFunc, strings.Join(values, ", "), tableName)
}

func toUpdateCrudFormat(funcName, subFunc, fieldType, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "fieldType", fieldType, -1)
	return strings.Replace(b, "tableName", tableName, -1)
}

func (m *MetadataTable) ToUpdateCrudFormat(funcName, structPrefix, tableName string) (b string) {
	subFunc := structPrefix + funcName
	fieldType := ""
	if v := m.PrimaryKey(); v != nil {
		fieldType = v.TypeOf()
	}

	return toUpdateCrudFormat(funcName, subFunc, fieldType, tableName)
}

func toSelectCrudFormat(funcName, subFunc, structName, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable) ToSelectCrudFormat(funcName, structPrefix, tableName string) (b string) {
	structName := "[]*" + structPrefix + m.ToUpperCase()
	subFunc := structPrefix + funcName
	return toSelectCrudFormat(funcName, subFunc, structName, tableName)
}

func toInsertCrudFormat(funcName, structName, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(insertCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable) ToInsertCrudFormat(funcName, structPrefix, tableName string) (b string) {
	structName := structPrefix + m.ToUpperCase()
	return toInsertCrudFormat(funcName, structName, tableName)
}

func toWhereCrudFormat(funcName, subFunc, structName, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(whereCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return strings.Replace(b, "subFunc", subFunc, -1)
}

func (m *MetadataTable) ToWhereCrudFormat(funcName, structPrefix, tableName string) (b string) {
	structName := "[]*" + structPrefix + m.ToUpperCase()
	subFunc := structPrefix + funcName
	return toWhereCrudFormat(funcName, subFunc, structName, tableName)
}
