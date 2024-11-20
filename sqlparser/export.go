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

func ExportCrudFormatFile(modName, componentName, pkgName, commandFile, commonFile, storeFile, rootDir string, source *Database) {
	if componentName == "" {
		componentName = "components"
	}
	if pkgName == "" {
		pkgName = "databases"
	}
	if commandFile == "" {
		commandFile = "command.go"
	}
	if commonFile == "" {
		commonFile = "common.go"
	}
	if storeFile == "" {
		storeFile = "store.go"
	}
	pkgMod := fmt.Sprintf("%s/%s/%s", modName, componentName, pkgName)
	sqlMod := fmt.Sprintf("%s/%s/%s", modName, componentName, "mysql")
	commandImport := fmt.Sprintf("import (\n\t\"fmt\"\n\t\"%s\"\n\t\"%s\"\n)", pkgMod, sqlMod)
	commonImport := fmt.Sprintf("import (\n\t\"database/sql\"\n\t\"%s\"\n\t\"%s\"\n)", pkgMod, sqlMod)
	storeImport := fmt.Sprintf("import (\n\t\"%s/common\"\n\t\"%s\"\n\t\"sync\"\n)", modName, pkgMod)
	compare := fmt.Sprintf("package %v\n\n", pkgName)
	values := fmt.Sprintf("package %s\n\nimport(\n\t\"strconv\"\n\t\"fmt\"\n)", pkgName)
	parsetIntFormat := "func ParseInt(s string) int {\n\treturn int(ParseInt64(s))\n}"
	parsetInt64Format := "func ParseInt64(s string) int64 {\n\td, err := strconv.ParseInt(s, 10, 64)\n\tif err != nil {\n\t\treturn 0\n\t}\n\treturn d\n}"
	parsetFloatFormat := "func ParseFloat(s string) float64 {\n\td, err := strconv.ParseFloat(s, 64)\n\tif err != nil {\n\t\treturn 0\n\t}\n\treturn d\n}"
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
		structArray = append(structArray, source.Tables[i].ToStructFormat("json", "label"))
		tableName := fmt.Sprintf("%s.%sTable", pkgName, source.Tables[i].ToUpperCase())
		values += fmt.Sprintf("\t%sTable=\"%s\"\n", source.Tables[i].ToUpperCase(), source.Tables[i].Name)
		compareArray = append(compareArray, source.Tables[i].ToStructCompareFormat("s", "d", "Compare"))
		structName := fmt.Sprintf("%s.%s", pkgName, source.Tables[i].ToUpperCase())
		middleName := source.Tables[i].ToLowerCase()
		commandFileName := path.Join(rootDir, pkgName, middleName, commandFile)
		commonFileName := path.Join(rootDir, pkgName, middleName, commonFile)
		storeFileName := path.Join(rootDir, pkgName, middleName, storeFile)
		go func(pkgName, importHead, queryFunc, parserFunc, structPrefix, structName, databasePrefix, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToSelectSQLFormat("selectTable") + "\n\n"
			b += data.ToInsertSQLFormat("insert", structPrefix, structName) + "\n\n"
			b += data.ToUpdateSQLFormat("update") + "\n\n"
			b += data.ToRemoveSQLFormat("remove") + "\n\n"
			b += data.ToWhereSQLFormat("where") + "\n\n"
			b += data.ToQuerySQLFormat(queryFunc, "elements", structName) + "\n\n"
			b += data.ToParserSQLFormat(parserFunc, structPrefix, structName, databasePrefix) + "\n\n"
			b += data.ToSubSelectSQLFormat("by")
			b += data.ToSetSQLFormat("set")

			ch <- WriteFile(b, fileName)
		}(middleName, commandImport, "query", "parser", "element", structName, pkgName, commandFileName, source.Tables[i])
		go func(pkgName, importHead, InsertFunc, SelectFunc, UpdateFunc, UpdateTickerFunc, RemoveFunc, WhereFunc, ByFunc, SetFunc, queryFunc, structName, databasePrefix, tableName, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToInsertCrudFormat(InsertFunc, "insert", "element", structName, tableName) + "\n\n"
			b += data.ToSelectCrudFormat(SelectFunc, queryFunc, "selectTable", structName, tableName) + "\n\n"
			b += data.ToUpdateCrudFormat(UpdateFunc, "update", tableName) + "\n\n"
			if data.HasField("updated_at") {
				b += data.ToUpdateTickerCrudFormat(UpdateTickerFunc, queryFunc, fmt.Sprintf("%s.%s", databasePrefix, UpdateTickerFunc), structName, tableName) + "\n\n"
			}
			b += data.ToRemoveCrudFormat(RemoveFunc, "remove", tableName) + "\n\n"
			b += data.ToWhereCrudFormat(WhereFunc, queryFunc, "where", structName, tableName) + "\n\n"
			b += data.ToSubSelectCrudFormat(ByFunc, queryFunc, "by", structName, tableName)
			b += data.ToSetCrudFormat(SetFunc, "set", tableName) + "\n\n"

			ch <- WriteFile(b, fileName)
		}(middleName, commonImport, "Insert", "Select", "Update", "UpdateTicker", "Remove", "Where", "By", "Set", "query", structName, pkgName, tableName, commonFileName, source.Tables[i])
		go func(pkgName, importHead, newFunc, mapFunc, selectFunc, updateFunc, compareFunc, subSelectFunc, compareStruct, structPrefix, structName, tableName string, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToStoreFormat(newFunc, mapFunc, selectFunc, updateFunc, compareFunc, subSelectFunc, compareStruct, structPrefix, structName, tableName)

			ch <- WriteFile(b, fileName)
		}(middleName, storeImport, "NewStore", "Mapping", "Select", "UpdateTicker", "Compare", "By", structName, "store", "Store", tableName, storeFileName, source.Tables[i])
	}

	values += ")\n\n"

	values += parsetIntFormat + "\n\n"
	values += parsetInt64Format + "\n\n"
	values += parsetFloatFormat + "\n\n"
	values += selectFormat + "\n\n"
	values += whereFormat + "\n\n"
	if source.HasField("updated_at") {
		values += updateTickerFormat + "\n\n"
	}
	values += strings.Join(structArray, "\n\n")
	if err := WriteFile(values, path.Join(rootDir, pkgName, commonFile)); err != nil {
		fmt.Printf("[%s]ExportCrudFormatFile: %v\n", commonFile, err)
	}
	if err := WriteFile(compare+strings.Join(compareArray, "\n\n"), path.Join(rootDir, pkgName, "compare.go")); err != nil {
		fmt.Printf("[%s]ExportCrudFormatFile: %v\n", "compare.go", err)
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			if err != nil {
				fmt.Printf("[error]ExportCrudFormatFile: %v\n", err.Error())
			}
		}
	}
	fmt.Printf("[success]ExportCrudFormatFile: %d\n", count)
}

