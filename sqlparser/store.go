package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	storeStructFormat = "dHlwZSBTdG9yZSBzdHJ1Y3QgewoJc3luYy5SV011dGV4Cgl0YWJsZSAgIHN0cmluZwoJc3RvcmUgICBjb21tb24uQnVja2V0Cm1hcHBpbmdGaWVsZAoJVXBkYXRlZCBib29sCglQYXRjaCAgIFtdaW50ZXJmYWNle30KfQ"
	newStoreFuncFormat = "ZnVuYyBmdW5jTmFtZShjIGNvbW1vbi5CdWNrZXQpIChzdG9yZSAqU3RvcmUpIHsKICAgIHN0b3JlID0gJlN0b3JlewogICAgICAgIHRhYmxlOiAgIHRhYmxlTmFtZSwKICAgICAgICBzdG9yZTogICBjLAptYWtlRmllbGQKICAgICAgICBVcGRhdGVkOiBmYWxzZSwKICAgIH0KICAgIGVsZW1lbnRzIDo9IHF1ZXJ5KHN1YkZ1bmMoc3RvcmUudGFibGUpKQogICAgaWYgZWxlbWVudHMgPT0gbmlsIHx8IGxlbihlbGVtZW50cykgPD0gMCB7CiAgICAgICAgcmV0dXJuCiAgICB9CiAgICBmb3IgaSA6PSByYW5nZSBlbGVtZW50cyB7CiAgICAgICAgZWxlbWVudCA6PSBlbGVtZW50c1tpXQogICAgICAgIHN0b3JlLnN0b3JlLlB1dChlbGVtZW50LmluZGV4RmllbGQsIGVsZW1lbnQpCm1hcHBpbmdGaWVsZAogICAgfQoKICAgIHJldHVybgp9"
	mapStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZSgpIHsKbWFrZUZpZWxkCglzdG9yZS5zdG9yZS5DYWxsYmFja0Z1bmMoZnVuYyh2IGludGVyZmFjZXt9KSB7CgkJaWYgdiAhPSBuaWwgewoJCQllbGVtZW50IDo9IHYuKCpzdHJ1Y3ROYW1lKQptYXBwaW5nRmllbGQKCQl9Cgl9KQpzdG9yZUZpZWxkCn0"
	updateStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShkYXRldGltZSBzdHJpbmcpIHsKICAgIHN0b3JlLlVwZGF0ZWQgPSBmYWxzZQogICAgc3RvcmUuUGF0Y2ggPSBuaWwKICAgIGVsZW1lbnRzIDo9IHF1ZXJ5KHN1YkZ1bmMoZGF0ZXRpbWUsIHN0b3JlLnRhYmxlKSkKICAgIGlmIGVsZW1lbnRzID09IG5pbCB8fCBsZW4oZWxlbWVudHMpIDw9IDAgewogICAgICAgIHJldHVybgogICAgfQogICAgZm9yIGkgOj0gMDsgaSA8IGxlbihlbGVtZW50cyk7IGkrKyB7CiAgICAgICAgZWxlbWVudCA6PSBlbGVtZW50c1tpXQogICAgICAgIGlmICFzdG9yZS5jb21wYXJlRnVuYyhlbGVtZW50KSB7CiAgICAgICAgICAgIHN0b3JlLnN0b3JlLlB1dChlbGVtZW50LmluZGV4RmllbGQsIGVsZW1lbnQpCiAgICAgICAgICAgIHN0b3JlLlBhdGNoID0gYXBwZW5kKHN0b3JlLlBhdGNoLCBlbGVtZW50KQogICAgICAgICAgICBpZiAhc3RvcmUuVXBkYXRlZCB7CiAgICAgICAgICAgICAgICBzdG9yZS5VcGRhdGVkID0gdHJ1ZQogICAgICAgICAgICB9CiAgICAgICAgfQogICAgfQptYXBGaWVsZAp9"
	compareStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lKSBib29sIHsKICAgIGlmIHYgOj0gc3RvcmUuc3RvcmUuR2V0KGVsZW1lbnQuaW5kZXhGaWVsZCk7IHYgIT0gbmlsIHsKICAgICAgICByZXR1cm4gdi4oKnN0cnVjdE5hbWUpLmNvbXBhcmVGdW5jKGVsZW1lbnQpCiAgICB9CgogICAgcmV0dXJuIGZhbHNlCn0"
	standardStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBHZXQoa2V5IGZpZWxkVHlwZSkgaW50ZXJmYWNle30gewogICAgcmV0dXJuIHN0b3JlLnN0b3JlLkdldChrZXkpCn0KCmZ1bmMgKHN0b3JlICpTdG9yZSkgUmVtb3ZlKGtleSBmaWVsZFR5cGUpIHsKICAgIHN0b3JlLnN0b3JlLlJlbW92ZShrZXkpCn0KCmZ1bmMgKHN0b3JlICpTdG9yZSkgVmFsdWVzKCkgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKICAgIHN0b3JlLnN0b3JlLkNhbGxiYWNrRnVuYyhmdW5jKHYgaW50ZXJmYWNle30pIHsKICAgICAgICBpZiB2ICE9IG5pbCB7CiAgICAgICAgICAgIGVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCB2Ligqc3RydWN0TmFtZSkpCiAgICAgICAgfQogICAgfSkKICAgIHJldHVybgp9CgpmdW5jIChzdG9yZSAqU3RvcmUpIFRvSlNPTigpIChbXWJ5dGUsIGVycm9yKSB7CiAgICByZXR1cm4gc3RvcmUuc3RvcmUuVG9KU09OKCkKfQoKZnVuYyAoc3RvcmUgKlN0b3JlKSBCeUlkKGlkIGZpZWxkVHlwZSkgKnN0cnVjdE5hbWUgewogICAgaWYgdiA6PSBzdG9yZS5HZXQoaWQpOyB2ICE9IG5pbCB7CiAgICAgICAgcmV0dXJuIHYuKCpzdHJ1Y3ROYW1lKQogICAgfQoKICAgIHJldHVybiBuaWwKfQ"
	selectStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShmaWVsZE5hbWUgZmlsZWRUeXBlKSAqc3RydWN0TmFtZSB7CiAgICBpZiBpLCBmb3VuZCA6PSBzdG9yZS5tYXBGaWVsZDsgZm91bmQgewogICAgICAgIGlmIHYgOj0gc3RvcmUuR2V0KGkpOyB2ICE9IG5pbCB7CiAgICAgICAgICAgIHJldHVybiB2Ligqc3RydWN0TmFtZSkKICAgICAgICB9CiAgICB9CgogICAgcmV0dXJuIG5pbAp9"
	selectStoreCallbackFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSAoZWxlbWVudHMgW10qc3RydWN0TmFtZSkgewoJc3RvcmUuc3RvcmUuQ2FsbGJhY2tGdW5jKGZ1bmModiBpbnRlcmZhY2V7fSkgewoJCWlmIHYgIT0gbmlsIHsKCQkJaWYgdi4oKnN0cnVjdE5hbWUpLnN0cnVjdEZpZWxkID09IGZpZWxkTmFtZSB7CgkJCQllbGVtZW50cyA9IGFwcGVuZChlbGVtZW50cywgdi4oKnN0cnVjdE5hbWUpKQoJCQl9CgkJfQoJfSkKCXJldHVybgp9"
)

