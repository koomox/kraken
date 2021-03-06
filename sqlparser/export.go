package sqlparser

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func MkdirAll(p string) (err error) {
	if _, err = os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, os.ModePerm); err != nil {
			return
		}
	}
	return
}

func ExportCrudFormatFile(modName, pkgName, commandFile, commonFile, storeFile, rootDir string, source *Database) {
	commandImport := fmt.Sprintf("import (\n\t\"fmt\"\n\t\"%s/component/database\"\n\t\"%s/component/mysql\"\n)", modName, modName)
	commonImport := fmt.Sprintf("import (\n\t\"database/sql\"\n\t\"%s/component/database\"\n\t\"%s/component/mysql\"\n)", modName, modName)
	storeImport := fmt.Sprintf("import (\n\t\"%s/common\"\n\t\"%s/component/database\"\n\t\"sync\"\n)", modName, modName)
	compare := fmt.Sprintf("package %v\n\n", pkgName)
	values := fmt.Sprintf("package %s\n\nimport(\n\t\"strconv\"\n\t\"fmt\"\n)", pkgName)
	parsetIntFormat := "func ParseInt(s string) int {\n\treturn int(ParseInt64(s))\n}"
	parsetInt64Format := "func ParseInt64(s string) int64 {\n\td, err := strconv.ParseInt(s, 10, 64)\n\tif err != nil {\n\t\treturn 0\n\t}\n\treturn d\n}"
	selectFormat := "func Select(table string) string {\n\treturn fmt.Sprintf(`SELECT * FROM %v`, table)\n}"
	whereFormat := "func Where(command string, table string) string {\n\treturn fmt.Sprintf(`SELECT * FROM %v WHERE %v`, table, command)\n}"
	updateTickerFormat := "func UpdateTicker(updated_at string, table string) string {\n\treturn fmt.Sprintf(`SELECT * FROM %v WHERE updated_at > \"%v\"`, table, updated_at)\n}"
	values += "\n\nconst (\n"
	var structArray []string
	var compareArray []string

	count := 0
	ch := make(chan error, 3)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 3
		structArray = append(structArray, source.Tables[i].ToStructFormat("json"))
		tableName := fmt.Sprintf("%s.%sTable", pkgName, source.Tables[i].ToUpperCase())
		values += fmt.Sprintf("\t%sTable=\"%s\"\n", source.Tables[i].ToUpperCase(), source.Tables[i].Name)
		compareArray = append(compareArray, source.Tables[i].ToStructCompareFormat("s", "d", "Compare"))
		structName := fmt.Sprintf("%s.%s", pkgName, source.Tables[i].ToUpperCase())
		middleName := source.Tables[i].ToLowerCase()
		commandFileName := path.Join(rootDir, pkgName, middleName, commandFile)
		commonFileName := path.Join(rootDir, pkgName, middleName, commonFile)
		storeFileName := path.Join(rootDir, pkgName, middleName, storeFile)
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

	values += ")\n\n"

	values += parsetIntFormat + "\n\n"
	values += parsetInt64Format + "\n\n"
	values += selectFormat + "\n\n"
	values += whereFormat + "\n\n"
	values += updateTickerFormat + "\n\n"
	values += strings.Join(structArray, "\n\n")
	if err := WriteFile(values, path.Join(rootDir, pkgName, commonFile)); err != nil {
		fmt.Printf("[%s]ExportCrudFormatFile: %v\n", commonFile, err)
	}
	if err := WriteFile(compare + strings.Join(compareArray, "\n\n"), path.Join(rootDir, pkgName, "compare.go")); err != nil {
		fmt.Printf("[%s]ExportCrudFormatFile: %v\n", "compare.go", err)
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			fmt.Printf("[%d]ExportCrudFormatFile: %v\n", i+1, err)
		}
	}
}

