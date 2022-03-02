package sqlparser

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

const (
	parsetIntImportFormat = "aW1wb3J0ICgKICAgICJzdHJjb252Igop"
	parsetIntFuncFormat   = "ZnVuYyBQYXJzZUludDY0KHMgc3RyaW5nKSBpbnQ2NCB7CglkLCBlcnIgOj0gc3RyY29udi5QYXJzZUludChzLCAxMCwgNjQpCglpZiBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gMAoJfQoJcmV0dXJuIGQKfQoKZnVuYyBQYXJzZUludChzIHN0cmluZykgaW50IHsKCXJldHVybiBpbnQoUGFyc2VJbnQ2NChzKSkKfQ"
	structFuncFormat      = "ZnVuYyBTZWxlY3QodGFibGUgc3RyaW5nKSBzdHJpbmcgewogICAgcmV0dXJuIGZtdC5TcHJpbnRmKGBTRUxFQ1QgKiBGUk9NICV2YCwgdGFibGUpCn0KCmZ1bmMgV2hlcmUoY29tbWFuZCBzdHJpbmcsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgU0VMRUNUICogRlJPTSAldiBXSEVSRSAldmAsIHRhYmxlLCBjb21tYW5kKQp9CgpmdW5jIFVwZGF0ZShjb21tYW5kIHN0cmluZywgaWQgaWRUeXBlLCB0YWJsZSBzdHJpbmcpIHN0cmluZyB7CiAgICByZXR1cm4gZm10LlNwcmludGYoYFVQREFURSAldiBTRVQgJXYgV0hFUkUgaWQ9JXZgLCB0YWJsZSwgY29tbWFuZCwgaWQpCn0KCmZ1bmMgUmVtb3ZlKGlkIGlkVHlwZSwgdXBkYXRlZF9ieSB1cGRhdGVkVHlwZSwgdXBkYXRlZF9hdCBzdHJpbmcsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgVVBEQVRFICV2IFNFVCBkZWxldGVkPTEsIHVwZGF0ZWRfYnk9JXYsIHVwZGF0ZWRfYXQ9IiV2IiBXSEVSRSBpZD0ldmAsIHRhYmxlLCB1cGRhdGVkX2J5LCB1cGRhdGVkX2F0LCBpZCkKfQoKZnVuYyBVcGRhdGVUaWNrZXIodXBkYXRlZF9hdCBzdHJpbmcsIHRhYmxlIHN0cmluZykgc3RyaW5nIHsKICAgIHJldHVybiBmbXQuU3ByaW50ZihgU0VMRUNUICogRlJPTSAldiBXSEVSRSB1cGRhdGVkX2F0ID4gIiV2ImAsIHRhYmxlLCB1cGRhdGVkX2F0KQp9"
)

func MkdirAll(p string) (err error) {
	if _, err = os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, os.ModePerm); err != nil {
			return
		}
	}
	return
}

func ExportFrontendColumnsFormatFile(head, foot, columnsName string, fileName string, data MetadataTable) error {
	element := head + "\n\n"
	element += data.ToFrontendColumnsFormat(columnsName) + "\n\n"
	element += foot
	return WriteFile(element, fileName)
}

func ExportForntendParseFormatFile(pkgName, importHead, funcPrefix, tagName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToForntendParseFormat(funcPrefix) + "\n\n"
	element += data.ToStructFormat(tagName)
	return WriteFile(element, fileName)
}

func ExportModelFormatFile(pkgName, importHead, createFunc, compreFunc, updateFunc, removeFunc, whereFunc, selectFunc, structPrefix, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToCreateModelFuncFormat(createFunc, structPrefix) + "\n\n"
	element += data.ToCompareModelFuncFormat(compreFunc, structPrefix) + "\n\n"
	element += data.ToUpdateModelFuncFormat(updateFunc, updateFunc) + "\n\n"
	element += data.ToRemoveModelFuncFormat(removeFunc, removeFunc) + "\n\n"
	element += data.ToWhereModelFuncFormat(whereFunc, structPrefix) + "\n\n"
	element += data.ToSelectModelFuncFormat(selectFunc, structPrefix)
	return WriteFile(element, fileName)
}

func ExportStorageSubFormatFile(pkgName, importHead, selectPrefix, structPrefix, fileName string, data MetadataTable) error {
	if data.PrimaryKeyLen() != 1 {
		return nil
	}
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToSelectStorageFuncFormat(selectPrefix, structPrefix)
	return WriteFile(element, fileName)
}

