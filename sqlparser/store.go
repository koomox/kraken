package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	storeStructFormat               = "dHlwZSBTdG9yZSBzdHJ1Y3QgewogICAgc3luYy5SV011dGV4CiAgICB0YWJsZSAgIHN0cmluZwogICAgc3RvcmUgICBjb21tb24uQnVja2V0CmNvbnRlbnRGaWVsZCAgICBVcGRhdGVkIGJvb2wKICAgIFBhdGNoICAgW11pbnRlcmZhY2V7fQp9"
	newStoreFuncFormat              = "ZnVuYyBmdW5jTmFtZShjIGNvbW1vbi5CdWNrZXQpIChzdG9yZSAqU3RvcmUpIHsKICAgIHN0b3JlID0gJlN0b3JlewogICAgICAgIHRhYmxlOiAgIHRhYmxlTmFtZSwKICAgICAgICBzdG9yZTogICBjLAptYWtlRmllbGQKICAgICAgICBVcGRhdGVkOiBmYWxzZSwKICAgIH0KICAgIGVsZW1lbnRzIDo9IHF1ZXJ5KHN1YkZ1bmMoc3RvcmUudGFibGUpKQogICAgaWYgZWxlbWVudHMgPT0gbmlsIHx8IGxlbihlbGVtZW50cykgPD0gMCB7CiAgICAgICAgcmV0dXJuCiAgICB9CiAgICBmb3IgaSA6PSByYW5nZSBlbGVtZW50cyB7CiAgICAgICAgZWxlbWVudCA6PSBlbGVtZW50c1tpXQogICAgICAgIHN0b3JlLnN0b3JlLlB1dChlbGVtZW50LmluZGV4RmllbGQsIGVsZW1lbnQpCm1hcHBpbmdGaWVsZAogICAgfQoKICAgIHJldHVybgp9"
	mapStoreFuncFormat              = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZSgpIHsKbWFrZUZpZWxkCglzdG9yZS5zdG9yZS5DYWxsYmFja0Z1bmMoZnVuYyh2IGludGVyZmFjZXt9KSB7CgkJaWYgdiAhPSBuaWwgewoJCQllbGVtZW50IDo9IHYuKCpzdHJ1Y3ROYW1lKQptYXBwaW5nRmllbGQKCQl9Cgl9KQpzdG9yZUZpZWxkCn0"
	updateStoreFuncFormat           = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShkYXRldGltZSBzdHJpbmcpIHsKICAgIHN0b3JlLlVwZGF0ZWQgPSBmYWxzZQogICAgc3RvcmUuUGF0Y2ggPSBuaWwKICAgIGVsZW1lbnRzIDo9IHF1ZXJ5KHN1YkZ1bmMoZGF0ZXRpbWUsIHN0b3JlLnRhYmxlKSkKICAgIGlmIGVsZW1lbnRzID09IG5pbCB8fCBsZW4oZWxlbWVudHMpIDw9IDAgewogICAgICAgIHJldHVybgogICAgfQogICAgZm9yIGkgOj0gMDsgaSA8IGxlbihlbGVtZW50cyk7IGkrKyB7CiAgICAgICAgZWxlbWVudCA6PSBlbGVtZW50c1tpXQogICAgICAgIGlmICFzdG9yZS5jb21wYXJlRnVuYyhlbGVtZW50KSB7CiAgICAgICAgICAgIHN0b3JlLnN0b3JlLlB1dChlbGVtZW50LmluZGV4RmllbGQsIGVsZW1lbnQpCiAgICAgICAgICAgIHN0b3JlLlBhdGNoID0gYXBwZW5kKHN0b3JlLlBhdGNoLCBlbGVtZW50KQogICAgICAgICAgICBpZiAhc3RvcmUuVXBkYXRlZCB7CiAgICAgICAgICAgICAgICBzdG9yZS5VcGRhdGVkID0gdHJ1ZQogICAgICAgICAgICB9CiAgICAgICAgfQogICAgfQptYXBGaWVsZAp9"
	compareStoreFuncFormat          = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShlbGVtZW50ICpzdHJ1Y3ROYW1lKSBib29sIHsKICAgIGlmIHYgOj0gc3RvcmUuc3RvcmUuR2V0KGVsZW1lbnQuaW5kZXhGaWVsZCk7IHYgIT0gbmlsIHsKICAgICAgICByZXR1cm4gdi4oKnN0cnVjdE5hbWUpLmNvbXBhcmVGdW5jKGVsZW1lbnQpCiAgICB9CgogICAgcmV0dXJuIGZhbHNlCn0"
	standardStoreFuncFormat         = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBHZXQoYXJnc0ZpZWxkKSBpbnRlcmZhY2V7fSB7CiAgICByZXR1cm4gc3RvcmUuc3RvcmUuR2V0KGtleXNGaWVsZCkKfQoKZnVuYyAoc3RvcmUgKlN0b3JlKSBSZW1vdmUoYXJnc0ZpZWxkKSB7CiAgICBzdG9yZS5zdG9yZS5SZW1vdmUoa2V5c0ZpZWxkKQp9CgpmdW5jIChzdG9yZSAqU3RvcmUpIFZhbHVlcygpIChlbGVtZW50cyBbXSpzdHJ1Y3ROYW1lKSB7CiAgICBzdG9yZS5zdG9yZS5DYWxsYmFja0Z1bmMoZnVuYyh2IGludGVyZmFjZXt9KSB7CiAgICAgICAgaWYgdiAhPSBuaWwgewogICAgICAgICAgICBlbGVtZW50cyA9IGFwcGVuZChlbGVtZW50cywgdi4oKnN0cnVjdE5hbWUpKQogICAgICAgIH0KICAgIH0pCiAgICByZXR1cm4KfQoKZnVuYyAoc3RvcmUgKlN0b3JlKSBUb0pTT04oKSAoW11ieXRlLCBlcnJvcikgewogICAgcmV0dXJuIHN0b3JlLnN0b3JlLlRvSlNPTigpCn0KCmZ1bmMgKHN0b3JlICpTdG9yZSkgZnVuY05hbWUoYXJnc0ZpZWxkKSAqc3RydWN0TmFtZSB7CiAgICBpZiB2IDo9IHN0b3JlLkdldChrZXlzRmllbGQpOyB2ICE9IG5pbCB7CiAgICAgICAgcmV0dXJuIHYuKCpzdHJ1Y3ROYW1lKQogICAgfQoKICAgIHJldHVybiBuaWwKfQ"
	selectStoreFuncFormat           = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShmaWVsZE5hbWUgZmlsZWRUeXBlKSAqc3RydWN0TmFtZSB7CiAgICBpZiBpLCBmb3VuZCA6PSBzdG9yZS5tYXBGaWVsZDsgZm91bmQgewogICAgICAgIGlmIHYgOj0gc3RvcmUuR2V0KGkpOyB2ICE9IG5pbCB7CiAgICAgICAgICAgIHJldHVybiB2Ligqc3RydWN0TmFtZSkKICAgICAgICB9CiAgICB9CgogICAgcmV0dXJuIG5pbAp9"
	selectStorePrimaryKeyFuncFormat = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShhcmdzRmllbGQpICpzdHJ1Y3ROYW1lIHsKICAgIHN0b3JlLnN0b3JlLkNhbGxiYWNrRnVuYyhmdW5jKHYgaW50ZXJmYWNle30pIHsKICAgICAgICBpZiB2ICE9IG5pbCB7CiAgICAgICAgICAgIGlmIGZvcm1hdEZpZWxkIHsKICAgICAgICAgICAgICAgIHJldHVybiB2Ligqc3RydWN0TmFtZSkKICAgICAgICAgICAgfQogICAgICAgIH0KICAgIH0pCiAgICByZXR1cm4gbmlsCn0"
	selectStoreCallbackFuncFormat   = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShhcmdzRmllbGQpIChlbGVtZW50cyBbXSpzdHJ1Y3ROYW1lKSB7CiAgICBzdG9yZS5zdG9yZS5DYWxsYmFja0Z1bmMoZnVuYyh2IGludGVyZmFjZXt9KSB7CiAgICAgICAgaWYgdiAhPSBuaWwgewogICAgICAgICAgICBpZiBmb3JtYXRGaWVsZCB7CiAgICAgICAgICAgICAgICBlbGVtZW50cyA9IGFwcGVuZChlbGVtZW50cywgdi4oKnN0cnVjdE5hbWUpKQogICAgICAgICAgICB9CiAgICAgICAgfQogICAgfSkKICAgIHJldHVybgp9"
)

