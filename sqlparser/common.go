package sqlparser

import (
	"strings"
)

var (
	mySQLKeywords = []string{
		"CREATE",
		"TABLE",
		"IF",
		"NOT",
		"EXISTS",
		"DEFAULT",
		"ASC",
		"PRIMARY",
		"KEY",
		"UNIQUE",
		"INDEX",
		"ENGINE",
		"InnoDB",
		"DEFAULT",
		"CHARACTER",
		"SET",
		"utf8mb4",
		"COLLATE",
		"utf8mb4_unicode_ci",
		"DESC",
		"DATABASE",
		"DROP",
		"USER",
		"USE",
		"GRANT",
		"SELECT",
		"NULL",
		"AUTO_INCREMENT",
	}

	dataTypeKeywords = []string{
		"INT",
		"TINYINT",
		"SMALLINT",
		"MEDIUMINT",
		"BIGINT",
		"FLOAT",
		"DOUBLE",
		"DECIMAL",
		"DATE",
		"DATETIME",
		"TIMESTAMP",
		"TIME",
		"YEAR",
		"CHAR",
		"VARCHAR",
		"TEXT",
		"TINYTEXT",
		"MEDIUMTEXT",
		"LONGTEXT",
		"BLOB",
		"BINARY",
		"VARBINARY",
		"ENUM",
	}
)

type MetadataTable struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name         string
	DataType     string
	Unique       bool
	PrimaryKey   bool
	AutoIncrment bool
}

func (f *MetadataTable) Keys() (elements []string) {
	for i := range f.Fields {
		elements = append(elements, f.Fields[i].Name)
	}
	return
}

func (f *MetadataTable) SetPrimaryKey(s string) {
	for i := range f.Fields {
		if strings.EqualFold(f.Fields[i].Name, s) {
			f.Fields[i].PrimaryKey = true
			break
		}
	}
}

func (f *MetadataTable) PrimaryKey() *Field {
	for i := range f.Fields {
		if f.Fields[i].PrimaryKey {
			return &f.Fields[i]
		}
	}
	return nil
}

func (f *MetadataTable) PrimaryKeyLen() (n int) {
	for i := range f.Fields {
		if f.Fields[i].PrimaryKey {
			n++
		}
	}
	return
}

func (f *Field) TypeOf() string {
	switch f.DataType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
		return "int"
	case "BIGINT":
		return "int64"
	default:
		return "string"
	}
}

func (f *Field) ValueOf() string {
	switch f.DataType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "FLOAT", "DOUBLE":
		return "%v"
	case "BIGINT":
		return "%v"
	default:
		return `"%v"`
	}
}

func (f *MetadataTable) ToUpperCase() string {
	return toFieldUpperFormat(f.Name)
}

func (f *MetadataTable) ToLowerCase() string {
	return toFieldLowerFormat(f.Name)
}

func (f *MetadataTable) ToLowerCamelCase() string {
	return toLowerCamelFormat(f.Name)
}

func (f *Field) ToUpperCase() string {
	return toFieldUpperFormat(f.Name)
}

func (f *Field) ToLowerCase() string {
	return toFieldLowerFormat(f.Name)
}

func (f *Field) ToLowerCamelCase() string {
	return toLowerCamelFormat(f.Name)
}

func findDataTypeString(s string) string {
	for i, v := range dataTypeKeywords {
		if strings.EqualFold(s, v) || strings.HasPrefix(s, v) {
			return dataTypeKeywords[i]
		}
	}
	return ""
}

func findKeywordString(s string) string {
	for i, v := range mySQLKeywords {
		if strings.EqualFold(s, v) {
			return mySQLKeywords[i]
		}
	}
	return ""
}
