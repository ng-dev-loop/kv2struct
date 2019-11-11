package kv2struct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	mTagKey           = "json"
	mHaveDefaultValue = false
)

// 需要实现的接口
type InterfaceKeyValue interface {
	GetString(key string, def ...string) string
}

//
func SetTagKey(key string) {
	mTagKey = key
}

//Unmarshal url.Values to struct
func Unmarshal(dataSource InterfaceKeyValue, object interface{}) error {
	val := reflect.ValueOf(object)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("Unmarshal() expects struct input. ")
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("Unmarshal() expects struct input. ")
	}
	return reflectValueFromTag(dataSource, val)
}

func reflectValueFromTag(dataSource InterfaceKeyValue, value reflect.Value) error {
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		err := reflectFieldFromTag(dataSource, field, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func reflectFieldFromTag(dataSource InterfaceKeyValue, field reflect.StructField, fieldValue reflect.Value) error {

	tag := field.Tag.Get(mTagKey)
	if tag == "-" {
		return nil
	}

	fieldStringData := ""
	kind := fieldValue.Kind()

	// fmt.Printf("field %v type %v\n", field.Name, kind.String())
	if kind != reflect.Struct {
		fieldStringData = getTagValue(dataSource, tag, field.Name)
	}

	switch kind {
	case reflect.String:
		fieldValue.SetString(fieldStringData)
	case reflect.Bool:
		if len(fieldStringData) == 0 {
			fieldStringData = "false"
		}
		b, err := strconv.ParseBool(fieldStringData)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"cast bool has error, expect type: %v ,value: %v ,query key: %v, field: %v",
				fieldValue.Type(), fieldStringData, tag, field.Name))
		}
		fieldValue.SetBool(b)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if len(fieldStringData) == 0 {
			fieldStringData = "0"
		}
		n, err := strconv.ParseUint(fieldStringData, 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"cast uint has error, expect type: %v ,value: %v ,query key: %v, field: %v, error: %v",
				fieldValue.Type(), fieldStringData, tag, field.Name, err))
		}
		if fieldValue.OverflowUint(n) {
			return errors.New(fmt.Sprintf(
				"cast uint has error, expect type: %v ,value: %v ,query key: %v, field: %v, overflow",
				fieldValue.Type(), fieldStringData, tag, field.Name))
		}

		fieldValue.SetUint(n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(fieldStringData) == 0 {
			fieldStringData = "0"
		}
		n, err := strconv.ParseInt(fieldStringData, 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"cast int has error, expect type: %v ,value: %v ,query key: %v, field: %v, error: %v",
				fieldValue.Type(), fieldStringData, tag, field.Name, err))
		}
		if fieldValue.OverflowInt(n) {
			return errors.New(fmt.Sprintf(
				"cast int has error, expect type: %v ,value: %v ,query key: %v, field: %v, overflow",
				fieldValue.Type(), fieldStringData, tag, field.Name))
		}
		fieldValue.SetInt(n)
	case reflect.Float32, reflect.Float64:
		if len(fieldStringData) == 0 {
			fieldStringData = "0"
		}
		n, err := strconv.ParseFloat(fieldStringData, fieldValue.Type().Bits())
		if err != nil {
			return errors.New(fmt.Sprintf(
				"cast float has error, expect type: %v ,value: %v ,query key: %v, field: %v, error: %v",
				fieldValue.Type(), fieldStringData, tag, field.Name, err))
		}
		if fieldValue.OverflowFloat(n) {
			return errors.New(fmt.Sprintf(
				"cast float has error, expect type: %v ,value: %v ,query key: %v, field: %v, overflow",
				fieldValue.Type(), fieldStringData, field.Name, tag))
		}
		fieldValue.SetFloat(n)
	case reflect.Struct:
		err := reflectValueFromTag(dataSource, fieldValue)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"cast struct has error, expect type: %v ,value: %v ,query key: %v, field: %v, overflow",
				fieldValue.Type(), fieldStringData, field.Name, tag))
		}
	case reflect.Ptr:
		fmt.Printf(fmt.Sprintf(
			"unsupported ptr type: %v ,value: %v ,query key: %v, field: %v\n",
			fieldValue.Type(), fieldStringData, tag, field.Name))
	case reflect.Array:
		fmt.Printf("unsupported ptr type: %v ,value: %v ,query key: %v, field: %v\n",
			fieldValue.Type(), fieldStringData, tag, field.Name)
	default:
		fmt.Printf("unsupported type: %v ,value: %v ,query key: %v, field: %v\n",
			fieldValue.Type(), fieldStringData, tag, field.Name)
	}

	return nil
}

//get val, if absent get from tag default val
func getTagValue(dataSource InterfaceKeyValue, tag string, defaultTag string) string {
	name, tagOptions := parseTag(tag)
	if len(name) == 0 {
		name = defaultTag
	}
	uv := dataSource.GetString(name)
	if mHaveDefaultValue {
		optsLen := len(tagOptions)
		if optsLen > 0 {
			if optsLen == 1 && uv == "" {
				uv = tagOptions[0]
			}
		}
	}
	return uv
}

func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}