func (m *MetadataTable) ToStoreFormat(mapFunc, newFunc, selectFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, typeField, tableName string) (b string) {
	b = m.toStoreStructFormat()
	b += "\n\n"
	b += m.toNewStoreFuncFormat(newFunc, selectFunc, structPrefix, tableName)
	b += "\n\n"
	b += m.toMapStoreFuncFormat(mapFunc, structPrefix)
	b += "\n\n"
	b += m.toUpdateStoreFuncFormat(updateFunc, compareFunc, mapFunc, structPrefix)
	b += "\n\n"
	b += m.toCompareStoreFuncFormat(compareFunc, compareStructFunc, structPrefix)
	b += "\n\n"
	b += m.toStandardStoreFuncFormat(selectPrefix, structPrefix)
	b += m.toSelectStoreFuncFormat(selectPrefix, structPrefix)
	return
}

func toStoreStructFormat(contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(storeStructFormat)
	return strings.Replace(string(fieldFormat), "contentField", contentField, -1)
}

func (m *MetadataTable) toStoreStructFormat() (b string) {
	typeField := ""
	if v := m.PrimaryKey(); v != nil {
		typeField = v.TypeOf()
	}
	fieldsLen := len(m.Fields)
	var elements []string
	counter := 0
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		counter++
		elements = append(elements, fmt.Sprintf("\t%vMapping map[%v]%v\n", toFieldUpperFormat(m.Fields[i].Name), m.Fields[i].TypeOf(), typeField))
	}

	if counter == 0 {
		return toStoreStructFormat("")
	}
	return toStoreStructFormat(strings.Join(elements, ""))
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

