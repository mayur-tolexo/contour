package favour

import (
	"errors"
	"reflect"
	"time"
)

//struct tags
const (
	SQLTag     = "sql"
	JSONTag    = "json"
	GORMTag    = "gorm"
	DefaultTag = "default"
	Null       = "null"
)

//MergeMap will merge two maps
func MergeMap(a, b map[string]interface{}) {
	for k, v := range b {
		a[k] = v
	}
}

//GetAVal will return first non empty string value
func GetAVal(val ...string) string {
	for _, v := range val {
		if v != "" {
			return v
		}
	}
	return ""
}

//GetIVal will return first non empty int value
func GetIVal(val ...int) int {
	for _, v := range val {
		if v != 0 {
			return v
		}
	}
	return 0
}

//GetFVal will return first non empty float64 value
func GetFVal(val ...float64) float64 {
	for _, v := range val {
		if v != 0 {
			return v
		}
	}
	return 0
}

//GetFieldVal Convert reflect value to its corrosponding data type
func GetFieldVal(val reflect.Value) (castValue interface{}, err error) {
	switch val.Kind() {
	case reflect.String:
		castValue = val.String()
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		castValue = val.Int()
	case reflect.Float32, reflect.Float64:
		castValue = val.Float()
	case reflect.Map, reflect.Slice, reflect.Struct, reflect.Interface:
		castValue = val.Interface()
	default:
		err = errors.New("GetFieldVal: Invalid Filed Kind")
	}
	return
}

//IsDefault Check whether value is default
func IsDefault(val reflect.Value) (isDefault bool) {
	if !val.IsValid() {
		isDefault = true
	} else {
		zero := reflect.Zero(val.Type())
		fKind := val.Kind()

		if tVal, ok := val.Interface().(time.Time); ok {
			isDefault = tVal.IsZero()
		} else if fKind == reflect.Map || fKind == reflect.Slice {
			if val.Len() == 0 {
				isDefault = true
			}
		} else if val.Interface() == zero.Interface() {
			isDefault = reflect.DeepEqual(val, zero)
		}
	}
	return
}
