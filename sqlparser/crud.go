package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	insertFormat     = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CglyZXR1cm4gZm10LlNwcmludGYoYElOU0VSVCBJTlRPICV2KGtleXNGaWVsZCkgVkFMVUVTKHZhbHVlc0ZpZWxkKWAsIHRhYmxlLCBlbGVtZW50c0ZpZWxkKQp9"
	updateFormat     = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZywgYXJnc0ZpZWxkLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CiAgICByZXR1cm4gZm10LlNwcmludGYoYFVQREFURSAldiBTRVQgJXYgV0hFUkUgZm9ybWF0RmllbGRgLCB0YWJsZSwgY29tbWFuZCwga2V5c0ZpZWxkKQp9"
	removeFormat     = "ZnVuYyBmdW5jTmFtZShhcmdzRmllbGQsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgVVBEQVRFICV2IFNFVCBkZWxldGVkPTEsIHVwZGF0ZWRfYnk9JXYsIHVwZGF0ZWRfYXQ9IiV2IiBXSEVSRSBmb3JtYXRGaWVsZGAsIHRhYmxlLCBrZXlzRmllbGQpCn0"
	queryFormat      = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKCWRhdGEsIGxlbmd0aCA6PSBteXNxbC5RdWVyeShjb21tYW5kKQoJaWYgZGF0YSA9PSBuaWwgfHwgbGVuZ3RoIDw9IDAgewoJCXJldHVybgoJfQoJYiA6PSAqZGF0YQoJZm9yIGkgOj0gMDsgaSA8IGxlbmd0aDsgaSsrIHsKCQllbGVtZW50IDo9IHBhcnNlcihiW2ldKQoJCWVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCBlbGVtZW50KQoJfQoJcmV0dXJuCn0"
	parserFormat     = "ZnVuYyBmdW5jTmFtZSh2YWx1ZXNGaWVsZCBtYXBbc3RyaW5nXXN0cmluZykgKnN0cnVjdE5hbWUgewoJcmV0dXJuICZzdHJ1Y3ROYW1lewoJCWNvbnRlbnRGaWVsZAoJfQp9"
	selectFormat     = "ZnVuYyBmdW5jTmFtZShhcmdzRmllbGQsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgU0VMRUNUICogRlJPTSAldiBXSEVSRSBmb3JtYXRGaWVsZGAsIHRhYmxlLCBrZXlzRmllbGQpCn0"
	insertCrudFormat = "ZnVuYyBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lKSAoc3FsLlJlc3VsdCwgZXJyb3IpIHsKCXJldHVybiBteXNxbC5FeGVjKGluc2VydChlbGVtZW50LCB0YWJsZU5hbWUpKQp9"
	selectCrudFormat = "ZnVuYyBmdW5jTmFtZSgpIHN0cnVjdE5hbWUgewoJcmV0dXJuIHF1ZXJ5KHN1YkZ1bmModGFibGVOYW1lKSkKfQ"
	updateCrudFormat = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZywgYXJnc0ZpZWxkKSAoc3FsLlJlc3VsdCwgZXJyb3IpIHsKICAgIHJldHVybiBteXNxbC5FeGVjKHN1YkZ1bmMoY29tbWFuZCwga2V5c0ZpZWxkLCB0YWJsZU5hbWUpKQp9"
	removeCrudFormat = "ZnVuYyBmdW5jTmFtZShhcmdzRmllbGQpIChzcWwuUmVzdWx0LCBlcnJvcikgewogICAgcmV0dXJuIG15c3FsLkV4ZWMoc3ViRnVuYyhrZXlzRmllbGQsIHRhYmxlTmFtZSkpCn0"
	whereCrudFormat  = "ZnVuYyBmdW5jTmFtZShjb21tYW5kIHN0cmluZykgc3RydWN0TmFtZSB7CglyZXR1cm4gcXVlcnkoc3ViRnVuYyhjb21tYW5kLCB0YWJsZU5hbWUpKQp9"
	publicSubFormat  = "ZnVuYyBmdW5jTmFtZShhcmdzRmllbGQpIFtdKnN0cnVjdE5hbWUgewogICAgcmV0dXJuIHF1ZXJ5KHN1YkZ1bmMoa2V5c0ZpZWxkLCB0YWJsZU5hbWUpKQp9"
)

func (m *MetadataTable) ToInsertFormat(structPrefix, funcName string) (b string) {
	return m.toInsertFormat(structPrefix, funcName)
}

func (m *MetadataTable) ToUpdateFormat(funcName string) (b string) {
	var args []string
	var keys []string
	var format []string
	counter := 0
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			counter++
			keys = append(keys, m.Fields[i].Name)
			args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
			format = append(format, fmt.Sprintf(`%v=%v`, m.Fields[i].Name, m.Fields[i].ValueOf()))
		}
	}
	if counter == 0 {
		return
	}
	return toUpdateFormat(funcName, strings.Join(args, ", "), strings.Join(format, " AND "), strings.Join(keys, ", "))
}

