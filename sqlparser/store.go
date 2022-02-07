package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	storeStructFormat = "dHlwZSBTdG9yZSBzdHJ1Y3QgewoJc3luYy5SV011dGV4Cgl0YWJsZSAgIHN0cmluZwoJc3RvcmUgICBjb21tb24uQnVja2V0Cm1hcHBpbmdGaWVsZAoJVXBkYXRlZCBib29sCglQYXRjaCAgIFtdaW50ZXJmYWNle30KfQ"
	newStoreFuncFormat = "ZnVuYyBOZXdTdG9yZShjIGNvbW1vbi5CdWNrZXQpIChzdG9yZSAqU3RvcmUpIHsKCXN0b3JlID0gJlN0b3JlewoJCXRhYmxlOiAgIHRhYmxlTmFtZSwKCQlzdG9yZTogICBjLAptYWtlRmllbGQKCQlVcGRhdGVkOiBmYWxzZSwKCX0KCWVsZW1lbnRzIDo9IHF1ZXJ5KGRhdGFiYXNlLlNlbGVjdFRhYmxlKHN0b3JlLnRhYmxlKSkKCWlmIGVsZW1lbnRzID09IG5pbCB8fCBsZW4oZWxlbWVudHMpIDw9IDAgewoJCXJldHVybgoJfQoJZm9yIGksIF8gOj0gcmFuZ2UgZWxlbWVudHMgewoJCWVsZW1lbnQgOj0gZWxlbWVudHNbaV0KCQlzdG9yZS5zdG9yZS5QdXQoZWxlbWVudC5JZCwgZWxlbWVudCkKbWFwcGluZ0ZpZWxkCgl9CgoJcmV0dXJuCn0"
	mapStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZSgpIHsKbWFrZUZpZWxkCglzdG9yZS5zdG9yZS5DYWxsYmFja0Z1bmMoZnVuYyh2IGludGVyZmFjZXt9KSB7CgkJaWYgdiAhPSBuaWwgewoJCQllbGVtZW50IDo9IHYuKCpzdHJ1Y3ROYW1lKQptYXBwaW5nRmllbGQKCQl9Cgl9KQpzdG9yZUZpZWxkCn0"
	updateStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShkYXRldGltZSBzdHJpbmcpIHsKCXN0b3JlLlVwZGF0ZWQgPSBmYWxzZQoJc3RvcmUuUGF0Y2ggPSBuaWwKCWVsZW1lbnRzIDo9IHF1ZXJ5KGRhdGFiYXNlLlNlbGVjdFVwZGF0ZWQoZGF0ZXRpbWUsIHN0b3JlLnRhYmxlKSkKCWlmIGVsZW1lbnRzID09IG5pbCB8fCBsZW4oZWxlbWVudHMpIDw9IDAgewoJCXJldHVybgoJfQoJZm9yIGkgOj0gMDsgaSA8IGxlbihlbGVtZW50cyk7IGkrKyB7CgkJZWxlbWVudCA6PSBlbGVtZW50c1tpXQoJCWlmICFzdG9yZS5jb21wYXJlRnVuYyhlbGVtZW50KSB7CgkJCXN0b3JlLnN0b3JlLlB1dChlbGVtZW50LklkLCBlbGVtZW50KQoJCQlzdG9yZS5QYXRjaCA9IGFwcGVuZChzdG9yZS5QYXRjaCwgZWxlbWVudCkKCQkJaWYgIXN0b3JlLlVwZGF0ZWQgewoJCQkJc3RvcmUuVXBkYXRlZCA9IHRydWUKCQkJfQoJCX0KCX0KCWlmIHN0b3JlLlVwZGF0ZWQgewoJCXN0b3JlLm1hcEZ1bmMoKQoJfQp9"
	compareStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lKSBib29sIHsKCWlmIHYgOj0gc3RvcmUuc3RvcmUuR2V0KGVsZW1lbnQuSWQpOyB2ICE9IG5pbCB7CgkJaWYgdi4oKnN0cnVjdE5hbWUpLmNvbXBhcmVGdW5jKGVsZW1lbnQpIHsKCQkJcmV0dXJuIHRydWUKCQl9Cgl9CgoJcmV0dXJuIGZhbHNlCn0"
	standardStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBHZXQoa2V5IGludCkgaW50ZXJmYWNle30gewoJcmV0dXJuIHN0b3JlLnN0b3JlLkdldChrZXkpCn0KCmZ1bmMgKHN0b3JlICpTdG9yZSkgUmVtb3ZlKGtleSBpbnQpIHsKCXN0b3JlLnN0b3JlLlJlbW92ZShrZXkpCn0KCmZ1bmMgKHN0b3JlICpTdG9yZSkgVmFsdWVzKCkgKGVsZW1lbnRzIFtdKnN0cnVjdE5hbWUpIHsKCXN0b3JlLnN0b3JlLkNhbGxiYWNrRnVuYyhmdW5jKHYgaW50ZXJmYWNle30pIHsKCQlpZiB2ICE9IG5pbCB7CgkJCWVsZW1lbnRzID0gYXBwZW5kKGVsZW1lbnRzLCB2Ligqc3RydWN0TmFtZSkpCgkJfQoJfSkKCXJldHVybgp9CgpmdW5jIChzdG9yZSAqU3RvcmUpIFRvSlNPTigpIChbXWJ5dGUsIGVycm9yKSB7CglyZXR1cm4gc3RvcmUuc3RvcmUuVG9KU09OKCkKfQoKZnVuYyAoc3RvcmUgKlN0b3JlKSBCeUlkKGlkIGludCkgKnN0cnVjdE5hbWUgewoJaWYgdiA6PSBzdG9yZS5HZXQoaWQpOyB2ICE9IG5pbCB7CgkJcmV0dXJuIHYuKCpzdHJ1Y3ROYW1lKQoJfQoKCXJldHVybiBuaWwKfQ"
	selectStoreFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShmaWVsZE5hbWUgZmlsZWRUeXBlKSAqc3RydWN0TmFtZSB7CglpZiBpLCBmb3VuZCA6PSBzdG9yZS5tYXBGaWVsZDsgZm91bmQgewoJCWlmIHYgOj0gci5HZXQoaSk7IHYgIT0gbmlsIHsKCQkJcmV0dXJuIHYuKCpzdHJ1Y3ROYW1lKQoJCX0KCX0KCglyZXR1cm4gbmlsCn0"
)

