package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	createModelFuncFormat = "ZnVuYyBmdW5jTmFtZShwYXJhbXNGaWVsZCkgKHNxbC5SZXN1bHQsIGVycm9yKSB7CglyZXR1cm4gZmllbGROYW1lLkluc2VydCgmc3RydWN0TmFtZXsKY29udGVudEZpZWxkCgl9KQp9"
	compareModelFuncFormat = "ZnVuYyBmdW5jTmFtZShwYXJhbXNGaWVsZCwgZWxlbWVudCAqc3RydWN0TmFtZSkgc3RyaW5nIHsKICAgIHZhciBjb21tYW5kIFtdc3RyaW5nCmNvbnRlbnRGaWVsZAogICAgaWYgY29tbWFuZCA9PSBuaWwgfHwgbGVuKGNvbW1hbmQpID09IDAgewogICAgICAgIHJldHVybiAiIgogICAgfQoKICAgIHJldHVybiBzdHJpbmdzLkpvaW4oY29tbWFuZCwgIiwgIikKfQ"
	compareModelSubFuncFormat = "ICAgIGlmIGZpZWxkTmFtZSAhPSBlbGVtZW50LnN0cnVjdEZpZWxkIHsKICAgICAgICBjb21tYW5kID0gYXBwZW5kKGNvbW1hbmQsIGZtdC5TcHJpbnRmKGBmaWVsZE5hbWU9ZmllbGRUeXBlYCwgZmllbGROYW1lKSkKICAgIH0"
	updateModelFuncFormat = "ZnVuYyBmdW5jTmFtZShwYXJhbXNGaWVsZCkgKHNxbC5SZXN1bHQsIGVycm9yKSB7CmNvbnRlbnRGaWVsZAogICAgcmV0dXJuIGZpZWxkTmFtZS51cGRhdGVGdW5jKHBhcmFtc1VwZGF0ZSkKfQ"
	updateModelSubFuncFormat = "CWNvbW1hbmQgKz0gZm10LlNwcmludGYoYCwgZmllbGROYW1lPWZpZWxkVHlwZWAsIGZpZWxkTmFtZSk"
	removeModelFuncFormat = "ZnVuYyBmdW5jTmFtZShpZCBpbnQpIChzcWwuUmVzdWx0LCBlcnJvcikgewoJcmV0dXJuIGZpZWxkTmFtZS5zdWJGdW5jKGlkKQp9"
	selectModelFuncFormat = "ZnVuYyBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSByZXN1bHRGaWVsZCB7CglyZXR1cm4gc3RydWN0RmllbGQuc3ViRnVuYyhmaWVsZE5hbWUpCn0"
)

func toCreateModelFuncFormat(funcName, structName, paramsField, fieldName, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(createModelFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "paramsField", paramsField, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func (m *MetadataTable)ToCreateModelFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + m.ToUpperCase()
	funcName := funcPrefix + m.ToUpperCase()
	fieldsLen := len(m.Fields)
	var params []string
	var elements []string
	createdBy := false
	createdAt := false
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "id":
			continue
		case "status", "deleted":
			elements = append(elements, fmt.Sprintf("\t\t%v: 0,", m.Fields[i].ToUpperCase()))
			continue
		case "created_by", "updated_by":
			createdBy = true
			elements = append(elements, fmt.Sprintf("\t\t%v: %v,", m.Fields[i].ToUpperCase(), "created_by"))
			continue
		case "created_at", "updated_at":
			createdAt = true
			elements = append(elements, fmt.Sprintf("\t\t%v: %v,", m.Fields[i].ToUpperCase(), "created_at"))
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int"))
			elements = append(elements, fmt.Sprintf("\t\t%v: %v,", m.Fields[i].ToUpperCase(), m.Fields[i].Name))
		case "BIGINT":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int64"))
			elements = append(elements, fmt.Sprintf("\t\t%v: %v,", m.Fields[i].ToUpperCase(), m.Fields[i].Name))
		default:
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "string"))
			elements = append(elements, fmt.Sprintf("\t\t%v: %v,", m.Fields[i].ToUpperCase(), m.Fields[i].Name))
		}
	}
	if createdBy {
		params = append(params, fmt.Sprintf("%v %v", "created_by", "string"))
	}
	if createdAt {
		params = append(params, fmt.Sprintf("%v %v", "created_at", "string"))
	}

	return toCreateModelFuncFormat(funcName, structName, strings.Join(params, ", "), m.ToLowerCase(), strings.Join(elements, "\n"))
}

func toCompareModelFuncFormat(funcName, structName, paramsField, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareModelFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "paramsField", paramsField, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func toCompareModelSubFuncFormat(structField, fieldName, fieldType string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareModelSubFuncFormat)
	b = strings.Replace(string(fieldFormat), "fieldName", fieldName, -1)
	b = strings.Replace(b, "structField", structField, -1)
	return strings.Replace(b, "fieldType", fieldType, -1)
}

