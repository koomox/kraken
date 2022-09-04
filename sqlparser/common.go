package sqlparser

import (
	"fmt"
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

type Database struct {
	Tables []*MetadataTable
}

type MetadataTable struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name         string
	DataType     string
	Unique       bool
	PrimaryKey   bool
	AutoIncrment bool
}

func (source *Database) ToString() (s string) {
	for k := range source.Tables {
		s += fmt.Sprintf("Table: %s\n", source.Tables[k].Name)
		for i := range source.Tables[k].Fields {
			s += fmt.Sprintf("\tField: %s %s", source.Tables[k].Fields[i].Name, source.Tables[k].Fields[i].DataType)
			if source.Tables[k].Fields[i].Unique {
				s += " UNIQUE"
			}
			if source.Tables[k].Fields[i].PrimaryKey {
				s += " PRIMARY KEY"
			}
			if source.Tables[k].Fields[i].AutoIncrment {
				s += " AutoIncrment"
			}
			s += "\n"
		}
	}
	return
}

func (f *MetadataTable) PrimaryKey() (elements []*Field) {
	for i := range f.Fields {
		if f.Fields[i].PrimaryKey {
			elements = append(elements, f.Fields[i])
		}
	}
	return
}

func (f *MetadataTable) SetPrimaryKey(name string) {
	for i := range f.Fields {
		if strings.EqualFold(f.Fields[i].Name, name) {
			f.Fields[i].PrimaryKey = true
			break
		}
	}
}

func (f *MetadataTable) PrimaryKeyLen() (n int) {
	for i := range f.Fields {
		if f.Fields[i].PrimaryKey {
			n++
		}
	}
	return
}

func (f *MetadataTable) TypeOf() string {
	keys := f.PrimaryKey()
	switch len(f.PrimaryKey()) {
	case 1:
		return keys[0].TypeOf()
	default:
		return "string"
	}
}

func (f *MetadataTable) Id() string {
	var ids []string
	var idx []string
	keys := f.PrimaryKey()
	switch len(f.PrimaryKey()) {
	case 1:
		return fmt.Sprintf("element.%s", keys[0].ToUpperCase())
	default:
		for i := range keys {
			ids = append(ids, "%v")
			idx = append(idx, fmt.Sprintf("element.%s", keys[i].ToUpperCase()))
		}
		return "fmt.Sprintf(" + `"` + strings.Join(ids, "-") + `", ` + strings.Join(idx, ", ") + ")"
	}
}

func (f *MetadataTable) id() string {
	var ids []string
	var idx []string
	keys := f.PrimaryKey()
	switch len(f.PrimaryKey()) {
	case 1:
		return fmt.Sprintf("%s", keys[0].ToLowerCase())
	default:
		for i := range keys {
			ids = append(ids, "%v")
			idx = append(idx, fmt.Sprintf("%s", keys[i].ToLowerCase()))
		}
		return "fmt.Sprintf(" + `"` + strings.Join(ids, "-") + `", ` + strings.Join(idx, ", ") + ")"
	}
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