func ExportStorageFormatFile(modName, pkgName, component, database, commonFile, rootDir string, source *Database) {
	store := "store"
	Store := "Store"
	importHead := fmt.Sprintf("import (\n\t\"%s/%s/%s\"\n)", modName, component, database)
	values := fmt.Sprintf("package %s\n\nimport (\n\t\"encoding/json\"\n", pkgName)
	values += fmt.Sprintf("\t\"%s/common/cache\"\n\t\"%s/%s/%s\"\n", modName, modName, component, database)

	var importArray []string
	var structArray []string
	var storeArray []string
	var updateArray []string
	count := 0
	ch := make(chan error, 1)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 1
		importArray = append(importArray, fmt.Sprintf("\t\"%s/%s/%s/%s\"\n", modName, component, database, source.Tables[i].ToLowerCase()))
		structArray = append(structArray, fmt.Sprintf("\t%s *%s.%s\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), Store))
		updateArray = append(updateArray, fmt.Sprintf("\t%s.%s.UpdateTicker(datetime)\n", store, source.Tables[i].ToUpperCase()))

		switch source.Tables[i].TypeOf() {
		case "int":
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "cache", "NewWithIntComparator"))
		case "int64":
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "cache", "NewWithInt64Comparator"))
		default:
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "cache", "NewWithStringComparator"))
		}

		fName := path.Join(rootDir, pkgName, source.Tables[i].ToLowerCase()+".go")
		go func(pkgName, importHead, fromPrefix, selectPrefix, databasePrefix, storePrefix, StorePrefix, currentPrefix, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToSelectStorageFuncFormat(fromPrefix, selectPrefix, databasePrefix, storePrefix, StorePrefix, currentPrefix)

			ch <- WriteFile(b, fileName)
		}(pkgName, importHead, "From", "By", database, "store", "Store", "current", fName, source.Tables[i])
	}

	values += strings.Join(importArray, "")
	values += "\t\"sync\"\n\t\"time\"\n)\n\n"
	values += fmt.Sprintf("type %s struct {\n\tsync.RWMutex\n", Store)
	values += strings.Join(structArray, "")
	values += "\tFix []interface{}\n}\n\n"
	values += fmt.Sprintf("var (\n\tcurrent = &%s{}\n)\n\n", Store)
	values += "func Initial() {\n\tcurrent = NewStore()\n}\n\n"
	values += fmt.Sprintf("func Background() *%s {\n\treturn current\n}\n\n", Store)
	values += fmt.Sprintf("func NewStore() *%s {\n", Store)
	values += fmt.Sprintf("\treturn &%s{\n", Store)
	values += strings.Join(storeArray, "")
	values += "\t}\n}\n\n"
	values += fmt.Sprintf("func (%s *%s) UpdateTicker(datetime string) {\n", store, Store)
	values += strings.Join(updateArray, "")
	values += "}"

	filename := path.Join(rootDir, pkgName, commonFile)
	if err := WriteFile(values, filename); err != nil {
		fmt.Printf("[%s]ExportStorageFormatFile: %v\n", filename, err)
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			fmt.Printf("[%d]ExportStorageFormatFile: %v\n", i+1, err)
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

func ExportModelFormatFile(modName, pkgName, component, database, rootDir string, source *Database) {
	count := 0
	ch := make(chan error, 1)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 1
		fName := path.Join(rootDir, pkgName, source.Tables[i].ToLowerCase()+".go")
		importHead := fmt.Sprintf("import (\n\t\"fmt\"\n\t\"database/sql\"\n\t\"%s/%s/%s/%s\"\n\t\"%s/%s/%s\"\n\t\"strings\"\n)", modName, component, database, source.Tables[i].ToLowerCase(), modName, component, database)
		go func(pkgName, importHead, createFunc, insertFunc, compareFunc, updateFunc, removeFunc, whereFunc, fromPrefix, selectPrefix, databasePrefix, fileName string, data *MetadataTable){
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToCreateModelFuncFormat(createFunc, insertFunc, databasePrefix) + "\n\n"
			b += data.ToCompareModelFuncFormat(compareFunc, "element", databasePrefix) + "\n\n"
			b += data.ToUpdateModelFuncFormat(updateFunc) + "\n\n"
			b += data.ToRemoveModelFuncFormat(removeFunc) + "\n\n"
			b += data.ToWhereModelFuncFormat(whereFunc, databasePrefix) + "\n\n"
			b += data.ToSelectModelFuncFormat(fromPrefix, selectPrefix, databasePrefix)

			ch <- WriteFile(b, fileName)
		}(pkgName, importHead, "Create", "Insert", "Compare", "Update", "Remove", "Where", "From", "By", database, fName, source.Tables[i])
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			fmt.Printf("[%d]ExportModelFormatFile: %v\n", i+1, err)
		}
	}
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