func (m *MetadataTable)ToStoreFormat(mapFunc, newFunc, selectFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, typeField, tableName string) (b string) {
	b = m.toStoreStructFormat(typeField)
	b +="\n\n"
	b += m.toNewStoreFuncFormat(newFunc, selectFunc, structPrefix, typeField, tableName)
	b += "\n\n"
	b += m.toMapStoreFuncFormat(mapFunc, typeField, structPrefix)
	b += "\n\n"
	b += m.toUpdateStoreFuncFormat(updateFunc, compareFunc, mapFunc, structPrefix)
	b += "\n\n"
	b += m.toCompareStoreFuncFormat(compareFunc, compareStructFunc, structPrefix)
	b += "\n\n"
	b += m.toStandardStoreFuncFormat(structPrefix, typeField)
	b += m.toSelectStoreFuncFormat(selectPrefix, structPrefix)
	return
}

func toStoreStructFormat(mappingField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(storeStructFormat)
	return strings.Replace(string(fieldFormat), "mappingField", mappingField, -1)
}

func (m *MetadataTable)toStoreStructFormat(typeField string) (b string) {
	fieldsLen := len(m.Fields)
	var elements []string
	mapCount := 0
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || i == 0 {
			continue
		}
		mapCount++
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]%v", toFieldUpperFormat(m.Fields[i].Name), "int", typeField))
		case "BIGINT":
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]%v", toFieldUpperFormat(m.Fields[i].Name), "int64", typeField))
		default:
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]%v", toFieldUpperFormat(m.Fields[i].Name), "string", typeField))
		}
	}

	if mapCount == 0 {
		return toStoreStructFormat("")
	}
	return toStoreStructFormat(strings.Join(elements, "\n"))
}

func toNewStoreFuncFormat(funcName, subFunc, makeField, mappingField, indexField, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(newStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "makeField", makeField, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	b = strings.Replace(b, "funcName", funcName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "indexField", indexField, -1)
	return strings.Replace(b, "mappingField", mappingField, -1)
}

func (m *MetadataTable)toNewStoreFuncFormat(funcName, subFunc, structPrefix, typeField, tableName string) (b string) {
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	mapCount := 0
	indexField := toFieldUpperFormat(m.Fields[0].Name)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || i == 0 {
			continue
		}
		mapCount++
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]%v),", toFieldUpperFormat(m.Fields[i].Name), "int", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[element.%v] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		case "BIGINT":
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]%v),", toFieldUpperFormat(m.Fields[i].Name), "int64", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[element.%v] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		default:
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]%v),", toFieldUpperFormat(m.Fields[i].Name), "string", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[database.ToLower(element.%v)] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		}
	}

	makeValues := ""
	mappingValues := ""
	if mapCount != 0 {
		makeValues = strings.Join(makeField, "\n")
		mappingValues = strings.Join(mappingField, "\n")
	}
	return toNewStoreFuncFormat(funcName, structPrefix + subFunc, makeValues, mappingValues, indexField, tableName)
}