func ExportStorageFormatFile(pkgName, importHead, importPrefix, structName, fieldSuffix, newFunc, updateFunc, fileName string, data []MetadataTable) error {
	element := "package " + pkgName + "\n\n" + ToImportStorageFormat(importHead, importPrefix, data) + "\n\n"
	element += ToStructStorageFormat(structName, fieldSuffix, data) + "\n\n"
	element += ToInitialStorageFuncFormat() + "\n\n"
	element += ToNewStorageFuncFormat(newFunc, newFunc, structName, data) + "\n\n"
	element += ToUpdateStorageFuncFormat(updateFunc, structName, data) + "\n\n"
	return WriteFile(element, fileName)
}

func ExportStoreFormatFile(pkgName, importHead, mapFunc, newFunc, selectFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, tableName, fileName string, data MetadataTable) error {
	if data.PrimaryKeyLen() != 1 {
		return nil
	}
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	typeField := ""
	if len(data.Fields) > 0 {
		switch data.Fields[0].DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
			typeField = "int"
		case "BIGINT":
			typeField = "int64"
		default:
			typeField = "string"
		}
	}
	element += data.ToStoreFormat(mapFunc, newFunc, selectFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, typeField, tableName)

	return WriteFile(element, fileName)
}

func ExportPublicCrudFormatFile(pkgName, importHead, insertFunc, selectFunc, updateFunc, removeFunc, whereFunc, funcPrefix, subPrefix, structPrefix, tableName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToInsertCrudFormat(insertFunc, structPrefix, tableName) + "\n\n"
	element += data.ToSelectCrudFormat(selectFunc, structPrefix, tableName) + "\n\n"
	element += data.ToUpdateCrudFormat(updateFunc, "update", tableName) + "\n\n"
	element += data.ToRemoveCrudFormat(removeFunc, "remove", tableName) + "\n\n"
	element += data.ToWhereCrudFormat(whereFunc, structPrefix, tableName)
	element += data.ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName)

	return WriteFile(element, fileName)
}

func ExportCrudFormatFile(pkgName, importHead, structPrefix, insertName, queryName, parserName, selectName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToInsertFormat(structPrefix, insertName) + "\n\n"
	element += data.ToUpdateFormat("update") + "\n\n"
	element += data.ToRemoveFormat("remove") + "\n\n"
	element += data.ToQueryFormat(structPrefix, queryName) + "\n\n"
	element += data.ToParserFormat("element", structPrefix, parserName) + "\n\n"
	element += data.ToSelectFuncFormat("by")

	return WriteFile(element, fileName)
}

func ExportStructFormatFile(pkgName, tagName, fileName string, data []MetadataTable) error {
	tableSuffix := "Table"
	importField := []string{"strconv", "fmt"}
	values := fmt.Sprintf("package %v\n\n", pkgName)
	values += "import (\n"
	for i := range importField {
		values += fmt.Sprintf("\t%v%v%v\n", `"`, importField[i], `"`)
	}
	values += ")\n\n"
	values += "const (\n"
	for i := range data {
		values += fmt.Sprintf("\t%v%v=%v%v%v\n", data[i].ToUpperCase(), tableSuffix, `"`, data[i].Name, `"`)
	}
	values += ")\n\n"
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parsetIntFuncFormat)
	values += string(fieldFormat) +"\n"
	// select func
	values += "\nfunc Select(table string) string {\n"
	values += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %v`, table)\n", "%v")
	values += "}\n"

	// where func
	values += "\nfunc Where(command string, table string) string {\n"
	values += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %v WHERE %v`, table, command)\n", `%v`, "%v")
	values += "}\n"

	// update func
	values += "\nfunc UpdateTicker(updated_at string, table string) string {\n"
	values += fmt.Sprintf("\treturn fmt.Sprintf(`SELECT * FROM %v WHERE updated_at > %v`, table, updated_at)\n", `%v`, `"%v"`)
	values += "}\n"

	for i := range data {
		values += "\n"
		values += data[i].ToStructFormat(tagName)
	}

	return WriteFile(values, fileName)
}

func ExportStructCompareFormatFile(pkgName, src, dst, funcName, fileName string, data []MetadataTable) error {
	values := fmt.Sprintf("package %v\n\n", pkgName)
	for i := range data {
		if data[i].Name == "" {
			continue
		}
		values += data[i].ToStructCompareFormat(src, dst, funcName) +"\n\n"
	}

	return WriteFile(values, fileName)
}

func ExportFile(filename, tagField string, data []MetadataTable) error {
	var element string
	for i, _ := range data {
		element += "\n\n"
		element += data[i].ToStructFormat(tagField)
	}

	return WriteFile(element, filename)
}

func WriteFile(element string, filename string) error {
	if err := MkdirAll(filepath.Dir(filename)); err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(element)
	return err
}
