package sqlparser

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	parsetIntImportFormat = "aW1wb3J0ICgKICAgICJzdHJjb252Igop"
	parsetIntFuncFormat = "ZnVuYyBQYXJzZUludDY0KHMgc3RyaW5nKSBpbnQ2NCB7CglkLCBlcnIgOj0gc3RyY29udi5QYXJzZUludChzLCAxMCwgNjQpCglpZiBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gMAoJfQoJcmV0dXJuIGQKfQoKZnVuYyBQYXJzZUludChzIHN0cmluZykgaW50IHsKCXJldHVybiBpbnQoUGFyc2VJbnQ2NChzKSkKfQ")

func MkdirAll(p string) (err error) {
	if _, err = os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, os.ModePerm); err != nil {
			return
		}
	}
	return
}

func ExportModelFormatFile(pkgName, importHead, createFunc, compreFunc, updateFunc, removeFunc, selectFunc, structPrefix, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToCreateModelFuncFormat(createFunc, structPrefix) + "\n\n"
	element += data.ToCompareModelFuncFormat(compreFunc, structPrefix) + "\n\n"
	element += data.ToUpdateModelFuncFormat(updateFunc, updateFunc) + "\n\n"
	element += data.ToRemoveModelFuncFormat(removeFunc, removeFunc) + "\n\n"
	element += data.ToSelectModelFuncFormat(selectFunc, structPrefix) + "\n\n"
	return WriteFile(element, fileName)
}

func ExportStorageSubFormatFile(pkgName, importHead, selectPrefix, structPrefix, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToSelectStorageFuncFormat(selectPrefix, structPrefix)
	return WriteFile(element, fileName)
}

func ExportStorageFormatFile(pkgName, importHead, importPrefix, structName, fieldSuffix, newFunc, updateFunc, fileName string, data []MetadataTable) error {
	element := "package " + pkgName + "\n\n" + ToImportStorageFormat(importHead, importPrefix ,data) + "\n\n"
	element += ToStructStorageFormat(structName, fieldSuffix, data) + "\n\n"
	element += ToInitialStorageFuncFormat() + "\n\n"
	element += ToNewStorageFuncFormat(newFunc, newFunc, structName, data) + "\n\n"
	element += ToUpdateStorageFuncFormat(updateFunc, structName, data) + "\n\n"
	return WriteFile(element, fileName)
}

func ExportStoreFormatFile(pkgName, importHead, mapFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, tableName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToStoreFormat(mapFunc, updateFunc, compareFunc, compareStructFunc, selectPrefix, structPrefix, tableName)

	return WriteFile(element, fileName)
}

func ExportPublicCrudFormatFile(pkgName, importHead, insertFunc, selectFunc, updateFunc, funcPrefix, subPrefix, structPrefix, tableName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToPublicCrudFormat(insertFunc, selectFunc,  updateFunc, structPrefix, tableName)
	element += data.ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName)

	return WriteFile(element, fileName)
}

func ExportCrudFormatFile(pkgName, importHead, structPrefix, insertName, queryName, parserName, selectName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToInsertFormat(structPrefix, insertName) + "\n\n"
	element += data.ToQueryFormat(structPrefix, queryName) + "\n\n"
	element += data.ToParserFormat("element", structPrefix, parserName) + "\n\n"
	element += data.ToSelectFuncFormat(selectName)

	return WriteFile(element, fileName)
}

func ExportStructFormatFile(pkgName, tagName, fileName string, data []MetadataTable) error {
	importHead, _ := base64.RawStdEncoding.DecodeString(parsetIntImportFormat)
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(parsetIntFuncFormat)
	tableSuffix := "Table"
	var tables []string

	var element string
	for i, _ := range data {
		if data[i].Name == "" {
			continue
		}
		tables = append(tables, fmt.Sprintf(`%v="%v"`, data[i].ToUpperCase()+tableSuffix, data[i].Name))
		element += "\n\n"
		element += data[i].ToStructFormat(tagName)
	}
	element = "package " + pkgName + "\n\n" + string(importHead) + "\n\nconst (\n\t" + strings.Join(tables, "\n\t") + "\n)\n\n" + string(fieldFormat) + element

	return WriteFile(element, fileName)
}

func ExportStructCompareFormatFile(pkgName, funcName, fileName string, data []MetadataTable) error {
	var element string
	for i, _ := range data {
		if data[i].Name == "" {
			continue
		}
		element += "\n\n"
		element += data[i].ToStructCompareFormat(funcName)
	}

	element = "package " + pkgName + element

	return WriteFile(element, fileName)
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