func ExportStorageFormatFile(modName, componentName, pkgName, dbPackageName, commonFile, rootDir string, source *Database) {
	store := "store"
	Store := "Store"
	importHead := fmt.Sprintf("import (\n\t\"%s/%s/%s\"\n)", modName, componentName, dbPackageName)
	values := fmt.Sprintf("package %s\n\nimport (\n\t\"encoding/json\"\n", pkgName)
	values += fmt.Sprintf("\t\"%s/common/memory\"\n\t\"%s/%s/%s\"\n", modName, modName, componentName, dbPackageName)

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
		importArray = append(importArray, fmt.Sprintf("\t\"%s/%s/%s/%s\"\n", modName, componentName, dbPackageName, source.Tables[i].ToLowerCase()))
		structArray = append(structArray, fmt.Sprintf("\t%s *%s.%s\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), Store))
		updateArray = append(updateArray, fmt.Sprintf("\t%s.%s.UpdateTicker(datetime)\n", store, source.Tables[i].ToUpperCase()))

		switch source.Tables[i].TypeOf() {
		case "int":
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "memory", "NewWithIntComparator"))
		case "int64":
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "memory", "NewWithInt64Comparator"))
		default:
			storeArray = append(storeArray, fmt.Sprintf("\t\t%s:%s.NewStore(%s.%s()),\n", source.Tables[i].ToUpperCase(), source.Tables[i].ToLowerCase(), "memory", "NewWithStringComparator"))
		}

		fName := path.Join(rootDir, pkgName, source.Tables[i].ToLowerCase()+".go")
		go func(pkgName, importHead, fromPrefix, selectPrefix, dbPackageName, storePrefix, StorePrefix, currentPrefix, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToSelectStorageFuncFormat(fromPrefix, selectPrefix, dbPackageName, storePrefix, StorePrefix, currentPrefix)

			ch <- WriteFile(b, fileName)
		}(pkgName, importHead, "From", "By", dbPackageName, "store", "Store", "current", fName, source.Tables[i])
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
			if err != nil {
				fmt.Printf("[error]ExportStorageFormatFile: %v\n", err)
			}
		}
	}
	fmt.Printf("[success]ExportStorageFormatFile: %d\n", count)
}

func ExportFrontendColumnsFormatFile(exportDir, rootDir string, source *Database) {
	columnsName := "columnsIndex"
	count := 0
	ch := make(chan error, 1)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 1
		fName := path.Join(rootDir, exportDir, source.Tables[i].ToLowerCase()+".js")
		go func(columnsName, fileName string, data *MetadataTable) {
			b := data.ToFrontendColumnsFormat(columnsName) + "\n\n"
			ch <- WriteFile(b, fileName)
		}(columnsName, fName, source.Tables[i])
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			if err != nil {
				fmt.Printf("[error]ExportFrontendColumnsFormatFile: %v\n", err)
			}
		}
	}
	fmt.Printf("[success]ExportFrontendColumnsFormatFile: %d\n", count)
}

