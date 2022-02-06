package sqlparser

import (
	"os"
    "path/filepath"
)

func MkdirAll(p string) (err error) {
    if _, err = os.Stat(p); os.IsNotExist(err){
        if err = os.MkdirAll(p, os.ModePerm); err != nil {
            return
        }
    }
    return
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
    f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_TRUNC, os.ModePerm)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.WriteString(element)
    return err
}