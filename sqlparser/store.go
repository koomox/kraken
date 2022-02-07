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
)

func (m *MetadataTable)ToStoreFormat(mapFunc, structPrefix, tableName string) (b string) {
	b = m.toStoreStructFormat()
	b +="\n\n"
	b += m.toNewStoreFuncFormat(tableName)
	b += "\n\n"
	b += m.toMapStoreFuncFormat(mapFunc, structPrefix)
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
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]int", toFieldUpperFormat(m.Fields[i].Name), "int"))
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
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]int),", toFieldUpperFormat(m.Fields[i].Name), "int"))
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
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE":
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]int)", toFieldUpperFormat(m.Fields[i].Name), "int"))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		default:
			makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]int)", toFieldUpperFormat(m.Fields[i].Name), "string"))
			mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[strings.ToLower(element.%v)] = element.Id", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
		}
		storeField = append(storeField, fmt.Sprintf("\tstore.%vMapping = %vMapping", toFieldUpperFormat(m.Fields[i].Name), toFieldUpperFormat(m.Fields[i].Name)))
	}

	return toMapStoreFuncFormat(funcName, structName, strings.Join(makeField, "\n"), strings.Join(mappingField, "\n"), strings.Join(storeField, "\n"))
}