func (m *MetadataTable) ToRemoveFormat(funcName string) (b string) {
	var args []string
	var keys []string
	var format []string
	counter := 0
	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at":
			keys = append(keys, m.Fields[i].Name)
			args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
		default:
			if m.Fields[i].PrimaryKey {
				counter++
				keys = append(keys, m.Fields[i].Name)
				args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
				format = append(format, fmt.Sprintf(`%v=%v`, m.Fields[i].Name, m.Fields[i].ValueOf()))
			}
		}

	}
	if counter == 0 {
		return
	}
	keysField := strings.Join(keys[len(keys)-2:], ", ")
	keysField += ", " + strings.Join(keys[:len(keys)-2], ", ")
	return toRemoveFormat(funcName, strings.Join(args, ", "), strings.Join(format, " AND "), keysField)
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

func toUpdateFormat(funcName, argsField, formatField, keysField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "formatField", formatField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	return
}

func toRemoveFormat(funcName, argsField, formatField, keysField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(removeFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "formatField", formatField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
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
		if m.Fields[i].AutoIncrment {
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			values = append(values, "%v")
		default:
			values = append(values, `"%v"`)
		}
		keys = append(keys, m.Fields[i].Name)
		elements = append(elements, elementPrefix+m.Fields[i].ToUpperCase())
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
	var args []string
	var keys []string
	var format []string
	var subFunc []string
	counter := 0
	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "created_by":
			b += "\n\n"
			b += toSelectFuncFormat(funcName+m.Fields[i].ToUpperCase(), fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()), fmt.Sprintf("%v=%v", m.Fields[i].Name, m.Fields[i].ValueOf()), m.Fields[i].Name)
		default:
			if m.Fields[i].PrimaryKey {
				counter++
				subFunc = append(subFunc, m.Fields[i].ToUpperCase())
				keys = append(keys, m.Fields[i].Name)
				args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
				format = append(format, fmt.Sprintf("%v=%v", m.Fields[i].Name, m.Fields[i].ValueOf()))
			}
			if m.Fields[i].PrimaryKey || m.Fields[i].Unique {
				b += "\n\n"
				b += toSelectFuncFormat(funcName+m.Fields[i].ToUpperCase(), fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()), fmt.Sprintf("%v=%v", m.Fields[i].Name, m.Fields[i].ValueOf()), m.Fields[i].Name)
			}
		}
	}
	if counter < 2 {
		return
	}
	b += "\n\n"
	b += toSelectFuncFormat(funcName+strings.Join(subFunc, "And"), strings.Join(args, ", "), strings.Join(format, " AND "), strings.Join(keys, ", "))
	return
}

func toSelectFuncFormat(funcName, argsField, formatField, keysField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "formatField", formatField, -1)
	return strings.Replace(b, "keysField", keysField, -1)
}

func (m *MetadataTable) ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	var args []string
	var keys []string
	var subFunc []string
	counter := 0
	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "created_by":
			b += "\n\n"
			b += toPublicSubCrudFormat(funcPrefix+m.Fields[i].ToUpperCase(), subPrefix+m.Fields[i].ToUpperCase(), fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()), m.Fields[i].Name, structName, tableName)
		default:
			if m.Fields[i].PrimaryKey {
				counter++
				subFunc = append(subFunc, m.Fields[i].ToUpperCase())
				keys = append(keys, m.Fields[i].Name)
				args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
			}
			if m.Fields[i].PrimaryKey || m.Fields[i].Unique {
				b += "\n\n"
				b += toPublicSubCrudFormat(funcPrefix+m.Fields[i].ToUpperCase(), subPrefix+m.Fields[i].ToUpperCase(), fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()), m.Fields[i].Name, structName, tableName)
			}
		}
	}
	if counter < 2 {
		return
	}
	b += "\n\n"
	b += toPublicSubCrudFormat(funcPrefix+strings.Join(subFunc, "And"), subPrefix+strings.Join(subFunc, "And"), strings.Join(args, ", "), strings.Join(keys, ", "), structName, tableName)
	return
}

func toPublicSubCrudFormat(funcName, subFunc, argsField, keysField, structName, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(publicSubFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return strings.Replace(b, "tableName", tableName, -1)
}

func toRemoveCrudFormat(funcName, subFunc, argsField, keysField, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(removeCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	return strings.Replace(b, "tableName", tableName, -1)
}

func (m *MetadataTable) ToRemoveCrudFormat(funcName, subFunc, tableName string) (b string) {
	var args []string
	var keys []string
	counter := 0
	for i := range m.Fields {
		switch m.Fields[i].Name {
		case "updated_by", "updated_at":
			keys = append(keys, m.Fields[i].Name)
			args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
		default:
			if m.Fields[i].PrimaryKey {
				counter++
				counter++
				keys = append(keys, m.Fields[i].Name)
				args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
			}
		}
	}
	if counter == 0 {
		return
	}

	return toRemoveCrudFormat(funcName, subFunc, strings.Join(args, ", "), strings.Join(keys, ", "), tableName)
}

func toUpdateCrudFormat(funcName, subFunc, argsField, keysField, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateCrudFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	return strings.Replace(b, "tableName", tableName, -1)
}

func (m *MetadataTable) ToUpdateCrudFormat(funcName, subFunc, tableName string) (b string) {
	var args []string
	var keys []string
	counter := 0
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			counter++
			keys = append(keys, m.Fields[i].Name)
			args = append(args, fmt.Sprintf("%v %v", m.Fields[i].Name, m.Fields[i].TypeOf()))
		}
	}
	if counter == 0 {
		return
	}

	return toUpdateCrudFormat(funcName, subFunc, strings.Join(args, ", "), strings.Join(keys, ", "), tableName)
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
