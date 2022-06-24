package sqlparser

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
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

func ExportCrudFormatFile(pkgName, commandFile, commonFile, storeFile, baseDir string, source *Database) {
	commandImport := fmt.Sprintf("import (\n\t\"fmt\"\n\t\"%s/component/database\"\n\t\"%s/component/mysql\"\n)", pkgName, pkgName)
	commonImport := fmt.Sprintf("import (\n\t\"database/sql\"\n\t\"%s/component/database\"\n\t\"%s/component/mysql\"\n)", pkgName, pkgName)
	storeImport := fmt.Sprintf("import (\n\t\"%s/common\"\n\t\"%s/component/database\"\n\t\"sync\"\n)", pkgName, pkgName)

	count := 0
	ch := make(chan error, 3)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 3
		tableName := fmt.Sprintf("%s.%sTable", pkgName, source.Tables[i].ToUpperCase())
		structName := fmt.Sprintf("%s.%s", pkgName, source.Tables[i].ToUpperCase())
		middleName := source.Tables[i].ToLowerCase()
		commandFileName := path.Join(baseDir, middleName, commandFile)
		commonFileName := path.Join(baseDir, middleName, commonFile)
		storeFileName := path.Join(baseDir, middleName, storeFile)
		go func(pkgName, importHead, insertFunc, updateFunc, removeFunc, queryFunc, parserFunc, selectFunc, structPrefix, structName, databasePrefix, fileName string, data *MetadataTable){
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToInsertSQLFormat(insertFunc, structPrefix, structName) + "\n\n"
			b += data.ToUpdateSQLFormat(updateFunc) + "\n\n"
			b += data.ToRemoveSQLFormat(removeFunc) + "\n\n"
			b += data.ToQuerySQLFormat(queryFunc, "elements", structName) + "\n\n"
			b += data.ToParserSQLFormat(parserFunc, structPrefix, structName, databasePrefix) + "\n\n"
			b += data.ToSubSelectSQLFormat(selectFunc)

			ch <- WriteFile(b, fileName)
		}(middleName, commandImport, "insert", "update", "remove", "query", "parser", "by", "element", structName, pkgName, commandFileName, source.Tables[i])
		go func(pkgName, importHead, InsertFunc, insertFunc, SelectFunc, selectFunc, UpdateFunc, updateFunc, UpdateTickerFunc, RemoveFunc, removeFunc, WhereFunc, ByFunc, byFunc, queryFunc, structName, databasePrefix, tableName, fileName string, data *MetadataTable){
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToInsertCrudFormat(InsertFunc, insertFunc, "element", structName, tableName) + "\n\n"
			b += data.ToSelectCrudFormat(SelectFunc, queryFunc, fmt.Sprintf("%s.%s", databasePrefix, SelectFunc), structName, tableName) + "\n\n"
			b += data.ToUpdateCrudFormat(UpdateFunc, updateFunc, tableName) + "\n\n"
			b += data.ToUpdateTickerCrudFormat(UpdateTickerFunc, queryFunc, fmt.Sprintf("%s.%s", databasePrefix, UpdateTickerFunc), structName, tableName) + "\n\n"
			b += data.ToRemoveCrudFormat(RemoveFunc, removeFunc, tableName) + "\n\n"
			b += data.ToWhereCrudFormat(WhereFunc, queryFunc, fmt.Sprintf("%s.%s", databasePrefix, WhereFunc), structName, tableName) +"\n\n"
			b += data.ToSubSelectCrudFormat(ByFunc, queryFunc, byFunc, structName, tableName)

			ch <- WriteFile(b, fileName)
		}(middleName, commonImport, "Insert", "insert", "Select", "select", "Update", "update", "UpdateTicker", "Remove", "remove", "Where", "By", "by", "query", structName, pkgName, tableName, commonFileName, source.Tables[i])
		go func(pkgName, importHead, newFunc, mapFunc, selectFunc, updateFunc, compareFunc, subSelectFunc, compareStruct, structPrefix, structName, tableName string, fileName string, data *MetadataTable){
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToStoreFormat(newFunc, mapFunc, selectFunc, updateFunc, compareFunc, subSelectFunc, compareStruct, structPrefix, structName, tableName)

			ch <- WriteFile(b, fileName)
		}(middleName, storeImport, "NewStore", "Mapping", "Select", "UpdateTicker", "Compare", "By", structName, "store", "Store", tableName, storeFileName, source.Tables[i])
	}
	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			fmt.Printf("[%d]ExportCrudFormatFile: %v\n", i+1, err)
		}
	}
}

