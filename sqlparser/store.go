package sqlparser

import (
	"fmt"
	"strings"
)

func (m *MetadataTable) ToStoreFormat(newFunc, mapFunc, selectFunc, updateFunc, compareFunc, subSelectFunc, compareStruct, structPrefix, structName, tableName string) (b string) {
	b = m.toStoreStructFormat()
	b += "\n\n"
	b += m.toNewStoreFuncFormat(newFunc, selectFunc, structPrefix, structName)
	b += "\n\n"
	isVaild := false
	for i := range m.Fields {
		if !m.Fields[i].AutoIncrment && !m.Fields[i].PrimaryKey && m.Fields[i].Unique {
			isVaild = true
			break
		}
	}
	if isVaild {
		b += m.toMapStoreFuncFormat(mapFunc, structPrefix, structName, compareStruct)
		b += "\n\n"
	}
	b += m.toUpdateStoreFuncFormat(updateFunc, updateFunc, compareFunc, mapFunc, structPrefix, structName)
	b += "\n\n"
	b += m.toCompareStoreFuncFormat(compareFunc, compareFunc, structPrefix, structName, "element", compareStruct)
	b += "\n\n"
	b += m.toStandardStoreFuncFormat(structPrefix, structName, "elements", compareStruct)
	b += "\n\n"
	b += m.toSelectStoreFuncFormat(subSelectFunc, structPrefix, structName, compareStruct)
	return
}

func (m *MetadataTable) toStoreStructFormat() (b string) {
	fieldsLen := len(m.Fields)
	b = "type Store struct {\n\tsync.RWMutex\n\tstore common.Memory\n\t"
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("%vMapping map[%v]%v\n\t", toFieldUpperFormat(m.Fields[i].Name), m.Fields[i].TypeOf(), m.TypeOf())
	}
	b += "Updated bool\n\tPatch []interface{}\n}"
	return b
}

func (m *MetadataTable) toNewStoreFuncFormat(funcName, selectFunc, structPrefix, structName string) (b string) {
	fieldsLen := len(m.Fields)
	b = fmt.Sprintf("func %s(c common.Memory) (%s *%s) {\n\t%s = &%s{\n", funcName, structPrefix, structName, structPrefix, structName)
	b += fmt.Sprintf("\t\t%s:   c,\n", structPrefix)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("\t\t%vMapping: make(map[%v]%v),\n", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), m.TypeOf())
	}
	b += "\t\tUpdated: false,\n\t}\n"
	b += fmt.Sprintf("\telements := %s()\n", selectFunc)
	b += "\tif elements == nil || len(elements) <= 0 {\n\t\treturn\n\t}\n"
	b += "\tfor i := range elements {\n\t\telement := elements[i]\n"
	b += fmt.Sprintf("\t\t%s.store.Put(%s, element)\n", structPrefix, m.Id())
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("\t\t%s.%sMapping[element.%s] = %s\n", structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase(), m.Id())
	}
	b += "\t}\n\treturn\n}"
	return
}

func (m *MetadataTable) toMapStoreFuncFormat(funcName, structPrefix, structName, compareStruct string) (b string) {
	fieldsLen := len(m.Fields)
	b = fmt.Sprintf("func (%s *%s) %s() {\n", structPrefix, structName, funcName)
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("\t%sMapping := make(map[%s]%s)\n", m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), m.TypeOf())
	}
	b += fmt.Sprintf("\t%s.store.CallbackFunc(func(v interface{}) {\n", structPrefix)
	b += "\t\tif v != nil {\n"
	b += fmt.Sprintf("\t\t\telement := v.(*%s)\n", compareStruct)
	b += fmt.Sprintf("\t\t\tidx := %s\n", m.Id())
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("\t\t\t%sMapping[element.%s] = idx\n", m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase())
	}
	b += "\t\t}\n\t})\n"
	for i := 0; i < fieldsLen; i++ {
		if !m.Fields[i].Unique || m.Fields[i].PrimaryKey || m.Fields[i].AutoIncrment {
			continue
		}
		b += fmt.Sprintf("\t%s.%vMapping = %vMapping\n", structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase())
	}
	b += "}"
	return
}

