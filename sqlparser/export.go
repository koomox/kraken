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

func ExportPublicCrudFormatFile(pkgName, importHead, insertFunc, selectFunc, compareFunc, updateFunc, funcPrefix, subPrefix, structPrefix, tableName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToPublicCrudFormat(insertFunc, selectFunc, compareFunc, updateFunc, structPrefix, tableName)
	element += data.ToPublicSubCrudFormat(funcPrefix, subPrefix, structPrefix, tableName)

	return WriteFile(element, fileName)
}

func ExportCrudFormatFile(pkgName, importHead, structPrefix, insertName, queryName, parserName, computedName, fileName string, data MetadataTable) error {
	element := "package " + pkgName + "\n\n" + importHead + "\n\n"
	element += data.ToInsertFormat(structPrefix, insertName) + "\n\n"
	element += data.ToQueryFormat(structPrefix, queryName) + "\n\n"
	element += data.ToParserFormat("element", structPrefix, parserName) + "\n\n"
	element += data.ToCompareCrudFormat(structPrefix, computedName)
	element += data.ToSelectFuncFormat("select")

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
		tables = append(tables, fmt.Sprintf(`%v="%v"`, toFieldUpperFormat(data[i].Name)+tableSuffix, data[i].Name))
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