func (m *MetadataTable)ToStoreFormat(mapFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, tableName string) (b string) {
	b = m.toStoreStructFormat()
	b +="\n\n"
	b += m.toNewStoreFuncFormat(tableName)
	b += "\n\n"
	b += m.toMapStoreFuncFormat(mapFunc, structPrefix)
	b += "\n\n"
	b += m.toUpdateStoreFuncFormat(updateFunc, compareFunc, mapFunc)
	b += "\n\n"
	b += m.toCompareStoreFuncFormat(compareFunc, compareStructFunc, structPrefix)
	b += "\n\n"
	b += m.toStandardStoreFuncFormat(structPrefix)
	b += m.toSelectStoreFuncFormat(selectPrefix, structPrefix)
	return
}

func toStoreStructFormat(mappingField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(storeStructFormat)
	return strings.Replace(string(fieldFormat), "mappingField", mappingField, -1)
}

func (m *MetadataTable)toStoreStructFormat() (b string) {
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique {
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]int", toFieldUpperFormat(m.Fields[i].Name), "int"))
		case "BIGINT":
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]int", toFieldUpperFormat(m.Fields[i].Name), "int64"))
		default:
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]int", toFieldUpperFormat(m.Fields[i].Name), "string"))
		}
	}

	return toStoreStructFormat(strings.Join(elements, "\n"))
}

