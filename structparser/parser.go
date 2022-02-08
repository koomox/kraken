package structparser

import (
	"reflect"
	"encoding/base64"
	"strings"
)

func toLabelFormat(labelField, fieldName, hiddenField, visibleField, writableField, updateWritableField, updateVisibleField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(labelFieldFormat)
	b = strings.Replace(string(fieldFormat), "labelField", labelField, -1)
	b = strings.Replace(b, "fieldName", fieldName, -1)
	b = strings.Replace(b, "hiddenField", hiddenField, -1)
	b = strings.Replace(b, "visibleField", visibleField, -1)
	b = strings.Replace(b, "writableField", writableField, -1)
	b = strings.Replace(b, "updateWritableField", updateWritableField, -1)
	return strings.Replace(b, "updateVisibleField", updateVisibleField, -1)
}

func toColumnsFormat(columnsField, contentField string) (b string) {
	fieldFormat, _ := base64.RawStdEncoding.DecodeString(columnsFieldFormat)
	b = strings.Replace(string(fieldFormat), "columnsField", columnsField, -1)
	return strings.Replace(b, "contentField", contentField, -1)
}

func ToFrontendColumnsFormat(v interface{}, columnsName, tagName string) (b string) {
	ref := reflect.ValueOf(v).Elem()
	if columnsName == "" {
		columnsName = reflect.TypeOf(v).Elem().Name()
	}
	var elements []string
	for i := 0; i < ref.NumField(); i++ {
		element := ref.Type().Field(i)
		switch toLowerCase(element.Name) {
		case "id", "uid", "username":
			elements = append(elements, toLabelFormat(element.Name, element.Tag.Get(tagName), "false", "true", "false", "false", "true"))
		case "status", "deleted", "created_by", "updated_by", "created_at", "updated_at":
			elements = append(elements, toLabelFormat(element.Name, element.Tag.Get(tagName), "true", "true", "false", "false", "true"))
		default:
			elements = append(elements, toLabelFormat(element.Name, element.Tag.Get(tagName), "false", "true", "true", "true", "true"))
		}
	}

	return toColumnsFormat(columnsName, strings.Join(elements, "\n"))
}