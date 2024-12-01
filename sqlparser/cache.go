package sqlparser

import (
	"fmt"
	"strings"
)

func (m *MetadataTable) ToCacheStructFormat(cacheName, recordName, databasePrefix string) (b string) {
	b = fmt.Sprintf("type %s struct {\n\tsync.RWMutex\n\t", cacheName)
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += fmt.Sprintf("records map[%s]*%s.%s\n\t", m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
			b += fmt.Sprintf("data []*%s.%s\n\t", databasePrefix, m.ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex || m.Fields[i].Unique {
			b += fmt.Sprintf("%s map[%s]*%s.%s\n\t", m.Fields[i].Name, m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
			continue
		}
	}
	b += fmt.Sprintf("rawData []byte\n\tversion string\n\tchanges []*%s\n}", recordName)
	return b
}

func (m *MetadataTable) ToRecordStructFormat(recordName, databasePrefix string) (b string) {
	return fmt.Sprintf("type %s struct {\n\tAction string\t`json:\"action\"`\n\tRecord *%s.%s\t`json:\"record\"`\n}", recordName, databasePrefix, m.ToUpperCase())
}

func (m *MetadataTable) ToNewCacheFuncFormat(funcName, selectFunc, structName, databasePrefix string) (b string) {
	b = fmt.Sprintf("func %s() *%s {\n\tdata := &%s{\n", funcName, structName, structName)
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += fmt.Sprintf("\t\trecords: make(map[%s]*%s.%s),\n", m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex {
			b += fmt.Sprintf("\t\t%s: make(map[%s]*%s.%s),\n", m.Fields[i].Name, m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
		}
	}
	b += "\t\tversion: \"\",\n\t}\n"
	b += fmt.Sprintf("\telements := %s()\n", selectFunc)
	b += "\tif elements == nil || len(elements) <= 0 {\n\t\treturn data\n\t}\n"
	b += "\tfor _, el := range elements {\n"
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += fmt.Sprintf("\t\tdata.records[el.%s] = el\n", m.Fields[i].ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex || m.Fields[i].Unique {
			b += fmt.Sprintf("\t\tdata.%s[el.%s] = el\n", m.Fields[i].Name, m.Fields[i].ToUpperCase())
			continue
		}
	}
	b += "\t}\n\tvalues := data.Values()\n\tdata.data = values\n\treturn data\n}"
	return
}

func (m *MetadataTable) ToGetCacheFuncFormat(funcName, cacheName, databasePrefix string) string {
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			return fmt.Sprintf("func (cache *%s) %s(key %s) *%s.%s {\tcache.RLock()\n\tdefer cache.RUnlock()\n\n\n\tif el, found := cache.records[key]; found {\n\t\treturn el\n\t}\n\treturn nil\n}", cacheName, funcName, m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
		}
	}
	return ""
}

func (m *MetadataTable) ToRemoveCacheFuncFormat(funcName, cacheName, databasePrefix string) string {
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			return fmt.Sprintf("func (cache *%s) %s(key %s) {\tcache.Lock()\n\tdefer cache.Unlock()\n\n\n\tdelete(cache.records, key)\n}", cacheName, funcName, m.Fields[i].TypeOf())
		}
	}
	return ""
}

func (m *MetadataTable) ToValuesCacheFuncFormat(funcName, cacheName, databasePrefix string) (b string) {
	b = fmt.Sprintf("func (cache *%s) %s() []*%s.%s {\n", cacheName, funcName, databasePrefix, m.ToUpperCase())
	b += fmt.Sprintf("\telements := make([]*%s.%s, 0, len(cache.records))\n", databasePrefix, m.ToUpperCase())
	b += "\n\tcache.RLock()\n\tfor _, el := range cache.records {\n\t\telements = append(elements, el)\n\t}\n\tcache.RUnlock()\n\n"
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += fmt.Sprintf("\tsort.Slice(elements, func(i, j int) bool {\n\t\treturn elements[i].%s < elements[j].%s\n\t})\n\n", m.Fields[i].ToUpperCase(), m.Fields[i].ToUpperCase())
		}
	}
	b += "\treturn elements\n}"
	return b
}

func (m *MetadataTable) ToJSONCacheFuncFormat(funcName, valuesFunc, cacheName string) (b string) {
	return fmt.Sprintf("func (cache *%s) %s() ([]byte, error) {\n\telements := cache.%s()\n\n\tdata, err := json.Marshal(elements)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn data, nil\n}", cacheName, funcName, valuesFunc)
}

