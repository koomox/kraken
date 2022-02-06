package sqlparser

import (
	"os"
)

func ExportFile(filename, tagField string, data []MetadataTable) error {
	f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_TRUNC, os.ModePerm)
    if err != nil {
    	return err
    }
    defer f.Close()

    var b string
    for i, _ := range data {
    	b += "\n\n"
    	b += data[i].ToStructFormat(tagField)
    }

    _, err = f.WriteString(b)
    return err
}

func WriteFile(b string, filename string) error {
    f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_TRUNC, os.ModePerm)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.WriteString(b)
    return err
}