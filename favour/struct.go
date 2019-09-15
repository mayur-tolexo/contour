package favour

import (
	"errors"
	"reflect"
	"strings"
)

//StructTagValue will return struct model tag field and value from specific tag
func StructTagValue(input interface{}, tag string) (fields map[string]interface{}, err error) {

	var refObj reflect.Value
	tag = GetAVal(tag, JSONTag)
	fields = make(map[string]interface{})
	refObj, err = getStructRefObj(input)

	if err == nil {
		for i := 0; i < refObj.NumField(); i++ {
			refField := refObj.Field(i)
			refType := refObj.Type().Field(i)

			//if not exported
			if refType.Name[0] > 'Z' {
				continue
			}

			//checking embeded anonymous structures
			if refType.Anonymous && refField.Kind() == reflect.Struct {
				var embdFields map[string]interface{}
				if embdFields, err = StructTagValue(refField.Interface(), tag); err != nil {
					break
				}
				MergeMap(fields, embdFields)
			} else if err = setFiledVal(tag, refType, refField, fields); err != nil {
				break
			}
		}
	}
	return
}

//getRefObj will return struct reflect obj
func getStructRefObj(input interface{}) (refObj reflect.Value, err error) {
	refObj = reflect.ValueOf(input)
	if refObj.Kind() == reflect.Ptr {
		refObj = refObj.Elem()
	}
	if refObj.Kind() != reflect.Struct || !refObj.IsValid() {
		err = errors.New("No a struct")
	}
	return
}

//setFiledVal will set struct current field value
func setFiledVal(tag string, refType reflect.StructField, refField reflect.Value,
	fields map[string]interface{}) (err error) {

	if col, exists := refType.Tag.Lookup(tag); exists {
		isDef := IsDefaultVal(refField)
		if col == "-" || (strings.Contains(col, ",omitempty") && isDef) {
			return
		}
		if tag == GORMTag {
			if sqlCol := refType.Tag.Get(SQLTag); sqlCol == "-" {
				return
			}
			col = strings.TrimPrefix(col, "column:")
		}
		col = strings.Split(col, ",")[0]
		dVal := refType.Tag.Get(DefaultTag)

		//if default tag set to null
		if strings.ToLower(dVal) == Null && isDef {
			fields[col] = nil
		} else {
			fields[col], err = GetFieldVal(refField)
		}
	}
	return
}
