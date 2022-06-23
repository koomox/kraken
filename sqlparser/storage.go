package sqlparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	importStorageFormat      = "aW1wb3J0ICgKCSJlbmNvZGluZy9qc29uIgppbXBvcnRGaWVsZAoJInN5bmMiCgkidGltZSIKKQ"
	structStorageFormat      = "dHlwZSBzdHJ1Y3ROYW1lIHN0cnVjdCB7CglzeW5jLlJXTXV0ZXgKc3RydWN0RmllbGQKCUZpeCBbXWludGVyZmFjZXt9Cn0"
	newStorageFuncFormat     = "ZnVuYyBmdW5jTmFtZSgpIChzdG9yZSAqc3RydWN0TmFtZSkgewoJc3RvcmUgPSAmc3RydWN0TmFtZXsKY29udGVudEZpZWxkCgl9CgoJcmV0dXJuCn0"
	updateStorageFuncFormat  = "ZnVuYyAoc3RvcmUgKnN0cnVjdE5hbWUpIGZ1bmNOYW1lKGRhdGV0aW1lIHN0cmluZykgewpjb250ZW50RmllbGQKfQ"
	initialStorageFuncFormat = "dmFyICgKCWN1cnJlbnQgPSAmU3RvcmV7fQopCgpmdW5jIEluaXRpYWwoKSB7CgljdXJyZW50ID0gTmV3U3RvcmUoKQp9"
	selectStorageFuncFormat  = "ZnVuYyAoc3RvcmUgKlN0b3JlKSBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSBzdHJ1Y3ROYW1lIHsKICAgIHJldHVybiBzdG9yZS5zdHJ1Y3RGaWVsZC5zdWJGdW5jKGZpZWxkTmFtZSkKfQoKZnVuYyBmdW5jTmFtZShmaWVsZE5hbWUgZmllbGRUeXBlKSBzdHJ1Y3ROYW1lIHsKICAgIHJldHVybiBjdXJyZW50LmZ1bmNOYW1lKGZpZWxkTmFtZSkKfQ"
)

func toImportStorageFormat(importField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(importStorageFormat)
	return strings.Replace(string(fieldFormat), "importField", importField, -1)
}

func ToImportStorageFormat(importHead, importPrefix string, data []*MetadataTable) (b string) {
	var elements []string
	for i, _ := range data {
		elements = append(elements, "\t"+fmt.Sprintf(`"%v/%v"`, importPrefix, data[i].ToLowerCase()))
	}
	return toImportStorageFormat(importHead + "\n" + strings.Join(elements, "\n"))
}

func toStructStorageFormat(structName, structField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(structStorageFormat)
	b = strings.Replace(string(fieldFormat), "structName", structName, -1)
	return strings.Replace(b, "structField", structField, -1)
}

func ToStructStorageFormat(structName, fieldSuffix string, data []*MetadataTable) (b string) {
	var elements []string
	for i, _ := range data {
		if data[i].PrimaryKeyLen() != 1 {
			continue
		}
		elements = append(elements, "\t"+fmt.Sprintf(`%v *%v%v`, data[i].ToUpperCase(), data[i].ToLowerCase(), fieldSuffix))
	}
	return toStructStorageFormat(structName, strings.Join(elements, "\n"))
}

func toNewStorageFuncFormat(funcName, structName, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(newStorageFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func ToNewStorageFuncFormat(funcName, newFunc, structName string, data []*MetadataTable) (b string) {
	var elements []string
	for i, _ := range data {
		if data[i].PrimaryKeyLen() != 1 {
			continue
		}
		elements = append(elements, fmt.Sprintf("\t\t%v:%v.%v(bucket.NewStore()),", data[i].ToUpperCase(), data[i].ToLowerCase(), newFunc))
	}
	return toNewStorageFuncFormat(funcName, structName, strings.Join(elements, "\n"))
}

func toUpdateStorageFuncFormat(funcName, structName, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(updateStorageFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func ToUpdateStorageFuncFormat(funcName, structName string, data []*MetadataTable) (b string) {
	var elements []string
	for i, _ := range data {
		if data[i].PrimaryKeyLen() != 1 {
			continue
		}
		elements = append(elements, fmt.Sprintf("\tstore.%v.%v(datetime)", data[i].ToUpperCase(), funcName))
	}
	return toUpdateStorageFuncFormat(funcName, structName, strings.Join(elements, "\n"))
}

func ToInitialStorageFuncFormat() (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(initialStorageFuncFormat)
	return string(fieldFormat)
}

func toSelectStorageFuncFormat(subFunc, structField, fieldName, fieldType, funcName, structName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(selectStorageFuncFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "subFunc", subFunc, -1)
	b = strings.Replace(b, "structField", structField, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	return strings.Replace(b, "fieldType", fieldType, -1)
}

func (m *MetadataTable) ToSelectStorageFuncFormat(selectPrefix, structPrefix string) (b string) {
	structName := "*" + structPrefix + m.ToUpperCase()
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].PrimaryKey && !m.Fields[i].Unique {
			continue
		}
		subFunc := selectPrefix + m.Fields[i].ToUpperCase()
		funcName := "From" + m.ToUpperCase() + selectPrefix + m.Fields[i].ToUpperCase()
		fieldType := ""
		switch m.Fields[i].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			fieldType = "int"
		case "BIGINT":
			fieldType = "int64"
		default:
			fieldType = "string"
		}
		elements = append(elements, toSelectStorageFuncFormat(subFunc, m.ToUpperCase(), m.Fields[i].Name, fieldType, funcName, structName))
	}

	return strings.Join(elements, "\n\n")
}
