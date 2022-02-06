package sqlparser

import (
	"os"
	"path/filepath"
	"fmt"
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

func ExportStructFormatFile(pkgName, tagName, fileName string, data []MetadataTable) error {
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
	element = "package " + pkgName + "\n\nconst (\n\t" + strings.Join(tables, "\n\t") + "\n)" + element

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
