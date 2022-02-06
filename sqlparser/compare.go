package sqlparser

import (
	"encoding/base64"
	"strings"
)

const (
	compareFormat    = "ZnVuYyAobyAqc3RydWN0TmFtZSkgZnVuY05hbWUoZWxlbWVudCAqc3RydWN0TmFtZSkgYm9vbCB7CmNvbnRlbnRGaWVsZAoJcmV0dXJuIHRydWUKfQ"
	compareSubFormat = "CWlmIG8uZmllbGROYW1lICE9IGVsZW1lbnQuZmllbGROYW1lIHsKCQlyZXR1cm4gZmFsc2UKCX0"
)

func (m *MetadataTable) ToStructCompareFormat(funcName string) string {
	return m.toStructCompare(funcName)
}

func toStructCompare(contentField, structName, funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareFormat)
	b = strings.Replace(string(fieldFormat), "funcName", funcName, -1)
	b = strings.Replace(b, "structName", structName, -1)
	b = strings.Replace(b, "contentField", contentField, -1)
	return
}

func (m *MetadataTable) toStructCompare(funcName string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(compareSubFormat)
	structName := toFieldUpperFormat(m.Name)
	fieldsLen := len(m.Fields)
	var elements []string
	for i := 0; i < fieldsLen; i++ {
		element := strings.Replace(string(fieldFormat), "fieldName", toFieldUpperFormat(m.Fields[i].Name), -1)
		elements = append(elements, element)
	}

	contentField := strings.Join(elements, "\n")
	return toStructCompare(contentField, structName, funcName)
}