func toMapStoreFuncFormat(funcName, structName, makeField, mappingField, storeField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(mapStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "makeField", makeField, -1)
	b = strings.Replace(b, "mappingField", mappingField, -1)
	return strings.Replace(b, "storeField", storeField, -1)
}

func (m *MetadataTable)toMapStoreFuncFormat(funcName, typeField, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	var storeField []string
	mapCount := 0
	indexField := toFieldUpperFormat(m.Fields[0].Name)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || i == 0 {
			continue
		}
		mapCount++
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]%v)", toFieldUpperFormat(m.Fields[i].Name), "int", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		case "BIGINT":
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]%v)", toFieldUpperFormat(m.Fields[i].Name), "int64", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		default:
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]%v)", toFieldUpperFormat(m.Fields[i].Name), "string", typeField))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[database.ToLower(element.%v)] = element.%v", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name), indexField))
		}
		storeField = append(storeField, fmt.Sprintf("\tstore.%vMapping = %vMapping", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
	}

	if mapCount == 0 {
		return ""
	}
	return toMapStoreFuncFormat(funcName, structName, strings.Join(makeField, "\n"), strings.Join(mappingField, "\n"), strings.Join(storeField, "\n"))
}

func toUpdateStoreFuncFormat(funcName, compareFunc, subFunc, indexField, mapField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "compareFunc", compareFunc, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "indexField", indexField, -1)
	return strings.Replace(b, "mapField", mapField, -1)
}

func (m *MetadataTable)toUpdateStoreFuncFormat(funcName, compareFunc, mapFunc, structPrefix string) (b string) {
	fieldsLen := len(m.Fields)
	mapCount := 0
	indexField := toFieldUpperFormat(m.Fields[0].Name)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || i == 0 {
			continue
		}
		mapCount++
	}
	mapField := ""
	if mapCount != 0 {
		mapField = fmt.Sprintf("\tif store.Updated {\n\t\tstore.%v()\n\t}", mapFunc)
	}
	return toUpdateStoreFuncFormat(funcName, compareFunc, structPrefix + funcName, indexField, mapField)
}

func toCompareStoreFuncFormat(funcName, compareFunc, indexField, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "compareFunc", compareFunc, -1)
	b = strings.Replace(b, "indexField", indexField, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable)toCompareStoreFuncFormat(funcName, compareFunc, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	indexField := toFieldUpperFormat(m.Fields[0].Name)
	return toCompareStoreFuncFormat(funcName, compareFunc, indexField, structName)
}

func toStandardStoreFuncFormat(fieldType, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(standardStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "structName", structName, -1)
	return strings.Replace(b, "fieldType", fieldType, -1)
}

func (m *MetadataTable)toStandardStoreFuncFormat(structPrefix, typeField string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	return toStandardStoreFuncFormat(typeField, structName)
}


func toSelectStoreFuncFormat(mapField, fieldName, filedType, funcName, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "filedType", filedType, -1)
	return strings.Replace(b, "mapField", mapField, -1)
}

func toSelectStoreCallbackFuncFormat(funcName, fieldName, fieldType, structField, structName string)(b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStoreCallbackFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "fieldType", fieldType, -1)
	return strings.Replace(b, "structField", structField, -1)
}

func (m *MetadataTable)toSelectStoreFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if m.Fields[i].Name == "id" || i == 0 {
			continue
		}
		if m.Fields[i].Name != "created_by" && !m.Fields[i].Unique {
			continue
		}
		b += "\n\n"
		fieldType := ""
		mapField := ""
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			fieldType = "int"
			mapField = fmt.Sprintf("%vMapping[%v]", m.Fields[i].ToUpperCase(), m.Fields[i].Name)
		case "BIGINT":
			fieldType = "int64"
			mapField = fmt.Sprintf("%vMapping[%v]", m.Fields[i].ToUpperCase(), m.Fields[i].Name)
		default:
			fieldType = "string"
			mapField = fmt.Sprintf("%vMapping[database.ToLower(%v)]", m.Fields[i].ToUpperCase(), m.Fields[i].Name)
		}
		funcName := funcPrefix + m.Fields[i].ToUpperCase()
		switch m.Fields[i].Name {
		case "created_by":
			b += toSelectStoreCallbackFuncFormat(funcName, m.Fields[i].Name, fieldType, m.Fields[i].ToUpperCase(), structName)
		default:
			b += toSelectStoreFuncFormat(mapField, m.Fields[i].Name, fieldType, funcName, structName)
		}
	}

	return
}