func ExportFrontendColumnsFormatFile(head, foot, columnsName string, fileName string, data *MetadataTable) error {
	element := head + "\n\n"
	element += data.ToFrontendColumnsFormat(columnsName) + "\n\n"
	element += foot
	return WriteFile(element, fileName)
}

func ExportForntendParseFormatFile(pkgName, importHead, funcPrefix, tagName, fileName string, data *MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToForntendParseFormat(funcPrefix) + "\n\n"
	element += data.ToStructFormat(tagName)
	return WriteFile(element, fileName)
}

func ExportModelFormatFile(pkgName, importHead, createFunc, compreFunc, updateFunc, removeFunc, whereFunc, selectFunc, structPrefix, fileName string, data *MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToCreateModelFuncFormat(createFunc, structPrefix) + "\n\n"
	element += data.ToCompareModelFuncFormat(compreFunc, structPrefix) + "\n\n"
	element += data.ToUpdateModelFuncFormat(updateFunc, updateFunc) + "\n\n"
	element += data.ToRemoveModelFuncFormat(removeFunc, removeFunc) + "\n\n"
	element += data.ToWhereModelFuncFormat(whereFunc, structPrefix) + "\n\n"
	element += data.ToSelectModelFuncFormat(selectFunc, structPrefix)
	return WriteFile(element, fileName)
}

func ExportStorageSubFormatFile(pkgName, importHead, selectPrefix, structPrefix, fileName string, data *MetadataTable) error {
	if data.PrimaryKeyLen() != 1 {
		return nil
	}
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToSelectStorageFuncFormat(selectPrefix, structPrefix)
	return WriteFile(element, fileName)
}

func ExportStorageFormatFile(pkgName, importHead, importPrefix, structName, fieldSuffix, newFunc, updateFunc, fileName string, data []*MetadataTable) error {
	element := "package " + pkgName + "\n\n" + ToImportStorageFormat(importHead, importPrefix, data) + "\n\n"
	element += ToStructStorageFormat(structName, fieldSuffix, data) + "\n\n"
	element += ToInitialStorageFuncFormat() + "\n\n"
	element += ToNewStorageFuncFormat(newFunc, newFunc, structName, data) + "\n\n"
	element += ToUpdateStorageFuncFormat(updateFunc, structName, data) + "\n\n"
	return WriteFile(element, fileName)
}

func ExportStructFormatFile(pkgName, tagName, fileName string, data []*MetadataTable) error {
	tableSuffix := "Table"
	importField := []string{"strconv", "fmt"}
	values := fmt.Sprintf("package %v\n\n", pkgName)
	values += toImportFormat(importField...) + "\n"
	values += "const (\n"
	for i := range data {
		values += fmt.Sprintf("\t%v%v=%v%v%v\n", data[i].ToUpperCase(), tableSuffix, `"`, data[i].Name, `"`)
	}
	values += ")\n\n"
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parsetIntFuncFormat)
	values += string(fieldFormat) + "\n"
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

func ExportStructCompareFormatFile(pkgName, src, dst, funcName, fileName string, data []*MetadataTable) error {
	values := fmt.Sprintf("package %v\n\n", pkgName)
	for i := range data {
		if data[i].Name == "" {
			continue
		}
		values += data[i].ToStructCompareFormat(src, dst, funcName) + "\n\n"
	}

	return WriteFile(values, fileName)
}

func ExportFile(filename, tagField string, data []*MetadataTable) error {
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