func (m *MetadataTable) ToSyncCacheFuncFormat(funcName, selectFunc, cacheName, recordName, databasePrefix string) (b string) {
	b = fmt.Sprintf("func (cache *%s) %s() {\n", cacheName, funcName)
	b += fmt.Sprintf("\tresult := %s()\n", selectFunc)
	b += "\tif result == nil || len(result) <= 0 {\n\t\treturn\n\t}\n"
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			key := fmt.Sprintf("item.%s", m.Fields[i].ToUpperCase())
			b += fmt.Sprintf("\telements := make(map[%s]*%s.%s)\n", m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase())
			b += fmt.Sprintf("\tfor _, item := range result {\n\t\telements[%s] = item\n\t}\n\thasChanges := false\n\tcache.changes = cache.changes[:0]\n", key)
			b += fmt.Sprintf("\tfor _, item := range elements {\n\t\tif el, found := cache.records[%s]; found {\n\t\t\tif cache.Compare(el) {\n\t\t\t\thasChanges = true\n\t\t\t\tcache.changes = append(cache.changes, &%s{Action: \"update\", Record: item})\n\t\t\t}\n\t\t} else {\n\t\t\thasChanges = true\n\t\t\tcache.changes = append(cache.changes, &%s{Action: \"add\", Record: item})\n\t\t}\n\t}\n", key, recordName, recordName)
			b += fmt.Sprintf("\tfor _, item := range cache.records {\n\t\tif _, found := elements[%s]; !found {\n\t\t\thasChanges = true\n\t\t\tcache.changes = append(cache.changes, &%s{Action: \"remove\", Record: item})\n\t\t}\n\t}\n", key, recordName)
		}
	}
	b += "\tif !hasChanges {\n\treturn\n\t}\n"
	b += "\tcache.Lock()\n\tfor _, change := range cache.changes {\n\t\tswitch change.Action {\n"
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += "\t\tcase \"add\":\n"
			b += fmt.Sprintf("\t\t\tcache.records[change.Record.%s] = change.Record\n", m.Fields[i].ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex || m.Fields[i].Unique {
			b += fmt.Sprintf("\t\t\tcache.%s[change.Record.%s] = change.Record\n", m.Fields[i].Name, m.Fields[i].ToUpperCase())
		}
	}
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += "\t\tcase \"update\":\n"
			b += fmt.Sprintf("\t\t\tcache.records[change.Record.%s] = change.Record\n", m.Fields[i].ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex || m.Fields[i].Unique {
			b += fmt.Sprintf("\t\t\tcache.%s[change.Record.%s] = change.Record\n", m.Fields[i].Name, m.Fields[i].ToUpperCase())
		}
	}
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += "\t\tcase \"remove\":\n"
			b += fmt.Sprintf("\t\t\tdelete(cache.records, change.Record.%s)\n", m.Fields[i].ToUpperCase())
			continue
		}
		if m.Fields[i].HasIndex || m.Fields[i].Unique {
			b += fmt.Sprintf("\t\t\tdelete(cache.%s, change.Record.%s)\n", m.Fields[i].Name, m.Fields[i].ToUpperCase())
		}
	}
	b += "\t\t}\n\t}\n\tcache.Unlock()\n"
	b += "\n\tvalues := cache.Values()\n\tcache.data = values\n}"

	return
}

func (m *MetadataTable) ToCompareCacheFuncFormat(funcName, compareFunc, cacheName, databasePrefix string) (b string) {
	b = fmt.Sprintf("func (cache *%s) %s(element *%s.%s) bool {\n", cacheName, funcName, databasePrefix, m.ToUpperCase())
	b += "\tcache.RLock()\n\tdefer cache.RUnlock()\n\n"
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			b += fmt.Sprintf("\tif el, found := cache.records[element.%s]; found {\n", m.Fields[i].ToUpperCase())
		}
	}

	b += fmt.Sprintf("\t\treturn el.%s(element)\n\t}\n\treturn false\n}", compareFunc)
	return
}

func (m *MetadataTable) ToSubSelectCacheFuncFormat(funcPrefix, cacheName, databasePrefix string) string {
	var result []string
	for i := range m.Fields {
		if m.Fields[i].PrimaryKey {
			result = append(result, fmt.Sprintf("func (cache *%s) %s%s(%s %s) *%s.%s {\n\tif el, found := cache.records[%s]; found {\n\t\treturn el\n\t}\n\treturn nil\n}", cacheName, funcPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].Name, m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase(), m.Fields[i].Name))
			continue
		}
		if m.Fields[i].HasIndex {
			result = append(result, fmt.Sprintf("func (cache *%s) %s%s(key %s) *%s.%s {\n\tif el, found := cache.%s[key]; found {\n\t\treturn el\n\t}\n\treturn nil\n}", cacheName, funcPrefix, m.Fields[i].ToUpperCase(), m.Fields[i].TypeOf(), databasePrefix, m.ToUpperCase(), m.Fields[i].Name))
		}
	}
	return strings.Join(result, "\n\n")
}

func (m *MetadataTable) ToDataCacheFuncFormat(funcName, cacheName, databasePrefix string) string {
	return fmt.Sprintf("func (cache *%s) %s() []*%s.%s {\n\treturn cache.data\n}", cacheName, funcName, databasePrefix, m.ToUpperCase())
}