func toNewStoreFuncFormat(makeField, mappingField, tableName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(newStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "makeField", makeField, -1)
	b = strings.Replace(b, "tableName", tableName, -1)
	return strings.Replace(b, "mappingField", mappingField, -1)
}

func (m *MetadataTable)toNewStoreFuncFormat(tableName string) (b string) {
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique {
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]int),", toFieldUpperFormat(m.Fields[i].Name), "int"))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[element.%v] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		case "BIGINT":
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]int),", toFieldUpperFormat(m.Fields[i].Name), "int64"))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[element.%v] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		default:
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]int),", toFieldUpperFormat(m.Fields[i].Name), "string"))
			mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[strings.ToLower(element.%v)] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		}
	}

	return toNewStoreFuncFormat(strings.Join(makeField, "\n"), strings.Join(mappingField, "\n"), tableName)
}

func toMapStoreFuncFormat(funcName, structName, makeField, mappingField, storeField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(mapStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "makeField", makeField, -1)
	b = strings.Replace(b, "mappingField", mappingField, -1)
	return strings.Replace(b, "storeField", storeField, -1)
}

func (m *MetadataTable)toMapStoreFuncFormat(funcName, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	var storeField []string
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique {
			continue
		}
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]int)", toFieldUpperFormat(m.Fields[i].Name), "int"))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		case "BIGINT":
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]int)", toFieldUpperFormat(m.Fields[i].Name), "int64"))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		default:
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]int)", toFieldUpperFormat(m.Fields[i].Name), "string"))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[strings.ToLower(element.%v)] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		}
		storeField = append(storeField, fmt.Sprintf("\tstore.%vMapping = %vMapping", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
	}

	return toMapStoreFuncFormat(funcName, structName, strings.Join(makeField, "\n"), strings.Join(mappingField, "\n"), strings.Join(storeField, "\n"))
}

func toUpdateStoreFuncFormat(funcName, compareFunc, mapFunc string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "compareFunc", compareFunc, -1)
	return strings.Replace(b, "mapFunc", mapFunc, -1)
}

func (m *MetadataTable)toUpdateStoreFuncFormat(funcName, compareFunc, mapFunc string) (b string) {
	return toUpdateStoreFuncFormat(funcName, compareFunc, mapFunc)
}

func toCompareStoreFuncFormat(funcName, compareFunc, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "compareFunc", compareFunc, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable)toCompareStoreFuncFormat(funcName, compareFunc, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	return toCompareStoreFuncFormat(funcName, compareFunc, structName)
}

func toStandardStoreFuncFormat(structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(standardStoreFuncFormat)
	return strings.Replace(string(fieldFormat), "structName", structName, -1)
}

func (m *MetadataTable)toStandardStoreFuncFormat(structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	return toStandardStoreFuncFormat(structName)
}


func toSelectStoreFuncFormat(mapField, fieldName, filedType, funcName, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "filedType", filedType, -1)
	return strings.Replace(b, "mapField", mapField, -1)
}

func (m *MetadataTable)toSelectStoreFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique {
			continue
		}
		b += "\n\n"
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			b += toSelectStoreFuncFormat(fmt.Sprintf("%vMapping[%v]", m.Fields[i].ToUpperCase(), m.Fields[i].Name), m.Fields[i].Name, "int", funcPrefix + m.Fields[i].ToUpperCase(), structName)
		case "BIGINT":
			b += toSelectStoreFuncFormat(fmt.Sprintf("%vMapping[%v]", m.Fields[i].ToUpperCase(), m.Fields[i].Name), m.Fields[i].Name, "int64", funcPrefix + m.Fields[i].ToUpperCase(), structName)
		default:
			b += toSelectStoreFuncFormat(fmt.Sprintf("%vMapping[strings.ToLower(%v)]", m.Fields[i].ToUpperCase(), m.Fields[i].Name), m.Fields[i].Name, "string", funcPrefix + m.Fields[i].ToUpperCase(), structName)
		}
	}

	return
}