func ExportForntendUnmarshalJSONFormatFile(modName, componentName, pkgName, rootDir string, source *Database) {
	count := 1
	ch := make(chan error, 1)

	go func(pkgName, fileName string){
		b := fmt.Sprintf("package %s\n\nimport(\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io\"\n\t\"strconv\"\n)\n\n", pkgName)
		b += "func toInt(s string) int {\n\treturn int(toInt64(s))\n}\n\n"
		b += "func toInt64(s string) int64 {\n\td, err := strconv.ParseInt(s, 10, 64)\n\tif err != nil {\n\t\treturn 0\n\t}\n\treturn d\n}\n\n"
		b += "func toFloat(s string) float64 {\n\td, err := strconv.ParseFloat(s, 64)\n\tif err != nil {\n\t\treturn 0\n\t}\n\treturn d\n}\n\n"
		b += "func UnmarshalAndTransform[T any](reader io.Reader, convert func(map[string]interface{}) T) (T, error) {\n\tresult := make(map[string]interface{})\n\tif err := json.NewDecoder(reader).Decode(&result); err != nil {\n\t\tvar zeroValue T\n\t\treturn zeroValue, fmt.Errorf(\"failed to decode JSON: %w\", err)\n\t}\n\n\treturn convert(result), nil\n}"

		ch <- WriteFile(b, fileName)
	}(pkgName, path.Join(rootDir, pkgName, "common.go"))

	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 1
		fName := path.Join(rootDir, pkgName, source.Tables[i].ToLowerCase()+".go")
		importHead := "import (\n\t\"fmt\"\n\t\"strings\"\n)"
		go func(pkgName, importHead, funcPrefix, elementName, safeName, tagName, labelName, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToForntendUnmarshalJSONFormat(funcPrefix, data.ToUpperCase(), elementName) + "\n\n"
			b += data.ToStructFormat(tagName, labelName) + "\n\n"
			b += data.ToStructSafeFormat(safeName, tagName, labelName)

			ch <- WriteFile(b, fileName)
		}(pkgName, importHead, "MapTo", "element", "Safe", "json", "label", fName, source.Tables[i])
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			if err != nil {
				fmt.Printf("[error]ExportForntendParseFormatFile: %v\n", err)
			}
		}
	}
	fmt.Printf("[success]ExportForntendParseFormatFile: %d\n", count)
}

func ExportModelFormatFile(modName, componentName, pkgName, dbPackageName, rootDir string, source *Database) {
	count := 0
	ch := make(chan error, 1)
	for i := range source.Tables {
		if source.Tables[i].Name == "" {
			continue
		}
		count += 1
		fName := path.Join(rootDir, pkgName, source.Tables[i].ToLowerCase()+".go")
		importHead := fmt.Sprintf("import (\n\t\"fmt\"\n\t\"database/sql\"\n\t\"%s/%s/%s/%s\"\n\t\"%s/%s/%s\"\n\t\"strings\"\n)", modName, componentName, dbPackageName, source.Tables[i].ToLowerCase(), modName, componentName, dbPackageName)
		go func(pkgName, importHead, createFunc, insertFunc, compareFunc, selectTableFunc, updateFunc, setFunc, removeFunc, whereFunc, fromPrefix, selectPrefix, dbPackageName, fileName string, data *MetadataTable) {
			b := fmt.Sprintf("package %s\n\n%s\n\n", pkgName, importHead)
			b += data.ToCreateModelFuncFormat(createFunc, insertFunc, dbPackageName) + "\n\n"
			b += data.ToCompareModelFuncFormat(compareFunc, "element", dbPackageName) + "\n\n"
			b += data.ToSelectTableModelFuncFormat(selectTableFunc, "Table", dbPackageName) + "\n\n"
			b += data.ToUpdateModelFuncFormat(updateFunc) + "\n\n"
			b += data.ToRemoveModelFuncFormat(removeFunc) + "\n\n"
			b += data.ToWhereModelFuncFormat(whereFunc, dbPackageName) + "\n\n"
			b += data.ToSelectModelFuncFormat(fromPrefix, selectPrefix, dbPackageName) + "\n"
			b += data.ToSetModelFuncFormat(updateFunc, setFunc)

			ch <- WriteFile(b, fileName)
		}(pkgName, importHead, "Create", "Insert", "Compare", "Select", "Update", "Set", "Remove", "Where", "From", "By", dbPackageName, fName, source.Tables[i])
	}

	for i := 0; i < count; i++ {
		select {
		case err := <-ch:
			if err != nil {
				fmt.Printf("[error]ExportModelFormatFile: %v\n", err)
			}
		}
	}
	fmt.Printf("[success]ExportModelFormatFile: %d\n", count)
}

func ExportFile(filename, tagField, labelField string, data []*MetadataTable) error {
	var element string
	for i := range data {
		element += "\n\n"
		element += data[i].ToStructFormat(tagField, labelField)
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
