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
		"COMMENT",
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

	queryKeywords = []string{
		"id",
		"userid",
		"username",
		"created_by",
	}

	requiredUpdateKeywords = []string{
		"updated_by",
		"updated_at",
	}

	requiredCreatedKeywords = []string{
		"created_by",
		"created_at",
	}
)

type Database struct {
	Tables []*MetadataTable
}

type MetadataTable struct {
	Name             string
	Fields           []*Field
	RequiredUpdate   bool
	RequiredCreated  bool
	RequiredDatetime bool
	UpdateTimeField  string
	HasIndex         bool
	IndexFields      [][]string
}

type Field struct {
	Name             string
	DataType         string
	Comment          string
	Unique           bool
	PrimaryKey       bool
	AutoIncrment     bool
	HasComment       bool
	HasQuery         bool
	HasIndex         bool
	RequiredUpdate   bool
	RequiredCreated  bool
	RequiredDatetime bool
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

func (source *Database) EnableQueryFields(words ...string) {
	fields := make(map[string]bool, len(words))
	for _, word := range words {
		fields[word] = true
	}

	for idx := range source.Tables {
		for i := range source.Tables[idx].Fields {
			if _, found := fields[source.Tables[idx].Fields[i].Name]; found {
				source.Tables[idx].Fields[i].HasQuery = true
			}
		}
	}
}

func (source *Database) EnableRequiredUpdateFields(words ...string) {
	fields := make(map[string]bool, len(words))
	for _, word := range words {
		fields[word] = true
	}

	for idx := range source.Tables {
		for i := range source.Tables[idx].Fields {
			if _, found := fields[source.Tables[idx].Fields[i].Name]; found {
				source.Tables[idx].RequiredUpdate = true
				source.Tables[idx].Fields[i].RequiredUpdate = true
				if source.Tables[idx].Fields[i].TypeOf() == "string" {
					source.Tables[idx].UpdateTimeField = source.Tables[idx].Fields[i].Name
				}
			}
		}
	}
}

func (source *Database) EnableRequiredCreatedFields(words ...string) {
	fields := make(map[string]bool, len(words))
	for _, word := range words {
		fields[word] = true
	}

	for idx := range source.Tables {
		for i := range source.Tables[idx].Fields {
			if _, found := fields[source.Tables[idx].Fields[i].Name]; found {
				source.Tables[idx].RequiredCreated = true
				source.Tables[idx].Fields[i].RequiredCreated = true
			}
		}
	}
}

func (source *Database) EnableRequiredDatetimeFields() {
	for idx := range source.Tables {
		for i := range source.Tables[idx].Fields {
			if source.Tables[idx].Fields[i].DataType == "DATETIME" || source.Tables[idx].Fields[i].DataType == "DATE" {
				source.Tables[idx].RequiredDatetime = true
				source.Tables[idx].Fields[i].RequiredDatetime = true
			}
		}
	}
}

func (source *Database) EnableIndexFields(words map[string][]string) {
	for i := range source.Tables {
		if fields, found := words[source.Tables[i].Name]; found {
			source.Tables[i].EnableIndexFields(fields...)
		}
	}
}

func (source *MetadataTable) EnableIndexFields(words ...string) {
	fields := make(map[string]bool, len(words))
	for _, word := range words {
		fields[word] = true
	}

	for i := range source.Fields {
		if _, found := fields[source.Fields[i].Name]; found {
			source.Fields[i].HasIndex = true
		}
	}
}

func (source *Database) EnableMultiIndexFields(words map[string][][]string) {
	for i := range source.Tables {
		if fields, found := words[source.Tables[i].Name]; found {
			source.Tables[i].HasIndex = true
			source.Tables[i].IndexFields = fields
		}
	}
}

func (source *Database) HasField(name string) bool {
	for i := range source.Tables {
		if source.Tables[i].HasField(name) {
			return true
		}
	}
	return false
}

func (f *MetadataTable) HasField(name string) bool {
	for i := range f.Fields {
		if strings.EqualFold(f.Fields[i].Name, name) {
			return true
		}
	}
	return false
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

func (m *MetadataTable) Id() string {
	var names, formats []string
	for _, field := range m.Fields {
		if field.PrimaryKey {
			names = append(names, fmt.Sprintf("element.%s", field.ToUpperCase()))
			formats = append(formats, "%v")
		}
	}
	return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", strings.Join(formats, "-"), strings.Join(names, ", "))
}

func (m *MetadataTable) id() string {
	var names, formats []string
	for _, field := range m.Fields {
		if field.PrimaryKey {
			names = append(names, fmt.Sprintf("%s", field.ToLowerCase()))
			formats = append(formats, "%v")
		}
	}
	return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", strings.Join(formats, "-"), strings.Join(names, ", "))
}

func (m *MetadataTable) extractFieldFormat(filter func(field *Field) bool) (names, types, formats []string) {
	for _, field := range m.Fields {
		if filter(field) {
			names = append(names, field.Name)
			types = append(types, fmt.Sprintf("%s %s", field.Name, field.TypeOf()))
			formats = append(formats, fmt.Sprintf(`%s=%v`, field.Name, field.ValueOf()))
		}
	}
	return names, types, formats
}

func (m *MetadataTable) ExtractUpdateFieldFormat() ([]string, []string, []string) {
	return m.extractFieldFormat(func(field *Field) bool {
		return field.RequiredUpdate
	})
}

func (m *MetadataTable) ExtractPrimaryFieldFormat() ([]string, []string, []string) {
	return m.extractFieldFormat(func(field *Field) bool {
		return field.PrimaryKey
	})
}

func (m *MetadataTable) ExtractPrimaryAndUpdateFieldFormat() ([]string, []string, []string) {
	return m.extractFieldFormat(func(field *Field) bool {
		return field.PrimaryKey || field.RequiredUpdate
	})
}

func (f *Field) TypeOf() string {
	switch f.DataType {
	case "TINYINT", "SMALLINT", "MEDIUMINT":
		return "int"
	case "INT", "BIGINT":
		return "int64"
	case "FLOAT", "DOUBLE", "DECIMAL":
		return "float64"
	default:
		return "string"
	}
}

func (f *Field) TypeSafeOf() string {
	switch f.DataType {
	case "TINYINT", "SMALLINT", "MEDIUMINT":
		return "int"
	case "INT", "BIGINT":
		return "string"
	case "FLOAT", "DOUBLE", "DECIMAL":
		return "float64"
	default:
		return "string"
	}
}

func (f *Field) ValueOf() string {
	switch f.DataType {
	case "TINYINT", "SMALLINT", "MEDIUMINT":
		return "%v"
	case "INT", "BIGINT":
		return "%v"
	case "FLOAT", "DOUBLE", "DECIMAL":
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