func (m *MetadataTable)ToCompareModelFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + m.ToUpperCase()
	funcName := funcPrefix + m.ToUpperCase()
	fieldsLen := len(m.Fields)
	var params []string
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "id", "username", "uid", "created_by", "updated_by", "created_at", "updated_at":
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int"))
			elements = append(elements, toCompareModelSubFuncFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "%v"))
		case "BIGINT":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int64"))
			elements = append(elements, toCompareModelSubFuncFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, "%v"))
		default:
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "string"))
			elements = append(elements, toCompareModelSubFuncFormat(m.Fields[i].ToUpperCase(), m.Fields[i].Name, `"%v"`))
		}
	}

	return toCompareModelFuncFormat(funcName, structName, strings.Join(params, ", "), strings.Join(elements, "\n"))
}

func toUpdateModelFuncFormat(funcName, paramsField, fieldName, updateFunc, contentField, paramsUpdate string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateModelFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "paramsField", paramsField, -1)
	b = strings.Replace(b, "contentField", contentField, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "updateFunc", updateFunc, -1)
	return strings.Replace(b, "paramsUpdate", paramsUpdate, -1)
}

func toUpdateModelSubFuncFormat(fieldName, fieldType string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateModelSubFuncFormat)
	b = strings.Replace(string(fieldFormat), "fieldName", fieldName, -1)
	return strings.Replace(b, "fieldType", fieldType, -1)
}

func (m *MetadataTable)ToUpdateModelFuncFormat(funcPrefix, updateFunc string) (b string) {
	funcName := funcPrefix + m.ToUpperCase()
	fieldsLen := len(m.Fields)
	var params []string
	var elements []string
	var param string
	for i := 0; i < fieldsLen; i++ {
		switch m.Fields[i].Name {
		case "id", "updated_by", "updated_at":
		default:
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int"))
			if m.Fields[i].Name != "id" {
				elements = append(elements, toUpdateModelSubFuncFormat(m.Fields[i].Name, "%v"))
			}
		case "BIGINT":
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "int64"))
			if m.Fields[i].Name != "id" {
				elements = append(elements, toUpdateModelSubFuncFormat(m.Fields[i].Name, "%v"))
			}
		default:
			params = append(params, fmt.Sprintf("%v %v", m.Fields[i].Name, "string"))
			if m.Fields[i].Name != "id" {
				elements = append(elements, toUpdateModelSubFuncFormat(m.Fields[i].Name, `"%v"`))
			}
		}
	}
	if params == nil || len(params) == 0 {
		param = "command string"
	} else {
		param = strings.Join(params, ", ") + ", command string"
	}
 
	return toUpdateModelFuncFormat(funcName, param, m.ToLowerCase(), updateFunc, strings.Join(elements, "\n"), "command, id")
}

func toRemoveModelFuncFormat(funcName, fieldName, subFunc string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(removeModelFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	return strings.Replace(b, "subFunc", subFunc, -1)
}

func (m *MetadataTable)ToRemoveModelFuncFormat(funcPrefix, removeFunc string)(b string) {
	funcName := funcPrefix + m.ToUpperCase()
	return toRemoveModelFuncFormat(funcName, m.ToLowerCase(), removeFunc)
}

func toSelectModelFuncFormat(funcName, fieldName, fieldType, resultField, structField, subFunc string)(b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectModelFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "fieldType", fieldType, -1)
	b = strings.Replace(b, "resultField", resultField, -1)
	b = strings.Replace(b, "structField", structField, -1)
	return strings.Replace(b, "subFunc", subFunc, -1)
}

func (m *MetadataTable)ToSelectModelFuncFormat(funcPrefix, structPrefix string) (b string) {
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].Name != "id" && m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}
		funcName := "From" + m.ToUpperCase() + funcPrefix + m.Fields[i].ToUpperCase()
		resultField := "[]*" + structPrefix + m.ToUpperCase()
		subFunc := funcPrefix + m.Fields[i].ToUpperCase()
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, toSelectModelFuncFormat(funcName, m.Fields[i].Name, "int", resultField, m.ToLowerCase(), subFunc))
		case "BIGINT":
			elements = append(elements, toSelectModelFuncFormat(funcName, m.Fields[i].Name, "int64", resultField, m.ToLowerCase(), subFunc))
		default:
			elements = append(elements, toSelectModelFuncFormat(funcName, m.Fields[i].Name, "string", resultField, m.ToLowerCase(), subFunc))
		}
	}

	return strings.Join(elements, "\n\n")
}