func (m *MetadataTable) toUpdateStoreFuncFormat(funcName, updateFunc, compareFunc, mapFunc, structPrefix, structName string) (b string) {
	b = fmt.Sprintf("func (%s *%s) %s(datetime string) {\n", structPrefix, structName, funcName)
	b += fmt.Sprintf("\t%s.Updated = false\n\t%s.Patch = nil\n", structPrefix, structPrefix)
	b += fmt.Sprintf("\telements := %s(datetime)\n", updateFunc)
	b += "\tif elements == nil || len(elements) <= 0 {\n\t\treturn\n\t}\n"
	b += "\tfor i := 0; i < len(elements); i++ {\n"
	b += "\t\telement := elements[i]\n"
	b += fmt.Sprintf("\t\tif !%s.%s(element) {\n", structPrefix, compareFunc)
	b += fmt.Sprintf("\t\t\t%s.store.Put(%s, element)\n", structPrefix, m.Id())
	b += fmt.Sprintf("\t\t\t%s.Patch = append(%s.Patch, element)\n", structPrefix, structPrefix)
	b += fmt.Sprintf("\t\t\tif !%s.Updated {\n\t\t\t\t%s.Updated = true\n\t\t\t}\n", structPrefix, structPrefix)
	b += "\t\t}\n\t}"
	isVaild := false
	for i := range m.Fields {
		if !m.Fields[i].AutoIncrment && !m.Fields[i].PrimaryKey && m.Fields[i].Unique {
			isVaild = true
			break
		}
	}
	if isVaild {
		b += fmt.Sprintf("\n\tif %s.Updated {\n\t\t%s.%s()\n\t}", structPrefix, structPrefix, mapFunc)
	}
	b += "\n}"

	return
}

func (m *MetadataTable) toCompareStoreFuncFormat(funcName, compareFunc, srcName, srcStruct, dstName, dstStruct string) (b string) {
	b = fmt.Sprintf("func (%s *%s) %s(%s *%s) bool {\n", srcName, srcStruct, funcName, dstName, dstStruct)
	b += fmt.Sprintf("\tif v := %s.store.Get(%s); v != nil {\n", srcName, m.Id())
	b += fmt.Sprintf("\t\treturn v.(*%s).%s(%s)\n\t}\n\treturn false\n}", dstStruct, compareFunc, dstName)
	return
}

func (m *MetadataTable) toStandardStoreFuncFormat(srcName, srcStruct, dstName, dstStruct string) (b string) {
	b = fmt.Sprintf("func (%s *%s) Get(key %v) *%s {\n\tif v := %s.store.Get(key); v != nil {\n\t\treturn v.(*%s)\n\t}\n\treturn nil\n}\n", srcName, srcStruct, m.TypeOf(), dstStruct, srcName, dstStruct)
	b += fmt.Sprintf("\nfunc (%s *%s) Remove(key %v) {\n\t%s.store.Remove(key)\n}\n", srcName, srcStruct, m.TypeOf(), srcName)
	b += fmt.Sprintf("\nfunc (%s *%s) Values() (%s []*%s) {\n\t%s.store.CallbackFunc(func(v interface{}) {\n\t\tif v != nil {\n\t\t\t%s = append(%s, v.(*%s))\n\t\t}\n\t})\n\treturn\n}\n", srcName, srcStruct, dstName, dstStruct, srcName, dstName, dstName, dstStruct)
	b += fmt.Sprintf("\nfunc (%s *%s) ToJSON() ([]byte, error) {\n\treturn %s.store.ToJSON()\n}", srcName, srcStruct, srcName)
	return
}

func (m *MetadataTable) toSelectStoreFuncFormat(funcPrefix, structPrefix, structName, compareStruct string) (b string) {
	var idx []string
	var ids []string
	keys := m.PrimaryKey()
	switch len(keys) {
	case 1:
		b = fmt.Sprintf("func (%s *%s) %s%s(%s %s) *%s {\n", structPrefix, structName, funcPrefix, keys[0].ToUpperCase(), keys[0].ToLowerCase(), keys[0].TypeOf(), compareStruct)
	default:
		for _, v := range keys {
			idx = append(idx, fmt.Sprintf("%s", v.ToUpperCase()))
			ids = append(ids, fmt.Sprintf("%s %s", v.ToLowerCase(), v.TypeOf()))
		}
		b = fmt.Sprintf("func (%s *%s) %s%s(%s) *%s {\n", structPrefix, structName, funcPrefix, strings.Join(idx, "And"), strings.Join(ids, ", "), compareStruct)
	}
	b += fmt.Sprintf("\treturn %s.Get(%s)\n}", structPrefix, m.id())

	for i := range m.Fields {
		if !m.Fields[i].PrimaryKey && m.Fields[i].Unique {
			b += fmt.Sprintf("\n\nfunc (%s *%s) %s%s(%s %s) *%s {\n", structPrefix, structName, funcPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].ToLowerCase(), m.Fields[i].TypeOf(), compareStruct)
			b += fmt.Sprintf("\tif i, found := %s.%sMapping[%s]; found {\n", structPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].ToLowerCase())
			b += fmt.Sprintf("\t\treturn %s.Get(i)\n", structPrefix)
			b += "\t}\n\treturn nil\n}"
		}
	}
	return
}