func (m *MetadataTable) toNewStoreFuncFormat(funcName, subFunc, structPrefix, tableName string) (b string) {
	keyType := ""
	keyName := ""
	if v := m.PrimaryKey(); v != nil {
		keyType = v.TypeOf()
		keyName = v.ToUpperCase()
	}
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	counter := 0
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || i == 0 {
			continue
		}
		counter++
		makeField = append(makeField, fmt.Sprintf("\t\t%vMapping: make(map[%v]%v),", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), keyType))
		mappingField = append(mappingField, fmt.Sprintf("\t\tstore.%vMapping[element.%v] = element.%v", m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase(), keyName))
	}

	makeValues := ""
	mappingValues := ""
	if counter != 0 {
		makeValues = strings.Join(makeField, "\n")
		mappingValues = strings.Join(mappingField, "\n")
	}
	return toNewStoreFuncFormat(funcName, structPrefix+subFunc, makeValues, mappingValues, keyName, tableName)
}

func toMapStoreFuncFormat(funcName, structName, makeField, mappingField, storeField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(mapStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "makeField", makeField, -1)
	b = strings.Replace(b, "mappingField", mappingField, -1)
	return strings.Replace(b, "storeField", storeField, -1)
}

func (m *MetadataTable) toMapStoreFuncFormat(funcName, structPrefix string) (b string) {
	keyType := ""
	keyName := ""
	if v := m.PrimaryKey(); v != nil {
		keyType = v.TypeOf()
		keyName = v.ToUpperCase()
	}
	structName := structPrefix + toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var makeField []string
	var mappingField []string
	var storeField []string
	counter := 0
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		counter++
		makeField = append(makeField, fmt.Sprintf("\t%vMapping := make(map[%v]%v)", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), keyType))
		mappingField = append(mappingField, fmt.Sprintf("\t\t\t%vMapping[element.%v] = element.%v", m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase(), keyName))
		storeField = append(storeField, fmt.Sprintf("\tstore.%vMapping = %vMapping", m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase()))
	}

	if counter == 0 {
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

func (m *MetadataTable) toUpdateStoreFuncFormat(funcName, compareFunc, mapFunc, structPrefix string) (b string) {
	fieldsLen := len(m.Fields)
	mapCount := 0
	keyName := ""
	if v := m.PrimaryKey(); v != nil {
		keyName = v.ToUpperCase()
	}
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		mapCount++
	}
	mapField := ""
	if mapCount != 0 {
		mapField = fmt.Sprintf("\tif store.Updated {\n\t\tstore.%v()\n\t}", mapFunc)
	}
	return toUpdateStoreFuncFormat(funcName, compareFunc, structPrefix+funcName, keyName, mapField)
}

func toCompareStoreFuncFormat(funcName, compareFunc, indexField, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "compareFunc", compareFunc, -1)
	b = strings.Replace(b, "indexField", indexField, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable) toCompareStoreFuncFormat(funcName, compareFunc, structPrefix string) (b string) {
	structName := structPrefix + toFieldUpperFormat(m.Name)
	indexField := toFieldUpperFormat(m.Fields[0].Name)
	return toCompareStoreFuncFormat(funcName, compareFunc, indexField, structName)
}

func toStandardStoreFuncFormat(funcName, argsField, keysField, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(standardStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "keysField", keysField, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable) toStandardStoreFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + m.ToUpperCase()
	key := m.PrimaryKey()
	funcName := funcPrefix + key.ToUpperCase()
	if key == nil {
		return
	}
	return toStandardStoreFuncFormat(funcName, fmt.Sprintf("%v %v", key.ToLowerCamelCase(), key.TypeOf()), key.ToLowerCamelCase(), structName)
}

func toSelectStoreFuncFormat(mapField, fieldName, filedType, funcName, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStoreFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "filedType", filedType, -1)
	return strings.Replace(b, "mapField", mapField, -1)
}

func toSelectStoreCallbackFuncFormat(funcName, argsField, formatField, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStoreCallbackFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "formatField", formatField, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func toSelectStorePrimaryKeyFuncFormat(funcName, argsField, formatField, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStorePrimaryKeyFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "argsField", argsField, -1)
	b = strings.Replace(b, "formatField", formatField, -1)
	return strings.Replace(b, "structName", structName, -1)
}

func (m *MetadataTable) toSelectStoreFuncFormat(funcPrefix, structPrefix string) (b string) {
	structName := structPrefix + m.ToUpperCase()
	for i := range m.Fields {
		if !m.Fields[i].PrimaryKey && m.Fields[i].Unique {
			b += "\n\n"
			mapField := fmt.Sprintf("%vMapping[%v]", m.Fields[i].ToUpperCase(), m.Fields[i].Name)
			funcName := funcPrefix + m.Fields[i].ToUpperCase()
			b += toSelectStoreFuncFormat(mapField, m.Fields[i].Name, m.Fields[i].TypeOf(), funcName, structName)
		}
	}
	return
}
