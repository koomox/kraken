package sqlparser

import (
	"fmt"
)

func (m *MetadataTable) ToStructCompareFormat(src, dst, funcName string) (b string) {
	b = fmt.Sprintf("func (%v *%v) %v(%v *%v) bool {\n", src, m.ToUpperCase(), funcName, dst, m.ToUpperCase())
	for i := range m.Fields {
		b += fmt.Sprintf("\tif %v.%v != %v.%v {\n\t\treturn false\n\t}\n", src, m.Fields[i].ToUpperCase(), dst, m.Fields[i].ToUpperCase())
	}
	b += "\treturn true\n}\n"
	return
}